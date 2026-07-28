package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aige/aimodel"
	"github.com/aige/api/internal/middleware"
	"github.com/aige/audit"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AIModelHandler struct {
	ModelRepo aimodel.Repository
	AuditRepo audit.Repository
	Pool      *pgxpool.Pool
}

func statsInterval(value string) string {
	switch value {
	case "today":
		return "1 day"
	case "7d":
		return "7 days"
	case "30d":
		return "30 days"
	default:
		return "7 days"
	}
}

func (h *AIModelHandler) logAction(r *http.Request, action, resource string, detail map[string]interface{}) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		adminID = r.RemoteAddr
	}
	h.AuditRepo.Create(r.Context(), &audit.Log{
		UserID:    adminID,
		Action:    action,
		Resource:  resource,
		Detail:    detail,
		IPAddress: r.RemoteAddr,
	})
}

func (h *AIModelHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	rangeParam := statsInterval(r.URL.Query().Get("range"))

	var totalCalls, totalTokens, failedCount int
	var totalCost float64
	err := h.Pool.QueryRow(r.Context(),
		`SELECT
			COUNT(*) as total_calls,
			COALESCE(SUM(COALESCE(a.tokens_used, 0)), 0) as total_tokens,
			COALESCE(SUM(COALESCE(a.cost, 0)), 0) as total_cost,
			COALESCE(SUM(CASE WHEN COALESCE((a.analysis->>'api_success')::boolean, true) THEN 0 ELSE 1 END), 0) as failed_count
		 FROM ai_answers a
		 JOIN ai_tasks t ON t.id = a.task_id
		 WHERE a.created_at > NOW() - $1::interval`,
		rangeParam,
	).Scan(&totalCalls, &totalTokens, &totalCost, &failedCount)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to query stats",
		})
		return
	}

	successRate := 0.0
	failureRate := 0.0
	if totalCalls > 0 {
		failureRate = float64(failedCount) / float64(totalCalls) * 100
		successRate = 100.0 - failureRate
	}

	rows, err := h.Pool.Query(r.Context(),
		`SELECT
			COALESCE(am.name, a.model, 'unknown') as model_name,
			COUNT(*) as calls,
			COALESCE(SUM(COALESCE(a.tokens_used, 0)), 0) as tokens,
			COALESCE(SUM(COALESCE(a.cost, 0)), 0) as cost,
			COALESCE(SUM(CASE WHEN COALESCE((a.analysis->>'api_success')::boolean, true) THEN 1 ELSE 0 END) * 100.0 / NULLIF(COUNT(*), 0), 0) as success_rate,
			COALESCE(AVG(NULLIF((a.analysis->>'latency_ms')::numeric, 0)), 0) as avg_latency_ms
		 FROM ai_answers a
		 JOIN ai_tasks t ON t.id = a.task_id
			LEFT JOIN ai_models am ON am.model = a.model OR am.name = a.model
		 WHERE a.created_at > NOW() - $1::interval
		 GROUP BY COALESCE(am.name, a.model, 'unknown')
		 ORDER BY calls DESC`,
		rangeParam,
	)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to query model breakdown",
		})
		return
	}
	defer rows.Close()

	type modelStat struct {
		Model       string  `json:"model"`
		Calls       int     `json:"calls"`
		Tokens      int     `json:"tokens"`
		Cost        float64 `json:"cost"`
		SuccessRate float64 `json:"success_rate"`
		AvgLatency  float64 `json:"avg_latency_ms"`
	}
	var breakdown []modelStat
	for rows.Next() {
		var s modelStat
		if err := rows.Scan(&s.Model, &s.Calls, &s.Tokens, &s.Cost, &s.SuccessRate, &s.AvgLatency); err != nil {
			continue
		}
		breakdown = append(breakdown, s)
	}
	if err := rows.Err(); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to read model breakdown",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"total_calls":  totalCalls,
			"total_tokens": totalTokens,
			"total_cost":   totalCost,
			"success_rate": successRate,
			"failure_rate": failureRate,
			"breakdown":    breakdown,
		},
	})
}

func (h *AIModelHandler) ListEnabled(w http.ResponseWriter, r *http.Request) {
	models, err := h.ModelRepo.ListEnabled(r.Context())
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	for i := range models {
		models[i].APIKey = maskAPIKey(models[i].APIKey)
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": models,
	})
}

func (h *AIModelHandler) List(w http.ResponseWriter, r *http.Request) {
	models, err := h.ModelRepo.List(r.Context())
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to list models"})
		return
	}
	for _, m := range models {
		m.APIKey = maskAPIKey(m.APIKey)
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": models})
}

func (h *AIModelHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	m, err := h.ModelRepo.GetByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{"success": false, "message": "model not found"})
		return
	}
	m.APIKey = maskAPIKey(m.APIKey)
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": m})
}

func normalizeBaseURL(url string) string {
	u := strings.TrimRight(url, "/")
	if idx := strings.LastIndex(u, "/chat/completions"); idx > 0 {
		u = u[:idx]
	}
	return strings.TrimRight(u, "/")
}

func (h *AIModelHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req aimodel.CreateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "invalid request body"})
		return
	}
	req.BaseURL = normalizeBaseURL(req.BaseURL)
	if req.Name == "" || req.Provider == "" || req.Model == "" || req.BaseURL == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "name, provider, model and base_url are required"})
		return
	}
	m, err := h.ModelRepo.Create(r.Context(), &req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to create model"})
		return
	}
	h.logAction(r, "create", "ai_model", map[string]interface{}{"model_id": m.ID, "name": m.Name, "provider": m.Provider})
	m.APIKey = maskAPIKey(m.APIKey)
	jsonResponse(w, http.StatusCreated, map[string]interface{}{"success": true, "data": m})
}

func (h *AIModelHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req aimodel.UpdateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "invalid request body"})
		return
	}
	if req.BaseURL != nil {
		n := normalizeBaseURL(*req.BaseURL)
		req.BaseURL = &n
	}
	m, err := h.ModelRepo.Update(r.Context(), id, &req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to update model"})
		return
	}
	h.logAction(r, "update", "ai_model", map[string]interface{}{"model_id": m.ID, "name": m.Name})
	m.APIKey = maskAPIKey(m.APIKey)
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": m})
}

func (h *AIModelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	m, getErr := h.ModelRepo.GetByID(r.Context(), id)
	if err := h.ModelRepo.Delete(r.Context(), id); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	name := id
	if getErr == nil && m != nil {
		name = m.Name
	}
	h.logAction(r, "delete", "ai_model", map[string]interface{}{"model_id": id, "name": name})
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "message": "model deleted"})
}

func (h *AIModelHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	m, err := h.ModelRepo.GetByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{"success": false, "message": "model not found"})
		return
	}

	var body struct {
		APIKey string `json:"api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err == nil && body.APIKey != "" {
		m.APIKey = body.APIKey
	}

	status, testMsg := testProviderConnection(m)
	if err := h.ModelRepo.UpdateTestStatus(r.Context(), id, status); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to update test status"})
		return
	}

	success := status == "success"
	msg := testMsg
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{"success": success, "status": status, "message": msg},
	})
}

func (h *AIModelHandler) Discover(w http.ResponseWriter, r *http.Request) {
	var req aimodel.DiscoverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "invalid request body"})
		return
	}
	if req.BaseURL == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "base_url is required"})
		return
	}

	apiKey := req.APIKey

	// If model_id is provided and api_key is empty/masked, look up the stored key
	if apiKey == "" || strings.Contains(apiKey, "****") {
		if req.ModelID != "" {
			m, err := h.ModelRepo.GetByID(r.Context(), req.ModelID)
			if err == nil && m.APIKey != "" {
				apiKey = m.APIKey
			}
		}
	}

	if apiKey == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "api_key is required"})
		return
	}

	models, err := fetchModels(req.BaseURL, apiKey)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": fmt.Sprintf("failed to discover models: %v", err)})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": models})
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

type openAIModelListResponse struct {
	Data []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

func modelsEndpoint(baseURL string) string {
	return strings.TrimRight(baseURL, "/") + "/models"
}

func fetchModels(baseURL, apiKey string) ([]aimodel.DiscoveredModel, error) {
	url := modelsEndpoint(normalizeBaseURL(baseURL))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error.Message != "" {
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, apiErr.Error.Message)
		}
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var list openAIModelListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	result := make([]aimodel.DiscoveredModel, 0, len(list.Data))
	for _, d := range list.Data {
		result = append(result, aimodel.DiscoveredModel{ID: d.ID, Object: d.Object, OwnedBy: d.OwnedBy})
	}
	return result, nil
}

func completionEndpoint(baseURL string) string {
	u := strings.TrimRight(baseURL, "/")
	if strings.HasSuffix(u, "/chat/completions") {
		return u
	}
	return u + "/chat/completions"
}

func testProviderConnection(m *aimodel.AIModel) (string, string) {
	if m.APIKey == "" {
		return "failed", "API Key 为空，请填写 API Key"
	}
	url := completionEndpoint(m.BaseURL)

	payload, _ := json.Marshal(map[string]interface{}{
		"model":      m.Model,
		"messages":   []map[string]string{{"role": "user", "content": "hi"}},
		"max_tokens": 1,
	})

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return "failed", "构造请求失败"
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "failed", fmt.Sprintf("无法连接到服务器：%v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiErr struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	json.Unmarshal(body, &apiErr)

	if resp.StatusCode == http.StatusOK {
		return "success", "连接成功，API 响应正常"
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return "failed", "API Key 无效，请检查 API Key"
	case http.StatusForbidden:
		return "failed", "API Key 无权限访问该接口"
	case http.StatusNotFound:
		return "success", "服务器连接正常（当前模型在此 API 中不存在）"
	case http.StatusBadRequest:
		return "success", "服务器连接正常（请求参数需要调整）"
	}

	if apiErr.Error.Type == "insufficient_quota" {
		return "success", "连接成功（API 额度不足）"
	}
	if apiErr.Error.Type == "rate_limit_exceeded" {
		return "success", "连接成功（API 频率受限，稍后重试）"
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		msg := "服务器已连接"
		if apiErr.Error.Message != "" {
			msg = msg + "，响应：" + apiErr.Error.Message
		}
		return "success", msg
	}

	return "failed", fmt.Sprintf("服务器返回异常状态码 %d", resp.StatusCode)
}

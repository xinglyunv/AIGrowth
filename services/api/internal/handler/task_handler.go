package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aige/ai-engine/analyzer"
	"github.com/aige/ai-engine/generator"
	"github.com/aige/ai-engine/models"
	"github.com/aige/ai-engine/scorer"
	"github.com/aige/aimodel"
	"github.com/aige/api/internal/middleware"
	"github.com/aige/competitor"
	"github.com/aige/project"
	"github.com/aige/report"
	"github.com/aige/task"
	"github.com/aige/user"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	TaskRepo    task.Repository
	ProjectRepo project.Repository
	ModelRepo   aimodel.Repository
	AIProvider  aimodel.AIProvider
	UserRepo    user.Repository
	CompRepo    competitor.Repository
	ReportRepo  report.Repository
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req task.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	req.ProjectID = strings.TrimSpace(req.ProjectID)
	if req.ProjectID == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "project_id is required",
		})
		return
	}

	p, err := h.ProjectRepo.FindByID(r.Context(), req.ProjectID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if p == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "project not found",
		})
		return
	}
	if p.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	models := task.NormalizeModels(req.Model)
	if len(models) == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "请至少选择一个 AI 模型",
		})
		return
	}
	requiredCredits := task.CreditCost(models)
	u, err := h.UserRepo.DeductCredits(r.Context(), userID, requiredCredits)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if u == nil {
		jsonResponse(w, http.StatusPaymentRequired, map[string]interface{}{
			"success": false, "message": "积分不足，请先充值",
		})
		return
	}

	t, err := h.TaskRepo.Create(r.Context(), req, userID)
	if err != nil {
		if _, refundErr := h.UserRepo.AddCredits(r.Context(), userID, requiredCredits); refundErr != nil {
			log.Printf("ERROR refunding credits after task creation failure: %v", refundErr)
		}
		log.Printf("ERROR creating task: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    t,
	})
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	projectID := r.URL.Query().Get("project_id")

	var tasks []task.AITask
	var total int
	var err error

	if projectID != "" {
		p, err := h.ProjectRepo.FindByID(r.Context(), projectID)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		if p == nil || p.UserID != userID {
			jsonResponse(w, http.StatusForbidden, map[string]interface{}{
				"success": false, "message": "access denied",
			})
			return
		}
		tasks, total, err = h.TaskRepo.ListByProject(r.Context(), projectID, offset, limit)
	} else {
		tasks, total, err = h.TaskRepo.ListByUser(r.Context(), userID, offset, limit)
	}
	if err != nil {
		log.Printf("ERROR listing tasks: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tasks,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")
	t, err := h.TaskRepo.FindByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if t == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "task not found",
		})
		return
	}
	if t.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    t,
	})
}

func (h *TaskHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	t, err := h.TaskRepo.FindByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if t == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "task not found",
		})
		return
	}
	if t.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	report, err := h.TaskRepo.GetReport(r.Context(), id)
	if err != nil {
		log.Printf("ERROR getting report: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if report == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "report not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
	})
}

func (h *TaskHandler) GetComparisonReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	t, err := h.TaskRepo.FindByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if t == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "task not found",
		})
		return
	}
	if t.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	report, err := h.TaskRepo.GetComparisonReport(r.Context(), id)
	if err != nil {
		log.Printf("ERROR getting comparison report: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if report == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "report not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
	})
}

func (h *TaskHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{"success": false, "message": "unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	t, err := h.TaskRepo.FindByID(r.Context(), id)
	if err != nil || t == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{"success": false, "message": "task not found"})
		return
	}
	if t.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{"success": false, "message": "access denied"})
		return
	}
	if t.Status != "pending" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "task is not in pending status"})
		return
	}

	if err := h.TaskRepo.UpdateStatus(r.Context(), id, "running", 0, ""); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}

	p, err := h.ProjectRepo.FindByID(r.Context(), t.ProjectID)
	if err != nil || p == nil {
		h.TaskRepo.UpdateStatus(r.Context(), id, "failed", 0, "project not found")
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "project not found"})
		return
	}

	enabledModels, err := h.ModelRepo.ListEnabled(r.Context())
	if err != nil || len(enabledModels) == 0 {
		h.TaskRepo.UpdateStatus(r.Context(), id, "failed", 0, "no enabled AI models")
		jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": map[string]interface{}{
			"id": id, "status": "failed", "error_message": "no enabled AI models configured",
		}})
		return
	}

	selectedModels := enabledModels
	if t.Model != "" {
		taskModelNames := task.NormalizeModels(t.Model)
		modelNameSet := make(map[string]bool)
		for _, tn := range taskModelNames {
			modelNameSet[strings.TrimSpace(tn)] = true
		}
		var filtered []*aimodel.AIModel
		for _, m := range enabledModels {
			if modelNameSet[m.Model] || modelNameSet[m.Name] {
				filtered = append(filtered, m)
			}
		}
		if len(filtered) > 0 {
			selectedModels = filtered
		}
	}
	if len(selectedModels) == 0 {
		h.TaskRepo.UpdateStatus(r.Context(), id, "failed", 0, "没有匹配的可用模型")
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "没有匹配的可用模型"})
		return
	}

	// Launch background processing and return immediately
	go h.processTask(id, t.ProjectID, p, selectedModels, userID)

	jsonResponse(w, http.StatusAccepted, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":     id,
			"status": "running",
		},
	})
}

func (h *TaskHandler) processTask(id, projectID string, p *project.BrandProject, selectedModels []*aimodel.AIModel, userID string) {
	ctx := context.Background()

	confirmedCompetitors := make([]*competitor.Competitor, 0)
	var err error
	if h.CompRepo != nil {
		confirmedCompetitors, err = h.CompRepo.ListByProjectID(ctx, p.ID)
		if err != nil {
			log.Printf("WARN listing competitors for task %s: %v", id, err)
			confirmedCompetitors = nil
		}
	}
	competitorNames := make([]string, 0, len(confirmedCompetitors))
	for _, comp := range confirmedCompetitors {
		if comp != nil && strings.TrimSpace(comp.Name) != "" && !strings.EqualFold(comp.Name, p.Name) {
			competitorNames = append(competitorNames, strings.TrimSpace(comp.Name))
		}
	}

	questions := generator.GenerateQuestions(p.Industry, p.Name, 5)
	if err := h.TaskRepo.SetQuestionsCount(ctx, id, len(questions)*len(selectedModels)); err != nil {
		h.TaskRepo.UpdateStatus(ctx, id, "failed", 0, "无法初始化检测进度")
		log.Printf("ERROR SetQuestionsCount for task %s: %v", id, err)
		return
	}
	completedCount := 0
	type competitorAggregate struct {
		count      int
		bestRank   int
		advantages []string
	}
	competitorAggregates := make(map[string]*competitorAggregate)

	for _, q := range questions {
		for _, m := range selectedModels {
			startedAt := time.Now()
			formattedPrompt := buildCompetitivePrompt(p.Name, competitorNames, q)

			msg := []aimodel.ChatMessage{
				{Role: "system", Content: "你是一个专业的品牌分析工具。请严格按照要求的 JSON 格式回复。"},
				{Role: "user", Content: formattedPrompt},
			}

			respContent, apiErr := h.AIProvider.Chat(ctx, m, msg)
			if apiErr != nil {
				log.Printf("ERROR calling model %s: %v", m.Name, apiErr)
				respContent = fmt.Sprintf(`{"brand_mentioned": false, "answer": "API调用失败: %v", "sentiment": "neutral", "rank_position": 0}`, apiErr)
			}

			var parsed competitivePromptResponse
			cleaned := strings.TrimSpace(respContent)
			cleaned = strings.TrimPrefix(cleaned, "```json")
			cleaned = strings.TrimPrefix(cleaned, "```")
			cleaned = strings.TrimSuffix(cleaned, "```")
			cleaned = strings.TrimSpace(cleaned)

			answerText := respContent
			brandMentioned := false
			sentiment := "neutral"
			rankPos := 0

			var entities []models.EntityAnalysis
			if err := json.Unmarshal([]byte(cleaned), &parsed); err == nil {
				answerText = parsed.Answer
				localTargets := make([]analyzer.EntityTarget, 0, len(competitorNames))
				for _, name := range competitorNames {
					localTargets = append(localTargets, analyzer.EntityTarget{Name: name})
				}
				localEntities := analyzer.AnalyzeEntities(answerText, analyzer.EntityTarget{Name: p.Name}, localTargets)
				entities = mergeEntityAnalysis(parsed.Entities, localEntities)
				brandMentioned = localEntities[0].Mentioned
				sentiment = analyzer.AnalyzeSentiment(answerText, p.Name)
				if sentiment == "neutral" && parsed.Sentiment != "" {
					sentiment = parsed.Sentiment
				}
				rankPos = parsed.RankPosition
				for _, entity := range entities {
					if entity.Role == "target" && entity.RankPosition > 0 {
						rankPos = entity.RankPosition
					}
				}
			} else {
				brandMentioned = analyzer.DetectBrand(respContent, p.Name).Mentioned
				sentiment = analyzer.AnalyzeSentiment(respContent, p.Name)
				localTargets := make([]analyzer.EntityTarget, 0, len(competitorNames))
				for _, name := range competitorNames {
					localTargets = append(localTargets, analyzer.EntityTarget{Name: name})
				}
				entities = mergeEntityAnalysis(nil, analyzer.AnalyzeEntities(respContent, analyzer.EntityTarget{Name: p.Name}, localTargets))
			}

			rankPosVal := rankPos
			competitors := analyzer.DetectCompetitors(answerText, p.Name)
			if len(entities) > 0 {
				competitors = nil
				for _, entity := range entities {
					if entity.Role == "competitor" && entity.Mentioned {
						competitors = append(competitors, entity.Name)
					}
				}
			}

			answerRecord := &task.AIAnswer{
				TaskID:         id,
				Question:       q,
				Answer:         answerText,
				Model:          m.Model,
				BrandMentioned: brandMentioned,
				Sentiment:      sentiment,
				RankPosition:   &rankPosVal,
				Analysis: map[string]interface{}{
					"model_name":          m.Name,
					"model_provider":      m.Provider,
					"api_success":         apiErr == nil,
					"competitor_mentions": len(competitors),
					"latency_ms":          time.Since(startedAt).Milliseconds(),
					"entities":            entities,
					"comparison":          parsed.Comparison,
				},
			}

			if err := h.TaskRepo.SaveAnswer(ctx, answerRecord); err != nil {
				log.Printf("ERROR saving answer: %v", err)
				h.TaskRepo.UpdateStatus(ctx, id, "failed", completedCount, err.Error())
				return
			}
			for _, entity := range entities {
				if entity.Role != "competitor" || !entity.Mentioned {
					continue
				}
				aggregate := competitorAggregates[entity.Name]
				if aggregate == nil {
					aggregate = &competitorAggregate{}
					competitorAggregates[entity.Name] = aggregate
				}
				aggregate.count += entity.MentionCount
				if entity.RankPosition > 0 && (aggregate.bestRank == 0 || entity.RankPosition < aggregate.bestRank) {
					aggregate.bestRank = entity.RankPosition
				}
				aggregate.advantages = append(aggregate.advantages, entity.Advantages...)
			}
			completedCount++
		}
	}

	allAnswers, _ := h.TaskRepo.GetAnswers(ctx, id)
	brandMentions := 0
	topRankCount := 0
	var totalSentiment float64
	sentimentCount := 0
	for _, a := range allAnswers {
		if a.BrandMentioned {
			brandMentions++
		}
		if a.RankPosition != nil && *a.RankPosition <= 3 {
			topRankCount++
		}
		switch a.Sentiment {
		case "positive":
			totalSentiment += 0.8
			sentimentCount++
		case "negative":
			totalSentiment += 0.2
			sentimentCount++
		case "neutral":
			totalSentiment += 0.5
			sentimentCount++
		}
	}
	avgSentiment := 0.5
	if sentimentCount > 0 {
		avgSentiment = totalSentiment / float64(sentimentCount)
	}
	mentionRate := 0.0
	if len(allAnswers) > 0 {
		mentionRate = float64(brandMentions) / float64(len(allAnswers))
	}

	analysisResult := &models.AnalysisResult{
		TaskID:           id,
		BrandName:        p.Name,
		BrandMentions:    brandMentions,
		TotalQuestions:   len(allAnswers),
		MentionRate:      mentionRate,
		ModelCoverage:    len(selectedModels),
		TopRankCount:     topRankCount,
		AverageSentiment: avgSentiment,
	}
	for name, aggregate := range competitorAggregates {
		analysisResult.Competitors = append(analysisResult.Competitors, models.CompetitorMention{
			Name: name, Count: aggregate.count, TopRank: aggregate.bestRank,
			Advantages: strings.Join(uniqueStrings(aggregate.advantages), "；"),
		})
		if h.CompRepo != nil {
			if err := h.CompRepo.UpsertSummary(ctx, p.ID, name, aggregate.count, aggregate.bestRank, strings.Join(uniqueStrings(aggregate.advantages), "；"), map[string]interface{}{"task_id": id}); err != nil {
				log.Printf("WARN persisting competitor %s for task %s: %v", name, id, err)
			}
		}
	}
	visibilityScore := scorer.CalculateVisibilityScore(analysisResult)
	competitiveMetrics := calculateCompetitiveMetrics(allAnswers, p.Name)
	analysisResult.CompetitiveScore = scorer.CalculateCompetitiveScore(competitiveMetrics)
	if h.ReportRepo != nil {
		content := map[string]interface{}{
			"task_id":           id,
			"target_brand":      p.Name,
			"visibility_score":  visibilityScore,
			"competitive_score": analysisResult.CompetitiveScore,
			"brand_mentions":    brandMentions,
			"total_answers":     len(allAnswers),
			"competitors":       analysisResult.Competitors,
			"partial_failures":  completedCount < len(questions)*len(selectedModels),
		}
		if err := h.ReportRepo.Create(ctx, &report.Report{
			ProjectID: p.ID, TaskID: id, UserID: userID, Title: p.Name + " 竞争分析报告",
			Type: "competitive", VisibilityScore: visibilityScore, Content: content,
			Summary: fmt.Sprintf("目标品牌提及率 %.1f%%，识别到 %d 个竞争品牌", mentionRate*100, len(analysisResult.Competitors)), Status: "completed",
		}); err != nil {
			log.Printf("WARN creating report for task %s: %v", id, err)
		}
	}

	status := "completed"

	if err := h.TaskRepo.UpdateStatus(ctx, id, status, completedCount, ""); err != nil {
		log.Printf("ERROR updating task status: %v", err)
	}

	log.Printf("Task %s completed: visibility=%v competitive=%.1f", id, visibilityScore, analysisResult.CompetitiveScore)
}

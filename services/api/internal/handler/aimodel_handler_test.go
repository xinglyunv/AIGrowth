package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchModelsIncludesProviderErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"invalid api key"}}`))
	}))
	defer server.Close()

	_, err := fetchModels(server.URL+"/v1", "bad-key")
	if err == nil || !strings.Contains(err.Error(), "invalid api key") {
		t.Fatalf("expected provider error message, got %v", err)
	}
}

func TestFetchModelsReturnsDiscoveredModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"data":[{"id":"model-a","object":"model","owned_by":"test"}]}`)
	}))
	defer server.Close()

	models, err := fetchModels(server.URL+"/v1/chat/completions", "test-key")
	if err != nil {
		t.Fatalf("fetchModels returned error: %v", err)
	}
	if len(models) != 1 || models[0].ID != "model-a" {
		t.Fatalf("unexpected models: %+v", models)
	}
}

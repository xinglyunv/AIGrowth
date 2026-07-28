package handler

import (
	"testing"

	"github.com/aige/ai-engine/analyzer"
	"github.com/aige/ai-engine/models"
)

func TestMergeEntityAnalysisTrustsLocalEvidence(t *testing.T) {
	modelEntities := []models.EntityAnalysis{{Name: "豆包", Role: "target", Mentioned: true}}
	local := analyzer.AnalyzeEntities("豆包适合中文写作。", analyzer.EntityTarget{Name: "豆包"}, nil)

	merged := mergeEntityAnalysis(modelEntities, local)
	if len(merged) != 1 || merged[0].Evidence == "" || merged[0].Validation != "verified" {
		t.Fatalf("expected local evidence, got %+v", merged)
	}
}

func TestCompetitivePromptIncludesEveryConfirmedCompetitor(t *testing.T) {
	prompt := buildCompetitivePrompt("豆包", []string{"ChatGPT", "Claude"}, "中文 AI 助手有哪些？")
	for _, expected := range []string{"豆包", "ChatGPT", "Claude", "中文 AI 助手有哪些？"} {
		if !contains(prompt, expected) {
			t.Fatalf("prompt does not include %q: %s", expected, prompt)
		}
	}
}

func contains(value, expected string) bool {
	for i := 0; i+len(expected) <= len(value); i++ {
		if value[i:i+len(expected)] == expected {
			return true
		}
	}
	return false
}

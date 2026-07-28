package scorer

import (
	"testing"

	"github.com/aige/ai-engine/models"
)

func TestCalculateCompetitiveScoreUsesConfiguredWeights(t *testing.T) {
	score := CalculateCompetitiveScore(models.CompetitiveMetrics{
		MentionRate:          1,
		RankScore:            100,
		ReasonQuality:        1,
		Completeness:         1,
		InformationAdvantage: 1,
	})

	if score != 100 {
		t.Fatalf("expected full score, got %d", score)
	}
}

func TestCalculateCompetitiveScoreClampsIncompleteMetrics(t *testing.T) {
	score := CalculateCompetitiveScore(models.CompetitiveMetrics{})
	if score != 0 {
		t.Fatalf("expected empty metrics to score 0, got %d", score)
	}
}

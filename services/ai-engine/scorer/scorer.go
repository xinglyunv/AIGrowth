package scorer

import (
	"math"

	"github.com/aige/ai-engine/models"
)

func CalculateVisibilityScore(result *models.AnalysisResult) int {
	if result == nil || result.TotalQuestions == 0 {
		return 0
	}

	mentionRateScore := int(math.Round(result.MentionRate * 30))
	if mentionRateScore > 30 {
		mentionRateScore = 30
	}

	topRankScore := 0
	if result.TopRankCount > 0 {
		ratio := float64(result.TopRankCount) / float64(result.TotalQuestions)
		topRankScore = int(math.Round(ratio * 25))
	}
	if topRankScore > 25 {
		topRankScore = 25
	}

	modelCoverageScore := result.ModelCoverage * 4
	if modelCoverageScore > 20 {
		modelCoverageScore = 20
	}

	sentimentScore := 0
	if result.AverageSentiment > 0.3 {
		sentimentScore = int(math.Round(result.AverageSentiment * 15))
	}
	if sentimentScore > 15 {
		sentimentScore = 15
	}

	competitorGapScore := 0
	if len(result.Competitors) == 0 {
		competitorGapScore = 10
	} else {
		totalMentions := 0
		for _, c := range result.Competitors {
			totalMentions += c.Count
		}
		if totalMentions == 0 {
			competitorGapScore = 10
		} else if totalMentions < result.BrandMentions {
			competitorGapScore = 8
		} else if totalMentions == result.BrandMentions {
			competitorGapScore = 5
		} else {
			competitorGapScore = 2
		}
	}
	if competitorGapScore > 10 {
		competitorGapScore = 10
	}

	score := mentionRateScore + topRankScore + modelCoverageScore + sentimentScore + competitorGapScore
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

func CalculateCompetitiveScore(metrics models.CompetitiveMetrics) int {
	clamp := func(value float64) float64 {
		if value < 0 {
			return 0
		}
		if value > 1 {
			return 1
		}
		return value
	}
	score := clamp(metrics.MentionRate)*30 +
		clamp(metrics.RankScore)*25 +
		clamp(metrics.ReasonQuality)*20 +
		clamp(metrics.Completeness)*15 +
		clamp(metrics.InformationAdvantage)*10
	return int(math.Round(score))
}

func CalculateTrend(current, previous int) string {
	if current > previous {
		return "up"
	}
	if current < previous {
		return "down"
	}
	return "stable"
}

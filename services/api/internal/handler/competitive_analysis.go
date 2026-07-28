package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aige/ai-engine/analyzer"
	"github.com/aige/ai-engine/models"
	"github.com/aige/task"
)

type competitivePromptResponse struct {
	Answer       string                  `json:"answer"`
	Entities     []models.EntityAnalysis `json:"entities"`
	Comparison   map[string]interface{}  `json:"comparison"`
	BrandLegacy  bool                    `json:"brand_mentioned"`
	Sentiment    string                  `json:"sentiment"`
	RankPosition int                     `json:"rank_position"`
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func buildCompetitivePrompt(target string, competitors []string, question string) string {
	competitorJSON, _ := json.Marshal(competitors)
	return fmt.Sprintf(`你正在回答真实用户问题，请先自然回答，再输出结构化分析。

用户问题：%s
目标品牌：%s
竞争品牌：%s

要求：不刻意推广品牌，依据你的知识自然回答。结构化分析必须只记录回答中实际出现的品牌；每个品牌提供角色、提及次数、推荐顺序、证据片段、优势和劣势。

请严格输出 JSON：
{
  "answer": "自然回答",
  "entities": [{"name":"品牌名称","role":"target|competitor","mentioned":true,"mention_count":1,"rank_position":1,"sentiment":"positive|neutral|negative","evidence":"原文证据","advantages":[],"disadvantages":[]}],
  "comparison": {"winner":"品牌名称","dimensions":[]}
}`, question, target, string(competitorJSON))
}

func mergeEntityAnalysis(modelEntities []models.EntityAnalysis, local []analyzer.EntityResult) []models.EntityAnalysis {
	byName := make(map[string]models.EntityAnalysis, len(modelEntities))
	for _, entity := range modelEntities {
		byName[entity.Name] = entity
	}
	merged := make([]models.EntityAnalysis, 0, len(local))
	for _, result := range local {
		entity := byName[result.Name]
		entity.Name = result.Name
		entity.Role = result.Role
		entity.Mentioned = result.Mentioned
		entity.MentionCount = result.MentionCount
		entity.Position = result.Position
		entity.Evidence = result.Evidence
		entity.Validation = result.Validation
		merged = append(merged, entity)
	}
	return merged
}

func calculateCompetitiveMetrics(answers []task.AIAnswer, target string) models.CompetitiveMetrics {
	if len(answers) == 0 {
		return models.CompetitiveMetrics{}
	}
	mentioned := 0
	rankTotal := 0.0
	rankCount := 0
	evidenceCount := 0
	completeCount := 0
	for _, answer := range answers {
		var analysis struct {
			Entities []models.EntityAnalysis `json:"entities"`
		}
		encoded, _ := json.Marshal(answer.Analysis)
		_ = json.Unmarshal(encoded, &analysis)
		for _, entity := range analysis.Entities {
			if entity.Role != "target" || entity.Name != target {
				continue
			}
			if !entity.Mentioned {
				continue
			}
			mentioned++
			if entity.RankPosition > 0 {
				rankCount++
				switch {
				case entity.RankPosition == 1:
					rankTotal += 1
				case entity.RankPosition <= 3:
					rankTotal += 0.8
				default:
					rankTotal += 0.4
				}
			}
			if entity.Evidence != "" {
				evidenceCount++
			}
			if len(entity.Advantages) > 0 || len(entity.Disadvantages) > 0 {
				completeCount++
			}
		}
	}
	metrics := models.CompetitiveMetrics{MentionRate: float64(mentioned) / float64(len(answers))}
	if rankCount > 0 {
		metrics.RankScore = rankTotal / float64(rankCount)
	}
	if mentioned > 0 {
		metrics.ReasonQuality = float64(evidenceCount) / float64(mentioned)
		metrics.Completeness = float64(completeCount) / float64(mentioned)
	}
	metrics.InformationAdvantage = metrics.Completeness
	return metrics
}

package analyzer

import (
	"regexp"
	"strings"

	"github.com/aige/ai-engine/models"
)

type EntityTarget struct {
	Name    string
	Aliases []string
}

type EntityResult struct {
	Name         string
	Role         string
	Mentioned    bool
	MentionCount int
	Position     int
	Evidence     string
	Validation   string
}

func AnalyzeEntities(text string, target EntityTarget, competitors []EntityTarget) []EntityResult {
	targets := make([]EntityTarget, 0, len(competitors)+1)
	targets = append(targets, target)
	targets = append(targets, competitors...)
	results := make([]EntityResult, 0, len(targets))
	for index, entity := range targets {
		role := "competitor"
		if index == 0 {
			role = "target"
		}
		result := EntityResult{Name: entity.Name, Role: role, Position: -1, Validation: "absent"}
		for _, candidate := range append([]string{entity.Name}, entity.Aliases...) {
			positions := matchPositions(text, candidate)
			if len(positions) == 0 {
				continue
			}
			if !result.Mentioned || positions[0] < result.Position {
				result.Position = positions[0]
				result.Evidence = extractContext(text, positions[0], len(candidate))
			}
			result.Mentioned = true
			result.MentionCount += len(positions)
			result.Validation = "verified"
		}
		results = append(results, result)
	}
	return results
}

func ToModelEntities(results []EntityResult) []models.EntityAnalysis {
	entities := make([]models.EntityAnalysis, 0, len(results))
	for _, result := range results {
		entities = append(entities, models.EntityAnalysis{
			Name: result.Name, Role: result.Role, Mentioned: result.Mentioned,
			MentionCount: result.MentionCount, Position: result.Position,
			Evidence: result.Evidence, Validation: result.Validation,
		})
	}
	return entities
}

func matchPositions(text, candidate string) []int {
	if text == "" || candidate == "" {
		return nil
	}
	if isASCIIWord(candidate) && len([]rune(candidate)) <= 3 {
		pattern := `(?i)(^|[^a-z0-9])` + regexp.QuoteMeta(candidate) + `([^a-z0-9]|$)`
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringIndex(text, -1)
		positions := make([]int, 0, len(matches))
		for _, match := range matches {
			position := strings.Index(strings.ToLower(text[match[0]:match[1]]), strings.ToLower(candidate))
			positions = append(positions, match[0]+position)
		}
		return positions
	}
	lowerText := strings.ToLower(text)
	lowerCandidate := strings.ToLower(candidate)
	positions := make([]int, 0)
	for offset := 0; offset < len(lowerText); {
		position := strings.Index(lowerText[offset:], lowerCandidate)
		if position < 0 {
			break
		}
		position += offset
		positions = append(positions, position)
		offset = position + len(lowerCandidate)
	}
	return positions
}

func isASCIIWord(value string) bool {
	for _, char := range value {
		if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || char == ' ') {
			return false
		}
	}
	return true
}

type BrandMentionResult struct {
	Mentioned bool
	Position  int
	Context   string
}

type BrandDetector struct {
	BrandName  string
	Synonyms   []string
	ExactMatch bool
}

func DetectBrand(text string, brandName string) *BrandMentionResult {
	if text == "" || brandName == "" {
		return &BrandMentionResult{}
	}

	lowerText := strings.ToLower(text)
	lowerBrand := strings.ToLower(brandName)

	pos := strings.Index(lowerText, lowerBrand)
	if pos >= 0 {
		start := pos - 30
		if start < 0 {
			start = 0
		}
		end := pos + len(brandName) + 30
		if end > len(text) {
			end = len(text)
		}
		context := text[start:end]
		return &BrandMentionResult{
			Mentioned: true,
			Position:  pos,
			Context:   context,
		}
	}

	words := strings.Fields(lowerBrand)
	for _, word := range words {
		if len(word) <= 2 {
			continue
		}
		pos := strings.Index(lowerText, word)
		if pos >= 0 {
			context := extractContext(text, pos, len(word))
			return &BrandMentionResult{
				Mentioned: true,
				Position:  pos,
				Context:   context,
			}
		}
	}

	return &BrandMentionResult{}
}

func DetectCompetitors(text string, brandName string) []string {
	if text == "" {
		return nil
	}

	knownBrands := []string{
		"OpenAI", "Google", "Microsoft", "Meta", "Amazon", "Apple",
		"百度", "阿里", "腾讯", "字节跳动", "华为", "科大讯飞",
		"商汤", "旷视", "寒武纪", "第四范式",
	}

	lowerText := strings.ToLower(text)
	lowerBrand := strings.ToLower(brandName)
	var competitors []string

	seen := make(map[string]bool)
	for _, cb := range knownBrands {
		lowerCB := strings.ToLower(cb)
		if lowerCB == lowerBrand {
			continue
		}
		if seen[lowerCB] {
			continue
		}
		if strings.Contains(lowerText, lowerCB) {
			competitors = append(competitors, cb)
			seen[lowerCB] = true
		}
	}

	return competitors
}

func AnalyzeSentiment(text string, brandName string) string {
	if text == "" {
		return "neutral"
	}

	lowerText := strings.ToLower(text)

	positiveWords := []string{
		"优秀", "出色", "领先", "创新", "好评", "满意", "推荐",
		"excellent", "great", "outstanding", "leading", "innovative",
		"positive", "recommend", "best", "top", "amazing",
		"强大", "可靠", "高效", "便捷", "优质", "专业",
		"robust", "reliable", "efficient", "convenient", "professional",
	}

	negativeWords := []string{
		"差劲", "糟糕", "落后", "不满", "投诉", "差评", "失望",
		"terrible", "poor", "bad", "worst", "awful", "disappointing",
		"落后", "不稳定", "昂贵", "复杂", "难用",
		"unstable", "expensive", "complicated", "useless", "horrible",
	}

	positiveCount := 0
	for _, w := range positiveWords {
		if strings.Contains(lowerText, w) {
			positiveCount++
		}
	}

	negativeCount := 0
	for _, w := range negativeWords {
		if strings.Contains(lowerText, w) {
			negativeCount++
		}
	}

	if positiveCount > negativeCount {
		return "positive"
	}
	if negativeCount > positiveCount {
		return "negative"
	}
	return "neutral"
}

func extractContext(text string, pos int, length int) string {
	start := pos - 30
	if start < 0 {
		start = 0
	}
	end := pos + length + 30
	if end > len(text) {
		end = len(text)
	}
	return text[start:end]
}

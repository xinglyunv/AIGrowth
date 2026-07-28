package models

type Question struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type Answer struct {
	Question       string                 `json:"question"`
	Answer         string                 `json:"answer"`
	Model          string                 `json:"model"`
	BrandMentioned bool                   `json:"brand_mentioned"`
	Sentiment      string                 `json:"sentiment"`
	RankPosition   int                    `json:"rank_position"`
	Analysis       map[string]interface{} `json:"analysis"`
}

type AnalysisResult struct {
	TaskID           string              `json:"task_id"`
	BrandName        string              `json:"brand_name"`
	BrandMentions    int                 `json:"brand_mentions"`
	TotalQuestions   int                 `json:"total_questions"`
	MentionRate      float64             `json:"mention_rate"`
	ModelCoverage    int                 `json:"model_coverage"`
	TopRankCount     int                 `json:"top_rank_count"`
	AverageSentiment float64             `json:"average_sentiment"`
	Competitors      []CompetitorMention `json:"competitors"`
	VisibilityScore  int                 `json:"visibility_score"`
	CompetitiveScore int                 `json:"competitive_score"`
}

type CompetitiveMetrics struct {
	MentionRate          float64
	RankScore            float64
	ReasonQuality        float64
	Completeness         float64
	InformationAdvantage float64
}

type EntityAnalysis struct {
	Name          string   `json:"name"`
	Role          string   `json:"role"`
	Mentioned     bool     `json:"mentioned"`
	MentionCount  int      `json:"mention_count"`
	Position      int      `json:"position"`
	RankPosition  int      `json:"rank_position"`
	Evidence      string   `json:"evidence,omitempty"`
	Sentiment     string   `json:"sentiment,omitempty"`
	Advantages    []string `json:"advantages,omitempty"`
	Disadvantages []string `json:"disadvantages,omitempty"`
	Validation    string   `json:"validation,omitempty"`
}

type ComparisonDimension struct {
	Name        string            `json:"name"`
	Target      string            `json:"target,omitempty"`
	Competitors map[string]string `json:"competitors,omitempty"`
	Reason      string            `json:"reason,omitempty"`
}

type CompetitorMention struct {
	Name       string `json:"name"`
	Count      int    `json:"count"`
	TopRank    int    `json:"top_rank"`
	Advantages string `json:"advantages"`
}

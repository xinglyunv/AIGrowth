package prompts

var promptTemplates = map[string]string{
	"brand_analysis": `你是一个专业的品牌分析助手。请分析以下问题中是否提及了品牌"%s"。

问题：%s

请按以下 JSON 格式回复（不要加任何其他内容）：
{
  "brand_mentioned": true/false,
  "answer": "对该问题的详细回答（200字以内）",
  "sentiment": "positive/neutral/negative",
  "rank_position": 1-10的数字（1=最靠前，10=最靠后，若品牌未被提及则填0）
}

其中 brand_mentioned 为 true 当且仅当回答中明确提到了该品牌名称。sentiment 表示该回答对品牌的整体情感倾向。`,

	"question_generation": `你是一个专业的品牌可见度分析专家。请为以下品牌生成关于其市场可见度的问题。

品牌名称：%s
行业：%s
问题数量：%d

请生成关于该品牌在AI对话和搜索结果中可见度的相关问题。问题应该覆盖：
1. 品牌知名度和认知度
2. 产品和服务质量
3. 用户满意度和口碑
4. 与竞品的比较
5. 技术创新和未来发展

请按以下 JSON 格式回复（不要加任何其他内容）：
{
  "questions": ["问题1", "问题2", "问题3", ...]
}`,

	"sentiment_analysis": `你是一个专业的品牌情感分析助手。请分析以下文本对品牌"%s"的情感倾向。

文本：%s

请按以下 JSON 格式回复（不要加任何其他内容）：
{
  "sentiment": "positive/neutral/negative",
  "score": 0-100的分数（越高越正面，中性为50分）,
  "key_points": ["关键观点1", "关键观点2", ...],
  "explanation": "情感判断的简要解释"
}`,

	"competitor_analysis": `你是一个专业的竞品分析助手。请分析以下文本中提到的与品牌"%s"相关的竞争对手信息。

文本：%s

请按以下 JSON 格式回复（不要加任何其他内容）：
{
  "competitors": [
    {
      "name": "竞品名称",
      "advantages": "竞品的优势描述",
      "threat_level": "high/medium/low"
    }
  ],
  "brand_position": "该品牌相对于竞品的市场定位"
}`,

	"optimization_advice": `你是一个品牌可见度优化顾问。基于以下分析数据，为品牌"%s"提供可见度优化建议。

品牌名称：%s
行业：%s
品牌提及次数：%d
总问题数：%d
提及率：%.1f%%
覆盖模型数：%d
排名前三次数：%d
平均情感分：%.1f
可见度评分：%d/100

竞品信息：
%s

请按以下 JSON 格式回复（不要加任何其他内容）：
{
  "recommendations": [
    {
      "priority": "high/medium/low",
      "category": "内容策略/SEO优化/品牌建设/竞品应对/其他",
      "suggestion": "具体建议",
      "expected_impact": "预期效果"
    }
  ],
  "summary": "总体评估和建议摘要"
}`,
}

func GetPrompt(purpose string) string {
	template, ok := promptTemplates[purpose]
	if !ok {
		return ""
	}
	return template
}

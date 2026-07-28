package analyzer

import "testing"

func TestAnalyzeEntitiesMatchesAliasesAndPreservesEvidence(t *testing.T) {
	answer := "豆包适合中文写作，ChatGPT 的生态更丰富，豆包也有不错的移动端体验。"
	entities := AnalyzeEntities(answer, EntityTarget{Name: "豆包", Aliases: []string{"字节豆包"}}, []EntityTarget{
		{Name: "ChatGPT", Aliases: []string{"GPT"}},
	})

	if len(entities) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(entities))
	}
	if !entities[0].Mentioned || entities[0].MentionCount != 2 || entities[0].Role != "target" {
		t.Fatalf("unexpected target entity: %+v", entities[0])
	}
	if entities[0].Evidence == "" || entities[0].Position < 0 {
		t.Fatalf("expected target evidence and position: %+v", entities[0])
	}
	if !entities[1].Mentioned || entities[1].Role != "competitor" {
		t.Fatalf("unexpected competitor entity: %+v", entities[1])
	}
}

func TestAnalyzeEntitiesDoesNotMatchShortAliasInsideWord(t *testing.T) {
	entities := AnalyzeEntities("GPTX 是一个内部代号。", EntityTarget{Name: "豆包"}, []EntityTarget{{Name: "GPT", Aliases: []string{"AI"}}})

	if len(entities) != 2 || entities[1].Mentioned {
		t.Fatalf("expected no competitor mention, got %+v", entities)
	}
}

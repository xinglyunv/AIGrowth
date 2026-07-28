# Competitive Detection Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the P0 competitive detection loop with confirmed competitor snapshots, structured entity analysis, deterministic scoring, progress tracking, and persisted report data.

**Architecture:** Keep the existing Go services and JSONB answer analysis. Extend the analyzer with project-scoped brand entities, make the task handler send target and confirmed competitors in one structured prompt, validate model output against the original answer, then aggregate and persist deterministic metrics.

**Tech Stack:** Go, PostgreSQL, pgx, existing AI provider, existing task and competitor repositories.

## Global Constraints

- Preserve the existing REST API response envelope.
- Keep SQL parameterized.
- Do not introduce new dependencies.
- Keep user-facing errors in Simplified Chinese.
- Never write secrets into prompts, logs, fixtures, or source files.
- Target brand and confirmed competitors must be analyzed in the same execution unit.

### Task 1: Define Competitive Analysis Contracts

**Files:**
- Modify: `services/ai-engine/models/models.go`
- Modify: `services/competitor/model.go`
- Test: `services/ai-engine/models/models_test.go`

- [ ] Add typed entity, evidence, comparison, and competitive score structures.
- [ ] Add competitor aliases and role fields with JSON compatibility.
- [ ] Add tests for JSON round-trip and role validation.
- [ ] Run `go test ./...` in `services/ai-engine` and `services/competitor`.

### Task 2: Replace Hard-Coded Competitor Detection

**Files:**
- Modify: `services/ai-engine/analyzer/analyzer.go`
- Test: `services/ai-engine/analyzer/analyzer_test.go`

- [ ] Add `AnalyzeEntities(answer, target, competitors)` using exact names and aliases.
- [ ] Return evidence, occurrence count, first position, context, role, and validation state.
- [ ] Preserve `DetectBrand` and `DetectCompetitors` compatibility for existing callers.
- [ ] Add failing tests for aliases, target exclusion, multiple mentions, and substring false positives.
- [ ] Run analyzer tests and verify the new tests fail before implementation.

### Task 3: Persist Project Competitor Snapshots

**Files:**
- Modify: `services/competitor/repository.go`
- Modify: `services/api/internal/handler/competitor_handler.go`
- Modify: `services/api/internal/router/router.go`
- Modify: `services/database/migrations/000015_competitor_aliases.sql`
- Test: `services/api/internal/handler/competitor_handler_test.go`

- [ ] Add project-owned create, update, and delete endpoints.
- [ ] Validate name, type, project ownership, and aliases before persistence.
- [ ] Add alias persistence with a migration that keeps existing rows valid.
- [ ] Add handler tests for ownership and invalid competitor input.

### Task 4: Add Task Brand Snapshot and Progress Contracts

**Files:**
- Modify: `services/task/model.go`
- Modify: `services/task/repository.go`
- Modify: `services/api/internal/handler/task_handler.go`
- Modify: `services/database/migrations/000016_task_competitor_snapshots.sql`
- Test: `services/task/model_test.go`

- [ ] Store target brand and confirmed competitor snapshot with each task.
- [x] Set total execution units to question count multiplied by selected model count.
- [x] Add atomic pending-to-running claim semantics.
- [x] Update progress after every execution unit.
- [ ] Add tests for total units, progress boundaries, and duplicate execution prevention.

### Task 5: Implement Structured Target and Competitor Analysis

**Files:**
- Modify: `services/ai-engine/prompts/prompts.go`
- Modify: `services/api/internal/handler/task_handler.go`
- Modify: `services/ai-engine/scorer/scorer.go`
- Test: `services/ai-engine/scorer/scorer_test.go`

- [x] Send target and snapshot competitors in one structured prompt.
- [x] Parse the answer and entity list with safe JSON cleanup.
- [x] Reconcile model entities with local analyzer evidence.
- [x] Store complete entity analysis in `AIAnswer.Analysis`.
- [x] Add deterministic appearance, ranking, reason quality, completeness, and information scores.
- [ ] Add scorer tests for missing, listed, and first-ranked brands.

### Task 6: Aggregate Competitive Results and Persist Reports

**Files:**
- Modify: `services/api/internal/handler/task_handler.go`
- Modify: `services/task/repository.go`
- Modify: `services/report/model.go`
- Modify: `services/report/repository.go`
- Test: `services/api/internal/handler/task_handler_test.go`

- [ ] Aggregate target and competitor mention rate, rank, sentiment, and evidence.
- [x] Upsert project competitor summaries after task completion.
- [x] Persist a completed report containing competitive score and gap facts.
- [x] Keep partial model failures visible in report metadata.
- [ ] Add handler tests for complete and partially failed tasks.

### Task 7: Verify the P0 Loop

- [x] Run all Go tests in the affected modules.
- [ ] Run migration validation against the development database.
- [ ] Verify task creation, execution, progress, answer analysis, competitor summary, and report endpoints.
- [x] Run `git diff --check`.

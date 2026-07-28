package task

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const taskCols = `id, project_id, user_id, model, status, questions_count, completed_count, COALESCE(total_tokens, 0) AS total_tokens, COALESCE(total_cost, 0) AS total_cost, COALESCE(error_message, '') AS error_message, started_at, finished_at, created_at, updated_at`

const answerCols = `id, task_id, question, answer, model, brand_mentioned, COALESCE(sentiment, '') AS sentiment, rank_position, analysis, created_at`

type DashboardStats struct {
	Projects        int      `json:"projects"`
	Tasks           int      `json:"tasks"`
	Completed       int      `json:"completed"`
	VisibilityScore int      `json:"visibilityScore"`
	RecentTasks     []AITask `json:"recentTasks"`
}

type Repository interface {
	Create(ctx context.Context, req CreateTaskRequest, userID string) (*AITask, error)
	FindByID(ctx context.Context, id string) (*AITask, error)
	ListByProject(ctx context.Context, projectID string, offset, limit int) ([]AITask, int, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]AITask, int, error)
	GetAnswers(ctx context.Context, taskID string) ([]AIAnswer, error)
	GetReport(ctx context.Context, taskID string) (*TaskReport, error)
	GetComparisonReport(ctx context.Context, taskID string) (*ComparisonReport, error)
	UpdateStatus(ctx context.Context, id string, status string, completedCount int, errorMsg string) error
	SetQuestionsCount(ctx context.Context, id string, count int) error
	SaveAnswer(ctx context.Context, answer *AIAnswer) error
	GetDashboardStats(ctx context.Context, userID string) (*DashboardStats, error)
	ListAll(ctx context.Context, offset, limit int) ([]AITask, int, error)
	Delete(ctx context.Context, id string) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, req CreateTaskRequest, userID string) (*AITask, error) {
	model := req.Model
	if model == "" {
		model = "gpt-4"
	}

	t := &AITask{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO ai_tasks (project_id, user_id, model, status, questions_count, completed_count)
		 VALUES ($1, $2, $3, 'pending', 0, 0)
		 RETURNING `+taskCols,
		req.ProjectID, userID, model,
	).Scan(
		&t.ID, &t.ProjectID, &t.UserID, &t.Model, &t.Status,
		&t.QuestionsCount, &t.CompletedCount, &t.TotalTokens, &t.TotalCost, &t.ErrorMessage,
		&t.StartedAt, &t.FinishedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	setProgress(t)
	return t, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*AITask, error) {
	t := &AITask{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+taskCols+` FROM ai_tasks WHERE id = $1`, id,
	).Scan(
		&t.ID, &t.ProjectID, &t.UserID, &t.Model, &t.Status,
		&t.QuestionsCount, &t.CompletedCount, &t.TotalTokens, &t.TotalCost, &t.ErrorMessage,
		&t.StartedAt, &t.FinishedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find task by id: %w", err)
	}
	setProgress(t)
	return t, nil
}

func (r *PostgresRepository) ListByProject(ctx context.Context, projectID string, offset, limit int) ([]AITask, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ai_tasks WHERE project_id = $1`, projectID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count tasks by project: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.project_id, t.user_id, t.model, t.status, t.questions_count, t.completed_count, COALESCE(t.total_tokens, 0), COALESCE(t.total_cost, 0), COALESCE(t.error_message, '') AS error_message, t.started_at, t.finished_at, t.created_at, t.updated_at, COALESCE(bp.name, '') AS project_name FROM ai_tasks t
		 LEFT JOIN brand_projects bp ON bp.id = t.project_id
		 WHERE t.project_id = $1
		 ORDER BY t.created_at DESC LIMIT $2 OFFSET $3`,
		projectID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list tasks by project: %w", err)
	}
	defer rows.Close()

	var tasks []AITask
	for rows.Next() {
		var t AITask
		err := rows.Scan(
			&t.ID, &t.ProjectID, &t.UserID, &t.Model, &t.Status,
			&t.QuestionsCount, &t.CompletedCount, &t.TotalTokens, &t.TotalCost, &t.ErrorMessage,
			&t.StartedAt, &t.FinishedAt, &t.CreatedAt, &t.UpdatedAt,
			&t.ProjectName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan task: %w", err)
		}
		setProgress(&t)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []AITask{}
	}
	return tasks, total, nil
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string, offset, limit int) ([]AITask, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ai_tasks WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count tasks by user: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.project_id, t.user_id, t.model, t.status, t.questions_count, t.completed_count, COALESCE(t.total_tokens, 0), COALESCE(t.total_cost, 0), COALESCE(t.error_message, '') AS error_message, t.started_at, t.finished_at, t.created_at, t.updated_at, COALESCE(bp.name, '') AS project_name FROM ai_tasks t
		 LEFT JOIN brand_projects bp ON bp.id = t.project_id
		 WHERE t.user_id = $1
		 ORDER BY t.created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list tasks by user: %w", err)
	}
	defer rows.Close()

	var tasks []AITask
	for rows.Next() {
		var t AITask
		err := rows.Scan(
			&t.ID, &t.ProjectID, &t.UserID, &t.Model, &t.Status,
			&t.QuestionsCount, &t.CompletedCount, &t.TotalTokens, &t.TotalCost, &t.ErrorMessage,
			&t.StartedAt, &t.FinishedAt, &t.CreatedAt, &t.UpdatedAt,
			&t.ProjectName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan task: %w", err)
		}
		setProgress(&t)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []AITask{}
	}
	return tasks, total, nil
}

func (r *PostgresRepository) GetAnswers(ctx context.Context, taskID string) ([]AIAnswer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+answerCols+` FROM ai_answers WHERE task_id = $1 ORDER BY created_at ASC`, taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("get answers: %w", err)
	}
	defer rows.Close()

	var answers []AIAnswer
	for rows.Next() {
		var a AIAnswer
		err := rows.Scan(
			&a.ID, &a.TaskID, &a.Question, &a.Answer, &a.Model,
			&a.BrandMentioned, &a.Sentiment, &a.RankPosition, &a.Analysis, &a.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan answer: %w", err)
		}
		answers = append(answers, a)
	}
	if answers == nil {
		answers = []AIAnswer{}
	}
	return answers, nil
}

func (r *PostgresRepository) GetReport(ctx context.Context, taskID string) (*TaskReport, error) {
	task, err := r.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	var projectBrief ProjectBrief
	err = r.pool.QueryRow(ctx,
		`SELECT id, name, COALESCE(website, '') AS website, industry FROM brand_projects WHERE id = $1`,
		task.ProjectID,
	).Scan(&projectBrief.ID, &projectBrief.Name, &projectBrief.Website, &projectBrief.Industry)
	if err != nil {
		return nil, fmt.Errorf("get project brief: %w", err)
	}

	answers, err := r.GetAnswers(ctx, taskID)
	if err != nil {
		return nil, err
	}

	brandMentions := 0
	for _, a := range answers {
		if a.BrandMentioned {
			brandMentions++
		}
	}

	totalQuestions := len(answers)
	visibilityScore := 0
	if totalQuestions > 0 {
		visibilityScore = brandMentions * 100 / totalQuestions
	}

	recommendations := generateRecommendations(answers, visibilityScore)

	return &TaskReport{
		Task:            *task,
		Project:         projectBrief,
		Answers:         answers,
		TotalQuestions:  totalQuestions,
		BrandMentions:   brandMentions,
		VisibilityScore: visibilityScore,
		Recommendations: recommendations,
	}, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status string, completedCount int, errorMsg string) error {
	now := time.Now()

	switch status {
	case "running":
		result, err := r.pool.Exec(ctx,
			`UPDATE ai_tasks SET status = $2, started_at = $3, updated_at = $3 WHERE id = $1 AND status = 'pending'`,
			id, status, now,
		)
		if err != nil {
			return fmt.Errorf("update task status to running: %w", err)
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("task %s is already running or finished", id)
		}
	case "completed":
		_, err := r.pool.Exec(ctx,
			`UPDATE ai_tasks SET status = $2, completed_count = $3, finished_at = $4, updated_at = $4 WHERE id = $1`,
			id, status, completedCount, now,
		)
		if err != nil {
			return fmt.Errorf("update task status to completed: %w", err)
		}
	case "failed":
		_, err := r.pool.Exec(ctx,
			`UPDATE ai_tasks SET status = $2, error_message = $3, finished_at = $4, updated_at = $4 WHERE id = $1`,
			id, status, errorMsg, now,
		)
		if err != nil {
			return fmt.Errorf("update task status to failed: %w", err)
		}
	default:
		_, err := r.pool.Exec(ctx,
			`UPDATE ai_tasks SET status = $2, updated_at = NOW() WHERE id = $1`,
			id, status,
		)
		if err != nil {
			return fmt.Errorf("update task status: %w", err)
		}
	}

	return nil
}

func (r *PostgresRepository) SetQuestionsCount(ctx context.Context, id string, count int) error {
	if count < 0 {
		count = 0
	}
	_, err := r.pool.Exec(ctx,
		`UPDATE ai_tasks SET questions_count = $2, updated_at = NOW() WHERE id = $1`, id, count)
	if err != nil {
		return fmt.Errorf("set task questions count: %w", err)
	}
	return nil
}

func (r *PostgresRepository) SaveAnswer(ctx context.Context, answer *AIAnswer) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO ai_answers (task_id, question, answer, model, brand_mentioned, sentiment, rank_position, analysis, tokens_used, cost)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, created_at`,
		answer.TaskID, answer.Question, answer.Answer, answer.Model,
		answer.BrandMentioned, answer.Sentiment, answer.RankPosition, answer.Analysis,
		answer.TokensUsed, answer.Cost,
	).Scan(&answer.ID, &answer.CreatedAt)
	if err != nil {
		return fmt.Errorf("save answer: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetDashboardStats(ctx context.Context, userID string) (*DashboardStats, error) {
	var totalProjects int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM brand_projects WHERE user_id = $1 AND status != 'archived'`, userID,
	).Scan(&totalProjects)
	if err != nil {
		return nil, fmt.Errorf("count projects: %w", err)
	}

	var totalTasks int
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ai_tasks WHERE user_id = $1`, userID,
	).Scan(&totalTasks)
	if err != nil {
		return nil, fmt.Errorf("count tasks: %w", err)
	}

	var completedTasks int
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ai_tasks WHERE user_id = $1 AND status = 'completed'`, userID,
	).Scan(&completedTasks)
	if err != nil {
		return nil, fmt.Errorf("count completed tasks: %w", err)
	}

	rows, err := r.pool.Query(ctx, `
		SELECT COALESCE(
			SUM(CASE WHEN a.brand_mentioned THEN 1 ELSE 0 END) * 100 / NULLIF(COUNT(*), 0)
		, 0)
		FROM ai_tasks t
		JOIN ai_answers a ON a.task_id = t.id
		WHERE t.user_id = $1 AND t.status = 'completed'
		GROUP BY t.id`, userID)
	if err != nil {
		return nil, fmt.Errorf("query visibility scores: %w", err)
	}
	defer rows.Close()

	var scores []int
	for rows.Next() {
		var score int
		if err := rows.Scan(&score); err != nil {
			return nil, fmt.Errorf("scan visibility score: %w", err)
		}
		scores = append(scores, score)
	}

	avgScore := 0
	if len(scores) > 0 {
		sum := 0
		for _, s := range scores {
			sum += s
		}
		avgScore = sum / len(scores)
	}

	recentTasks, _, err := r.ListByUser(ctx, userID, 0, 5)
	if err != nil {
		return nil, fmt.Errorf("list recent tasks: %w", err)
	}

	return &DashboardStats{
		Projects:        totalProjects,
		Tasks:           totalTasks,
		Completed:       completedTasks,
		VisibilityScore: avgScore,
		RecentTasks:     recentTasks,
	}, nil
}

func (r *PostgresRepository) ListAll(ctx context.Context, offset, limit int) ([]AITask, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM ai_tasks`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count all tasks: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.project_id, t.user_id, t.model, t.status, t.questions_count, t.completed_count, COALESCE(t.total_tokens, 0), COALESCE(t.total_cost, 0), COALESCE(t.error_message, '') AS error_message, t.started_at, t.finished_at, t.created_at, t.updated_at, COALESCE(bp.name, '') AS project_name, COALESCE(u.username, '') AS username FROM ai_tasks t
		 LEFT JOIN brand_projects bp ON bp.id = t.project_id
		 LEFT JOIN users u ON u.id = t.user_id
		 ORDER BY t.created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list all tasks: %w", err)
	}
	defer rows.Close()

	var tasks []AITask
	for rows.Next() {
		var t AITask
		err := rows.Scan(
			&t.ID, &t.ProjectID, &t.UserID, &t.Model, &t.Status,
			&t.QuestionsCount, &t.CompletedCount, &t.TotalTokens, &t.TotalCost, &t.ErrorMessage,
			&t.StartedAt, &t.FinishedAt, &t.CreatedAt, &t.UpdatedAt,
			&t.ProjectName, &t.Username,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan task: %w", err)
		}
		setProgress(&t)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []AITask{}
	}
	return tasks, total, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin delete task: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM ai_answers WHERE task_id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete task answers: %w", err)
	}
	result, err := tx.Exec(ctx, `DELETE FROM ai_tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("task %s not found", id)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit delete task: %w", err)
	}
	return nil
}

func setProgress(t *AITask) {
	if t.QuestionsCount <= 0 {
		t.Progress = 0
		if t.Status == "completed" {
			t.Progress = 100
		}
		return
	}
	progress := t.CompletedCount * 100 / t.QuestionsCount
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	t.Progress = progress
}

func (r *PostgresRepository) GetComparisonReport(ctx context.Context, taskID string) (*ComparisonReport, error) {
	t, err := r.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	answers, err := r.GetAnswers(ctx, taskID)
	if err != nil {
		return nil, err
	}

	modelAnswers := make(map[string][]AIAnswer)
	for _, a := range answers {
		modelAnswers[a.Model] = append(modelAnswers[a.Model], a)
	}

	var results []ModelComparisonResult
	for model, ans := range modelAnswers {
		brandMentions := 0
		for _, a := range ans {
			if a.BrandMentioned {
				brandMentions++
			}
		}
		score := 0.0
		if len(ans) > 0 {
			score = float64(brandMentions*100) / float64(len(ans))
		}
		results = append(results, ModelComparisonResult{
			Model:   model,
			Status:  "completed",
			Answers: ans,
			Score:   score,
		})
	}

	var projectBrief ProjectBrief
	err = r.pool.QueryRow(ctx,
		`SELECT id, name, COALESCE(website, '') AS website, industry FROM brand_projects WHERE id = $1`,
		t.ProjectID,
	).Scan(&projectBrief.ID, &projectBrief.Name, &projectBrief.Website, &projectBrief.Industry)
	if err != nil {
		return nil, fmt.Errorf("get project brief: %w", err)
	}

	brandMentions := 0
	for _, a := range answers {
		if a.BrandMentioned {
			brandMentions++
		}
	}
	totalQuestions := len(answers)
	visibilityScore := 0
	if totalQuestions > 0 {
		visibilityScore = brandMentions * 100 / totalQuestions
	}
	recommendations := generateRecommendations(answers, visibilityScore)

	return &ComparisonReport{
		Task: TaskReport{
			Task:            *t,
			Project:         projectBrief,
			Answers:         answers,
			TotalQuestions:  totalQuestions,
			BrandMentions:   brandMentions,
			VisibilityScore: visibilityScore,
			Recommendations: recommendations,
		},
		Results: results,
	}, nil
}

func generateRecommendations(answers []AIAnswer, visibilityScore int) []string {
	var recs []string

	if visibilityScore < 30 {
		recs = append(recs, "品牌可见度较低，建议增加品牌在 AI 对话中的提及频率，优化品牌关键词策略")
	} else if visibilityScore < 70 {
		recs = append(recs, "品牌可见度中等，建议持续优化内容策略，进一步提升品牌曝光")
	} else {
		recs = append(recs, "品牌可见度较高，建议保持当前策略并关注竞品动态")
	}

	negativeCount := 0
	neutralCount := 0
	positiveCount := 0
	for _, a := range answers {
		switch a.Sentiment {
		case "positive":
			positiveCount++
		case "negative":
			negativeCount++
		case "neutral":
			neutralCount++
		}
	}

	if negativeCount > positiveCount {
		recs = append(recs, "负面情感占比较高，建议排查用户痛点并优化产品体验")
	}
	if neutralCount > 0 && negativeCount == 0 && positiveCount == 0 {
		recs = append(recs, "情感倾向以中性为主，建议增强品牌情感化表达")
	}
	if positiveCount > negativeCount && positiveCount > 0 {
		recs = append(recs, "正面情感占优，建议将正面反馈转化为营销素材")
	}

	topRankCount := 0
	for _, a := range answers {
		if a.RankPosition != nil && *a.RankPosition <= 3 {
			topRankCount++
		}
	}
	if len(answers) > 0 && topRankCount < len(answers)/2 {
		recs = append(recs, "排名位置偏低，建议加强 SEO 和内容营销以提升搜索结果排名")
	}

	if len(recs) == 0 {
		recs = append(recs, "品牌表现良好，建议定期监测以保持优势")
	}

	return recs
}

-- AI-Growth-Engine Phase 2 迁移
-- 品牌项目管理

-- ============================================
-- 品牌项目表
-- ============================================
CREATE TABLE brand_projects (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id),
    name            VARCHAR(255) NOT NULL,
    website         VARCHAR(500),
    industry        VARCHAR(100) NOT NULL,
    description     TEXT,
    keywords        TEXT[],
    service_area    VARCHAR(200),
    target_users    TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bp_user ON brand_projects(user_id);
CREATE INDEX idx_bp_status ON brand_projects(status);

-- ============================================
-- AI 检测任务表
-- ============================================
CREATE TABLE ai_tasks (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id      UUID NOT NULL REFERENCES brand_projects(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    model           VARCHAR(50) NOT NULL DEFAULT 'gpt-4',
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    questions_count INTEGER DEFAULT 0,
    completed_count INTEGER DEFAULT 0,
    error_message   TEXT,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ait_project ON ai_tasks(project_id);
CREATE INDEX idx_ait_user ON ai_tasks(user_id);
CREATE INDEX idx_ait_status ON ai_tasks(status);

-- ============================================
-- AI 回答表
-- ============================================
CREATE TABLE ai_answers (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id         UUID NOT NULL REFERENCES ai_tasks(id),
    question        TEXT NOT NULL,
    answer          TEXT NOT NULL,
    model           VARCHAR(50) NOT NULL,
    brand_mentioned BOOLEAN DEFAULT false,
    sentiment       VARCHAR(20),
    rank_position   INTEGER,
    analysis        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_aia_task ON ai_answers(task_id);

-- ============================================
-- updated_at 触发器
-- ============================================
CREATE TRIGGER brand_projects_updated_at
    BEFORE UPDATE ON brand_projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER ai_tasks_updated_at
    BEFORE UPDATE ON ai_tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

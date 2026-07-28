-- Plans table
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    monthly_price DECIMAL(10,2) DEFAULT 0,
    yearly_price DECIMAL(10,2) DEFAULT 0,
    max_projects INTEGER DEFAULT 1,
    max_ai_queries INTEGER DEFAULT 10,
    max_reports INTEGER DEFAULT 5,
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Seed plans
INSERT INTO plans (name, code, description, monthly_price, yearly_price, max_projects, max_ai_queries, sort_order)
VALUES
('免费版', 'free', '个人体验，了解品牌 AI 可见度', 0, 0, 1, 10, 1),
('专业版', 'pro', '适合中小企业和个人品牌', 199, 1990, 10, 100, 2),
('企业版', 'business', '适合营销团队和代理机构', 499, 4990, 50, 500, 3),
('定制版', 'enterprise', '适合大型企业和私有部署需求', 0, 0, -1, -1, 4)
ON CONFLICT (code) DO NOTHING;

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    order_no VARCHAR(100) NOT NULL UNIQUE,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'CNY',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(50),
    payment_time TIMESTAMPTZ,
    description TEXT,
    plan_id UUID REFERENCES plans(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

-- Prompts table
CREATE TABLE IF NOT EXISTS prompts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    purpose VARCHAR(100) NOT NULL,
    version VARCHAR(20) DEFAULT 'v1.0',
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    created_by UUID REFERENCES admins(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Brand infos table
CREATE TABLE IF NOT EXISTS brand_infos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES brand_projects(id) ON DELETE CASCADE,
    brand_intro TEXT,
    product_intro TEXT,
    service_intro TEXT,
    faq JSONB DEFAULT '[]',
    advantages TEXT[] DEFAULT '{}',
    cases JSONB DEFAULT '[]',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id)
);

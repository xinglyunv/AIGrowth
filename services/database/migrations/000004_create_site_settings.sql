CREATE TABLE IF NOT EXISTS site_settings (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO site_settings (key, value) VALUES
    ('site_name', 'AI Growth Engine'),
    ('site_title', 'AI 品牌可见度分析平台'),
    ('site_description', 'AI 品牌可见度分析与增长优化 SaaS 平台'),
    ('logo_url', ''),
    ('footer_text', 'AI Growth Engine. All rights reserved.')
ON CONFLICT (key) DO NOTHING;

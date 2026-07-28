INSERT INTO site_settings (key, value) VALUES
    ('contact_email', 'contact@aige.com'),
    ('contact_address', '北京市海淀区中关村科技园区'),
    ('contact_phone', '400-888-8888'),
    ('working_hours', '周一至周五 9:00 - 18:00'),
    ('stat_companies', '100+'),
    ('stat_models', '10+'),
    ('stat_reports', '1000+'),
    ('stat_accuracy', '99.5%'),
    ('hero_tagline', 'AI 品牌可见度分析平台'),
    ('hero_title', '让 AI 认识你的品牌'),
    ('hero_subtitle', 'AI 品牌可见度分析与增长优化平台，了解 AI 如何看待你的品牌，发现增长机会')
ON CONFLICT (key) DO NOTHING;

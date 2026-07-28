export const SITE_THEMES = ['syntro', 'saas-landing', 'saas-ui', 'hikari', 'cruip', 'nextly'] as const;
export type SiteTheme = (typeof SITE_THEMES)[number];

export interface SiteConfig {
  site_name: string;
  site_title: string;
  site_description: string;
  site_theme: SiteTheme;
  logo_url: string;
  footer_text: string;
  contact_email: string;
  contact_address: string;
  contact_phone: string;
  working_hours: string;
  navigation: {
    items: { label: string; href: string }[];
    login_label: string;
    login_href: string;
    cta_label: string;
    cta_href: string;
  };
  hero: {
    tagline: string;
    title: string;
    subtitle: string;
    primary_label: string;
    primary_href: string;
    secondary_label: string;
    secondary_href: string;
  };
  features: { eyebrow: string; title: string; description: string; href: string; highlights: string[]; icon: string }[];
  stats: { items: { value: string; label: string }[] };
  trust: { eyebrow: string; items: string[] };
  footer: { description: string; columns: { title: string; items: { label: string; href: string }[] }[]; copyright: string };
  dashboard: { score: string; score_delta: string; mentions: string; recommendation_rate: string; competitor_gap: string; trend_label: string };
  seo: { title: string; description: string; og_image: string };
  theme: { primary: string; accent: string; surface: string; text: string; muted: string };
}

export const defaultSiteConfig: SiteConfig = {
  site_name: 'AI Growth Engine',
  site_title: 'AI 品牌可见度分析平台',
  site_description: 'AI 品牌可见度分析与增长优化 SaaS 平台',
  site_theme: 'syntro',
  logo_url: '',
  footer_text: 'AI Growth Engine',
  contact_email: 'contact@aige.com',
  contact_address: '北京市海淀区中关村科技园区',
  contact_phone: '400-888-8888',
  working_hours: '周一至周五 9:00 - 18:00',
  navigation: {
    items: [{ label: '功能', href: '/features' }, { label: '定价', href: '/pricing' }, { label: '资源', href: '/resources' }, { label: '关于', href: '/about' }, { label: '联系我们', href: '/contact' }],
    login_label: '登录', login_href: '/app', cta_label: '开始使用', cta_href: '/app',
  },
  hero: {
    tagline: 'AI 品牌可见度分析平台', title: '让 AI 认识你的品牌', subtitle: 'AI 品牌可见度分析与增长优化平台，了解 AI 如何看待你的品牌，发现增长机会',
    primary_label: '开始分析', primary_href: '/app', secondary_label: '查看方法', secondary_href: '#features',
  },
  features: [
    { eyebrow: '品牌认知', title: '看见 AI 如何描述你的品牌', description: '从多模型、多场景、多轮对话中还原品牌真实认知。', href: '/features', highlights: ['跨模型品牌提及率对比', '品牌推荐情感分析', '行业关联度识别'], icon: 'brain' },
    { eyebrow: '增长机会', title: '找到影响推荐的关键因素', description: '定位内容、产品和口碑中的增长缺口，形成优先级清晰的行动方案。', href: '/features', highlights: ['多品牌可见度横向对比', '竞品优势策略拆解', '行业差距量化分析'], icon: 'chart' },
    { eyebrow: '持续追踪', title: '让每次优化都有数据反馈', description: '持续监测品牌在 AI 答案中的表现变化，验证每项策略的真实影响。', href: '/features', highlights: ['自动化持续监控', '异常变化告警推送', '历史趋势对比分析'], icon: 'trend' },
  ],
  stats: { items: [{ value: '100+', label: '品牌持续追踪' }, { value: '10+', label: '主流 AI 模型' }, { value: '1,000+', label: '分析报告生成' }, { value: '99.5%', label: '数据准确率' }] },
  trust: { eyebrow: '被增长团队选择', items: ['ChatGPT', 'Claude', 'Gemini', 'DeepSeek', '通义千问'] },
  footer: { description: '帮助品牌建立可被 AI 理解、信任和推荐的增长系统。', columns: [], copyright: '© 2026 AI Growth Engine. All rights reserved.' },
  dashboard: { score: '82.4', score_delta: '+18.6%', mentions: '1,284', recommendation_rate: '64.8%', competitor_gap: '-12.4', trend_label: '过去 30 天' },
  seo: { title: 'AI 品牌可见度分析平台', description: 'AI 品牌可见度分析与增长优化 SaaS 平台', og_image: '' },
  theme: { primary: '#6d5efc', accent: '#2dd4bf', surface: '#f8fafc', text: '#0f172a', muted: '#64748b' },
};

function mergeConfig(data: Partial<SiteConfig>): SiteConfig {
  const siteTheme = SITE_THEMES.includes(data.site_theme as SiteTheme) ? data.site_theme as SiteTheme : defaultSiteConfig.site_theme;
  const theme = { ...defaultSiteConfig.theme, ...data.theme };
  for (const key of Object.keys(theme) as (keyof SiteConfig['theme'])[]) {
    if (!/^#[0-9a-f]{6}$/i.test(theme[key])) theme[key] = defaultSiteConfig.theme[key];
  }
  return {
    ...defaultSiteConfig, ...data, site_theme: siteTheme,
    navigation: { ...defaultSiteConfig.navigation, ...data.navigation, items: data.navigation?.items?.length ? data.navigation.items : defaultSiteConfig.navigation.items },
    hero: { ...defaultSiteConfig.hero, ...data.hero },
    features: data.features?.length ? data.features : defaultSiteConfig.features,
    stats: { ...defaultSiteConfig.stats, ...data.stats, items: data.stats?.items?.length ? data.stats.items : defaultSiteConfig.stats.items },
    trust: { ...defaultSiteConfig.trust, ...data.trust, items: data.trust?.items?.length ? data.trust.items : defaultSiteConfig.trust.items },
    footer: { ...defaultSiteConfig.footer, ...data.footer },
    dashboard: { ...defaultSiteConfig.dashboard, ...data.dashboard },
    seo: { ...defaultSiteConfig.seo, ...data.seo },
    theme,
  };
}

export async function loadSiteConfig(): Promise<SiteConfig> {
  try {
    const response = await fetch('http://localhost:8080/api/v1/settings');
    if (!response.ok) return defaultSiteConfig;
    const body = await response.json();
    return body.success && body.data ? mergeConfig(body.data) : defaultSiteConfig;
  } catch {
    return defaultSiteConfig;
  }
}

package setting

import (
	"fmt"
	"regexp"
)

type NavItem struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

type NavigationConfig struct {
	Items      []NavItem `json:"items"`
	LoginLabel string    `json:"login_label"`
	LoginHref  string    `json:"login_href"`
	CTALabel   string    `json:"cta_label"`
	CTAHref    string    `json:"cta_href"`
}

type HeroConfig struct {
	Tagline        string `json:"tagline"`
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	PrimaryLabel   string `json:"primary_label"`
	PrimaryHref    string `json:"primary_href"`
	SecondaryLabel string `json:"secondary_label"`
	SecondaryHref  string `json:"secondary_href"`
}

type FeatureItem struct {
	Eyebrow     string   `json:"eyebrow"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Href        string   `json:"href"`
	Highlights  []string `json:"highlights"`
	Icon        string   `json:"icon"`
}

type StatItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type StatsConfig struct {
	Items []StatItem `json:"items"`
}

type TrustConfig struct {
	Eyebrow string   `json:"eyebrow"`
	Items   []string `json:"items"`
}

type FooterColumn struct {
	Title string    `json:"title"`
	Items []NavItem `json:"items"`
}

type FooterConfig struct {
	Description string         `json:"description"`
	Columns     []FooterColumn `json:"columns"`
	Copyright   string         `json:"copyright"`
}

type DashboardConfig struct {
	Score              string `json:"score"`
	ScoreDelta         string `json:"score_delta"`
	Mentions           string `json:"mentions"`
	RecommendationRate string `json:"recommendation_rate"`
	CompetitorGap      string `json:"competitor_gap"`
	TrendLabel         string `json:"trend_label"`
}

type SEOConfig struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OGImage     string `json:"og_image"`
}

type ThemeConfig struct {
	Primary string `json:"primary"`
	Accent  string `json:"accent"`
	Surface string `json:"surface"`
	Text    string `json:"text"`
	Muted   string `json:"muted"`
}

func normalizeTheme(theme string) string {
	switch theme {
	case "syntro", "saas-landing", "saas-ui", "hikari", "cruip", "nextly":
		return theme
	default:
		return "syntro"
	}
}

func ValidateSiteConfig(config *SiteConfig) error {
	if config == nil {
		return fmt.Errorf("site config is required")
	}
	colorPattern := regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	colors := map[string]string{
		"primary": config.Theme.Primary, "accent": config.Theme.Accent, "surface": config.Theme.Surface,
		"text": config.Theme.Text, "muted": config.Theme.Muted,
	}
	for name, value := range colors {
		if !colorPattern.MatchString(value) {
			return fmt.Errorf("theme %s must be a six-digit hex color", name)
		}
	}
	if config.SiteTheme != normalizeTheme(config.SiteTheme) {
		return fmt.Errorf("unsupported site theme %q", config.SiteTheme)
	}
	return nil
}

type Setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SiteConfig struct {
	SiteName        string           `json:"site_name"`
	SiteTitle       string           `json:"site_title"`
	SiteDescription string           `json:"site_description"`
	SiteTheme       string           `json:"site_theme"`
	LogoURL         string           `json:"logo_url"`
	FooterText      string           `json:"footer_text"`
	Navigation      NavigationConfig `json:"navigation"`
	Hero            HeroConfig       `json:"hero"`
	Features        []FeatureItem    `json:"features"`
	Stats           StatsConfig      `json:"stats"`
	Trust           TrustConfig      `json:"trust"`
	Footer          FooterConfig     `json:"footer"`
	Dashboard       DashboardConfig  `json:"dashboard"`
	SEO             SEOConfig        `json:"seo"`
	Theme           ThemeConfig      `json:"theme"`

	ContactEmail   string `json:"contact_email"`
	ContactAddress string `json:"contact_address"`
	ContactPhone   string `json:"contact_phone"`
	WorkingHours   string `json:"working_hours"`

	StatCompanies string `json:"stat_companies"`
	StatModels    string `json:"stat_models"`
	StatReports   string `json:"stat_reports"`
	StatAccuracy  string `json:"stat_accuracy"`

	HeroTagline  string `json:"hero_tagline"`
	HeroTitle    string `json:"hero_title"`
	HeroSubtitle string `json:"hero_subtitle"`

	AllowRegistration string `json:"allow_registration"`

	SmtpHost     string `json:"smtp_host"`
	SmtpPort     string `json:"smtp_port"`
	SmtpUser     string `json:"smtp_user"`
	SmtpPassword string `json:"smtp_password"`
	SmtpFrom     string `json:"smtp_from"`

	SmsProvider    string `json:"sms_provider"`
	SmsAccessKey   string `json:"sms_access_key"`
	SmsSecretKey   string `json:"sms_secret_key"`
	SmsSignName    string `json:"sms_sign_name"`
	SmsbaoUsername string `json:"smsbao_username"`
	SmsbaoPassword string `json:"smsbao_password"`
}

func defaultSiteConfig() *SiteConfig {
	return &SiteConfig{
		SiteName:          "AI Growth Engine",
		SiteTitle:         "AI 品牌可见度分析平台",
		SiteDescription:   "AI 品牌可见度分析与增长优化 SaaS 平台",
		SiteTheme:         "syntro",
		FooterText:        "AI Growth Engine. All rights reserved.",
		ContactEmail:      "contact@aige.com",
		ContactAddress:    "北京市海淀区中关村科技园区",
		ContactPhone:      "400-888-8888",
		WorkingHours:      "周一至周五 9:00 - 18:00",
		StatCompanies:     "100+",
		StatModels:        "10+",
		StatReports:       "1000+",
		StatAccuracy:      "99.5%",
		HeroTagline:       "AI 品牌可见度分析平台",
		HeroTitle:         "让 AI 认识你的品牌",
		HeroSubtitle:      "AI 品牌可见度分析与增长优化平台，了解 AI 如何看待你的品牌，发现增长机会",
		AllowRegistration: "true",
		Navigation: NavigationConfig{
			Items:      []NavItem{{Label: "产品能力", Href: "#features"}, {Label: "分析方法", Href: "#method"}, {Label: "客户案例", Href: "#customers"}, {Label: "资源中心", Href: "#resources"}},
			LoginLabel: "登录", LoginHref: "/login", CTALabel: "开始分析", CTAHref: "/register",
		},
		Hero: HeroConfig{
			Tagline: "AI 品牌可见度分析平台", Title: "让 AI 认识你的品牌", Subtitle: "AI 品牌可见度分析与增长优化平台，了解 AI 如何看待你的品牌，发现增长机会",
			PrimaryLabel: "开始分析", PrimaryHref: "/register", SecondaryLabel: "查看方法", SecondaryHref: "#method",
		},
		Features: []FeatureItem{
			{Eyebrow: "品牌认知", Title: "看见 AI 如何描述你的品牌", Description: "从多模型、多场景、多轮对话中还原品牌真实认知。", Href: "#visibility", Highlights: []string{"跨模型品牌提及率对比", "品牌推荐情感分析", "行业关联度识别"}, Icon: "brain"},
			{Eyebrow: "增长机会", Title: "找到影响推荐的关键因素", Description: "定位内容、产品和口碑中的增长缺口，形成优先级清晰的行动方案。", Href: "#growth", Highlights: []string{"多品牌可见度横向对比", "竞品优势策略拆解", "行业差距量化分析"}, Icon: "chart"},
			{Eyebrow: "持续追踪", Title: "让每次优化都有数据反馈", Description: "持续监测品牌在 AI 答案中的表现变化，验证每项策略的真实影响。", Href: "#tracking", Highlights: []string{"自动化持续监控", "异常变化告警推送", "历史趋势对比分析"}, Icon: "trend"},
		},
		Stats:     StatsConfig{Items: []StatItem{{Value: "100+", Label: "品牌持续追踪"}, {Value: "10+", Label: "主流 AI 模型"}, {Value: "1,000+", Label: "分析报告生成"}, {Value: "99.5%", Label: "数据准确率"}}},
		Trust:     TrustConfig{Eyebrow: "被增长团队选择", Items: []string{"产品团队", "市场团队", "内容团队", "品牌团队"}},
		Footer:    FooterConfig{Description: "帮助品牌建立可被 AI 理解、信任和推荐的增长系统。", Copyright: "© 2026 AI Growth Engine. All rights reserved."},
		Dashboard: DashboardConfig{Score: "82.4", ScoreDelta: "+18.6%", Mentions: "1,284", RecommendationRate: "64.8%", CompetitorGap: "-12.4", TrendLabel: "过去 30 天"},
		SEO:       SEOConfig{Title: "AI 品牌可见度分析平台", Description: "AI 品牌可见度分析与增长优化 SaaS 平台"},
		Theme:     ThemeConfig{Primary: "#6d5efc", Accent: "#2dd4bf", Surface: "#f8fafc", Text: "#0f172a", Muted: "#64748b"},
	}
}

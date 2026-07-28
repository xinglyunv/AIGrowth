package setting

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetAll(ctx context.Context) (*SiteConfig, error)
	GetByKey(ctx context.Context, key string) (string, error)
	Update(ctx context.Context, cfg *SiteConfig) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) GetByKey(ctx context.Context, key string) (string, error) {
	var value string
	err := r.pool.QueryRow(ctx, `SELECT value FROM site_settings WHERE key = $1`, key).Scan(&value)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get setting by key: %w", err)
	}
	return value, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) (*SiteConfig, error) {
	rows, err := r.pool.Query(ctx, `SELECT key, value FROM site_settings`)
	if err != nil {
		return nil, fmt.Errorf("get all settings: %w", err)
	}
	defer rows.Close()

	values := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, fmt.Errorf("scan setting row: %w", err)
		}
		values[k] = v
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate settings rows: %w", err)
	}

	cfg := defaultSiteConfig()

	fields := map[string]*string{
		"site_name":          &cfg.SiteName,
		"site_title":         &cfg.SiteTitle,
		"site_description":   &cfg.SiteDescription,
		"site_theme":         &cfg.SiteTheme,
		"logo_url":           &cfg.LogoURL,
		"footer_text":        &cfg.FooterText,
		"contact_email":      &cfg.ContactEmail,
		"contact_address":    &cfg.ContactAddress,
		"contact_phone":      &cfg.ContactPhone,
		"working_hours":      &cfg.WorkingHours,
		"stat_companies":     &cfg.StatCompanies,
		"stat_models":        &cfg.StatModels,
		"stat_reports":       &cfg.StatReports,
		"stat_accuracy":      &cfg.StatAccuracy,
		"hero_tagline":       &cfg.HeroTagline,
		"hero_title":         &cfg.HeroTitle,
		"hero_subtitle":      &cfg.HeroSubtitle,
		"allow_registration": &cfg.AllowRegistration,
		"smtp_host":          &cfg.SmtpHost,
		"smtp_port":          &cfg.SmtpPort,
		"smtp_user":          &cfg.SmtpUser,
		"smtp_password":      &cfg.SmtpPassword,
		"smtp_from":          &cfg.SmtpFrom,
		"sms_provider":       &cfg.SmsProvider,
		"sms_access_key":     &cfg.SmsAccessKey,
		"sms_secret_key":     &cfg.SmsSecretKey,
		"sms_sign_name":      &cfg.SmsSignName,
		"smsbao_username":    &cfg.SmsbaoUsername,
		"smsbao_password":    &cfg.SmsbaoPassword,
	}

	for k, ptr := range fields {
		if v, ok := values[k]; ok {
			*ptr = v
		}
	}
	cfg.SiteTheme = normalizeTheme(cfg.SiteTheme)

	jsonFields := map[string]any{
		"navigation": &cfg.Navigation,
		"hero":       &cfg.Hero,
		"features":   &cfg.Features,
		"stats":      &cfg.Stats,
		"trust":      &cfg.Trust,
		"footer":     &cfg.Footer,
		"dashboard":  &cfg.Dashboard,
		"seo":        &cfg.SEO,
		"theme":      &cfg.Theme,
	}
	for key, target := range jsonFields {
		if value, ok := values[key]; ok && value != "" {
			if err := json.Unmarshal([]byte(value), target); err != nil {
				return nil, fmt.Errorf("decode setting %s: %w", key, err)
			}
		}
	}

	return cfg, nil
}

func (r *PostgresRepository) Update(ctx context.Context, cfg *SiteConfig) error {
	upsert := `INSERT INTO site_settings (key, value) VALUES ($1, $2)
	           ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()`

	entries := map[string]string{
		"site_name":          cfg.SiteName,
		"site_title":         cfg.SiteTitle,
		"site_description":   cfg.SiteDescription,
		"site_theme":         cfg.SiteTheme,
		"logo_url":           cfg.LogoURL,
		"footer_text":        cfg.FooterText,
		"contact_email":      cfg.ContactEmail,
		"contact_address":    cfg.ContactAddress,
		"contact_phone":      cfg.ContactPhone,
		"working_hours":      cfg.WorkingHours,
		"stat_companies":     cfg.StatCompanies,
		"stat_models":        cfg.StatModels,
		"stat_reports":       cfg.StatReports,
		"stat_accuracy":      cfg.StatAccuracy,
		"hero_tagline":       cfg.HeroTagline,
		"hero_title":         cfg.HeroTitle,
		"hero_subtitle":      cfg.HeroSubtitle,
		"allow_registration": cfg.AllowRegistration,
		"smtp_host":          cfg.SmtpHost,
		"smtp_port":          cfg.SmtpPort,
		"smtp_user":          cfg.SmtpUser,
		"smtp_password":      cfg.SmtpPassword,
		"smtp_from":          cfg.SmtpFrom,
		"sms_provider":       cfg.SmsProvider,
		"sms_access_key":     cfg.SmsAccessKey,
		"sms_secret_key":     cfg.SmsSecretKey,
		"sms_sign_name":      cfg.SmsSignName,
		"smsbao_username":    cfg.SmsbaoUsername,
		"smsbao_password":    cfg.SmsbaoPassword,
	}

	jsonValues := map[string]any{
		"navigation": cfg.Navigation,
		"hero":       cfg.Hero,
		"features":   cfg.Features,
		"stats":      cfg.Stats,
		"trust":      cfg.Trust,
		"footer":     cfg.Footer,
		"dashboard":  cfg.Dashboard,
		"seo":        cfg.SEO,
		"theme":      cfg.Theme,
	}
	for key, value := range jsonValues {
		encoded, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("encode setting %s: %w", key, err)
		}
		entries[key] = string(encoded)
	}

	for k, v := range entries {
		if _, err := r.pool.Exec(ctx, upsert, k, v); err != nil {
			return fmt.Errorf("upsert setting %s: %w", k, err)
		}
	}

	return nil
}

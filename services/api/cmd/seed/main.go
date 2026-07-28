package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aige/admin"
	"github.com/aige/setting"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/ai_growth_engine?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	repo := admin.NewPostgresRepository(pool)

	existing, _ := repo.FindByEmail(context.Background(), "admin@aige.com")
	if existing == nil {
		a, err := repo.Create(context.Background(), admin.CreateAdminRequest{
			Username: "admin",
			Email:    "admin@aige.com",
			Password: "admin123",
			Role:     "superadmin",
		})
		if err != nil {
			log.Fatalf("create admin: %v", err)
		}
		fmt.Printf("Admin created: %s / %s (role: %s, id: %s)\n", a.Email, "admin123", a.Role, a.ID)
	} else {
		fmt.Println("Admin already exists, skipping admin seed.")
	}

	// Seed default settings
	settingRepo := setting.NewPostgresRepository(pool)
	existingCfg, _ := settingRepo.GetAll(context.Background())
	if existingCfg != nil && existingCfg.SiteName != "AI Growth Engine" {
		// Only seed if config was never initialized
		// Actually, GetAll always returns defaults, so we check if allow_registration exists
	}

	upsert := `INSERT INTO site_settings (key, value) VALUES ($1, $2)
	           ON CONFLICT (key) DO NOTHING`

	defaults := map[string]string{
		"site_name":           "AI Growth Engine",
		"site_title":          "AI 品牌可见度分析平台",
		"site_description":    "AI 品牌可见度分析与增长优化 SaaS 平台",
		"footer_text":         "AI Growth Engine. All rights reserved.",
		"contact_email":       "contact@aige.com",
		"contact_address":     "北京市海淀区中关村科技园区",
		"contact_phone":       "400-888-8888",
		"working_hours":       "周一至周五 9:00 - 18:00",
		"stat_companies":      "100+",
		"stat_models":         "10+",
		"stat_reports":        "1000+",
		"stat_accuracy":       "99.5%",
		"hero_tagline":        "AI 品牌可见度分析平台",
		"hero_title":          "让 AI 认识你的品牌",
		"hero_subtitle":       "AI 品牌可见度分析与增长优化平台，了解 AI 如何看待你的品牌，发现增长机会",
		"allow_registration":  "true",
	}

	for k, v := range defaults {
		if _, err := pool.Exec(context.Background(), upsert, k, v); err != nil {
			log.Printf("seed setting %s: %v", k, err)
		} else {
			fmt.Printf("Setting seeded: %s = %s\n", k, v)
		}
	}
}

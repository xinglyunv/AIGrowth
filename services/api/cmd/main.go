package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aige/admin"
	"github.com/aige/aimodel"
	"github.com/aige/api/internal/config"
	"github.com/aige/api/internal/router"
	"github.com/aige/audit"
	"github.com/aige/cdk"
	"github.com/aige/competitor"
	"github.com/aige/contact"
	"github.com/aige/notification"
	"github.com/aige/order"
	"github.com/aige/payment"
	"github.com/aige/plan"
	"github.com/aige/project"
	"github.com/aige/report"
	"github.com/aige/requestlog"
	"github.com/aige/setting"
	"github.com/aige/space"
	"github.com/aige/task"
	"github.com/aige/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	// Store pool reference in config for admin dashboard
	cfg.Pool = pool

	// Create repositories
	userRepo := user.NewPostgresRepository(pool)
	projectRepo := project.NewPostgresRepository(pool)
	taskRepo := task.NewPostgresRepository(pool)
	adminRepo := admin.NewPostgresRepository(pool)
	settingRepo := setting.NewPostgresRepository(pool)
	contactRepo := contact.NewPostgresRepository(pool)
	aimodelRepo := aimodel.NewPostgresRepository(pool)
	auditRepo := audit.NewPostgresRepository(pool)
	reportRepo := report.NewPostgresRepository(pool)
	notifRepo := notification.NewPostgresRepository(pool)
	compRepo := competitor.NewPostgresRepository(pool)
	planRepo := plan.NewPostgresRepository(pool)
	orderRepo := order.NewPostgresRepository(pool)
	paymentRepo := payment.NewPostgresRepository(pool)
	cdkRepo := cdk.NewPostgresRepository(pool)
	requestLogRepo := requestlog.NewPostgresRepository(pool)
	spaceRepo := space.NewPostgresRepository(pool)

	// Create router
	r := router.New(cfg, userRepo, projectRepo, taskRepo, adminRepo, settingRepo, contactRepo, aimodelRepo, auditRepo, reportRepo, notifRepo, compRepo, planRepo, orderRepo, paymentRepo, cdkRepo, requestLogRepo, spaceRepo)

	// Start HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("server starting on port %s", cfg.ServerPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
	log.Println("server stopped")
}

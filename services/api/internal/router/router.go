package router

import (
	"github.com/aige/admin"
	"github.com/aige/aimodel"
	"github.com/aige/api/internal/config"
	"github.com/aige/api/internal/handler"
	"github.com/aige/api/internal/middleware"
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
	"github.com/go-chi/chi/v5"
)

// New creates and configures the chi router with all routes.
func New(cfg *config.Config, userRepo user.Repository, projectRepo project.Repository, taskRepo task.Repository, adminRepo admin.Repository, settingRepo setting.Repository, contactRepo contact.Repository, aimodelRepo aimodel.Repository, auditRepo audit.Repository, reportRepo report.Repository, notifRepo notification.Repository, compRepo competitor.Repository, planRepo plan.Repository, orderRepo order.Repository, paymentRepo payment.Repository, cdkRepo cdk.Repository, requestLogRepo requestlog.Repository, spaceRepo space.Repository) chi.Router {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger(requestLogRepo))

	// Handlers
	aiProvider := aimodel.NewOpenAIProvider()
	authH := &handler.AuthHandler{UserRepo: userRepo, SettingRepo: settingRepo, JWTSecret: cfg.JWTSecret}
	userH := &handler.UserHandler{UserRepo: userRepo, JWTSecret: cfg.JWTSecret}
	projectH := &handler.ProjectHandler{ProjectRepo: projectRepo, ModelRepo: aimodelRepo, AIProvider: aiProvider, UserRepo: userRepo}
	taskH := &handler.TaskHandler{TaskRepo: taskRepo, ProjectRepo: projectRepo, ModelRepo: aimodelRepo, AIProvider: aiProvider, UserRepo: userRepo, CompRepo: compRepo, ReportRepo: reportRepo}
	dashboardH := &handler.DashboardHandler{TaskRepo: taskRepo, UserRepo: userRepo, ProjectRepo: projectRepo}
	adminAuthH := &handler.AdminAuthHandler{AdminRepo: adminRepo, JWTSecret: cfg.JWTSecret}
	adminH := &handler.AdminHandler{AdminRepo: adminRepo}
	adminDashboardH := &handler.AdminDashboardHandler{Pool: cfg.Pool, UserRepo: userRepo, ProjectRepo: projectRepo, TaskRepo: taskRepo}
	settingH := &handler.SettingHandler{SettingRepo: settingRepo}
	contactH := &handler.ContactHandler{ContactRepo: contactRepo}
	aimodelH := &handler.AIModelHandler{ModelRepo: aimodelRepo, AuditRepo: auditRepo, Pool: cfg.Pool}
	auditH := &handler.AuditLogHandler{AuditRepo: auditRepo}
	reportH := &handler.ReportHandler{ReportRepo: reportRepo}
	notifH := &handler.NotificationHandler{NotifRepo: notifRepo}
	compH := &handler.CompetitorHandler{CompRepo: compRepo, ProjectRepo: projectRepo}
	adminUserH := &handler.AdminUserHandler{Pool: cfg.Pool, AuditRepo: auditRepo}
	adminTaskH := &handler.AdminTaskHandler{TaskRepo: taskRepo}
	adminPromptH := &handler.AdminPromptHandler{Pool: cfg.Pool}
	adminPlanH := &handler.AdminPlanHandler{PlanRepo: planRepo}
	adminOrderH := &handler.AdminOrderHandler{OrderRepo: orderRepo}
	adminPaymentH := &handler.AdminPaymentHandler{PaymentRepo: paymentRepo}
	adminCDKH := &handler.AdminCDKHandler{CDKRepo: cdkRepo}
	adminRequestLogH := handler.NewAdminRequestLogHandler(requestLogRepo)
	spaceH := &handler.SpaceHandler{SpaceRepo: spaceRepo}

	billingH := &handler.BillingHandler{
		PlanRepo:    planRepo,
		OrderRepo:   orderRepo,
		UserRepo:    userRepo,
		CDKRepo:     cdkRepo,
		PaymentRepo: paymentRepo,
	}

	// Public settings routes
	r.Get("/api/v1/settings", settingH.Get)

	// Public contact routes
	r.Post("/api/v1/contact", contactH.Create)

	// Public billing routes (no auth - payment callback)
	r.Post("/api/v1/billing/notify", billingH.PaymentNotify)

	// Auth routes (no JWT required)
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authH.Register)
		r.Post("/login", authH.Login)
		r.Post("/send-code", authH.SendCode)
		r.Post("/reset-password", authH.ResetPassword)
	})

	// Admin auth routes (no admin JWT required)
	r.Post("/api/v1/admin/auth/login", adminAuthH.Login)

	// Protected routes (JWT required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))

		r.Route("/api/v1/users", func(r chi.Router) {
			r.Get("/me", userH.GetMe)
			r.Put("/me", userH.UpdateMe)
			r.Put("/me/password", userH.ChangePassword)
		})

		r.Route("/api/v1/projects", func(r chi.Router) {
			r.Post("/", projectH.Create)
			r.Get("/", projectH.List)
			r.Get("/{id}", projectH.Get)
			r.Put("/{id}", projectH.Update)
			r.Delete("/{id}", projectH.Delete)
		})

		r.Route("/api/v1/tasks", func(r chi.Router) {
			r.Post("/", taskH.Create)
			r.Get("/", taskH.List)
			r.Get("/{id}", taskH.Get)
			r.Get("/{id}/report", taskH.GetReport)
			r.Get("/{id}/comparison", taskH.GetComparisonReport)
			r.Post("/{id}/execute", taskH.Execute)
		})

		r.Get("/api/v1/models", aimodelH.ListEnabled)
		r.Get("/api/v1/dashboard/stats", dashboardH.GetStats)

		r.Route("/api/v1/spaces", func(r chi.Router) {
			r.Get("/", spaceH.List)
			r.Post("/", spaceH.Create)
			r.Get("/current", spaceH.Current)
			r.Put("/{id}/current", spaceH.SetCurrent)
			r.Get("/{id}/members", spaceH.Members)
			r.Post("/{id}/members/invite", spaceH.Invite)
			r.Patch("/{id}/members/{userID}", spaceH.UpdateMemberRole)
			r.Delete("/{id}/members/{userID}", spaceH.RemoveMember)
		})

		r.Route("/api/v1/reports", func(r chi.Router) {
			r.Get("/", reportH.List)
			r.Get("/{id}", reportH.GetByID)
		})

		r.Get("/api/v1/projects/{projectId}/reports", reportH.ListByProject)
		r.Get("/api/v1/projects/{projectId}/competitors", compH.ListByProject)

		// Billing routes (JWT required)
		r.Get("/api/v1/plans", billingH.GetPlans)
		r.Route("/api/v1/billing", func(r chi.Router) {
			r.Post("/orders", billingH.CreateOrder)
			r.Get("/orders", billingH.GetOrders)
			r.Post("/cdk/redeem", billingH.RedeemCDK)
			r.Get("/credits", billingH.GetCredits)
		})

		r.Route("/api/v1/notifications", func(r chi.Router) {
			r.Get("/", notifH.List)
			r.Put("/{id}/read", notifH.MarkRead)
			r.Put("/read-all", notifH.MarkAllRead)
			r.Get("/unread-count", notifH.CountUnread)
		})
	})

	// Admin protected routes (admin JWT required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AdminJWTAuth(cfg.JWTSecret))

		r.Get("/api/v1/admin/auth/me", adminAuthH.Me)

		r.Route("/api/v1/admin/admins", func(r chi.Router) {
			r.Get("/", adminH.List)
			r.Post("/", adminH.Create)
			r.Get("/{id}", adminH.Get)
			r.Put("/{id}", adminH.Update)
			r.Delete("/{id}", adminH.Delete)
		})

		r.Get("/api/v1/admin/dashboard/stats", adminDashboardH.GetStats)

		r.Get("/api/v1/admin/settings", settingH.Get)
		r.Put("/api/v1/admin/settings", settingH.Update)

		r.Route("/api/v1/admin/contact", func(r chi.Router) {
			r.Get("/", contactH.List)
			r.Put("/{id}/read", contactH.MarkRead)
			r.Delete("/{id}", contactH.Delete)
		})

		r.Route("/api/v1/admin/models", func(r chi.Router) {
			r.Get("/", aimodelH.List)
			r.Post("/", aimodelH.Create)
			r.Post("/discover", aimodelH.Discover)
			r.Get("/stats", aimodelH.GetStats)
			r.Get("/{id}", aimodelH.Get)
			r.Put("/{id}", aimodelH.Update)
			r.Delete("/{id}", aimodelH.Delete)
			r.Post("/{id}/test", aimodelH.TestConnection)
		})

		r.Get("/api/v1/admin/logs", auditH.List)
		r.Get("/api/v1/admin/request-logs", adminRequestLogH.List)

		// Admin user management
		r.Get("/api/v1/admin/users", adminUserH.ListUsers)
		r.Get("/api/v1/admin/users/{id}", adminUserH.GetUser)
		r.Get("/api/v1/admin/users/{id}/projects", adminUserH.GetUserProjects)
		r.Put("/api/v1/admin/users/{id}", adminUserH.UpdateUser)

		// Admin task management
		r.Get("/api/v1/admin/tasks", adminTaskH.ListTasks)
		r.Post("/api/v1/admin/tasks/{id}/retry", adminTaskH.RetryTask)
		r.Delete("/api/v1/admin/tasks/{id}", adminTaskH.DeleteTask)

		// Admin prompt management
		r.Get("/api/v1/admin/prompts", adminPromptH.ListPrompts)
		r.Post("/api/v1/admin/prompts", adminPromptH.CreatePrompt)
		r.Get("/api/v1/admin/prompts/{id}", adminPromptH.GetPrompt)
		r.Put("/api/v1/admin/prompts/{id}", adminPromptH.UpdatePrompt)
		r.Post("/api/v1/admin/prompts/{id}/publish", adminPromptH.PublishPrompt)

		// Admin plan management
		r.Get("/api/v1/admin/plans", adminPlanH.ListPlans)
		r.Post("/api/v1/admin/plans", adminPlanH.CreatePlan)
		r.Put("/api/v1/admin/plans/{id}", adminPlanH.UpdatePlan)

		// Admin order management
		r.Get("/api/v1/admin/orders", adminOrderH.ListOrders)

		// Admin payment config
		r.Get("/api/v1/admin/payment/config", adminPaymentH.GetConfig)
		r.Put("/api/v1/admin/payment/config", adminPaymentH.UpdateConfig)

		// Admin CDK management
		r.Route("/api/v1/admin/cdk", func(r chi.Router) {
			r.Get("/", adminCDKH.ListCDK)
			r.Post("/", adminCDKH.CreateCDK)
			r.Put("/{id}", adminCDKH.UpdateCDK)
			r.Get("/{id}/usages", adminCDKH.GetCDKUsages)
		})
	})

	return r
}

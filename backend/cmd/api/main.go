package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ggampp/ordens_servico/backend/internal/auth"
	"github.com/ggampp/ordens_servico/backend/internal/config"
	"github.com/ggampp/ordens_servico/backend/internal/database"
	"github.com/ggampp/ordens_servico/backend/internal/handler"
	"github.com/ggampp/ordens_servico/backend/internal/repository"
	"github.com/ggampp/ordens_servico/backend/internal/service"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	ctx := context.Background()

	// Database connection + migrations (PostgreSQL + PostGIS via DATABASE_URL).
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	// Wire dependencies (repository -> service -> handler).
	jwtManager := auth.NewManager(cfg.JWTSecret, cfg.JWTExpiry)

	employeeRepo := repository.NewEmployeeRepository(pool)
	orderRepo := repository.NewServiceOrderRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	dashboardRepo := repository.NewDashboardRepository(pool)
	mapRepo := repository.NewMapRepository(pool)

	authSvc := service.NewAuthService(userRepo, jwtManager)
	if err := authSvc.SeedAdmin(ctx, cfg.SeedAdminEmail, cfg.SeedAdminPass); err != nil {
		slog.Error("seed admin failed", "error", err)
		os.Exit(1)
	}

	handlers := handler.Handlers{
		Auth:      handler.NewAuthHandler(authSvc),
		Employee:  handler.NewEmployeeHandler(service.NewEmployeeService(employeeRepo)),
		Order:     handler.NewServiceOrderHandler(service.NewServiceOrderService(orderRepo, employeeRepo)),
		Map:       handler.NewMapHandler(service.NewMapService(orderRepo, mapRepo)),
		Dashboard: handler.NewDashboardHandler(service.NewDashboardService(dashboardRepo)),
	}

	router := handler.NewRouter(handlers, jwtManager, cfg.StaticDir)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Graceful shutdown.
	go func() {
		slog.Info("server listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

// setupLogger configures the global structured JSON logger.
func setupLogger(level string) {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	slog.SetDefault(logger)
}

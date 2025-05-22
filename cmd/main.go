package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/config"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/database"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/handler"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/middleware"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/repository"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/routes"
	"gitlab.com/timkado/api/daisi-rest-postgres/internal/service"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load config via Viper
	cfg := config.LoadConfig()

	// Initialize Zap logger (ISO8601 timestamps, "timestamp" key)
	log := logger.NewLogger()
	defer log.Sync()

	// Connect to Postgres using the full DSN from env (PG_DSN)
	if err := database.ConnectGORM(cfg.PgDsn, log); err != nil {
		log.Fatal("Cannot initialize database", zap.Error(err))
	}

	// Repo + Service Registration
	agentRepo := repository.NewAgentRepository()
	chatRepo := repository.NewChatRepository()
	messageRepo := repository.NewMessageRepository()
	agentSvc := service.NewAgentService(agentRepo)
	chatSvc := service.NewChatService(chatRepo)
	messageSvc := service.NewMessageService(messageRepo)

	handler.RegisterAgentService(agentSvc)
	handler.RegisterChatService(chatSvc)
	handler.RegisterMessageService(messageSvc)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Daisi REST Postgres API",
	})

	// Global middleware
	app.Use(middleware.Recover())     // panic recovery
	app.Use(middleware.RequestID())   // inject X-Request-ID
	app.Use(middleware.Helmet())      // security headers
	app.Use(middleware.CORS())        // CORS support
	app.Use(middleware.Idempotency()) // idempotent POST/PUT/DELETE

	// Health endpoints
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/readyz", func(c *fiber.Ctx) error {
		sqlDB, err := database.DB.DB()
		if err != nil {
			return c.
				Status(fiber.StatusServiceUnavailable).
				JSON(fiber.Map{"status": "db unavailable", "error": err.Error()})
		}
		if err := sqlDB.PingContext(context.Background()); err != nil {
			return c.
				Status(fiber.StatusServiceUnavailable).
				JSON(fiber.Map{"status": "db unavailable", "error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// Register all /api/v1 routes
	routes.RegisterV1Routes(app)

	// Start server in goroutine
	go func() {
		log.Info("Listening on port " + cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatal("Failed to listen", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info("Shutting down server...")

	// Give active requests up to 5s to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use ShutdownWithContext to respect our timeout
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("Error during shutdown", zap.Error(err))
	}
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	in_memory "github.com/kasasunil/chat_app/database/in-memory"

	"github.com/kasasunil/chat_app/bootstrap"
	"github.com/kasasunil/chat_app/config"
	"github.com/kasasunil/chat_app/controller"
	"github.com/kasasunil/chat_app/internal/middleware"
	"github.com/kasasunil/chat_app/internal/pkg/logger"
	"github.com/kasasunil/chat_app/internal/services/websocket"
)

func main() {
	configPath := os.Getenv(config.EnvConfigPath)
	if configPath == "" {
		configPath = config.DefaultConfigPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// Initialize logger with defaults first
		logger.InitLogger("info", "")
		logger.Warn(logger.TraceConfigLoadFailed, configPath, err)
		cfg = config.NewConfig()
	}

	// Initialize logger with config
	if err := logger.InitLogger(cfg.Logging.Level, ""); err != nil {
		// Fallback to default logger if initialization fails
		logger.Warn(logger.TraceLoggerInitFailed, err)
	}

	store := in_memory.NewStore() // For product use actual db instance.
	wsManager := websocket.NewMockWebSocketManager()
	handler := controller.NewHandler(store, wsManager)

	// Setup demo data
	bootstrap.SetupDemoData(store, wsManager)

	// Initialize authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Setup routes
	router := bootstrap.SetupRouter(handler, authMiddleware)

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info(logger.TraceServerStarting, cfg.Server.Host, cfg.Server.Port)
	logger.Info("API Endpoints (all require authentication):")
	logger.Info("  POST   /api/v1/sendMessage")
	logger.Info("  POST   /api/v1/ack/delivered")
	logger.Info("  POST   /api/v1/ack/read")
	logger.Info("  GET    /api/v1/conversations/{destinationId}/messages")
	logger.Info("  GET    /api/v1/users/{userId}/conversations")
	logger.Info("  GET    /api/v1/search/{userId}?query=xxx")
	//logger.Info("Authentication: Basic Auth with credentials from conf/config.toml")
	logger.Info("Run the demo test to see the system in action!")

	// Start server with graceful shutdown
	bootstrap.StartServerWithGracefulShutdown(server)
}

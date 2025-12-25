package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kasasunil/chat_app/controller"
	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/middleware"
	"github.com/kasasunil/chat_app/internal/pkg/logger"
	"github.com/kasasunil/chat_app/internal/services/websocket"

	"github.com/gorilla/mux"
)

// SetupDemoData creates sample users, groups, and simulates connections
// Accepts interface for store (following Go best practices)
func SetupDemoData(store database.Repository, wsManager websocket.WebSocketManager) {
	// Create users
	user1 := &database.User{
		ID:    "user1",
		Name:  "Alice",
		Email: "alice@example.com",
	}
	user2 := &database.User{
		ID:    "user2",
		Name:  "Bob",
		Email: "bob@example.com",
	}
	user3 := &database.User{
		ID:    "user3",
		Name:  "Charlie",
		Email: "charlie@example.com",
	}

	store.CreateUser(user1)
	store.CreateUser(user2)
	store.CreateUser(user3)

	// Create group
	group1 := &database.Group{
		ID:          "group1",
		Name:        "Project Team",
		Description: "Team chat for project discussions",
		CreatedBy:   "user1",
	}
	store.CreateGroup(group1)
	store.AddGroupMember("group1", "user1")
	store.AddGroupMember("group1", "user2")
	store.AddGroupMember("group1", "user3")

	// Simulate active connections
	wsManager.AddConnection("user1", "conn1")
	wsManager.AddConnection("user2", "conn2")
	wsManager.AddConnection("user3", "conn3")

	logger.Info(logger.TraceDemoDataInitialized)
	logger.Info("  Users: Alice (user1), Bob (user2), Charlie (user3)")
	logger.Info("  Group: Project Team (group1)")
	logger.Info("  All users are connected")
}

// SetupRouter initializes and configures all routes
func SetupRouter(handler *controller.Handler, authMiddleware *middleware.AuthMiddleware) *mux.Router {
	router := mux.NewRouter()

	// Public routes (no authentication required)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Protected routes (authentication required)
	// API v1 routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(authMiddleware.Authenticate)

	apiRouter.HandleFunc("/sendMessage", handler.SendMessage).Methods("POST")
	apiRouter.HandleFunc("/ack/delivered", handler.AckDelivered).Methods("POST")
	apiRouter.HandleFunc("/ack/read", handler.AckRead).Methods("POST")
	apiRouter.HandleFunc("/conversations/{destinationId}/messages", handler.GetMessages).Methods("GET")
	apiRouter.HandleFunc("/users/{userId}/conversations", handler.GetUserConversations).Methods("GET")
	apiRouter.HandleFunc("/search/{userId}", handler.SearchMessages).Methods("GET")
	return router
}

// StartServerWithGracefulShutdown starts the server and handles graceful shutdown
func StartServerWithGracefulShutdown(server *http.Server) {
	// Channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	wg := &sync.WaitGroup{}

	// Goroutine 1: Start the server
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Server is starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(logger.TraceServerFailed, err)
		}
	}()

	// Goroutine 2: Handle graceful shutdown (waits for OS signal)
	// This goroutine blocks on <-sigChan until a signal is received
	// It will only execute shutdown logic ONCE when signal is received
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Wait for interrupt signal from OS (blocks here until signal is received)
		sig := <-sigChan
		logger.Info("Received signal: %v. Initiating graceful shutdown...", sig)

		// Create shutdown context with timeout, so that all the idle connections will destroy, for prod we can increase this time.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Shutdown the server (this will cause server goroutine's ListenAndServe to return)
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown: %v", err)
		} else {
			logger.Info("Server gracefully stopped")
		}
	}()

	logger.Info("Server started!!!")
	// Wait for both goroutines to complete
	// - Server goroutine completes when ListenAndServe returns (after Shutdown is called)
	// - Shutdown goroutine completes after handling the signal and shutting down the server
	wg.Wait()
	logger.Info("Server shutdown complete")
}

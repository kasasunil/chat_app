package controller

import (
	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/services/search"
	"github.com/kasasunil/chat_app/internal/services/websocket"
)

// Handler contains all HTTP handlers
type Handler struct {
	store         database.Repository
	wsManager     websocket.WebSocketManager
	searchService *search.SearchService
}

// NewHandler creates a new handler instance
func NewHandler(store database.Repository, wsManager websocket.WebSocketManager) *Handler {
	return &Handler{
		store:         store,
		wsManager:     wsManager,
		searchService: search.NewSearchService(store),
	}
}

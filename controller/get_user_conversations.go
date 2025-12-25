package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/middleware"
	"github.com/kasasunil/chat_app/internal/pkg/errors"

	"github.com/gorilla/mux"
)

// ConversationListItem represents a conversation in the list view
type ConversationListItem struct {
	ConversationID   string                    `json:"conversation_id"`
	DestinationID    string                    `json:"destination_id"`
	ConversationType database.ConversationType `json:"conversation_type"`
	LastMessage      *database.Message         `json:"last_message,omitempty"`
	UnreadCount      int                       `json:"unread_count"`
	UpdatedAt        time.Time                 `json:"updated_at"`
}

// GetUserConversationsResponse represents the response for fetching user conversations
type GetUserConversationsResponse struct {
	Conversations []ConversationListItem `json:"conversations"`
}

// GetUserConversations handles GET /users/{userId}/conversations
func (h *Handler) GetUserConversations(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	authenticatedUserID := middleware.GetUserID(r)
	if authenticatedUserID == "" {
		respondWithError(w, errors.ErrAuthRequired)
		return
	}

	vars := mux.Vars(r)
	requestedUserID := vars["userId"]

	// Validate that requested user_id matches authenticated user
	if requestedUserID != authenticatedUserID {
		respondWithError(w, errors.ErrForbiddenAccessDenied)
		return
	}

	userID := authenticatedUserID

	conversations, err := h.store.GetUserConversations(userID)
	if err != nil {
		respondWithError(w, errors.ErrInternalError)
		return
	}

	items := make([]ConversationListItem, 0, len(conversations))
	for _, conv := range conversations {
		// Get last message for this conversation
		messages, _, _ := h.store.GetMessages(conv.DestinationID, 1, "")
		var lastMessage *database.Message
		if len(messages) > 0 {
			lastMessage = messages[0]
		}

		// Calculate unread count (messages sent to user that haven't been read)
		// A message is unread if:
		// 1. User is not the sender
		// 2. There's no MessageRead entry for this user and message
		unreadCount := 0
		allMessages, _, _ := h.store.GetMessages(conv.DestinationID, 1000, "")
		for _, msg := range allMessages {
			// Only count messages not sent by the user
			if msg.SenderID != userID {
				reads := h.store.GetMessageReads(msg.ID)
				hasRead := false
				// Check if this user has read this message
				for _, read := range reads {
					if read.UserID == userID {
						hasRead = true
						break
					}
				}
				// If no read receipt exists for this user, it's unread
				if !hasRead {
					unreadCount++
				}
			}
		}

		items = append(items, ConversationListItem{
			ConversationID:   conv.ID,
			DestinationID:    conv.DestinationID,
			ConversationType: conv.ConversationType,
			LastMessage:      lastMessage,
			UnreadCount:      unreadCount,
			UpdatedAt:        conv.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(GetUserConversationsResponse{
		Conversations: items,
	})
}

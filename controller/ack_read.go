package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/middleware"
	"github.com/kasasunil/chat_app/internal/pkg/errors"
	"github.com/kasasunil/chat_app/internal/pkg/logger"
)

// AckReadRequest represents the request to acknowledge read
type AckReadRequest struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
}

// AckReadResponse represents the response after acknowledging read
type AckReadResponse struct {
	Message string `json:"message"`
}

// AckRead handles POST /ack/read
func (h *Handler) AckRead(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	authenticatedUserID := middleware.GetUserID(r)
	if authenticatedUserID == "" {
		respondWithError(w, errors.ErrAuthRequired)
		return
	}

	var req AckReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, errors.ErrInvalidRequest)
		return
	}

	// Use authenticated user (override if user_id provided for backward compatibility)
	userID := authenticatedUserID
	if req.UserID != "" {
		if req.UserID != authenticatedUserID {
			respondWithError(w, errors.ErrForbiddenAccessDenied)
			return
		}
		userID = req.UserID
	}

	message, err := h.store.GetMessage(req.MessageID)
	if err != nil {
		respondWithError(w, errors.ErrMessageNotFound)
		return
	}

	// Verify user is the recipient
	isRecipient := false
	if message.ConversationType == database.ConversationTypeOneToOne {
		isRecipient = message.DestinationID == userID
	} else {
		isRecipient = h.store.IsGroupMember(message.DestinationID, userID)
	}

	if !isRecipient {
		respondWithError(w, errors.ErrNotMessageRecipient)
		return
	}

	// Create read receipt
	if err := h.store.CreateMessageRead(req.MessageID, userID); err != nil {
		logger.Error(logger.TraceMessageReadFailed, req.MessageID, userID, err)
		respondWithError(w, errors.ErrInternalError)
		return
	}

	logger.Info(logger.TraceMessageRead, req.MessageID, userID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AckReadResponse{
		Message: "Message read",
	})
}

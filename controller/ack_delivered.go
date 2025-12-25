package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/middleware"
	"github.com/kasasunil/chat_app/internal/pkg/errors"
)

// AckDeliveredRequest represents the request to acknowledge delivery
type AckDeliveredRequest struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
}

// AckDeliveredResponse represents the response after acknowledging delivery
type AckDeliveredResponse struct {
	Message string `json:"message"`
}

// AckDelivered handles POST /ack/delivered
func (h *Handler) AckDelivered(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	authenticatedUserID := middleware.GetUserID(r)
	if authenticatedUserID == "" {
		respondWithError(w, errors.ErrAuthRequired)
		return
	}

	var req AckDeliveredRequest
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

	// Update message status to DELIVERED
	if message.Status == database.StatusSent {
		if err := h.store.UpdateMessageStatus(req.MessageID, database.StatusDelivered); err != nil {
			respondWithError(w, errors.ErrInternalError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AckDeliveredResponse{
		Message: "Message delivered",
	})
}

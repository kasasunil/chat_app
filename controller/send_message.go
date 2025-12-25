package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/middleware"
	"github.com/kasasunil/chat_app/internal/pkg/errors"
	"github.com/kasasunil/chat_app/internal/pkg/logger"
	"github.com/kasasunil/chat_app/internal/pkg/utils"
)

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	SenderID      string `json:"sender_id"`
	DestinationID string `json:"destination_id"`
	Message       string `json:"message"`
}

// SendMessageResponse represents the response after sending a message
type SendMessageResponse struct {
	MessageID string                 `json:"message_id"`
	Status    database.MessageStatus `json:"status"`
}

// SendMessage handles POST /sendMessage
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	authenticatedUserID := middleware.GetUserID(r)
	if authenticatedUserID == "" {
		logger.Warn("Unauthorized request to SendMessage")
		respondWithError(w, errors.ErrAuthRequired)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid request body in SendMessage: %v", err)
		respondWithError(w, errors.ErrInvalidRequest)
		return
	}

	if utils.IsEmpty(req.Message) {
		logger.Warn(logger.TraceValidationFailed, FieldMessage, "empty")
		respondWithError(w, errors.ErrMessageEmpty)
		return
	}

	// Use authenticated user as sender (override if sender_id provided for backward compatibility)
	senderID := authenticatedUserID
	if req.SenderID != "" {
		// Validate that provided sender_id matches authenticated user
		if req.SenderID != authenticatedUserID {
			respondWithError(w, errors.ErrForbiddenAccessDenied)
			return
		}
		senderID = req.SenderID
	}

	// Validate sender exists
	_, err := h.store.GetUser(senderID)
	if err != nil {
		respondWithError(w, errors.ErrSenderNotFound)
		return
	}

	// Determine conversation type
	convType := database.ConversationTypeOneToOne
	_, err = h.store.GetGroup(req.DestinationID)
	if err == nil {
		// Destination is a group
		convType = database.ConversationTypeGroup
		// Verify sender is a member
		if !h.store.IsGroupMember(req.DestinationID, senderID) {
			respondWithError(w, errors.ErrNotGroupMember)
			return
		}
	} else {
		// Verify destination user exists
		_, err = h.store.GetUser(req.DestinationID)
		if err != nil {
			respondWithError(w, errors.ErrDestinationNotFound)
			return
		}
	}

	// Create message
	message := &database.Message{
		ID:               utils.GenerateID(),
		SenderID:         senderID,
		DestinationID:    req.DestinationID,
		MessageText:      utils.SanitizeString(req.Message, 10000),
		Status:           database.StatusSent,
		ConversationType: convType,
	}

	logger.Info(logger.TraceMessageSent, senderID, req.DestinationID, convType, message.ID)

	if err := h.store.CreateMessage(message); err != nil {
		respondWithError(w, errors.ErrInternalError)
		return
	}

	// Simulate sending message via WebSocket
	if convType == database.ConversationTypeOneToOne {
		h.wsManager.SendMessage(req.DestinationID, message)
	} else {
		// For group messages, send to all members
		group, _ := h.store.GetGroup(req.DestinationID)
		if group != nil {
			// In real implementation, would iterate through group members
			// For now, we just simulate the send
			h.wsManager.SendMessage(req.DestinationID, message)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created - new resource created
	json.NewEncoder(w).Encode(SendMessageResponse{
		MessageID: message.ID,
		Status:    message.Status,
	})
}

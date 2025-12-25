package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/pkg/errors"

	"github.com/gorilla/mux"
)

// GetMessagesResponse represents the response for fetching messages
type GetMessagesResponse struct {
	Messages   []*database.Message `json:"messages"`
	NextCursor string              `json:"next_cursor,omitempty"`
	HasMore    bool                `json:"has_more"`
}

// GetMessages handles GET /conversations/{destinationId}/messages
// This route fetches conversations of a single one-one messages / grp messages.
func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	destinationID := vars["destinationId"]

	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	messages, nextCursor, err := h.store.GetMessages(destinationID, limit, cursor)
	if err != nil {
		respondWithError(w, errors.ErrInternalError)
		return
	}

	hasMore := nextCursor != ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(GetMessagesResponse{
		Messages:   messages,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	})
}

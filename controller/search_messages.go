package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/middleware"
	"github.com/kasasunil/chat_app/internal/pkg/errors"
	"github.com/kasasunil/chat_app/internal/pkg/utils"

	"github.com/gorilla/mux"
)

// SearchMessagesResponse represents the response for search
type SearchMessagesResponse struct {
	Results []*database.Message `json:"results"`
	Query   string              `json:"query"`
}

// SearchMessages handles GET /search/{userId}?query=xxx
func (h *Handler) SearchMessages(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query().Get(FieldQuery)
	if utils.IsEmpty(query) {
		respondWithError(w, errors.ErrSearchQueryRequired)
		return
	}

	results, err := h.searchService.SearchMessages(userID, query)
	if err != nil {
		respondWithError(w, errors.ErrSearchFailed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(SearchMessagesResponse{
		Results: results,
		Query:   query,
	})
}

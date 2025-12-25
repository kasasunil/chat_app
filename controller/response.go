package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kasasunil/chat_app/internal/pkg/errors"
)

// respondWithError sends an error response using the error package
func respondWithError(w http.ResponseWriter, err *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus)
	json.NewEncoder(w).Encode(err.ToJSONResponse())
}

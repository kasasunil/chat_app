package search

import "github.com/kasasunil/chat_app/database"

// SearchService provides search functionality
type SearchService struct {
	store database.Repository
}

// NewSearchService creates a new search service
// Accepts interface, returns struct (following Go best practices)
func NewSearchService(store database.Repository) *SearchService {
	return &SearchService{
		store: store,
	}
}

// SearchMessages performs keyword search across messages
func (s *SearchService) SearchMessages(userID, query string) ([]*database.Message, error) {
	return s.store.SearchMessages(userID, query)
}

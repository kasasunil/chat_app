package in_memory

import (
	"fmt"
	"github.com/kasasunil/chat_app/database"
	"time"
)

// User operations
func (s *MemoryStore) CreateUser(user *database.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStore) GetUser(userID string) (*database.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

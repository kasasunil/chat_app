package in_memory

import (
	"fmt"
	"github.com/kasasunil/chat_app/database"
	"time"
)

// Group operations
func (s *MemoryStore) CreateGroup(group *database.Group) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.groups[group.ID]; exists {
		return fmt.Errorf("group already exists")
	}

	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	s.groups[group.ID] = group
	s.groupMembers[group.ID] = make(map[string]bool)
	return nil
}

func (s *MemoryStore) GetGroup(groupID string) (*database.Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	group, exists := s.groups[groupID]
	if !exists {
		return nil, fmt.Errorf("group not found")
	}
	return group, nil
}

func (s *MemoryStore) AddGroupMember(groupID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.groups[groupID]; !exists {
		return fmt.Errorf("group not found")
	}

	if s.groupMembers[groupID] == nil {
		s.groupMembers[groupID] = make(map[string]bool)
	}
	s.groupMembers[groupID][userID] = true
	return nil
}

func (s *MemoryStore) IsGroupMember(groupID, userID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.groupMembers[groupID][userID]
}

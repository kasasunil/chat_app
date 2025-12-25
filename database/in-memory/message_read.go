package in_memory

import (
	"fmt"
	"time"

	"github.com/kasasunil/chat_app/database"
	"github.com/kasasunil/chat_app/internal/pkg/utils"
)

// MessageRead operations
func (s *MemoryStore) CreateMessageRead(messageID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.messageReads[messageID] == nil {
		s.messageReads[messageID] = make(map[string]*database.MessageRead)
	}

	if _, exists := s.messageReads[messageID][userID]; exists {
		return nil // Already read
	}

	mr := &database.MessageRead{
		ID:        fmt.Sprintf("mr_%s_%s", messageID, userID),
		MessageID: messageID,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.messageReads[messageID][userID] = mr

	// Update message status to READ
	// Directly access messages map since we already have the write lock
	// (Calling GetMessage would cause deadlock as it tries to acquire read lock)
	found := false
	for _, messages := range s.messages {
		for _, msg := range messages {
			if msg.ID == messageID {
				msg.Status = database.StatusRead
				msg.UpdatedAt = time.Now()
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	return nil
}

func (s *MemoryStore) GetMessageReads(messageID string) []*database.MessageRead {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reads, exists := s.messageReads[messageID]
	if !exists {
		return []*database.MessageRead{}
	}

	result := make([]*database.MessageRead, 0, len(reads))
	for _, mr := range reads {
		result = append(result, mr)
	}
	return result
}

// UserConversation operations
func (s *MemoryStore) GetUserConversations(userID string) ([]*database.UserConversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conversations, exists := s.userConversations[userID]
	if !exists {
		return []*database.UserConversation{}, nil
	}

	// Sort by updated_at desc
	result := make([]*database.UserConversation, len(conversations))
	copy(result, conversations)

	// Simple sort by updated_at (newest first)
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].UpdatedAt.Before(result[j].UpdatedAt) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// Search operations
func (s *MemoryStore) SearchMessages(userID, query string) ([]*database.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]*database.Message, 0)

	// Search in all messages where user is sender or recipient
	for _, messages := range s.messages {
		for _, msg := range messages {
			// Check if user is part of this conversation
			isParticipant := false
			if msg.SenderID == userID {
				isParticipant = true
			} else if msg.ConversationType == database.ConversationTypeOneToOne && msg.DestinationID == userID {
				isParticipant = true
			} else if msg.ConversationType == database.ConversationTypeGroup {
				if s.groupMembers[msg.DestinationID][userID] {
					isParticipant = true
				}
			}

			if isParticipant {
				// Simple keyword search (case-insensitive)
				if utils.ContainsString(msg.MessageText, query) {
					results = append(results, msg)
				}
			}
		}
	}

	return results, nil
}

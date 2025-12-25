package in_memory

import (
	"fmt"
	"time"

	"github.com/kasasunil/chat_app/database"
)

// Message operations
func (s *MemoryStore) CreateMessage(message *database.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	message.Status = database.StatusSent

	destKey := message.DestinationID
	if s.messages[destKey] == nil {
		s.messages[destKey] = make([]*database.Message, 0)
	}
	s.messages[destKey] = append(s.messages[destKey], message)

	// Update user conversations for sender
	s.updateUserConversation(message.SenderID, message.DestinationID, message.ConversationType, message)

	// Update user conversations for recipient(s)
	if message.ConversationType == database.ConversationTypeOneToOne {
		s.updateUserConversation(message.DestinationID, message.SenderID, message.ConversationType, message)
	} else {
		// For group messages, update all group members
		if members, exists := s.groupMembers[message.DestinationID]; exists {
			for memberID := range members {
				if memberID != message.SenderID {
					s.updateUserConversation(memberID, message.DestinationID, message.ConversationType, message)
				}
			}
		}
	}

	return nil
}

func (s *MemoryStore) updateUserConversation(userID, destinationID string, convType database.ConversationType, message *database.Message) {
	if s.userConversations[userID] == nil {
		s.userConversations[userID] = make([]*database.UserConversation, 0)
	}

	// Check if conversation already exists
	var found bool
	for _, uc := range s.userConversations[userID] {
		if uc.DestinationID == destinationID && uc.ConversationType == convType {
			uc.UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		uc := &database.UserConversation{
			ID:               fmt.Sprintf("uc_%s_%s", userID, destinationID),
			UserID:           userID,
			DestinationID:    destinationID,
			ConversationType: convType,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		s.userConversations[userID] = append(s.userConversations[userID], uc)
	}
}

func (s *MemoryStore) GetMessages(destinationID string, limit int, cursor string) ([]*database.Message, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages, exists := s.messages[destinationID]
	if !exists {
		return []*database.Message{}, "", nil
	}

	// Reverse to get newest first, then apply cursor
	reversed := make([]*database.Message, len(messages))
	for i := range messages {
		reversed[len(messages)-1-i] = messages[i]
	}

	// Apply cursor if provided
	startIdx := 0
	if cursor != "" {
		for i, msg := range reversed {
			if msg.ID == cursor {
				startIdx = i + 1
				break
			}
		}
	}

	// Apply limit
	endIdx := startIdx + limit
	if endIdx > len(reversed) {
		endIdx = len(reversed)
	}

	result := reversed[startIdx:endIdx]

	// Get next cursor (first message of next page, if exists)
	nextCursor := ""
	if endIdx < len(reversed) {
		nextCursor = reversed[endIdx].ID
	}

	return result, nextCursor, nil
}

func (s *MemoryStore) GetMessage(messageID string) (*database.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, messages := range s.messages {
		for _, msg := range messages {
			if msg.ID == messageID {
				return msg, nil
			}
		}
	}
	return nil, fmt.Errorf("message not found")
}

func (s *MemoryStore) UpdateMessageStatus(messageID string, status database.MessageStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, messages := range s.messages {
		for _, msg := range messages {
			if msg.ID == messageID {
				msg.Status = status
				msg.UpdatedAt = time.Now()
				return nil
			}
		}
	}
	return fmt.Errorf("message not found")
}

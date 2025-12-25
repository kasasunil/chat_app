package in_memory

import (
	"sync"

	"github.com/kasasunil/chat_app/database"
)

// MemoryStore is an in-memory implementation of the Repository interface
type MemoryStore struct {
	mu                sync.RWMutex
	users             map[string]*database.User
	groups            map[string]*database.Group
	groupMembers      map[string]map[string]bool                  // groupID -> userID -> bool
	messages          map[string][]*database.Message              // destinationID -> messages
	userConversations map[string][]*database.UserConversation     // userID -> conversations
	messageReads      map[string]map[string]*database.MessageRead // messageID -> userID -> MessageRead
}

// NewStore creates a new in-memory store
// Returns struct (following "accept interfaces, return structs" principle)
func NewStore() *MemoryStore {
	return &MemoryStore{
		users:             make(map[string]*database.User),
		groups:            make(map[string]*database.Group),
		groupMembers:      make(map[string]map[string]bool),
		messages:          make(map[string][]*database.Message),
		userConversations: make(map[string][]*database.UserConversation),
		messageReads:      make(map[string]map[string]*database.MessageRead),
	}
}

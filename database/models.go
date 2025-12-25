package database

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Group represents a group in the system
type Group struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConversationType represents the type of conversation
type ConversationType string

const (
	ConversationTypeOneToOne ConversationType = "one-one"
	ConversationTypeGroup    ConversationType = "group"
)

// MessageStatus represents the status of a message
type MessageStatus string

const (
	StatusSent     MessageStatus = "SENT"     // Single tick (✓)
	StatusDelivered MessageStatus = "DELIVERED" // Double tick (✓✓)
	StatusRead     MessageStatus = "READ"     // Blue double tick (✓✓)
)

// UserConversation represents a user's view of a conversation
// Used to efficiently build the chat list screen
type UserConversation struct {
	ID               string           `json:"id"`
	UserID           string           `json:"user_id"`
	DestinationID    string           `json:"destination_id"`
	ConversationType ConversationType `json:"conversation_type"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

// Message represents a message in the system
type Message struct {
	ID               string           `json:"id"`
	SenderID         string           `json:"sender_id"`
	DestinationID    string           `json:"destination_id"`
	MessageText      string           `json:"message_text"`
	Status           MessageStatus    `json:"status"`
	ConversationType ConversationType `json:"conversation_type"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

// MessageRead represents a read receipt for a message
// Stores viewers of each conversation
type MessageRead struct {
	ID        string    `json:"id"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


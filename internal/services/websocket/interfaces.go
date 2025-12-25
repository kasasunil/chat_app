package websocket

import "github.com/kasasunil/chat_app/database"

// WebSocketManager defines the interface for WebSocket operations
// This is mocked - real WebSocket implementation is not required
type WebSocketManager interface {
	// SendMessage sends a message to a user's active connections
	SendMessage(userID string, message *database.Message) error

	// AckDelivered acknowledges that a message was delivered to a user's inbox
	AckDelivered(userID string, messageID string) error

	// AckRead acknowledges that a user has read a message
	AckRead(userID string, messageID string) error

	// AddConnection adds a connection for a user (simulating multiple devices)
	AddConnection(userID string, connectionID string)

	// RemoveConnection removes a connection for a user
	RemoveConnection(userID string, connectionID string)

	// IsUserConnected checks if a user has any active connections
	IsUserConnected(userID string) bool
}

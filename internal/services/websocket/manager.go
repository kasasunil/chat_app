package websocket

import "github.com/kasasunil/chat_app/database"

// MockWebSocketManager is an in-memory implementation of WebSocketManager
type MockWebSocketManager struct {
	connections map[string]map[string]bool // userID -> connectionID -> bool
}

// NewMockWebSocketManager creates a new mock WebSocket manager
func NewMockWebSocketManager() *MockWebSocketManager {
	return &MockWebSocketManager{
		connections: make(map[string]map[string]bool),
	}
}

// AddConnection adds a connection for a user
func (m *MockWebSocketManager) AddConnection(userID string, connectionID string) {
	if m.connections[userID] == nil {
		m.connections[userID] = make(map[string]bool)
	}
	m.connections[userID][connectionID] = true
}

// RemoveConnection removes a connection for a user
func (m *MockWebSocketManager) RemoveConnection(userID string, connectionID string) {
	if m.connections[userID] != nil {
		delete(m.connections[userID], connectionID)
		if len(m.connections[userID]) == 0 {
			delete(m.connections, userID)
		}
	}
}

// IsUserConnected checks if a user has any active connections
func (m *MockWebSocketManager) IsUserConnected(userID string) bool {
	conns, exists := m.connections[userID]
	return exists && len(conns) > 0
}

// SendMessage simulates sending a message to a user's active connections
// In a real implementation, this would push the message through WebSocket
// Here, we just simulate that the message was sent
func (m *MockWebSocketManager) SendMessage(userID string, message *database.Message) error {
	// Simulate message delivery - in real implementation, this would
	// push message through WebSocket connections
	if m.IsUserConnected(userID) {
		// Message would be delivered to user's inbox
		// For group messages, this would fan out to all group members
		return nil
	}
	// User is offline - message will be delivered when they come online
	return nil
}

// AckDelivered simulates a delivery acknowledgment from the client
// In a real implementation, the client would send this ACK over WebSocket
func (m *MockWebSocketManager) AckDelivered(userID string, messageID string) error {
	// This simulates the client sending a delivery ACK
	// In real implementation, this would come from WebSocket handler
	return nil
}

// AckRead simulates a read acknowledgment from the client
// In a real implementation, the client would send this ACK over WebSocket
func (m *MockWebSocketManager) AckRead(userID string, messageID string) error {
	// This simulates the client sending a read ACK
	// In real implementation, this would come from WebSocket handler
	return nil
}

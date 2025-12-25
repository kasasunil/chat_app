package database

// Repository defines the interface for database operations
// This allows for easy swapping between in-memory and real database implementations
// without changing any code in handlers, services, or other parts of the application.
//
// Any implementation (in-memory, PostgreSQL, MySQL, etc.) must implement this interface.
// The current implementation is MemoryStore which provides in-memory storage.
//
// To add a new database implementation:
//  1. Create a new struct that implements all methods in this interface
//  2. Update main.go to use the new implementation instead of NewStore()
//  3. No other code changes are needed
type Repository interface {
	// User operations
	CreateUser(user *User) error
	GetUser(userID string) (*User, error)

	// Group operations
	CreateGroup(group *Group) error
	GetGroup(groupID string) (*Group, error)
	AddGroupMember(groupID, userID string) error
	IsGroupMember(groupID, userID string) bool

	// Message operations
	CreateMessage(message *Message) error
	GetMessage(messageID string) (*Message, error)
	GetMessages(destinationID string, limit int, cursor string) ([]*Message, string, error)
	UpdateMessageStatus(messageID string, status MessageStatus) error

	// MessageRead operations
	CreateMessageRead(messageID, userID string) error
	GetMessageReads(messageID string) []*MessageRead

	// UserConversation operations
	GetUserConversations(userID string) ([]*UserConversation, error)

	// Search operations
	SearchMessages(userID, query string) ([]*Message, error)
}

// Production Environment Architecture:
// ------------------------------------
// In a production environment, once we have an actual database, the architecture will be
// restructured as follows:
//
// 1. Repository Interface (Base CRUD):
//   - Will define common CRUD operations (Create, Read, Update, Delete) that are
//     model-agnostic and work with any entity type
//   - Example: Create(entity interface{}), Update(id string, entity interface{}), etc.
//
// 2. Model-Specific Interfaces:
//   - Each model (User, Group, Message, etc.) will have its own interface
//   - These interfaces will embed the base Repository interface
//   - They will define model-specific methods (e.g., GetUserByEmail, GetMessagesByConversation)
//
// 3. Implementation Pattern:
//   - The Repository implementation will provide the core CRUD methods
//   - Each model's repository will embed this base Repository
//   - Model-specific methods will be wrapper functions that internally call the
//     base Repository CRUD methods with appropriate queries/filters
//
// Example Structure:
//
//	type BaseRepository interface {
//	    Create(entity interface{}) error
//	    Update(id string, entity interface{}) error
//	    Delete(id string) error
//	    FindByID(id string, entity interface{}) error
//	}
//
//	type UserRepository interface {
//	    BaseRepository  // Embed base CRUD
//	    GetUserByEmail(email string) (*User, error)  // Model-specific method
//	}
//
//	type UserRepo struct {
//	    BaseRepository  // Embed base implementation
//	}
//
//	func (r *UserRepo) GetUserByEmail(email string) (*User, error) {
//	    // Wrapper that calls base Repository methods with email filter
//	    return r.FindByFilter("email", email, &User{})
//	}
//
// This pattern provides:
// - Separation of concerns (base CRUD vs model-specific logic)
// - Code reusability (common CRUD operations)
// - Type safety (model-specific interfaces)
// - Flexibility (easy to add new models without duplicating CRUD code)

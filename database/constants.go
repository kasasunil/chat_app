package database

// Message status constants
const (
	StatusSentString      = "SENT"
	StatusDeliveredString = "DELIVERED"
	StatusReadString      = "READ"
)

// Conversation type constants
const (
	ConversationTypeOneToOneString = "one-to-one"
	ConversationTypeGroupString    = "group"
)

// Default values
const (
	DefaultMessageLimit      = 50
	MaxMessageLimit          = 100
	DefaultConversationLimit = 50
	MaxConversationLimit     = 100
	MaxMessageLength         = 10000
	MaxGroupMembers          = 100
)

// Database operation names
const (
	OpCreateUser           = "CreateUser"
	OpGetUser              = "GetUser"
	OpCreateGroup          = "CreateGroup"
	OpGetGroup             = "GetGroup"
	OpAddGroupMember       = "AddGroupMember"
	OpIsGroupMember        = "IsGroupMember"
	OpCreateMessage        = "CreateMessage"
	OpGetMessage           = "GetMessage"
	OpGetMessages          = "GetMessages"
	OpUpdateMessageStatus  = "UpdateMessageStatus"
	OpCreateMessageRead    = "CreateMessageRead"
	OpGetMessageReads      = "GetMessageReads"
	OpGetUserConversations = "GetUserConversations"
	OpSearchMessages       = "SearchMessages"
)

// Error messages
const (
	ErrUserAlreadyExists  = "user already exists"
	ErrUserNotFound       = "user not found"
	ErrGroupAlreadyExists = "group already exists"
	ErrGroupNotFound      = "group not found"
	ErrMessageNotFound    = "message not found"
	ErrInvalidCursor      = "invalid cursor"
)

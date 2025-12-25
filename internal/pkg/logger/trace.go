package logger

// Trace messages for server operations
const (
	TraceServerStarting    = "Server starting on %s:%s"
	TraceServerStarted     = "Server started successfully"
	TraceServerShutdown    = "Server shutting down"
	TraceServerFailed      = "Server failed to start: %v"
	TraceConfigLoaded      = "Configuration loaded from %s"
	TraceConfigLoadFailed  = "Failed to load config from %s: %v. Using defaults."
	TraceLoggerInitialized = "Logger initialized with level: %s"
	TraceLoggerInitFailed  = "Failed to initialize logger: %v. Using defaults."
)

// Trace messages for authentication
const (
	TraceAuthSuccess       = "User authenticated: %s"
	TraceAuthFailed        = "Authentication failed for username: %s"
	TraceAuthHeaderMissing = "Authorization header missing"
	TraceAuthInvalidFormat = "Invalid authorization header format"
	TraceAuthInvalidBase64 = "Invalid base64 encoding in authorization header"
	TraceAuthInvalidCreds  = "Invalid credentials format"
)

// Trace messages for user operations
const (
	TraceUserCreated       = "User created: id=%s, name=%s"
	TraceUserNotFound      = "User not found: id=%s"
	TraceUserAlreadyExists = "User already exists: id=%s"
	TraceUserRetrieved     = "User retrieved: id=%s"
)

// Trace messages for group operations
const (
	TraceGroupCreated       = "Group created: id=%s, name=%s"
	TraceGroupNotFound      = "Group not found: id=%s"
	TraceGroupMemberAdded   = "Member added to group: groupId=%s, userId=%s"
	TraceGroupMemberRemoved = "Member removed from group: groupId=%s, userId=%s"
	TraceGroupMemberCheck   = "Checking group membership: groupId=%s, userId=%s"
)

// Trace messages for message operations
const (
	TraceMessageSent          = "Message sent: sender=%s, destination=%s, type=%s, id=%s"
	TraceMessageCreated       = "Message created: id=%s, sender=%s, destination=%s"
	TraceMessageNotFound      = "Message not found: id=%s"
	TraceMessageStatusUpdated = "Message status updated: id=%s, status=%s"
	TraceMessageDelivered     = "Message delivered: id=%s, user=%s"
	TraceMessageRead          = "Message read acknowledged: messageID=%s, userID=%s"
	TraceMessageReadFailed    = "Failed to create read receipt: messageID=%s, userID=%s, error=%v"
	TraceMessageFetch         = "Fetching messages: destination=%s, limit=%d, cursor=%s"
	TraceMessageFetched       = "Messages fetched: destination=%s, count=%d"
)

// Trace messages for conversation operations
const (
	TraceConversationCreated     = "Conversation created: id=%s, user=%s, destination=%s"
	TraceConversationUpdated     = "Conversation updated: id=%s, user=%s"
	TraceConversationListFetched = "Conversation list fetched: user=%s, count=%d"
	TraceConversationNotFound    = "Conversation not found: id=%s"
)

// Trace messages for search operations
const (
	TraceSearchStarted   = "Search started: user=%s, query=%s"
	TraceSearchCompleted = "Search completed: user=%s, query=%s, results=%d"
	TraceSearchFailed    = "Search failed: user=%s, query=%s, error=%v"
)

// Trace messages for WebSocket operations
const (
	TraceWSConnectionAdded   = "WebSocket connection added: user=%s, connection=%s"
	TraceWSConnectionRemoved = "WebSocket connection removed: user=%s, connection=%s"
	TraceWSMessageSent       = "WebSocket message sent: user=%s, messageId=%s"
	TraceWSUserConnected     = "User connected: user=%s"
	TraceWSUserDisconnected  = "User disconnected: user=%s"
)

// Trace messages for database operations
const (
	TraceDBOperationStart   = "Database operation started: operation=%s"
	TraceDBOperationSuccess = "Database operation succeeded: operation=%s"
	TraceDBOperationFailed  = "Database operation failed: operation=%s, error=%v"
)

// Trace messages for API requests
const (
	TraceAPIRequestReceived = "API request received: method=%s, path=%s, user=%s"
	TraceAPIRequestSuccess  = "API request succeeded: method=%s, path=%s, status=%d"
	TraceAPIRequestFailed   = "API request failed: method=%s, path=%s, status=%d, error=%v"
)

// Trace messages for validation
const (
	TraceValidationFailed = "Validation failed: field=%s, error=%s"
	TraceValidationPassed = "Validation passed: field=%s"
)

// Trace messages for demo/initialization
const (
	TraceDemoDataInitialized = "Demo data initialized"
	TraceDemoUserCreated     = "Demo user created: id=%s, name=%s"
	TraceDemoGroupCreated    = "Demo group created: id=%s, name=%s"
	TraceDemoConnectionAdded = "Demo connection added: user=%s"
)

// Trace messages for errors
const (
	TraceErrorOccurred  = "Error occurred: code=%s, message=%s"
	TraceErrorHandled   = "Error handled: code=%s, status=%d"
	TraceErrorUnhandled = "Unhandled error: %v"
)

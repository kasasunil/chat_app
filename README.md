# Chat System - Production-Ready Implementation

A production-ready chat system implementation in Go with mocked WebSocket behavior, in-memory storage, and support for one-to-one and group messaging with read receipts. The codebase is structured for easy migration to production environments with real databases.

## Technical Specifications

For detailed technical specifications and design decisions, see the [Design Document](https://docs.google.com/document/d/1aLcT7AQUlreDYTA3JZsbcxSK-iRDs3fvHFjJ8d-REAg/edit?tab=t.0).

## Features

- ✅ One-to-one and group chat
- ✅ Message lifecycle: SENT (✓) → DELIVERED (✓✓) → READ (✓✓ blue)
- ✅ Cursor-based pagination for message fetching
- ✅ Conversation list view with accurate unread counts
- ✅ Keyword search across messages (case-insensitive)
- ✅ Mocked WebSocket implementation (via APIs)
- ✅ In-memory data storage (no external dependencies)
- ✅ TOML-based configuration management
- ✅ Basic Authentication middleware
- ✅ Graceful shutdown handling
- ✅ Centralized error handling with structured error codes
- ✅ Singleton logger with configurable levels
- ✅ Clean separation of concerns (handlers, services, bootstrap)

## Project Structure

```
chat_app/
├── bootstrap/              # Server initialization and graceful shutdown
│   └── bootstrap.go
├── cmd/
│   ├── server/            # Main server application
│   │   └── main.go
│   └── demo/              # End-to-end demo/test
│       └── main.go
├── conf/
│   └── config.toml        # Configuration file (can be overridden with prod.toml)
├── config/                # Configuration structs and loader
│   ├── config.go
│   └── constants.go
├── controller/            # HTTP handlers (one file per handler)
│   ├── handler.go         # Handler struct and initialization
│   ├── response.go        # Common error response helper
│   ├── send_message.go
│   ├── ack_delivered.go
│   ├── ack_read.go
│   ├── get_messages.go
│   ├── get_user_conversations.go
│   ├── search_messages.go
│   └── constants.go
├── database/              # Data models and repository interface
│   ├── store_interface.go # Repository interface (for easy DB migration)
│   ├── models.go          # Data models
│   ├── constants.go
│   └── in-memory/         # In-memory implementation
│       ├── store.go
│       ├── user.go
│       ├── group.go
│       ├── message.go
│       └── message_read.go
├── internal/
│   ├── middleware/        # Authentication middleware
│   │   ├── auth.go
│   │   └── constants.go
│   ├── pkg/
│   │   ├── errors/        # Centralized error codes and handling
│   │   │   └── errors.go
│   │   ├── logger/        # Singleton logger
│   │   │   ├── logger.go
│   │   │   └── trace.go   # Log message constants
│   │   └── utils/         # Common utility functions
│   │       └── utils.go
│   └── services/
│       ├── search/        # Message search service
│       │   └── search.go
│       └── websocket/     # Mocked WebSocket manager
│           ├── interfaces.go
│           └── manager.go
├── go.mod
├── README.md
├── TESTING.md             # Comprehensive testing guide
└── AUTHENTICATION.md      # Authentication guide
```

## Prerequisites

- Go 1.21 or later
- No external databases or services required (for in-memory mode)

## Installation

1. Clone or navigate to the project directory:
```bash
cd chat_app
```

2. Install dependencies:
```bash
go mod download
```

3. Configure the application:
   - Edit `conf/config.toml` to customize settings
   - Default configuration will be used if file is not found
   - For production, create `conf/prod.toml` and set `CONFIG_PATH` environment variable

## Running the Server

Start the server:
```bash
go run cmd/server/main.go
```

Or build and run:
```bash
go build -o bin/server ./cmd/server
./bin/server
```

The server will start on port `8080` by default (configurable in `conf/config.toml`).

### Graceful Shutdown

The server supports graceful shutdown:
- Press `Ctrl+C` or send `SIGTERM`/`SIGINT` to initiate shutdown
- Server will stop accepting new connections
- Existing requests will complete (up to 30 seconds timeout)
- All connections are closed gracefully before exit

## Running the Demo

In a separate terminal, run the end-to-end demo:
```bash
go run cmd/demo/main.go
```

The demo will:
1. Create users (Alice, Bob, Charlie)
2. Create a group (Project Team)
3. Send messages (one-to-one and group)
4. Simulate delivery ACKs
5. Simulate read ACKs
6. Fetch messages with pagination
7. Fetch conversation lists
8. Perform message searches

## Configuration

The application uses TOML configuration files located in the `conf/` directory. The config file contains:

- **Server settings**: port, host, read/write timeouts, idle timeout
- **Authentication settings**: Basic Auth credentials (username/password pairs)
- **Database settings**: mode (memory), max connections
- **Logging configuration**: level (debug, info, warn, error), format
- **Feature flags**: enable search, enable group chat, max message length, max group members

See `conf/config.toml` for the complete configuration structure.

### Environment-Specific Configuration

To use a custom config file location, set the `CONFIG_PATH` environment variable:
```bash
# Development
CONFIG_PATH=conf/config.toml go run cmd/server/main.go

# Production
CONFIG_PATH=conf/prod.toml go run cmd/server/main.go
```

## Authentication

All API endpoints (except `/health`) require authentication. See [AUTHENTICATION.md](AUTHENTICATION.md) for detailed authentication guide, including:
- Authentication method (Basic Auth)
- Configuration setup
- Usage examples
- Why Basic Auth was chosen over JWT
- Production considerations

## API Endpoints

**Note:** All endpoints require authentication (except `/health`). See [AUTHENTICATION.md](AUTHENTICATION.md) for details.

### 1. Health Check
**GET** `/health`

No authentication required. Returns `OK` if server is running.

### 2. Send Message
**POST** `/api/v1/sendMessage` or `/sendMessage`

Request body:
```json
{
  "sender_id": "user1",
  "destination_id": "user2",
  "message": "Hello!"
}
```

Response:
```json
{
  "message_id": "1234567890_...",
  "status": "SENT"
}
```

### 3. Acknowledge Delivery
**POST** `/api/v1/ack/delivered` or `/ack/delivered`

Request body:
```json
{
  "user_id": "user2",
  "message_id": "1234567890_..."
}
```

Response:
```json
{
  "message": "Message delivered"
}
```

### 4. Acknowledge Read
**POST** `/api/v1/ack/read` or `/ack/read`

Request body:
```json
{
  "user_id": "user2",
  "message_id": "1234567890_..."
}
```

Response:
```json
{
  "message": "Message read"
}
```

### 5. Fetch Messages
**GET** `/api/v1/conversations/{destinationId}/messages?cursor={cursor}&limit={limit}`

Query parameters:
- `cursor` (optional): Message ID for pagination
- `limit` (optional): Number of messages to fetch (default: 50)

Response:
```json
{
  "messages": [
    {
      "id": "...",
      "sender_id": "user1",
      "destination_id": "user2",
      "message_text": "Hello!",
      "status": "SENT",
      "conversation_type": "one-one",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "next_cursor": "cursor_string",
  "has_more": true
}
```

**Note:** Messages are returned in descending order (newest first).

### 6. Get User Conversations
**GET** `/api/v1/users/{userId}/conversations`

Response:
```json
{
  "conversations": [
    {
      "conversation_id": "uc_123",
      "destination_id": "user2",
      "conversation_type": "one-one",
      "last_message": {
        "id": "...",
        "message_text": "Hello!",
        "status": "SENT",
        "created_at": "2024-01-15T10:30:00Z"
      },
      "unread_count": 3,
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 7. Search Messages
**GET** `/api/v1/search/{userId}?query=hello`

Query parameters:
- `query` (required): Search keyword (case-insensitive)

Response:
```json
{
  "results": [
    {
      "id": "...",
      "sender_id": "user1",
      "destination_id": "user2",
      "message_text": "Hello Bob!",
      "conversation_type": "one-one",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "query": "hello"
}
```

## Error Handling

All errors follow a consistent JSON response format:

```json
{
  "error": {
    "code": "BAD_REQUEST_MESSAGE_EMPTY",
    "message": "Message cannot be empty"
  }
}
```

### Error Code Classification

Error codes are prefixed to indicate their HTTP status class:

- **BAD_REQUEST_***: 4xx client errors (400)
- **UNAUTHORIZED_***: Authentication errors (401)
- **FORBIDDEN_***: Authorization errors (403)
- **NOT_FOUND_***: Resource not found (404)
- **CONFLICT_***: Resource conflicts (409)
- **SERVER_ERROR_***: 5xx server errors (500)

This makes debugging easier as error codes clearly indicate the type of error.

## Message Status Flow

1. **SENT (✓)**: Message is created and stored. Status: `SENT`
2. **DELIVERED (✓✓)**: Recipient acknowledges delivery via `/ack/delivered`. Status: `DELIVERED`
3. **READ (✓✓ blue)**: Recipient acknowledges read via `/ack/read`. Status: `READ`

For group messages, read receipts are tracked per user in the `MessageReads` table.

## Testing

See [TESTING.md](TESTING.md) for comprehensive testing instructions including:
- Quick test script usage
- Step-by-step API testing
- Authentication testing
- Message lifecycle testing
- Pagination testing
- Search functionality testing
- Graceful shutdown testing
- Example cURL commands for all endpoints

## Architecture & Design Decisions

### 1. Repository Pattern
- **Interface-based design**: `database.Repository` interface allows easy swapping between implementations
- **In-memory implementation**: Current implementation uses Go maps and slices
- **Production migration**: Simply implement the `Repository` interface with your database (PostgreSQL, MySQL, etc.) and update `main.go` initialization - no business logic changes needed

### 2. Handler Organization
- **One handler per file**: Each API endpoint has its own file for better maintainability
- **Common error handling**: All handlers use `respondWithError()` for consistent error responses
- **Separation of concerns**: Handlers, services, and middleware are cleanly separated

### 3. Bootstrap Package
- **Router initialization**: All routes are configured in `bootstrap.SetupRouter()`
- **Graceful shutdown**: Server lifecycle management in `bootstrap.StartServerWithGracefulShutdown()`
- **Clean main**: `main.go` focuses on configuration and initialization

### 4. Error Management
- **Centralized errors**: All error codes defined in `internal/pkg/errors/errors.go`
- **Structured codes**: Error codes prefixed by HTTP status class (BAD_REQUEST_*, SERVER_ERROR_*, etc.)
- **Consistent format**: All errors return the same JSON structure

### 5. Logging
- **Singleton logger**: Configurable log levels (debug, info, warn, error)
- **Trace constants**: All log messages defined as constants in `logger/trace.go` for better debugging
- **Structured logging**: JSON format support for production environments

### 6. Configuration Management
- **TOML-based**: Human-readable configuration format
- **Environment-specific**: Easy to switch between dev/prod configs
- **Type-safe**: Configuration loaded into Go structs with validation

### 7. Message Ordering
- **Newest first**: Messages are returned in descending order by creation time
- **Cursor-based pagination**: Efficient pagination using message IDs as cursors

### 8. Unread Counts
- **Accurate tracking**: Unread count based on `MessageRead` entries, not message status
- **Per-user tracking**: Each user's read receipts tracked separately

## Code Quality

- ✅ Clean separation of concerns (models, store, handlers, services, bootstrap)
- ✅ Production-ready code structure
- ✅ Comprehensive error handling with structured error codes
- ✅ Centralized logging with configurable levels
- ✅ Graceful shutdown handling
- ✅ Interface-based design for easy database migration
- ✅ One handler per file for maintainability
- ✅ Common utility functions in shared package
- ✅ Constants for all fixed values
- ✅ No over-engineering - simple and maintainable

## Production Readiness

This codebase is **production-ready** and can be easily deployed to production environments. Here's why:

### 1. Environment Configuration Management
- **TOML file structure**: Production environment values can be easily updated by creating `conf/prod.toml`
- **Environment variables**: Support for `CONFIG_PATH` to switch between dev/prod configs
- **No code changes**: Configuration changes don't require code modifications or rebuilds

### 2. Database Migration Path
- **Interface-based design**: The `database.Repository` interface abstracts all database operations
- **Easy migration**: To add a real database (SQL/PostgreSQL/MySQL):
  1. Implement the `Repository` interface with your database client
  2. Update `main.go` to initialize the new database implementation
  3. **No business logic changes needed** - all handlers, services, and controllers work with the interface
- **Wrapper pattern**: The interface acts as a wrapper, so you only need to implement the CRUD operations for your database

### 3. Structured Error Handling
- **Error code classification**: Error codes are separated with proper prefixes (BAD_REQUEST_*, SERVER_ERROR_*, etc.)
- **Business error messages**: Each error has a clear, business-friendly message
- **Easy debugging**: Error codes clearly indicate the type and source of errors
- **Consistent format**: All errors follow the same JSON response structure

### 4. Graceful Shutdown
- **Signal handling**: Server listens for SIGTERM/SIGINT signals
- **Connection cleanup**: When a pod is killed, the server:
  - Stops accepting new HTTP connections
  - Waits for existing requests to complete (configurable timeout)
  - Closes idle HTTP connections
  - Can be extended to close database connections and other resources gracefully
- **Zero-downtime deployments**: Proper shutdown ensures no requests are dropped during deployments

### 5. Additional Production Features
- **Health check endpoint**: `/health` endpoint for load balancer health checks
- **Request timeouts**: Configurable read/write/idle timeouts prevent resource exhaustion
- **Structured logging**: JSON log format support for log aggregation systems
- **Authentication middleware**: Centralized authentication with Basic Auth support
- **Error tracking**: Centralized error handling makes it easy to integrate with error tracking services
- **Code organization**: Clean structure makes it easy to add monitoring, metrics, and tracing

### 6. Scalability Considerations
- **Stateless design**: In-memory store can be replaced with distributed storage
- **Interface abstraction**: Easy to add caching layers (Redis) between handlers and database
- **Service separation**: Search and WebSocket services are separate, making it easy to scale independently

## Migration to Production

### Step 1: Create Production Config
```bash
cp conf/config.toml conf/prod.toml
# Edit prod.toml with production values
```

### Step 2: Implement Database Repository
```go
// database/postgres/repository.go
type PostgresRepository struct {
    db *sql.DB
}

func (r *PostgresRepository) CreateUser(user *database.User) error {
    // Implement PostgreSQL create logic
}

// Implement all other Repository interface methods
```

### Step 3: Update main.go
```go
// Replace in-memory store with database implementation
// store := in_memory.NewStore()
store := postgres.NewRepository(dbConnection)
```

That's it! No other code changes needed.

## Notes

- **Current implementation**: Uses in-memory storage (data is lost on server restart)
- **WebSocket mocking**: WebSocket behavior is simulated through APIs, not real WebSocket connections
- **Search**: Basic case-insensitive keyword search (can be enhanced with full-text search engines)
- **Authentication**: Basic Auth is implemented (can be extended to JWT, OAuth, etc.)

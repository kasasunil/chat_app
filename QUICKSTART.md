# Quick Start Guide

## Prerequisites

- Go 1.21 or later installed
- Terminal/Command line access

## Setup

1. **Install dependencies:**
```bash
cd chat_app
go mod download
```

If `go mod download` doesn't work, try:
```bash
go get github.com/gorilla/mux
```

## Running the System

### Step 1: Start the Server

In one terminal:
```bash
go run cmd/server/main.go
```

You should see:
```
Server starting on port 8080...
API Endpoints:
  POST   /sendMessage
  POST   /ack/delivered
  POST   /ack/read
  GET    /conversations/{destinationId}/messages
  GET    /users/{userId}/conversations
  GET    /search/{userId}?query=xxx

Demo data initialized:
  Users: Alice (user1), Bob (user2), Charlie (user3)
  Group: Project Team (group1)
  All users are connected
```

### Step 2: Run the Demo

In another terminal:
```bash
go run cmd/demo/main.go
```

The demo will:
1. Send messages between users
2. Simulate delivery acknowledgments
3. Simulate read acknowledgments
4. Fetch messages with pagination
5. Fetch conversation lists
6. Perform message searches

## Expected Output

You should see output showing:
- Message sending with status transitions (SENT → DELIVERED → READ)
- Pagination working correctly
- Conversation lists with unread counts
- Search results

## Manual Testing

You can also test manually using curl:

### Send a message:
```bash
curl -X POST http://localhost:8080/sendMessage \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Hello!"}'
```

### Acknowledge delivery:
```bash
curl -X POST http://localhost:8080/ack/delivered \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user2","message_id":"1"}'
```

### Acknowledge read:
```bash
curl -X POST http://localhost:8080/ack/read \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user2","message_id":"1"}'
```

### Fetch messages:
```bash
curl http://localhost:8080/conversations/user2/messages
```

### Get conversations:
```bash
curl http://localhost:8080/users/user1/conversations
```

### Search:
```bash
curl "http://localhost:8080/search/user1?query=hello"
```

## Troubleshooting

If you get import errors:
1. Make sure you're in the `chat_app` directory
2. Run `go mod tidy`
3. Run `go mod download`

If the server doesn't start:
- Check if port 8080 is already in use
- Modify the port in `config/config.go` if needed


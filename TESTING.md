# Testing and Validation Guide

This guide provides two types of testing approaches: **One-Click Testing** (automated) and **Manual Testing** (step-by-step).

## Prerequisites

- Go 1.21 or later installed
- `curl` command-line tool (for manual testing)
- Terminal access

## One-Click Testing (Automated)

The fastest way to test all functionality is using the automated demo script.

### Steps:

1. **Start the Server** (Terminal 1):
```bash
cd /Users/k.rsunilkumar/Development/chat_app
go run cmd/server/main.go
```

Wait for the server to start. You should see:
```
[INFO] Server starting on 0.0.0.0:8080
[INFO] Server started!!!
```

2. **Run the Demo Script** (Terminal 2):
```bash
cd /Users/k.rsunilkumar/Development/chat_app
go run cmd/demo/main.go
```

### What the Demo Tests:

The demo script automatically tests:
- ✅ User and group creation (pre-seeded by server)
- ✅ Sending one-to-one messages
- ✅ Sending group messages
- ✅ Delivery acknowledgments (SENT → DELIVERED)
- ✅ Read acknowledgments (DELIVERED → READ)
- ✅ Message fetching with pagination
- ✅ Conversation list retrieval
- ✅ Message search functionality

### Expected Output:

```
==============================================
Chat System Demo - End-to-End Test
==============================================

=== STEP 1: Users and Groups ===
Users: user1 (Alice), user2 (Bob), user3 (Charlie)
Group: group1 (Project Team)

=== STEP 2: Send Messages ===
[user1 → user2] Message sent: Hello Bob! How are you?
[user2 → user1] Message sent: Hi Alice! I'm doing great, thanks!
[user1 → group1] Message sent: Hello team! Let's discuss the project.
[user2 → group1] Message sent: Sure, I'm ready to discuss.

=== STEP 3: Simulate Delivery ACK ===
[user2] Delivery ACK for message ...
[user1] Delivery ACK for message ...
...

=== STEP 4: Simulate Read ACK ===
[user2] Read ACK for message ...
[user1] Read ACK for message ...
...

=== STEP 5: Fetch Messages with Pagination ===
...

=== STEP 6: Fetch Conversation List ===
...

=== STEP 7: Search Messages ===
...

==============================================
Demo completed successfully!
==============================================
```

If you see "Demo completed successfully!", all automated tests passed.

---

## Manual Testing (Step-by-Step)

For detailed validation and understanding of each API, follow these manual test cases.

### Setup

1. **Start the Server**:
```bash
cd /Users/k.rsunilkumar/Development/chat_app
go run cmd/server/main.go
```

2. **Verify Server is Running**:
```bash
curl -X GET http://localhost:8080/health
```
**Expected:** `OK`

### Authentication Credentials

The system uses Basic Authentication. Credentials are configured in `conf/config.toml`:
- `user1` / `password1` (Alice)
- `user2` / `password2` (Bob)
- `user3` / `password3` (Charlie)

**Helper for Basic Auth:**
```bash
# Create a helper variable (optional, for easier testing)
AUTH_USER1="Authorization: Basic $(echo -n 'user1:password1' | base64)"
AUTH_USER2="Authorization: Basic $(echo -n 'user2:password2' | base64)"
AUTH_USER3="Authorization: Basic $(echo -n 'user3:password3' | base64)"
```

---

## Test Case 1: Happy Path - One-to-One Message with Read Receipts

**Objective:** Test complete message lifecycle for one-to-one conversation, including status transitions.

### Step 1.1: Send One-to-One Message

**API:** `POST /api/v1/sendMessage`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "user1",
    "destination_id": "user2",
    "message": "Hello Bob! This is a test message."
  }'
```

**Expected Response (200 OK):**
```json
{
  "message_id": "1234567890_...",
  "status": "SENT"
}
```

**✅ Validation:**
- Status code: `200`
- Response contains `message_id`
- Status is `SENT`

**Save the `message_id` for next steps!**

### Step 1.2: Verify Message Status (SENT)

**API:** `GET /api/v1/conversations/{destinationId}/messages`

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/conversations/user2/messages?limit=1" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "messages": [
    {
      "id": "1234567890_...",
      "sender_id": "user1",
      "destination_id": "user2",
      "message_text": "Hello Bob! This is a test message.",
      "status": "SENT",
      "conversation_type": "one-one",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "next_cursor": "",
  "has_more": false
}
```

**✅ Validation:**
- Status code: `200`
- Message status is `SENT`
- Message text matches what was sent
- Messages are ordered newest first (this message should be first)

### Step 1.3: Acknowledge Delivery

**API:** `POST /api/v1/ack/delivered`

**Request:**
```bash
# Replace MESSAGE_ID with the actual message_id from Step 1.1
curl -X POST http://localhost:8080/api/v1/ack/delivered \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user2",
    "message_id": "MESSAGE_ID"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Message delivered"
}
```

**✅ Validation:**
- Status code: `200`
- Response message confirms delivery

### Step 1.4: Verify Message Status (DELIVERED)

**API:** `GET /api/v1/conversations/{destinationId}/messages`

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/conversations/user2/messages?limit=1" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "messages": [
    {
      "id": "1234567890_...",
      "status": "DELIVERED",
      ...
    }
  ]
}
```

**✅ Validation:**
- Status code: `200`
- Message status changed from `SENT` to `DELIVERED`

### Step 1.5: Acknowledge Read

**API:** `POST /api/v1/ack/read`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/ack/read \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user2",
    "message_id": "MESSAGE_ID"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Message read"
}
```

**✅ Validation:**
- Status code: `200`
- Response message confirms read

### Step 1.6: Verify Message Status (READ)

**API:** `GET /api/v1/conversations/{destinationId}/messages`

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/conversations/user2/messages?limit=1" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "messages": [
    {
      "id": "1234567890_...",
      "status": "READ",
      ...
    }
  ]
}
```

**✅ Validation:**
- Status code: `200`
- Message status changed from `DELIVERED` to `READ`
- **Complete lifecycle verified:** SENT → DELIVERED → READ

### Step 1.7: Verify Unread Count

**API:** `GET /api/v1/users/{userId}/conversations`

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/user2/conversations \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "conversations": [
    {
      "conversation_id": "...",
      "destination_id": "user1",
      "conversation_type": "one-one",
      "last_message": {
        "id": "1234567890_...",
        "status": "READ",
        ...
      },
      "unread_count": 0,
      ...
    }
  ]
}
```

**✅ Validation:**
- Status code: `200`
- `unread_count` is `0` (message has been read)
- Last message status is `READ`

---

## Test Case 2: Happy Path - Group Message with Multiple Read Receipts

**Objective:** Test group message functionality with multiple users acknowledging delivery and read.

### Step 2.1: Send Group Message

**API:** `POST /api/v1/sendMessage`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "user1",
    "destination_id": "group1",
    "message": "Hello team! This is a group message for testing."
  }'
```

**Expected Response (200 OK):**
```json
{
  "message_id": "GROUP_MSG_ID_...",
  "status": "SENT"
}
```

**✅ Validation:**
- Status code: `200`
- Response contains `message_id`
- Status is `SENT`

**Save the `message_id` for next steps!**

### Step 2.2: Verify Group Message

**API:** `GET /api/v1/conversations/{destinationId}/messages`

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/conversations/group1/messages?limit=1" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "messages": [
    {
      "id": "GROUP_MSG_ID_...",
      "sender_id": "user1",
      "destination_id": "group1",
      "message_text": "Hello team! This is a group message for testing.",
      "status": "SENT",
      "conversation_type": "group",
      ...
    }
  ]
}
```

**✅ Validation:**
- Status code: `200`
- `conversation_type` is `group`
- Message is visible to all group members

### Step 2.3: User2 Acknowledges Delivery

**API:** `POST /api/v1/ack/delivered`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/ack/delivered \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user2",
    "message_id": "GROUP_MSG_ID"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Message delivered"
}
```

**✅ Validation:**
- Status code: `200`

### Step 2.4: User3 Acknowledges Delivery

**API:** `POST /api/v1/ack/delivered`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/ack/delivered \
  -H "Authorization: Basic $(echo -n 'user3:password3' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user3",
    "message_id": "GROUP_MSG_ID"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Message delivered"
}
```

**✅ Validation:**
- Status code: `200`
- Multiple users can acknowledge the same message

### Step 2.5: User2 Acknowledges Read

**API:** `POST /api/v1/ack/read`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/ack/read \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user2",
    "message_id": "GROUP_MSG_ID"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Message read"
}
```

**✅ Validation:**
- Status code: `200`

### Step 2.6: User3 Acknowledges Read

**API:** `POST /api/v1/ack/read`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/ack/read \
  -H "Authorization: Basic $(echo -n 'user3:password3' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user3",
    "message_id": "GROUP_MSG_ID"
  }'
```

**Expected Response (200 OK):**
```json
{
  "message": "Message read"
}
```

**✅ Validation:**
- Status code: `200`
- Multiple users can acknowledge read for the same message

### Step 2.7: Verify Read Receipts (Per User)

**API:** `GET /api/v1/users/{userId}/conversations`

**Request (as user2):**
```bash
curl -X GET http://localhost:8080/api/v1/users/user2/conversations \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "conversations": [
    {
      "destination_id": "group1",
      "conversation_type": "group",
      "unread_count": 0,
      ...
    }
  ]
}
```

**✅ Validation:**
- Status code: `200`
- `unread_count` is `0` for user2 (they read the message)

**Request (as user3):**
```bash
curl -X GET http://localhost:8080/api/v1/users/user3/conversations \
  -H "Authorization: Basic $(echo -n 'user3:password3' | base64)"
```

**✅ Validation:**
- Status code: `200`
- `unread_count` is `0` for user3 (they also read the message)
- **Each user's read receipt is tracked separately**

---

## Test Case 3: Security and Authorization - Access Control

**Objective:** Verify that users cannot access other users' data and proper error handling.

### Step 3.1: Attempt to Access Another User's Conversations

**API:** `GET /api/v1/users/{userId}/conversations`

**Request:** User1 trying to access User2's conversations
```bash
curl -X GET http://localhost:8080/api/v1/users/user2/conversations \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (403 Forbidden):**
```json
{
  "error": {
    "code": "FORBIDDEN_ACCESS_DENIED",
    "message": "Access denied"
  }
}
```

**✅ Validation:**
- Status code: `403 Forbidden`
- Error code indicates forbidden access
- User cannot access another user's conversations

### Step 3.2: Attempt to Search Another User's Messages

**API:** `GET /api/v1/search/{userId}?query=xxx`

**Request:** User1 trying to search User2's messages
```bash
curl -X GET "http://localhost:8080/api/v1/search/user2?query=hello" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (403 Forbidden):**
```json
{
  "error": {
    "code": "FORBIDDEN_ACCESS_DENIED",
    "message": "Access denied"
  }
}
```

**✅ Validation:**
- Status code: `403 Forbidden`
- User cannot search another user's messages

### Step 3.3: Attempt to Acknowledge Message as Non-Recipient

**API:** `POST /api/v1/ack/read`

**Scenario:** User1 sends a message to User2. User3 tries to acknowledge it.

**Step 3.3a: User1 sends message to User2**
```bash
MESSAGE_ID=$(curl -s -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "user1",
    "destination_id": "user2",
    "message": "Private message for Bob"
  }' | jq -r '.message_id')

echo "Message ID: $MESSAGE_ID"
```

**Step 3.3b: User3 tries to acknowledge (should fail)**
```bash
curl -X POST http://localhost:8080/api/v1/ack/read \
  -H "Authorization: Basic $(echo -n 'user3:password3' | base64)" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"user3\",
    \"message_id\": \"$MESSAGE_ID\"
  }"
```

**Expected Response (403 Forbidden):**
```json
{
  "error": {
    "code": "FORBIDDEN_NOT_MESSAGE_RECIPIENT",
    "message": "User is not the recipient of this message"
  }
}
```

**✅ Validation:**
- Status code: `403 Forbidden`
- Error message clearly indicates user is not the recipient
- **Security: Users can only acknowledge messages they received**

### Step 3.4: Attempt to Send Message with Mismatched Sender ID

**API:** `POST /api/v1/sendMessage`

**Request:** User1 authenticated but trying to send as User2
```bash
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "user2",
    "destination_id": "user3",
    "message": "Trying to impersonate user2"
  }'
```

**Expected Response (403 Forbidden):**
```json
{
  "error": {
    "code": "FORBIDDEN_ACCESS_DENIED",
    "message": "Access denied"
  }
}
```

**✅ Validation:**
- Status code: `403 Forbidden`
- **Security: Users cannot impersonate other users**

---

## Test Case 4: Pagination and Message Ordering

**Objective:** Verify cursor-based pagination and message ordering (newest first).

### Step 4.1: Send Multiple Messages

**API:** `POST /api/v1/sendMessage`

Send 3 messages in sequence:
```bash
# Message 1
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Message 1"}'

sleep 1

# Message 2
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Message 2"}'

sleep 1

# Message 3
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Message 3"}'
```

### Step 4.2: Verify Message Ordering (Newest First)

**API:** `GET /api/v1/conversations/{destinationId}/messages?limit=3`

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/conversations/user2/messages?limit=3" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "messages": [
    {
      "message_text": "Message 3",
      "created_at": "2024-01-15T10:33:00Z"
    },
    {
      "message_text": "Message 2",
      "created_at": "2024-01-15T10:32:00Z"
    },
    {
      "message_text": "Message 1",
      "created_at": "2024-01-15T10:31:00Z"
    }
  ],
  "next_cursor": "...",
  "has_more": false
}
```

**✅ Validation:**
- Status code: `200`
- Messages are ordered newest first (Message 3, then 2, then 1)
- `created_at` timestamps are in descending order

### Step 4.3: Test Pagination with Cursor

**API:** `GET /api/v1/conversations/{destinationId}/messages?cursor={cursor}&limit=2`

**Step 4.3a: Get first page**
```bash
RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/conversations/user2/messages?limit=2" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)")

echo "$RESPONSE" | jq '.'
```

**Expected:**
- Returns 2 messages (newest 2)
- `has_more: true` (if there are more messages)
- `next_cursor` contains the message ID of the last message in this page

**Step 4.3b: Get next page using cursor**
```bash
CURSOR=$(echo "$RESPONSE" | jq -r '.next_cursor')

curl -X GET "http://localhost:8080/api/v1/conversations/user2/messages?cursor=$CURSOR&limit=2" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "messages": [
    {
      "message_text": "Message 1",
      ...
    }
  ],
  "next_cursor": "",
  "has_more": false
}
```

**✅ Validation:**
- Status code: `200`
- Returns messages after the cursor
- `has_more: false` when no more messages
- **Pagination works correctly**

---

## Test Case 5: Search Functionality

**Objective:** Verify message search works correctly and is case-insensitive.

### Step 5.1: Send Messages with Different Keywords

**API:** `POST /api/v1/sendMessage`

```bash
# Message with "project" keyword
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Let'\''s discuss the project details"}'

# Message with "hello" keyword
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Hello Bob! How are you?"}'
```

### Step 5.2: Search for Messages (Case-Insensitive)

**API:** `GET /api/v1/search/{userId}?query=xxx`

**Request (lowercase):**
```bash
curl -X GET "http://localhost:8080/api/v1/search/user1?query=project" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "results": [
    {
      "message_text": "Let's discuss the project details",
      ...
    }
  ],
  "query": "project"
}
```

**Request (uppercase):**
```bash
curl -X GET "http://localhost:8080/api/v1/search/user1?query=PROJECT" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "results": [
    {
      "message_text": "Let's discuss the project details",
      ...
    }
  ],
  "query": "PROJECT"
}
```

**✅ Validation:**
- Status code: `200`
- Both lowercase and uppercase queries return the same results
- **Search is case-insensitive**
- Only returns messages where the user is a participant

### Step 5.3: Search with Empty Query (Should Fail)

**API:** `GET /api/v1/search/{userId}?query=`

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/search/user1?query=" \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)"
```

**Expected Response (400 Bad Request):**
```json
{
  "error": {
    "code": "BAD_REQUEST_SEARCH_QUERY_REQUIRED",
    "message": "Query parameter is required"
  }
}
```

**✅ Validation:**
- Status code: `400 Bad Request`
- Error code indicates missing query parameter

---

## Test Case 6: Error Handling and Validation

**Objective:** Verify proper error handling for invalid requests.

### Step 6.1: Send Message with Empty Text

**API:** `POST /api/v1/sendMessage`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "user1",
    "destination_id": "user2",
    "message": ""
  }'
```

**Expected Response (400 Bad Request):**
```json
{
  "error": {
    "code": "BAD_REQUEST_MESSAGE_EMPTY",
    "message": "Message cannot be empty"
  }
}
```

**✅ Validation:**
- Status code: `400 Bad Request`
- Error code clearly indicates the issue

### Step 6.2: Send Message to Non-Existent User

**API:** `POST /api/v1/sendMessage`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "user1",
    "destination_id": "nonexistent_user",
    "message": "Hello"
  }'
```

**Expected Response (404 Not Found):**
```json
{
  "error": {
    "code": "NOT_FOUND_DESTINATION_NOT_FOUND",
    "message": "Destination not found"
  }
}
```

**✅ Validation:**
- Status code: `404 Not Found`
- Error indicates destination doesn't exist

### Step 6.3: Acknowledge Non-Existent Message

**API:** `POST /api/v1/ack/read`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/ack/read \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user2",
    "message_id": "nonexistent_message_id"
  }'
```

**Expected Response (404 Not Found):**
```json
{
  "error": {
    "code": "NOT_FOUND_MESSAGE_NOT_FOUND",
    "message": "Message not found"
  }
}
```

**✅ Validation:**
- Status code: `404 Not Found`
- Error indicates message doesn't exist

### Step 6.4: Missing Authentication

**API:** `GET /api/v1/users/{userId}/conversations`

**Request (no Authorization header):**
```bash
curl -X GET http://localhost:8080/api/v1/users/user1/conversations
```

**Expected Response (401 Unauthorized):**
```json
{
  "error": {
    "code": "UNAUTHORIZED_AUTH_REQUIRED",
    "message": "Authorization header required"
  }
}
```

**✅ Validation:**
- Status code: `401 Unauthorized`
- Error indicates authentication is required

### Step 6.5: Invalid Authentication Credentials

**API:** `GET /api/v1/users/{userId}/conversations`

**Request (invalid credentials):**
```bash
curl -X GET http://localhost:8080/api/v1/users/user1/conversations \
  -H "Authorization: Basic $(echo -n 'invalid:password' | base64)"
```

**Expected Response (401 Unauthorized):**
```json
{
  "error": {
    "code": "UNAUTHORIZED_INVALID_CREDENTIALS",
    "message": "Invalid username or password"
  }
}
```

**✅ Validation:**
- Status code: `401 Unauthorized`
- Error indicates invalid credentials

---

## Test Case 7: Idempotency - Multiple ACKs for Same Message

**Objective:** Verify that multiple acknowledgments for the same message are handled correctly (idempotent).

### Step 7.1: Send Message and Acknowledge Read

**Step 7.1a: Send message**
```bash
MESSAGE_ID=$(curl -s -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Idempotency test"}' \
  | jq -r '.message_id')
```

**Step 7.1b: First read acknowledgment**
```bash
curl -X POST http://localhost:8080/api/v1/ack/read \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"user2\",\"message_id\":\"$MESSAGE_ID\"}"
```

**Expected Response (200 OK):**
```json
{
  "message": "Message read"
}
```

### Step 7.2: Acknowledge Read Again (Idempotency Test)

**API:** `POST /api/v1/ack/read`

**Request (same message, same user):**
```bash
curl -X POST http://localhost:8080/api/v1/ack/read \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"user2\",\"message_id\":\"$MESSAGE_ID\"}"
```

**Expected Response (200 OK):**
```json
{
  "message": "Message read"
}
```

**✅ Validation:**
- Status code: `200 OK`
- **Idempotent:** Multiple acknowledgments don't cause errors
- System handles duplicate acknowledgments gracefully

---

## Test Case 8: Unread Count Accuracy

**Objective:** Verify unread counts are calculated correctly based on MessageRead entries.

### Step 8.1: Send Multiple Messages

**API:** `POST /api/v1/sendMessage`

Send 3 messages from user1 to user2:
```bash
# Message 1
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Unread test 1"}'

sleep 1

# Message 2
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Unread test 2"}'

sleep 1

# Message 3
MSG3_ID=$(curl -s -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic $(echo -n 'user1:password1' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"sender_id":"user1","destination_id":"user2","message":"Unread test 3"}' \
  | jq -r '.message_id')
```

### Step 8.2: Check Unread Count (Should be 3)

**API:** `GET /api/v1/users/{userId}/conversations`

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/user2/conversations \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "conversations": [
    {
      "destination_id": "user1",
      "unread_count": 3,
      ...
    }
  ]
}
```

**✅ Validation:**
- Status code: `200`
- `unread_count` is `3` (all 3 messages are unread)

### Step 8.3: Read One Message

**API:** `POST /api/v1/ack/read`

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/ack/read \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"user2\",\"message_id\":\"$MSG3_ID\"}"
```

### Step 8.4: Check Unread Count Again (Should be 2)

**API:** `GET /api/v1/users/{userId}/conversations`

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/user2/conversations \
  -H "Authorization: Basic $(echo -n 'user2:password2' | base64)"
```

**Expected Response (200 OK):**
```json
{
  "conversations": [
    {
      "destination_id": "user1",
      "unread_count": 2,
      ...
    }
  ]
}
```

**✅ Validation:**
- Status code: `200`
- `unread_count` decreased from `3` to `2`
- **Unread count is accurate and updates correctly**

---

## Test Case 9: Graceful Shutdown

**Objective:** Verify the server handles shutdown signals gracefully.

### Step 9.1: Start Server

**Terminal 1:**
```bash
go run cmd/server/main.go
```

Wait for: `[INFO] Server started!!!`

### Step 9.2: Send Test Request

**Terminal 2:**
```bash
curl -X GET http://localhost:8080/health
```

**Expected:** `OK`

### Step 9.3: Trigger Graceful Shutdown

**Terminal 1:** Press `Ctrl+C` or send SIGTERM:
```bash
# In another terminal
pkill -TERM -f "go run cmd/server/main.go"
```

### Step 9.4: Verify Shutdown Behavior

**Expected Server Output:**
```
[INFO] Received signal: interrupt. Initiating graceful shutdown...
[INFO] Server gracefully stopped
[INFO] Server shutdown complete
```

**✅ Validation:**
- Server receives the signal
- Logs shutdown initiation
- Waits for existing requests to complete
- Closes connections gracefully
- Exits cleanly

### Step 9.5: Verify New Requests Are Rejected

While server is shutting down, try:
```bash
curl -X GET http://localhost:8080/health
```

**Expected:** Connection refused or timeout (server is shutting down)

---

## Quick Reference: All Test Cases

| Test Case | Description | Key Validations |
|-----------|-------------|-----------------|
| **1** | One-to-One Message Lifecycle | SENT → DELIVERED → READ status transitions |
| **2** | Group Message with Multiple ACKs | Multiple users can acknowledge, per-user tracking |
| **3** | Security and Authorization | Users cannot access other users' data |
| **4** | Pagination and Ordering | Messages ordered newest first, cursor pagination works |
| **5** | Search Functionality | Case-insensitive search, only user's messages |
| **6** | Error Handling | Proper error codes and messages for invalid requests |
| **7** | Idempotency | Multiple ACKs don't cause errors |
| **8** | Unread Count Accuracy | Count based on MessageRead entries, updates correctly |
| **9** | Graceful Shutdown | Server handles SIGTERM/SIGINT gracefully |

---

## Troubleshooting

### Server Won't Start
- Check if port 8080 is already in use: `lsof -i :8080`
- Verify `conf/config.toml` exists and is valid
- Check Go version: `go version` (should be 1.21+)

### Authentication Fails
- Verify credentials in `conf/config.toml`
- Check Base64 encoding: `echo -n 'user1:password1' | base64`
- Ensure `Authorization` header format: `Basic <base64_credentials>`

### Messages Not Appearing
- Verify sender and destination IDs exist
- Check conversation type (one-one vs group)
- Ensure user is a member of the group (for group messages)

### Unread Counts Incorrect
- Verify MessageRead entries are created when acknowledging read
- Check that unread count only counts messages without MessageRead entry for that user

---

## Summary

This testing guide provides:
- ✅ **One-click testing** via automated demo script
- ✅ **Detailed manual test cases** with step-by-step instructions
- ✅ **Clear validation criteria** for each step
- ✅ **Security and authorization testing**
- ✅ **Error handling validation**
- ✅ **Idempotency verification**

Follow these test cases to thoroughly validate all functionality of the chat application.

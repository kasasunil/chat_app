# Authentication Guide

## Overview

The chat system implements Basic Authentication middleware that protects all API endpoints (except `/health`). The system uses username/password authentication configured in the TOML file.

## Why Basic Auth Instead of JWT?

**JWT (JSON Web Tokens) would be ideal for this type of application**, but we've chosen Basic Authentication for the following reasons:

1. **No User Registration/Signup APIs**: The application doesn't expose signup or initial customer registration APIs. All users are pre-seeded during server bootup (see `bootstrap.SetupDemoData()`).

2. **Simplicity for Current Use Case**: Since all users are known at startup and configured in the TOML file, Basic Auth provides a simpler authentication mechanism without the overhead of token generation, validation, and refresh logic.

3. **Easy Migration to JWT**: The authentication middleware is designed to be easily extensible. For production environments, JWT can be implemented with minimal changes:
   - Add JWT token generation endpoint (if user registration is added)
   - Update the middleware to validate JWT tokens instead of Basic Auth
   - The rest of the application (handlers, services) remains unchanged

4. **Production Ready**: The current Basic Auth implementation is production-ready when used over HTTPS. For production, you can:
   - Keep Basic Auth if users are managed externally
   - Switch to JWT by updating only the middleware layer
   - Add OAuth/OIDC integration using the same middleware pattern

The authentication architecture is designed to be **swappable** - you can replace the authentication method without changing any business logic.

## Configuration

Authentication credentials are configured in `conf/config.toml`:

```toml
[auth]

[auth.client1]
username = "user1"
password = "password1"

[auth.client2]
username = "user2"
password = "password2"

[auth.client3]
username = "user3"
password = "password3"
```

You can add as many clients as needed by adding more `[auth.clientN]` sections.

## Authentication Method

### Basic Authentication

Send credentials using HTTP Basic Authentication in the `Authorization` header:

```
Authorization: Basic <base64(username:password)>
```

The middleware iterates through all configured auth clients and checks if the provided username and password match any of them.

**Example:**
```bash
# Using curl with -u flag (automatically encodes to Basic auth)
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -u user1:password1 \
  -H "Content-Type: application/json" \
  -d '{"destination_id":"user2","message":"Hello!"}'
```

**Manual Basic Auth:**
```bash
# Encode username:password to base64
# user1:password1 -> dXNlcjE6cGFzc3dvcmQx

curl -X POST http://localhost:8080/api/v1/sendMessage \
  -H "Authorization: Basic dXNlcjE6cGFzc3dvcmQx" \
  -H "Content-Type: application/json" \
  -d '{"destination_id":"user2","message":"Hello!"}'
```

## Protected Endpoints

All endpoints except `/health` require authentication:

- `POST /api/v1/sendMessage`
- `POST /api/v1/ack/delivered`
- `POST /api/v1/ack/read`
- `GET /api/v1/conversations/{destinationId}/messages`
- `GET /api/v1/users/{userId}/conversations`
- `GET /api/v1/search/{userId}?query=xxx`

Legacy routes (without `/api/v1` prefix) are also protected:
- `POST /sendMessage`
- `POST /ack/delivered`
- `POST /ack/read`
- `GET /conversations/{destinationId}/messages`
- `GET /users/{userId}/conversations`
- `GET /search/{userId}?query=xxx`

## Security Features

1. **User Context**: The authenticated user ID is stored in the request context and used by handlers
2. **Authorization Checks**: Handlers validate that users can only access their own data
3. **Sender Validation**: When sending messages, the authenticated user is used as the sender (request body sender_id is validated to match)
4. **Resource Access**: Users can only fetch their own conversations and search their own messages

## Error Responses

All authentication errors follow the common error response format:

### Unauthorized (401)
```json
{
  "error": {
    "code": "UNAUTHORIZED_AUTH_REQUIRED",
    "message": "Authorization header required"
  }
}
```

Other possible 401 errors:
- `UNAUTHORIZED_INVALID_CREDENTIALS`: Invalid username or password
- `UNAUTHORIZED_INVALID_AUTH_FORMAT`: Invalid authorization header format
- `UNAUTHORIZED_INVALID_BASE64`: Invalid base64 encoding
- `UNAUTHORIZED_INVALID_CREDENTIALS_FORMAT`: Invalid credentials format

### Forbidden (403)
```json
{
  "error": {
    "code": "FORBIDDEN_ACCESS_DENIED",
    "message": "Access denied"
  }
}
```

See the [Error Handling](#error-handling) section in README.md for more details on error code classification.

## Testing Authentication

### Test with Basic Auth:
```bash
# Send message
curl -X POST http://localhost:8080/api/v1/sendMessage \
  -u user1:password1 \
  -H "Content-Type: application/json" \
  -d '{"destination_id":"user2","message":"Test message"}'

# Get conversations
curl -u user1:password1 \
  http://localhost:8080/api/v1/users/user1/conversations

# Acknowledge delivery
curl -X POST http://localhost:8080/api/v1/ack/delivered \
  -u user2:password2 \
  -H "Content-Type: application/json" \
  -d '{"message_id":"123"}'

# Search messages
curl -u user1:password1 \
  "http://localhost:8080/api/v1/search/user1?query=hello"
```

## Production Considerations

For production use, consider:

1. **HTTPS**: Always use HTTPS in production (Basic Auth sends credentials in base64, which is easily decoded)
2. **Secret Management**: Store passwords as hashed values, not plain text
3. **Password Hashing**: Use bcrypt or similar for password storage
4. **User Database**: Store credentials in a secure database instead of config files
5. **Rate Limiting**: Add rate limiting to prevent brute force attacks
6. **Account Lockout**: Implement account lockout after failed login attempts
7. **Environment Variables**: Use environment variables or secret managers for sensitive credentials

## Configuration File Location

By default, the config file is loaded from `conf/config.toml`. You can override this with the `CONFIG_PATH` environment variable:

```bash
CONFIG_PATH=/path/to/config.toml go run cmd/server/main.go
```


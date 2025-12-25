package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents a system error code
type ErrorCode string

// Error code prefixes for classification
const (
	// 4xx Client Errors
	PrefixBadRequest   = "BAD_REQUEST"
	PrefixUnauthorized = "UNAUTHORIZED"
	PrefixForbidden    = "FORBIDDEN"
	PrefixNotFound     = "NOT_FOUND"
	PrefixConflict     = "CONFLICT"

	// 5xx Server Errors
	PrefixServerError = "SERVER_ERROR"
)

const (
	// 4xx - Bad Request errors
	ErrCodeBadRequestInvalidRequest      ErrorCode = PrefixBadRequest + "_INVALID_REQUEST"
	ErrCodeBadRequestValidationError     ErrorCode = PrefixBadRequest + "_VALIDATION_ERROR"
	ErrCodeBadRequestMessageEmpty        ErrorCode = PrefixBadRequest + "_MESSAGE_EMPTY"
	ErrCodeBadRequestMessageTooLong      ErrorCode = PrefixBadRequest + "_MESSAGE_TOO_LONG"
	ErrCodeBadRequestInvalidRecipient    ErrorCode = PrefixBadRequest + "_INVALID_RECIPIENT"
	ErrCodeBadRequestSearchQueryRequired ErrorCode = PrefixBadRequest + "_SEARCH_QUERY_REQUIRED"
	ErrCodeBadRequestGroupMemberLimit    ErrorCode = PrefixBadRequest + "_GROUP_MEMBER_LIMIT"
	ErrCodeBadRequestInvalidConversation ErrorCode = PrefixBadRequest + "_INVALID_CONVERSATION"

	// 4xx - Unauthorized errors
	ErrCodeUnauthorizedAuthRequired             ErrorCode = PrefixUnauthorized + "_AUTH_REQUIRED"
	ErrCodeUnauthorizedInvalidCredentials       ErrorCode = PrefixUnauthorized + "_INVALID_CREDENTIALS"
	ErrCodeUnauthorizedInvalidAuthFormat        ErrorCode = PrefixUnauthorized + "_INVALID_AUTH_FORMAT"
	ErrCodeUnauthorizedInvalidBase64            ErrorCode = PrefixUnauthorized + "_INVALID_BASE64"
	ErrCodeUnauthorizedInvalidCredentialsFormat ErrorCode = PrefixUnauthorized + "_INVALID_CREDENTIALS_FORMAT"

	// 4xx - Forbidden errors
	ErrCodeForbiddenAccessDenied        ErrorCode = PrefixForbidden + "_ACCESS_DENIED"
	ErrCodeForbiddenNotGroupMember      ErrorCode = PrefixForbidden + "_NOT_GROUP_MEMBER"
	ErrCodeForbiddenNotMessageRecipient ErrorCode = PrefixForbidden + "_NOT_MESSAGE_RECIPIENT"

	// 4xx - Not Found errors
	ErrCodeNotFoundResourceNotFound     ErrorCode = PrefixNotFound + "_RESOURCE_NOT_FOUND"
	ErrCodeNotFoundUserNotFound         ErrorCode = PrefixNotFound + "_USER_NOT_FOUND"
	ErrCodeNotFoundGroupNotFound        ErrorCode = PrefixNotFound + "_GROUP_NOT_FOUND"
	ErrCodeNotFoundMessageNotFound      ErrorCode = PrefixNotFound + "_MESSAGE_NOT_FOUND"
	ErrCodeNotFoundConversationNotFound ErrorCode = PrefixNotFound + "_CONVERSATION_NOT_FOUND"
	ErrCodeNotFoundSenderNotFound       ErrorCode = PrefixNotFound + "_SENDER_NOT_FOUND"
	ErrCodeNotFoundDestinationNotFound  ErrorCode = PrefixNotFound + "_DESTINATION_NOT_FOUND"

	// 4xx - Conflict errors
	ErrCodeConflictUserAlreadyExists  ErrorCode = PrefixConflict + "_USER_ALREADY_EXISTS"
	ErrCodeConflictGroupAlreadyExists ErrorCode = PrefixConflict + "_GROUP_ALREADY_EXISTS"

	// 5xx - Server Errors
	ErrCodeServerErrorInternalError          ErrorCode = PrefixServerError + "_INTERNAL_ERROR"
	ErrCodeServerErrorSearchFailed           ErrorCode = PrefixServerError + "_SEARCH_FAILED"
	ErrCodeServerErrorConfigNotFound         ErrorCode = PrefixServerError + "_CONFIG_NOT_FOUND"
	ErrCodeServerErrorConfigInvalid          ErrorCode = PrefixServerError + "_CONFIG_INVALID"
	ErrCodeServerErrorConfigValidationFailed ErrorCode = PrefixServerError + "_CONFIG_VALIDATION_FAILED"
)

// AppError represents an application error
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	HTTPStatus int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Details:    make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	e.Details[key] = value
	return e
}

// Predefined errors - 4xx Client Errors
var (
	// Bad Request (400)
	ErrInvalidRequest      = NewAppError(ErrCodeBadRequestInvalidRequest, "Invalid request", http.StatusBadRequest)
	ErrValidationError     = NewAppError(ErrCodeBadRequestValidationError, "Validation error", http.StatusBadRequest)
	ErrMessageEmpty        = NewAppError(ErrCodeBadRequestMessageEmpty, "Message cannot be empty", http.StatusBadRequest)
	ErrMessageTooLong      = NewAppError(ErrCodeBadRequestMessageTooLong, "Message exceeds maximum length", http.StatusBadRequest)
	ErrInvalidRecipient    = NewAppError(ErrCodeBadRequestInvalidRecipient, "Invalid recipient", http.StatusBadRequest)
	ErrSearchQueryRequired = NewAppError(ErrCodeBadRequestSearchQueryRequired, "Query parameter is required", http.StatusBadRequest)
	ErrGroupMemberLimit    = NewAppError(ErrCodeBadRequestGroupMemberLimit, "Group member limit exceeded", http.StatusBadRequest)
	ErrInvalidConversation = NewAppError(ErrCodeBadRequestInvalidConversation, "Invalid conversation", http.StatusBadRequest)

	// Unauthorized (401)
	ErrAuthRequired             = NewAppError(ErrCodeUnauthorizedAuthRequired, "Authorization header required", http.StatusUnauthorized)
	ErrInvalidCredentials       = NewAppError(ErrCodeUnauthorizedInvalidCredentials, "Invalid username or password", http.StatusUnauthorized)
	ErrInvalidAuthFormat        = NewAppError(ErrCodeUnauthorizedInvalidAuthFormat, "Invalid authorization header format. Expected: Basic <base64(username:password)>", http.StatusUnauthorized)
	ErrInvalidBase64            = NewAppError(ErrCodeUnauthorizedInvalidBase64, "Invalid base64 encoding in authorization header", http.StatusUnauthorized)
	ErrInvalidCredentialsFormat = NewAppError(ErrCodeUnauthorizedInvalidCredentialsFormat, "Invalid credentials format. Expected: username:password", http.StatusUnauthorized)

	// Forbidden (403)
	ErrForbiddenAccessDenied = NewAppError(ErrCodeForbiddenAccessDenied, "Access denied", http.StatusForbidden)
	ErrNotGroupMember        = NewAppError(ErrCodeForbiddenNotGroupMember, "User is not a member of this group", http.StatusForbidden)
	ErrNotMessageRecipient   = NewAppError(ErrCodeForbiddenNotMessageRecipient, "User is not the recipient of this message", http.StatusForbidden)

	// Not Found (404)
	ErrNotFound             = NewAppError(ErrCodeNotFoundResourceNotFound, "Resource not found", http.StatusNotFound)
	ErrUserNotFound         = NewAppError(ErrCodeNotFoundUserNotFound, "User not found", http.StatusNotFound)
	ErrGroupNotFound        = NewAppError(ErrCodeNotFoundGroupNotFound, "Group not found", http.StatusNotFound)
	ErrMessageNotFound      = NewAppError(ErrCodeNotFoundMessageNotFound, "Message not found", http.StatusNotFound)
	ErrConversationNotFound = NewAppError(ErrCodeNotFoundConversationNotFound, "Conversation not found", http.StatusNotFound)
	ErrSenderNotFound       = NewAppError(ErrCodeNotFoundSenderNotFound, "Sender not found", http.StatusNotFound)
	ErrDestinationNotFound  = NewAppError(ErrCodeNotFoundDestinationNotFound, "Destination not found", http.StatusNotFound)

	// Conflict (409)
	ErrUserAlreadyExists  = NewAppError(ErrCodeConflictUserAlreadyExists, "User already exists", http.StatusConflict)
	ErrGroupAlreadyExists = NewAppError(ErrCodeConflictGroupAlreadyExists, "Group already exists", http.StatusConflict)
)

// Predefined errors - 5xx Server Errors
var (
	ErrInternalError          = NewAppError(ErrCodeServerErrorInternalError, "An internal error occurred", http.StatusInternalServerError)
	ErrSearchFailed           = NewAppError(ErrCodeServerErrorSearchFailed, "Search failed", http.StatusInternalServerError)
	ErrConfigNotFound         = NewAppError(ErrCodeServerErrorConfigNotFound, "Config file not found", http.StatusInternalServerError)
	ErrConfigInvalid          = NewAppError(ErrCodeServerErrorConfigInvalid, "Invalid config file", http.StatusInternalServerError)
	ErrConfigValidationFailed = NewAppError(ErrCodeServerErrorConfigValidationFailed, "Config validation failed", http.StatusInternalServerError)
)

// ToJSONResponse converts error to common JSON response format
func (e *AppError) ToJSONResponse() map[string]interface{} {
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    string(e.Code),
			"message": e.Message,
		},
	}
	if len(e.Details) > 0 {
		response["error"].(map[string]interface{})["details"] = e.Details
	}
	return response
}

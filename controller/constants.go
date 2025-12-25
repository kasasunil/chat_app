package controller

// API endpoint paths
const (
	EndpointSendMessage          = "/api/v1/sendMessage"
	EndpointAckDelivered         = "/api/v1/ack/delivered"
	EndpointAckRead              = "/api/v1/ack/read"
	EndpointGetMessages          = "/api/v1/conversations/{destinationId}/messages"
	EndpointGetUserConversations = "/api/v1/users/{userId}/conversations"
	EndpointSearchMessages       = "/api/v1/search/{userId}"
	EndpointHealth               = "/health"
)

// Legacy endpoint paths
const (
	LegacyEndpointSendMessage          = "/sendMessage"
	LegacyEndpointAckDelivered         = "/ack/delivered"
	LegacyEndpointAckRead              = "/ack/read"
	LegacyEndpointGetMessages          = "/conversations/{destinationId}/messages"
	LegacyEndpointGetUserConversations = "/users/{userId}/conversations"
	LegacyEndpointSearchMessages       = "/search/{userId}"
)

// HTTP methods
const (
	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodDELETE = "DELETE"
)

// Default values
const (
	DefaultMessageLimit      = 50
	MaxMessageLimit          = 100
	DefaultConversationLimit = 50
	MaxConversationLimit     = 100
)

// Request field names
const (
	FieldSenderID      = "sender_id"
	FieldDestinationID = "destination_id"
	FieldMessage       = "message"
	FieldMessageID     = "message_id"
	FieldUserID        = "user_id"
	FieldQuery         = "query"
	FieldCursor        = "cursor"
	FieldLimit         = "limit"
)

// Response messages
const (
	MsgDeliveryAcknowledged = "Delivery acknowledged"
	MsgReadAcknowledged     = "Read acknowledged"
	MsgHealthOK             = "OK"
)

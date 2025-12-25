package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("==============================================")
	fmt.Println("Chat System Demo - End-to-End Test")
	fmt.Println("==============================================")
	fmt.Println()

	// Wait for server to be ready
	fmt.Println("Waiting for server to be ready...")
	time.Sleep(2 * time.Second)

	// Note: Users and groups are already created by the server
	fmt.Println("\n=== STEP 1: Users and Groups ===")
	fmt.Println("Users: user1 (Alice), user2 (Bob), user3 (Charlie)")
	fmt.Println("Group: group1 (Project Team)")

	fmt.Println("\n=== STEP 2: Send Messages ===")

	message1ID := sendMessage("user1", "user2", "Hello Bob! How are you?")
	time.Sleep(100 * time.Millisecond)
	message2ID := sendMessage("user2", "user1", "Hi Alice! I'm doing great, thanks!")
	time.Sleep(100 * time.Millisecond)
	message3ID := sendMessage("user1", "group1", "Hello team! Let's discuss the project.")
	time.Sleep(100 * time.Millisecond)
	message4ID := sendMessage("user2", "group1", "Sure, I'm ready to discuss.")

	fmt.Println("\n=== STEP 3: Simulate Delivery ACK ===")
	ackDelivered("user2", message1ID)
	time.Sleep(100 * time.Millisecond)
	ackDelivered("user1", message2ID)
	time.Sleep(100 * time.Millisecond)
	ackDelivered("user2", message3ID)
	time.Sleep(100 * time.Millisecond)
	ackDelivered("user3", message3ID)
	time.Sleep(100 * time.Millisecond)
	ackDelivered("user1", message4ID)
	time.Sleep(100 * time.Millisecond)
	ackDelivered("user3", message4ID)

	fmt.Println("\n=== STEP 4: Simulate Read ACK ===")
	ackRead("user2", message1ID)
	time.Sleep(100 * time.Millisecond)
	ackRead("user1", message2ID)
	time.Sleep(100 * time.Millisecond)
	ackRead("user2", message3ID)
	time.Sleep(100 * time.Millisecond)
	ackRead("user3", message3ID)

	fmt.Println("\n=== STEP 5: Fetch Messages with Pagination ===")
	fetchMessages("user1", "user2", "") // Fetch messages in conversation with user2 (as user1)
	time.Sleep(100 * time.Millisecond)
	fetchMessages("user1", "group1", "") // Fetch group messages (as user1)

	fmt.Println("\n=== STEP 6: Fetch Conversation List ===")
	fetchConversations("user1")
	time.Sleep(100 * time.Millisecond)
	fetchConversations("user2")

	fmt.Println("\n=== STEP 7: Search Messages ===")
	searchMessages("user1", "project")
	time.Sleep(100 * time.Millisecond)
	searchMessages("user2", "hello")

	fmt.Println("\n==============================================")
	fmt.Println("Demo completed successfully!")
	fmt.Println("==============================================")
}

func sendMessage(senderID, destID, message string) string {
	reqBody := map[string]string{
		"destination_id": destID,
		"message":        message,
	}
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/sendMessage", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Use basic auth - username from senderID, password from config
	// For demo, we'll use the username as the senderID
	auth := getBasicAuth(senderID)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	messageID := result["message_id"].(string)
	status := result["status"].(string)

	senderName := getSenderName(senderID)
	destName := getDestName(destID)
	fmt.Printf("  [%s → %s] Message sent: %s (ID: %s, Status: %s ✓)\n",
		senderName, destName, message, messageID, status)
	return messageID
}

func ackDelivered(userID, messageID string) {
	reqBody := map[string]string{
		"message_id": messageID,
	}
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/ack/delivered", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	auth := getBasicAuth(userID)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error acking delivery: %v\n", err)
		return
	}
	defer resp.Body.Close()

	userName := getUserName(userID)
	fmt.Printf("  [%s] Delivery ACK for message %s (Status: ✓✓ DELIVERED)\n", userName, messageID)
}

func ackRead(userID, messageID string) {
	reqBody := map[string]string{
		"message_id": messageID,
	}
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/ack/read", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	auth := getBasicAuth(userID)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error acking read: %v\n", err)
		return
	}
	defer resp.Body.Close()

	userName := getUserName(userID)
	fmt.Printf("  [%s] Read ACK for message %s (Status: ✓✓ READ)\n", userName, messageID)
}

func fetchMessages(userID, destID string, cursor string) {
	url := fmt.Sprintf("%s/api/v1/conversations/%s/messages", baseURL, destID)
	if cursor != "" {
		url += "?cursor=" + cursor
	}

	req, _ := http.NewRequest("GET", url, nil)
	auth := getBasicAuth(userID)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching messages: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	messages := result["messages"].([]interface{})
	hasMore := result["has_more"].(bool)

	destName := getDestName(destID)
	fmt.Printf("  Messages for %s (%s):\n", destName, destID)
	for i, msg := range messages {
		m := msg.(map[string]interface{})
		senderName := getSenderName(m["sender_id"].(string))
		status := m["status"].(string)
		statusIcon := getStatusIcon(status)
		fmt.Printf("    %d. [%s] %s %s\n", i+1, senderName, m["message_text"], statusIcon)
	}
	if hasMore {
		fmt.Printf("    ... (more messages available)\n")
	}
}

func fetchConversations(userID string) {
	url := fmt.Sprintf("%s/api/v1/users/%s/conversations", baseURL, userID)

	req, _ := http.NewRequest("GET", url, nil)
	auth := getBasicAuth(userID)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching conversations: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	conversations := result["conversations"].([]interface{})
	userName := getUserName(userID)
	fmt.Printf("  Conversations for %s (%s):\n", userName, userID)
	for i, conv := range conversations {
		c := conv.(map[string]interface{})
		convType := c["conversation_type"].(string)
		destID := c["destination_id"].(string)
		unreadCount := int(c["unread_count"].(float64))

		destName := getDestName(destID)
		lastMsg := ""
		if c["last_message"] != nil {
			lm := c["last_message"].(map[string]interface{})
			lastMsg = lm["message_text"].(string)
		}

		fmt.Printf("    %d. [%s] %s - Last: \"%s\" (Unread: %d)\n",
			i+1, convType, destName, lastMsg, unreadCount)
	}
}

func searchMessages(userID, query string) {
	url := fmt.Sprintf("%s/api/v1/search/%s?query=%s", baseURL, userID, query)

	req, _ := http.NewRequest("GET", url, nil)
	auth := getBasicAuth(userID)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	results := result["results"].([]interface{})
	userName := getUserName(userID)
	fmt.Printf("  Search results for %s (query: \"%s\"): %d results\n", userName, query, len(results))
	for i, msg := range results {
		m := msg.(map[string]interface{})
		senderName := getSenderName(m["sender_id"].(string))
		fmt.Printf("    %d. [%s] %s\n", i+1, senderName, m["message_text"])
	}
}

func getSenderName(userID string) string {
	names := map[string]string{
		"user1": "Alice",
		"user2": "Bob",
		"user3": "Charlie",
	}
	return names[userID]
}

func getUserName(userID string) string {
	return getSenderName(userID)
}

func getDestName(destID string) string {
	if destID == "group1" {
		return "Project Team"
	}
	return getSenderName(destID)
}

func getStatusIcon(status string) string {
	switch status {
	case "SENT":
		return "✓"
	case "DELIVERED":
		return "✓✓"
	case "READ":
		return "✓✓ (blue)"
	default:
		return ""
	}
}

// getBasicAuth returns Basic auth header value for a user
// Maps userID to username/password from config
func getBasicAuth(userID string) string {
	// Map userID to username/password
	// user1 -> user1/password1, user2 -> user2/password2, etc.
	username := userID
	password := fmt.Sprintf("password%s", strings.TrimPrefix(userID, "user"))

	credentials := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return fmt.Sprintf("Basic %s", encoded)
}

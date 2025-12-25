package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// GenerateID generates a unique ID using timestamp and random bytes
func GenerateID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomStr := base64.URLEncoding.EncodeToString(randomBytes)
	return fmt.Sprintf("%d_%s", timestamp, randomStr)
}

// GenerateShortID generates a shorter unique ID
func GenerateShortID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := base64.URLEncoding.EncodeToString(randomBytes)
	return fmt.Sprintf("%d_%s", timestamp, randomStr)
}

// Contains checks if a string slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// SanitizeString removes leading/trailing whitespace and limits length
func SanitizeString(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if maxLen > 0 && len(s) > maxLen {
		s = s[:maxLen]
	}
	return s
}

// ContainsString performs case-insensitive substring search in text
// Returns true if query is found in text (case-insensitive)
func ContainsString(text, query string) bool {
	if len(query) == 0 {
		return false
	}
	textLower := toLowerString(text)
	queryLower := toLowerString(query)
	return containsSubstring(textLower, queryLower)
}

// toLowerString converts a string to lowercase (custom implementation)
func toLowerString(s string) string {
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result = append(result, r+'a'-'A')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// containsSubstring checks if substr exists in text
func containsSubstring(text, substr string) bool {
	if len(substr) > len(text) {
		return false
	}
	for i := 0; i <= len(text)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if text[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

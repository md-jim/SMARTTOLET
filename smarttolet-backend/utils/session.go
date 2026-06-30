package utils

import (
	"fmt"
	"time"
)

var sessions = make(map[string]int)

func GenerateToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func SetSession(token string, userID int) {
	sessions[token] = userID
	fmt.Printf("✅ Session saved - Token: %s, UserID: %d\n", token, userID)
}

func GetUserID(token string) (int, bool) {
	userID, exists := sessions[token]
	fmt.Printf("🔍 Session check - Token: %s, Exists: %v, UserID: %d\n", token, exists, userID)
	return userID, exists
}

// RemoveSession - Delete session
func RemoveSession(token string) {
	delete(sessions, token)
	fmt.Printf("🗑️ Session removed - Token: %s\n", token)
}

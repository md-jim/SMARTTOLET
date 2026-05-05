package utils

import "fmt"
import "time"

var sessions = make(map[string]int)

func GenerateToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func SetSession(token string, userID int) {
	sessions[token] = userID
}

func GetUserID(token string) (int, bool) {
	userID, exists := sessions[token]
	return userID, exists
}

func RemoveSession(token string) {
	delete(sessions, token)
}
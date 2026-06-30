package controllers

import (
	"encoding/json"
	"net/http"
	"smarttolet-backend/config"
	"smarttolet-backend/models"
	"smarttolet-backend/utils"
	"strings"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)

	var existingID int
	config.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingID)
	if existingID != 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email already registered"})
		return
	}

	result, _ := config.DB.Exec("INSERT INTO users (name, email, password, phone) VALUES (?, ?, ?, ?)",
		req.Name, req.Email, req.Password, req.Phone)

	userID, _ := result.LastInsertId()
	token := utils.GenerateToken()
	utils.SetSession(token, int(userID))

	var user models.User
	config.DB.QueryRow("SELECT id, name, email, phone, role, is_active, is_admin, created_at FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.IsAdmin, &user.CreatedAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.LoginResponse{
		Message: "Registration successful",
		Token:   token,
		User:    user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	var user models.User
	err := config.DB.QueryRow("SELECT id, name, email, password, phone, role, is_active, is_admin, created_at FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Role, &user.IsActive, &user.IsAdmin, &user.CreatedAt)

	if err != nil || user.Password != req.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	token := utils.GenerateToken()
	utils.SetSession(token, user.ID)
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		Message: "Login successful",
		Token:   token,
		User:    user,
	})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	userID, exists := utils.GetUserID(token)
	if !exists {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var user models.User
	config.DB.QueryRow("SELECT id, name, email, phone, role, is_active, is_admin, created_at FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.IsAdmin, &user.CreatedAt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]models.User{"user": user})
}

// UpdateProfile - Update user profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	userID, exists := utils.GetUserID(token)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return
	}

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	// Check if email is already taken by another user
	var existingEmail string
	config.DB.QueryRow("SELECT email FROM users WHERE email = ? AND id != ?", req.Email, userID).Scan(&existingEmail)
	if existingEmail != "" {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email already taken"})
		return
	}

	// Update user
	_, err = config.DB.Exec(
		"UPDATE users SET name = ?, phone = ?, email = ? WHERE id = ?",
		req.Name, req.Phone, req.Email, userID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update profile"})
		return
	}

	// Get updated user
	var user models.User
	config.DB.QueryRow("SELECT id, name, email, phone, role, is_active, is_admin, created_at FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.IsAdmin, &user.CreatedAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// ChangePassword - Change user password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	userID, exists := utils.GetUserID(token)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	// Get current user password
	var currentPassword string
	config.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&currentPassword)

	// Verify current password
	if currentPassword != req.CurrentPassword {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Current password is incorrect"})
		return
	}

	// Update password
	_, err = config.DB.Exec(
		"UPDATE users SET password = ? WHERE id = ?",
		req.NewPassword, userID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to change password"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}

// DeleteAccount - Delete user account
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	userID, exists := utils.GetUserID(token)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return
	}

	// Delete user
	_, err := config.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete account"})
		return
	}

	// Remove session
	utils.RemoveSession(token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Account deleted successfully",
	})
}

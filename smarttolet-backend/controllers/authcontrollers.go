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

	// Check if user exists
	var existingID int
	err := config.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingID)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email already registered"})
		return
	}

	// Insert user
	result, err := config.DB.Exec("INSERT INTO users (name, email, password, phone) VALUES (?, ?, ?, ?)",
		req.Name, req.Email, req.Password, req.Phone)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}

	userID, _ := result.LastInsertId()
	token := utils.GenerateToken()
	utils.SetSession(token, int(userID))

	var user models.User
	config.DB.QueryRow("SELECT id, name, email, phone, role, is_active, created_at FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.LoginResponse{
		Message: "✅ Registration successful",
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
	err := config.DB.QueryRow("SELECT id, name, email, password, phone, role, is_active, created_at FROM users WHERE email = ?", req.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	if user.Password != req.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	token := utils.GenerateToken()
	utils.SetSession(token, user.ID)
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		Message: "✅ Login successful",
		Token:   token,
		User:    user,
	})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, exists := utils.GetUserID(token)
	if !exists {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var user models.User
	config.DB.QueryRow("SELECT id, name, email, phone, role, is_active, created_at FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]models.User{"user": user})
}
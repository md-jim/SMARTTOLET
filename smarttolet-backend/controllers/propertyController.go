package controllers

import (
	"encoding/json"
	"net/http"
	"smarttolet-backend/config"
	"smarttolet-backend/models"
	"smarttolet-backend/utils"
	"strings"
)

func AddProperty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, exists := utils.GetUserID(token)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token. Please login again"})
		return
	}

	var req models.AddPropertyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request data"})
		return
	}

	// Insert property
	result, err := config.DB.Exec(`
		INSERT INTO properties (
			user_id, title, description, price, property_type, category,
			division, district, upazila, area, address, bedrooms, bathrooms,
			area_size, furnished, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')`,
		userID, req.Title, req.Description, req.Price, req.PropertyType,
		req.Category, req.Division, req.District, req.Upazila, req.Area,
		req.Address, req.Bedrooms, req.Bathrooms, req.AreaSize, req.Furnished)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add property: " + err.Error()})
		return
	}

	propertyID, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "✅ Property added successfully! Waiting for admin approval.",
		"property_id": propertyID,
	})
}

func GetMyProperties(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, exists := utils.GetUserID(token)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login"})
		return
	}

	rows, err := config.DB.Query(`
		SELECT id, title, price, property_type, category, division, district, 
		       upazila, area, status, created_at 
		FROM properties WHERE user_id = ? ORDER BY created_at DESC`, userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch properties"})
		return
	}
	defer rows.Close()

	var properties []map[string]interface{}
	for rows.Next() {
		var id int
		var title, propertyType, category, division, district, upazila, area, status, createdAt string
		var price float64

		rows.Scan(&id, &title, &price, &propertyType, &category, &division, &district, &upazila, &area, &status, &createdAt)

		properties = append(properties, map[string]interface{}{
			"id":            id,
			"title":         title,
			"price":         price,
			"property_type": propertyType,
			"category":      category,
			"division":      division,
			"district":      district,
			"upazila":       upazila,
			"area":          area,
			"status":        status,
			"created_at":    createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"properties": properties})
}

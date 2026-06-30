package controllers

import (
	"encoding/json"
	"fmt"
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

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, exists := utils.GetUserID(token)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Please login first"})
		return
	}

	var req models.AddPropertyRequest
	json.NewDecoder(r.Body).Decode(&req)

	imagesJSON, _ := json.Marshal(req.Images)

	_, err := config.DB.Exec(`
		INSERT INTO properties (
			user_id, title, description, price, property_type, category,
			division, district, upazila, area, address, bedrooms, bathrooms,
			area_size, furnished, images, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')`,
		userID, req.Title, req.Description, req.Price, req.PropertyType,
		req.Category, req.Division, req.District, req.Upazila, req.Area,
		req.Address, req.Bedrooms, req.Bathrooms, req.AreaSize, req.Furnished, string(imagesJSON))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Property added successfully"})
}

func GetAllProperties(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := config.DB.Query("SELECT id, title, price, property_type, division, district, upazila, area, status FROM properties ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var properties []map[string]interface{}
	for rows.Next() {
		var id int
		var title, propertyType, division, district, upazila, area, status string
		var price float64
		rows.Scan(&id, &title, &price, &propertyType, &division, &district, &upazila, &area, &status)

		properties = append(properties, map[string]interface{}{
			"id": id, "title": title, "price": price,
			"property_type": propertyType, "division": division,
			"district": district, "upazila": upazila, "area": area, "status": status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"properties": properties})
}

func SearchProperties(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	division := r.URL.Query().Get("division")
	district := r.URL.Query().Get("district")
	upazila := r.URL.Query().Get("upazila")
	area := r.URL.Query().Get("area")

	fmt.Println("=== SEARCH DEBUG ===")
	fmt.Println("Division:", division)
	fmt.Println("District:", district)
	fmt.Println("Upazila:", upazila)
	fmt.Println("Area:", area)

	query := "SELECT id, title, price, property_type, division, district, upazila, area, status FROM properties WHERE 1=1"
	var args []interface{}

	if division != "" {
		query += " AND division = ?"
		args = append(args, division)
	}
	if district != "" {
		query += " AND district = ?"
		args = append(args, district)
	}
	if upazila != "" {
		query += " AND upazila = ?"
		args = append(args, upazila)
	}
	if area != "" {
		query += " AND area = ?"
		args = append(args, area)
	}

	query += " ORDER BY created_at DESC"

	fmt.Println("Query:", query)
	fmt.Println("Args:", args)

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var properties []map[string]interface{}
	for rows.Next() {
		var id int
		var title, propertyType, division, district, upazila, area, status string
		var price float64
		rows.Scan(&id, &title, &price, &propertyType, &division, &district, &upazila, &area, &status)

		properties = append(properties, map[string]interface{}{
			"id": id, "title": title, "price": price,
			"property_type": propertyType, "division": division,
			"district": district, "upazila": upazila, "area": area, "status": status,
		})
	}

	fmt.Println("Found:", len(properties))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"properties": properties, "count": len(properties)})
}

// GetPropertyByID - Get single property details
func GetPropertyByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get property ID from URL
	parts := strings.Split(r.URL.Path, "/")
	propertyID := parts[len(parts)-1]

	if propertyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Property ID required"})
		return
	}

	var property models.Property
	var imagesStr string

	err := config.DB.QueryRow(`
		SELECT id, user_id, title, description, price, property_type, category,
		division, district, upazila, area, address, bedrooms, bathrooms,
		area_size, furnished, images, status, created_at
		FROM properties WHERE id = ?`, propertyID).
		Scan(&property.ID, &property.UserID, &property.Title, &property.Description, &property.Price,
			&property.PropertyType, &property.Category, &property.Division, &property.District,
			&property.Upazila, &property.Area, &property.Address, &property.Bedrooms, &property.Bathrooms,
			&property.AreaSize, &property.Furnished, &imagesStr, &property.Status, &property.CreatedAt)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Property not found"})
		return
	}

	property.Images = imagesStr

	// Get owner info
	var owner struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	config.DB.QueryRow("SELECT id, name, email, phone FROM users WHERE id = ?", property.UserID).
		Scan(&owner.ID, &owner.Name, &owner.Email, &owner.Phone)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"property": property,
		"owner":    owner,
	})
}

package controllers

import (
	"encoding/json"
	"net/http"
	"smarttolet-backend/config"
	"strconv"
	"strings"
)

// Get Dashboard Statistics
func GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var totalUsers, totalProperties, pendingProperties, approvedProperties, rejectedProperties, totalRent, totalSale int

	config.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	config.DB.QueryRow("SELECT COUNT(*) FROM properties").Scan(&totalProperties)
	config.DB.QueryRow("SELECT COUNT(*) FROM properties WHERE status = 'pending'").Scan(&pendingProperties)
	config.DB.QueryRow("SELECT COUNT(*) FROM properties WHERE status = 'approved'").Scan(&approvedProperties)
	config.DB.QueryRow("SELECT COUNT(*) FROM properties WHERE status = 'rejected'").Scan(&rejectedProperties)
	config.DB.QueryRow("SELECT COUNT(*) FROM properties WHERE property_type = 'rent'").Scan(&totalRent)
	config.DB.QueryRow("SELECT COUNT(*) FROM properties WHERE property_type = 'sale'").Scan(&totalSale)

	stats := map[string]interface{}{
		"total_users":         totalUsers,
		"total_properties":    totalProperties,
		"pending_properties":  pendingProperties,
		"approved_properties": approvedProperties,
		"rejected_properties": rejectedProperties,
		"total_rent":          totalRent,
		"total_sale":          totalSale,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Get All Users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := config.DB.Query(`
		SELECT id, name, email, phone, role, is_active, is_admin, created_at 
		FROM users ORDER BY created_at DESC`)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email, phone, role, createdAt string
		var isActive, isAdmin bool
		rows.Scan(&id, &name, &email, &phone, &role, &isActive, &isAdmin, &createdAt)

		users = append(users, map[string]interface{}{
			"id": id, "name": name, "email": email, "phone": phone,
			"role": role, "is_active": isActive, "is_admin": isAdmin, "created_at": createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}

// Update User Status (Block/Unblock)
func UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	userID := parts[len(parts)-1]

	var req struct {
		IsActive bool `json:"is_active"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	_, err := config.DB.Exec("UPDATE users SET is_active = ? WHERE id = ?", req.IsActive, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update user status"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "User status updated successfully",
		"is_active": req.IsActive,
	})
}

// Make User Admin
func MakeAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	userID := parts[len(parts)-1]

	var req struct {
		IsAdmin bool `json:"is_admin"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	_, err := config.DB.Exec("UPDATE users SET is_admin = ? WHERE id = ?", req.IsAdmin, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update admin status"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "Admin status updated successfully",
		"is_admin": req.IsAdmin,
	})
}

// Delete User
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	userID := parts[len(parts)-1]

	_, err := config.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

// Get All Properties (Admin)
func GetAllPropertiesAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := config.DB.Query(`
		SELECT id, user_id, title, price, property_type, division, district, area, status, created_at 
		FROM properties ORDER BY created_at DESC`)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch properties"})
		return
	}
	defer rows.Close()

	var properties []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var title, propertyType, division, district, area, status, createdAt string
		var price float64
		rows.Scan(&id, &userID, &title, &price, &propertyType, &division, &district, &area, &status, &createdAt)

		properties = append(properties, map[string]interface{}{
			"id": id, "user_id": userID, "title": title, "price": price,
			"property_type": propertyType, "division": division,
			"district": district, "area": area, "status": status, "created_at": createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"properties": properties})
}

// Update Property Status (Approve/Reject)
func UpdatePropertyStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	propertyID, _ := strconv.Atoi(parts[len(parts)-1])

	var req struct {
		Status string `json:"status"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	_, err := config.DB.Exec("UPDATE properties SET status = ? WHERE id = ?", req.Status, propertyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update property status"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Property status updated successfully",
		"status":  req.Status,
	})
}

// Delete Property (Admin)
func DeletePropertyAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	propertyID := parts[len(parts)-1]

	_, err := config.DB.Exec("DELETE FROM properties WHERE id = ?", propertyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete property"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Property deleted successfully",
	})
}

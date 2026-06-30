package main

import (
	"fmt"
	"log"
	"net/http"
	"smarttolet-backend/config"
	"smarttolet-backend/controllers"
	"smarttolet-backend/middleware"
)

func main() {
	config.ConnectDB()
	defer config.DB.Close()

	// Public routes
	http.HandleFunc("/api/auth/register", middleware.CorsMiddleware(controllers.Register))
	http.HandleFunc("/api/auth/login", middleware.CorsMiddleware(controllers.Login))
	http.HandleFunc("/api/properties/all", middleware.CorsMiddleware(controllers.GetAllProperties))
	http.HandleFunc("/api/properties/search", middleware.CorsMiddleware(controllers.SearchProperties))
	http.HandleFunc("/api/properties/", middleware.CorsMiddleware(controllers.GetPropertyByID))

	// User routes (authentication required)
	http.HandleFunc("/api/profile", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.GetProfile)))
	http.HandleFunc("/api/profile/update", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.UpdateProfile)))
	http.HandleFunc("/api/profile/changepassword", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.ChangePassword)))
	http.HandleFunc("/api/profile/delete", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.DeleteAccount)))
	http.HandleFunc("/api/properties/add", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.AddProperty)))

	// Admin routes
	http.HandleFunc("/api/admin/dashboard", middleware.CorsMiddleware(controllers.GetDashboardStats))
	http.HandleFunc("/api/admin/users", middleware.CorsMiddleware(controllers.GetAllUsers))
	http.HandleFunc("/api/admin/users/status/", middleware.CorsMiddleware(controllers.UpdateUserStatus))
	http.HandleFunc("/api/admin/users/makeadmin/", middleware.CorsMiddleware(controllers.MakeAdmin))
	http.HandleFunc("/api/admin/users/delete/", middleware.CorsMiddleware(controllers.DeleteUser))
	http.HandleFunc("/api/admin/properties", middleware.CorsMiddleware(controllers.GetAllPropertiesAdmin))
	http.HandleFunc("/api/admin/properties/status/", middleware.CorsMiddleware(controllers.UpdatePropertyStatus))
	http.HandleFunc("/api/admin/properties/delete/", middleware.CorsMiddleware(controllers.DeletePropertyAdmin))

	// Home route - MUST BE LAST
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "SmarTOlet API Running!", "status": "online"}`))
	})

	fmt.Println("========================================")
	fmt.Println("🚀 SmarTOlet Server Running")
	fmt.Println("========================================")
	fmt.Println("Server: http://localhost:8080")
	fmt.Println("")
	fmt.Println("📋 User Profile Endpoints:")
	fmt.Println("   GET    /api/profile           - Get profile")
	fmt.Println("   PUT    /api/profile/update    - Update profile")
	fmt.Println("   PUT    /api/profile/changepassword - Change password")
	fmt.Println("   DELETE /api/profile/delete    - Delete account")
	fmt.Println("========================================")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

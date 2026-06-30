package routes

import (
	"net/http"
	"smarttolet-backend/controllers"
	"smarttolet-backend/middleware"
)

func SetupRoutes() {
	// Public routes
	http.HandleFunc("/api/auth/register", middleware.CorsMiddleware(controllers.Register))
	http.HandleFunc("/api/auth/login", middleware.CorsMiddleware(controllers.Login))
	http.HandleFunc("/api/properties/all", middleware.CorsMiddleware(controllers.GetAllProperties))
	http.HandleFunc("/api/properties/search", middleware.CorsMiddleware(controllers.SearchProperties))

	// User routes
	http.HandleFunc("/api/profile", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.GetProfile)))
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
}

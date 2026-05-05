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

	// Protected routes
	http.HandleFunc("/api/profile", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.GetProfile)))

	// Property routes (add these)
	http.HandleFunc("/api/properties/add", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.AddProperty)))
	http.HandleFunc("/api/properties/my", middleware.CorsMiddleware(middleware.AuthMiddleware(controllers.GetMyProperties)))
}

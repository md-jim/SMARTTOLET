package main

import (
	"fmt"
	"log"
	"net/http"
	"smarttolet-backend/config"
	"smarttolet-backend/routes"
)

func main() {
	// Connect to database
	config.ConnectDB()
	defer config.DB.Close()

	// Setup all routes
	routes.SetupRoutes()

	// Home route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "🚀 SmarTOlet API Running!", "status": "online"}`))
	})

	fmt.Println("========================================")
	fmt.Println("🚀 SmarTOlet Authentication Server")
	fmt.Println("========================================")
	fmt.Println("Server: http://localhost:8080")
	fmt.Println("")
	fmt.Println("📋 Public Endpoints:")
	fmt.Println("   POST   /api/auth/register  - Create new account")
	fmt.Println("   POST   /api/auth/login     - Login user")
	fmt.Println("")
	fmt.Println("📋 Protected Endpoints:")
	fmt.Println("   GET    /api/profile        - Get user profile (requires token)")
	fmt.Println("========================================")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	// Using your password: 1234
	dsn := "root:1234@tcp(127.0.0.1:3306)/smarttolet?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ MySQL connection error:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ MySQL ping failed. Error:", err)
	}

	fmt.Println("✅ MySQL connected successfully!")

	// Create users table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		phone VARCHAR(20),
		role VARCHAR(50) DEFAULT 'user',
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("❌ Failed to create users table:", err)
	}
	fmt.Println("✅ Users table ready!")
	// Create properties table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS properties (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    property_type ENUM('rent', 'sale') NOT NULL,
    category VARCHAR(50) NOT NULL,
    division VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    upazila VARCHAR(100) NOT NULL,
    area VARCHAR(100) NOT NULL,
    address TEXT,
    bedrooms INT DEFAULT 0,
    bathrooms INT DEFAULT 0,
    area_size INT DEFAULT 0,
    furnished ENUM('furnished', 'semi_furnished', 'unfurnished') DEFAULT 'unfurnished',
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`)
	if err != nil {
		log.Fatal("❌ Failed to create properties table:", err)
	}
	fmt.Println("✅ Properties table ready!")
}

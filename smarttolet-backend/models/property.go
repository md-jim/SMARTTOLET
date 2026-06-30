package models

type Property struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	PropertyType string  `json:"property_type"`
	Category     string  `json:"category"`
	Division     string  `json:"division"`
	District     string  `json:"district"`
	Upazila      string  `json:"upazila"`
	Area         string  `json:"area"`
	Address      string  `json:"address"`
	Bedrooms     int     `json:"bedrooms"`
	Bathrooms    int     `json:"bathrooms"`
	AreaSize     int     `json:"area_size"`
	Furnished    string  `json:"furnished"`
	Images       string  `json:"images"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

type AddPropertyRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	PropertyType string   `json:"property_type"`
	Category     string   `json:"category"`
	Division     string   `json:"division"`
	District     string   `json:"district"`
	Upazila      string   `json:"upazila"`
	Area         string   `json:"area"`
	Address      string   `json:"address"`
	Bedrooms     int      `json:"bedrooms"`
	Bathrooms    int      `json:"bathrooms"`
	AreaSize     int      `json:"area_size"`
	Furnished    string   `json:"furnished"`
	Images       []string `json:"images"`
}

package entities

// Location represents model for coordinates with Lat:Long
type Location struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"long" binding:"required"`
}

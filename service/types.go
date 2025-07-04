package service

// SchoolInfo represents the structure of school data used by services.
type SchoolInfo struct {
	SchoolID   string `json:"school_id"`
	SchoolName string `json:"school_name"`
	// Add CityName, CityId if needed at the service level
	// CityName string `json:"city_name,omitempty"`
	// CityId   string `json:"city_id,omitempty"`
}

package models

import "time"

// Log represents a single motorcycle dashboard scan entry.
type Log struct {
	ID          string     `json:"id"`
	OdometerKm  float64    `json:"odometer_km"`
	FuelLevel   *float64   `json:"fuel_level,omitempty"` // 0–100 %
	Notes       string     `json:"notes"`
	ImagePath   string     `json:"image_path,omitempty"`
	ParsedByOCR bool       `json:"parsed_by_ocr"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateLogRequest is the JSON body accepted when creating a new log.
type CreateLogRequest struct {
	OdometerKm float64  `json:"odometer_km" binding:"required,min=0"`
	FuelLevel  *float64 `json:"fuel_level,omitempty"`
	Notes      string   `json:"notes"`
}

// OCRResult holds the values extracted from a dashboard image.
type OCRResult struct {
	OdometerKm *float64 `json:"odometer_km,omitempty"`
	FuelLevel  *float64 `json:"fuel_level,omitempty"`
	Confidence float64  `json:"confidence"`
}

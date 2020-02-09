package model

import "time"

// Place represents a cool location.
type Place struct {
	ID          uint `gorm:"primary_key"`
	Name        string
	Description string
	Latitude    float64
	Longitude   float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

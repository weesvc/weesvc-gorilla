package model

import "github.com/jinzhu/gorm"

// Place represents a cool location.
type Place struct {
	gorm.Model

	Name        string
	Description string
	Latitude    float64
	Longitude   float64
}

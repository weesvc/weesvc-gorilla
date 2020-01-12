package migrations

import "github.com/jinzhu/gorm"

// Migration defines a script to be applied to a database.
type Migration struct {
	Number uint `gorm:"primary_key"`
	Name   string

	Forwards func(db *gorm.DB) error `gorm:"-"`
}

// Migrations is a slice of available scripts to be applied to the database.
var Migrations []*Migration

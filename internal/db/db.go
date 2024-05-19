// Package db contains implementations for accessing back-end databases.
package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// Initialize supported dialects
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Database represents the data access object.
type Database struct {
	*gorm.DB
}

// New creates a new instance of the data access object given configuration settings.
func New(config *Config) (*Database, error) {
	db, err := gorm.Open(config.Dialect, config.DatabaseURI)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to %s database", config.Dialect)
	}

	db.LogMode(config.Verbose)

	return &Database{db}, nil
}

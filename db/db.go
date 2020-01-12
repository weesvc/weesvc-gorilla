package db

import (
	"github.com/jinzhu/gorm"
	// Notes the database dialect to be used
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

// Database represents the data access object.
type Database struct {
	*gorm.DB
}

// New creates a new instance of the data access object given configuration settings.
func New(config *Config) (*Database, error) {
	db, err := gorm.Open("sqlite3", config.DatabaseURI)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	db.LogMode(config.Verbose)

	return &Database{db}, nil
}

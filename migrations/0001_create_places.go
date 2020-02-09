package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var createPlacesMigration0001 = &Migration{
	Number: 1,
	Name:   "Create places",
	Forwards: func(db *gorm.DB) error {
		const createPlacesSQL = `
			CREATE TABLE places(
				id INTEGER PRIMARY KEY,
				name TEXT UNIQUE NOT NULL,
				description TEXT,
				latitude REAL,
				longitude REAL,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL
			);
		`
		err := db.Exec(createPlacesSQL).Error
		return errors.Wrap(err, "unable to create places table")
	},
}

func init() {
	Migrations = append(Migrations, createPlacesMigration0001)
}

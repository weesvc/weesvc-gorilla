package db

import (
	"github.com/pkg/errors"

	"github.com/weesvc/weesvc-gorilla/model"
)

// GetPlaces retrieves all available places from the database.
func (db *Database) GetPlaces() ([]*model.Place, error) {
	var places []*model.Place
	return places, errors.Wrap(db.Find(&places).Error, "unable to find places")
}

// GetPlaceByID retrieves a single place given its identifier.
func (db *Database) GetPlaceByID(id uint) (*model.Place, error) {
	var place model.Place
	return &place, errors.Wrap(db.First(&place, id).Error, "unable to get place")
}

// CreatePlace add the provided place to the database.
func (db *Database) CreatePlace(place *model.Place) error {
	return errors.Wrap(db.Create(place).Error, "unable to create place")
}

// UpdatePlace updates the existing place in the database.
func (db *Database) UpdatePlace(place *model.Place) error {
	return errors.Wrap(db.Save(place).Error, "unable to update place")
}

// DeletePlaceByID removes a single place from the database given its identifier.
func (db *Database) DeletePlaceByID(id uint) error {
	return errors.Wrap(db.Delete(&model.Place{}, id).Error, "unable to delete place")
}

package app

import "github.com/weesvc/weesvc-gorilla/model"

// GetPlaces returns available places.
func (ctx *Context) GetPlaces() ([]*model.Place, error) {
	return ctx.Database.GetPlaces()
}

// GetPlaceByID returns the place specified by the provided identifier.
func (ctx *Context) GetPlaceByID(id uint) (*model.Place, error) {
	place, err := ctx.Database.GetPlaceByID(id)
	if err != nil {
		return nil, err
	}

	return place, nil
}

// CreatePlace persists the provided place.
func (ctx *Context) CreatePlace(place *model.Place) error {
	if err := ctx.validatePlace(place); err != nil {
		return err
	}

	return ctx.Database.CreatePlace(place)
}

const maxPlaceNameLength = 100

func (ctx *Context) validatePlace(place *model.Place) *ValidationError {
	if len(place.Name) > maxPlaceNameLength {
		return &ValidationError{"name is too long"}
	}

	return nil
}

// UpdatePlace saves changes made to the provided place.
func (ctx *Context) UpdatePlace(place *model.Place) error {
	if err := ctx.validatePlace(place); err != nil {
		return err
	}

	return ctx.Database.UpdatePlace(place)
}

// DeletePlaceByID removes the place from storage given the identifier.
func (ctx *Context) DeletePlaceByID(id uint) error {
	_, err := ctx.GetPlaceByID(id)
	if err != nil {
		return err
	}

	return ctx.Database.DeletePlaceByID(id)
}

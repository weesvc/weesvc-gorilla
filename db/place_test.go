package db

import (
	"context"
	"log"
	"testing"

	"github.com/weesvc/weesvc-gorilla/model"
	"github.com/weesvc/weesvc-gorilla/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_GetPlaces(t *testing.T) {
	t.Parallel()
	placeDB := setupDatabase(t)

	places, err := placeDB.GetPlaces()
	assert.NoError(t, err)
	assert.Equal(t, 10, len(places))
}

func TestDatabase_GetPlaceByID(t *testing.T) {
	t.Parallel()
	placeDB := setupDatabase(t)

	fetchID := uint(6)
	place, err := placeDB.GetPlaceByID(fetchID)
	if assert.NoError(t, err) {
		assert.Equal(t, fetchID, place.ID)
		assert.Equal(t, "MIA", place.Name)
		assert.Equal(t, "Miami International Airport, FL, USA", place.Description)
		assert.Equal(t, 25.79516, place.Latitude)
		assert.Equal(t, -80.27959, place.Longitude)
		assert.NotNil(t, place.CreatedAt)
		assert.NotNil(t, place.UpdatedAt)
	}
}

func TestDatabase_CreatePlace(t *testing.T) {
	t.Parallel()
	placeDB := setupDatabase(t)

	newPlace := &model.Place{
		ID:          20,
		Name:        "Kerid Crater",
		Description: "Kerid Crater, Iceland",
		Latitude:    64.04126,
		Longitude:   -20.88530,
	}
	err := placeDB.CreatePlace(newPlace)
	if assert.NoError(t, err) {
		// Verify our inserted place
		created, err := placeDB.GetPlaceByID(newPlace.ID)
		if assert.NoError(t, err) {
			assert.Equal(t, newPlace.ID, created.ID)
			assert.Equal(t, newPlace.Name, created.Name)
			assert.Equal(t, newPlace.Description, created.Description)
			assert.Equal(t, newPlace.Latitude, created.Latitude)
			assert.Equal(t, newPlace.Longitude, created.Longitude)
			assert.NotNil(t, created.CreatedAt)
			assert.NotNil(t, created.UpdatedAt)
		}
	}
}

func TestDatabase_UpdatePlace(t *testing.T) {
	t.Parallel()
	placeDB := setupDatabase(t)

	original, err := placeDB.GetPlaceByID(7)
	if assert.NoError(t, err) {
		changes := &model.Place{
			ID:          original.ID,
			Name:        "The Alamo",
			Description: "The Alamo, San Antonio, TX, USA",
			Latitude:    29.42590,
			Longitude:   -98.48625,
		}
		if assert.NoError(t, placeDB.UpdatePlace(changes)) {
			// Verify the updated place
			updated, err := placeDB.GetPlaceByID(original.ID)
			if assert.NoError(t, err) {
				assert.Equal(t, original.ID, updated.ID)
				assert.Equal(t, changes.Name, updated.Name)
				assert.Equal(t, changes.Description, updated.Description)
				assert.Equal(t, changes.Latitude, updated.Latitude)
				assert.Equal(t, changes.Longitude, updated.Longitude)
				assert.Equal(t, original.CreatedAt, updated.CreatedAt)
				assert.NotEqual(t, original.UpdatedAt, updated.UpdatedAt)
			}
		}
	}
}

func TestDatabase_DeletePlaceByID(t *testing.T) {
	t.Parallel()
	placeDB := setupDatabase(t)

	deleteID := uint(1)
	_, err := placeDB.GetPlaceByID(deleteID)
	if assert.NoError(t, err) {
		if assert.NoError(t, placeDB.DeletePlaceByID(deleteID)) {
			// Verify no longer retrievable
			_, err = placeDB.GetPlaceByID(deleteID)
			assert.EqualError(t, err, "unable to get place: record not found")
		}
	}
}

// setupDatabase creates an isolated `Database` instance backed by a Postgres Testcontainer.
func setupDatabase(t *testing.T) *Database {
	ctx := context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	placeDB, err := New(&Config{
		DatabaseURI: pgContainer.ConnectionString,
		Dialect:     "postgres",
		Verbose:     true,
	})
	if err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	return placeDB
}

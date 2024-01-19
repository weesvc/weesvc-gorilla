package db

import (
	"context"
	"log"
	"testing"

	"github.com/weesvc/weesvc-gorilla/model"
	"github.com/weesvc/weesvc-gorilla/testhelpers"

	"github.com/stretchr/testify/assert"
)

// Here we're testing each CRUD method independently, but we're using a fresh container
// for each test. This _may_ be ok, but could be a bit of overhead for CICD pipelines.
//
// With this method, we've refactored some container work into the `testhelpers` package.

func TestDatabase_GetPlaces(t *testing.T) {
	t.Parallel()
	placeDb := setupDatabase(t)

	places, err := placeDb.GetPlaces()
	assert.NoError(t, err)
	assert.Equal(t, 10, len(places))
}

func TestDatabase_GetPlaceByID(t *testing.T) {
	t.Parallel()
	placeDb := setupDatabase(t)

	fetchId := uint(6)
	place, err := placeDb.GetPlaceByID(fetchId)
	if assert.NoError(t, err) {
		assert.Equal(t, fetchId, place.ID)
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
	placeDb := setupDatabase(t)

	newPlace := &model.Place{
		ID:          20,
		Name:        "Kerid Crater",
		Description: "Kerid Crater, Iceland",
		Latitude:    64.04126,
		Longitude:   -20.88530,
	}
	err := placeDb.CreatePlace(newPlace)
	if assert.NoError(t, err) {
		// Verify our inserted place
		created, err := placeDb.GetPlaceByID(newPlace.ID)
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
	placeDb := setupDatabase(t)

	original, err := placeDb.GetPlaceByID(7)
	if assert.NoError(t, err) {
		changes := &model.Place{
			ID:          original.ID,
			Name:        "The Alamo",
			Description: "The Alamo, San Antonio, TX, USA",
			Latitude:    29.42590,
			Longitude:   -98.48625,
		}
		if assert.NoError(t, placeDb.UpdatePlace(changes)) {
			// Verify the updated place
			updated, err := placeDb.GetPlaceByID(original.ID)
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
	placeDb := setupDatabase(t)

	deleteID := uint(1)
	_, err := placeDb.GetPlaceByID(deleteID)
	if assert.NoError(t, err) {
		if assert.NoError(t, placeDb.DeletePlaceByID(deleteID)) {
			// Verify no longer retrievable
			_, err = placeDb.GetPlaceByID(deleteID)
			assert.EqualError(t, err, "unable to get place: record not found")
		}
	}
}

func setupDatabase(t *testing.T) *Database {
	ctx := context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	placeDb, err := New(&Config{
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

	return placeDb
}

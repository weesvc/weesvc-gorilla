package db

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/weesvc/weesvc-gorilla/model"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Here we're testing each CRUD method sequentially using a single container for the test.
// This is a bit more of a brute force test, but fastest to "market."

// TestDatabase provided as an example where we do ALL THE THINGS in a single testcase.
func TestDatabase(t *testing.T) {
	// _Real_ test is provided with `place_test.go`...this is just for learning.
	t.Skip()

	t.Parallel()
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if terr := pgContainer.Terminate(ctx); terr != nil {
			t.Fatalf("failed to terminate pgContainer: %s", terr)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	placeDB, err := New(&Config{
		DatabaseURI: connStr,
		Dialect:     "postgres",
		Verbose:     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check 1: Retrieve all places
	places, err := placeDB.GetPlaces()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(places), 10)

	// Check 2: Add a new place
	newPlace := &model.Place{
		ID:          20,
		Name:        "Kerid Crater",
		Description: "Kerid Crater, Iceland",
		Latitude:    64.04126,
		Longitude:   -20.88530,
	}
	err = placeDB.CreatePlace(newPlace)
	if err != nil {
		t.Fatal(err)
	}

	// Check 3: Retrieve the new place
	place, err := placeDB.GetPlaceByID(newPlace.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, place.ID, newPlace.ID)
	assert.Equal(t, place.Name, newPlace.Name)
	assert.Equal(t, place.Description, newPlace.Description)
	assert.Equal(t, place.Latitude, newPlace.Latitude)
	assert.Equal(t, place.Longitude, newPlace.Longitude)
	assert.NotNil(t, place.CreatedAt)
	assert.NotNil(t, place.UpdatedAt)

	// Check 4: Update a place
	updatedDescription := "UPDATED"
	err = placeDB.UpdatePlace(&model.Place{ID: newPlace.ID, Description: updatedDescription})
	if err != nil {
		t.Fatal(err)
	}
	place, _ = placeDB.GetPlaceByID(newPlace.ID)
	assert.Equal(t, updatedDescription, place.Description)
	assert.Greater(t, place.UpdatedAt, place.CreatedAt)

	// Check 5: Delete the place
	err = placeDB.DeletePlaceByID(newPlace.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = placeDB.GetPlaceByID(newPlace.ID)
	assert.EqualError(t, err, "unable to get place: record not found")
}

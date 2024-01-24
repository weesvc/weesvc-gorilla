package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/weesvc/weesvc-gorilla/model"
	"github.com/weesvc/weesvc-gorilla/testhelpers"
)

// Here we're testing each CRUD method independently, but we're using a single container
// for the entire test suite. This _may_ be ok, but you need to be wary of what is changing
// in the database as suite tests may collide.
//
// With this method, we've refactored some container work into the `testhelpers` package.

// TestDatabase_TestSuite executes the suite of tests.
func TestDatabase_TestSuite(t *testing.T) {
	// _Real_ test is provided with `place_test.go`...this is just for learning.
	t.Skip()

	t.Parallel()
	suite.Run(t, new(PlaceTestSuite))
}

// PlaceTestSuite contains shared state amongst all test(s) within the suite.
type PlaceTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	ctx         context.Context
	placeDB     *Database // Reference to our test fixture
}

// SetupSuite executes prior to any test(s) in order to prepare the shared state.
func (suite *PlaceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	placeDB, err := New(&Config{
		DatabaseURI: pgContainer.ConnectionString,
		Dialect:     "postgres",
		Verbose:     true,
	})
	if err != nil {
		log.Fatal(err)
	}
	suite.placeDB = placeDB
}

// TearDownSuite executes cleanup after all test(s) have run.
func (suite *PlaceTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *PlaceTestSuite) Test_GetPlaces() {
	places, err := suite.placeDB.GetPlaces()
	if assert.NoError(suite.T(), err) {
		assert.Greater(suite.T(), len(places), 0)
	}
}

func (suite *PlaceTestSuite) Test_GetPlaceByID() {
	fetchID := uint(6)
	place, err := suite.placeDB.GetPlaceByID(fetchID)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), fetchID, place.ID)
		assert.Equal(suite.T(), "MIA", place.Name)
		assert.Equal(suite.T(), "Miami International Airport, FL, USA", place.Description)
		assert.Equal(suite.T(), 25.79516, place.Latitude)
		assert.Equal(suite.T(), -80.27959, place.Longitude)
		assert.NotNil(suite.T(), place.CreatedAt)
		assert.NotNil(suite.T(), place.UpdatedAt)
	}
}

func (suite *PlaceTestSuite) Test_CreatePlace() {
	newPlace := &model.Place{
		ID:          20,
		Name:        "Kerid Crater",
		Description: "Kerid Crater, Iceland",
		Latitude:    64.04126,
		Longitude:   -20.88530,
	}
	err := suite.placeDB.CreatePlace(newPlace)
	if assert.NoError(suite.T(), err) {
		// Verify our inserted newPlace
		created, err := suite.placeDB.GetPlaceByID(newPlace.ID)
		if assert.NoError(suite.T(), err) {
			assert.Equal(suite.T(), newPlace.ID, created.ID)
			assert.Equal(suite.T(), newPlace.Name, created.Name)
			assert.Equal(suite.T(), newPlace.Description, created.Description)
			assert.Equal(suite.T(), newPlace.Latitude, created.Latitude)
			assert.Equal(suite.T(), newPlace.Longitude, created.Longitude)
			assert.NotNil(suite.T(), created.CreatedAt)
			assert.NotNil(suite.T(), created.UpdatedAt)
		}
	}
}

func (suite *PlaceTestSuite) Test_UpdatePlace() {
	original, err := suite.placeDB.GetPlaceByID(7)
	if assert.NoError(suite.T(), err) {
		changes := &model.Place{
			ID:          original.ID,
			Name:        "The Alamo",
			Description: "The Alamo, San Antonio, TX, USA",
			Latitude:    29.42590,
			Longitude:   -98.48625,
		}
		if assert.NoError(suite.T(), suite.placeDB.UpdatePlace(changes)) {
			// Verify the updated place
			updated, err := suite.placeDB.GetPlaceByID(original.ID)
			if assert.NoError(suite.T(), err) {
				assert.Equal(suite.T(), original.ID, updated.ID)
				assert.Equal(suite.T(), changes.Name, updated.Name)
				assert.Equal(suite.T(), changes.Description, updated.Description)
				assert.Equal(suite.T(), changes.Latitude, updated.Latitude)
				assert.Equal(suite.T(), changes.Longitude, updated.Longitude)
				assert.Equal(suite.T(), original.CreatedAt, updated.CreatedAt)
				assert.NotEqual(suite.T(), original.UpdatedAt, updated.UpdatedAt)
			}
		}
	}
}

func (suite *PlaceTestSuite) Test_DeletePlaceByID() {
	deleteID := uint(1)
	_, err := suite.placeDB.GetPlaceByID(deleteID)
	if assert.NoError(suite.T(), err) {
		if assert.NoError(suite.T(), suite.placeDB.DeletePlaceByID(deleteID)) {
			// Verify no longer retrievable
			_, err = suite.placeDB.GetPlaceByID(deleteID)
			assert.EqualError(suite.T(), err, "unable to get place: record not found")
		}
	}
}

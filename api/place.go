package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/weesvc/weesvc-gorilla/app"
	"github.com/weesvc/weesvc-gorilla/model"
)

func (a *API) getPlaces(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	places, err := ctx.GetPlaces()
	if err != nil {
		return err
	}

	data, err := json.Marshal(places)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

type createPlaceInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type createPlaceResponse struct {
	ID uint `json:"id"`
}

func (a *API) createPlace(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input createPlaceInput

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	place := &model.Place{
		Name:        input.Name,
		Description: input.Description,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
	}

	if err := ctx.CreatePlace(place); err != nil {
		return err
	}

	data, err := json.Marshal(&createPlaceResponse{ID: place.ID})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) getPlaceByID(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIDFromRequest(r)
	place, err := ctx.GetPlaceByID(id)
	if err != nil {
		return handleError(w, r, err)
	}

	data, err := json.Marshal(place)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

type updatePlaceInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
}

func (a *API) updatePlaceByID(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIDFromRequest(r)

	var input updatePlaceInput

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	existingPlace, err := ctx.GetPlaceByID(id)
	if err != nil {
		return handleError(w, r, err)
	}

	if input.Name != nil {
		existingPlace.Name = *input.Name
	}
	if input.Description != nil {
		existingPlace.Description = *input.Description
	}
	if input.Latitude != nil {
		existingPlace.Latitude = *input.Latitude
	}
	if input.Longitude != nil {
		existingPlace.Longitude = *input.Longitude
	}

	err = ctx.UpdatePlace(existingPlace)
	if err != nil {
		return err
	}

	data, err := json.Marshal(existingPlace)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) deletePlaceByID(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	id := getIDFromRequest(r)
	err := ctx.DeletePlaceByID(id)

	if err != nil {
		return handleError(w, r, err)
	}

	return &app.UserError{StatusCode: http.StatusOK, Message: "removed"}
}

func getIDFromRequest(r *http.Request) uint {
	vars := mux.Vars(r)
	id := vars["id"]

	intID, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return 0
	}

	return uint(intID)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) error {
	if strings.Contains(err.Error(), "record not found") {
		http.NotFound(w, r)
		return nil
	}
	return err
}

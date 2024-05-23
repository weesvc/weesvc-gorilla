package handlers

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/weesvc/weesvc-gorilla/internal/app"
	"github.com/weesvc/weesvc-gorilla/internal/server/views"
)

type PlacesHandler struct {
	Application *app.App
	Service     *app.Context
}

func NewPlacesHandler(app *app.App) *PlacesHandler {
	return &PlacesHandler{Application: app, Service: app.NewContext()}
}

func (h *PlacesHandler) GetPlaces(w http.ResponseWriter, r *http.Request) {
	places, err := h.Service.GetPlaces()
	if err != nil {
		panic(err)
	}

	templ.Handler(views.Layout(views.Places(places))).ServeHTTP(w, r)
}

func (h *PlacesHandler) SearchPlaces(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	places, err := h.Service.SearchPlaces(search)
	if err != nil {
		panic(err)
	}

	templ.Handler(views.PlaceRows(places)).ServeHTTP(w, r)
}

func (h *PlacesHandler) GetPlaceByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	place, err := h.Service.GetPlaceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templ.Handler(views.Layout(views.PlaceDetailsPage(place))).ServeHTTP(w, r)
}

func (h *PlacesHandler) GetPlaceDetails(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	place, err := h.Service.GetPlaceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templ.Handler(views.PlaceDetails(place)).ServeHTTP(w, r)
}

func (h *PlacesHandler) GetPlaceEditor(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	place, err := h.Service.GetPlaceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templ.Handler(views.PlaceEditor(place)).ServeHTTP(w, r)
}

func (h *PlacesHandler) UpdatePlaceByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	place, err := h.Service.GetPlaceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	place.Name = r.FormValue("name")
	place.Description = r.FormValue("description")
	place.Latitude, err = strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	place.Longitude, err = strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = h.Service.UpdatePlace(place)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templ.Handler(views.PlaceDetails(place)).ServeHTTP(w, r)
}

func (h *PlacesHandler) DeletePlaceByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = h.Service.DeletePlaceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func getIDFromRequest(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	id := vars["id"]

	intID, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		return 0, err
	}

	return uint(intID), nil
}

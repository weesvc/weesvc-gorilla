package handlers

import (
	"github.com/a-h/templ"
	"github.com/weesvc/weesvc-gorilla/internal/app"
	"github.com/weesvc/weesvc-gorilla/internal/server/views"
	"net/http"
)

type PlacesHandler struct {
	Application *app.App
}

func NewPlacesHandler(app *app.App) *PlacesHandler {
	return &PlacesHandler{app}
}

func (h *PlacesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svc := h.Application.NewContext()
	places, err := svc.GetPlaces()
	if err != nil {
		panic(err)
	}

	templ.Handler(views.Layout(views.Places(places))).ServeHTTP(w, r)
}

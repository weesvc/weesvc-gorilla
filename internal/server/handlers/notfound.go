package handlers

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/weesvc/weesvc-gorilla/internal/server/views"
)

type NotFoundHandler struct{}

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templ.Handler(
		views.Layout(views.NotFound("Not found.")),
		templ.WithStatus(http.StatusNotFound),
	).ServeHTTP(w, r)
}

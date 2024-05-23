package handlers

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/weesvc/weesvc-gorilla/internal/server/views"
)

type WelcomeHandler struct{}

func NewWelcomeHandler() *WelcomeHandler {
	return &WelcomeHandler{}
}

func (h *WelcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templ.Handler(views.Layout(views.Index("StLGo"))).ServeHTTP(w, r)
}

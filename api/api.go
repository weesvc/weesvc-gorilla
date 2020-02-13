package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/weesvc/weesvc-gorilla/app"
)

type statusCodeRecorder struct {
	http.ResponseWriter
	http.Hijacker
	StatusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// API represents the context for the public interface.
type API struct {
	App    *app.App
	Config *Config
}

// New creates a new API instance.
func New(a *app.App) (api *API, err error) {
	api = &API{App: a}
	api.Config, err = initConfig()
	if err != nil {
		return nil, err
	}
	return api, nil
}

// Init is where we define the routes our API will support.
func (a *API) Init(r *mux.Router) {
	r.Handle("/hello", a.handler(a.helloHandler))

	// place methods
	placesRouter := r.PathPrefix("/places").Subrouter()
	placesRouter.Handle("", a.handler(a.getPlaces)).Methods("GET")
	placesRouter.Handle("", a.handler(a.createPlace)).Methods("POST")
	placesRouter.Handle("/{id:[0-9]+}", a.handler(a.getPlaceByID)).Methods("GET")
	placesRouter.Handle("/{id:[0-9]+}", a.handler(a.updatePlaceByID)).Methods("PATCH")
	placesRouter.Handle("/{id:[0-9]+}", a.handler(a.deletePlaceByID)).Methods("DELETE")
}

func (a *API) handler(f func(*app.Context, http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024)

		beginTime := time.Now()
		traceID, _ := uuid.NewUUID()

		hijacker, _ := w.(http.Hijacker)
		w = &statusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		ctx := a.App.NewContext().WithRemoteAddress(a.addressForRequest(r)).WithTraceID(traceID)

		defer func() {
			statusCode := w.(*statusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)

			ctx.Logger.WithFields(logrus.Fields{
				"duration":       duration,
				"status_code":    statusCode,
				"remote_address": ctx.RemoteAddress,
				"trace_id":       ctx.TraceID,
			}).Info(r.Method + " " + r.URL.RequestURI())
		}()

		defer func() {
			if r := recover(); r != nil {
				ctx.Logger.Error(fmt.Errorf("%v: %s", r, debug.Stack()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		w.Header().Set("Content-Type", "application/json")

		if err := f(ctx, w, r); err != nil {
			if verr, ok := err.(*app.ValidationError); ok {
				data, err := json.Marshal(verr)
				if err == nil {
					w.WriteHeader(http.StatusBadRequest)
					_, err = w.Write(data)
				}

				if err != nil {
					ctx.Logger.Error(err)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			} else if uerr, ok := err.(*app.UserError); ok {
				data, err := json.Marshal(uerr)
				if err == nil {
					w.WriteHeader(uerr.StatusCode)
					_, err = w.Write(data)
				}

				if err != nil {
					ctx.Logger.Error(err)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			} else {
				ctx.Logger.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}
	})
}

func (a *API) helloHandler(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte(fmt.Sprintf(`{"hello":"world","remote_address":%q,"trace_id":%q}`, ctx.RemoteAddress, ctx.TraceID)))
	return err
}

func (a *API) addressForRequest(r *http.Request) string {
	addr := r.RemoteAddr
	return addr[:strings.LastIndex(addr, ":")]
}

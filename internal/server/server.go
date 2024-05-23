package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/weesvc/weesvc-gorilla/internal/server/api"
	"github.com/weesvc/weesvc-gorilla/internal/server/handlers"

	"github.com/weesvc/weesvc-gorilla/internal/config"

	gorilla "github.com/gorilla/handlers"

	"github.com/weesvc/weesvc-gorilla/internal/app"
)

//go:embed assets
var wwwroot embed.FS

// StartServer sets up HTTP routes and starts the application server.
func StartServer(config *config.Config) {
	router := mux.NewRouter()

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(recoveryMiddleware)

	// Initialize the application
	svc, err := app.New(config)
	if err != nil {
		panic(err)
	}
	restapi := api.New(svc)

	// Bind endpoints to handlers
	restapi.Init(router.PathPrefix("/api").Subrouter())

	router.Handle("/", handlers.NewWelcomeHandler()).Methods("GET")

	ph := handlers.NewPlacesHandler(svc)
	router.HandleFunc("/places", ph.GetPlaces).Methods("GET")
	// TODO router.HandleFunc("/places", ph.CreatePlace).Methods("POST")
	router.HandleFunc("/places/{id:[0-9]+}", ph.GetPlaceByID).Methods("GET")
	router.HandleFunc("/places/{id:[0-9]+}/edit", ph.GetPlaceEditor).Methods("GET")
	router.HandleFunc("/places/{id:[0-9]+}/cancel", ph.GetPlaceDetails).Methods("GET")
	router.HandleFunc("/places/{id:[0-9]+}", ph.UpdatePlaceByID).Methods("PUT")
	router.HandleFunc("/places/{id:[0-9]+}", ph.DeletePlaceByID).Methods("DELETE")
	router.HandleFunc("/places/search", ph.SearchPlaces).Methods("POST")

	router.PathPrefix("/assets/").Handler(handlers.NewStaticHandler(config, wwwroot, "assets"))

	router.NotFoundHandler = handlers.NewNotFoundHandler()

	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", config.Port),
		Handler:     gorilla.CombinedLoggingHandler(os.Stdout, router),
		ReadTimeout: 2 * time.Minute,
	}

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := s.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("server shutdown complete")
		} else if err != nil {
			slog.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	slog.Info(fmt.Sprintf("serving api at http://localhost:%d", config.Port))
	<-killSig

	slog.Info("server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		slog.Error("error shutting down", slog.Any("err", err))
		//nolint:gocritic
		os.Exit(1)
	}
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gorilla.RecoveryHandler(gorilla.PrintRecoveryStack(true))
		next.ServeHTTP(w, r)
	})
}

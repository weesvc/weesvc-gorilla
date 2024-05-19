package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/weesvc/weesvc-gorilla/internal/config"

	"github.com/pkg/errors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/weesvc/weesvc-gorilla/internal/api"
	"github.com/weesvc/weesvc-gorilla/internal/app"
)

// StartServer sets up HTTP routes and starts the application server.
func StartServer(config *config.Config) error {
	a, err := app.New(config)
	if err != nil {
		return err
	}

	api := api.New(a)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		slog.Info("signal caught. shutting down...")
		cancel()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		serveAPI(ctx, a, api)
	}()

	wg.Wait()
	return nil
}

func serveAPI(ctx context.Context, app *app.App, api *api.API) {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"}),
	)

	router := mux.NewRouter()
	api.Init(router.PathPrefix("/api").Subrouter())

	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", app.Config.Port),
		Handler:     cors(router),
		ReadTimeout: 2 * time.Minute,
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		//nolint:contextcheck
		if err := s.Shutdown(context.Background()); err != nil {
			slog.Error(err.Error())
		}
		close(done)
	}()

	slog.Info(fmt.Sprintf("serving api at http://127.0.0.1:%d", app.Config.Port))
	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server shutdown complete")
	} else if err != nil {
		slog.Error("server error", slog.Any("err", err))
		os.Exit(1)
	}

	<-done
}

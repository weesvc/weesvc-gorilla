// Package app provides the core service implementations.
package app

import (
	"github.com/sirupsen/logrus"

	"github.com/weesvc/weesvc-gorilla/db"
)

// App defines the main application state and behaviors.
type App struct {
	Database *db.Database
}

// NewContext creates context to bind to an incoming request.
func (a *App) NewContext() *Context {
	return &Context{
		Logger:   logrus.New(),
		Database: a.Database,
	}
}

// New constructs a new instance of the application.
func New() (app *App, err error) {
	app = &App{}

	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}

	app.Database, err = db.New(dbConfig)
	if err != nil {
		return nil, err
	}

	return app, err
}

// Close ensures cleanup of resources.
func (a *App) Close() error {
	return a.Database.Close()
}

// ValidationError defines a data-centric error.
type ValidationError struct {
	Message string `json:"message"`
}

// Error creates an error message from the underlying ValidationError.
func (e *ValidationError) Error() string {
	return e.Message
}

// UserError defines a user-centric error.
type UserError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

// Error creates an error message from the underlying UserError.
func (e *UserError) Error() string {
	return e.Message
}

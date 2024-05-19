package app

import (
	"log/slog"

	"github.com/google/uuid"

	"github.com/weesvc/weesvc-gorilla/internal/db"
)

// Context provides for a request-scoped context.
type Context struct {
	Logger        *slog.Logger
	RemoteAddress string
	TraceID       uuid.UUID
	Database      *db.Database
}

// WithLogger associates the provided logger to the request context.
func (ctx *Context) WithLogger(logger *slog.Logger) *Context {
	ret := *ctx
	ret.Logger = logger
	return &ret
}

// WithRemoteAddress associates the provided address to the request context.
func (ctx *Context) WithRemoteAddress(address string) *Context {
	ret := *ctx
	ret.RemoteAddress = address
	return &ret
}

// WithTraceID associates the provided UUID to the request context.
func (ctx *Context) WithTraceID(uuid uuid.UUID) *Context {
	ret := *ctx
	ret.TraceID = uuid
	return &ret
}

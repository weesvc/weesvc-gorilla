package app

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/weesvc/weesvc-gorilla/db"
)

// Context provides for a request-scoped context.
type Context struct {
	Logger        logrus.FieldLogger
	RemoteAddress string
	TraceID       uuid.UUID
	Database      *db.Database
}

// WithLogger associates the provided logger to the request context.
func (ctx *Context) WithLogger(logger logrus.FieldLogger) *Context {
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

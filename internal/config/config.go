// Package config provides the central configuration definition for application and services.
package config

// Config provides database connection information.
type Config struct {
	Dialect     string
	DatabaseURI string
	Verbose     bool
	// The port to bind the web application api to
	Port int
}

// NewConfig returns an initialized, but empty, configuration object.
func NewConfig() *Config {
	return &Config{}
}

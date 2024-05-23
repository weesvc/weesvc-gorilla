// Package config provides the central configuration definition for application and services.
package config

// Config provides database connection information.
type Config struct {
	Dialect     string
	DatabaseURI string
	Verbose     bool
	// Port to bind the web application api to
	Port int
	// ResourceCachingEnabled denotes caching headers are applied to static web resources
	// Tip: disable during development to immediately reflect resource changes, e.g. CSS and images
	ResourceCachingEnabled bool
}

// NewConfig returns an initialized, but empty, configuration object.
func NewConfig() *Config {
	return &Config{}
}

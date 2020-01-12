package db

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config provides database connection information.
type Config struct {
	DatabaseURI string
	Verbose     bool
}

// InitConfig initializes the database configuration from external settings.
func InitConfig() (*Config, error) {
	config := &Config{
		DatabaseURI: viper.GetString("DatabaseURI"),
		Verbose:     viper.GetBool("Verbose"),
	}
	if config.DatabaseURI == "" {
		return nil, fmt.Errorf("DatabaseURI must be set")
	}
	return config, nil
}

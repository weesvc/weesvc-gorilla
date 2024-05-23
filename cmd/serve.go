package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/weesvc/weesvc-gorilla/internal/config"

	"github.com/weesvc/weesvc-gorilla/internal/server"
)

func newServeCommand(config *config.Config) *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the application server",
		Run: func(_ *cobra.Command, _ []string) {
			server.StartServer(config)
		},
	}

	serveCmd.PersistentFlags().IntVarP(&config.Port, "server-port", "p", 9092, "port to access the api")
	serveCmd.PersistentFlags().StringVar(&config.Dialect, "dialect", "sqlite3", "database dialect")
	serveCmd.PersistentFlags().StringVar(&config.DatabaseURI, "database-uri", "", "database connection string")
	serveCmd.PersistentFlags().BoolVar(&config.ResourceCachingEnabled,
		"resource-caching-enabled",
		true,
		"enable browser caching of web resources")

	_ = viper.BindPFlag("Port", serveCmd.PersistentFlags().Lookup("server-port"))
	_ = viper.BindPFlag("Dialect", serveCmd.PersistentFlags().Lookup("dialect"))
	_ = viper.BindPFlag("DatabaseURI", serveCmd.PersistentFlags().Lookup("database-uri"))
	_ = viper.BindPFlag("ResourceCachingEnabled", serveCmd.PersistentFlags().Lookup("resource-caching-enabled"))

	return serveCmd
}

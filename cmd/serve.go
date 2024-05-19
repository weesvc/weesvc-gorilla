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
		RunE: func(_ *cobra.Command, _ []string) error {
			return server.StartServer(config)
		},
	}

	serveCmd.PersistentFlags().IntVarP(&config.Port, "api-port", "p", 9092, "port to access the api")
	serveCmd.PersistentFlags().StringVar(&config.Dialect, "dialect", "sqlite3", "database dialect")
	serveCmd.PersistentFlags().StringVar(&config.DatabaseURI, "database-uri", "", "database connection string")

	_ = viper.BindPFlag("Port", serveCmd.PersistentFlags().Lookup("api-port"))
	_ = viper.BindPFlag("Dialect", serveCmd.PersistentFlags().Lookup("dialect"))
	_ = viper.BindPFlag("DatabaseURI", serveCmd.PersistentFlags().Lookup("database-uri"))

	return serveCmd
}

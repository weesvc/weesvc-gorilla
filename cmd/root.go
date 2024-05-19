// Package cmd contains available commands for command line access.
package cmd

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/weesvc/weesvc-gorilla/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute parses and runs the command line.
func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	configFile := ""
	cfg := config.NewConfig()

	rootCmd := &cobra.Command{
		Use:   "weesvc",
		Short: "WeeService Application",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return initConfig(configFile, cfg)
		},
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Usage(); err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().BoolVar(&cfg.Verbose, "verbose", false, "verbose output")

	_ = viper.BindPFlag("Verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	subCommands := []func(*config.Config) *cobra.Command{
		newServeCommand,
		newMigrateCommand,
		newVersionCommand,
	}
	for _, sc := range subCommands {
		rootCmd.AddCommand(sc(cfg))
	}

	return rootCmd
}

func initConfig(configFile string, config *config.Config) error {
	viper.SetConfigFile(configFile)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		// If not found, we'll use defaults and suffer consequences otherwise
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return err
		}
	}

	// Enable overrides by environment variable
	// E.g. WEESVC_DIALECT will override `Dialect` within the configuration!
	viper.SetEnvPrefix("WEESVC")
	// Hashes and dots should be treated as underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// Populate our configuration object
	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

package cmd

import (
	"fmt"

	"github.com/weesvc/weesvc-gorilla/internal/config"

	"github.com/spf13/cobra"

	"github.com/weesvc/weesvc-gorilla/internal/env"
)

func newVersionCommand(_ *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("WeeService %v (%v)\n", env.Version, env.Revision)
		},
	}
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/weesvc/weesvc-gorilla/env"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("WeeService %v (%v)\n", env.Version, env.Revision)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

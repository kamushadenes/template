package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "template",
	Short: "A CLI tool for bootstrapping projects",
	Long:  `A CLI tool for bootstrapping projects using template files for different languages.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
}

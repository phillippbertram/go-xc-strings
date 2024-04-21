package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/config"
)

// Global variables to hold flag values for commands
var stringsReferencePath string
var swiftDirectory string
var ignorePatterns []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "xc-strings",
	Short:   "A tool for managing localization strings in Swift projects",
	Version: config.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

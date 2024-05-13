package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xcs",
	Short: "A tool for cleaning localization strings in Swift projects",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	// Configure the --version flag
	rootCmd.Version = fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s\n", config.Version, config.Commit, config.BuildDate)
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`) // Optional: custom format for version output

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

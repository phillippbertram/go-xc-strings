package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/phillippbertram/xc-strings/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xcs",
	Short: "A tool for cleaning localization strings in Swift projects",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .xcs.yaml)")

	// Configure the --version flag
	rootCmd.Version = fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s\n", config.Version, config.Commit, config.BuildDate)
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`) // Optional: custom format for version output

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

}

// will be set via flag (or not)
var cfgFile string

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search for config file in the current directory
		viper.AddConfigPath(".")

		// Search config in home directory with name ".xcs" (without extension).
		// home, err := os.UserHomeDir()
		// if err == nil {
		// 	viper.AddConfigPath(home)
		// }

		viper.SetConfigName(".xcs")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Printf("Error reading config file: %v\n", err)
	}

	// Unmarshal the config into the cfg struct
	if err := viper.Unmarshal(&config.Cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}

package cmd

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"

	"github.com/phillippbertram/xc-strings/internal"
	"github.com/phillippbertram/xc-strings/internal/constants"
	"github.com/phillippbertram/xc-strings/internal/localizable"

	"github.com/spf13/cobra"
)

type UnusedOptions struct {
	removeUnused    bool
	stringsPath     string
	swiftDirectory  string
	baseStringsPath string
	ignorePatterns  []string

	// TODO: dryRun bool
}

var unusedOptions UnusedOptions = UnusedOptions{
	ignorePatterns: constants.DefaultIgnorePatterns,
}

var unusedCmd = &cobra.Command{
	Use:   "unused [strings-path] -b <Localizable.strings> [-d <path to swift code>] [-i <ignore pattern>...]",
	Short: "Finds unused keys in .strings files",
	Long: heredoc.Doc(
		`Check for localization keys defined in a .strings file that are not used in any Swift file within a specified directory.`),
	Example: heredoc.Doc(`
		unused -b Localizable.strings
		unused -b Localizable.strings -d Sources/MyApp -i "Pods/*" "Carthage/*" "*.generated.swift"
	`),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		unusedOptions.stringsPath = args[0]

		if unusedOptions.baseStringsPath == "" {
			return fmt.Errorf("base Localizable.strings file is required")
		}

		if unusedOptions.swiftDirectory == "" {
			unusedOptions.swiftDirectory = "."
		}

		manager, err := localizable.NewStringsFileManager([]string{unusedOptions.stringsPath})
		if err != nil {
			return fmt.Errorf("error initializing strings manager: %w", err)
		}

		// Start a spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Suffix = " Searching for unused keys..."
		s.Start()

		keysForBaseStrings := manager.GetKeysForFile(unusedOptions.baseStringsPath)
		unusedKeys := internal.FindUnusedKeysInSwiftFiles(unusedOptions.swiftDirectory, keysForBaseStrings, unusedOptions.ignorePatterns)
		s.Stop()

		if len(unusedKeys) == 0 {
			color.Green("No unused keys found. 🚀")
			return nil
		}

		for _, key := range unusedKeys {
			fmt.Println(key)
		}
		color.Red("\nFound %d unused keys\n", len(unusedKeys))

		// TODO: remove unused keys from the .strings file
		if unusedOptions.removeUnused {
			color.Yellow("Removing unused is not yet implemented. 🚧")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(unusedCmd)
	unusedCmd.Flags().StringVarP(&unusedOptions.baseStringsPath, "base", "b", "", "Path to the base Localizable.strings file which is used as reference for finding unused keys (required)")
	unusedCmd.Flags().StringVarP(&unusedOptions.swiftDirectory, "swift-dir", "d", "", "Path to the directory containing Swift files (.)")
	unusedCmd.Flags().StringSliceVarP(&unusedOptions.ignorePatterns, "ignore", "i", constants.DefaultIgnorePatterns, "Glob patterns for files or directories to ignore")
	unusedCmd.Flags().BoolVar(&unusedOptions.removeUnused, "remove", false, "Remove unused keys from the .strings file")
}

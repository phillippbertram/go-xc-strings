package cmd

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/briandowns/spinner"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
	"phillipp.io/go-xc-strings/internal/constants"
	"phillipp.io/go-xc-strings/internal/localizable"
)

type UnusedOptions struct {
	removeUnused    bool
	stringsPath     string
	swiftDirectory  string
	baseStringsPath string
	ignorePatterns  []string
}

var unusedOptions UnusedOptions = UnusedOptions{
	ignorePatterns: constants.DefaultIgnorePatterns,
}

var unusedCmd = &cobra.Command{
	Use:   "unused -b <Localizable.strings> [-d <path to swift code>] [-i <ignore pattern>...]",
	Short: "Finds unused keys in .strings files",
	Long: heredoc.Doc(
		`Check for localization keys defined in a .strings file that are not used in any Swift file within a specified directory.`),
	Example: heredoc.Doc(`
		unused -b Localizable.strings
		unused -b Localizable.strings -d Sources/MyApp -i "Pods/*" "Carthage/*" "*.generated.swift"
	`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		manager, err := localizable.NewStringsFileManager(sortOptions.paths)
		if err != nil {
			return fmt.Errorf("error initializing strings manager: %w", err)
		}

		// Start a spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Suffix = " Searching for unused keys..."
		s.Start()

		keysForBaseStrings := manager.GetKeysForFile(unusedOptions.baseStringsPath)
		unusedKeys := internal.SearchKeysInSwiftFiles(unusedOptions.swiftDirectory, keysForBaseStrings, unusedOptions.ignorePatterns)
		for keys := range unusedKeys {
			fmt.Printf("Unused keys in %s:\n", keys)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(unusedCmd)
	unusedCmd.Flags().StringVar(&unusedOptions.baseStringsPath, "base", "b", "Path to the base Localizable.strings file which is used as reference for finding unused keys (required)")
	unusedCmd.Flags().StringVar(&unusedOptions.stringsPath, "strings", "p", "Path to the directory containing the Localizable.string files (.)")
	unusedCmd.Flags().StringVarP(&unusedOptions.swiftDirectory, "swift-dir", "d", "", "Path to the directory containing Swift files (.)")
	unusedCmd.Flags().StringSliceVarP(&unusedOptions.ignorePatterns, "ignore", "i", []string{}, "Glob patterns for files or directories to ignore")
	unusedCmd.Flags().BoolVar(&unusedOptions.removeUnused, "remove", false, "Remove unused keys from the .strings file")
}

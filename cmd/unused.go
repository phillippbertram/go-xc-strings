package cmd

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/briandowns/spinner"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var removeUnused bool

var unusedCmd = &cobra.Command{
	Use:   "unused -r <Localizable.strings> [-d <path to swift code>] [-i <ignore pattern>...]",
	Short: "Finds unused keys in .strings files",
	Long: heredoc.Doc(
		`Check for localization keys defined in a .strings file that are not used in any Swift file within a specified directory.`),
	Example: heredoc.Doc(`
		unused -r Localizable.strings
		unused -r Localizable.strings -d Sources/MyApp -i "Pods/*" "Carthage/*" "*.generated.swift"
	`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if swiftDirectory == "" {
			swiftDirectory = "."
		}
		if stringsReferencePath == "" {
			return fmt.Errorf("please specify the path to the .strings file and the directory containing Swift files")
		}

		// Start a spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Suffix = " Searching for unused keys..."
		s.Start()

		if removeUnused {
			if err := internal.CleanAndSortStringsFiles(stringsPath, stringsReferencePath, swiftDirectory, ignorePatterns, false); err != nil {
				return fmt.Errorf("error cleaning .strings files: %w", err)
			}
		} else {
			unusedKeys, err := internal.FindUnusedKeys(stringsReferencePath, swiftDirectory, ignorePatterns)
			s.Stop()
			if err != nil {
				return fmt.Errorf("error finding unused keys: %v", err)
			}

			if len(unusedKeys) > 0 {
				fmt.Print("The following keys are unused:\n\n")
				for _, key := range unusedKeys {
					fmt.Println(key)
				}
				fmt.Printf("\nFound %d unused keys.\n", len(unusedKeys))
			} else {
				fmt.Println("No unused keys found.")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(unusedCmd)
	unusedCmd.Flags().StringVar(&stringsPath, "strings", "", "Path to the directory containing the Localizable.string files (.)")
	unusedCmd.Flags().StringVarP(&stringsReferencePath, "reference", "r", "", "Path to the Localizable.strings file which is used as reference for finding unused keys (required)")
	unusedCmd.Flags().StringVarP(&swiftDirectory, "swift-dir", "d", "", "Path to the directory containing Swift files (.)")
	unusedCmd.Flags().StringSliceVarP(&ignorePatterns, "ignore", "i", []string{}, "Glob patterns for files or directories to ignore")
	unusedCmd.Flags().BoolVar(&removeUnused, "remove", false, "Remove unused keys from the .strings file")
}

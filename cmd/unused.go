package cmd

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/briandowns/spinner"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var findCmd = &cobra.Command{
	Use:   "unused -s <Localizable.strings> [-d <path to swift code>] [-i <ignore pattern>...]",
	Short: "Finds unused keys in .strings files",
	Long: heredoc.Doc(
		`Check for localization keys defined in a .strings file that are not used in any Swift file within a specified directory.`),
	Example: heredoc.Doc(`
		unused -r Localizable.strings
		unused -r Localizable.strings -d Sources/MyApp -i "Pods/*" "Carthage/*" "*.generated.swift"
	`),
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

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().StringVarP(&stringsReferencePath, "reference", "r", "", "Path to the Localizable.strings file which is used as reference for finding unused keys (required)")
	findCmd.Flags().StringVarP(&swiftDirectory, "swift-dir", "d", "", "Path to the directory containing Swift files (.)")
	findCmd.Flags().StringSliceVarP(&ignorePatterns, "ignore", "i", []string{}, "Glob patterns for files or directories to ignore")
}
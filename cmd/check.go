package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
	"phillipp.io/go-xc-strings/internal/constants"
	"phillipp.io/go-xc-strings/internal/localizable"
)

type CheckOptions struct {
	exitOnIssue     bool
	stringsPath     string
	swiftDirectory  string
	baseStringsPath string
	ignorePatterns  []string
}

var checkOptions CheckOptions = CheckOptions{
	exitOnIssue: true,
}

var checkCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Check for issues in .strings files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) > 0 {
			checkOptions.stringsPath = args[0]
		} else {
			checkOptions.stringsPath = constants.DefaultStringsGlob
		}

		if checkOptions.baseStringsPath == "" {
			return fmt.Errorf("base Localizable.strings file is required")
		}

		if checkOptions.swiftDirectory == "" {
			checkOptions.swiftDirectory = "."
		}

		manager, err := localizable.NewStringsFileManager([]string{checkOptions.stringsPath})
		if err != nil {
			return fmt.Errorf("error initializing strings manager: %w", err)
		}

		// Start a spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		var unsortedFiles []*localizable.StringsFile
		var filesWithDuplicates []*localizable.StringsFile
		var filesWithEmptyValues []*localizable.StringsFile

		for _, file := range manager.Files {

			// check if strings needs sorting
			if !file.IsSorted() || !file.IsSanitized() {
				unsortedFiles = append(unsortedFiles, file)
			}

			// check if strings has duplicates
			if file.HasDuplicates() {
				filesWithDuplicates = append(filesWithDuplicates, file)
			}

			// check if strings has empty values
			if file.HasEmptyValues() {
				filesWithEmptyValues = append(filesWithEmptyValues, file)
			}
		}

		// TODO: check if strings has unused keys
		keysForBaseStrings := manager.GetKeysForFile(checkOptions.baseStringsPath)
		unusedKeys := internal.FindUnusedKeysInSwiftFiles(checkOptions.swiftDirectory, keysForBaseStrings, checkOptions.ignorePatterns)

		if len(unsortedFiles) > 0 {
			color.Yellow("Unsorted files (%d): ", len(unsortedFiles))
			for _, file := range unsortedFiles {
				fmt.Println(file.Path)
			}
		}

		if len(filesWithDuplicates) > 0 {
			color.Yellow("Files with duplicates (%d):", len(filesWithDuplicates))
			for _, file := range filesWithDuplicates {
				fmt.Println(file.Path)
			}
		}

		if len(filesWithEmptyValues) > 0 {
			color.Yellow("Files with empty values (%d):", len(filesWithEmptyValues))
			for _, file := range filesWithEmptyValues {
				fmt.Println(file.Path)
			}
		}

		if len(unusedKeys) > 0 {
			color.Yellow("Unused keys (%d):\n", len(unusedKeys))
			for _, key := range unusedKeys {
				fmt.Println(key)
			}
		}

		anyIssuesOccurred := len(unsortedFiles) > 0 || len(filesWithDuplicates) > 0 || len(filesWithEmptyValues) > 0 || len(unusedKeys) > 0
		if !anyIssuesOccurred {
			color.Green("No issues found. 🚀")
			return nil
		}
		color.Red("Issues found. 🚧")

		if checkOptions.exitOnIssue && anyIssuesOccurred {
			os.Exit(1)
			panic("Should not reach here.")
		}

		fmt.Printf("Unsorted files: %d\n", len(unsortedFiles))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&checkOptions.baseStringsPath, "base", "b", "", "Path to the base Localizable.strings file which is used as reference for finding unused keys (required)")
	checkCmd.Flags().StringVarP(&checkOptions.swiftDirectory, "swift-dir", "d", "", "Path to the directory containing Swift files (.)")
	checkCmd.Flags().StringSliceVarP(&checkOptions.ignorePatterns, "ignore", "i", constants.DefaultIgnorePatterns, "Glob patterns for files or directories to ignore")
}

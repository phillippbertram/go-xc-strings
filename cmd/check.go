package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"

	"github.com/phillippbertram/xc-strings/internal"
	"github.com/phillippbertram/xc-strings/internal/constants"
	"github.com/phillippbertram/xc-strings/internal/localizable"

	"github.com/spf13/cobra"
)

type CheckOptions struct {
	exitOnIssue     bool
	stringsPath     string
	swiftDirectory  string
	baseStringsPath string
	ignorePatterns  []string

	checkSorting     bool
	checkDuplicates  bool
	checkEmptyValues bool
	checkUnused      bool
}

var checkOptions CheckOptions = CheckOptions{
	exitOnIssue: true,
}

var checkCmd = &cobra.Command{
	Use:   "check -b [path to base strings file] -d [path to Swift directory] [path to strings file(s)] ",
	Short: "Check for issues in .strings files",
	Example: heredoc.Doc(`
		# Run all checks (sorting, duplicates, empty values, unused keys):
		$ xcs check

		# Run only the sorting check:
		$ xcs check --check-sorting

		# Run the sorting and duplicate checks only:
		$ xcs check --check-sorting --check-duplicates

		# Specify a base Localizable.strings file and a Swift directory:
		$ xcs check -b ./Base.lproj/Localizable.strings -d ./Sources

		# Specify a strings file to check and ignore certain patterns:
		$ xcs check -b ./Base.lproj/Localizable.strings -d ./Sources --ignore "test/*,docs/*"
	`),
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) > 0 {
			checkOptions.stringsPath = args[0]
		} else {
			checkOptions.stringsPath = constants.DefaultStringsGlob
		}

		if checkOptions.baseStringsPath == "" && checkOptions.checkUnused {
			return fmt.Errorf("base Localizable.strings file is required for unused key check")
		}

		if checkOptions.swiftDirectory == "" {
			checkOptions.swiftDirectory = "."
		}
		// check if any `--check` flag was set
		checkSortingFlag := cmd.Flags().Changed("check-sorting")
		checkDuplicatesFlag := cmd.Flags().Changed("check-duplicates")
		checkEmptyValuesFlag := cmd.Flags().Changed("check-empty-values")
		checkUnusedFlag := cmd.Flags().Changed("check-unused")

		// if no `--check` flag was set, set all check.options to true
		if !checkSortingFlag && !checkDuplicatesFlag && !checkEmptyValuesFlag && !checkUnusedFlag {
			checkOptions.checkSorting = true
			checkOptions.checkDuplicates = true
			checkOptions.checkEmptyValues = true
			checkOptions.checkUnused = true
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

			// Check for sorting if enabled
			if checkOptions.checkSorting && (!file.IsSorted() || !file.IsSanitized()) {
				unsortedFiles = append(unsortedFiles, file)
			}

			// Check for duplicates if enabled
			if checkOptions.checkDuplicates && file.HasDuplicates() {
				filesWithDuplicates = append(filesWithDuplicates, file)
			}

			// Check for empty values if enabled
			if checkOptions.checkEmptyValues && file.HasEmptyValues() {
				filesWithEmptyValues = append(filesWithEmptyValues, file)
			}
		}

		// Check for unused keys if enabled
		var unusedKeys []string
		if checkOptions.checkUnused {
			keysForBaseStrings := manager.GetKeysForFile(checkOptions.baseStringsPath)
			unusedKeys = internal.FindUnusedKeysInSwiftFiles(checkOptions.swiftDirectory, keysForBaseStrings, checkOptions.ignorePatterns)
		}

		// Stop spinner before showing results
		s.Stop()

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
		if anyIssuesOccurred {
			color.Red("Issues found. 🚧")
			if checkOptions.exitOnIssue {
				os.Exit(1)
			}
		} else {
			color.Green("No issues found. 🚀")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&checkOptions.baseStringsPath, "base", "b", "", "Path to the base Localizable.strings file which is used as reference for finding unused keys (required)")
	checkCmd.Flags().StringVarP(&checkOptions.swiftDirectory, "swift-dir", "d", "", "Path to the directory containing Swift files (.)")
	checkCmd.Flags().StringSliceVarP(&checkOptions.ignorePatterns, "ignore", "i", constants.DefaultIgnorePatterns, "Glob patterns for files or directories to ignore")

	// Default all checks to false
	checkCmd.Flags().BoolVar(&checkOptions.checkSorting, "check-sorting", false, "Enable or disable sorting check")
	checkCmd.Flags().BoolVar(&checkOptions.checkDuplicates, "check-duplicates", false, "Enable or disable duplicate check")
	checkCmd.Flags().BoolVar(&checkOptions.checkEmptyValues, "check-empty-values", false, "Enable or disable empty values check")
	checkCmd.Flags().BoolVar(&checkOptions.checkUnused, "check-unused", false, "Enable or disable unused keys check")
}

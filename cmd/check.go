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

// Define constants for the different checks to avoid hardcoded strings
const (
	CheckSorting     = "sorting"
	CheckDuplicates  = "duplicates"
	CheckEmptyValues = "emptyValues"
	CheckUnused      = "unused"
)

// Define a list of all available checks
var allChecks = []string{
	CheckSorting,
	CheckDuplicates,
	CheckEmptyValues,
	CheckUnused,
}

// Define options for different checks and flags
type CheckOptions struct {
	exitOnIssue     bool
	stringsPath     string
	swiftDirectory  string
	baseStringsPath string
	ignorePatterns  []string
	includeChecks   []string
	excludeChecks   []string
}

// Initialize the CheckOptions struct
var checkOptions = CheckOptions{
	exitOnIssue: true,
}

var checkCmd = &cobra.Command{
	Use:   "check -b [path to base strings file] -d [path to Swift directory] [path to strings file(s)]",
	Short: "Check for issues in .strings files",
	Example: heredoc.Doc(`
		# Run all checks (sorting, duplicates, empty values, unused keys):
		$ ./xcs check

		# Include only sorting and duplicates checks:
		$ ./xcs check --include sorting --include duplicates

		# Exclude sorting check:
		$ ./xcs check --exclude sorting

		# Exclude both sorting and duplicate checks:
		$ ./xcs check --exclude sorting --exclude duplicates
	`),
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) > 0 {
			checkOptions.stringsPath = args[0]
		} else {
			checkOptions.stringsPath = constants.DefaultStringsGlob
		}

		// Check if the unused key check requires a base strings file
		if checkOptions.baseStringsPath == "" && contains(checkOptions.includeChecks, CheckUnused) {
			return fmt.Errorf("base Localizable.strings file is required for unused key check")
		}

		// Set a default swift directory if not specified
		if checkOptions.swiftDirectory == "" {
			checkOptions.swiftDirectory = "."
		}

		// Ensure that both `--include` and `--exclude` are not used simultaneously
		if len(checkOptions.includeChecks) > 0 && len(checkOptions.excludeChecks) > 0 {
			return fmt.Errorf("you cannot use both --include and --exclude flags at the same time")
		}

		// Create a map to track the active status of each check
		activeChecks := make(map[string]bool)
		for _, check := range allChecks {
			activeChecks[check] = true // Enable all checks by default
		}

		// If `--include` is used, only enable the specified checks
		if len(checkOptions.includeChecks) > 0 {
			// Disable all checks first
			for _, check := range allChecks {
				activeChecks[check] = false
			}
			// Enable only the included checks
			for _, includeCheck := range checkOptions.includeChecks {
				if _, ok := activeChecks[includeCheck]; ok {
					activeChecks[includeCheck] = true
				} else {
					return fmt.Errorf("unknown check to include: %s", includeCheck)
				}
			}
		}

		// If `--exclude` is used, disable the specified checks
		if len(checkOptions.excludeChecks) > 0 {
			for _, excludeCheck := range checkOptions.excludeChecks {
				if _, ok := activeChecks[excludeCheck]; ok {
					activeChecks[excludeCheck] = false
				} else {
					return fmt.Errorf("unknown check to exclude: %s", excludeCheck)
				}
			}
		}

		// Initialize the strings file manager
		manager, err := localizable.NewStringsFileManager([]string{checkOptions.stringsPath})
		if err != nil {
			return fmt.Errorf("error initializing strings manager: %w", err)
		}

		// Start a spinner to provide feedback while processing
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		var unsortedFiles []*localizable.StringsFile
		var filesWithDuplicates []*localizable.StringsFile
		var filesWithEmptyValues []*localizable.StringsFile

		// Perform the checks based on the active checks map
		for _, file := range manager.Files {

			// Check for sorting if enabled
			if activeChecks[CheckSorting] && (!file.IsSorted() || !file.IsSanitized()) {
				unsortedFiles = append(unsortedFiles, file)
			}

			// Check for duplicates if enabled
			if activeChecks[CheckDuplicates] && file.HasDuplicates() {
				filesWithDuplicates = append(filesWithDuplicates, file)
			}

			// Check for empty values if enabled
			if activeChecks[CheckEmptyValues] && file.HasEmptyValues() {
				filesWithEmptyValues = append(filesWithEmptyValues, file)
			}
		}

		// Check for unused keys if enabled
		var unusedKeys []string
		if activeChecks[CheckUnused] {
			keysForBaseStrings := manager.GetKeysForFile(checkOptions.baseStringsPath)
			unusedKeys = internal.FindUnusedKeysInSwiftFiles(checkOptions.swiftDirectory, keysForBaseStrings, checkOptions.ignorePatterns)
		}

		// Stop the spinner after processing
		s.Stop()

		// Display results for the various checks
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

		// Determine if any issues were found and handle the exit status
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

	// Flags for include and exclude lists
	availableChecks := fmt.Sprintf("%s", allChecks)
	checkCmd.Flags().StringSliceVar(&checkOptions.includeChecks, "include", []string{}, fmt.Sprintf("List of checks to include (%s)", availableChecks))
	checkCmd.Flags().StringSliceVar(&checkOptions.excludeChecks, "exclude", []string{}, fmt.Sprintf("List of checks to exclude (%s)", availableChecks))
}

package cmd

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var skipSort bool
var stringsPath string

// TODO: var dryRun bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans unused localization keys from .strings files",
	Long: `This command searches for unused localization keys within .strings files
in the specified directory and subdirectories, removes them, and optionally sorts the files.`,
	Example: heredoc.Doc(`
		# clean all .strings files in the current directory and its subdirectories
		clean
	`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if swiftDirectory == "" {
			swiftDirectory = "."
		}
		if stringsPath == "" {
			stringsPath = "."
		}

		if stringsReferencePath == "" {
			return fmt.Errorf("please specify the path to the .strings file and the directory containing Swift files")
		}

		// if dryRun {
		// 	fmt.Println("Running in dry-run mode. No changes will be made.")
		// }

		// Start a spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Suffix = " Searching for unused keys..."
		s.Start()

		// Clean and optionally sort the .strings files
		if err := internal.CleanAndSortStringsFiles(stringsPath, stringsReferencePath, swiftDirectory, ignorePatterns, !skipSort); err != nil {
			return fmt.Errorf("error cleaning .strings files: %w", err)
		}

		fmt.Println("Cleaning and sorting completed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringVarP(&stringsReferencePath, "reference", "r", "", "Path to the Localizable.strings file which is used as reference for finding unused keys (required)")
	cleanCmd.Flags().StringVarP(&swiftDirectory, "swift-dir", "d", "", "Path to the directory containing Swift files (.)")
	cleanCmd.Flags().StringVar(&stringsPath, "strings", "", "Path to the directory containing the Localizable.string files (.)")
	cleanCmd.Flags().StringSliceVarP(&ignorePatterns, "ignore", "i", []string{}, "Glob patterns for files or directories to ignore")
	cleanCmd.Flags().BoolVarP(&skipSort, "skip-sort", "s", false, "Skip sort the .strings files")
	// TODO: cleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate the changes without applying them")
}

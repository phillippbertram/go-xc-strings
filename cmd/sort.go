package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var sortCmd = &cobra.Command{
	Use:   "sort [path]",
	Short: "Sorts and groups keys in .strings files",
	Long: `Sorts keys alphabetically in a .strings file and groups them by prefix.
If a directory path is provided, it sorts all .strings files within that directory.
If a file path is provided, it sorts that specific file.`,
	Example: heredoc.Doc(`
		# sort all .strings files in the current directory and its subdirectories
		sort

		# sort all .strings files in a specific directory
		sort path/to/directory

		# sort a specific .strings file
		sort Localizable.strings
	`),
	Args: cobra.MaximumNArgs(1), // Allows zero or one argument only
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine the path to sort
		var path string
		if len(args) == 1 {
			path = args[0]
		} else {
			// If no argument is provided, use the current directory
			var err error
			path, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("error getting current directory: %w", err)
			}
		}

		// Make sure the path is absolute
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("error resolving path: %w", err)
		}

		// Sort the files
		err = internal.SortStringsFiles(absPath)
		if err != nil {
			return fmt.Errorf("error sorting .strings files: %w", err)
		}

		fmt.Println("Sorting completed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)
}

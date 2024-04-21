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
	Long: heredoc.Doc(`
	Sorts keys alphabetically in a .strings file and groups them by prefix.
	If a directory path is provided, it sorts all .strings files within that directory.
	If a file path is provided, it sorts that specific file.
	`),
	Example: heredoc.Doc(`
		# sort all .strings files in the current directory and its subdirectories
		sort

		# sort all .strings files in a specific directory
		sort path/to/directory

		# sort a specific .strings file
		sort path1/Localizable.strings path2/InfoPlist.strings
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine the path to sort
		var paths []string
		if len(args) == 0 {
			// If no argument is provided, use the current directory
			var err error
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error getting current directory: %w", err)
			}
			paths = []string{wd}

		} else {
			paths = append(paths, args...)
		}

		for _, pattern := range paths {
			// Expand the glob pattern to match files and directories
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return fmt.Errorf("error processing glob pattern '%s': %w", pattern, err)
			}
			for _, match := range matches {
				err := internal.SortStringsFiles(match)
				if err != nil {
					return fmt.Errorf("error sorting files in '%s': %w", match, err)
				}
			}
		}
		fmt.Println("Sorting completed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)
}

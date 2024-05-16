package cmd

import (
	"fmt"

	"github.com/phillippbertram/xc-strings/internal/constants"
	"github.com/phillippbertram/xc-strings/internal/localizable"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type SortOptions struct {
	paths        []string
	dryRun       bool
	skipSanitize bool
}

var sortOptions SortOptions = SortOptions{
	paths: []string{constants.DefaultStringsGlob},
}

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

		if len(args) != 0 {
			sortOptions.paths = args
		}

		manager, err := localizable.NewStringsFileManager(sortOptions.paths)
		if err != nil {
			return fmt.Errorf("error initializing strings manager: %w", err)
		}

		if sortOptions.dryRun {
			color.Yellow("Running in dry-run mode. No changes will be made.\n")
		}

		manager.Sort()

		if !sortOptions.skipSanitize {
			manager.Sanitize()
		} else {
			fmt.Printf("Skipping sanitizing the file\n")
		}

		if !sortOptions.dryRun {
			manager.Save()
		} else {
			color.Yellow("Dry-run completed. No changes were made.\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)

	sortCmd.Flags().BoolVar(&sortOptions.dryRun, "dry-run", false, "Prints the changes without writing them to the file")
	sortCmd.Flags().BoolVar(&sortOptions.skipSanitize, "skip-sanitize", false, "Skips sanitizing the file after sorting")
}

package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal/constants"
	"phillipp.io/go-xc-strings/internal/localizable"
)

// Script to find duplicate keys in .strings files
// for f in <PATH>/*.lproj/Localizable.strings; do echo "$f:"; sed '/^$/d' "$f" | sort | uniq -cd; echo; done

type DuplicatesOptions struct {
	paths            []string
	removeDuplicates bool
	dryRun           bool
}

var duplicatesOptions DuplicatesOptions = DuplicatesOptions{
	paths: []string{constants.DefaultStringsGlob},
}

var duplicatesCmd = &cobra.Command{
	Use:   "duplicates [path]",
	Short: "Find duplicate keys in .strings files",
	Long:  `Finds and lists duplicate keys in .strings files located at the specified path. The command can be used for a single file or a directory recursively.`,
	Example: heredoc.Doc(`
		# find duplicate keys in all .strings files in the current directory and its subdirectories
		duplicates

		# find duplicate keys in all .strings files in the specified directory
		duplicates path/to/directory

		# remove all but the last occurrence of each duplicate key
		duplicates --remove
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			sortOptions.paths = args
		}

		manager, err := localizable.NewStringsFileManager(sortOptions.paths)
		if err != nil {
			return fmt.Errorf("error initializing strings manager: %w", err)
		}

		duplicates := manager.FindDuplicates()
		if len(duplicates) == 0 {
			color.Green("No duplicate keys found.")
			return nil
		}

		for file, dup := range duplicates {
			fmt.Printf("Duplicates in %s:\n", file)
			for key, l := range dup.Duplicates {
				fmt.Printf("%d - %s\n", len(l), key)

				// TODO: print the lines if wanted?
			}
			fmt.Println()
		}

		if duplicatesOptions.removeDuplicates {
			// remove all but the last occurrence of each duplicate key
			for _, file := range manager.Files {
				removedLines := file.RemoveDuplicatesKeepLast()
				fmt.Printf("Removed %d duplicates in %s\n", len(removedLines), file.Path)

				if !duplicatesOptions.dryRun {
					_ = file.Save()
				}
			}
		}

		if duplicatesOptions.dryRun {
			color.Yellow("Dry-run completed. No changes were made.")
		} else {
			color.Green("All duplicates removed successfully.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(duplicatesCmd)
	duplicatesCmd.Flags().BoolVar(&duplicatesOptions.removeDuplicates, "remove", false, "Remove all but the last occurrence of each duplicate key")
	duplicatesCmd.Flags().BoolVar(&duplicatesOptions.dryRun, "dry-run", false, "Prints the changes without writing them to the file")
}

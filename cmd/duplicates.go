package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var removeDuplicates bool

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
		// Determine the path to sort
		var path string
		if len(args) == 0 {
			// If no argument is provided, use the current directory
			var err error
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error getting current directory: %w", err)
			}
			path = wd

		} else {
			path = args[0]
		}

		duplicates, err := internal.FindDuplicates(path)
		if err != nil {
			return err
		}
		if len(duplicates) == 0 {
			color.Green("No duplicate keys found.")
			return nil
		}

		fileColor := color.New(color.FgCyan, color.Bold)
		keyColor := color.New(color.FgYellow)
		valueColor := color.New(color.FgGreen)

		if removeDuplicates {
			err = internal.RemoveDuplicatesKeepLast(path, duplicates)
			if err != nil {
				return fmt.Errorf("failed to remove duplicates: %w", err)
			}
			color.Green("Duplicates removed successfully.")
		} else {

			for file, keys := range duplicates {
				fileColor.Printf("Duplicates in %s:\n", file)
				for key, values := range keys {
					keyColor.Printf("%s:\n", key)
					for _, value := range values {
						valueColor.Printf("-> %s\n ", value)
					}
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(duplicatesCmd)
	duplicatesCmd.Flags().BoolVar(&removeDuplicates, "remove", false, "Remove all but the last occurrence of each duplicate key")
}

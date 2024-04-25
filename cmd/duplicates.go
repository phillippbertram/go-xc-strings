package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var removeDuplicates bool

var duplicatesCmd = &cobra.Command{
	Use:   "duplicates [path]",
	Short: "Find duplicate keys in .strings files",
	Long:  `Finds and lists duplicate keys in .strings files located at the specified path. The command can be used for a single file or a directory recursively.`,
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
			fmt.Println("No duplicate keys found.")
			return nil
		}

		if removeDuplicates {
			err = internal.RemoveDuplicates(path, duplicates)
			if err != nil {
				return fmt.Errorf("failed to remove duplicates: %w", err)
			}
			fmt.Println("Duplicates removed successfully.")
		} else {
			if len(duplicates) == 0 {
				fmt.Println("No duplicate keys found.")
				return nil
			}
			for file, keys := range duplicates {
				fmt.Printf("Duplicates in %s:\n", file)
				for key, value := range keys {
					fmt.Printf("%s = %s\n", key, value)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(duplicatesCmd)
	duplicatesCmd.Flags().BoolVar(&removeDuplicates, "remove", false, "Remove all but the first occurrence of each duplicate key")
}

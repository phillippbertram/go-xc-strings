package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var excludeLanguages []string

var removeCmd = &cobra.Command{
	Use:   "remove [key] [directory]",
	Short: "Removes a localization key from all .strings files",
	Example: heredoc.Doc(`
	    # removes the key "key_name" from all .strings files 
		remove "key_name"

		# removes the key "key_name" from all .strings files in the specified directory
		remove "key_name" -d path/to/directory

		# removes the key "key_name" from the specified .strings file
		remove "key_name" path/to/Localizable.strings
	`),
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		directory := args[1]

		fileNames, err := internal.RemoveKeyFromAllStringsFiles(key, directory)
		if err != nil {
			return fmt.Errorf("failed to remove key: %w", err)
		}

		if len(fileNames) == 0 {
			fmt.Println("Key not found in any .strings files")
			return nil
		}

		fmt.Printf("Key removed successfully from the following %d files:\n", len(fileNames))
		for _, fileName := range fileNames {
			fmt.Println(fileName)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringArrayVarP(&excludeLanguages, "exclude", "e", []string{}, "Exclude languages from the operation")
}

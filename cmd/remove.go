package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var excludeLanguages []string

var removeCmd = &cobra.Command{
	Use:   "remove [key] -p [directory]",
	Short: "Removes a localization key from all .strings files",
	Example: heredoc.Doc(`
	    # removes the key "key_name" from all .strings files recursively in the current directory
		remove "key_name"

		# removes the key "key_name" from all .strings files in the specified directory
		remove "key_name" -p path/to/directory

		# removes multiple the keys from the specified .strings file
		remove "key_name" "key_name2" "key_name_3" -p path/to/Localizable.strings
	`),
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keys := args

		for _, key := range keys {
			fileNames, err := internal.RemoveKeyFromAllStringsFiles(key, stringsPath, excludeLanguages)
			if err != nil {
				return fmt.Errorf("failed to remove key: %w", err)
			}

			if len(fileNames) == 0 {
				fmt.Println("Key not found in any .strings files")
				return nil
			}

			fmt.Printf("Key [%s] removed successfully from the following %d files:\n", key, len(fileNames))
			for _, fileName := range fileNames {
				fmt.Println(fileName)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringArrayVarP(&excludeLanguages, "exclude", "e", []string{}, "Exclude languages from the operation")
	removeCmd.Flags().StringVarP(&stringsPath, "strings", "p", ".", "Path to the directory containing the .strings files")
}

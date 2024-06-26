package cmd

import (
	"fmt"

	"github.com/phillippbertram/xc-strings/internal/constants"
	"github.com/phillippbertram/xc-strings/internal/localizable"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type KeysOptions struct {
	path       string
	keys       []string
	removeKeys bool
	dryRun     bool

	// TODO: excludeLanguages []string
}

var keysOptions KeysOptions = KeysOptions{
	path: constants.DefaultStringsGlob,
}

var findKeysCmd = &cobra.Command{
	Use:   "keys [keys] -p [directory]",
	Short: "Finds keys from all .strings files",
	Example: heredoc.Doc(`
	    # finds the key "key_name" from all .strings files recursively in the current directory
		find "key_name"

		# removes the key "key_name" from all .strings files in the specified directory
		remove "key_name" path/to/directory --remove

		# removes multiple the keys from the specified .strings file
		remove "key_name" "key_name2" "key_name_3" path/to/Localizable.strings --remove
	`),
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) > 1 {
			keysOptions.keys = args[:len(args)-1]
		}

		keysOptions.path = args[len(args)-1]
		manager, err := localizable.NewStringsFileManager([]string{keysOptions.path})
		if err != nil {
			return fmt.Errorf("error initializing strings manager: %w", err)
		}

		var keys []string
		if len(keysOptions.keys) == 0 {
			keys = manager.GetAllKeys()
			for _, key := range keys {
				fmt.Printf("%s\n", key)
			}
			color.Green("Found %d unique keys in %d files\n", len(keys), len(manager.Files))
			return nil
		}

		for _, file := range manager.Files {

			for _, key := range keysOptions.keys {
				foundLines := file.GetLinesForKey(key)

				if len(foundLines) == 0 {
					fmt.Printf("Key [%s] not found] in %s\n", key, file.Path)
					continue
				}

				if keysOptions.removeKeys {
					removed := file.RemoveKey(key)
					fmt.Printf("Key [%s] removed [%dx] in %s\n", key, len(removed), file.Path)

					if !keysOptions.dryRun {
						if err := file.Save(); err != nil {
							return fmt.Errorf("error saving file: %w", err)
						}
					}
				} else {
					fmt.Printf("Key [%s] found [%dx] in %s\n", key, len(foundLines), file.Path)
				}

			}

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findKeysCmd)

	// TODO: exclude certain language files
	// removeCmd.Flags().StringArrayVarP(&excludeLanguages, "exclude", "e", []string{}, "Exclude languages from the operation")
	findKeysCmd.Flags().BoolVar(&keysOptions.removeKeys, "remove", false, "Remove the key from the .strings file")
	findKeysCmd.Flags().BoolVar(&keysOptions.dryRun, "dry-run", false, "Run the command without making any changes")
}

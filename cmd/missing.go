package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/phillippbertram/xc-strings/internal/localizable"
	"github.com/spf13/cobra"
)

type MissingCmdOptions struct {
	baseStringsPath string
	stringsPath     string
}

var missingOptions MissingCmdOptions = MissingCmdOptions{}

var missingCmd = &cobra.Command{
	Use:   "missing [strings-path] -b <base Localizable.strings>",
	Short: "Find missing translations in the strings files",
	Example: heredoc.Doc(`
		# find missing translations in all .strings files in the current directory and its subdirectories
		xcs missing App/Resources -b App/Resources/en.lproj/Localizable.strings
	`),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		missingOptions.stringsPath = args[0]
		return findMissingKeys(missingOptions)
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)
	missingCmd.Flags().StringVarP(&missingOptions.baseStringsPath, "base", "b", "", "Path to the base Localizable.strings file which is used as reference for finding unused keys (required)")
}

func findMissingKeys(opts MissingCmdOptions) error {
	manager, err := localizable.NewStringsFileManager([]string{missingOptions.stringsPath})
	if err != nil {
		return fmt.Errorf("error initializing strings manager: %w", err)
	}

	baseFile := manager.GetFile(opts.baseStringsPath)
	baseKeys := manager.GetKeysForFile(opts.baseStringsPath)
	var missingTranslations map[string][]localizable.Line = make(map[string][]localizable.Line)

	for _, file := range manager.Files {
		// Skip the base file
		if file.Path == opts.baseStringsPath {
			continue
		}
		// fmt.Printf("Checking %s:", file.Path)

		keys := manager.GetKeysForFile(file.Path)
		// fmt.Printf(" %d Keys\n", len(keys))

		for _, key := range baseKeys {
			if !contains(keys, key) {
				if missingTranslations[file.Path] == nil {
					missingTranslations[file.Path] = make([]localizable.Line, 0)
				}
				line := baseFile.GetLinesForKey(key)[0]
				missingTranslations[file.Path] = append(missingTranslations[file.Path], line)
			}
		}
	}

	if len(missingTranslations) == 0 {
		color.Green("No missing translations found.\n")
		return nil
	}

	for file, lines := range missingTranslations {
		color.Yellow("%d Missing translations in %s:\n", len(lines), file)
		for _, key := range lines {
			fmt.Println(key.Text)
		}
		fmt.Println()
	}

	return nil
}

func contains(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

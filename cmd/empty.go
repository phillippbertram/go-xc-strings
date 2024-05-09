package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal/constants"
	"phillipp.io/go-xc-strings/internal/localizable"
)

// command to find empty translation values

type EmptyOptions struct {
	path string
}

var emptyOptions EmptyOptions = EmptyOptions{
	path: constants.DefaultStringsGlob,
}

var emptyCmd = &cobra.Command{
	Use:   "empty [path]",
	Short: "Find empty translation values in .strings files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		emptyOptions.path = args[0]

		manager, err := localizable.NewStringsFileManager([]string{emptyOptions.path})
		if err != nil {
			return fmt.Errorf("error initializing strings manager: %w", err)
		}

		for _, file := range manager.Files {
			lines := file.Lines
			for _, line := range lines {
				if line.Key != "" {
					if line.Value == "" {
						fmt.Printf("Empty translation value in %s: %s\n", file.Path, line.Key)
					}
				}

			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(emptyCmd)
}

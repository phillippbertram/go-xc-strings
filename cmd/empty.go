package cmd

import (
	"fmt"

	"github.com/phillippbertram/xc-strings/internal/constants"
	"github.com/phillippbertram/xc-strings/internal/localizable"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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

		for idx, file := range manager.Files {
			fmt.Printf("Checking %s\n", file.Path)
			emptyLines := file.EmptyValues()

			for _, line := range emptyLines {
				color.Yellow("Empty translation for: %s\n", line.Key)
			}

			if idx < len(manager.Files)-1 {
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(emptyCmd)
}

package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

// TODO: implement

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans unused localization keys from .strings files",
	Long: `Searches for unused localization keys within .strings files
in the specified directory and subdirectories, removes them, and also sorts the files.`,
	Example: heredoc.Doc(`
		# clean all .strings files in the current directory and its subdirectories
		clean
	`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not implemented 😱")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

}

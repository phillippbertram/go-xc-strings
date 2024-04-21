package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sorts and groups keys in .strings files",
	Long: `Sorts keys alphabetically in a .strings file and groups them by prefix. 
	Keys with the same prefix are grouped together, and an empty line is added between different groups.`,
	Run: func(cmd *cobra.Command, args []string) {
		if stringsPath == "" {
			fmt.Println("Please specify the path to the .strings file.")
			os.Exit(1)
		}
		err := internal.SortStringsFile(stringsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sorting .strings file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("The .strings file has been sorted and grouped successfully.")
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.PersistentFlags().StringVarP(&stringsPath, "strings", "s", "", "Path to the Localizable.strings file (required)")
}

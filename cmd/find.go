package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Finds unused keys in .strings files",
	Long: `This command checks for localization keys defined in a .strings file that are 
not used in any Swift file within a specified directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if stringsPath == "" || swiftDirectory == "" {
			fmt.Println("Please specify both the strings file path and the Swift files directory.")
			os.Exit(1)
		}
		unusedKeys, err := internal.FindUnusedKeys(stringsPath, swiftDirectory, ignorePatterns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding unused keys: %v\n", err)
			os.Exit(1)
		}
		if len(unusedKeys) > 0 {
			fmt.Println("The following keys are unused:")
			for _, key := range unusedKeys {
				fmt.Println(key)
			}
		} else {
			fmt.Println("No unused keys found.")
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.PersistentFlags().StringVarP(&stringsPath, "strings", "s", "", "Path to the Localizable.strings file (required)")
	findCmd.PersistentFlags().StringVarP(&swiftDirectory, "swift-dir", "d", "", "Path to the directory containing Swift files (required)")
	findCmd.PersistentFlags().StringSliceVarP(&ignorePatterns, "ignore", "i", []string{}, "Glob patterns for files or directories to ignore")
}

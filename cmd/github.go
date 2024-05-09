package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"phillipp.io/go-xc-strings/internal/constants"
)

type GithubOptions struct {
	showRelease bool
}

var githubOptions GithubOptions

var githubCmd = &cobra.Command{
	Use:   "gh",
	Short: "Show the current GitHub repository",
	Example: heredoc.Doc(`
		# show the current GitHub repository
		gh

		# show releases
		gh --releases
	`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if githubOptions.showRelease {
			return browser.OpenURL(constants.GithubReleasesPage)
		} else {
			return browser.OpenURL(constants.GithubPage)
		}
	},
}

func init() {
	rootCmd.AddCommand(githubCmd)
	githubCmd.Flags().BoolVar(&githubOptions.showRelease, "releases", false, "Show releases")

}

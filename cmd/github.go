package cmd

import (
	"github.com/phillippbertram/xc-strings/internal/constants"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
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
		// print config
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

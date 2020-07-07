package cmd

import (
	"github.com/Masterminds/semver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// New create a new root command for 'kitt'.
func New(version semver.Version) *cobra.Command {
	root := &cobra.Command{
		Use:          "kitt",
		Short:        "KUDO Index transfer tool",
		Long:         "Synchronize KUDO operators from remote sources",
		SilenceUsage: true,
	}

	root.AddCommand(updateCmd())
	root.AddCommand(versionCmd(version))

	verbose := root.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if verbose != nil && *verbose {
			log.SetLevel(log.DebugLevel)
		}
	}

	return root
}

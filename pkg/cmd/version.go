package cmd

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

func versionCmd(version semver.Version) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Args:  cobra.NoArgs,
		Short: "Print the version",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Printf("kitt %s\n", version.String())
	}

	return cmd
}

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/kudobuilder/kitt/pkg/update"
)

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [operator.yaml...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Update a repository with operators",
		Long: `KUDO repositories consist of a collection of indexed operator package tarballs.
kitt creates or updates such a repository by resolving a list of operator
references and creating an operator package tarball for each reference.`,
	}

	force := cmd.Flags().BoolP("force", "f", false, "force update of operators that are already indexed")

	repoPath := cmd.Flags().String("repository", ".", "path to the operator repository")

	if err := cmd.MarkFlagDirname("repository"); err != nil {
		panic(err)
	}

	repoURL := cmd.Flags().String("repository_url", "", "URL of the operator repository to set in \"index.yaml\"")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return update.Update(cmd.Context(), update.FromFiles(args), *repoPath, *repoURL, *force)
	}

	return cmd
}

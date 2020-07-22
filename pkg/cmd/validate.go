package cmd

import (
	"github.com/spf13/cobra"

	"github.com/kudobuilder/kitt/pkg/loader"
	"github.com/kudobuilder/kitt/pkg/validate"
)

func validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [operator.yaml...]",
		Short: "Validate operator references",
		Long: `Run various validation checks that ensure the consistency and validity of the
operator references as well as their referenced operator packages.`,
	}

	strict := cmd.Flags().Bool("strict", false, "treat warnings as errors")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return validate.Validate(cmd.Context(), loader.FromFiles(args), *strict)
	}

	return cmd
}

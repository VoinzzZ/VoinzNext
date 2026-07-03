package cli

import (
	"github.com/spf13/cobra"
)

// ponytail: implement actual feature addition logic; see the validFeatures list in the original stub

var addCmd = &cobra.Command{
	Use:    "add [feature]",
	Short:  "Add a feature to an existing VoinzNext project",
	Long:   "Add a feature to an existing VoinzNext project.",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

package cli

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VoinzNext to the latest version",
	Long:  "Check for updates and install the latest version of VoinzNext.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("  Checking for updates...")

		updateCmd := exec.Command("go", "install", "github.com/VoinzzZ/VoinzNext/cmd/voinznest@latest")
		output, err := updateCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("update failed: %w\n%s", err, string(output))
		}

		fmt.Println("  ✅ VoinzNext has been updated to the latest version!")
		return nil
	},
}

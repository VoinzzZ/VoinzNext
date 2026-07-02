package cli

import (
	"fmt"
	"os/exec"

	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VoinzNext to the latest version",
	Long:  "Check for updates and install the latest version of VoinzNext.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("  %s Checking for updates...\n", style.SprintCyan("●"))

		updateCmd := exec.Command("go", "install", "github.com/VoinzzZ/VoinzNext/cmd/voinznest@latest")
		output, err := updateCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("  %s Update failed: %v\n", style.SprintRed("✘"), err)
			fmt.Printf("  %s\n\n", style.Dimmed(string(output)))
			return err
		}

		fmt.Printf("  %s VoinzNext has been updated to the latest version!\n\n", style.SprintGreen("✔"))
		return nil
	},
}

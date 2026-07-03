package cli

import (
	"fmt"

	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/spf13/cobra"
)

var (
	Version   = "0.5.0"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print VoinzNext version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		fmt.Printf("  %s  %s %s\n", style.SprintCyan("◆"), style.Value("VoinzNext"), Version)
		fmt.Printf("  %s  %s %s\n", style.Dimmed("●"), style.Dimmed("Commit:"), style.Dimmed(GitCommit))
		fmt.Printf("  %s  %s %s\n", style.Dimmed("●"), style.Dimmed("Built:"), style.Dimmed(BuildTime))
		fmt.Println()
	},
}

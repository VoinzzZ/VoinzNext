package cli

import (
	"fmt"
	"strings"

	"github.com/VoinzzZ/VoinzNext/internal/config"
	"github.com/VoinzzZ/VoinzNext/internal/registry"
	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tech stack options",
	Long:  "Display all available technologies and frameworks that can be selected during the interactive survey.",
	RunE: func(cmd *cobra.Command, args []string) error {
		style.Banner("VoinzNext - Available Tech Stacks", "Choose your preferred technologies")

		for _, q := range registry.Questions {
			if len(q.Options) == 0 {
				continue
			}

			fmt.Printf("  %s %s\n", style.Label("◆"), style.Value(q.Message))

			for _, o := range q.Options {
				mark := " "
				highlight := style.Dimmed
				if o.ID == q.Default {
					mark = style.SprintGreen("★")
					highlight = style.Value
				}
				name := fmt.Sprintf("%-22s", o.Name)
				fmt.Printf("    %s %s %s\n", mark, highlight(name), style.Dimmed(o.Description))
			}
			fmt.Println()
		}

		fmt.Printf("  %s  %s\n\n", style.SprintGreen("★ Legend"), style.Dimmed("default option"))
		return nil
	},
}

func getDefaultDisplay(q config.Question) string {
	for _, o := range q.Options {
		if o.ID == q.Default {
			return o.Name
		}
	}
	return q.Default
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
}

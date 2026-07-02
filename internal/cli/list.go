package cli

import (
	"fmt"

	"github.com/VoinzzZ/VoinzNext/internal/config"
	"github.com/VoinzzZ/VoinzNext/internal/registry"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tech stack options",
	Long:  "Display all available technologies and frameworks that can be selected during the interactive survey.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println()
		fmt.Println("  ╔══════════════════════════════════════════╗")
		fmt.Println("  ║    VoinzNext - Available Tech Stacks    ║")
		fmt.Println("  ╚══════════════════════════════════════════╝")
		fmt.Println()

		for _, q := range registry.Questions {
			if len(q.Options) == 0 {
				continue
			}
			fmt.Printf("  📦 %s\n", q.Message)
			fmt.Printf("     Default: %s\n", getDefaultDisplay(q))
			for _, o := range q.Options {
				mark := " "
				if o.ID == q.Default {
					mark = "★"
				}
				fmt.Printf("     %s %-20s %s\n", mark, o.Name, o.Description)
			}
			fmt.Println()
		}

		fmt.Println("  Legend: ★ = default option")
		fmt.Println()
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

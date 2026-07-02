package cli

import (
	"fmt"

	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/spf13/cobra"
)

var validFeatures = []string{
	"prisma", "drizzle",
	"nextauth", "lucia", "clerk",
	"trpc", "shadcn",
	"vitest", "jest", "playwright",
}

var addCmd = &cobra.Command{
	Use:   "add [feature]",
	Short: "Add a feature to an existing VoinzNext project",
	Long: `Add additional features to an already-generated project.
Supports: prisma, drizzle, nextauth, lucia, clerk, trpc, shadcn, vitest, jest, playwright

Examples:
  voinznext add prisma     - Add Prisma ORM to existing project
  voinznext add nextauth   - Add NextAuth.js authentication
  voinznext add shadcn     - Add shadcn/ui components
  voinznext add trpc       - Add tRPC API pattern`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Printf("  %s Please specify a feature to add\n\n", style.SprintYellow("●"))
			fmt.Printf("  %s\n", style.Label("Available features:"))
			for _, f := range validFeatures {
				fmt.Printf("    %s %s\n", style.SprintCyan("·"), style.Value(f))
			}
			fmt.Println()
			return nil
		}

		feature := args[0]
		validMap := make(map[string]bool)
		for _, f := range validFeatures {
			validMap[f] = true
		}

		if !validMap[feature] {
			fmt.Printf("  %s Unknown feature %s\n\n", style.SprintRed("✘"), style.Value(feature))
			fmt.Printf("  %s\n", style.Label("Available features:"))
			for _, f := range validFeatures {
				fmt.Printf("    %s %s\n", style.SprintCyan("·"), style.Value(f))
			}
			fmt.Println()
			return nil
		}

		fmt.Printf("  %s Adding %s...\n", style.SprintCyan("●"), style.Value(feature))
		fmt.Printf("  %s Feature %s is coming soon!\n", style.SprintYellow("●"), style.Value(feature))
		fmt.Printf("  %s For now, run %s to generate a new project with this feature.\n",
			style.SprintCyan("●"),
			style.Value("voinznext init"))
		fmt.Println()

		return nil
	},
}

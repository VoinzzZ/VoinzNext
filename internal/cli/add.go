package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [feature]",
	Short: "Add a feature to an existing VoinzNext project",
	Long: `Add additional features to an already-generated project.
Supports: prisma, drizzle, nextauth, lucia, clerk, trpc, shadcn, vitest, jest, playwright

Examples:
  voinznest add prisma     - Add Prisma ORM to existing project
  voinznest add nextauth   - Add NextAuth.js authentication
  voinznest add shadcn     - Add shadcn/ui components
  voinznest add trpc       - Add tRPC API pattern`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please specify a feature to add\n\nAvailable features: prisma, drizzle, nextauth, lucia, clerk, trpc, shadcn, vitest, jest, playwright")
		}

		feature := args[0]
		fmt.Printf("  Adding %q to project...\n", feature)

		validFeatures := map[string]bool{
			"prisma": true, "drizzle": true,
			"nextauth": true, "lucia": true, "clerk": true,
			"trpc": true, "shadcn": true,
			"vitest": true, "jest": true, "playwright": true,
		}

		if !validFeatures[feature] {
			return fmt.Errorf("unknown feature: %q\n\nAvailable features: prisma, drizzle, nextauth, lucia, clerk, trpc, shadcn, vitest, jest, playwright", feature)
		}

		fmt.Printf("  ✅ Feature %q is not yet implemented (coming soon!)\n", feature)
		fmt.Println("  For now, re-run `voinznest init` to generate a new project with this feature.")
		return nil
	},
}

package cli

import (
	"fmt"

	"github.com/VoinzzZ/VoinzNext/internal/generator"
	"github.com/VoinzzZ/VoinzNext/internal/survey"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Starts interactive survey and generates a new Next.js project",
	Long: `Runs an interactive survey asking about your preferred tech stack,
then generates a complete Next.js project with all configurations,
components, dependencies, and .env.example files.

Usage:
  voinznest init

The survey will ask about:
  - Project name, router type, language
  - CSS framework, UI library
  - Database ORM, authentication
  - API pattern, testing framework
  - Docker, ESLint/Prettier, Git setup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := survey.RunSurvey()
		if err != nil {
			return fmt.Errorf("survey failed: %w", err)
		}

		fmt.Printf("\n  Generating project: %s\n\n", cfg.ProjectName)

		gen := generator.New(cfg)
		if err := gen.Generate(); err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}

		if err := gen.PostGenerate(); err != nil {
			return fmt.Errorf("post-generation failed: %w", err)
		}

		fmt.Printf("\n  ✅ Project %q created successfully!\n\n", cfg.ProjectName)
		fmt.Printf("  Next steps:\n")
		fmt.Printf("    $ cd %s\n", cfg.ProjectName)
		fmt.Printf("    $ %s install\n", cfg.PackageManager)
		fmt.Printf("    $ %s run dev\n\n", cfg.PackageManager)

		return nil
	},
}

func init() {
	initCmd.Flags().BoolP("yes", "y", false, "Skip prompts and use defaults")
}

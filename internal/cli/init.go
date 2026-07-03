package cli

import (
	"errors"
	"fmt"

	"github.com/VoinzzZ/VoinzNext/internal/config"
	"github.com/VoinzzZ/VoinzNext/internal/generator"
	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/VoinzzZ/VoinzNext/internal/survey"
	surveylib "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Start interactive survey and generate a new Next.js project",
	Long: `Run an interactive survey asking about your preferred tech stack,
then generate a complete Next.js project with all configurations,
components, dependencies, and .env.example files.

Usage:
  voinznext init
  voinznext init my-app --yes

The survey will ask about:
  - Project name, router type, language
  - CSS framework, UI library
  - Database ORM, authentication
  - API pattern, testing framework
  - Docker, ESLint/Prettier, Git setup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skipPrompts, _ := cmd.Flags().GetBool("yes")

		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}

		cfg, err := survey.RunSurvey(skipPrompts, projectName)
		if err != nil {
			style.ErrorBanner(fmt.Errorf("survey failed: %w", err))
			return err
		}

		if skipPrompts {
			cfg.Overwrite = true
		}

		gen := generator.New(cfg)
		if err := gen.Generate(); err != nil {
			if errors.Is(err, config.ErrDirNotEmpty) {
				var overwrite bool
				prompt := &surveylib.Confirm{
					Message: fmt.Sprintf("Directory %q already exists and is not empty. Overwrite?", cfg.ProjectName),
					Default: false,
				}
				if askErr := surveylib.AskOne(prompt, &overwrite); askErr != nil {
					return askErr
				}
				if !overwrite {
					fmt.Printf("  %s Aborted.\n", style.SprintYellow("●"))
					return nil
				}
				cfg.Overwrite = true
				gen = generator.New(cfg)
				if err := gen.Generate(); err != nil {
					style.ErrorBanner(fmt.Errorf("generation failed: %w", err))
					return err
				}
			} else {
				style.ErrorBanner(fmt.Errorf("generation failed: %w", err))
				return err
			}
		}

		if err := gen.PostGenerate(); err != nil {
			style.ErrorBanner(fmt.Errorf("post-generation failed: %w", err))
			return err
		}

		style.SuccessBanner(cfg.ProjectName)

		style.NextSteps([][2]string{
			{"➜", fmt.Sprintf("cd %s", style.Value(cfg.ProjectName))},
			{"➜", fmt.Sprintf("%s install", style.Value(cfg.PackageManager))},
			{"➜", fmt.Sprintf("%s run dev", style.Value(cfg.PackageManager))},
		})

		return nil
	},
}

func init() {
	initCmd.Flags().BoolP("yes", "y", false, "Skip prompts and use defaults")
}

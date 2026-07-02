package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "voinznest",
	Short: style.Dimmed("VoinzNext - Interactive Next.js Starter Generator"),
	Long: style.Dimmed(`VoinzNext is a CLI tool that helps you scaffold a Next.js project
with your preferred tech stack through an interactive survey.

It generates a complete project structure with all configurations,
components, and dependencies pre-configured.`),
	Run: func(cmd *cobra.Command, args []string) {
		style.Banner("VoinzNext - CLI", "Interactive Next.js Starter Generator")

		fmt.Printf("  %s %s\n", style.Label("◆"), style.Value("VoinzNext helps you scaffold Next.js projects"))
		fmt.Printf("  %s %s\n\n", style.Label("◆"), style.Value("Answer a few questions and get a production-ready project."))

		fmt.Printf("  %s\n", style.Label("Available commands:"))
		fmt.Printf("    %s  %-12s %s\n", style.SprintCyan("●"), "init", style.Dimmed("Start interactive survey and generate project"))
		fmt.Printf("    %s  %-12s %s\n", style.SprintCyan("●"), "list", style.Dimmed("Show available tech stack options"))
		fmt.Printf("    %s  %-12s %s\n", style.SprintCyan("●"), "add", style.Dimmed("Add a feature to existing project"))
		fmt.Printf("    %s  %-12s %s\n", style.SprintCyan("●"), "update", style.Dimmed("Update VoinzNext to latest version"))
		fmt.Printf("    %s  %-12s %s\n", style.SprintCyan("●"), "version", style.Dimmed("Print version information"))
		fmt.Println()
		fmt.Printf("  %s %s\n", style.Dimmed("Try:"), style.Value("voinznest init"))
		fmt.Println()
	},
}

func Execute() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		style.Banner("VoinzNext - Help", "Usage information")
		fmt.Printf("  %s %s %s\n\n", style.Dimmed("Use"), style.Value("voinznest [command] --help"), style.Dimmed("for command details"))

		fmt.Printf("  %s\n", style.Label("Commands:"))
		for _, c := range cmd.Commands() {
			if c.Name() == "completion" || c.Name() == "help" || c.Hidden {
				continue
			}
			name := fmt.Sprintf("%-12s", c.Name())
			desc := c.Short
			desc = strings.TrimPrefix(desc, "VoinzNext - ")
			fmt.Printf("    %s  %s %s\n", style.SprintCyan("●"), style.Value(name), style.Dimmed(desc))
		}
		fmt.Println()
		fmt.Printf("  %s %s %s\n", style.Dimmed("Flags:"), style.Value("-h, --help"), style.Dimmed("show help"))
		fmt.Printf("  %s %s %s\n", style.Dimmed("Usage:"), style.Value("voinznest [command] [flags]"), "")
		fmt.Println()
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", style.SprintRed("✘"), err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

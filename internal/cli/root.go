package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "voinznest",
	Short: "VoinzNext - Interactive Next.js Starter Generator",
	Long: `VoinzNext is a CLI tool that helps you scaffold a Next.js project
with your preferred tech stack through an interactive survey.

It generates a complete project structure with all configurations,
components, and dependencies pre-configured. Just answer a few
questions and get started!`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
}

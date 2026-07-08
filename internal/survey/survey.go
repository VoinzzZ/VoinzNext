package survey

import (
	"fmt"
	"os"
	"strings"

	"github.com/VoinzzZ/VoinzNext/internal/config"
	"github.com/VoinzzZ/VoinzNext/internal/registry"
	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/AlecAivazis/survey/v2"
)

func RunSurvey(skipPrompts bool, projectName string) (*config.ProjectConfig, error) {
	cfg := &config.ProjectConfig{}

	if projectName != "" {
		if err := validateProjectName(projectName); err != nil {
			return nil, fmt.Errorf("invalid project name: %w", err)
		}
		cfg.ProjectName = projectName
	} else {
		cfg.ProjectName = "my-next-app"
	}

	dir, _ := os.Getwd()
	cfg.ProjectDir = dir + string(os.PathSeparator) + cfg.ProjectName

	if skipPrompts {
		applyDefaults(cfg)
		style.Banner("VoinzNext - Project Setup", "Using default configuration (--yes)")
		return cfg, nil
	}

	style.Banner("VoinzNext - Project Setup", "Interactive Next.js Starter Generator")

	for _, q := range registry.Questions {
		if err := askQuestion(cfg, q); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func applyDefaults(cfg *config.ProjectConfig) {
	for _, q := range registry.Questions {
		if q.Key == "ProjectName" {
			continue
		}
		setConfigValue(cfg, q.Key, getDefaultName(q), q.Options)

		switch q.Key {
		case "Docker":
			cfg.Docker = q.Default == "true"
		case "ESLintPrettier":
			cfg.ESLintPrettier = q.Default == "true"
		case "InitGit":
			cfg.InitGit = q.Default == "true"
		}
	}
	// DatabaseORM depends on DatabaseType
	if cfg.DatabaseType == "none" {
		cfg.DatabaseORM = "none"
	}
}

func askQuestion(cfg *config.ProjectConfig, q config.Question) error {
	switch q.Key {
	case "ProjectName":
		var name string
		prompt := &survey.Input{
			Message: style.Label(q.Message),
			Default: q.Default,
		}
		if err := survey.AskOne(prompt, &name, survey.WithValidator(validateProjectName)); err != nil {
			return err
		}
		cfg.ProjectName = name
		dir, _ := os.Getwd()
		cfg.ProjectDir = dir + string(os.PathSeparator) + name

	case "Docker", "ESLintPrettier", "InitGit":
		var val string
		options := make([]string, len(q.Options))
		for i, o := range q.Options {
			options[i] = o.Name
		}
		prompt := &survey.Select{
			Message: style.Label(q.Message),
			Options: options,
			Default: getDefaultName(q),
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return err
		}
		boolVal := strings.Contains(strings.ToLower(val), "yes")
		switch q.Key {
		case "Docker":
			cfg.Docker = boolVal
		case "ESLintPrettier":
			cfg.ESLintPrettier = boolVal
		case "InitGit":
			cfg.InitGit = boolVal
		}

	case "DatabaseORM":
		if cfg.DatabaseType == "none" {
			cfg.DatabaseORM = "none"
			return nil
		}
		var filtered []config.Option
		for _, o := range q.Options {
			if cfg.DatabaseType == "mongodb" && o.ID == "drizzle" {
				continue
			}
			filtered = append(filtered, o)
		}
		options := make([]string, len(filtered))
		for i, o := range filtered {
			options[i] = o.Name
		}
		var val string
		prompt := &survey.Select{
			Message: style.Label(q.Message),
			Options: options,
			Default: getDefaultName(config.Question{Options: filtered, Default: q.Default}),
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return err
		}
		setConfigValue(cfg, q.Key, val, filtered)

	default:
		var val string
		options := make([]string, len(q.Options))
		for i, o := range q.Options {
			options[i] = o.Name
		}
		prompt := &survey.Select{
			Message: style.Label(q.Message),
			Options: options,
			Default: getDefaultName(q),
		}
		if err := survey.AskOne(prompt, &val); err != nil {
			return err
		}
		setConfigValue(cfg, q.Key, val, q.Options)
	}

	return nil
}

func getDefaultName(q config.Question) string {
	for _, o := range q.Options {
		if o.ID == q.Default {
			return o.Name
		}
	}
	if len(q.Options) > 0 {
		return q.Options[0].Name
	}
	return ""
}

func setConfigValue(cfg *config.ProjectConfig, key string, selectedName string, options []config.Option) {
	var id string
	for _, o := range options {
		if o.Name == selectedName {
			id = o.ID
			break
		}
	}

	switch key {
	case "Router":
		cfg.Router = id
	case "Language":
		cfg.Language = id
	case "PackageManager":
		cfg.PackageManager = id
	case "CSSFramework":
		cfg.CSSFramework = id
	case "UILibrary":
		cfg.UILibrary = id
	case "DatabaseType":
		cfg.DatabaseType = id
	case "DatabaseORM":
		cfg.DatabaseORM = id
	case "Auth":
		cfg.Auth = id
	case "APIPattern":
		cfg.APIPattern = id
	case "Testing":
		cfg.Testing = id
	}
}

// validateProjectName checks if a project name is valid for filesystem usage.
// It rejects empty strings, illegal characters, and Windows reserved names.
func validateProjectName(val interface{}) error {
	name, ok := val.(string)
	if !ok {
		return fmt.Errorf("project name must be a string")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check for illegal filesystem characters
	illegalChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range illegalChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("project name cannot contain: %s", strings.Join(illegalChars, " "))
		}
	}

	// Check for Windows reserved names
	reserved := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}
	upper := strings.ToUpper(name)
	for _, r := range reserved {
		if upper == r || strings.HasPrefix(upper, r+".") {
			return fmt.Errorf("'%s' is a reserved name on Windows", name)
		}
	}

	return nil
}

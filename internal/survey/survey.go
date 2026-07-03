package survey

import (
	"os"
	"strings"

	"github.com/VoinzzZ/VoinzNext/internal/config"
	"github.com/VoinzzZ/VoinzNext/internal/registry"
	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/AlecAivazis/survey/v2"
)

func RunSurvey() (*config.ProjectConfig, error) {
	cfg := &config.ProjectConfig{}

	style.Banner("VoinzNext - Project Setup", "Interactive Next.js Starter Generator")

	for _, q := range registry.Questions {
		if err := askQuestion(cfg, q); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func askQuestion(cfg *config.ProjectConfig, q config.Question) error {
	switch q.Key {
	case "ProjectName":
		var name string
		prompt := &survey.Input{
			Message: style.Label(q.Message),
			Default: q.Default,
		}
		if err := survey.AskOne(prompt, &name, survey.WithValidator(survey.Required)); err != nil {
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

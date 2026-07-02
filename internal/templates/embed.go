package templates

import (
	"embed"
)

//go:embed files/*
var TemplateFS embed.FS

//go:embed files/.gitignore
var gitignoreContent string

func GetGitIgnore() string {
	return gitignoreContent
}

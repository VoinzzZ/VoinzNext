package config

import "errors"

// ErrDirNotEmpty is returned when the target project directory already exists
// and contains files. The caller should prompt the user for confirmation and
// set Overwrite = true before retrying.
var ErrDirNotEmpty = errors.New("project directory already exists and is not empty")

type ProjectConfig struct {
	ProjectName    string
	Router         string
	Language       string
	PackageManager string
	CSSFramework   string
	UILibrary      string
	DatabaseType   string
	DatabaseORM    string
	Auth           string
	APIPattern     string
	Testing        string
	Docker         bool
	ESLintPrettier bool
	InitGit        bool
	Overwrite      bool
	ProjectDir     string
}

type Option struct {
	ID          string
	Name        string
	Description string
}

type Question struct {
	Key     string
	Message string
	Options []Option
	Default string
}

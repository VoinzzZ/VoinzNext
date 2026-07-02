package config

type ProjectConfig struct {
	ProjectName     string
	Router          string
	Language        string
	PackageManager  string
	CSSFramework    string
	UILibrary       string
	DatabaseORM     string
	Auth            string
	APIPattern      string
	Testing         string
	Docker          bool
	ESLintPrettier  bool
	InitGit         bool
	ProjectDir      string
}

type Option struct {
	ID          string
	Name        string
	Description string
}

type Question struct {
	Key         string
	Message     string
	Options     []Option
	Default     string
}

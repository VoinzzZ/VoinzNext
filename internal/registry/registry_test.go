package registry

import (
	"testing"

	"github.com/VoinzzZ/VoinzNext/internal/config"
)

func TestGetDependencies_TypeScript(t *testing.T) {
	cfg := &config.ProjectConfig{
		Router:         "app",
		Language:       "typescript",
		CSSFramework:   "tailwind",
		DatabaseORM:    "none",
		Auth:           "none",
		APIPattern:     "none",
		Testing:        "none",
		ESLintPrettier: false,
	}

	deps := GetDependencies(cfg)

	tests := []struct {
		key      string
		expected bool
	}{
		{"next", true},
		{"react", true},
		{"react-dom", true},
		{"dev:typescript", true},
		{"dev:@types/react", true},
		{"dev:tailwindcss", true},
	}

	for _, tt := range tests {
		_, ok := deps[tt.key]
		if ok != tt.expected {
			t.Errorf("dependency %s: expected %v, got %v", tt.key, tt.expected, ok)
		}
	}
}

func TestGetDependencies_FullStack(t *testing.T) {
	cfg := &config.ProjectConfig{
		Router:         "app",
		Language:       "typescript",
		CSSFramework:   "tailwind",
		UILibrary:      "shadcn",
		DatabaseORM:    "prisma",
		Auth:           "nextauth",
		APIPattern:     "trpc",
		Testing:        "vitest",
		Docker:         true,
		ESLintPrettier: true,
	}

	deps := GetDependencies(cfg)

	required := []string{
		"next", "react", "react-dom",
		"@prisma/client",
		"dev:prisma",
		"next-auth",
		"@trpc/server", "@trpc/client", "@trpc/react-query",
		"@tanstack/react-query",
		"dev:lucide-react",
		"dev:class-variance-authority",
		"dev:vitest",
		"dev:@testing-library/react",
		"dev:eslint",
		"dev:prettier",
	}

	for _, key := range required {
		if _, ok := deps[key]; !ok {
			t.Errorf("expected dependency %s to be present, but it's missing", key)
		}
	}
}

func TestQuestionsCount(t *testing.T) {
	if len(Questions) == 0 {
		t.Error("expected at least one question in registry")
	}
}

func TestQuestionsHaveKeys(t *testing.T) {
	for i, q := range Questions {
		if q.Key == "" {
			t.Errorf("question %d has empty key", i)
		}
		if q.Message == "" {
			t.Errorf("question %d has empty message", i)
		}
	}
}

package generator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/VoinzzZ/VoinzNext/internal/config"
)

func TestGenerator_Generate(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "test-project")

	cfg := &config.ProjectConfig{
		ProjectName:    "test-project",
		ProjectDir:     projectDir,
		Router:         "app",
		Language:       "typescript",
		PackageManager: "pnpm",
		CSSFramework:   "tailwind",
		UILibrary:      "shadcn",
		DatabaseType:   "postgresql",
		DatabaseORM:    "prisma",
		Auth:           "nextauth",
		APIPattern:     "trpc",
		Testing:        "vitest",
		Docker:         true,
		ESLintPrettier: true,
		InitGit:        false,
	}

	g := New(cfg)
	if err := g.Generate(); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	expectedFiles := []string{
		"package.json",
		"next.config.mjs",
		"tsconfig.json",
		"tailwind.config.ts",
		"postcss.config.js",
		"Dockerfile",
		"docker-compose.yml",
		".env.example",
		"README.md",
		"src/styles/globals.css",
		"src/app/layout.tsx",
		"src/app/page.tsx",
		"src/lib/utils.ts",
		"src/components/ui/button.tsx",
		"prisma/schema.prisma",
		"src/db/index.ts",
		"src/server/trpc.ts",
		"vitest.config.ts",
		".eslintrc.js",
		".prettierrc",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(projectDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist, but it doesn't", f)
		}
	}

	pkgPath := filepath.Join(projectDir, "package.json")
	if data, err := os.ReadFile(pkgPath); err != nil {
		t.Errorf("failed to read package.json: %v", err)
	} else if !json.Valid(data) {
		t.Errorf("package.json is not valid JSON:\n%s", string(data))
	}

	dirs := []string{
		"src/components/ui",
		"src/app",
		"src/lib",
		"src/styles",
		"src/db",
		"src/server",
		"prisma",
	}

	for _, d := range dirs {
		path := filepath.Join(projectDir, d)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected directory %s to exist: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", d)
		}
	}
}

func TestGenerator_GenerateMinimal(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "minimal-app")

	cfg := &config.ProjectConfig{
		ProjectName:    "minimal-app",
		ProjectDir:     projectDir,
		Router:         "app",
		Language:       "javascript",
		PackageManager: "npm",
		CSSFramework:   "none",
		UILibrary:      "none",
		DatabaseType:   "none",
		DatabaseORM:    "none",
		Auth:           "none",
		APIPattern:     "none",
		Testing:        "none",
		Docker:         false,
		ESLintPrettier: false,
		InitGit:        false,
	}

	g := New(cfg)
	if err := g.Generate(); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	expectedFiles := []string{
		"package.json",
		"next.config.js",
		"README.md",
		"src/styles/globals.css",
		"src/app/layout.jsx",
		"src/app/page.jsx",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(projectDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist, but it doesn't", f)
		}
	}

	notExpected := []string{
		"tsconfig.json",
		"tailwind.config.js",
		"Dockerfile",
		"prisma/schema.prisma",
		"src/server/trpc.ts",
		"vitest.config.ts",
	}

	for _, f := range notExpected {
		path := filepath.Join(projectDir, f)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("expected file %s to NOT exist, but it does", f)
		}
	}
}

package generator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func TestGenerator_PackageJSON_Deterministic(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "det-test")

	cfg := &config.ProjectConfig{
		ProjectName:    "det-test",
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

	// Generate twice and compare package.json output
	g := New(cfg)
	if err := g.Generate(); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	first, err := os.ReadFile(filepath.Join(projectDir, "package.json"))
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}

	// Regenerate
	cfg.ProjectDir = filepath.Join(tmpDir, "det-test-2")
	g2 := New(cfg)
	if err := g2.Generate(); err != nil {
		t.Fatalf("Generate() 2nd error: %v", err)
	}
	second, err := os.ReadFile(filepath.Join(cfg.ProjectDir, "package.json"))
	if err != nil {
		t.Fatalf("read package.json 2nd: %v", err)
	}

	if string(first) != string(second) {
		t.Errorf("package.json output is non-deterministic:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}

	// Verify pnpm build-script approvals are written via supported project config files.
	workspaceConfig, err := os.ReadFile(filepath.Join(cfg.ProjectDir, "pnpm-workspace.yaml"))
	if err != nil {
		t.Fatalf("read pnpm-workspace.yaml: %v", err)
	}
	for _, required := range []string{"allowBuilds:", "esbuild: true", "sharp: true", "unrs-resolver: true"} {
		if !strings.Contains(string(workspaceConfig), required) {
			t.Errorf("pnpm-workspace.yaml missing %q", required)
		}
	}

	npmrc, err := os.ReadFile(filepath.Join(cfg.ProjectDir, ".npmrc"))
	if err != nil {
		t.Fatalf("read .npmrc: %v", err)
	}
	if !strings.Contains(string(npmrc), "enable-pre-post-scripts=true") {
		t.Error(".npmrc missing enable-pre-post-scripts=true")
	}

	// Verify keys in dependencies are sorted
	var full map[string]interface{}
	json.Unmarshal(first, &full)
	content := string(first)
	if deps, ok := full["dependencies"].(map[string]interface{}); ok && len(deps) > 1 {
		// Check that keys appear in sorted order in the raw JSON
		lastIdx := -1
		for _, key := range sortedKeys(deps) {
			idx := strings.Index(content, "\""+key+"\"")
			if idx <= lastIdx {
				t.Errorf("dependencies keys not sorted: %s appeared before expected position", key)
			}
			lastIdx = idx
		}
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// already need to sort to verify
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
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
		"jsconfig.json",
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

func TestGenerator_JavaScriptMode_NoTypeScriptSyntax(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "js-app")

	cfg := &config.ProjectConfig{
		ProjectName:    "js-app",
		ProjectDir:     projectDir,
		Router:         "app",
		Language:       "javascript",
		PackageManager: "npm",
		CSSFramework:   "tailwind",
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

	// App Router layout.jsx must not contain TypeScript syntax
	layoutPath := filepath.Join(projectDir, "src", "app", "layout.jsx")
	layoutData, err := os.ReadFile(layoutPath)
	if err != nil {
		t.Fatalf("read layout.jsx: %v", err)
	}
	layout := string(layoutData)

	tsForbidden := []string{
		"import type",
		": Metadata",
		"Readonly<",
		"React.ReactNode",
		": React.",
	}
	for _, ts := range tsForbidden {
		if strings.Contains(layout, ts) {
			t.Errorf("layout.jsx contains TypeScript syntax %q", ts)
		}
	}

	// Pages Router: test _app.jsx
	projectDir2 := filepath.Join(tmpDir, "js-pages-app")
	cfg2 := &config.ProjectConfig{
		ProjectName:    "js-pages-app",
		ProjectDir:     projectDir2,
		Router:         "pages",
		Language:       "javascript",
		PackageManager: "npm",
		CSSFramework:   "tailwind",
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

	g2 := New(cfg2)
	if err := g2.Generate(); err != nil {
		t.Fatalf("Generate() pages error: %v", err)
	}

	appPath := filepath.Join(projectDir2, "src", "pages", "_app.jsx")
	appData, err := os.ReadFile(appPath)
	if err != nil {
		t.Fatalf("read _app.jsx: %v", err)
	}
	app := string(appData)

	appTsForbidden := []string{
		"import type",
		": AppProps",
		"AppProps",
	}
	for _, ts := range appTsForbidden {
		if strings.Contains(app, ts) {
			t.Errorf("_app.jsx contains TypeScript syntax %q", ts)
		}
	}
}

func TestGenerator_ErrDirNotEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "existing-project")

	// Create dir with a file inside to simulate existing project
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "important.txt"), []byte("do not delete"), 0644)

	cfg := &config.ProjectConfig{
		ProjectName:    "existing-project",
		ProjectDir:     projectDir,
		Router:         "app",
		Language:       "typescript",
		PackageManager: "pnpm",
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
		Overwrite:      false, // not allowed
	}

	g := New(cfg)
	err := g.Generate()

	// Must return ErrDirNotEmpty
	if err == nil {
		t.Fatal("expected ErrDirNotEmpty, got nil")
	}
	if err != config.ErrDirNotEmpty {
		t.Fatalf("expected ErrDirNotEmpty, got: %v", err)
	}

	// Original file must still exist
	if _, statErr := os.Stat(filepath.Join(projectDir, "important.txt")); os.IsNotExist(statErr) {
		t.Error("important.txt was deleted without confirmation!")
	}
}

func TestGenerator_OverwriteExistingDir(t *testing.T) {
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "overwrite-project")

	// Create dir with a file inside
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "old-file.txt"), []byte("old"), 0644)

	cfg := &config.ProjectConfig{
		ProjectName:    "overwrite-project",
		ProjectDir:     projectDir,
		Router:         "app",
		Language:       "typescript",
		PackageManager: "pnpm",
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
		Overwrite:      true, // explicitly allowed
	}

	g := New(cfg)
	if err := g.Generate(); err != nil {
		t.Fatalf("Generate() with Overwrite=true error: %v", err)
	}

	// Old file should be gone
	if _, err := os.Stat(filepath.Join(projectDir, "old-file.txt")); !os.IsNotExist(err) {
		t.Error("old-file.txt still exists after overwrite")
	}

	// New project files should exist
	if _, err := os.Stat(filepath.Join(projectDir, "package.json")); os.IsNotExist(err) {
		t.Error("package.json not created after overwrite")
	}
}

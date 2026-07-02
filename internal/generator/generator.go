package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/VoinzzZ/VoinzNext/internal/config"
	"github.com/VoinzzZ/VoinzNext/internal/registry"
	"github.com/VoinzzZ/VoinzNext/internal/style"
	"github.com/VoinzzZ/VoinzNext/internal/templates"
)

type Generator struct {
	cfg *config.ProjectConfig
}

func New(cfg *config.ProjectConfig) *Generator {
	return &Generator{cfg: cfg}
}

func (g *Generator) Generate() error {
	dir := g.cfg.ProjectDir

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create project directory: %w", err)
	}

	fmt.Printf("  %s %s\n", style.Label("◆"), fmt.Sprintf("Scaffolding %s...", style.Value(g.cfg.ProjectName)))
	fmt.Println()

	steps := []struct {
		name string
		fn   func(string) error
	}{
		{"Directory structure", g.scaffoldDirs},
		{"package.json with dependencies", g.writePackageJSON},
		{"Config files (next, ts, tailwind, postcss)", g.writeConfigFiles},
		{"Source files (layout, pages, components)", g.writeSourceFiles},
		{".env.example", g.writeEnvFile},
		{"README.md", g.writeReadme},
		{"Dockerfile & docker-compose", g.writeDockerFiles},
		{"Database schema & client", g.writeDatabaseFiles},
		{"Auth provider setup", g.writeAuthFiles},
		{"API layer (tRPC)", g.writeAPIFiles},
		{"Test framework config", g.writeTestFiles},
		{"UI component library", g.writeUIFiles},
		{"ESLint & Prettier config", g.writeLintFiles},
	}

	for _, step := range steps {
		style.StepRunning(step.name)
		if err := step.fn(dir); err != nil {
			style.StepError(step.name, err)
			return fmt.Errorf("%s: %w", step.name, err)
		}
		style.StepDone(step.name)
	}

	return nil
}

func (g *Generator) scaffoldDirs(dir string) error {
	dirs := []string{
		"src/app",
		"src/components",
		"src/lib",
		"src/styles",
		"public",
	}

	if g.cfg.Testing != "none" {
		dirs = append(dirs, "src/tests", "src/__tests__")
	}
	if g.cfg.DatabaseORM != "none" {
		dirs = append(dirs, "src/db")
	}
	if g.cfg.Auth != "none" {
		dirs = append(dirs, "src/auth")
	}
	if g.cfg.APIPattern != "none" {
		dirs = append(dirs, "src/server")
	}
	if g.cfg.UILibrary == "shadcn" {
		dirs = append(dirs, "src/components/ui")
	}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(dir, d), 0755); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) writePackageJSON(dir string) error {
	deps := registry.GetDependencies(g.cfg)

	pkg := map[string]interface{}{
		"name":            g.cfg.ProjectName,
		"version":         "0.1.0",
		"private":         true,
		"scripts":         g.getScripts(),
		"dependencies":    map[string]string{},
		"devDependencies": map[string]string{},
	}

	for k, v := range deps {
		if strings.HasPrefix(k, "dev:") {
			pkg["devDependencies"].(map[string]string)[k[4:]] = v
		} else {
			pkg["dependencies"].(map[string]string)[k] = v
		}
	}

	return writeJSON(filepath.Join(dir, "package.json"), pkg)
}

func (g *Generator) getScripts() map[string]string {
	scripts := map[string]string{
		"dev":   "next dev",
		"build": "next build",
		"start": "next start",
		"lint":  "next lint",
	}

	if g.cfg.Testing == "vitest" {
		scripts["test"] = "vitest"
		scripts["test:ui"] = "vitest --ui"
	} else if g.cfg.Testing == "jest" {
		scripts["test"] = "jest"
	} else if g.cfg.Testing == "playwright" {
		scripts["test:e2e"] = "playwright test"
	}

	if g.cfg.DatabaseORM == "prisma" {
		scripts["db:generate"] = "prisma generate"
		scripts["db:push"] = "prisma db push"
		scripts["db:studio"] = "prisma studio"
		scripts["db:migrate"] = "prisma migrate dev"
	} else if g.cfg.DatabaseORM == "drizzle" {
		scripts["db:generate"] = "drizzle-kit generate"
		scripts["db:push"] = "drizzle-kit push"
		scripts["db:studio"] = "drizzle-kit studio"
		scripts["db:migrate"] = "drizzle-kit migrate"
	}

	if g.cfg.ESLintPrettier {
		scripts["format"] = "prettier --write \"src/**/*.{ts,tsx,js,jsx,json,css,md}\""
	}

	return scripts
}

func (g *Generator) writeConfigFiles(dir string) error {
	ext := ".js"
	if g.cfg.Language == "typescript" {
		ext = ".ts"
	}

	if err := writeFile(filepath.Join(dir, "next.config"+ext), readTemplateFile("next.config"+ext)); err != nil {
		return err
	}

	if g.cfg.Language == "typescript" {
		if err := writeFile(filepath.Join(dir, "tsconfig.json"), readTemplateFile("tsconfig.json")); err != nil {
			return err
		}
	}

	if g.cfg.CSSFramework == "tailwind" {
		if err := writeFile(filepath.Join(dir, "tailwind.config"+ext), g.renderTailwindConfig()); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(dir, "postcss.config.js"), readTemplateFile("postcss.config.js")); err != nil {
			return err
		}
	}

	if err := writeFile(filepath.Join(dir, ".npmrc"), readTemplateFile(".npmrc")); err != nil {
		return err
	}

	return nil
}

func (g *Generator) renderTailwindConfig() string {
	return readTemplateFile("tailwind.config.ts")
}

func (g *Generator) writeSourceFiles(dir string) error {
	cssContent := readTemplateFile("globals.css")
	srcStyles := filepath.Join(dir, "src", "styles", "globals.css")
	if g.cfg.CSSFramework == "none" {
		cssContent = strings.ReplaceAll(cssContent, "@tailwind base;\n@tailwind components;\n@tailwind utilities;\n", "")
	}
	if err := writeFile(srcStyles, cssContent); err != nil {
		return err
	}

	ext := ".tsx"
	if g.cfg.Language == "javascript" {
		ext = ".jsx"
	}

	if g.cfg.Router == "app" {
		if err := writeFile(filepath.Join(dir, "src", "app", "layout"+ext), defaultLayout(g.cfg)); err != nil {
			return err
		}
		pageContent := g.renderPage()
		if err := writeFile(filepath.Join(dir, "src", "app", "page"+ext), pageContent); err != nil {
			return err
		}
	} else {
		if err := writeFile(filepath.Join(dir, "src", "pages", "_app."+ext), defaultApp(g.cfg)); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(dir, "src", "pages", "index."+ext), defaultPage(g.cfg)); err != nil {
			return err
		}
	}

	if g.cfg.UILibrary == "shadcn" {
		utilsContent := readTemplateFile("lib-utils.ts")
		if err := writeFile(filepath.Join(dir, "src", "lib", "utils.ts"), utilsContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderPage() string {
	if g.cfg.Language == "javascript" {
		return readTemplateFile("page.js")
	}
	return readTemplateFile("page.tsx")
}

func defaultLayout(cfg *config.ProjectConfig) string {
	interLine := `import { Inter } from "next/font/google";`
	interClass := `className={inter.className}`
	cssImport := `import "@/styles/globals.css";`

	if cfg.CSSFramework == "none" {
		cssImport = ""
	}

	return fmt.Sprintf(`import type { Metadata } from "next";
%s
%s

export const metadata: Metadata = {
  title: "%s",
  description: "Generated by VoinzNext",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body %s>{children}</body>
    </html>
  );
}
`, interLine, cssImport, cfg.ProjectName, interClass)
}

func defaultApp(cfg *config.ProjectConfig) string {
	return `import type { AppProps } from "next/app";
import "@/styles/globals.css";

export default function App({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />;
}
`
}

func defaultPage(cfg *config.ProjectConfig) string {
	return `export default function Home() {
  return (
    <div>
      <h1>Welcome to ` + cfg.ProjectName + `</h1>
    </div>
  );
}
`
}

func (g *Generator) writeEnvFile(dir string) error {
	envContent := readTemplateFile(".env.example")
	return writeFile(filepath.Join(dir, ".env.example"), envContent)
}

func (g *Generator) writeReadme(dir string) error {
	readme := fmt.Sprintf(`# %s

This project was generated with [VoinzNext](https://github.com/VoinzzZ/VoinzNext).

## Tech Stack
- **Router:** %s Router
- **Language:** %s
- **CSS:** %s
- **UI Library:** %s
- **Database:** %s
- **Auth:** %s
- **API:** %s
- **Testing:** %s
- **Package Manager:** %s

## Getting Started

First, install dependencies:

`+"```bash"+`
%s install
`+"```"+`

Then, run the development server:

`+"```bash"+`
%s run dev
`+"```"+`

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Environment Variables

Copy the example env file and fill in your values:

`+"```bash"+`
cp .env.example .env
`+"```"+`

## Learn More

- [Next.js Documentation](https://nextjs.org/docs)
`,
		g.cfg.ProjectName,
		g.cfg.Router,
		g.cfg.Language,
		g.cfg.CSSFramework,
		g.cfg.UILibrary,
		g.cfg.DatabaseORM,
		g.cfg.Auth,
		g.cfg.APIPattern,
		g.cfg.Testing,
		g.cfg.PackageManager,
		g.cfg.PackageManager,
		g.cfg.PackageManager,
	)

	return writeFile(filepath.Join(dir, "README.md"), readme)
}

func (g *Generator) writeDockerFiles(dir string) error {
	if !g.cfg.Docker {
		return nil
	}

	if err := writeFile(filepath.Join(dir, "Dockerfile"), readTemplateFile("Dockerfile")); err != nil {
		return err
	}

	if g.cfg.DatabaseORM != "none" || g.cfg.Auth != "none" {
		content := readTemplateFile("docker-compose.yml")
		content = strings.ReplaceAll(content, "${PROJECT_NAME}", g.cfg.ProjectName)
		if err := writeFile(filepath.Join(dir, "docker-compose.yml"), content); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) writeDatabaseFiles(dir string) error {
	switch g.cfg.DatabaseORM {
	case "prisma":
		prismaDir := filepath.Join(dir, "prisma")
		if err := os.MkdirAll(prismaDir, 0755); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(prismaDir, "schema.prisma"), readTemplateFile("schema.prisma")); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(dir, "src", "db", "index.ts"), readTemplateFile("db-prisma.ts")); err != nil {
			return err
		}

	case "drizzle":
		if err := writeFile(filepath.Join(dir, "drizzle.config.ts"), readTemplateFile("drizzle.config.ts")); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(dir, "src", "db", "index.ts"), readTemplateFile("db-drizzle.ts")); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(dir, "src", "db", "schema.ts"), readTemplateFile("schema-drizzle.ts")); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) writeAuthFiles(dir string) error {
	if g.cfg.Auth == "nextauth" {
		authDir := filepath.Join(dir, "src", "auth")
		if err := writeFile(filepath.Join(authDir, "index.ts"), readTemplateFile("auth-nextauth.ts")); err != nil {
			return err
		}
	}

	if g.cfg.Auth == "clerk" {
		content := `import { ClerkProvider } from "@clerk/nextjs";

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return <ClerkProvider>{children}</ClerkProvider>;
}`
		if err := writeFile(filepath.Join(dir, "src", "auth", "provider.tsx"), content); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) writeAPIFiles(dir string) error {
	if g.cfg.APIPattern == "trpc" {
		serverDir := filepath.Join(dir, "src", "server")
		if err := writeFile(filepath.Join(serverDir, "trpc.ts"), readTemplateFile("trpc-server.ts")); err != nil {
			return err
		}

		apiDir := filepath.Join(dir, "src", "app", "api", "trpc")
		if g.cfg.Router == "pages" {
			apiDir = filepath.Join(dir, "src", "pages", "api", "trpc")
		}
		if err := os.MkdirAll(apiDir, 0755); err != nil {
			return err
		}

		routeContent := `import { fetchRequestHandler } from "@trpc/server/adapters/fetch";
import { appRouter } from "@/server/trpc";

const handler = (req: Request) =>
  fetchRequestHandler({
    endpoint: "/api/trpc",
    req,
    router: appRouter,
    createContext: () => ({}),
  });

export { handler as GET, handler as POST };
`
		if err := writeFile(filepath.Join(apiDir, "[...trpc]", "route.ts"), routeContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) writeTestFiles(dir string) error {
	switch g.cfg.Testing {
	case "vitest":
		if err := writeFile(filepath.Join(dir, "vitest.config.ts"), readTemplateFile("vitest.config.ts")); err != nil {
			return err
		}
		setupContent := `import "@testing-library/jest-dom";`
		if err := writeFile(filepath.Join(dir, "src", "tests", "setup.ts"), setupContent); err != nil {
			return err
		}

	case "jest":
		if err := writeFile(filepath.Join(dir, "jest.config.ts"), readTemplateFile("jest.config.ts")); err != nil {
			return err
		}
		setupContent := `import "@testing-library/jest-dom";`
		if err := writeFile(filepath.Join(dir, "src", "tests", "setup.ts"), setupContent); err != nil {
			return err
		}

	case "playwright":
		pwDir := filepath.Join(dir, "e2e")
		if err := os.MkdirAll(pwDir, 0755); err != nil {
			return err
		}
		testContent := `import { test, expect } from "@playwright/test";

test("homepage loads", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/.*/);
});
`
		if err := writeFile(filepath.Join(pwDir, "home.spec.ts"), testContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) writeUIFiles(dir string) error {
	if g.cfg.UILibrary != "shadcn" {
		return nil
	}

	buttonContent := `import { forwardRef, ButtonHTMLAttributes } from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        destructive: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
        outline: "border border-input bg-background hover:bg-accent hover:text-accent-foreground",
        secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80",
        ghost: "hover:bg-accent hover:text-accent-foreground",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 rounded-md px-3",
        lg: "h-11 rounded-md px-8",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

export interface ButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => {
    return (
      <button
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";

export { Button, buttonVariants };
`
	if err := writeFile(filepath.Join(dir, "src", "components", "ui", "button.tsx"), buttonContent); err != nil {
		return err
	}

	return nil
}

func (g *Generator) writeLintFiles(dir string) error {
	if !g.cfg.ESLintPrettier {
		return nil
	}

	if err := writeFile(filepath.Join(dir, ".eslintrc.js"), readTemplateFile("eslint.config.js")); err != nil {
		return err
	}

	if err := writeFile(filepath.Join(dir, ".prettierrc"), readTemplateFile(".prettierrc")); err != nil {
		return err
	}

	return nil
}

func (g *Generator) PostGenerate() error {
	dir := g.cfg.ProjectDir

	if g.cfg.DatabaseORM == "prisma" {
		style.StepRunning("Generating Prisma client")
		if err := g.runCmd(dir, "npx", "prisma", "generate"); err != nil {
			style.StepWarn("Prisma generate failed", fmt.Sprintf("%v", err))
		} else {
			style.StepDone("Prisma client generated")
		}
	}

	if g.cfg.InitGit {
		style.StepRunning("Initializing git repository")
		gitDir := filepath.Join(dir, ".gitignore")
		if err := writeFile(gitDir, templates.GetGitIgnore()); err != nil {
			return err
		}
		if err := g.runCmd(dir, "git", "init"); err != nil {
			return fmt.Errorf("git init: %w", err)
		}
		if err := g.runCmd(dir, "git", "add", "."); err != nil {
			return err
		}
		if err := g.runCmd(dir, "git", "commit", "-m", "Initial commit: generated by VoinzNext"); err != nil {
			style.StepWarn("Git commit failed", "may need git config")
		} else {
			style.StepDone("Git repository initialized")
		}
	}

	return nil
}

func (g *Generator) runCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func readTemplateFile(name string) string {
	data, err := templates.TemplateFS.ReadFile("files/" + name)
	if err != nil {
		return ""
	}
	return string(data)
}

func writeFile(path, content string) error {
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func escapeJSONString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func writeJSON(path string, data interface{}) error {
	content := "{\n"
	switch m := data.(type) {
	case map[string]interface{}:
		if name, ok := m["name"]; ok {
			content += fmt.Sprintf("  \"name\": \"%s\",\n", escapeJSONString(fmt.Sprintf("%v", name)))
		}
		if ver, ok := m["version"]; ok {
			content += fmt.Sprintf("  \"version\": \"%s\",\n", escapeJSONString(fmt.Sprintf("%v", ver)))
		}
		if priv, ok := m["private"]; ok {
			content += fmt.Sprintf("  \"private\": %v,\n", priv)
		}
		if scripts, ok := m["scripts"].(map[string]string); ok {
			content += "  \"scripts\": {\n"
			first := true
			for k, v := range scripts {
				if !first {
					content += ",\n"
				}
				content += fmt.Sprintf("    \"%s\": \"%s\"", k, escapeJSONString(v))
				first = false
			}
			content += "\n  },\n"
		}
		if deps, ok := m["dependencies"].(map[string]string); ok && len(deps) > 0 {
			content += "  \"dependencies\": {\n"
			first := true
			for k, v := range deps {
				if !first {
					content += ",\n"
				}
				content += fmt.Sprintf("    \"%s\": \"%s\"", k, escapeJSONString(v))
				first = false
			}
			content += "\n  },\n"
		}
		if devDeps, ok := m["devDependencies"].(map[string]string); ok && len(devDeps) > 0 {
			content += "  \"devDependencies\": {\n"
			first := true
			for k, v := range devDeps {
				if !first {
					content += ",\n"
				}
				content += fmt.Sprintf("    \"%s\": \"%s\"", k, escapeJSONString(v))
				first = false
			}
			content += "\n  }\n"
		}
	}
	content += "}\n"
	return writeFile(path, content)
}

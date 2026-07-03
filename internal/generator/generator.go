package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

	if _, err := os.Stat(dir); err == nil {
		entries, _ := os.ReadDir(dir)
		if len(entries) > 0 {
			if !g.cfg.Overwrite {
				return config.ErrDirNotEmpty
			}
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("remove existing directory: %w", err)
			}
		}
	}
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
		{"Database setup", g.writeDatabaseFiles},
		{fmt.Sprintf("API layer (%s)", g.cfg.APIPattern), g.writeAPIFiles},
		{"Auth setup", g.writeAuthFiles},
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
	if g.cfg.DatabaseType != "none" {
		dirs = append(dirs, "src/db")
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

	depsMap := sortedMap{}
	devDepsMap := sortedMap{}

	for k, v := range deps {
		if strings.HasPrefix(k, "dev:") {
			devDepsMap = append(devDepsMap, sortedMapEntry{Key: k[4:], Value: v})
		} else {
			depsMap = append(depsMap, sortedMapEntry{Key: k, Value: v})
		}
	}
	sort.Sort(depsMap)
	sort.Sort(devDepsMap)

	pkg := orderedPackageJSON{
		Name:            g.cfg.ProjectName,
		Version:         "0.1.0",
		Private:         true,
		Scripts:         newSortedMapFromMap(g.getScripts()),
		Dependencies:    depsMap,
		DevDependencies: devDepsMap,
		Pnpm: &pnpmConfig{
			OnlyBuiltDependencies: []string{"unrs-resolver"},
		},
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
	nextCfgExt := ".js"
	twCfgExt := ".js"
	if g.cfg.Language == "typescript" {
		nextCfgExt = ".mjs"
		twCfgExt = ".ts"
	}

	if err := writeFile(filepath.Join(dir, "next.config"+nextCfgExt), readTemplateFile("next.config"+nextCfgExt)); err != nil {
		return err
	}

	if g.cfg.Language == "typescript" {
		if err := writeFile(filepath.Join(dir, "tsconfig.json"), readTemplateFile("tsconfig.json")); err != nil {
			return err
		}
	}

	if g.cfg.CSSFramework == "tailwind" {
		if err := writeFile(filepath.Join(dir, "tailwind.config"+twCfgExt), g.renderTailwindConfig()); err != nil {
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
	if g.cfg.Language == "javascript" {
		return `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}`
	}
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
		if err := writeFile(filepath.Join(dir, "src", "pages", "_app"+ext), defaultApp(g.cfg)); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(dir, "src", "pages", "index"+ext), defaultPage(g.cfg)); err != nil {
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
	interInit := `const inter = Inter({ subsets: ["latin"] });`
	interClass := `className={inter.className}`
	cssImport := `import "@/styles/globals.css";`

	if cfg.CSSFramework == "none" {
		cssImport = ""
	}

	if cfg.Language == "javascript" {
		return fmt.Sprintf(`%s
%s

%s

export const metadata = {
  title: "%s",
  description: "Generated by VoinzNext",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body %s>{children}</body>
    </html>
  );
}
`, interLine, cssImport, interInit, cfg.ProjectName, interClass)
	}

	return fmt.Sprintf(`import type { Metadata } from "next";
%s
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
`, interLine, cssImport, interInit, cfg.ProjectName, interClass)
}

func defaultApp(cfg *config.ProjectConfig) string {
	if cfg.Language == "javascript" {
		return `import "@/styles/globals.css";

export default function App({ Component, pageProps }) {
  return <Component {...pageProps} />;
}
`
	}

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
	content := "# Next.js\nNEXT_PUBLIC_APP_URL=http://localhost:3000\n\n"

	switch g.cfg.DatabaseType {
	case "mysql":
		content += "# Database (MySQL)\nDATABASE_URL=\"mysql://user:password@localhost:3306/" + g.cfg.ProjectName + "\"\n"
	case "postgresql":
		content += "# Database (PostgreSQL)\nDATABASE_URL=\"postgresql://user:password@localhost:5432/" + g.cfg.ProjectName + "\"\n"
	case "mongodb":
		content += "# Database (MongoDB)\nDATABASE_URL=\"mongodb://user:password@localhost:27017/" + g.cfg.ProjectName + "\"\n"
	case "supabase":
		content += "# Database (Supabase)\nDATABASE_URL=\"postgresql://user:password@db.supabase.co:5432/postgres\"\n"
		content += "NEXT_PUBLIC_SUPABASE_URL=your-project-url\n"
		content += "NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key\n"
	case "tidb":
		content += "# Database (TiDB)\nDATABASE_URL=\"mysql://user:password@host:4000/" + g.cfg.ProjectName + "?sslaccept=strict\"\n"
	}

	if g.cfg.APIPattern == "trpc" {
		content += "\n# tRPC\nNEXT_PUBLIC_TRPC_API_URL=http://localhost:3000/api/trpc\n"
	}

	if g.cfg.Auth == "nextauth" {
		content += "\n# NextAuth.js\nAUTH_SECRET=your-secret-key\n"
		content += "AUTH_URL=http://localhost:3000\n"
	} else if g.cfg.Auth == "clerk" {
		content += "\n# Clerk\nNEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=your-publishable-key\n"
		content += "CLERK_SECRET_KEY=your-secret-key\n"
	} else if g.cfg.Auth == "lucia" {
		content += "\n# Lucia Auth\nAUTH_SECRET=your-secret-key\n"
	}

	if g.cfg.Docker {
		content += "\n# Docker\nPOSTGRES_USER=user\nPOSTGRES_PASSWORD=password\nPOSTGRES_DB=" + g.cfg.ProjectName + "\n"
	}

	return writeFile(filepath.Join(dir, ".env.example"), content)
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
- **ORM:** %s
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
		g.cfg.DatabaseType,
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

	if g.cfg.DatabaseType != "none" || g.cfg.Auth != "none" {
		content := readTemplateFile("docker-compose.yml")
		content = strings.ReplaceAll(content, "${PROJECT_NAME}", g.cfg.ProjectName)
		if err := writeFile(filepath.Join(dir, "docker-compose.yml"), content); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) writeDatabaseFiles(dir string) error {
	if g.cfg.DatabaseType == "none" {
		return nil
	}

	ext := ".ts"
	if g.cfg.Language == "javascript" {
		ext = ".js"
	}

	switch g.cfg.DatabaseORM {
	case "prisma":
		prismaDir := filepath.Join(dir, "prisma")
		if err := os.MkdirAll(prismaDir, 0755); err != nil {
			return err
		}
		schema := g.renderPrismaSchema()
		if err := writeFile(filepath.Join(prismaDir, "schema.prisma"), schema); err != nil {
			return err
		}
		if g.cfg.DatabaseType == "tidb" {
			if err := writeFile(filepath.Join(dir, "src", "db", "index"+ext), readTemplateFile("db-prisma-tidb.ts")); err != nil {
				return err
			}
		} else {
			if err := writeFile(filepath.Join(dir, "src", "db", "index"+ext), readTemplateFile("db-prisma.ts")); err != nil {
				return err
			}
		}

	case "drizzle":
		if err := writeFile(filepath.Join(dir, "drizzle.config.ts"), readTemplateFile("drizzle.config.ts")); err != nil {
			return err
		}
		switch g.cfg.DatabaseType {
		case "mysql", "tidb":
			if err := writeFile(filepath.Join(dir, "src", "db", "index"+ext), readTemplateFile("db-drizzle-mysql.ts")); err != nil {
				return err
			}
		default:
			if err := writeFile(filepath.Join(dir, "src", "db", "index"+ext), readTemplateFile("db-drizzle.ts")); err != nil {
				return err
			}
		}
		if err := writeFile(filepath.Join(dir, "src", "db", "schema"+ext), readTemplateFile("schema-drizzle.ts")); err != nil {
			return err
		}

	default:
		dbContent := g.renderRawDriver(ext)
		if err := writeFile(filepath.Join(dir, "src", "db", "index"+ext), dbContent); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderPrismaSchema() string {
	provider := "postgresql"
	switch g.cfg.DatabaseType {
	case "mysql", "tidb":
		provider = "mysql"
	case "mongodb":
		provider = "mongodb"
	}
	return fmt.Sprintf(`generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "%s"
  url      = env("DATABASE_URL")
}

model User {
  id        String   @id @default(cuid())
  email     String   @unique
  name      String?
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}
`, provider)
}

func (g *Generator) renderRawDriver(ext string) string {
	switch g.cfg.DatabaseType {
	case "mysql":
		return `import mysql from "mysql2/promise";

const pool = mysql.createPool(process.env.DATABASE_URL!);
export default pool;
`
	case "postgresql":
		return `import { Pool } from "pg";

const pool = new Pool({ connectionString: process.env.DATABASE_URL });
export default pool;
`
	case "mongodb":
		return `import mongoose from "mongoose";

export async function connectDB() {
  if (mongoose.connection.readyState >= 1) return;
  return mongoose.connect(process.env.DATABASE_URL!);
}
`
	case "supabase":
		return `import { createClient } from "@supabase/supabase-js";

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL!;
const supabaseKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!;
export const supabase = createClient(supabaseUrl, supabaseKey);
`
	case "tidb":
		return `import { connect } from "@tidbcloud/serverless";

const conn = connect({ url: process.env.DATABASE_URL! });
export default conn;
`
	}
	return ""
}

func (g *Generator) writeAPIFiles(dir string) error {
	switch g.cfg.APIPattern {
	case "trpc":
		if err := g.writeTRPCFiles(dir); err != nil {
			return err
		}
	case "rest":
		if err := g.writeRESTFiles(dir); err != nil {
			return err
		}
	case "graphql":
		if err := g.writeGraphQLFiles(dir); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) writeTRPCFiles(dir string) error {
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
	return writeFile(filepath.Join(apiDir, "[...trpc]", "route.ts"), routeContent)
}

func (g *Generator) writeRESTFiles(dir string) error {
	apiDir := filepath.Join(dir, "src", "app", "api", "hello")
	if g.cfg.Router == "pages" {
		apiDir = filepath.Join(dir, "src", "pages", "api")
	}
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return err
	}

	if g.cfg.Router == "app" {
		content := `import { NextResponse } from "next/server";

export async function GET() {
  return NextResponse.json({ message: "Hello from VoinzNext!" });
}

export async function POST(req: Request) {
  const body = await req.json();
  return NextResponse.json({ received: body });
}
`
		return writeFile(filepath.Join(apiDir, "route.ts"), content)
	}

	content := `import type { NextApiRequest, NextApiResponse } from "next";

type Data = {
  message: string;
};

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<Data>,
) {
  res.status(200).json({ message: "Hello from VoinzNext!" });
}
`
	return writeFile(filepath.Join(apiDir, "hello.ts"), content)
}

func (g *Generator) writeGraphQLFiles(dir string) error {
	serverDir := filepath.Join(dir, "src", "server")
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return err
	}

	schemaContent := "import { createSchema } from \"graphql-yoga\";\nconst typeDefs = `\n  type Query {\n    hello: String!\n  }\n`;\nconst resolvers = {\n  Query: {\n    hello: () => \"Hello from VoinzNext!\",\n  },\n};\nexport const schema = createSchema({ typeDefs, resolvers });\n"
	if err := writeFile(filepath.Join(serverDir, "schema.ts"), schemaContent); err != nil {
		return err
	}

	routeDir := filepath.Join(dir, "src", "app", "api", "graphql")
	if g.cfg.Router == "pages" {
		routeDir = filepath.Join(dir, "src", "pages", "api")
	}
	if err := os.MkdirAll(routeDir, 0755); err != nil {
		return err
	}

	if g.cfg.Router == "app" {
		routeContent := `import { createYoga } from "graphql-yoga";
import { schema } from "@/server/schema";

const { handleRequest } = createYoga({ schema });

export { handleRequest as GET, handleRequest as POST };
`
		return writeFile(filepath.Join(routeDir, "route.ts"), routeContent)
	}

	routeContent := `import { createYoga } from "graphql-yoga";
import type { NextApiRequest, NextApiResponse } from "next";
import { schema } from "@/server/schema";

export default createYoga<{
  req: NextApiRequest;
  res: NextApiResponse;
}>({
  schema,
  graphqlEndpoint: "/api/graphql",
});
`
	return writeFile(filepath.Join(routeDir, "graphql.ts"), routeContent)
}

// ── Auth boilerplate ──

func (g *Generator) writeAuthFiles(dir string) error {
	switch g.cfg.Auth {
	case "nextauth":
		return g.writeNextAuthFiles(dir)
	case "clerk":
		return g.writeClerkFiles(dir)
	case "lucia":
		return g.writeLuciaFiles(dir)
	}
	return nil
}

func (g *Generator) writeNextAuthFiles(dir string) error {
	libContent := `import NextAuth from "next-auth";
import type { NextAuthConfig } from "next-auth";

export const authConfig: NextAuthConfig = {
  providers: [],
  callbacks: {},
};

export const { handlers, auth, signIn, signOut } = NextAuth(authConfig);
`
	if err := writeFile(filepath.Join(dir, "src", "lib", "auth.ts"), libContent); err != nil {
		return err
	}

	// App Router: src/app/api/auth/[...nextauth]/route.ts
	if g.cfg.Router == "app" {
		routeDir := filepath.Join(dir, "src", "app", "api", "auth", "[...nextauth]")
		if err := os.MkdirAll(routeDir, 0755); err != nil {
			return err
		}
		routeContent := `import { handlers } from "@/lib/auth";

export const { GET, POST } = handlers;
`
		return writeFile(filepath.Join(routeDir, "route.ts"), routeContent)
	}

	// Pages Router: src/pages/api/auth/[...nextauth].ts
	routeDir := filepath.Join(dir, "src", "pages", "api", "auth")
	if err := os.MkdirAll(routeDir, 0755); err != nil {
		return err
	}
	routeContent := `import NextAuth from "next-auth";
import { authConfig } from "@/lib/auth";

export default NextAuth(authConfig);
`
	return writeFile(filepath.Join(routeDir, "[...nextauth].ts"), routeContent)
}

func (g *Generator) writeClerkFiles(dir string) error {
	// App Router: middleware.ts wrapping
	content := `import { clerkMiddleware } from "@clerk/nextjs/server";

export default clerkMiddleware();

export const config = {
  matcher: [
    // Skip Next.js internals and all static files
    "/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)",
    "/(api|trpc)(.*)",
  ],
};
`
	// Middleware goes at project root for both router types
	return writeFile(filepath.Join(dir, "middleware.ts"), content)
}

func (g *Generator) writeLuciaFiles(dir string) error {
	libContent := `import { Lucia } from "lucia";

export const lucia = new Lucia();
`
	if err := writeFile(filepath.Join(dir, "src", "lib", "auth.ts"), libContent); err != nil {
		return err
	}

	// Minimal middleware
	content := `import { lucia } from "@/lib/auth";

export default async function middleware() {
  // Add Lucia session validation here
}

export const config = {
  matcher: "/((?!_next/static|_next/image|favicon.ico).*)",
};
`
	return writeFile(filepath.Join(dir, "middleware.ts"), content)
}

// ── Test files ──

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

// ── Ordered JSON types for deterministic package.json output ──

type sortedMapEntry struct {
	Key   string
	Value string
}

type sortedMap []sortedMapEntry

func (s sortedMap) Len() int           { return len(s) }
func (s sortedMap) Less(i, j int) bool { return s[i].Key < s[j].Key }
func (s sortedMap) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s sortedMap) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte("{}"), nil
	}
	var buf strings.Builder
	buf.WriteByte('{')
	for i, entry := range s {
		if i > 0 {
			buf.WriteByte(',')
		}
		key, _ := json.Marshal(entry.Key)
		val, _ := json.Marshal(entry.Value)
		buf.Write(key)
		buf.WriteByte(':')
		buf.Write(val)
	}
	buf.WriteByte('}')
	return []byte(buf.String()), nil
}

func newSortedMapFromMap(m map[string]string) sortedMap {
	s := make(sortedMap, 0, len(m))
	for k, v := range m {
		s = append(s, sortedMapEntry{Key: k, Value: v})
	}
	sort.Sort(s)
	return s
}

type pnpmConfig struct {
	OnlyBuiltDependencies []string `json:"onlyBuiltDependencies"`
}

type orderedPackageJSON struct {
	Name            string     `json:"name"`
	Version         string     `json:"version"`
	Private         bool       `json:"private"`
	Scripts         sortedMap  `json:"scripts"`
	Dependencies    sortedMap  `json:"dependencies"`
	DevDependencies sortedMap  `json:"devDependencies"`
	Pnpm            *pnpmConfig `json:"pnpm,omitempty"`
}

func writeJSON(path string, data interface{}) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	content = append(content, '\n')
	return writeFile(path, string(content))
}

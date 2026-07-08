package registry

import "github.com/VoinzzZ/VoinzNext/internal/config"

var Questions = []config.Question{
	{
		Key:     "ProjectName",
		Message: "What is your project name?",
		Options: nil,
		Default: "my-next-app",
	},
	{
		Key:     "Router",
		Message: "Which Next.js router would you like to use?",
		Options: []config.Option{
			{ID: "app", Name: "App Router", Description: "Next.js 13+ App Router (recommended)"},
			{ID: "pages", Name: "Pages Router", Description: "Classic Next.js Pages Router"},
		},
		Default: "app",
	},
	{
		Key:     "Language",
		Message: "Which language?",
		Options: []config.Option{
			{ID: "typescript", Name: "TypeScript", Description: "Type-safe JavaScript"},
			{ID: "javascript", Name: "JavaScript", Description: "Plain JavaScript"},
		},
		Default: "typescript",
	},
	{
		Key:     "PackageManager",
		Message: "Which package manager?",
		Options: []config.Option{
			{ID: "pnpm", Name: "pnpm", Description: "Fast, disk space efficient"},
			{ID: "npm", Name: "npm", Description: "Node package manager"},
			{ID: "yarn", Name: "Yarn", Description: "Yarn classic"},
		},
		Default: "pnpm",
	},
	{
		Key:     "CSSFramework",
		Message: "Which CSS framework?",
		Options: []config.Option{
			{ID: "tailwind", Name: "Tailwind CSS", Description: "Utility-first CSS framework"},
			{ID: "css-modules", Name: "CSS Modules", Description: "Scoped CSS by default"},
			{ID: "none", Name: "None", Description: "Plain CSS"},
		},
		Default: "tailwind",
	},
	{
		Key:     "UILibrary",
		Message: "Which UI component library?",
		Options: []config.Option{
			{ID: "shadcn", Name: "shadcn/ui", Description: "Beautifully designed components (Radix + Tailwind)"},
			{ID: "daisyui", Name: "daisyUI", Description: "Tailwind CSS component library"},
			{ID: "none", Name: "None", Description: "No UI library"},
		},
		Default: "shadcn",
	},
	{
		Key:     "DatabaseType",
		Message: "Which database do you want to use?",
		Options: []config.Option{
			{ID: "mysql", Name: "MySQL", Description: "Open-source relational database"},
			{ID: "postgresql", Name: "PostgreSQL", Description: "Advanced open-source relational database"},
			{ID: "mongodb", Name: "MongoDB", Description: "NoSQL document database"},
			{ID: "supabase", Name: "Supabase", Description: "Open-source Firebase alternative (PostgreSQL)"},
			{ID: "tidb", Name: "TiDB", Description: "MySQL-compatible distributed SQL database"},
			{ID: "none", Name: "None", Description: "No database"},
		},
		Default: "none",
	},
	{
		Key:     "DatabaseORM",
		Message: "Which database ORM?",
		Options: []config.Option{
			{ID: "prisma", Name: "Prisma", Description: "Next-gen ORM for Node.js & TypeScript"},
			{ID: "drizzle", Name: "Drizzle ORM", Description: "Lightweight, performant ORM"},
			{ID: "none", Name: "None", Description: "No ORM (use raw driver)"},
		},
		Default: "prisma",
	},
	{
		Key:     "Auth",
		Message: "Which authentication solution?",
		Options: []config.Option{
			{ID: "nextauth", Name: "NextAuth.js", Description: "Authentication for Next.js"},
			{ID: "lucia", Name: "Lucia", Description: "Simple, flexible auth"},
			{ID: "clerk", Name: "Clerk", Description: "Complete user management"},
			{ID: "none", Name: "None", Description: "No authentication"},
		},
		Default: "nextauth",
	},
	{
		Key:     "APIPattern",
		Message: "Which API pattern?",
		Options: []config.Option{
			{ID: "trpc", Name: "tRPC", Description: "End-to-end typesafe APIs"},
			{ID: "rest", Name: "REST API", Description: "Traditional REST endpoints"},
			{ID: "graphql", Name: "GraphQL", Description: "GraphQL with Apollo/Urql"},
			{ID: "none", Name: "None", Description: "No specific API pattern"},
		},
		Default: "trpc",
	},
	{
		Key:     "Testing",
		Message: "Which testing framework?",
		Options: []config.Option{
			{ID: "vitest", Name: "Vitest + Testing Library", Description: "Fast unit & integration tests"},
			{ID: "jest", Name: "Jest + Testing Library", Description: "Classic testing framework"},
			{ID: "playwright", Name: "Playwright", Description: "E2E testing"},
			{ID: "none", Name: "None", Description: "No testing setup"},
		},
		Default: "vitest",
	},
	{
		Key:     "Docker",
		Message: "Include Docker setup?",
		Options: []config.Option{
			{ID: "true", Name: "Yes", Description: "Dockerfile + docker-compose.yml"},
			{ID: "false", Name: "No", Description: "No Docker setup"},
		},
		Default: "true",
	},
	{
		Key:     "ESLintPrettier",
		Message: "Include ESLint + Prettier?",
		Options: []config.Option{
			{ID: "true", Name: "Yes", Description: "ESLint + Prettier configuration"},
			{ID: "false", Name: "No", Description: "No linter/formatter setup"},
		},
		Default: "true",
	},
	{
		Key:     "InitGit",
		Message: "Initialize git repository?",
		Options: []config.Option{
			{ID: "true", Name: "Yes", Description: "git init after project creation"},
			{ID: "false", Name: "No", Description: "No git initialization"},
		},
		Default: "true",
	},
}

func GetDependencies(cfg *config.ProjectConfig) map[string]string {
	deps := map[string]string{
		"next":      "^15.2.0",
		"react":     "^19.0.0",
		"react-dom": "^19.0.0",
	}

	devDeps := map[string]string{}

	if cfg.Language == "typescript" {
		devDeps["typescript"] = "^5.7.0"
		devDeps["@types/react"] = "^19.0.0"
		devDeps["@types/react-dom"] = "^19.0.0"
		devDeps["@types/node"] = "^22.10.0"
	}

	switch cfg.CSSFramework {
	case "tailwind":
		devDeps["tailwindcss"] = "^4.0.0"
		devDeps["postcss"] = "^8.5.0"       // Peer dependency for tailwindcss
		devDeps["autoprefixer"] = "^10.4.0" // Peer dependency for tailwindcss
	}

	switch cfg.UILibrary {
	case "shadcn":
		// shadcn/ui peer dependencies
		devDeps["lucide-react"] = "^0.460.0"
		devDeps["class-variance-authority"] = "^0.7.1"
		devDeps["clsx"] = "^2.1.1"
		devDeps["tailwind-merge"] = "^2.6.0"
		devDeps["tailwindcss-animate"] = "^1.0.7"
	case "daisyui":
		devDeps["daisyui"] = "^5.0.0"
	}

	switch cfg.DatabaseORM {
	case "prisma":
		devDeps["prisma"] = "^7.8.0"
		deps["@prisma/client"] = "^7.8.0"
		switch cfg.DatabaseType {
		case "tidb":
			deps["@tidbcloud/prisma-adapter"] = "^6.17.0"
			deps["@tidbcloud/serverless"] = "^0.3.0"
		case "mysql":
			deps["mysql2"] = "^3.22.0"
		}
	case "drizzle":
		deps["drizzle-orm"] = "^0.45.0"
		devDeps["drizzle-kit"] = "^0.32.0"
		switch cfg.DatabaseType {
		case "mysql", "tidb":
			deps["mysql2"] = "^3.22.0"
		default:
			deps["pg"] = "^8.22.0"
		}
	}

	if cfg.DatabaseType != "none" && cfg.DatabaseORM == "none" {
		switch cfg.DatabaseType {
		case "mysql":
			deps["mysql2"] = "^3.22.0"
		case "postgresql", "supabase":
			deps["pg"] = "^8.22.0"
		case "mongodb":
			deps["mongoose"] = "^9.7.0"
		case "tidb":
			deps["@tidbcloud/serverless"] = "^0.3.0"
		}
	}

	if cfg.DatabaseType == "supabase" {
		deps["@supabase/supabase-js"] = "^2.108.0"
	}

	switch cfg.Auth {
	case "nextauth":
		deps["next-auth"] = "^5.0.0"
	case "lucia":
		deps["lucia"] = "^3.2.0"
		devDeps["@lucia-auth/adapter-prisma"] = "^4.0.0"
	case "clerk":
		deps["@clerk/nextjs"] = "^5.0.0"
	}

	switch cfg.APIPattern {
	case "trpc":
		deps["@trpc/server"] = "^11.0.0"
		deps["@trpc/client"] = "^11.0.0"
		deps["@trpc/react-query"] = "^11.0.0"
		deps["@tanstack/react-query"] = "^5.60.0" // Peer dependency for @trpc/react-query
		devDeps["@trpc/next"] = "^11.0.0"
		devDeps["zod"] = "^3.23.0"
	case "graphql":
		deps["graphql"] = "^16.9.0"
		deps["graphql-yoga"] = "^5.6.0"
	}

	switch cfg.Testing {
	case "vitest":
		devDeps["vitest"] = "^3.0.0"
		devDeps["@testing-library/react"] = "^16.0.0"
		devDeps["@testing-library/jest-dom"] = "^6.5.0"
		devDeps["jsdom"] = "^25.0.0"
	case "jest":
		devDeps["jest"] = "^29.7.0"
		devDeps["@testing-library/react"] = "^16.0.0"
		devDeps["@testing-library/jest-dom"] = "^6.5.0"
		devDeps["jest-environment-jsdom"] = "^29.7.0"
		devDeps["@types/jest"] = "^29.5.0"
		devDeps["ts-jest"] = "^29.1.0"
		devDeps["babel-jest"] = "^29.7.0"
	case "playwright":
		devDeps["@playwright/test"] = "^1.47.0"
	}

	if cfg.ESLintPrettier {
		devDeps["eslint"] = "^9.0.0"
		devDeps["eslint-config-next"] = "^15.2.0"
		devDeps["prettier"] = "^3.5.0"
		devDeps["prettier-plugin-tailwindcss"] = "^0.6.0"
		if cfg.Language == "typescript" {
			devDeps["@typescript-eslint/eslint-plugin"] = "^7.0.0"
			devDeps["@typescript-eslint/parser"] = "^7.0.0"
		}
	}

	result := make(map[string]string)
	for k, v := range deps {
		result[k] = v
	}
	for k, v := range devDeps {
		result["dev:"+k] = v
	}

	return result
}

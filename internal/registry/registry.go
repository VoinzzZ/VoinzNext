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
		Key:     "DatabaseORM",
		Message: "Which database ORM?",
		Options: []config.Option{
			{ID: "prisma", Name: "Prisma", Description: "Next-gen ORM for Node.js & TypeScript"},
			{ID: "drizzle", Name: "Drizzle ORM", Description: "Lightweight, performant ORM"},
			{ID: "none", Name: "None", Description: "No database ORM"},
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
		"next": "^14.2.0",
		"react": "^18.3.0",
		"react-dom": "^18.3.0",
	}

	devDeps := map[string]string{}

	if cfg.Language == "typescript" {
		devDeps["typescript"] = "^5.6.0"
		devDeps["@types/react"] = "^18.3.0"
		devDeps["@types/react-dom"] = "^18.3.0"
		devDeps["@types/node"] = "^22.0.0"
	}

	switch cfg.CSSFramework {
	case "tailwind":
		devDeps["tailwindcss"] = "^3.4.0"
		devDeps["postcss"] = "^8.4.0"
		devDeps["autoprefixer"] = "^10.4.0"
	}

	switch cfg.UILibrary {
	case "shadcn":
		devDeps["lucide-react"] = "^0.400.0"
		devDeps["class-variance-authority"] = "^0.7.0"
		devDeps["clsx"] = "^2.1.0"
		devDeps["tailwind-merge"] = "^2.3.0"
		devDeps["tailwindcss-animate"] = "^1.0.7"
	case "daisyui":
		devDeps["daisyui"] = "^4.12.0"
	}

	switch cfg.DatabaseORM {
	case "prisma":
		devDeps["prisma"] = "^5.14.0"
		deps["@prisma/client"] = "^5.14.0"
	case "drizzle":
		deps["drizzle-orm"] = "^0.33.0"
		devDeps["drizzle-kit"] = "^0.24.0"
		if cfg.DatabaseORM == "drizzle" && cfg.APIPattern == "trpc" {
			deps["@planetscale/database"] = "^1.19.0"
		}
	}

	switch cfg.Auth {
	case "nextauth":
		deps["next-auth"] = "^4.24.0"
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
		deps["@tanstack/react-query"] = "^5.60.0"
		devDeps["@trpc/next"] = "^11.0.0"
		devDeps["zod"] = "^3.23.0"
	case "graphql":
		deps["graphql"] = "^16.9.0"
		deps["graphql-yoga"] = "^5.6.0"
	}

	switch cfg.Testing {
	case "vitest":
		devDeps["vitest"] = "^2.0.0"
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
		devDeps["eslint"] = "^8.57.0"
		devDeps["eslint-config-next"] = "^14.2.0"
		devDeps["prettier"] = "^3.3.0"
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

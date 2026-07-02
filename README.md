# VoinzNext

Interactive CLI tool for scaffolding Next.js projects with your preferred tech stack.

## Installation

```bash
go install github.com/VoinzzZ/VoinzNext/cmd/voinznest@latest
```

## Usage

```bash
voinznest init     # Start interactive survey & generate project
voinznest list     # Show available tech stacks
voinznest add      # Add feature to existing project (coming soon)
voinznest update   # Update to latest version
voinznest version  # Show version info
```

## Tech Stack Options

- **Router:** App Router / Pages Router
- **Language:** TypeScript / JavaScript
- **CSS:** Tailwind CSS / CSS Modules / None
- **UI Library:** shadcn/ui / daisyUI / None
- **Database:** Prisma / Drizzle / None
- **Auth:** NextAuth.js / Lucia / Clerk / None
- **API:** tRPC / REST / GraphQL / None
- **Testing:** Vitest / Jest / Playwright / None
- **Docker:** Yes / No
- **ESLint + Prettier:** Yes / No
- **Git init:** Yes / No

## Development

```bash
git clone https://github.com/VoinzzZ/VoinzNext.git
cd VoinzNext
go build -o bin/voinznest.exe ./cmd/voinznest/
go test ./...
```

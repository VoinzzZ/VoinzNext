# VoinzNext

Interactive CLI tool for scaffolding Next.js projects with your preferred tech stack.

## Installation

### Via npm / npx (recommended)

```bash
# Run directly without installing
npx voinznext init

# Or install globally
npm install -g voinznext
voinznext init
```

### Via Go (alternative)

```bash
go install github.com/VoinzzZ/VoinzNext/cmd/voinznext@latest
```

> The binary auto-updates itself via `voinznext update` — no need to reinstall.

## Usage

```bash
voinznext init     # Start interactive survey & generate project
voinznext list     # Show available tech stacks
voinznext add      # Add feature to existing project (coming soon)
voinznext update   # Update to latest version
voinznext version  # Show version info
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
go build -o bin/voinznext ./cmd/voinznext/
go test ./...
```

# Changelog

All notable changes to this project will be documented in this file.

## [0.5.0] - 2026-07-10

### Added
- **Interactive Init CLI**:
  - Support passing project name as argument to `voinznext init <project-name>` to skip the name prompt.
  - Confirmation prompt before overwriting non-empty project directories.
- **Update & Download Progress**:
  - Version comparison logic to accurately determine if update is available.
  - Interactive download progress bar for binary downloads in the CLI.
- **Generator Tests**:
  - Added test suite for JavaScript mode to verify generated files do not contain TypeScript syntax.

### Changed
- **Templates & Code Generation**:
  - Enhanced tRPC generation for both App Router and Pages Router.
  - Improved deterministic sorting and generation of dependencies in `package.json`.
  - Refined project directory handling and comments structure.
- **Self-Update Script**:
  - Improved self-replacement `.bat` script for reliable execution on Windows.

### Fixed

**Note:** v0.5.0 release was delayed to address critical compatibility and configuration issues discovered during pre-release audit. The following 6 fixes were applied:

1. **NextAuth v5 → v4 Compatibility**  
   Downgraded `next-auth` from `^5.0.0` to `^4.24.14` for stable Next.js 14 compatibility. Updated auth configuration from `NextAuthConfig` to `NextAuthOptions` pattern and adjusted route handlers accordingly.

2. **Tailwind v4 PostCSS Plugin**  
   Fixed Tailwind CSS v4 configuration by switching from `tailwindcss` to `@tailwindcss/postcss` plugin in `postcss.config.js` and added `@tailwindcss/postcss` to devDependencies.

3. **pnpm Build Scripts Configuration**  
   Moved build script approvals from `package.json` (`onlyBuiltDependencies`) to proper configuration files: `pnpm-workspace.yaml` (with `allowBuilds` for `esbuild`, `sharp`, `unrs-resolver`) and `.npmrc` (with `enable-pre-post-scripts=true` and `shell-emulator=true`).

4. **JavaScript Mode jsconfig.json**  
   Added `jsconfig.json` generation for JavaScript projects (previously only `tsconfig.json` was generated for TypeScript projects, leaving JS projects without proper IDE support).

5. **ESLint Configuration for JS Projects**  
   Simplified `eslint.config.js` template by removing TypeScript-specific parser and rules that caused errors in JavaScript-only projects.

6. **Version Comparison String Bug**  
   Fixed version comparison logic that used string comparison instead of proper semantic version parsing (e.g., "0.9.0" > "0.10.0" as strings), causing incorrect update detection.

---

*For older changes, refer to Git commits.*

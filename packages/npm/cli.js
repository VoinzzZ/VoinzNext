#!/usr/bin/env node
const { spawn } = require("child_process");
const { existsSync } = require("fs");
const { join } = require("path");
const { platform } = require("process");

const binaryName = platform === "win32" ? "voinznext.exe" : "voinznext";
const binaryPath = join(__dirname, "bin", binaryName);

if (!existsSync(binaryPath)) {
  console.error("  ✘ Binary not found. Run `npm install` or `npx voinznext` to download it.");
  process.exit(1);
}

const child = spawn(binaryPath, process.argv.slice(2), { stdio: "inherit" });

child.on("exit", (code) => process.exit(code));
child.on("error", (err) => {
  console.error("  ✘ Failed to run voinznext:", err.message);
  process.exit(1);
});

#!/usr/bin/env node
module.exports = { downloadLatest };

const { spawnSync } = require("child_process");
const { createWriteStream, existsSync, mkdirSync, readFileSync, writeFileSync, chmodSync, renameSync, unlinkSync } = require("fs");
const https = require("https");
const { join } = require("path");
const { platform, arch } = require("process");

const REPO = "VoinzzZ/VoinzNext";
const APP = "voinznext";
const CACHE_TTL = 3600000;

const PLATFORM_MAP = { win32: "windows", darwin: "darwin", linux: "linux" };
const ARCH_MAP = { x64: "amd64", arm64: "arm64" };

function getBinaryName() {
  const p = PLATFORM_MAP[platform];
  const a = ARCH_MAP[arch];
  if (!p || !a) return null;
  const name = `${APP}-${p}-${a}`;
  return platform === "win32" ? `${name}.exe` : name;
}

function getBinaryPath() {
  return join(__dirname, "bin", platform === "win32" ? "voinznext.exe" : "voinznext");
}

function getVersionPath() {
  return join(__dirname, "bin", ".version");
}

function fetchJson(url) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { "Accept": "application/vnd.github.v3+json", "User-Agent": "voinznext-updater" } }, (res) => {
      let data = "";
      res.on("data", (chunk) => (data += chunk));
      res.on("end", () => {
        try { resolve(JSON.parse(data)); }
        catch { reject(new Error("Failed to parse response")); }
      });
    }).on("error", reject);
  });
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { "User-Agent": "voinznext-updater" } }, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        download(res.headers.location, dest).then(resolve).catch(reject);
        return;
      }
      if (res.statusCode !== 200) {
        reject(new Error(`HTTP ${res.statusCode}`));
        return;
      }
      const file = createWriteStream(dest);
      res.pipe(file);
      file.on("finish", () => file.close(resolve));
      file.on("error", reject);
    }).on("error", reject);
  });
}

function readVersionCache() {
  try {
    const data = readFileSync(getVersionPath(), "utf8");
    return JSON.parse(data);
  } catch { return null; }
}

function writeVersionCache(tag) {
  try {
    mkdirSync(join(__dirname, "bin"), { recursive: true });
    writeFileSync(getVersionPath(), JSON.stringify({ tag, checked: Date.now() }));
  } catch {}
}

async function downloadLatest() {
  const binaryPath = getBinaryPath();
  const binaryName = getBinaryName();
  if (!binaryName) return;

  const release = await fetchJson(`https://api.github.com/repos/${REPO}/releases/latest`);
  const version = release.tag_name;
  const url = `https://github.com/${REPO}/releases/download/${version}/${binaryName}`;
  const tmpPath = binaryPath + ".new";

  mkdirSync(join(__dirname, "bin"), { recursive: true });
  await download(url, tmpPath);
  if (platform !== "win32") chmodSync(tmpPath, 0o755);

  try {
    if (platform === "win32") {
      renameSync(tmpPath, binaryPath);
    } else {
      if (existsSync(binaryPath)) unlinkSync(binaryPath);
      renameSync(tmpPath, binaryPath);
    }
  } catch (e) {}
  writeVersionCache(version);
}

async function ensureBinary() {
  const binaryPath = getBinaryPath();
  const cache = readVersionCache();
  const binaryExists = existsSync(binaryPath);

  if (!cache || !binaryExists) {
    console.log("  ● Downloading voinznext binary...");
    await downloadLatest();
    return;
  }

  if (Date.now() - cache.checked > CACHE_TTL) {
    try {
      const release = await fetchJson(`https://api.github.com/repos/${REPO}/releases/latest`);
      if (release.tag_name !== cache.tag) {
        console.log("  ● Updating voinznext binary...");
        await downloadLatest();
        return;
      }
      writeVersionCache(cache.tag);
    } catch {}
  }
}

async function main() {
  if (!process.env.VOINZNEXT_INSTALL) {
    await ensureBinary().catch(() => {});
  }

  if (process.env.VOINZNEXT_INSTALL) {
    process.exit(0);
  }

  const binaryPath = getBinaryPath();
  if (!existsSync(binaryPath)) {
    console.error("  ✘ Binary not found. Run `npm install -g voinznext@latest --force`.");
    process.exit(1);
  }

  const result = spawnSync(binaryPath, process.argv.slice(2), { stdio: "inherit" });
  process.exit(result.status);
}

if (require.main === module) main();

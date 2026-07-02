#!/usr/bin/env node
const { createWriteStream, existsSync, mkdirSync, chmodSync } = require("fs");
const https = require("https");
const { join } = require("path");
const { platform, arch } = require("process");

const REPO = "VoinzzZ/VoinzNext";
const APP = "voinznext";

const PLATFORM_MAP = { win32: "windows", darwin: "darwin", linux: "linux" };
const ARCH_MAP = { x64: "amd64", arm64: "arm64" };

function getPlatform() {
  const p = PLATFORM_MAP[platform];
  if (!p) throw new Error(`Unsupported platform: ${platform}`);
  return p;
}

function getArch() {
  const a = ARCH_MAP[arch];
  if (!a) throw new Error(`Unsupported architecture: ${arch}`);
  return a;
}

function getBinaryName() {
  const p = getPlatform();
  const a = getArch();
  const name = `${APP}-${p}-${a}`;
  return p === "windows" ? `${name}.exe` : name;
}

function fetchJson(url) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { "Accept": "application/vnd.github.v3+json", "User-Agent": "voinznext-installer" } }, (res) => {
      let data = "";
      res.on("data", (chunk) => (data += chunk));
      res.on("end", () => {
        try { resolve(JSON.parse(data)); }
        catch { reject(new Error(`Failed to parse response: ${data}`)); }
      });
    }).on("error", reject);
  });
}

function getLatestVersion() {
  return fetchJson(`https://api.github.com/repos/${REPO}/releases/latest`).then((json) => json.tag_name);
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { "User-Agent": "voinznext-installer" } }, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        download(res.headers.location, dest).then(resolve).catch(reject);
        return;
      }
      if (res.statusCode !== 200) {
        reject(new Error(`Download failed with status ${res.statusCode}`));
        return;
      }
      const file = createWriteStream(dest);
      res.pipe(file);
      file.on("finish", () => file.close(resolve));
      file.on("error", reject);
    }).on("error", reject);
  });
}

async function install() {
  const binDir = join(__dirname, "bin");
  if (!existsSync(binDir)) mkdirSync(binDir, { recursive: true });

  const binaryName = getBinaryName();
  const targetPath = join(binDir, "voinznext.exe");

  console.log(`  ● Downloading ${APP} binary for ${platform}-${arch}...`);

  const version = await getLatestVersion();
  const downloadUrl = `https://github.com/${REPO}/releases/download/${version}/${binaryName}`;

  await download(downloadUrl, targetPath);

  if (platform !== "win32") chmodSync(targetPath, 0o755);
  console.log(`  ✔ Installed to ${targetPath}`);
}

install().catch((err) => {
  console.error("  ✘ Installation failed:", err.message);
  process.exit(1);
});

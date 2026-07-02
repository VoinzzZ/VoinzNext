#!/usr/bin/env node
const { createWriteStream, existsSync, mkdirSync, chmodSync } = require("fs");
const { get } = require("https");
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

function getLatestVersion() {
  return new Promise((resolve, reject) => {
    const url = `https://api.github.com/repos/${REPO}/releases/latest`;
    get(url, { headers: { "Accept": "application/vnd.github.v3+json", "User-Agent": "voinznext-installer" } }, (res) => {
      let data = "";
      res.on("data", (chunk) => (data += chunk));
      res.on("end", () => {
        try {
          const json = JSON.parse(data);
          resolve(json.tag_name);
        } catch {
          reject(new Error(`Failed to parse latest release response: ${data}`));
        }
      });
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

  await new Promise((resolve, reject) => {
    get(downloadUrl, { headers: { "User-Agent": "voinznext-installer" } }, (res) => {
      if (res.statusCode !== 200) {
        reject(new Error(`Download failed with status ${res.statusCode}`));
        return;
      }
      const file = createWriteStream(targetPath);
      res.pipe(file);
      file.on("finish", () => file.close(resolve));
      file.on("error", reject);
    }).on("error", reject);
  });

  if (platform !== "win32") chmodSync(targetPath, 0o755);
  console.log(`  ✔ Installed to ${targetPath}`);
}

install().catch((err) => {
  console.error("  ✘ Installation failed:", err.message);
  process.exit(1);
});

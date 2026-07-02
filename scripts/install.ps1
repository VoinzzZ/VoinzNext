#!/usr/bin/env pwsh
<#
.SYNOPSIS
  Install VoinzNext CLI on Windows
.DESCRIPTION
  Downloads the latest VoinzNext binary from GitHub Releases and adds it to PATH.
.EXAMPLE
  powershell -c "irm https://raw.githubusercontent.com/VoinzzZ/VoinzNext/main/scripts/install.ps1 | iex"
#>

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$Repo   = "VoinzzZ/VoinzNext"
$App    = "voinznext"
$BinDir = "$HOME\.$App\bin"

# ── Colors ──
$Green  = "Green"
$Yellow = "Yellow"
$Cyan   = "Cyan"
$Red    = "Red"

function Write-Color($Color, $Text) {
  Write-Host $Text -ForegroundColor $Color
}

# ── Detect latest version ──
Write-Color $Cyan "● Fetching latest version..."

try {
  $ApiUrl  = "https://api.github.com/repos/$Repo/releases/latest"
  $Release = Invoke-RestMethod -Uri $ApiUrl -Headers @{ "Accept" = "application/vnd.github.v3+json" }
  $Version = $Release.tag_name
} catch {
  Write-Color $Red "✘ Failed to fetch latest release from GitHub."
  exit 1
}

Write-Color $Green "✔ Latest version: $Version"

# ── Detect architecture ──
$Arch = "amd64"
if ([Environment]::Is64BitOperatingSystem -eq $false) {
  Write-Color $Red "✘ 32-bit systems are not supported."
  exit 1
}

# ── Download binary ──
$FileName  = "$App-windows-$Arch.exe"
$DownloadUrl = "https://github.com/$Repo/releases/download/$Version/$FileName"
$TargetPath   = "$BinDir\$FileName"
$TempFile     = "$env:TEMP\$FileName"

Write-Color $Cyan "● Downloading $App $Version for Windows $Arch..."
try {
  Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile -UseBasicParsing
} catch {
  Write-Color $Red "✘ Download failed: $($_.Exception.Message)"
  exit 1
}

# ── Install ──
New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
Move-Item -Path $TempFile -Destination $TargetPath -Force

# Remove old binary if exists
$OldExe = "$BinDir\$App.exe"
if (Test-Path $OldExe) {
  Remove-Item $OldExe -Force
}

Write-Color $Green "✔ Binary installed to: $TargetPath"

# ── Add to PATH ──
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$BinDir*") {
  $NewPath = "$UserPath;$BinDir"
  [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
  $env:PATH = "$env:PATH;$BinDir"
  Write-Color $Green "✔ Added $BinDir to PATH (user-level)"
} else {
  Write-Color $Yellow "⚠ $BinDir already in PATH"
}

# ── Verify ──
try {
  $VersionOutput = & "$TargetPath" version
  Write-Color $Green "✔ Installation verified!"
  Write-Host ""
  Write-Color $Cyan "╭──────────────────────────────────────────╮"
  Write-Color $Cyan "│         VoinzNext installed!            │"
  Write-Color $Cyan "├──────────────────────────────────────────┤"
  Write-Color $Green "│  Try: voinznext init                    │"
  Write-Color $Cyan "╰──────────────────────────────────────────╯"
  Write-Host ""
  $VersionOutput
} catch {
  Write-Color $Red "✘ Verification failed."
  exit 1
}

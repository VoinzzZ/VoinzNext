#!/bin/sh
set -e

REPO="VoinzzZ/VoinzNext"
APP="voinznext"
BIN_DIR="$HOME/.$APP/bin"

# ── Colors ──
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

info()  { printf "  ${CYAN}●${NC} %s\n" "$1"; }
ok()    { printf "  ${GREEN}✔${NC} %s\n" "$1"; }
warn()  { printf "  ${YELLOW}⚠${NC} %s\n" "$1"; }
err()   { printf "  ${RED}✘${NC} %s\n" "$1"; exit 1; }

# ── Detect latest version ──
info "Fetching latest version..."
if ! command -v curl >/dev/null 2>&1; then
  err "curl is required. Install curl and try again."
fi

API_URL="https://api.github.com/repos/$REPO/releases/latest"
VERSION=$(curl -sL "$API_URL" | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$VERSION" ]; then
  err "Failed to fetch latest release from GitHub."
fi

ok "Latest version: $VERSION"

# ── Detect OS & arch ──
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux)   OS_ALIAS="linux"   ;;
  darwin)  OS_ALIAS="darwin"  ;;
  *)       err "Unsupported OS: $OS" ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH_ALIAS="amd64" ;;
  aarch64|arm64) ARCH_ALIAS="arm64" ;;
  *)            err "Unsupported architecture: $ARCH" ;;
esac

FILENAME="$APP-$OS_ALIAS-$ARCH_ALIAS"
if [ "$OS_ALIAS" = "windows" ]; then
  FILENAME="$FILENAME.exe"
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

# ── Download binary ──
info "Downloading $APP $VERSION for $OS ($ARCH)..."

TMP_FILE=$(mktemp)
trap 'rm -f "$TMP_FILE"' EXIT

if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"; then
  err "Download failed."
fi

# ── Install ──
mkdir -p "$BIN_DIR"
TARGET_PATH="$BIN_DIR/$APP"
mv "$TMP_FILE" "$TARGET_PATH"
chmod +x "$TARGET_PATH"

ok "Binary installed to: $TARGET_PATH"

# ── Add to PATH ──
SHELL_PROFILE=""
case "$SHELL" in
  */zsh) SHELL_PROFILE="${HOME}/.zshrc" ;;
  */bash)
    if [ -f "${HOME}/.bash_profile" ]; then
      SHELL_PROFILE="${HOME}/.bash_profile"
    else
      SHELL_PROFILE="${HOME}/.bashrc"
    fi
    ;;
  */fish) SHELL_PROFILE="${HOME}/.config/fish/config.fish" ;;
esac

if [ -n "$SHELL_PROFILE" ]; then
  if ! grep -q "$BIN_DIR" "$SHELL_PROFILE" 2>/dev/null; then
    printf '\nexport PATH="%s:$PATH"\n' "$BIN_DIR" >> "$SHELL_PROFILE"
    ok "Added $BIN_DIR to PATH in $SHELL_PROFILE"
  else
    warn "$BIN_DIR already in PATH"
  fi
else
  warn "Unknown shell. Add $BIN_DIR to your PATH manually."
fi

export PATH="$BIN_DIR:$PATH"

# ── Verify ──
if "$APP" version >/dev/null 2>&1; then
  ok "Installation verified!"
  printf "\n"
  printf "  ${CYAN}╭──────────────────────────────────────────╮${NC}\n"
  printf "  ${CYAN}│         VoinzNext installed!            │${NC}\n"
  printf "  ${CYAN}├──────────────────────────────────────────┤${NC}\n"
  printf "  ${GREEN}│  Try: voinznext init                    │${NC}\n"
  printf "  ${CYAN}╰──────────────────────────────────────────╯${NC}\n"
  printf "\n"
  "$APP" version
else
  err "Verification failed."
fi

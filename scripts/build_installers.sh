#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_DIR="./cmd/openpdfreader"

cd "$ROOT_DIR"

if ! command -v fyne-cross >/dev/null 2>&1; then
  echo "fyne-cross is not installed."
  echo "Install it with: go install github.com/fyne-io/fyne-cross@latest"
  exit 1
fi

echo "Building Linux installer artifacts..."
fyne-cross linux -arch=amd64 "$APP_DIR"

echo "Building Windows installer artifacts..."
fyne-cross windows -arch=amd64 "$APP_DIR"

echo ""
echo "Installer artifacts generated under:"
echo "  fyne-cross/dist/"

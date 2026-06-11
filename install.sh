#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "Go is required to install AI agent configuration." >&2
  echo "Install Go, then rerun ./install.sh." >&2
  exit 1
fi

(
  cd "$ROOT_DIR/hooks-go"
  go run ./cmd/agent-config --repo-root "$ROOT_DIR" install
)

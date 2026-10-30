#!/usr/bin/env bash
set -euo pipefail

# Runs every repository check; CI runs it on Linux and macOS.

cd "$(dirname "${BASH_SOURCE[0]}")/.."

if ! cmp -s AGENTS.md CLAUDE.md; then
  echo "AGENTS.md and CLAUDE.md differ; edit one and copy it over the other" >&2
  exit 1
fi

go -C hooks-go run ./cmd/agent-config validate-skills ../skills
go -C hooks-go vet ./...
unformatted=$(gofmt -l hooks-go)
if [ -n "$unformatted" ]; then
  echo "gofmt needed:" >&2
  echo "$unformatted" >&2
  exit 1
fi
go -C hooks-go test -count=1 ./...
bash scripts/check_go_targets.sh
for script in skills/*/scripts/*.sh; do
  sh -n "$script"
done
git diff --check

#!/usr/bin/env bash
set -euo pipefail

# Builds hooks-go for every supported GOOS: linux, darwin, and windows, the
# platforms named in the build tags under internal/securepath.
#
# Building ./... rather than the library packages alone is deliberate: it
# includes the cmd/... main packages, so the link step runs and resolves the
# //go:linkname syscall.openat bridge used on darwin. A failure there is
# invisible to `go test` on the host platform.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "Go is required to check build targets." >&2
  exit 1
fi

TARGETS=(
  linux/amd64
  linux/arm64
  darwin/amd64
  darwin/arm64
  windows/amd64
)

cd "$ROOT_DIR/hooks-go"

failed=()
for target in "${TARGETS[@]}"; do
  goos="${target%/*}"
  goarch="${target#*/}"
  printf '%-16s ' "$target"
  if output="$(GOOS="$goos" GOARCH="$goarch" go build ./... 2>&1)"; then
    echo "ok"
  else
    echo "FAIL"
    failed+=("$target")
    printf '%s\n' "$output" | sed 's/^/    /'
  fi
done

if [ ${#failed[@]} -gt 0 ]; then
  echo "" >&2
  echo "Failed targets: ${failed[*]}" >&2
  exit 1
fi

echo ""
echo "All ${#TARGETS[@]} targets built."

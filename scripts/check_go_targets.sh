#!/usr/bin/env bash
set -euo pipefail

# Builds hooks-go for every GOOS named in the per-platform build tags under
# internal/securepath.
#
# Building ./... rather than the library packages alone is deliberate: it
# includes the cmd/... main packages, so the link step runs and resolves the
# //go:linkname syscall.openat bridges used on darwin and openbsd. A failure
# there is invisible to `go test` on the host platform.
#
# ios is omitted because it requires external cgo linking and an iOS toolchain.
# Its openat bridge is the darwin file (//go:build darwin || ios), which the
# darwin targets below already compile.

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
  android/arm64
  freebsd/amd64
  netbsd/amd64
  dragonfly/amd64
  openbsd/amd64
  aix/ppc64
  solaris/amd64
  illumos/amd64
  windows/amd64
  plan9/amd64
  js/wasm
  wasip1/wasm
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

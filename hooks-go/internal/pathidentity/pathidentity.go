// Package pathidentity resolves directory paths to stable filesystem identities.
package pathidentity

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveDirectory returns the canonical path and file information for path.
// It resolves relative paths and symlinks, then verifies that the resolved
// path exists and names a directory.
func ResolveDirectory(path string) (string, os.FileInfo, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, fmt.Errorf("resolve absolute path %q: %w", path, err)
	}
	canonical, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", nil, fmt.Errorf("resolve symlinks %q: %w", path, err)
	}
	info, err := os.Stat(canonical)
	if err != nil {
		return "", nil, fmt.Errorf("stat directory %q: %w", path, err)
	}
	if !info.IsDir() {
		return "", nil, fmt.Errorf("path is not a directory: %s", path)
	}
	return canonical, info, nil
}

// NormalizeSystemAlias rewrites the fixed macOS compatibility symlinks /var and
// /tmp to their /private targets so a descriptor-relative no-follow traversal
// can walk them. It is the only symlink resolution such a traversal may
// perform: nothing else is rewritten, so user-controlled symlink components
// stay subject to O_NOFOLLOW. On hosts where the alias does not resolve to
// /private/<alias> this returns path unchanged.
func NormalizeSystemAlias(path string) string {
	if !filepath.IsAbs(path) {
		return path
	}
	for _, alias := range []string{"var", "tmp"} {
		prefix := string(filepath.Separator) + alias
		if path != prefix && !strings.HasPrefix(path, prefix+string(filepath.Separator)) {
			continue
		}
		resolved, err := filepath.EvalSymlinks(prefix)
		if err != nil || resolved != filepath.Join(string(filepath.Separator), "private", alias) {
			continue
		}
		return resolved + strings.TrimPrefix(path, prefix)
	}
	return path
}

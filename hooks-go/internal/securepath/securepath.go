// Package securepath opens directories and files without following symlinks.
//
// On targets whose kernel API is reachable from Go it walks a path one
// component at a time from a held directory descriptor using openat(2) with
// O_DIRECTORY|O_NOFOLLOW, so the kernel is never asked to resolve a complete
// caller-supplied path and an ancestor symlink or a concurrent directory
// rename cannot redirect traversal.
//
// On targets without that API the package validates every pathname component
// with Lstat and then opens by pathname. That is best-effort only: it rejects
// symlinks that exist at check time but leaves a same-user time-of-check to
// time-of-use window that the descriptor-relative implementation does not
// have. Callers that need the stronger guarantee must state the platform
// distinction rather than assume it holds everywhere.
//
// Callers supply their own trust anchor. This package canonicalizes nothing
// except the fixed macOS /var and /tmp compatibility aliases, which it
// delegates to pathidentity.NormalizeSystemAlias.
package securepath

# Git And Workspace Gates

The harness coordinates git/workspace awareness but does not own git workflow. Use `conventional-git-flow` for branch, commit, push, PR, and dependency/version closure rules.

## Workspace Awareness

Before side effects or handoff-sensitive work, understand:

- Current branch.
- Dirty worktree state.
- Staged versus unstaged changes.
- Untracked files.
- User-owned changes that must not be overwritten.

## Checkpoint Questions

Ask whether the active task needs:

- A clean starting status.
- A branch or worktree.
- A phase checkpoint.
- A pre-commit diff sanity gate.
- Dependency/version metadata closure.
- A rollback or recovery note.

## Boundaries

- Do not run destructive git commands without explicit approval.
- Do not stage broad patterns by default.
- Do not commit only logical source files while omitting required manifests, lockfiles, or version metadata.
- Do not turn a workspace gate into a PR workflow unless the user asked for commit, push, or PR work.

## Related Owners

- `conventional-git-flow`: commit and PR mechanics.
- `code-review`: bounded current-diff sanity.
- `policy-testing`: verification evidence required before commit readiness.
- A runtime external-lookup role, when available: compatibility lookup for major or incompatible
  dependency updates.

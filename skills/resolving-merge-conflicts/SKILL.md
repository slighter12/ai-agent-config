---
name: resolving-merge-conflicts
description: Resolve all conflicts in an in-progress merge or rebase and continue it to completion. Use when the user explicitly invokes this skill for an active conflict. Avoid when there is no merge or rebase in progress, or when only a read-only safety review is requested.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Resolving Merge Conflicts

Explicit invocation authorizes resolving the current merge or rebase, staging resolved paths, and running the required continue or commit loop.

Inspect both sides, the base, repository intent, and nearby tests before choosing a resolution. Preserve compatible changes from both sides where possible. Validate the resolved result narrowly.

Never abort the merge or rebase. Never discard unrelated work. Stop and ask only when the intended behavior cannot be inferred safely.

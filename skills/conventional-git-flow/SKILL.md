---
name: conventional-git-flow
description: Prepare and execute a conventional branch, commit, push, or pull-request workflow. Use when the user explicitly asks for one of those git actions or its message text. Avoid when resolving merge conflicts, reviewing code, or writing release notes is the primary task.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Conventional Git Flow

Inspect the current branch, diff, status, and repository conventions. Propose the smallest conventional branch name, commit message, or PR text that describes the actual change.

Use `<type>: <description>` for commit headers by default, omitting scope unless repository instructions or the user require it. Include the causal closure of the change: required manifests, lockfiles, and version metadata belong in the same scope. For major or incompatible dependency changes, check current official compatibility information and report an unknown compatibility risk when it cannot be confirmed.

Only perform commit, push, PR creation, merge, or rebase when the user explicitly authorizes that action. Never include unrelated work, rewrite history, or resolve an active conflict under this skill.

Report the exact git evidence and any unpushed or uncommitted state.

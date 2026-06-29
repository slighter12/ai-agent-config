# Release Readiness

Use this reference for `release-readiness` mode. The review is read-only and should produce only confirmed blockers and plausible confirmation items tied to release risk.

## Scope

Select the reviewed range before judging risk:

- Use an explicit PR, branch, commit, or `base..head` range when the user provides one.
- If no scope is provided, prefer the previous usable release tag to `HEAD`.
- If no usable release tag exists, review the latest five commits and state that this is a fallback.
- Do not mix dirty worktree changes into a committed-range review unless the user explicitly asks.

Always report the selected range, current branch, and whether the worktree is dirty.

## Read-Only Boundaries

Never deploy, tag, publish, run migrations, rotate secrets, clear or warm caches, upload assets, trigger CI/CD, or change remote infrastructure.

Use only inspection commands such as:

```bash
git status --short
git tag --merged HEAD --sort=-creatordate
git diff --name-status <base>..<head>
git diff --stat <base>..<head>
git log --oneline --decorate --no-merges <base>..<head>
git diff -U3 <base>..<head> -- <path>
rg -n "<pattern>" .
```

Use `gh pr view` or `gh pr diff` only when the user provided a PR and the environment already permits read-only GitHub access.

## Review Checklist

Inspect the diff for release-linked risks:

- migrations, schema changes, indexes, seeds, and backfills
- new or changed env vars, config reads, secrets, flags, ports, and runtime dependencies
- cache keys, TTLs, invalidation, prewarm, and backward compatibility
- queue producers/consumers, topics, DLQs, idempotency, and deploy order
- object storage, CDN, templates, certificates, file permissions, and generated assets
- public API, event, webhook, SDK, CLI, and data-contract compatibility
- rollback behavior when new code and old data or old code and new data coexist
- observability, alerts, and runbook gaps only when the diff introduces release-critical behavior

Report secrets safely: path, line, variable name, secret type, and a redacted hint only. Never print full secret values, tokens, passwords, certificates, cookies, or private keys.

## Dirty Worktree

By default, committed-range release reviews exclude uncommitted and untracked changes.

- If dirty release-relevant files exist, add a confirmation item saying they were excluded and must be reviewed, committed, or discarded before release.
- If the user asks to include dirty changes, inspect them with read-only diffs and label them as uncommitted evidence.

## Findings

Use:

- `P0`: blocker; release should not proceed.
- `P1`: release-critical confirmation required.
- `P2`: lower-risk confirmation item before release.

Every finding needs evidence: file path and line when available, commit/range context, or the command limitation that prevented confirmation.

Conclusion values:

- `BLOCKED`: any P0 finding exists.
- `NEEDS_CONFIRMATION`: no P0, but P1/P2 confirmation items remain.
- `NO_BLOCKER_FOUND`: no P0-P2 findings found from available evidence.

Neutral tool limits belong in `unable_to_verify`; they do not change the conclusion unless tied to a release-critical diff.

## Output

Return the normal review output plus:

- `reviewed_range`
- `release_conclusion`
- `release_findings`
- `deployment_order`: only when multi-service, migration, queue, cache, or public-contract ordering matters
- `unable_to_verify`

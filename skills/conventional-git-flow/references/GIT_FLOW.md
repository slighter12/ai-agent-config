# Git Flow

Use this workflow for branch, commit, push, and PR preparation.

## Branch Names

Use:

```text
<type>/<area>-<short-slug>
```

Examples:

```text
feat/skill-creator-portable-metadata-validation
fix/skill-creator-compatibility-list
docs/conventional-git-flow
```

Rules:

- `type` matches the commit type.
- `area` is the main affected subsystem or package.
- `short-slug` is lowercase kebab-case.
- Keep names stable, descriptive, and under 80 characters when possible.
- If the user provides an issue key, put it after the type: `feat/ABC-123-skill-creator-metadata`.

## PR Title

Use the same header style as the commit:

```text
<type>: <description>
```

Do not include scope by default.

## PR Body

Use:

```md
## Summary

- Area: <area>
- <change summary>

## Validation

- <check or "Not run: <reason>">

## Risks

- <risk or "None identified">
```

Keep the body factual. Do not invent tests, issue numbers, approvals, or deployment status.

## Command Sequence

For read-only preparation, inspect:

```bash
git branch --show-current
git status --short
git diff --stat
git diff
git diff --staged
git log --oneline -5
```

For approved execution, use the smallest necessary sequence:

```bash
git switch -c <branch>
git add <explicit paths>
git commit -m "<commit headline>" -m "<commit body>"
git push -u origin <branch>
gh pr create --title "<pr title>" --body "<pr body>"
```

Stage explicit paths from the intended logical change. Include dependency manifests, lockfiles, generated dependency metadata, and version files when `VERSION_AND_DEPENDENCY_CLOSURE.md` classifies them as required closure for the change.

If the commit does not need a body, use only `git commit -m "<commit headline>"`.

Never run the execution sequence until the user has approved the exact branch, files, commit message, PR title, PR body, and commands.

## Sandbox And Permission Failures

Treat provider sandbox approvals and git workflow approval as separate gates. A direct user request such as "commit this" counts as git workflow approval for staging and committing the intended files after inspection; still ask if the file set, message, or scope is ambiguous.

If the current environment is known to block `.git` writes, request provider escalation on the first approved side-effect attempt instead of running a sandbox-first probe. Known blockers include a prior `.git/index.lock`, `.lock`, `Operation not permitted`, or `Permission denied` failure in this repo/session.

If the user has approved the git workflow but a command fails because the provider sandbox cannot write `.git` state or access the network, request escalated execution for the same command instead of changing the approved workflow.

When a side-effect command fails:

- If the error mentions `.git`, `index.lock`, `.lock`, `Operation not permitted`, `Permission denied`, DNS resolution, or network access, classify it as a likely sandbox or environment permission failure first.
- Re-run the exact failed command with the provider's escalation mechanism when available, using a concise justification tied to that command.
- Do not change the approved branch name, branch `<type>/` prefix, staged file list, Conventional Commit type, commit message, remote, PR title, or PR body just because the first attempt hit a sandbox or network error.
- Only propose a different branch name after a concrete ref conflict is confirmed, such as an existing branch, tag, or ref that blocks the requested path.
- For slash-prefixed branch names such as `feat/<slug>`, `fix/<slug>`, `docs/<slug>`, or `ci/<slug>`, do not assume the slash or type prefix is invalid. Confirm a real ref conflict before falling back to a flat branch name or removing the type prefix.

Use read-only checks to diagnose a possible ref conflict:

```bash
git branch --list <branch>
git branch --list <prefix>
git branch -r --list origin/<branch>
git branch -r --list origin/<prefix>
git show-ref --verify refs/heads/<prefix>
git show-ref --verify refs/tags/<prefix>
```

If these checks do not prove a ref conflict, keep the approved branch name exactly, including its `<type>/` prefix, and retry the failed side-effect command with the necessary sandbox permission.

## Existing PR Review Feedback

When the current branch already has a published PR and the user is applying reviewer feedback, keep the PR history transparent:

- Inspect reviewer comments and the local diff before deciding what to change.
- Stage only the explicit files needed for the feedback.
- Create a new Conventional Commit follow-up for the review fix.
- Push normally to the PR branch.

Do not propose `git commit --amend`, history-rewriting rebase, or force push for published PR feedback unless the user explicitly asks to amend, squash, or rewrite history. If rewriting is requested, show the amended commit and `--force-with-lease` push plan before asking for approval.

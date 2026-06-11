# Validation

Validation is part of the workflow even when the agent only prepares text.

## Before Side Effects

Confirm:

- Current branch is understood.
- Working tree state is understood.
- Staged and unstaged changes are distinguished.
- Proposed branch name matches the change type and area.
- Commit headline uses `<type>: <description>` with no scope.
- PR body includes Summary, Validation, and Risks.
- Exact commands are shown to the user.

## Before Commit

Confirm:

- Only intended files are staged.
- Generated or cache files are not staged unless intentionally changed.
- Dependency manifests, lockfiles, version metadata, and generated dependency files are included when they are required to reproduce, build, test, or publish the logical change.
- Major or ecosystem-specific incompatible dependency/version updates have a compatibility lookup or an explicit unknown-risk note.
- Commit message matches Conventional Commits.
- Sandbox or permission failures did not change the approved Conventional Commit type, headline, staged files, branch name, branch `<type>/` prefix, remote, PR title, or PR body.
- Multi-change commits include a concise factual body with relevant change details.
- Published PR reviewer feedback is being committed as a new follow-up commit unless the user explicitly requested history rewriting.
- Breaking changes, issue refs, or co-authors are in the body/footer when needed.

## Before PR Creation

Confirm:

- Branch is pushed or the exact push command is ready.
- Push is not a force push unless the user explicitly requested history rewriting.
- Base branch is known or `gh pr create` will prompt safely.
- PR title and body match the approved text.
- GitHub CLI is installed and authenticated, or the workflow stops with text output.

## After Side Effects

Report:

- Commands run.
- Exit status or success summary.
- New branch name.
- Commit hash, if a commit was created.
- PR URL, if a PR was created.
- Any skipped validation and why.

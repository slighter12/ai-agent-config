---
name: conventional-git-flow
description: Prepare Conventional Commit git workflows for branches, commits, pushes, and pull requests. Use when the user asks to commit current changes, create or review a branch name, commit message, PR title, PR body, or end-to-end git change flow. Avoid when the task is only code review, release notes, changelog writing, or provider-specific slash command setup.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.14"
---

# Conventional Git Flow

## Purpose

Enable the agent to produce and, with explicit approval, execute a safe git workflow using Conventional Commits-style naming for commits and PRs.
This skill is the source of truth for reusable git workflow policy; provider-specific git execution roles should stay thin and reference this skill instead of duplicating the full workflow.

## Use When

- The user wants a branch name, commit message, PR title, or PR body.
- The user asks to commit current work using Conventional Commits.
- The user asks to commit, stage and commit, or create a commit for current changes.
- The user asks to create a branch and prepare or open a pull request.
- The task requires translating a diff or change summary into git workflow text.

## Avoid When

- The request is only code review or implementation with no git workflow output.
- The user needs release notes or changelog content rather than commit or PR text.
- A provider-specific slash command wrapper is being authored instead of the shared workflow.

## Workflow

1. Inspect relevant git context: branch, status, staged diff, unstaged diff, and recent commit style.
2. Identify whether dependency manifests, lockfiles, version metadata, or generated dependency files are part of the logical change closure.
3. Identify the change type, affected area, branch slug, commit headline, PR title, and PR body.
4. Keep commit and PR title format as `<type>: <description>` without a scope in the header.
5. Put the affected area in the branch name and PR body instead of the commit header.
6. When applying reviewer feedback on an already published PR, default to a new follow-up commit and normal push.
7. For text-only preparation or non-simple git workflows, present the proposed branch name, staged files, commit message, PR title, PR body, and exact commands before any side effect.
8. Treat a direct user request such as "commit this" as approval to stage and commit the intended files after inspection; ask only if the file set, message, or scope is ambiguous.
9. Treat each direct simple git action request, such as creating a branch, staging, committing, pushing, or opening a PR for already-existing changes, as an independent explicit standing user authorization to delegate execution to the configured git execution role. The main session must inspect enough context to provide a complete git execution handoff before delegation.
10. Once a simple git action has been delegated with a complete handoff, do not add an extra main-session confirmation loop and do not ask the delegate to redo policy or scope decisions.
11. After execution, report the delegate's final confirmation, including commit or PR result, final status, and any follow-up verification.

## Git Execution Handoff

For simple delegated git actions, pass a compact handoff packet instead of asking the delegate to rediscover workflow policy:

- `repo_path`: absolute repository path.
- `actions`: ordered list of `branch`, `commit`, `push`, or `pr`.
- `intended_files`: exact files to stage, or explicit approval that every currently changed file belongs.
- `excluded_files`: files that must remain unstaged or dirty.
- `commit_message`: subject and body decision when committing, including an explicit empty-body decision when no body should be used.
- `branch`: branch name to create or use when branch, push, or PR actions are requested.
- `remote`, `base`, `head`: remote and PR refs when push or PR actions are requested.
- `pr_title` and `pr_body`: final PR text when opening a PR.
- `constraints`: no source edits, no amend, no rebase, no force push, no `--no-verify`, no broad staging unless explicitly approved.

If the handoff is complete, the delegate should execute it directly and return success or the first hard failure. If the handoff is incomplete, the delegate should fail closed instead of re-reading this skill or inferring missing scope.

## Tool And Side-Effect Boundaries

- Read-only git inspection is allowed when the active environment policy permits command execution.
- Creating branches, staging files, committing, pushing, and opening pull requests are side effects and require explicit user approval.
- When executing approved git side effects that write `.git` state, request provider escalation on the first attempt if the current environment is known to block `.git` writes, such as a prior `.git/index.lock`, `.lock`, `Operation not permitted`, or `Permission denied` failure in this repo/session.
- Do not stage broad file patterns by default; prefer explicit file paths from `git status`.
- Commit message precedence is: explicit user-provided message, explicit parent handoff message, then a generated Conventional Commit subject from the current diff. Recent commit style is local convention evidence, not authority over explicit user or skill rules.
- Do not push or run `gh pr create` unless the target branch, remote, title, and body have been shown to the user.
- Do not use `git commit --amend`, history-rewriting rebase, or force push for an already published PR unless the user explicitly asks to rewrite history.
- If an approved side-effect command fails with `.git` lock, `Operation not permitted`, permission, DNS, or network errors, treat it as a likely sandbox or environment failure and request escalated execution for the same command before changing the approved workflow.
- Never change the approved Conventional Commit type, headline, staged files, branch name, branch `<type>/` prefix, remote, PR title, or PR body because of a sandbox or permission failure.
- Do not replace slash-prefixed branch names such as `feat/<slug>`, `fix/<slug>`, `docs/<slug>`, or `ci/<slug>` unless read-only checks prove a real ref conflict; permission failures are not ref conflicts.
- If GitHub CLI is unavailable or unauthenticated, return the PR title/body and the exact command the user can run later.
- Do not delegate non-simple work to the git execution role: implementation, diagnosis, review, release flow, history rewrite, merge conflict resolution, unclear file scope, or PRs with unsafe unknown base/remote/title/body stay in the main session until the git execution scope is fixed.
- If runtime tool policy blocks delegation despite this repo preference, stop after read-only inspection and state that blocker clearly instead of silently continuing as if delegation was used.
- For simple git action prompts, do not run mutating git commands in the main session unless the user explicitly asks to bypass the configured git execution route after the blocker is reported.
- A `UserPromptSubmit` hook reminder is reinforcement only. If the hook appears, treat it as confirmation of this skill rule; if it does not appear but the prompt directly asks for branch creation, commit, push, or PR creation, still follow this skill rule and delegate.
- When the provider runtime supports explicit delegate cleanup, close the completed git execution delegate after its final result is no longer needed.

## Output

Return:

- `summary`: what git workflow was prepared or executed.
- `branch`: proposed or created branch name.
- `commit_message`: Conventional Commit message.
- `pr_title`: PR title.
- `pr_body`: PR body with summary, validation, and risks.
- `commands`: exact commands proposed or executed.
- `approval_needed`: side effects waiting for user approval.
- `manual_verification`: checklist for confirming branch, commit, push, or PR state.

## Version History

- v0.1.0 (2026-05-04): Initial portable Conventional Commit git flow skill.
- v0.1.1 (2026-05-04): Expand commit request routing triggers.
- v0.1.2 (2026-05-05): Prefer concise commit bodies for multi-change commits.
- v0.1.3 (2026-05-08): Prefer follow-up commits for published PR review feedback.
- v0.1.4 (2026-05-08): Distinguish sandbox permission failures from branch naming conflicts before changing approved git commands.
- v0.1.5 (2026-05-13): Add dependency and version metadata closure rules for commit readiness.
- v0.1.6 (2026-05-22): Preserve approved git workflow details across sandbox retries and allow upfront escalation for known `.git` write restrictions.
- v0.1.7 (2026-06-01): Clarify configured git execution delegation as fixed commit execution, not branch or PR ownership.
- v0.1.8 (2026-06-03): Add explicit commit trigger phrases and clarify that plain git requests do not implicitly delegate.
- v0.1.9 (2026-06-03): Treat simple git workflow requests as standing authorization for configured git execution delegation.
- v0.1.10 (2026-06-03): Fail closed when simple git workflow delegation is blocked instead of falling back to main-session git mutation.
- v0.1.11 (2026-06-03): Treat branch creation, commit, push, and PR creation as independent simple git action delegation triggers, with hook reminders as non-authoritative reinforcement only.
- v0.1.12 (2026-06-03): Keep the shared skill provider-neutral, remove extra confirmation loops for unambiguous simple git action delegation, and require delegate cleanup when supported.
- v0.1.13 (2026-06-03): Clarify shared skill ownership over git workflow policy, keep provider role prompts thin, and define commit message precedence over recent commit style.
- v0.1.14 (2026-06-04): Add compact git execution handoff packets so delegates can execute complete simple actions without repeating policy discovery.

## References

- `references/INDEX.md` - Navigation for commit, branch, PR, and validation rules.
- `references/VERSION_AND_DEPENDENCY_CLOSURE.md` - Dependency, lockfile, and version metadata commit-readiness rules.

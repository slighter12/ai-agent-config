---
name: conventional-git-flow
description: Prepare Conventional Commit git workflows for branches, commits, pushes, and pull requests. Use when the user asks to commit current changes, create or review a branch name, commit message, PR title, PR body, or end-to-end git change flow. Avoid when the task is only code review, release notes, changelog writing, or command-wrapper setup.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.23"
---

# Conventional Git Flow

## Purpose

Enable the agent to produce and, with explicit approval, execute a safe git workflow using Conventional Commits-style naming for commits and PRs.
This skill is the source of truth for reusable git workflow policy; execution wrappers should stay thin and reference this skill instead of duplicating the full workflow.

## Use When

- The user wants a branch name, commit message, PR title, or PR body.
- The user asks to commit current work using Conventional Commits.
- The user asks to commit, stage and commit, or create a commit for current changes.
- The user asks to create a branch and prepare or open a pull request.
- The task requires translating a diff or change summary into git workflow text.

## Avoid When

- The request is only code review or implementation with no git workflow output.
- The user needs release notes or changelog content rather than commit or PR text.
- A command wrapper is being authored instead of the shared workflow.

## Workflow

1. For direct git actions, collect the read-only Git Context Pack before any side effect.
2. Preserve user-provided branch, file scope, commit message, PR title/body, remote, base, or head as explicit constraints.
3. If the current task is already a delegated git execution handoff, use the handoff context and execute only the requested git actions.
4. Keep commit and PR title format as `<type>: <description>` without a scope in the header.
5. Put the affected area in the branch name and PR body instead of the commit header.
6. When applying reviewer feedback on an already published PR, default to a new follow-up commit and normal push.
7. For non-delegated execution, use the Git Context Pack results and present proposed branch names, staged files, commit messages, PR text, and exact commands before any side effect unless the direct request already approved that exact action and intended scope.
8. If delegated execution fails because of incomplete context, unsafe scope, permission, auth, or network issues, follow the failure and escalation rules below instead of changing the approved workflow.
9. Treat a direct request such as "commit this" as approval for only the requested simple git action and intended files after inspection; ask if file set, message, or scope is ambiguous.
10. Do not run duplicate mutating git commands in both caller and delegate for the same requested action.
11. After execution, report the action result, final status, and any follow-up verification.

## Git Context Pack

Before git side effects, run this fixed read-only pack:

- `git status --short --branch`
- `git diff --name-status`
- `git diff --stat`
- `git diff --cached --name-status`
- `git diff --cached --stat`
- `git log -5 --oneline`

After choosing intended files, inspect only those files with `git diff -- <paths...>` and, when staged content exists, `git diff --cached -- <paths...>`. If a command is unavailable or fails, record that result and run only the narrow follow-up command needed to fill the missing context.

## Git Execution Handoff

For simple delegated git actions, pass a compact context handoff packet so the delegate can inspect and decide the git workflow:

- `repo_path`: absolute repository path.
- `original_user_request`: the exact user request that triggered the git action.
- `requested_actions`: ordered list of `branch`, `commit`, `push`, or `pr`.
- `explicit_user_constraints`: user-provided files, exclusions, branch, commit message, remote/base/head, PR title/body, or approval for all current changes.
- `git_context_pack_results`: optional outputs or failures from the Git Context Pack, plus any intended-file diff outputs already inspected.
- `known_context`: current task summary, known files touched in this session, and validation run or skipped.
- `safety_constraints`: no source edits, no amend, no rebase, no force push, no `--no-verify`, preserve user changes, and no broad staging unless explicitly approved.

If the handoff is complete, the delegate should use any fresh Git Context Pack results, choose safe intended files and git text, execute the requested actions, and return success or the first hard failure. If context is missing or stale, the delegate should run the fixed pack or the narrow intended-file diff commands itself without waiting for the parent to finish. If the handoff is incomplete or the diff scope is mixed, the delegate should fail closed instead of guessing.

## Tool And Side-Effect Boundaries

- Read-only git inspection is allowed when the active environment policy permits command execution.
- Creating branches, staging files, committing, pushing, and opening pull requests are side effects and require explicit user approval.
- When executing approved git side effects that write `.git` state, request provider escalation on the first attempt if the current environment is known to block `.git` writes, such as a prior `.git/index.lock`, `.lock`, `Operation not permitted`, or `Permission denied` failure in this repo/session.
- Do not stage broad file patterns by default; prefer explicit file paths from `git status`.
- Commit message precedence is: explicit user-provided message passed in the handoff, then a generated Conventional Commit subject from the current diff. Recent commit style is local convention evidence, not authority over explicit user or skill rules.
- Do not push or run `gh pr create` unless the action is explicitly requested or approved and the delegate has resolved the target branch, remote, title, and body.
- When running `gh pr create`, pass the explicit or resolved PR body through `--body-file` from a temporary file outside tracked paths. Do not inline multiline Markdown in `--body`.
- Do not use `git commit --amend`, history-rewriting rebase, or force push for an already published PR unless the user explicitly asks to rewrite history.
- If an approved side-effect command fails with `.git` lock, `Operation not permitted`, permission, DNS, or network errors, treat it as a likely sandbox or environment failure and request escalated execution for the same command before changing the approved workflow.
- Never change explicit user constraints or chosen Conventional Commit type, headline, staged files, branch name, branch `<type>/` prefix, remote, PR title, or PR body because of a sandbox or permission failure.
- Do not replace slash-prefixed branch names such as `feat/<slug>`, `fix/<slug>`, `docs/<slug>`, or `ci/<slug>` unless read-only checks prove a real ref conflict; permission failures are not ref conflicts.
- If GitHub CLI is unavailable or unauthenticated, return the PR title/body and the exact command the user can run later.
- Do not treat implementation, diagnosis, review, release flow, history rewrite, merge conflict resolution, or PRs with unsafe unknown base/remote/title/body as simple git execution until the git scope is fixed.
- Do not invent a git execution delegate when the caller did not provide a delegated handoff.
- This shared skill owns git workflow rules and handoff shape, not caller orchestration.

## Output

Return:

- `summary`: what git workflow was prepared or executed.
- `branch`: proposed or created branch name.
- `commit_message`: Conventional Commit message.
- `pr_title`: PR title.
- `pr_body`: PR body resolved after checking repository template/checklist sources; it preserves discovered requirements or uses fallback summary, validation, and risks sections only when no template/checklist exists.
- `commands`: exact commands proposed or executed.
- `approval_needed`: side effects waiting for user approval.
- `manual_verification`: checklist for confirming branch, commit, push, or PR state.

## Version History

- v0.1.0 (2026-05-04): Initial portable Conventional Commit git flow skill.
- v0.1.1 (2026-05-04): Expand commit request phrase handling.
- v0.1.2 (2026-05-05): Prefer concise commit bodies for multi-change commits.
- v0.1.3 (2026-05-08): Prefer follow-up commits for published PR review feedback.
- v0.1.4 (2026-05-08): Distinguish sandbox permission failures from branch naming conflicts before changing approved git commands.
- v0.1.5 (2026-05-13): Add dependency and version metadata closure rules for commit readiness.
- v0.1.6 (2026-05-22): Preserve approved git workflow details across sandbox retries and allow upfront escalation for known `.git` write restrictions.
- v0.1.7 (2026-06-01): Clarify delegated git execution as scoped to assigned actions.
- v0.1.8 (2026-06-03): Add explicit commit trigger phrases and clarify that plain git requests do not implicitly delegate.
- v0.1.9 (2026-06-03): Add delegated simple-git approval wording; superseded by current caller/handoff constraints.
- v0.1.10 (2026-06-03): Add blocked-delegation failure wording; superseded by current caller/handoff constraints.
- v0.1.11 (2026-06-03): Treat branch creation, commit, push, and PR creation as independent simple git actions.
- v0.1.12 (2026-06-03): Keep the shared skill portable and remove extra confirmation loops for unambiguous simple git action delegation.
- v0.1.13 (2026-06-03): Clarify shared skill ownership over git workflow policy, keep provider role prompts thin, and define commit message precedence over recent commit style.
- v0.1.14 (2026-06-04): Add compact git execution handoff packets for simple delegated actions.
- v0.1.15 (2026-06-23): Require PR creation to use body files so multiline Markdown stays intact.
- v0.1.16 (2026-06-23): Rebalance simple git delegation so the git execution role decides workflow details from a compact context handoff.
- v0.1.17 (2026-06-24): Simplify runtime boundary rules and keep caller orchestration outside the shared workflow.
- v0.1.18 (2026-06-25): Add fixed Git Context Pack.
- v0.1.19 (2026-06-25): Clarify read-only Git Context Pack collection for simple git actions.
- v0.1.20 (2026-06-26): Move caller orchestration out of the shared skill.
- v0.1.21 (2026-06-26): Remove caller-orchestration noise from the shared workflow.
- v0.1.22 (2026-06-26): Require PR bodies to preserve repository templates and required checklists.
- v0.1.23 (2026-06-26): Make PR template/checklist discovery an explicit gate before fallback bodies.

## References

- `references/INDEX.md` - Navigation for commit, branch, PR, and validation rules.
- `references/VERSION_AND_DEPENDENCY_CLOSURE.md` - Dependency, lockfile, and version metadata commit-readiness rules.

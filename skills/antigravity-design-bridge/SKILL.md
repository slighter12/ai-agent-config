---
name: antigravity-design-bridge
description: Coordinate optional Antigravity CLI assistance for UI or visual design critique and low-risk bounded edits. Use when UI, visual design, layout, UX work, rejected design quality, or redesign direction would benefit from Antigravity critique or direct assistance under strict file-scope boundaries. Avoid when Antigravity CLI is not needed, when the task is ordinary backend or non-visual implementation, or when safe file scope cannot be established.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.4.0"
---

# Antigravity Design Bridge

## Purpose

Use Antigravity CLI as an optional external design collaborator while keeping the primary agent responsible for scope control, review, final judgment, and user-facing delivery.

This skill supports two modes:

- `advisory`: Antigravity returns UI or UX critique, layout alternatives, visual hierarchy suggestions, or implementation guidance without editing files.
- `direct-edit`: Antigravity may modify only explicitly bounded low-risk UI/design files, and the primary agent must review the resulting diff before accepting or reporting the work.

## Use When

- The user explicitly asks to use Antigravity CLI for visual design, UI critique, UX alternatives, layout review, or design implementation assistance.
- The task benefits from a second model focused on visual or product-design judgment.
- A frontend or UI task has a clear visual/design goal and bounded target files where Antigravity can provide useful critique or low-risk edits.
- The primary agent should review Antigravity's output or diff before accepting any design changes.

## Avoid When

- The user asks for ordinary coding, backend work, commits, PRs, or current-change review.
- Antigravity CLI is not installed and installing or configuring it is outside the current task.
- Safe file scope cannot be established for direct edits.
- The task touches auth, security, data models, API contracts, migrations, infrastructure, dependencies, secrets, production data, or broad refactors.
- The task requires Antigravity changes to be accepted without primary-agent review.
- A more specific frontend or design policy skill should lead the work.

## Workflow

1. Confirm or infer the design goal, target screen or component, constraints, available artifacts, and expected output.
2. If `.design/` exists, read only the relevant design-pack files and include their constraints in the bounded prompt.
3. Check whether `agy` is available with a read-only command such as `command -v agy`.
4. If Antigravity CLI is unavailable, report that and continue with local review or ask whether installation should be handled separately.
5. Choose `advisory` unless every direct-edit gate below is satisfied:
   - The task is UI, visual design, layout, styling, or interaction polish.
   - The target files are already identified from the prompt or repo inspection.
   - The expected edit is low risk and does not require new dependencies, public API changes, data-flow rewrites, backend changes, migrations, or config changes.
   - The current worktree and target files can be checked so unrelated user changes are not overwritten.
   - The CLI can be run with `--sandbox` and `--add-dir` restricted to the smallest containing directory for the target files.
6. Use Gemini 3.5 Flash Medium for normal advisory critique. Use Gemini 3.5 Flash High for rejected design, redesign, or direct-edit work.
7. In `advisory` mode, pass a bounded prompt with `agy --sandbox --print-timeout <timeout> --print "<prompt>"` or `agy -p` that asks for critique or alternatives, not direct file mutation.
8. In `direct-edit` mode:
   - Inspect the target files and current git status before invoking Antigravity.
   - Use `agy --sandbox --add-dir <target-dir> --print-timeout <timeout> --print "<bounded edit prompt>"`.
   - Name the allowed files and design objective in the prompt.
   - Tell Antigravity not to edit other files, add dependencies, change configuration, or touch unrelated behavior.
9. Verify the effective model by checking the latest `~/.gemini/antigravity-cli/log/cli-*.log` entry for the propagated model label. If it does not match the intended Medium or High tier, report the mismatch and ask the user to switch Antigravity's selected model in its UI or interactive `/model` flow.
10. Review Antigravity's response and, for direct edits, inspect the actual diff. Accept, reject, or adapt recommendations based on repo constraints and the active design policy.
11. If Antigravity edits are directionally useful but rough, the primary agent may refine them before final delivery.

## Tool And Side-Effect Boundaries

- Do not install Antigravity CLI automatically.
- Do not give Antigravity credentials, secrets, production data, or unrelated repository context.
- Do not use `--dangerously-skip-permissions`.
- Do not claim that `agy --print` can force the model with a CLI flag unless the installed CLI exposes and verifies that flag.
- Do not let Antigravity mutate files outside the selected `--add-dir` or outside the named target files.
- Prefer short prompts with exact file paths, screenshots, local URLs, or design goals over dumping broad repo context.
- Treat Antigravity output and edits as candidate design assistance, not authoritative implementation.
- Reject or manually correct Antigravity edits that introduce broad refactors, new dependencies, hidden behavior changes, unrelated formatting churn, accessibility regressions, or non-visual scope.
- If target files contain unclear pre-existing user changes, use advisory mode or ask before direct editing.

## Output

Return:

- `summary`: what Antigravity was asked or why it was skipped.
- `mode`: advisory, direct-edit, unavailable, skipped, or failed.
- `antigravity_status`: available, unavailable, skipped, or failed.
- `model_check`: intended Medium or High tier and the observed log label, if Antigravity was invoked.
- `prompt_scope`: files, screenshots, URLs, or design goals shared with Antigravity.
- `changed_files`: files Antigravity modified in direct-edit mode, if any.
- `design_recommendations`: accepted recommendations or accepted edit direction only.
- `rejected_recommendations`: rejected or unsafe suggestions, if any.
- `manual_verification`: diff checks, visual checks, or screenshots needed after implementation.

## Version History

- v0.1.0 (2026-05-08): Initial optional Gemini CLI design bridge.
- v0.2.0 (2026-05-21): Rename bridge for Antigravity CLI migration and switch command detection to `agy`.
- v0.3.0 (2026-05-31): Add advisory and bounded direct-edit modes with autonomous low-risk selection gates.
- v0.4.0 (2026-05-31): Add `.design/` prompt context and Medium/High Antigravity model verification policy.

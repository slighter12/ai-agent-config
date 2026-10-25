---
name: pr
description: Format high-signal pull request bodies with a visual Summary, Before/After Evidence, and Merge Danger assessment (One-way vs. Two-way door + Blast Radius). Use whenever drafting or updating a PR description.
metadata:
  invocation: model
---

# Pull Request Body (`/pr`)

Write a PR body that lets a human reviewer grasp the shape of the change, see proof that it works, and judge merge risk in under 30 seconds.
Keep the body under about 40 lines, since it often becomes the squash commit body; move full SHAs, long logs, and file-by-file lists into a collapsed `<details>` block or a linked artifact.

## Pre-Check: Respect Repository PR Templates

1. Check the repository for a PR or merge request template in every standard location: `.github/pull_request_template.md`, `.github/PULL_REQUEST_TEMPLATE.md`, `.github/PULL_REQUEST_TEMPLATE/`, a root `PULL_REQUEST_TEMPLATE.md`, `docs/PULL_REQUEST_TEMPLATE.md`, and `.gitlab/merge_request_templates/`. Names are case-insensitive.
2. If a repo template exists, preserve its required headings/checklists and place `Summary`, `Evidence`, and `Merge Danger` cleanly inside or alongside its structure so repo automation never breaks.

## The Three Required Sections

Start directly at `## Summary` with no conversational preamble:

### 1. `## Summary` (A Visual, Not a Wall of Prose)
Pick the **smallest visual representation** that makes the structural change obvious before the reviewer reads a single line of the diff (usually 1 visual + 2–3 tight bullets). The bullets name what changes for the consumer or caller, and what the next maintainer inherits:
* A compact call tree or data-flow sketch (`Caller -> SproutSeam -> Adapter`)
* A before/after interface signature or state-transition table
* A file/module tree highlighting where a seam was deepened or sprouted
* A small Mermaid diagram (only when rendering on GitHub/GitLab web UI)

End the Summary with a one-line **Scope**: what the change covers and what it deliberately leaves out (a follow-up, a known gap).

### 2. `## Evidence` (Concrete Before & After Proof)
Never write `"All tests pass"` as a bare claim. Show the concrete observation that was false/missing before and is verified now:
* **Behavioral Seam Test**: The exact test case name(s) and assertion summary that failed on the base commit and pass now.
* **Runtime Verification (`Prove It Works`)**: The exact CLI command, API response, benchmark, or UI screenshot comparison captured before and after the change.

In `## Evidence` and `## Merge Danger`, every behavior, performance, or risk claim carries its evidence pointer or an explicit label: **Measured** (observed, with the pointer), **Inferred** (reasoned from code or docs), or **Unverified**.
For performance, lead with one primary number and its unit, `before -> after` (`p95 412ms -> 188ms`), and link the rest. An inconclusive run, or a check on a different surface than the claim (a unit test for a UI claim, staging for a production-only path), is not evidence: label the claim **Unverified**.

### 3. `## Merge Danger` (Door Type & Blast Radius)
Turn `"is this safe to merge?"` into an explicit, falsifiable claim:
* **Door Type**:
  * **Two-Way Door**: Reverting the commit cleanly restores previous state with no external side effects.
  * **One-Way Door**: Hard or impossible to reverse by `git revert` alone (e.g., destructive DB schema migration, persisted data format change, public contract removal, outbound external side effects). Flag one-way doors prominently so reviewers slow down.
* **Blast Radius**: Name the exact callers, endpoints, jobs, or user flows that could break if this change has a defect, and note any rollout guards (feature flags, expand-contract stage). Include the cross-boundary consumers the `blast-radius` skill finds beyond a symbol search.
* **Safety Fact**: Name the one fact the change is safe because of. Before writing this section, use the `blast-radius` skill to grade that fact on its `## Proof Ladder`. Levels 4–5 are **Measured**, 2–3 **Inferred**, 1 **Unverified**. Below level 4, report the fact as unproven the way the `blast-radius` skill's report does. If the change is a One-Way Door or touches a cross-boundary contract (serialized or wire formats, DB columns, FFI, feature flags) and the highest Safety Fact is below level 4, use the `blast-radius` skill for its full audit before the PR is opened.

## Plain Speech

Write every sentence with a concrete actor and verb ("the migration drops column `x`", "`retry()` now returns the last error") and attach evidence to every claim. Before finalizing the body, use the `unslop` skill on it, keeping identifiers, numbers, and evidence pointers unchanged. If the branch touches 3 or more files or has more than one commit, also use the `unslop` skill to audit the prose the branch adds relative to its base branch (docs and code comments), and report proposed rewrites instead of editing files the request did not name.

---
name: pr
description: Format high-signal pull request bodies with a visual Summary, Before/After Evidence, and Merge Danger assessment (One-way vs. Two-way door + Blast Radius). Use whenever drafting or updating a PR description.
metadata:
  invocation: model
---

# Pull Request Body (`/pr`)

Write a PR body that lets a human reviewer grasp the shape of the change, see proof that it works, and judge merge risk in under 30 seconds.

## Pre-Check: Respect Repository PR Templates

1. Check if `.github/pull_request_template.md` (or `.github/PULL_REQUEST_TEMPLATE/`) exists in the repository.
2. If a repo template exists, preserve its required headings/checklists and place `Summary`, `Evidence`, and `Merge Danger` cleanly inside or alongside its structure so repo automation never breaks.

## The Three Required Sections

Start directly at `## Summary` with no conversational preamble:

### 1. `## Summary` (A Visual, Not a Wall of Prose)
Pick the **smallest visual representation** that makes the structural change obvious before the reviewer reads a single line of the diff (usually 1 visual + 2–3 tight bullets):
* A compact call tree or data-flow sketch (`Caller -> SproutSeam -> Adapter`)
* A before/after interface signature or state-transition table
* A file/module tree highlighting where a seam was deepened or sprouted
* A small Mermaid diagram (only when rendering on GitHub/GitLab web UI)

### 2. `## Evidence` (Concrete Before & After Proof)
Never write `"All tests pass"` as a bare claim. Show the concrete observation that was false/missing before and is verified now:
* **Behavioral Seam Test**: The exact test case name(s) and assertion summary that failed on the base commit and pass now.
* **Runtime Verification (`Prove It Works`)**: The exact CLI command, API response, benchmark, or UI screenshot comparison captured before and after the change.

### 3. `## Merge Danger` (Door Type & Blast Radius)
Turn `"is this safe to merge?"` into an explicit, falsifiable claim:
* **Door Type**:
  * **Two-Way Door**: Reverting the commit cleanly restores previous state with no external side effects.
  * **One-Way Door**: Hard or impossible to reverse by `git revert` alone (e.g., destructive DB schema migration, persisted data format change, public contract removal, outbound external side effects). Flag one-way doors prominently so reviewers slow down.
* **Blast Radius**: Name the exact callers, endpoints, jobs, or user flows that could break if this change has a defect, and note any rollout guards (feature flags, expand-contract stage).

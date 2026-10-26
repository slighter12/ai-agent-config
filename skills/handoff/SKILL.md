---
name: handoff
description: Compact the current conversation into a handoff document for another agent to pick up.
argument-hint: "What will the next session be used for?"
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Handoff

Write a handoff document summarising the current conversation so a fresh agent can continue the work. Save to the temporary directory of the user's OS - not the current workspace: resolve it from `$TMPDIR`, falling back to `/tmp`, or `%TEMP%` on Windows. Name the file `handoff-<timestamp>-<slug>.md` and print its exact path so the user can pass it to the next session.

Include a "suggested skills" section in the document, naming which skills the next agent should call the Skill tool for. Name a user-invoked skill as the `/command` the user runs.

Do not duplicate content already captured in other artifacts (specs, plans, ADRs, issues, commits, diffs). Reference them by path or URL instead.

Redact any sensitive information, such as API keys, passwords, or personally identifiable information.

If the user passed arguments, treat them as a description of what the next session will focus on and tailor the doc accordingly.

## Local additions

* **Current stage and next action**: state where the work stands (for example "Grilling in progress: 2 branches settled, 1 open", "Spec settled: ready for Ticket 02", or "Debugging repro loop built, testing Hypothesis 2") and the next action.
* **Settled decisions**: record every choice the user already approved so the next session never re-asks it: settled `grilling` branches with any open `⚡ Tension` items and Spike / Prototype Prerequisites, the Agency Gate option chosen for each target, and the agreed test seams and evidence path for each slice.
* **Run state**: for an `implement` run, the integration branch and the branch it was created from (when used), the run's `BASE_SHA`, and the `decisions.tsv` path, so the next session reconciles landed tickets before dispatching.
* **Verified vs. unverified**: list facts confirmed by a `file:line` read or a command (cite it) apart from hypotheses the resuming agent must verify before acting on them.

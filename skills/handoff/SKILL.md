---
name: handoff
description: Compact the current conversation into a portable, stage-agnostic handoff document in the OS temp directory so a fresh session can resume without context loss.
argument-hint: What will the next session focus on?
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Handoff (Stage-Agnostic Context Transfer)

Write a self-contained handoff document summarizing the current conversation so a fresh agent session (in a new window, worktree, or harness) can continue the work seamlessly at **any** point in a workflow (mid-grilling, post-spec, mid-implementation, or mid-debugging).

## Output Location
Save the handoff file to the OS temporary directory (`$TMPDIR`, falling back to `/tmp` or `%TEMP%` on Windows) as `handoff-<timestamp>-<slug>.md` (or `.html` if the user requests browser viewing) — **never** inside the repository workspace. Print the exact file path so the user can pass it to the next window.

## Required Handoff Structure

1. **Current Stage & Next Immediate Action**:
   * State the exact current state (e.g., `"Grilling in progress — 2 branches settled, 1 open"`, `"Spec settled — ready for Ticket 02"`, or `"Debugging repro loop built, testing Hypothesis 2"`).
   * If the user passed arguments to `/handoff`, tailor the next session's focus to those arguments.
2. **Context Pointers (Never Duplicate Artifacts)**:
   * Link by exact file path, commit SHA, branch name, or issue number to existing specs, tickets, `GLOSSARY.md`, `docs/adr/`, and modified files instead of pasting their contents.
3. **Settled Decisions & User Agency Choices**:
   * Record every constraint or architectural choice the user already approved in this window so the next session never re-asks settled questions:
     * settled `grilling` branches, plus any open `⚡ Tension` items and **Spike / Prototype Prerequisites**;
     * the Agency Gate option — (a) Local minimal edit, (b) Sprout at the seam, (c) Characterization check + local tidy first, or (d) Migrate — chosen for each target;
     * the agreed test seams and Evidence Path (`Seam TDD` or `Runtime Verification`) for each slice.
4. **Verified Facts vs. Unverified Assumptions (Anti-Hallucination Guard)**:
   * **Verified**: Facts confirmed by reading specific `file:line` locations or running commands (cite the command/file).
   * **Unverified / Open**: Hypotheses or claims that the resuming agent MUST verify before acting on.
5. **Suggested Skill for the Next Window**:
   * Name the single skill (if any) the resuming agent should load first (e.g., `grilling`, `to-spec`, `implement`, `tdd`, `diagnosing-bugs`, `code-review`).
6. **Redaction**:
   * Redact all API keys, tokens, passwords, credentials, secrets, and personally identifiable information (`<REDACTED>`).

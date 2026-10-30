---
name: tidy
description: Simplify recently changed code or a named module for clarity while preserving behavior exactly, verifying after every step. Use when asked to simplify or clean up code, or for an agreed tidy step after tests are green.
metadata:
  invocation: model
---

# Tidy

Simplify code by reducing complexity while preserving exact behavior. The goal is not fewer lines - it's code that is easier to read, understand, modify, and debug. Every simplification must pass a simple test: "Would a new team member understand this faster than the original?"

`code-review` reports findings without editing; `improve-codebase-architecture` surveys the whole repository; `tidy` edits, but only the recent diff or the module the user names, and only in ways that preserve behavior.

## When to Use

- After a feature is working and tests pass, but the implementation feels heavier than it needs to be
- During code review when readability or complexity issues are flagged
- When you encounter deeply nested logic, long functions, or unclear names
- When refactoring code written under time pressure
- When consolidating related logic scattered across files
- After merging changes that introduced duplication or inconsistency

**When NOT to use:**

- Code is already clean and readable - don't simplify for the sake of it
- You don't understand what the code does yet - comprehend before you simplify
- The code is performance-critical and the "simpler" version would be measurably slower
- You're about to rewrite the module entirely - simplifying throwaway code wastes effort

## Boundaries

- Edit only on a change request: the user asked to simplify or clean up code, or agreed to a tidy step. Otherwise report candidates without editing.
- Scope: the changes since a fixed point the user names (or the uncommitted diff), or one named module.
- Before editing, check the request and the target. For an urgent fix or one the user wants with minimal blast radius, do not tidy; report that and stop. For an Active Migration Conflict or a Legacy Monolith, use the `codebase-design` skill and present its Agency Gate Menu instead of tidying.
- Never change public interfaces, error-return shape, state transitions, or I/O seams; those are structural contracts governed by Project Law.
- Add no dependencies.
- Leave changes uncommitted. Commits, branches, and PRs need their own approval.

## The Five Principles

### 1. Preserve Behavior Exactly

Don't change what the code does - only how it expresses it. All inputs, outputs, side effects, error behavior, and edge cases must remain identical. If you're not sure a simplification preserves behavior, don't make it.

Before every change, ask:

- Does this produce the same output for every input?
- Does this maintain the same error behavior?
- Does this preserve the same side effects and ordering?
- What baseline checks will reveal a behavior change?

### 2. Follow Project Conventions

Simplification means making code more consistent with the codebase, not imposing external preferences.

Before simplifying:

1. Read `AGENTS.md` / project conventions
2. Study how neighboring code handles similar patterns
3. Match the project's style for imports, naming, function style, error handling, and type annotations

Simplification that breaks project consistency is not simplification - it's churn.

### 3. Prefer Clarity Over Cleverness

Explicit code is better than compact code when the compact version requires a mental pause to parse.

- Replace nested ternaries with readable control flow
- Replace dense inline transforms with named intermediate steps when they clarify intent
- Keep helpful names even if they cost a few extra lines

### 4. Maintain Balance

Watch for over-simplification:

- Don't inline away names that carry meaning
- Don't merge unrelated logic into one larger function
- Don't remove abstractions that serve testability or extensibility
- Don't optimize for line count over comprehension

### 5. Scope to What Changed

Default to simplifying recently modified code. Avoid unrelated drive-by refactors unless explicitly asked.

## Process

### Step 1: Understand Before Touching

Before changing or removing anything, understand why it exists.

Answer:

- What is this code's responsibility?
- What calls it? What does it call?
- What are the edge cases and error paths?
- Are there tests that define expected behavior?
- Why might it have been written this way?

If you can't answer these, read more context first. Before deleting or restructuring logic that no test or comment explains, or any defensive guard (a null check, retry, timeout, clamp, or flag), use the `why` skill (`git log -S`, `git blame`, `docs/adr/`, linked issues) to confirm what constraint it originally guarded. Keep comments that explain an algorithmic, concurrency or ordering, or performance reason.

### Step 2: Establish the Baseline

Run the tests that cover the scope and record the command and result. If they fail before you start, stop and report. If a candidate lacks behavioral coverage, pause and ask the user whether to add a characterization test first (Agency Gate option (c)) or drop the candidate from this pass.

### Step 3: Look for Simplification Opportunities

Use the `codebase-design` skill and apply its §3 rules (the Deletion Test, Subtract Before You Add, the Laziness Protocol, the 30-Second Test, and Structure over branches) to the scope.

Signals:

- Deep nesting
- Long functions with mixed responsibilities
- Nested ternaries
- Boolean flag arguments
- Repeated conditionals
- Generic or misleading names
- Duplicated logic
- Dead code
- Wrappers or abstractions that add no value

### Step 4: Apply Changes Incrementally

Make one simplification at a time.

For each simplification:

1. Make the change
2. Rerun the baseline tests for the file or package the step touched to check preservation
3. Keep it only when the evidence supports preservation

Revert a change whose checks fail or whose preservation you cannot show; do not stack a second change on an unverified first. Separate refactoring from feature work whenever possible.

### Step 5: Verify the Result

After simplifying, confirm:

- The code is genuinely easier to understand
- The diff is clean and reviewable
- Project conventions still match
- No behavior, error handling, or side effects changed
- The whole Step 2 baseline, rerun once after the last step, still passes (record the command and result)

## Report

- **Changed**: each simplification with its location and the check that passed after it.
- **Skipped**: candidates you left alone, with the reason (no coverage, unclear purpose, performance-critical, out of scope).
- **Verdict**: Verified, Not verified, or Inconclusive for "behavior preserved", with the final command and result.

---
name: tdd
description: Language-agnostic, task-proportional test-first workflow with mechanical smell redlines and post-green test pruning. Use when building features, fixing reproducible bugs, or verifying behavioral seams.
metadata:
  invocation: model
---

# TDD & Task-Proportional Verification

Verify behavior through the tightest, most durable feedback loop that matches the architectural layer you are changing. Keep the test suite lean, behavioral, and 100% language-agnostic.

A good test verifies behavior through a public interface and reads like a specification of a capability, so it survives refactors. Before writing the first test, read [tests.md](tests.md) for good and bad examples; before adding any test double, read [mocking.md](mocking.md). When no seam was agreed before Step 3, or the test needs a new or changed exported interface, use the `codebase-design` skill to decide how deep the module is and where the seam belongs.

## Step 1: Pre-Flight Context, Agency & Layer Routing Gate

Before writing any test or production code:

1. **Read Project Law & Neighbor Tests**:
   * Check `AGENTS.md` / `CLAUDE.md`, `GLOSSARY.md`, and `docs/adr/` for project testing rules.
   * Open 1–2 existing test files in the same directory or package to observe the established test harness, assertion style, and boundary setup.
2. **Evaluate Project State & Yield Agency on Friction**:
   * **Active Migration Conflict** (a project ADR defines a new testing/architecture pattern while the target or its test file uses a deprecated one) and/or **Legacy Monolith ("Shit Mountain")** (the target has no tests and tightly couples I/O with logic): **stop**, use the `codebase-design` skill, and present its Agency Gate Menu once, with your recommendation. Under option (b), run this loop on the Sprout.
3. **Route by Task & Architectural Layer (When to Use TDD vs. Runtime Verification)**:

| Layer / Task Type | Required Verification Path | Why |
| :--- | :--- | :--- |
| **Domain / State / Calculation Logic** | **Seam-Level Red-Green Loop** | Deterministic inputs and state outputs yield high-leverage, refactor-safe behavioral tests. |
| **Reproducible Bug Fix** | **Reproduction-First Red-Green** | A failing test (or deterministic repro command) proves the bug's root cause before the fix lands. |
| **Thin Wiring / Glue / Config / CLI / UI** | **Runtime Verification (`Prove It Works`)** | Unit tests on pass-through glue only test mocks. Skip unit TDD here and verify by driving the real command/endpoint/build and capturing before/after output. |
| **Exploratory Spike / Prototype** | **Disposable Probe / Script** | Do not lock premature interfaces in tests until the spike settles the design. |

* Prefer no new unit test over a brittle, mock-heavy unit test — **provided** you satisfy the **Evidence Conservation Law** (every skipped unit test must be replaced by captured runtime verification output).

## Step 2: Mechanical Redlines (Check Syntactically on Every Test)

Apply these four low-reasoning structural checks to every test you write:

1. **Copy-Neighbor First (When Clean)**:
   * If the package already has a clean seam setup (e.g., an in-memory fake, table-driven harness, or dependency-injected boundary stub), copy that exact setup. Do not invent a second testing style in the same package.
2. **Single-Caller Helper Ban (No Test-Induced Damage)**:
   * Count the production callers of any function/method you want to test directly.
   * If you split an internal step out of an already-testable module into a function that has **only `1` caller in production code** solely so your test file could call it directly -> **STOP. Inline the function back into its caller** and write the test against the caller's public seam.
   * **Exemption**: a **Brownfield Sprout** chosen via Agency Gate option (b) to isolate new logic from an untestable legacy monolith is a deliberate seam. It may have one caller and is tested directly; do not inline it.
3. **No Call-Only / Mock-Heavy Assertions**:
   * Inspect your test's assertions and setup.
   * If your test sets up **more than 2 mocks/stubs**, OR if its assertions **only check that an internal dependency was called** (`mock.calledWith`) without asserting a returned value, state transition, or error result -> **STOP. Delete the unit test.** Either test at a higher seam or verify via runtime execution.
   * The request a true external system observes counts as an outcome when the contract promises it; assert it with a paired positive case ([tests.md](tests.md#the-two-edit-test)).
4. **Real-Input & Tautology Guard**:
   * **No Tautological Expected Values**: Expected values in assertions must be concrete literals or spec-derived outcomes, never recomputed by copying the production formula into the test.
   * **No Normalized Bug Inputs**: When replaying a real failing input or production call, never strip out flags, paths, or preconditions just to make the fixture easier to construct. If removing your production fix does not turn the test red, the test is vacuous.
   * **Explain the Number**: Before writing a test or benchmark that asserts a performance, latency, throughput, or allocation target, use the `diagnosing-bugs` skill and read its `PERF.md` checklist.

## Step 3: The Red-Green Vertical Loop

When running Seam-Level TDD:

1. **Confirm the Seam**: If the spec or ticket already pre-agreed the test seam, use it. Otherwise propose the public function/interface seam you will test at, give each proposed seam one line on what it catches and what it misses, and **wait for the user to confirm it** before writing the first test.
2. **One Tracer-Bullet Slice at a Time**:
   * **Red**: Write **one** behavioral test at the agreed seam. Run the single test file and confirm it fails on the expected behavior (not a syntax or import error).
   * **Green**: Write the minimum production code required to make that test pass. Do not write code for future test cases yet.
   * Repeat one behavior at a time. Never write a batch of 10 failing tests upfront (horizontal slicing): bulk tests verify *imagined* behavior and lock in test structure before you understand the implementation. Each test is a **tracer bullet** that responds to what the last cycle taught you.
   * **Attack the Premise**: If two fixes or tests built on the same premise fail, stop patching and use the `diagnosing-bugs` skill to attack the premise before the next attempt.
3. **Refactoring Is Not Part of the Loop**: Do not restructure code during red → green. Structural refactoring belongs to the review stage (`code-review`) or a separate, explicitly agreed tidy step. For any structural cleanup beyond renaming a local variable or helper in the file you just touched, use the `tidy` skill so every transformation stays behavior-preserving, test-gated, and separate from behavioral commits.

## Step 4: Post-Green Test Pruning Gate (Mandatory Before Commit)

Modern coding agents over-generate redundant tests during the red-green loop. Once all behaviors are green, audit every test case added in this session and **delete**:

1. **Same-Branch Duplicates**: If two test cases execute the exact same conditional branch (`if / else / switch` path) using different normal values (e.g., testing `"alice"` and `"bob"` when neither is a boundary condition), delete the duplicate. Keep only:
   * `1` canonical happy-path case per behavior, plus
   * `1` case per distinct branching edge, boundary limit (`0`, empty, max), or error mode.
2. **Verification Budget (One Owner per Claim)**: Each behavioral claim gets exactly one verification owner, either one seam test or one runtime check. Delete a second test that re-proves the same invariant from an inner unit and an outer wrapper.
3. **Transitional Scaffolding Tests**: Delete tiny stepping-stone tests written mid-loop whose coverage is now 100% subsumed by a richer behavioral test at the same seam.
4. **Two-Edit Check**: Rewrite or delete any surviving test that still passes when the behavior is broken, or fails when only the internals change ([The Two-Edit Test](tests.md#the-two-edit-test)).
5. Re-run the test file after pruning to confirm everything stays green.

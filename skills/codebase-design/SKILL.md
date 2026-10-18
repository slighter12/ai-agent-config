---
name: codebase-design
description: Shared vocabulary and architectural guardrails for designing deep modules, agent-legible boundaries, brownfield sprout seams, and scrap circuit breakers. Use when designing or reviewing interfaces, deciding where a seam goes, refactoring legacy coupling, or making code testable and AI-navigable.
metadata:
  invocation: model
---

# Codebase Design (Deep Modules & Agent-First Architecture)

Design **deep modules** with high **local legibility**: a lot of behavior behind a small, deterministic interface placed at a clean seam, where a change that looks right from a single file is right for the whole repository.

## 1. Pre-Flight Context & User Agency Gate

Before proposing or restructuring an interface:
1. **Check Project Law**: Read `AGENTS.md` / `CLAUDE.md`, `GLOSSARY.md`, and `docs/adr/`. Project ADRs always override default design heuristics.
2. **Separate structural contracts from cosmetics**: Error types and error-return shape, state transitions, I/O seams, and public interfaces are **structural contracts** governed by Project Law (`AGENTS.md` / ADRs). Only cosmetic conventions (naming, formatting, file layout, import order) follow the local file.
3. **Evaluate the Brownfield vs. Greenfield Reality**: When an ADR specifies a new pattern while the target module uses an older one (**Active Migration Conflict**), or the target is a tightly coupled legacy function with no clean seam (**Legacy Monolith ("Shit Mountain")**), stop and present the canonical **Agency Gate Menu** once, with your recommendation:
   * **(a) Local minimal edit** — follow the legacy file's pattern for a minimal diff, verified by a runtime check.
   * **(b) Sprout at the seam** — new behavior in a clean function/module that follows the ADR, tested with Seam TDD, leaving a 1–3 line call site (Section 4).
   * **(c) Characterization check + local tidy first** — lock current behavior with a characterization test, tidy locally, then change.
   * **(d) Migrate** — move the touched module/function to the new ADR standard.

## 2. Core Vocabulary

Use these terms consistently across specs, reviews, and code discussions:

* **Module**: Anything with an interface and an implementation (a function, struct/class, package, or subsystem).
* **Interface**: Everything a caller must know to use the module correctly: signature, invariants, ordering constraints, and error modes.
* **Implementation**: The internal logic hidden behind the interface.
* **Depth (Leverage)**: The amount of behavior a caller or test exercises per unit of interface complexity they must learn.
  * **Deep Module**: Small, simple interface + rich internal implementation.
  * **Shallow Module (Avoid)**: Large or pass-through interface + thin implementation that leaks internal steps to callers.
* **Seam** (Michael Feathers): The place where behavior can be observed or swapped (e.g., substituting an in-memory fake or stub for an external I/O dependency) without editing the module's internal code.
* **Adapter**: A concrete implementation that satisfies an interface at a seam. One adapter means a hypothetical seam; two adapters (or production I/O + deterministic test fake) justify a real seam.

## 3. Agent-First Architecture Principles (`@poteto` Boundary Discipline)

1. **Assume a Low-Context Next Contributor (Local Legibility)**:
   * Assume the next contributor is an agent that opens only the target file and 1–2 neighbors, copies the nearest pattern, and takes the shortest path that compiles.
   * Eliminate ambiguous dual patterns inside the same file. Push invariants down the **5 Layers of Correction**: from prose instructions -> tests -> linter/compiler checks -> data structures/types so invalid states cannot compile.
2. **Deterministic Core / Thin Mechanical Shell**:
   * Place input validation, parsing, type narrowing, and I/O orchestration at the thin outer boundary (HTTP/CLI/RPC handlers, storage adapters).
   * Keep domain rules and state transitions inside a deterministic boundary that takes plain inputs (or narrow port interfaces at the I/O edge) and returns explicit results/errors.
   * **The Deletion Test**: Imagine deleting an intermediate layer or helper. If complexity disappears instead of reappearing across `N` callers, that layer was shallow pass-through boilerplate — inline it.

## 4. Brownfield Evolution: The Sprout Pattern

When adding or changing behavior inside a legacy file that lacks tests or mixes I/O with domain logic:

1. **Do Not Imitate Structural Debt**: Match the local file's cosmetic conventions (naming, formatting), but **never** copy its God-function coupling, global state mutation, or a structural contract (such as raw string errors) that an ADR has replaced. The Sprout follows the ADR's contract; translate at the 1–3 line call site if the legacy caller still needs the old shape.
2. **Sprout at the Seam**:
   * Write the new calculation, validation, or state-transition logic as a clean, self-contained function or small module (**the Sprout**) with deterministic inputs and outputs.
   * Modify the legacy monolith by only **1–3 lines** at the call site to invoke the Sprout.
   * Write behavioral tests directly against the Sprout's interface; verify the 1–3 line legacy wiring via runtime check (`Prove It Works`).
   * A Sprout usually has exactly one production caller. That is expected: it is a deliberate seam chosen at the Agency Gate, so it is **exempt** from the Single-Caller Helper Ban and is not a Test-Induced Damage signal.

## 5. Scrap & Re-Architect Circuit Breaker (熔斷退回信號)

Treat any of the following observable smells during design or implementation as a **Scrap Signal** — stop patching and redesign the module seam (or consult the user):

1. **Workaround Smell**: You have to add a boolean flag, dummy parameter, or state bypass to work around a constraint imposed by your own abstraction. (A workaround is proof of over-design or a misplaced seam.)
2. **Pass-Through Cascade (Shotgun Surgery)**: Adding one field or behavior forces mechanical signature edits across 3+ intermediate layers that do no transformation of their own.
3. **Test-Induced Damage**: Testing a behavior requires either (a) splitting a single-caller internal step out of an already-testable module just to import it in a test (Brownfield Sprouts are exempt), or (b) wiring >2 mocks/stubs that mirror internal call order.

## Going Deeper

* For dependency categories and replace-don't-layer seam refactoring, read `DEEPENING.md` in this skill directory.
* For comparing alternative interface designs on high-ambiguity modules, read `DESIGN-IT-TWICE.md` in this skill directory.

Use provider-native subagents when their isolation or parallelism justifies the coordination cost. Otherwise produce the same independent designs sequentially in the current session, resetting to the stated constraints for each design before comparing them.

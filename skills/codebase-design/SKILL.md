---
name: codebase-design
description: Shared vocabulary and architectural guardrails for designing deep, agent-legible modules, brownfield sprout seams, and scrap circuit breakers. Use when designing or reviewing interfaces, deciding where a seam goes, or refactoring legacy coupling.
metadata:
  invocation: model
---

# Codebase Design (Deep Modules & Agent-First Architecture)

Design **deep modules**: a lot of behaviour behind a small interface, placed at a clean seam, testable through that interface. Use this language and these principles wherever code is being designed or restructured. The aim is leverage for callers, locality for maintainers, and testability for everyone. Design for high **local legibility** too: a change that looks right from a single file is right for the whole repository.

## 1. Pre-Flight Context & User Agency Gate

Before proposing or restructuring an interface:
1. **Check Project Law**: Read `AGENTS.md` / `CLAUDE.md`, `GLOSSARY.md`, and `docs/adr/`. Project ADRs always override default design heuristics.
2. **Separate structural contracts from cosmetics**: Error types and error-return shape, state transitions, I/O seams, and public interfaces are **structural contracts** governed by Project Law (`AGENTS.md` / ADRs). Only cosmetic conventions (naming, formatting, file layout, import order) follow the local file.
3. **Evaluate the Brownfield vs. Greenfield Reality**: When an ADR specifies a new pattern while the target module uses an older one (**Active Migration Conflict**), or the target is a tightly coupled legacy function with no clean seam (**Legacy Monolith ("Shit Mountain")**), stop and present the canonical **Agency Gate Menu** once, with your recommendation. This is the canonical menu; other skills present it from here:
   * **(a) Local minimal edit** — follow the legacy file's pattern for a minimal diff, verified by a runtime check.
   * **(b) Sprout at the seam** — put the new behavior in a clean function/module that follows the ADR, test it with Seam TDD, and leave a 1–3 line call site in the legacy file.
   * **(c) Characterization check + local tidy first** — lock current behavior with a characterization test, tidy locally, then make the change.
   * **(d) Migrate** — move the touched module/function to the new ADR standard.

   Section 4 defines how to build the Sprout for option (b).

## 2. Core Vocabulary

Use these terms exactly: don't substitute "component," "service," "API," or "boundary." Consistent language is the whole point.

* **Module**: anything with an interface and an implementation. Deliberately scale-agnostic: a function, class, package, or tier-spanning slice. _Avoid_: unit, component, service.
* **Interface**: everything a caller must know to use the module correctly: the type signature, but also invariants, ordering constraints, error modes, required configuration, and performance characteristics. _Avoid_: API, signature (too narrow, they refer only to the type-level surface).
* **Implementation**: what's inside a module, its body of code. Distinct from **Adapter**: a thing can be a small adapter with a large implementation (a Postgres repo) or a large adapter with a small implementation (an in-memory fake). Reach for "adapter" when the seam is the topic; "implementation" otherwise.
* **Depth**: leverage at the interface. The amount of behaviour a caller (or test) can exercise per unit of interface they have to learn. A module is **deep** when a large amount of behaviour sits behind a small interface, **shallow** when the interface is nearly as complex as the implementation.
* **Seam** _(Michael Feathers)_: a place where you can alter behaviour without editing in that place; the *location* at which a module's interface lives. Where to put the seam is its own design decision, distinct from what goes behind it. _Avoid_: boundary (overloaded with DDD's bounded context).
* **Adapter**: a concrete thing that satisfies an interface at a seam. Describes *role* (what slot it fills), not substance (what's inside).
* **Leverage**: what callers get from depth. More capability per unit of interface they learn. One implementation pays back across N call sites and M tests.
* **Locality**: what maintainers get from depth. Change, bugs, knowledge, and verification concentrate in one place rather than spreading across callers. Fix once, fixed everywhere.

### Deep vs Shallow

**Deep module** = small interface + lots of implementation:

```
┌─────────────────────┐
│   Small Interface   │  ← Few methods, simple params
├─────────────────────┤
│                     │
│  Deep Implementation│  ← Complex logic hidden
│                     │
└─────────────────────┘
```

**Shallow module** = large interface + little implementation (avoid):

```
┌─────────────────────────────────┐
│       Large Interface           │  ← Many methods, complex params
├─────────────────────────────────┤
│  Thin Implementation            │  ← Just passes through
└─────────────────────────────────┘
```

When designing an interface, ask:

* Can I reduce the number of methods?
* Can I simplify the parameters?
* Can I hide more complexity inside?

### Principles

* **Depth is a property of the interface, not the implementation.** A deep module can be internally composed of small, mockable, swappable parts; they just aren't part of the interface. A module can have **internal seams** (private to its implementation, used by its own tests) as well as the **external seam** at its interface.
* **The interface is the test surface.** Callers and tests cross the same seam. If you want to test *past* the interface, the module is probably the wrong shape.
* **One adapter means a hypothetical seam. Two adapters means a real one.** Don't introduce a seam unless something actually varies across it. A production I/O adapter plus a deterministic test fake counts as two.
* The deletion test lives in §3 with the other subtraction rules.

### Relationships

* A **Module** has exactly one **Interface** (the surface it presents to callers and tests).
* **Depth** is a property of a **Module**, measured against its **Interface**.
* A **Seam** is where a **Module**'s **Interface** lives.
* An **Adapter** sits at a **Seam** and satisfies the **Interface**.
* **Depth** produces **Leverage** for callers and **Locality** for maintainers.

## 3. Agent-First Architecture Principles (after `@poteto` principle-boundary-discipline)

1. **Assume a Low-Context Next Contributor (Local Legibility)**:
   * Assume the next contributor is an agent that opens only the target file and 1–2 neighbors, copies the nearest pattern, and takes the shortest path that compiles.
   * Eliminate ambiguous dual patterns inside the same file. Push invariants down the **5 Layers of Correction**: from prose instructions -> tests -> linter/compiler checks -> data structures/types so invalid states cannot compile.
   * **The 30-Second Test**: A new reader must be able to answer "where does X come from?" and "what can change X?" in under 30 seconds. If not, cut layers or shrink state scope: return values over mutation, locals over fields, fields over module state, module state over globals; derive instead of syncing.
2. **Deterministic Core / Thin Mechanical Shell**:
   * Place input validation, parsing, type narrowing, and I/O orchestration at the thin outer edge (HTTP/CLI/RPC handlers, storage adapters).
   * Keep domain rules and state transitions inside a deterministic core that takes plain inputs (or narrow port interfaces at the I/O edge) and returns explicit results/errors.
   * **The Deletion Test**: Imagine deleting the module. If complexity vanishes, it was a pass-through. If complexity reappears across N callers, it was earning its keep. Apply it to every intermediate layer or helper, and inline the pass-throughs.
   * **Subtract Before You Add**: When evolving a system, remove complexity first, then build. Default to subtraction. Make simplification a continual investment. Leave the design slightly simpler and more capable behind the same or smaller surface than you found it.
     * Sequence removal before construction
     * Cut before you polish (get to the minimum before investing in quality)
     * Design for observed usage, not speculative edge cases
     * No speculative validators, parsers, or guards beyond what the spec demands
     * Simplify prompts (remove redundant instructions, excessive templates)
     * When a reference has no novel content, delete it rather than leaving a stub
     * Before adding a parameter, flag, or wrapper, delete the dead paths and collapse the single-caller wrappers it would touch, so the change lands on a smaller base.
     * Subtraction stops at meaning: do not inline away a name that carries intent, merge unrelated logic into one function, or remove a seam that serves testability.
   * **Laziness Protocol**:
     * **Consolidate decisions.** Do not repeat the same choice in several places. Put it behind one source of truth and pass the result as a simple flag.
     * **Minimize the diff.** Make the smallest change that solves the problem. Fewer lines beat "elegant" boilerplate.
3. **Caller Usage & Data Shape First**:
   * Before writing signatures, write the most common caller's usage snippet, then the core data shape (types, state machine, or table registry) that makes illegal states unrepresentable. Parse external data into that shape once at the edge.
   * Derive the signatures from the caller usage, not the reverse.
   * **Structure over branches**: The signs you skipped the data shape are a new feature that grows an existing if/else chain by one more branch, a second boolean that must stay in sync with the first, or modules named after phases. Reach for a state machine, discriminated union, or registry instead; keep boring code when the current shape is clear, local, and unlikely to grow.
   * **Brand semantic primitives**: When two values share a primitive but mean different things (`UserId`, `OrderId`), give each its own type (a branded type, newtype, or Go defined type), validated once at creation. When a schema (protobuf, OpenAPI, GraphQL, a migration) owns a shape, derive the type from it instead of hand-rolling a parallel one.
4. **Designing for Testability**: Accept dependencies instead of constructing them (`processOrder(order, paymentGateway)`, not `new StripeGateway()` inside), and return results instead of producing side effects (`calculateDiscount(cart): Discount`, not `applyDiscount(cart): void` mutating the cart). Fewer methods and params mean fewer tests and simpler setup.

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
4. **Type Escape Hatches**: The design compiles only with `any`, casts, or optional fields that are always set in practice.
5. **The Lock Reflex**: You reach for a lock or a "busy" flag on state the design said was unshared. Separate the state before serializing access to it.
6. **Leaked Ordering**: Callers must know the abstraction's internal rules or call ordering to use it correctly.
7. **Repeated Deviations**: Two or more independent deviations from the spec or sketch share the same shape during implementation.

One instance can be noise; a repeated pattern is the signal.

## Rejected Framings

* **Depth as ratio of implementation-lines to interface-lines** (Ousterhout): rewards padding the implementation. We use depth-as-leverage instead.
* **"Interface" as a language `interface` keyword or a class's public methods**: too narrow: interface here includes every fact a caller must know.
* **"Boundary"**: overloaded with DDD's bounded context. Say **seam**, **interface**, or **edge** (for the outer I/O side).

## Going Deeper

* Before deciding how a dependency crosses a seam (its category, and whether tests replace it or layer over it), read [DEEPENING.md](DEEPENING.md).
* Before proposing, approving, or reviewing an interface, screen it against [DESIGN-RED-FLAGS.md](DESIGN-RED-FLAGS.md).
* When a Scrap Signal fired, the interface is exported across packages or modules, or the user asks for alternative designs, read [DESIGN-IT-TWICE.md](DESIGN-IT-TWICE.md) and follow its process.

Use provider-native subagents when their isolation or parallelism justifies the coordination cost. Otherwise produce the same independent designs sequentially in the current session, resetting to the stated constraints for each design before comparing them.

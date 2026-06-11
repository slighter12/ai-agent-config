---
description: "Bounded architecture discovery heuristics for planning-grill"
---

# Architecture Discovery

Use this reference when `planning-grill` is triggered by architecture shape, module boundary, coupling, seam, locality, or refactoring-opportunity prompts.

## Purpose

Find architecture improvement candidates before design or implementation. The output is a short candidate list and follow-up questions, not a refactor, formal audit report, or issue artifact.

## Vocabulary

- **Module**: a conceptual unit with a reason to change and a boundary that callers can understand.
- **Interface**: the public contract callers rely on, including function signatures, request/response shapes, events, or file formats.
- **Implementation**: private details hidden behind the interface.
- **Seam**: a boundary where behavior can be isolated, tested, replaced, or adapted.
- **Adapter**: code that translates across boundaries, protocols, frameworks, or external systems.
- **Depth**: how much useful behavior a module hides behind a small interface.
- **Leverage**: how much simplification or safety a boundary creates for callers.
- **Locality**: how close related knowledge, invariants, and changes stay to each other.

## Discovery Workflow

1. Bound the scan scope first: feature area, module, flow, layer, or caller set.
2. Inspect relevant entry points, callers, callees, tests, configs, docs, and any existing ADR or domain notes.
3. Map current architecture shape in plain terms: boundaries, dependencies, seams, adapters, data flow, and ownership.
4. Look for friction:
   - understanding one concept requires jumping across many unrelated files;
   - callers need to know implementation details;
   - a module mostly passes data through without hiding behavior;
   - tests are forced through the wrong seam;
   - domain terms or invariants are duplicated in several places;
   - external protocols or framework details leak into core logic;
   - changes are consistently non-local or require synchronized edits.
5. Apply the deletion test: if removing a module or abstraction would make the code simpler without losing behavior, it may be shallow.
6. Produce 2-5 candidates, not a full redesign.
7. For each candidate, state:
   - affected files, modules, or boundaries;
   - observed friction;
   - candidate direction;
   - expected benefit;
   - risk, confidence, and why now / why not now.
8. Ask the user to choose a candidate or confirm priorities before detailed design or implementation.

## Guardrails

- Do not make broad architecture claims without inspected facts.
- Do not design a new interface before showing why the current boundary is causing friction.
- Do not refactor, rename public APIs, move files, or add dependencies.
- Do not create HTML reports, formal audits, PRDs, ADRs, issues, tickets, or context docs unless explicitly requested.
- Do not replace language-specific architecture policy. Use Go, Rust, frontend, API, infra, or security policy when those details are primary.
- If the request is about current diff coherence, use `code-review` instead.
- If the request asks to implement an approved refactor, use `implement-change`.

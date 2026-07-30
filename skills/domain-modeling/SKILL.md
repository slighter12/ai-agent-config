---
name: domain-modeling
description: Capture durable domain language, invariants, boundaries, and architectural decisions from a discussion or codebase. Use when domain understanding must survive beyond the current session. Avoid when the information is temporary implementation detail.
metadata:
  invocation: model
---

# Domain Modeling

Record only concepts that change how future work should be reasoned about:

- ubiquitous terms and precise meanings;
- actors, responsibilities, and trust boundaries;
- entities, value objects, states, and transitions;
- invariants and failure semantics;
- bounded contexts and integration edges;
- decisions with meaningful alternatives.

Prefer updating existing project domain docs under `docs/agents/domain/`. Separate observed fact, agreed decision, and open question. Use a short ADR when a decision constrains future architecture.

Complete only with durable information that changes future reasoning. Every update distinguishes observed facts, agreed decisions, and open questions; code inventories and diagrams belong only when they add that durable context.

---
name: prototype
description: Build a fast, throwaway prototype to answer an empirical design, state-model, or UI question that talking cannot settle.
metadata:
  invocation: model
---

# Prototype (Throwaway Empirical Spike)

Build throwaway code to answer a specific design, algorithmic, or UX question when verbal `grilling` cannot settle it.

## Choose the Branch
1. **Logic / State / Algorithm Question**: Build a minimal self-contained runnable script, benchmark, or interactive single-file HTML harness that exercises the state transitions and prints/renders the full state after every action.
2. **UI / Interaction Question**: Build 2–3 distinct switchable variations on a local route or standalone page so the user can compare them directly.

## Rules
* **Clearly Marked Throwaway**: Name files/routes explicitly as `prototype` and skip production polish, abstractions, and unit test suites.
* **No Production Persistence**: Keep state in memory or a disposable scratch file.
* **Capture Decision & Clean Up**: Once the user inspects the prototype and makes the decision, commit the prototype only to a throwaway reference branch (never `main`), record the settled decision and any state/type snippet in the spec or ADR, and clean the working tree.

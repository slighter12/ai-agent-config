---
name: how
description: Explain how a subsystem or code path in this repository works, with file:line references, so it can be changed safely. Use for walkthroughs, runtime flow, and "where does this live" questions.
metadata:
  invocation: model
---

# How (Code Walkthrough)

Explain how code in this repository works at the depth a senior engineer needs to start changing it: enough for a working mental model, not an annotated copy of the source. `why` explains motivation and history; `codebase-design` decides how code should be shaped; `improve-codebase-architecture` surveys the whole repository for deepening; `research` covers external and primary-source questions. `how` describes what exists now.

This skill is read-only: it never edits files, the index, history, or external systems, and it writes nothing to the repository.

## 1. Size the Question

- **Simple** (one module, one function, a narrow question): explore and answer in this session.
- **Complex** (a subsystem across several files or services, a cross-cutting feature, an architectural overview): split it into 2–4 angles, such as the entry point and trigger, the data flow, the I/O edges and persistence, and the error and concurrency paths. Explore the angles one at a time in this session. Dispatch read-only `explorer` subagents in parallel only when their raw search output would bloat this window, and ask each to return the explorer contract below.

When in doubt, treat it as simple. If the scope is ambiguous, state your reading and proceed; the user can redirect.

## 2. Trace

- Start at a real entry point (a route, CLI command, job, event handler, or exported function) and follow the calls.
- Cite `file:line` for every claim. Read the code, not just the names.
- Treat comments and docs as claims to check against the code.
- Name where the subsystem connects to the rest of the codebase: its inputs, outputs, seams, and adapters, in `codebase-design` vocabulary.
- Mark anything you could not trace as an open question instead of guessing.

**Explorer contract** (when subagents are used): Components Found (name, `file:line`, one sentence each), Flow (each step's function, file, what it does, and what it calls next, with the data passed), Files Read, Connections to the rest of the codebase, Non-Obvious Things, and Open Questions.

## 3. Explain

Use these sections, dropping any that do not apply:

- **Overview**: one or two paragraphs on what this is, what it does, and why it exists. A reader should be able to stop here and decide whether to keep reading.
- **Key Concepts**: the types, services, and abstractions needed to follow the rest, briefly defined.
- **How It Works**: the longest section. Walk the flow in prose: what triggers it, each step, where data goes, and the decision points. Add a Mermaid or ASCII diagram only when components exchange messages or data changes shape across stages and prose alone is harder to follow.
- **Where Things Live**: a short file and directory map of what someone needs to start working here.
- **Gotchas**: surprising behavior, pitfalls, and easy mistakes.
- **Open Questions**: what you could not trace.

Write with concrete actors ("`UserService` calls `AuthClient.refresh()`", not "the service delegates to the client"). When something is complex, explain why it is complex; when it is simple, keep it short. If the real question turns out to be why the code exists, use the `why` skill; if it is where new code should go, use the `codebase-design` skill.

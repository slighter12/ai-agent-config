---
name: how
description: Explain how a subsystem or code path in this repository works, with file:line references, so it can be changed safely. Use for walkthroughs, runtime flow, and questions about where a flow or subsystem lives across files.
metadata:
  invocation: model
---

# How (Code Walkthrough)

Explore the codebase to answer "how does X work?" questions. Produce architectural explanations at the level of a senior engineer onboarding onto a subsystem, enough to build a working mental model, not so much that it reads like annotated source code. `why` explains motivation and history; `codebase-design` decides how code should be shaped; `improve-codebase-architecture` surveys the whole repository for deepening; `research` covers external and primary-source questions. `how` describes what exists now.

This skill is read-only: it never edits files, the index, history, or external systems, and it writes nothing to the repository.

## 1. Assess complexity

If the scope is ambiguous, state your interpretation and explore. The user can redirect.

- **Simple** (a single module, a small utility, a narrow question such as "how does function X work"): explore and explain in one pass in this session. Go to step 4.
- **Complex** (a subsystem spanning multiple files or services, a cross-cutting feature, a full architectural overview): go to step 2.

Treat the question as Complex when the answer spans two or more top-level directories or services; otherwise take the simple path.

## Trace Rules (Both Paths)

- Start at a real entry point (a route, CLI command, job, event handler, or exported function) and follow the calls.
- Cite `file:line` for every claim. Read the code, not just the names.
- Treat comments and docs as claims to check against the code.
- Name where the subsystem connects to the rest of the codebase: its inputs, outputs, seams, and adapters, in `codebase-design` vocabulary.
- Mark anything you could not trace as an open question instead of guessing.

## 2. Explore (complex only)

Decompose the question into 2 to 4 angles, each a distinct slice of the subsystem, such as the entry point and trigger, the data flow, the I/O edges and persistence, and the error and concurrency paths. Explore the angles one at a time in this session. Dispatch read-only `explorer` subagents in parallel only when their raw search output would bloat this window, each with the explorer contract below and its angle filled in.

### Explorer contract

Don't try to cover everything. Focus on your assigned angle and go deep. Return:

- **Components Found**: The key types, services, classes, and abstractions. For each: name, file path, and a one-sentence description of what it does.
- **Flow**: The execution flow step by step. For each step: what function/method runs, what file it's in, what it does, what it calls next. Include the data that flows between steps.
- **Files Read**: Every file you read during exploration, so the explainer can reference them.
- **Connections**: Where this subsystem connects to other parts of the codebase. The inputs and outputs.
- **Non-Obvious Things**: Anything surprising, historically motivated, or easy to get wrong. Things that look like they should work one way but work another.
- **Open Questions**: Anything you couldn't fully trace or understand. Be honest about gaps.

## 3. Synthesize (complex only)

Once every angle has returned, reconcile the findings. They will overlap in places and may occasionally contradict. Merge overlapping descriptions, resolve contradictions by checking the code yourself, and combine the separate slices into a unified picture. When explorers did the work, you shouldn't need to re-explore from scratch.

## 4. Explain

Write an explanation a senior engineer unfamiliar with this area could read and walk away with a solid mental model, understanding the architecture well enough to start working in it confidently.

Use this structure, adapted to what makes sense for the question. Not every section is needed for every question.

### Overview

1-2 paragraphs. What is this thing, what does it do, why does it exist. Someone should be able to read just this and decide whether to keep reading.

### Key Concepts

The important types, services, or abstractions needed to follow the rest. Brief definitions, not exhaustive.

### How It Works

The core of the explanation, and the longest section. Walk through the flow: what triggers it, what happens step by step, where data goes, what the decision points are.

Use prose, not pseudocode. Reference specific files and functions so the reader knows where to look, but don't dump large code blocks unless a snippet is essential to a point.

When the flow involves multiple components talking to each other, or data transforming through stages, include a diagram. Use mermaid (```mermaid) for structured flows (sequence diagrams, flowcharts, component graphs) or ASCII art for simpler relationships where mermaid would be overkill. Use your judgment. A diagram should clarify, not decorate. If prose covers the flow, skip the diagram.

### Where Things Live

A brief file/directory map. Just the ones someone would need to start working here.

### Gotchas

Non-obvious things, surprising behavior, historical context, pitfalls. Skip this section if there's nothing worth calling out.

### Open Questions

What you could not trace. Skip this section when every claim traced to code.

### Communication Style

- Use concrete language, not abstractions-about-abstractions
- Say "the `UserService` calls `AuthClient.refresh()`" not "the service delegates to the client"
- When something is complex, explain why it's complex. Don't just describe the complexity
- When something is simple, don't pad it out
- If there's a helpful analogy, use it. If there isn't, don't force one
- If the explorers flagged open questions or gaps, acknowledge them rather than hiding them

If the real question turns out to be why the code exists, use the `why` skill; if it is where new code should go, use the `codebase-design` skill.

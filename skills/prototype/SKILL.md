---
name: prototype
description: Build a fast, throwaway prototype to answer an empirical question that talking cannot settle. Use when a design, state-model, or UI question needs a runnable answer before committing to it.
metadata:
  invocation: model
---

# Prototype

A prototype is **throwaway code that answers a question**. The question decides the shape.

When `grilling` tags a branch as a **Spike / Prototype Prerequisite**, that branch's question is the question: copy it verbatim to the top of the prototype.

## Pick a branch

Identify which question is being answered — from the user's prompt or the surrounding code. If the prompt names neither logic or state nor appearance, ask the user once:

- **"Does this logic / state model feel right?"** → read [LOGIC.md](LOGIC.md) and follow its process.
- **"What should this look like?"** → read [UI.md](UI.md) and follow its process.

The two branches produce very different artifacts — getting this wrong wastes the whole prototype. In a non-interactive run, or when the user does not answer, default to whichever branch better matches the surrounding code (a backend module → logic; a page or component → UI) and state the assumption at the top of the prototype.

## Rules that apply to both

1. **Throwaway from day one, and clearly marked as such.** Locate the prototype code close to where it will actually be used (next to the module or page it's prototyping for) so context is obvious — but name it so a casual reader can see it's a prototype, not production. For throwaway UI routes, obey whatever routing convention the project already uses; don't invent a new top-level structure.
2. **Trivial to run.** A UI prototype starts from one command in the project's task runner — `pnpm <name>`, `python <path>`, `bun <path>`, etc. A logic demo is a single HTML file the user double-clicks. Either way, no thinking required to start it.
3. **No persistence by default.** State lives in memory. Persistence is the thing the prototype is _checking_, not something it should depend on. If the question explicitly involves a database, hit a scratch DB or a local file with a clear "PROTOTYPE — wipe me" name.
4. **Skip the polish.** No tests, no error handling beyond what makes the prototype _runnable_, no abstractions. The point is to learn something fast.
5. **Surface the state.** After every action (logic) or on every variant switch (UI), print or render the full relevant state so the user can see what changed.
6. **Capture it when done.** Fold any validated decision into the real code, or back into the spec by telling the user to run `/to-spec` (or into an ADR in `docs/adr/` when documentation changes are authorized) when implementation has not started. When the user explicitly authorizes branch and commit creation, capture the prototype itself as a **primary source** on a throwaway branch out of main and leave a context pointer to it on the implementation issue. Without that authorization, leave the prototype uncommitted and report the proposed branch and context pointer. Capture the answer too — the verdict and the question it settled — in an authorized issue or commit. The main branch keeps only the validated decision.

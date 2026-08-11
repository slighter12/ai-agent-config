# Domain Docs

## Before exploring

- Read root `CONTEXT.md`, or root `CONTEXT-MAP.md` and every context relevant to the task.
- Read ADRs under `docs/adr/` that affect the area. In a multi-context repo, also inspect context-local `docs/adr/` directories.
- Proceed silently when these files do not exist. `domain-modeling` creates them lazily only when a term or decision is resolved and the request authorizes documentation changes.

## Layout

A single-context repo uses root `CONTEXT.md` and `docs/adr/`. A multi-context repo uses root `CONTEXT-MAP.md`, which points to each context's `CONTEXT.md`; system-wide ADRs remain under root `docs/adr/`, while context-specific ADRs live beside that context.

## Consumer rules

- Use the glossary's canonical terms and avoid listed synonyms.
- Reconsider invented terminology before adding it; use `domain-modeling` only for a real domain gap.
- Surface conflicts with existing ADRs instead of silently overriding them.

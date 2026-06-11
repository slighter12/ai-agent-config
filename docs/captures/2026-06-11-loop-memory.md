# Loop Memory Capture

Date: 2026-06-11
Status: accepted

## Source Material

- User-provided text from Addy Osmani's X Article on loop engineering, originally linked from
  <https://x.com/addyosmani/status/2064127981161959567>. The public X Article URL was identified
  as <https://x.com/i/article/2064122477731852288>, but the article body was available in this
  session through the user's pasted text rather than public unauthenticated fetch.
- User-provided screenshot table comparing loop primitives across Codex app and Claude Code:
  automations, worktrees, skills, plugins/connectors, sub-agents, and state.
- Superpowers reference repository: <https://github.com/obra/superpowers>. Relevant observed
  patterns: methodology is expressed as skills, plans/specs are stored as artifacts, worktrees are
  used for isolation, and skill creation is a gated workflow rather than a default destination for
  all memory.
- Local repository context: `project-lifecycle` already owns lifecycle capture, workflow lessons,
  capture-worthy handoff notes, status/doc sync, and skill evolution routing through
  `skill-creator`.

## Accepted Decisions

- Loop memory belongs in `project-lifecycle` for v1. It should not be implemented as a standalone
  `loop-engineering` skill yet.
- Use an `Active State + Capture` model:
  - active state is short-term continuity for the next loop run;
  - long-lived memory is classified through lifecycle capture;
  - raw transcripts and full run logs are not committed by default.
- Project-specific loop state should live in the operated repo or its chosen external source of
  truth. This shared config repo should provide reusable lifecycle guidance and templates, not
  centralize every project's operational memory by default.
- Manual capture comes first. Future automations may trigger capture candidates, but the initial
  workflow should teach agents to create source-grounded capture packets by hand.
- Approved shared skill changes still route to `skill-creator`. Lifecycle capture classifies the
  candidate first; it does not silently mutate skills.

## Rejected Alternatives

- Standalone `loop-engineering` skill for v1: rejected because the immediate need is memory and
  lifecycle capture, which overlaps with `project-lifecycle`.
- Permanent ledger of every loop run: rejected as the default because it risks noisy, stale, or
  privacy-sensitive memory. It may be added only if a user explicitly requests an audit trail.
- Directly promoting stable loop observations into shared skills: rejected because Superpowers and
  this repo both treat skill evolution as a gated authoring workflow.

## Memory Classification

- Active state: current phase, owner, blockers, verification status, next action, and any draft
  patch/review status needed for the next run.
- Long-lived capture: accepted decisions, implementation pivots, status changes, documentation
  drift, capture-worthy handoff notes, workflow lessons, and approved skill evolution candidates.
- Not captured by default: full transcripts, raw run logs, speculative ideas, duplicated source
  material, and unsettled active debate.

## Open Questions

- Whether a future dedicated loop-design skill is needed for automation prompt authoring after
  project-lifecycle covers memory capture.
- Whether target repos should standardize active state paths, such as `docs/loops/<name>/STATE.md`,
  or leave path choice project-specific.

## Follow-Up Owners

- `project-lifecycle`: own loop memory, discussion records, active state checkpoint classification,
  and capture packet guidance.
- `skill-creator`: own any approved shared skill changes after lifecycle capture classifies them.
- `execution-harness`: remain responsible for multi-agent or multi-phase coordination inside one
  active run, not long-lived memory capture.

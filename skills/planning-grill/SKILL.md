---
name: planning-grill
description: "Clarify ambiguous requirements, product goals, design intent, architecture tradeoffs, prototype decisions, roadmap/phase planning, bounded codebase architecture discovery, and implementation-ready plans before coding. Use when the user asks to discuss direction, compare options, challenge assumptions, validate a plan, decide whether a prototype is worth building, inspect architecture shape/coupling/seams, split an approved design into text-only slices, or turn fuzzy intent into goals, constraints, success criteria, non-goals, open questions, architecture candidates, phases, or a compact handoff. Avoid when the task is direct implementation, approved prototype building, diagnosis, current-diff review, test-primary work, pure testing strategy, git workflow, lifecycle capture, issue/ticket publishing, formal report writing, orchestration, or policy-only guidance."
license: MIT
compatibility: [codex, claude, gemini, opencode]
metadata:
  version: "0.1.12"
---

# Planning Grill

## Purpose

Turn unclear requests, product ideas, design questions, architecture-shape concerns, roadmap/tier/phase planning, and broad workflow choices into decision-ready guidance without forcing a heavy process. This skill borrows the useful part of challenge-style planning: inspect first, separate facts from preferences, then ask only questions that materially change the plan.

## Use When

- The user asks to discuss direction, confirm a plan, compare options, challenge assumptions, or ask planning questions before coding.
- The user asks for product goals, requirements, environment concepts, design intent, or architecture tradeoffs before implementation.
- The user asks to inspect architecture shape, scan callers/callees, assess module boundaries, seams, coupling, locality, or refactoring opportunities before coding.
- The user needs to compare architecture styles or domain modeling approaches before choosing implementation shape.
- The request contains fuzzy intent, competing approaches, unresolved tradeoffs, or missing non-goals.
- The user wants fuzzy product or feature intent organized into goals, constraints, success criteria, non-goals, and open questions.
- The user wants a roadmap, tier, phase, milestone, implementation order, or text-only implementation slice breakdown.
- The user wants bounded architecture discovery that lists improvement candidates without editing files or producing a formal audit artifact.
- The user wants a compact implementation-ready plan, work breakdown, or handoff as planning output without editing files or creating external artifacts.
- The user asks to split an approved design or selected architecture plan into implementation issues, but does not ask to create, publish, or manage issue tracker artifacts.
- The user needs to decide whether a small throwaway prototype is worth building before implementation.

## Avoid When

- The user asks to build an approved throwaway prototype; use `implement-change`.
- The user asks for direct implementation and the requirement is already clear.
- The user asks for diagnosis, debugging, root cause, flaky behavior, or performance investigation; use `diagnose`.
- The user asks to review current changes; use `code-review`.
- The user asks for commits, branches, pushes, PRs, or release flow; use `conventional-git-flow`.
- Tests or executable evidence are the main product; use `verification-driven-change`.
- The main uncertainty is testing strategy, verification depth, or unit/integration/e2e selection; use `policy-testing`.
- The user asks to create, format, publish, or manage GitHub/Linear issues, tickets, PRDs, ADRs, or repo context docs as artifacts.
- The user asks to capture, sync, or persist an accepted decision, status change, handoff, or workflow lesson; use `project-lifecycle`.
- The user asks for a formal architecture audit report, repo-wide inventory artifact, compliance-style assessment, or document deliverable.
- The user asks for multi-phase, multi-agent orchestration, phase gates, workspace coordination, or long-running handoff; use `execution-harness`.
- A specific language, API, frontend, infra, or security policy should own the decision.

## Workflow

1. Inspect discoverable facts before asking: relevant files, config, schemas, existing patterns, and any existing context docs, ADRs, or domain notes. If those docs do not exist, continue without creating them.
2. State the current understanding as goal, success criteria, constraints, and likely non-goals.
3. If the request is architecture discovery, bound the scan scope and use `references/ARCHITECTURE_DISCOVERY.md` to inspect modules, callers, seams, coupling, tests, and locality. List candidates before proposing a design.
4. Treat DDD, Clean Architecture, modular monolith, feature-based layout, layered CRUD, and similar patterns as context-dependent style choices, not defaults. Compare them against project size, domain complexity, existing code shape, team conventions, and expected change pressure before recommending one.
5. Identify unresolved decisions and their dependencies. Walk the decision tree one branch at a time instead of asking unrelated questions in bulk.
6. Clarify overloaded or conflicting domain terms against inspected docs and code. When terminology or decisions stabilize, emit a `project-lifecycle` capture candidate instead of writing context docs or ADRs from planning.
7. Ask the smallest useful number of questions, preferably one at a time, with a recommended default.
8. If a question needs runnable evidence rather than more discussion, define the smallest throwaway prototype question, expected learning, and scope. Keep the prototype decision in planning; route approved prototype implementation to `implement-change` instead of creating the artifact here.
9. Convert answers into a compact decision recommendation, planning summary, architecture candidate list, roadmap/tier/phase breakdown, implementation-ready plan, text-only slices, or handoff when requested.
10. When planning reaches an accepted project-level decision, deferred scope, status change, or handoff-worthy outcome, emit a `project-lifecycle` capture candidate instead of silently relying on the user to remember docs sync.
11. Route follow-up work to narrower skills or agents when they clearly own the next step. If a decision is complete and should be captured, route to `project-lifecycle` instead of writing docs from planning.

## Tool And Side-Effect Boundaries

- Prefer read-only exploration until the user asks for implementation.
- Do not run broad tests, formatters, migrations, installers, or external-service commands just to clarify a plan.
- Do not create or modify GitHub/Linear issues, tickets, PRDs, ADRs, context docs, domain docs, or repo docs unless explicitly asked.
- Even when explicit persistence is requested, prefer routing accepted decision/status documentation to `project-lifecycle` so planning remains decision clarification rather than lifecycle capture.
- Do not turn architecture discovery into implementation, formal report writing, repo-wide inventory, or issue creation.
- Do not build prototypes in this skill. An approved throwaway prototype is implementation work owned by `implement-change`.
- Do not use this skill as a substitute for `execution-harness`; use it to clarify what should happen before a heavier workflow begins.

## Output

Return:

- `understanding`: goal, constraints, and success criteria.
- `facts`: repository or source facts already inspected.
- `decisions`: unresolved choices and recommended defaults.
- `architecture_candidates`: bounded improvement candidates, when architecture discovery was requested.
- `lifecycle_capture_candidate`: accepted decision, deferred scope, status update, or handoff-worthy outcome for `project-lifecycle`, or `none`.
- `plan`: decision-ready recommendations, roadmap tiers, phases, implementation-ready next steps, text-only slices, or compact handoff when enough information is available.
- `routing`: narrower skill or agent handoffs, if any.
- `manual_verification`: checklist for validating the plan or design.

## Version History

- v0.1.0 (2026-05-15): Initial lightweight design and requirement clarification workflow.
- v0.1.1 (2026-05-15): Quote frontmatter description so YAML parsers accept routing keywords with colons.
- v0.1.2 (2026-05-18): Tighten routing around design and requirement clarification rather than implementation or diagnosis.
- v0.1.3 (2026-05-19): Keep planning triggers broad while tightening artifact, issue, testing, orchestration, and architecture-audit boundaries.
- v0.1.4 (2026-05-21): Add bounded architecture discovery mode and reference routing for architecture-shape prompts.
- v0.1.5 (2026-05-21): Route current-change review boundaries to renamed `code-review`.
- v0.1.6 (2026-05-21): Add roadmap, tier, phase, and text-only implementation slice routing while preserving issue artifact boundaries.
- v0.1.7 (2026-05-22): Route accepted decision-document synchronization to the then-separate decision-doc sync workflow.
- v0.1.8 (2026-05-28): Emit project lifecycle capture candidates for accepted decisions and planning outcomes.
- v0.1.9 (2026-06-29): Add decision-tree, domain-term sharpening, and prototype-decision guidance while keeping approved prototype artifacts in `implement-change`.
- v0.1.10 (2026-06-30): Add context-dependent architecture style selection guidance without defaulting to DDD or Clean Architecture.
- v0.1.11 (2026-06-30): Remove retired lifecycle skill id from history wording.
- v0.1.12 (2026-07-13): Declare OpenCode compatibility for the provider-neutral planning workflow.

## References

- `references/ARCHITECTURE_DISCOVERY.md` - Bounded architecture-shape discovery and candidate-listing heuristics.

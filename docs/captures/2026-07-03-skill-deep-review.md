# Skill Deep-Review Closeout Capture

Date: 2026-07-03
Status: accepted

## Source Material

- Five-category deep review of all 20 shared skills and their reference files, run across one
  Claude Code session series ending 2026-07-03: task-type, coordination/lifecycle, policy, domain,
  and meta categories.
- Fix commits on `main` dated 2026-07-03; per-skill change details are recorded in git history and
  each skill's `Version History`, and are intentionally not duplicated here.
- `README.md` routing note added 2026-07-03 explaining that harness activation does not imply
  delegated orchestration.

## Accepted Decisions

- The review series is complete. All material and should-fix findings across the five categories
  were fixed and re-verified in this series; no open must-fix items remain.
- Fix details live in git history and skill `Version History` entries; this capture records only
  deferred decisions and drift that the repo would not otherwise remember.

## Deferred Decisions

Each item below is a deliberate defer with rationale, not an unnoticed gap. Revisit only when its
trigger fires.

- `execution-harness` delegated-orchestrator trigger wording keeps the broad phrase
  "orchestrator, team, harness, or multi-agent planning" in `ACTIVATION_MODEL.md` and
  `agents/orchestrator.md`, despite tension with the "harness activation is not automatic
  delegation" anti-pattern. Rationale: sole-user setup does not use vague harness wording to force
  delegation, and `README.md` documents the design intent for human readers.
  Trigger to revisit: an observed misroute where asking for a harness spawns the orchestrator
  role; then narrow the wording to "delegated orchestrator or team-lead role".
- `goal-context` remains manual-trigger by prose only; the Claude surface does not apply the
  `disable-model-invocation: true` overlay that `skill-creator` overlay rules would suggest.
  Rationale: description-level guard has been sufficient so far.
  Trigger to revisit: any observed auto-invocation of `goal-context` without an explicit user
  request.
- The first-screen rule ("build the usable workflow as the first screen") intentionally exists in
  both `policy-frontend` (product baseline) and `design-art-direction` (design rule). It is the
  one visual-adjacent rule excluded from the 2026-07-03 ownership split that moved visual detail
  to `design-art-direction`. Keep both copies in sync when either changes.
- Prose-first output style is authoring guidance only (`skill-creator` cross-provider defaults and
  `output-patterns.md`); no runtime baseline rule was added to `AGENTS.md`/`CLAUDE.md`.
  Rationale: scope of that patch was authoring standards; runtime skills' `Output` field lists
  have not been observed producing YAML-like dumps in practice.
  Trigger to revisit: runtime responses from task skills start presenting field dumps instead of
  prose; then add one baseline line to `AGENTS.md`/`CLAUDE.md`.

## Documentation Drift

- `docs/skill-routing-eval.md` L09 (line ~368) still expects
  "Report `lifecycle_capture_candidate: none`", but `implement-change` v0.1.11 now omits lifecycle
  capture output when no signal exists. The probe's expected outcome needs rewording.
- `docs/skill-routing-eval.md` PC06 (line ~263) states "Rust policy says run tests by default" as
  the conflict premise, but `RUST_TESTING.md` v0.2.0 aligned Rust execution defaults with the
  global no-execution baseline. The probe premise is stale; either mark it as a hypothetical
  conflict or rewrite it around a real current precedence tension.
- `docs/skill-guardrail-audit.md` was checked and remains consistent with current
  `AUTHORING_STANDARDS.md` principles; no drift found.

## Memory Classification

- Active state: none needed; the review series is closed with no in-flight phase.
- Long-lived capture: the deferred decisions and documentation drift items above.
- Not captured: per-skill fix details (owned by git history and Version History), review process
  transcripts, and workflow lessons already legislated into
  `skill-creator/references/AUTHORING_STANDARDS.md` (guardrails-versus-policy) and
  `SKILL_EVOLUTION_GATE.md` (self-contained packaging).

## Open Questions

- None blocking. The two `skill-routing-eval.md` probe updates are mechanical follow-ups awaiting
  an owner pass.

## Follow-Up Owners

- `implement-change`: apply the two `docs/skill-routing-eval.md` probe rewordings (L09, PC06) when
  approved.
- `project-lifecycle`: re-check the four deferred decisions if any listed trigger fires.

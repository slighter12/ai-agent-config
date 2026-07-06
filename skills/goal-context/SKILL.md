---
name: goal-context
description: Audit, create, or repair a stable GOAL.md-style Goal Brief, then return one copy-ready goal_launch_context for starting goal mode. Use when the user explicitly invokes goal-context while drafting, fixing, or auditing a goal/handoff document for agent runtime goal use, especially when context and goal were mixed. Avoid when the user wants execution steps, implementation now, ordinary planning, lifecycle capture, code review, git workflow, diagnosis, or background memory capture.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.2.17"
---

# Goal Context

## Purpose

Audit, write, revise, or review goal handoff material so the user ends with two clear artifacts:

- `GOAL.md` or equivalent Goal Brief: a stable document that can live in a repo or be pasted as the
  target brief.
- `goal_launch_context`: one copy-ready block the user can paste into a fresh session or goal mode
  to start from the brief.

The Goal Brief separates two information sets:

- `Context`: background, current state, sources, constraints, non-goals, and open questions.
- `Goal`: objective, deliverables, metric type, evaluator contract, acceptance criteria,
  verification plan, stop or escalate rule, and done condition.

The launch context is not a second goal brief. It may contain a compact `/goal` command, but the user
should not have to understand separate `runtime_goal_command` and `next_session_context_prompt`
fields. Treat those as internal concepts and return a single `goal_launch_context` unless the user
explicitly asks for separate fields.

Goal mode has prompt-length pressure. Keep durable background, full acceptance criteria, source
lists, and stable constraints in the Goal Brief. Keep the launch context short: brief path, transient
current state, the goal-mode activation command, acceptance gaps, and completion evidence to report.
The launch context is also a runtime completion contract: it must name a checkable completion
condition, require verification evidence to be reported in the conversation, and include a bounded
stop condition instead of relying only on "see the brief."
When validation lineage includes runtime gates, conversation-surfaced evidence is required but is
not sufficient by itself; it must include or point to durable runtime evidence such as exact
commands, exit codes, an attempt ledger, and log paths or embedded output excerpts.

`Validation lineage evidence` means checked sources document a replayable gate chain, acceptance
protocol, result document, checkpoint or artifact dependency, or prior final label. When that
evidence exists, keep durable lineage in the Goal Brief and keep the launch context limited to the
current-stage attempt, transient gap, and evidence to report.

This is a manual-trigger skill. Do not use it only because the conversation mentions goals,
planning, context, or Markdown. Use it only when the user explicitly asks to invoke this skill.

This skill shapes the Goal Brief and launch context. It does not execute the goal, activate a
runtime goal tool, prescribe implementation steps, perform lifecycle capture, commit, or implement
the requested work. Lineage rerun semantics are a completion contract to write into the brief or
launcher, not work this skill performs in the current session.

When the selected output includes a Goal Brief, the brief audit is mandatory and must be
source-grounded. The launch context is downstream of the audit and must not substitute for fixing
weak, stale, or unsupported brief content.

This skill follows the provider-neutral loop contract vocabulary from `policy-core`: Goal Briefs
must define a checkable completion contract without assuming that Codex, Claude, Gemini, or another
runtime shares the same evaluator, memory, isolation, or tool access. Codex goal mode can continue a
persisted thread goal and the model can mark it complete; Claude goal mode can evaluate surfaced
conversation evidence. Neither behavior removes the need for a source-grounded completion contract.

## Use When

- The user explicitly invokes `goal-context`, `$goal-context`, `/goal-context`, or asks to run this
  manual skill.
- After explicit invocation, the user wants a goal brief, goal document, or goal handoff Markdown
  format drafted, revised, or reviewed.
- After explicit invocation, the user says prior goal documents mixed up context and goal
  information.
- After explicit invocation, the user wants acceptance criteria rewritten so they are observable,
  verifiable, and useful for an agent or reviewer.
- After explicit invocation, the user needs a copy-ready context to start Codex goal mode, Claude
  Code, or another agent runtime.
- After explicit invocation, the user needs a fresh session to continue partially completed work
  without relying on previous chat history.

## Avoid When

- The user wants execution steps, implementation, testing, review, commit, or diagnosis in the
  current session instead of receiving a Goal Brief and launch context.
- The user wants broad requirements, architecture, or roadmap planning; use `planning-grill`.
- The user wants to record a completed decision, phase closeout, or project memory update; use
  `project-lifecycle`.
- The user asks for ordinary documentation edits unrelated to goal brief format; use
  `implement-change`.
- The user asks for code review, git operations, diagnosis, or test-primary work.
- The task would require inventing objectives, owners, status, deadlines, evidence, or acceptance
  criteria.

## Loop Contract Minimums

Every Goal Brief or launch context for long-running agent use must include or point to these
contract parts:

- `Claim Boundary`: `Accepted Scope`, `Non-Claims`, and `Rejected Evidence`.
- `Metric Type`: one of `deterministic_validator`, `frozen_metric`, `evidence_review`,
  `subjective_rubric`, or `human_review_required`.
- `Evaluator Contract`: what evidence can prove the goal and where that evidence must be surfaced.
- `Stop Or Escalate`: retry, turn, blocker, tool-failure, budget, and human-review stop rules.
- `Completion Audit`: the final check must return to authoritative current state, not a progress
  ledger, transcript memory, or model confidence.

## Memory Hygiene

- Treat memory as evidence to inspect, not context to trust. Repo state, current user instructions,
  and checked artifacts override remembered summaries.
- Keep run records and decision records traceable, but do not auto-inject them into launch context
  unless they are needed for the current goal's completion contract.
- Persistent memory candidates must name provenance, creation date, scope, revocation path, and why
  they are stable across sessions. One-off preferences, failed hypotheses, stale diagnoses, and
  temporary implementation details do not qualify.
- When memory could bias repeated judgment, preserve a review packet or evidence pointer instead of
  turning the memory into default context.

If a deterministic validator or frozen metric exists, make it the primary completion evidence. If no
quantified, replayable, or calibrated standard exists, the goal must stop at a human-review packet
and must not allow the agent to self-certify completion. LLM judgment is supporting evidence only
unless the user explicitly accepts LLM-as-judge.

## Workflow

1. Confirm the manual invocation. If this skill was loaded without an explicit user request to run
   `goal-context`, stop and ask the user to invoke it manually.
2. Determine the requested behavior:
   - `audit-only`: use when the user asks to confirm, review, inspect, or check a brief. This means
     no file edits unless the user also asks to fix or rewrite them; it still returns
     `goal_launch_context` unless the user explicitly asks for `brief-only`.
   - `brief-and-launch`: default when the user asks to write, repair, rewrite, prepare, or finish the
     handoff. Audit the brief, revise the target Markdown when needed, then return `goal_launch_context`.
   - `launch-only`: use when the Goal Brief is already stable and the user only needs copy-ready
     startup context.
   - `brief-only`: use only when the user explicitly asks for just the stable goal document.
3. Resolve the target Goal Brief path before reading or editing:
   - Use the exact path named by the user.
   - Otherwise, if a repo root is discoverable and root `GOAL.md` already exists, use that file.
   - Otherwise, ask for the target path or inline-output format instead of guessing.
   - Do not create a new root `GOAL.md` by default.
4. Read the user's rough request, existing draft, or referenced file. Inspect only the sources needed
   to preserve factual context and avoid contradicting existing project state.
5. Run the Goal Readiness Gate before drafting or repairing launch text:
   - If roadmap, product direction, architecture tradeoff, or success criteria are not settled, stop
     and route to `planning-grill` with the smallest set of questions or handoff facts needed to
     converge direction. Do not write a Goal Brief that invents direction.
   - If direction is settled but scope boundary, evidence, metric type, evaluator contract, or stop
     rule is missing, ask only material questions needed to make the goal evaluable. Do not impose a
     fixed question count.
   - If no deterministic validator, frozen metric, evidence review criteria, or calibrated rubric
     exists, set `Metric Type` to `human_review_required` and make the done condition produce a
     review packet and stop for human judgment.
6. Read `references/GOAL_CONTEXT_TEMPLATE.md` before drafting, auditing, revising, or launching. It
   is the single source of truth for section rules, audit checks, acceptance criteria rules, launch
   context shape, response shape, and self-review checks.
7. Sort every useful fact into `Goal Brief` or `goal_launch_context` using the template's section
   rules. Put stable context and durable validation lineage in the brief; put only transient startup
   state, acceptance gaps, current-stage-first hints, and completion evidence requirements in the
   launcher.
8. Audit the Goal Brief when the selected behavior includes a brief. Use the template's Goal Brief
   Audit Rules as the checklist and do not treat `format: pass` as source-grounded content audit.
9. Rewrite the Goal Brief using the template when the selected behavior includes edits and audit
   findings require changes. Keep background tight and move future behavior out of `Context`.
10. If required information is missing, record it under `Open Questions` instead of inventing facts,
   behavior, evidence, or success conditions.
11. Generate `goal_launch_context` for every behavior except explicit `brief-only`. Use the
    template's Goal Launch Context Rules and default structure; do not maintain a second launch shape
    in this file. Reference the resolved Goal Brief path in the launcher.
12. Apply the Launch Output Gate before final response:
    - If behavior is not explicit `brief-only`, include `## Goal Launch Context`.
    - If a valid launcher can be created, include exactly one fenced `md` block whose first line
      starts with `/goal`.
    - If required facts are missing, include `goal_launch_context: blocked` with the missing facts
      instead of silently omitting the section.
13. Format the final response with the template's Response Shape Rules. Do not expose
    `runtime_goal_command` and `next_session_context_prompt` as separate top-level outputs unless the
    user explicitly asks for that lower-level split.
14. Stop after the goal document or review and launch context are complete. Do not proceed into
    implementation.

## Tool And Side-Effect Boundaries

- This skill may create or edit Markdown goal brief files only after the user explicitly requests
  this manual skill and the target path or requested output format is clear.
- Do not modify code, configs, tests, git state, provider settings, or external systems.
- Do not execute the goal, activate `/goal`, create persistent runtime goals, or perform lifecycle
  capture from this skill.
- Do not embed transient startup prompts, session status, `/goal` startup commands, or handoff
  chatter into the Goal Brief unless the information is stable project context that belongs in
  `Context`.
- Do not treat `/goal` text as global memory. It is a thread/session completion contract and should
  point back to the stable Goal Brief when background context is needed.
- Do not assume provider memory, goal evaluators, context isolation, or tool access unless the
  runtime explicitly provides it. Use repo artifacts, Goal Briefs, handoff references, and surfaced
  evidence for cross-session continuity.
- Do not auto-promote logs, transcripts, active-state notes, or memories into launch context. Include
  only the smallest sourced state needed to evaluate the current goal.
- Do not let the launch context become a background essay or duplicate brief; keep background, full
  acceptance criteria, source lists, and detailed verification plans in the Goal Brief.
- Do not copy durable validation lineage, full protocol commands, downstream commands, acceptance
  rules, or fallback rules into launch context. Reference the Goal Brief and include at most one
  documented command needed to regenerate the current stage's missing input.
- Do not make producing `goal_launch_context`, `runtime_goal_command`, or
  `next_session_context_prompt` a project acceptance criterion unless the user is explicitly
  authoring or testing this skill.
- Do not let launch output mask missing Goal Brief revisions; if the brief audit finds issues,
  revise the brief or report the blocker.
- Do not infer a brief is correct because it is well-structured; verify major claims against
  checked sources or repo state.
- Do not make the launch context a step-by-step implementation plan unless the user explicitly asks
  for a procedure-oriented handoff.
- Do not run programs unless required to inspect an existing goal document or the user explicitly
  asks.
- Do not overwrite an existing Goal Brief blindly; read it first and preserve still-valid facts.
- Keep generated goal documents free of secrets, credentials, private tokens, and unnecessary
  personal data.

## Output

Return a concise Markdown response. Use this shape unless the user explicitly asks for a different
format:

- `Summary`: whether the Goal Brief was drafted, revised, reviewed, unchanged, or blocked, with exact
  Markdown paths changed if any.
- `Brief Audit`: `format`, compact bullet `source_checks`, `brief_audit`, `brief_changes`, and
  `claim_boundary_review`, `metric_type_review`, `acceptance_criteria_review`, and
  `completion_contract_review`.
- `Goal Launch Context`: required for every non-`brief-only` response. Include one fenced `md` block
  named `goal_launch_context`, or `goal_launch_context: blocked` with missing facts when a valid
  launcher cannot be created.
- `Verification`: checks run, assumptions, open questions, and manual verification checklist.

Do not return a long YAML-like list of every internal field. Do not expose `runtime_goal_command` and
`next_session_context_prompt` as separate top-level outputs unless the user explicitly requests that
lower-level split. Use `references/GOAL_CONTEXT_TEMPLATE.md` for the detailed response and
`goal_launch_context` shapes.

## Version History

- v0.2.17 (2026-07-06): Add deterministic target Goal Brief path resolution so the skill asks
  instead of guessing when no explicit path or existing root `GOAL.md` is available.
- v0.2.16 (2026-07-06): Require `goal_launch_context` for all non-brief-only outputs and add a
  launch output gate so audit-only responses do not omit the launcher.
- v0.2.15 (2026-07-03): Move duplicated audit, launch, response, and self-review detail to the
  template reference as the single runtime source.
- v0.2.14 (2026-06-30): Add memory hygiene rules that separate records from automatically injected context.
- v0.2.13 (2026-06-21): Adopt provider-neutral loop contract rules with Goal Readiness Gate,
  mandatory claim boundaries, metric type classification, evaluator contracts, authoritative
  completion audits, and human-review stops when no objective or calibrated standard exists.
- v0.2.12 (2026-06-13): Require durable runtime evidence for runtime-gated validation lineage,
  including exact commands, exit codes, reached/skipped attempt ledgers, and log paths or embedded
  output excerpts instead of transcript-only claims.
- v0.2.11 (2026-06-13): Add runtime completion contract requirements to launch output: a
  falsifiable runtime check, conversation-surfaced evidence, and bounded retry or turn stop so goal
  mode does not over-trust "see the brief" or continue unbounded after a failed check.
- v0.2.10 (2026-06-13): Add replayable validation lineage, current-stage-first rerun semantics, and
  regenerable artifact handling for projects with documented gate chains or result lineage.
- v0.2.9 (2026-06-10): Require prior-plan acceptance evidence requirements to be preserved when
  converting plans into Goal Briefs.
- v0.2.8 (2026-06-10): Use the Goal Brief Objective as the launch objective summary while pointing
  to the brief as acceptance source, avoiding both `Read GOAL.md` and overly meta "goal defined in"
  wording.
- v0.2.7 (2026-06-10): Make launch commands delegate the objective to the Goal Brief instead of
  renaming the goal into a narrower audit, review, or verification task.
- v0.2.6 (2026-06-10): Align launch wording with `/goal <objective>` by replacing the default
  `/goal Read` shape with completion-objective wording that points to the Goal Brief as the
  acceptance brief.
- v0.2.5 (2026-06-10): Require fenced launch contexts whose first line is `/goal`, switch source
  checks to readable bullets, and prevent duplicate launch-context warning text.
- v0.2.4 (2026-06-10): Collapse runtime command and next-session prompt into one user-facing
  `goal_launch_context`, clarify audit-only versus rewrite behavior, and require concise Markdown
  response formatting instead of long field dumps.
- v0.2.3 (2026-06-10): Require source-grounded Goal Brief audits with `source_checks` before
  claiming pass or `none_needed`.
- v0.2.2 (2026-06-10): Require Goal Brief content audit and visible brief changes before runtime
  handoff outputs when the selected mode includes a brief.
- v0.2.1 (2026-06-10): Add prompt-budget rules that keep `/goal` compact, keep durable detail in
  the Goal Brief, and prevent next-session prompts from duplicating the brief.
- v0.2.0 (2026-06-10): Recenter the skill on runtime `/goal` completion contracts with explicit
  verification surface, boundaries, iteration policy, blocked stop condition, and surfaced evidence.
- v0.1.4 (2026-06-10): Reframe next-session prompts around required outcomes, acceptance gaps, and
  evidence instead of execution steps or implementation procedures.
- v0.1.3 (2026-06-09): Clarify dual artifacts, output modes, next-session prompt output, and
  self-review checks for fresh-session handoff quality.
- v0.1.2 (2026-06-09): Require a copy-ready next-session context prompt while keeping transient
  session context out of the goal brief.
- v0.1.1 (2026-06-09): Reframe as a manual goal document format skill with separated Context and
  Goal sections plus acceptance criteria quality rules.
- v0.1.0 (2026-06-09): Initial manual-trigger skill for goal-specific Markdown context creation.

## References

- `references/GOAL_CONTEXT_TEMPLATE.md` - Markdown structure, launch context shape, and acceptance
  criteria rules for Goal Brief documents.
- `references/EVAL.md` - Authoring-only prompt checklist for changing this skill.

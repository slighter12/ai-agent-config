---
name: goal-context
description: Audit, create, or repair a stable GOAL.md-style Goal Brief, then return one copy-ready goal_launch_context for starting goal mode. Use when the user explicitly invokes goal-context while drafting, fixing, or auditing a goal/handoff document for agent runtime goal use, especially when context and goal were mixed. Avoid when the user wants execution steps, implementation now, ordinary planning, lifecycle capture, code review, git workflow, diagnosis, or background memory capture.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.2.14"
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
   - `audit-only`: use when the user asks to confirm, review, inspect, or check a brief. Do not edit
     files unless the user also asks to fix or rewrite them.
   - `brief-and-launch`: default when the user asks to write, repair, rewrite, prepare, or finish the
     handoff. Audit the brief, revise the target Markdown when needed, then return `goal_launch_context`.
   - `launch-only`: use when the Goal Brief is already stable and the user only needs copy-ready
     startup context.
   - `brief-only`: use only when the user explicitly asks for just the stable goal document.
3. Read the user's rough request, existing draft, or referenced file. Inspect only the sources needed
   to preserve factual context and avoid contradicting existing project state.
4. Run the Goal Readiness Gate before drafting or repairing launch text:
   - If roadmap, product direction, architecture tradeoff, or success criteria are not settled, stop
     and route to `planning-grill` with the smallest set of questions or handoff facts needed to
     converge direction. Do not write a Goal Brief that invents direction.
   - If direction is settled but scope boundary, evidence, metric type, evaluator contract, or stop
     rule is missing, ask only material questions needed to make the goal evaluable. Do not impose a
     fixed question count.
   - If no deterministic validator, frozen metric, evidence review criteria, or calibrated rubric
     exists, set `Metric Type` to `human_review_required` and make the done condition produce a
     review packet and stop for human judgment.
5. Sort every useful fact into `Goal Brief` or `Launch Context`:
   - Put stable past and present facts in `Context`: why this exists, what is true now, exact
     sources, constraints, non-goals, and open questions.
   - Put `Claim Boundary` in `Context`: accepted scope, non-claims, and rejected evidence. Rejected
     evidence must name why prior or tempting evidence cannot support the accepted claim.
   - Put durable validation lineage in `Context` only when validation lineage evidence exists.
     Capture current stage, direct predecessor, relevant labels, required inputs, reproduction path,
     and evidence sources using the template.
   - Put desired future outcomes in `Goal`: objective, deliverables, metric type, evaluator
     contract, acceptance criteria, verification plan, stop or escalate rule, and done condition.
   - Put transient startup facts in `goal_launch_context`: current worktree/session status, partial
     acceptance evidence, acceptance gaps, current-stage-first attempt, necessary predecessor
     fallback, the compact `/goal` activation command, and completion evidence requirements.
6. Audit the Goal Brief when the selected behavior includes a brief. Treat this as a source-grounded
   content gate, not a formatting check:
   - Inspect the referenced files, repo state, or conversation facts needed to verify major claims.
   - `Context`: factual, stable, sourced, and free of future deliverables.
   - `Claim Boundary`: includes `Accepted Scope`, `Non-Claims`, and `Rejected Evidence`; the
     accepted scope is not broader than the evidence supports.
   - `Validation Lineage`: required only when validation lineage evidence exists. Audit it against
     the template fields and treat runtime-affecting omissions as `revised` or `blocked`.
   - `Goal`: future outcome only, not background narrative.
   - `Deliverables`: actual project artifacts, behavior, or decisions, not handoff metadata.
   - `Metric Type`: uses a deterministic validator or frozen metric when available; requires a
     human-review packet when no quantified or calibrated standard exists.
   - `Evaluator Contract`: says what evidence can prove completion and where it must be surfaced.
   - `Acceptance Criteria`: map to deliverables and to evidence available from checked sources.
   - `Verification Plan`: can prove the acceptance criteria and does not omit required source or repo
     state checks.
   - `Done Condition`: exactly matches acceptance completion, verification evidence, and stated
     non-goals. When `Metric Type` is `human_review_required`, done means producing the review packet
     and stopping for human judgment, not self-certifying the final quality claim.
   - `Runtime Completion Check`: names a command, file inspection, manual check, or evidence
     condition that can falsify completion, requires surfaced evidence in the conversation, and
     includes a retry or turn ceiling when repeated attempts are possible.
   - `Runtime Gate Evidence`: when validation lineage includes runtime gates, requires durable
     evidence: exact commands, actual and expected exit code semantics, reached/skipped attempt
     ledger, and log paths or embedded output excerpts. Transcript-only claims are insufficient
     unless a result document embeds the necessary excerpts.
   - `Open Questions`: contains unresolved facts instead of silently assuming them.
   - `Launch Context Placement`: launch text, current session status, and `/goal` startup commands
     are not embedded in the brief unless the user explicitly asks for that file shape.
   - If validation lineage evidence exists, a launch context that only delegates completion to the
     brief without naming a runtime check or bounded stop is `revised`, or `blocked` when sources do
     not support a checkable condition.
   - If runtime gates are documented but durable runtime evidence is missing or reduced to "it
     passed" transcript claims, use `brief_audit: revised`; use `blocked` or `open_questions` when
     checked sources cannot establish durable evidence requirements.
   - If sources were not inspected enough to verify a major claim, do not mark the audit as pass;
     report `brief_audit: blocked` or record the missing source under `open_questions`.
7. Rewrite the Goal Brief using `references/GOAL_CONTEXT_TEMPLATE.md` when the selected behavior
   includes edits and audit findings require changes. Keep background tight and move future behavior
   out of `Context`.
8. If the audit passes without required edits, set `brief_changes` to `none_needed` only when
   `source_checks` names the sources checked and the specific major claims they support. Do not use
   `format: pass` as a substitute for source-grounded content audit.
9. Write acceptance criteria that validate the real required outcome, not the implementation steps
   or just the document format:
   - Use `Given / When / Then / Evidence` for behavior-oriented criteria.
   - Use checklist criteria only for pure documentation work, and still include `Evidence`.
   - Include at least one primary outcome criterion.
   - Add a boundary, failure, permission, empty-state, duplicate-state, external-service, or
     non-goal criterion when that risk is relevant.
   - When converting a prior plan into a Goal Brief, preserve every acceptance-relevant evidence
     requirement from the plan. Do not drop evidence requirements merely to avoid step-by-step
     implementation detail.
10. Apply the template's current-stage-first semantics when validation lineage evidence exists. Do not
   invent lineage, convert unrelated historical labels into required reruns, or require a full-chain
   rerun unless the brief explicitly requires it or current evidence questions dependency validity.
11. If required information is missing, record it under `Open Questions` instead of inventing facts,
   behavior, evidence, or success conditions.
12. Generate `goal_launch_context` when the selected behavior includes launch output:
    - Return one copy-ready fenced `md` block.
    - The first line inside the fenced block must be a compact `/goal <objective>` command so it can
      be pasted directly into goal mode.
    - The `/goal` command must use the Goal Brief's `Objective` as the primary objective summary.
      Copy the objective exactly when it is short enough; otherwise tightly summarize it without
      changing its meaning.
    - The same `/goal` line must point to the Goal Brief as the acceptance brief. Prefer
      `Use <GOAL.md> as the acceptance brief`.
    - The same `/goal` line must name a runtime check, command, inspection, or evidence condition
      that can falsify completion; require that result evidence to be reported in the conversation;
      and include a brief-defined or default bounded stop condition.
    - Do not use `/goal Read <GOAL.md>` as the default wording. Reading the brief is an implied
      setup action; the `/goal` line should describe the objective and completion contract.
    - Do not use `Complete the goal defined in <GOAL.md>` as the default wording. It is too abstract
      and makes the brief itself sound like the goal.
    - Do not replace the brief's objective with a narrower meta-task such as "audit", "review", or
      "documentation verification" unless that is the exact user-authored goal in the brief.
    - Include the exact Goal Brief path or output location.
    - Include a compact `/goal` command that references the Goal Brief instead of copying it.
    - Include only transient current-state facts and acceptance gaps needed to start without prior
      chat history.
    - When validation lineage exists, include only the current stage to try first, at most one
      documented command needed to regenerate the current stage's missing input, the direct
      predecessor fallback, and any brief-defined retry or turn ceiling.
    - When validation lineage includes runtime gates, require completion evidence to report exact
      commands, exit codes, reached/skipped attempt ledger, and log paths or embedded output
      excerpts. Do not ask for full logs when a focused excerpt proves the gate result.
    - Point to the brief for objective, deliverables, acceptance criteria, verification plan,
      durable validation lineage, protocol commands, acceptance rules, fallback rules, constraints,
      non-goals, open questions, and done condition instead of restating them.
    - Include the expected completion evidence in the response.
    - Require the next agent to surface verification evidence in the conversation; do not assume the
      runtime can inspect files, diffs, or command output unless those results are reported.
    - Require final completion audit against authoritative current state. Do not let progress ledgers,
      active state notes, or prior transcript claims prove completion by themselves.
    - When `Metric Type` is `human_review_required`, require a human-review packet and stop for
      human judgment instead of calling the goal complete.
    - Include a clear warning not to paste launch context into the Goal Brief, exactly once.
    - Do not include step-by-step implementation instructions unless a specific procedure is itself
      part of the acceptance criteria.
13. Format the final response for a user, not as an internal field dump:
    - Use short Markdown sections, usually `Summary`, `Brief Audit`, `Goal Launch Context`, and
      `Verification`.
    - Show `source_checks` as compact bullets with source, checked claim, and result. Avoid wide
      tables because long file names wrap poorly in CLI output.
    - Keep `brief_audit` to one outcome plus only material findings.
    - Keep `brief_changes` to changed sections or `none_needed` with one short rationale.
    - Do not expose `runtime_goal_command` and `next_session_context_prompt` as separate top-level
      fields unless the user explicitly asks for them.
    - Avoid repeating the full acceptance criteria, source list, or verification plan outside the
      Goal Brief.
14. Self-review before returning:
    - `Context` contains no future deliverables or acceptance criteria.
    - `Claim Boundary` has `Accepted Scope`, `Non-Claims`, and `Rejected Evidence`.
    - The accepted scope does not claim broad generalization when evidence only supports scaffolded
      synthetic compositional recombination or another narrower condition.
    - `Goal` contains no long background narrative.
    - `Metric Type` is present and matches available evidence. Missing objective standards route to
      human review instead of agent self-certification.
    - `Evaluator Contract` and `Stop Or Escalate` are present or explicitly blocked as missing.
    - `brief_audit` is present when the selected behavior includes a brief.
    - `source_checks` is present when the selected behavior includes a brief and `brief_audit` is
      `pass` or `revised`.
    - `format: pass` only means section separation is correct; it does not satisfy `brief_audit`.
    - `brief_audit: pass` is invalid unless `source_checks` names the evidence used.
    - `brief_changes: none_needed` is invalid unless major brief claims have source-backed rationale.
    - Launch output was generated only after required brief audit and brief edits.
    - `brief_changes` is not empty when audit findings require edits.
    - `goal_launch_context` references the Goal Brief and does not duplicate long context,
      acceptance criteria, source lists, or verification plans.
    - `goal_launch_context` tells the fresh session which `/goal` command to activate or use.
    - The first line inside the `goal_launch_context` fenced block starts with `/goal`.
    - The first line uses the Goal Brief's `Objective` as the objective summary, points to the brief
      as the acceptance brief, names a runtime completion check, requires surfaced evidence, and does
      not use `Read` as the default verb.
    - The first line does not use `Complete the goal defined in <GOAL.md>` as the default wording.
    - The first line does not narrow the goal to an audit, review, or verification task unless that
      narrower task is the actual brief objective.
    - `goal_launch_context` is a launcher and does not restate the full Goal Brief.
    - `goal_launch_context` is not embedded in the Goal Brief unless explicitly requested.
    - `goal_launch_context` does not prescribe how to implement the solution unless that procedure
      is part of the requirement.
    - `goal_launch_context` does not rely only on "use the brief" for completion; it includes a
      checkable condition and bounded stop.
    - Runtime-gated lineage includes durable runtime evidence requirements; transcript-only evidence
      is not accepted unless necessary output excerpts are embedded in a result document.
    - Prior-plan acceptance evidence is fully mapped into the Goal Brief's acceptance criteria,
      verification plan, done condition, or explicit acceptance gaps.
    - Validation lineage is present only when checked sources show it exists, and the brief captures
      the template fields as far as sources support them.
    - Missing generated or transient artifacts use documented regeneration before terminal blockers,
      and unrelated rejected branches are not treated as required reruns.
    - Acceptance criteria validate the real goal outcome and include observable `Then` and concrete
      `Evidence`.
    - Vague quality words such as "better", "complete", "normal", "reasonable", and "proper" are
      replaced with concrete outcomes or moved to `Open Questions`.
15. Stop after the goal document or review and launch context are complete. Do not proceed into
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

Return a concise Markdown response. Prefer this shape:

- `Summary`: whether the Goal Brief was drafted, revised, reviewed, unchanged, or blocked, with exact
  Markdown paths changed if any.
- `Brief Audit`: `format`, compact bullet `source_checks`, `brief_audit`, `brief_changes`, and
  `claim_boundary_review`, `metric_type_review`, `acceptance_criteria_review`, and
  `completion_contract_review`.
- `Goal Launch Context`: one fenced `md` block named `goal_launch_context`, unless the selected
  behavior excludes launch output.
- `Verification`: checks run, assumptions, open questions, and manual verification checklist.

Do not return a long YAML-like list of every internal field. Do not expose `runtime_goal_command` and
`next_session_context_prompt` as separate top-level outputs unless the user explicitly requests that
lower-level split.

`goal_launch_context` should use this compact form:

```md
/goal <objective copied or tightly summarized from the Goal Brief>. Use `<goal brief path>` as the acceptance brief. Runtime check: <named command, inspection, or evidence condition>. Report that evidence in conversation. Stop and report blocked after <brief-defined ceiling or default bounded attempts>. Stay within the brief's constraints and non-goals.

Current context:
- <transient current-state facts only>

Current-stage first attempt:
- <current stage to rerun/complete first; at most one documented command to regenerate its missing input if needed; direct predecessor fallback if blocked; retry/turn ceiling if specified>

Acceptance gaps:
- <missing evidence or "Use the brief to verify all acceptance criteria.">

Completion evidence to report:
- Files touched
- Acceptance status
- Verification results
- Metric type and evaluator contract result
- Completion audit against authoritative current state
- Runtime gate evidence when applicable: exact command, exit code, reached/skipped attempt ledger, and log path or embedded output excerpt
- Blockers or follow-up

Do not paste this launch context into the Goal Brief. Use the brief for stable context, durable
validation lineage, objective, deliverables, acceptance criteria, verification plan, constraints,
non-goals, open questions, and done condition.
```

## Version History

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

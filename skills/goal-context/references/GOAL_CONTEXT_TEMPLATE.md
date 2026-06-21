# Goal Brief Template

Use this template when drafting, revising, or reviewing goal documents for agent handoff. The
template intentionally separates stable `Context` from future `Goal` so background facts do not get
mixed with outcomes, acceptance criteria, or launch prompts.

```md
# Goal Brief: <short title>

## Context

### Background

<Why this goal exists. Past decisions, conversation facts, or project state only.>

### Current State

<What is true now. Include stable status, blockers, and known constraints.>

### Claim Boundary

- Accepted Scope: <what this stage is allowed to claim>
- Non-Claims: <nearby stronger claims that are not established>
- Rejected Evidence: <prior, tempting, stale, or too-broad evidence that must not count, with why>

### Validation Lineage

<Optional. Include only when checked sources show validation lineage evidence. Otherwise use `None
found in checked sources`.>

- Current stage: <stage to complete or rerun first>
- Direct dependency predecessor: <immediate predecessor stage or `None documented`>
- Relevant prior final labels: <accepted/rejected/inconclusive/blocked labels that affect context>
- Required input artifacts: <project-defined transient artifact path, generated checkpoint/log
  location, or durable source artifact, using generic descriptions when paths are unstable>
- Regeneration/reproduction path: <documented regeneration command, reproduction command, or
  `None documented`>
- Retry/turn ceiling: <brief-defined limit, or default bounded attempt policy>
- Runtime gate evidence: <when runtime gates exist, exact command(s), exit code(s),
  reached/skipped attempt ledger, and log paths or embedded output excerpts>
- Result evidence sources: <result docs, logs, acceptance records, or reviewed evidence>

### Relevant Sources

- <Exact file path, issue, doc, source, or conversation fact.>

### Constraints

- <Policy, technical, scope, time, tooling, or environment constraints.>

### Non-Goals

- <Nearby outcomes that must not be pulled into this goal.>

### Open Questions

- <Questions that affect correctness, or `None`.>

## Goal

### Objective

<One concrete outcome sentence.>

### Deliverables

- <Artifact, behavior, doc, code change, or decision expected from the work.>

### Metric Type

<Use one: deterministic_validator | frozen_metric | evidence_review | subjective_rubric |
human_review_required. Explain the selected standard and why broader or weaker evidence is
insufficient.>

### Evaluator Contract

<What evidence proves completion, where it must be surfaced, and whether a human reviewer is
required. Do not assume provider memory, an independent evaluator, context isolation, or hidden tool
access.>

### Acceptance Criteria

- [ ] AC-1: <happy path / primary outcome>
  Given <starting context or precondition>
  When <one action, event, or completed change occurs>
  Then <observable result or artifact exists>
  Evidence: <test, review check, manual check, file inspection, or demo>

- [ ] AC-2: <failure, boundary, permission, or non-goal guard>
  Given <edge condition or excluded condition>
  When <one action, event, or attempted change occurs>
  Then <observable safe result>
  Evidence: <test, review check, manual check, file inspection, or demo>

### Verification Plan

- <Smallest useful command, review, manual checklist, or evidence path.>
- Runtime completion check: <single command, file inspection, manual check, or named evidence
  condition that can falsify completion.>
- Durable runtime evidence required: <for runtime gates, exact command(s), exit code(s),
  reached/skipped attempt ledger, and log paths or embedded output excerpts.>
- Completion audit source: <authoritative current-state source to inspect before claiming done.>

### Stop Or Escalate

<Retry, turn, budget, blocker, tool-failure, and human-review stop rules. If Metric Type is
`human_review_required`, stop after producing the review packet and wait for human judgment.>

### Done Condition

<The exact condition where the agent should stop and report complete, including surfaced runtime
completion evidence, durable runtime evidence for runtime gates, authoritative completion audit,
human-review packet when required, and any retry/turn ceiling.>
```

## Section Rules

- `Context` contains only background, current state, constraints, sources, non-goals, and unknowns.
  Do not put future deliverables or acceptance criteria here.
- `Claim Boundary` is required for long-running goal use. It must include `Accepted Scope`,
  `Non-Claims`, and `Rejected Evidence`.
- `Goal` contains only objective, deliverables, metric type, evaluator contract, acceptance
  criteria, verification, stop or escalate rules, and done condition. Do not put long background
  narrative here.
- `Metric Type`, `Evaluator Contract`, and `Stop Or Escalate` are part of the completion contract.
  They are not provider runtime features.
- Replace placeholders instead of leaving angle-bracket text in the final file.
- Use absolute dates for status and decisions.
- Cite exact local paths when the context depends on repo files.
- Do not include secrets, credentials, private tokens, or unverifiable claims.
- Do not include session launch context, current worktree chatter, `/goal` startup text, or
  copy-ready next-session prompts in the Goal Brief. Return those in `goal_launch_context` instead.
- `Validation lineage evidence` means checked sources document a replayable gate chain, acceptance
  protocol, result document, checkpoint or artifact dependency, or prior final label. Do not invent
  lineage for projects that have only ordinary verification steps.
- Keep durable lineage, protocol commands, acceptance rules, fallback rules, and result evidence
  sources in the Goal Brief. Keep the launch context limited to the transient current-stage rerun,
  immediate acceptance gap, necessary predecessor fallback, and completion evidence.

## Goal Brief Audit Rules

When the selected behavior includes drafting, revising, or reviewing a Goal Brief, audit the brief
before generating launch output. Treat this as a source-grounded content gate, not a formatting
check.

Audit these sections:

- `Context`: factual, stable, sourced, and free of future deliverables.
- `Claim Boundary`: includes accepted scope, non-claims, and rejected evidence. Accepted scope must
  not be broader than the evidence supports.
- `Validation Lineage`: required only when validation lineage evidence exists. It must identify the
  current stage, direct dependency predecessor, relevant prior final labels, required input artifacts,
  documented regeneration or reproduction command, retry/turn ceiling, and result evidence sources
  as far as sources support them.
- `Goal`: future outcomes only, not background narrative.
- `Deliverables`: actual project artifacts, behavior, or decisions, not handoff metadata.
- `Metric Type`: uses deterministic validators or frozen metrics when available. If no quantified,
  replayable, evidence-review, or calibrated rubric standard exists, it must be
  `human_review_required`.
- `Evaluator Contract`: states what evidence can prove completion and where the evidence must be
  surfaced. It must not assume provider memory, hidden evaluator state, or unsurfaced tool output.
- `Acceptance Criteria`: map to deliverables, validate the real goal, and include observable
  evidence available from checked sources.
- `Verification Plan`: can prove the acceptance criteria and does not omit required source or repo
  state checks.
- `Done Condition`: exactly matches acceptance completion, verification evidence, and stated
  non-goals, and includes a bounded stop condition when repeated attempts are possible.
- `Stop Or Escalate`: covers retry, turn, budget, blocker, tool-failure, and human-review stops.
- `Runtime Completion Check`: names a command, file inspection, manual check, or evidence condition
  that can falsify completion and can be reported in the conversation.
- `Runtime Gate Evidence`: when validation lineage includes runtime gates, requires durable
  evidence: exact commands, actual and expected exit code semantics, reached/skipped attempt ledger,
  and log paths or embedded output excerpts.
- `Open Questions`: contains unresolved facts instead of silently assuming them.
- `Launch Context Placement`: launch prompts and current session status are not embedded in the
  brief unless the user explicitly asks for that file shape.

For every major claim in the brief, verify it against the provided sources or repo state:

- Current State claims must match referenced files, git state, or conversation facts.
- Validation Lineage claims must match documented gates, protocols, result docs, artifact
  dependencies, reproduction commands, and final labels.
- Deliverables must correspond to actual expected files, decisions, or artifacts.
- Acceptance Criteria must map to deliverables and to evidence from the referenced sources.
- Verification Plan must be sufficient to prove the ACs.
- Done Condition must not add requirements outside AC completion and stated non-goals.
- Metric Type and Evaluator Contract must not let the agent self-certify subjective or unquantified
  claims. Human review is required unless an objective validator, frozen metric, evidence-review
  protocol, or calibrated rubric exists.
- Runtime Completion Check must be specific enough for a fresh session to know what evidence to
  report, not only a pointer to the brief.
- Runtime Gate Evidence must be durable enough for review after the session. Transcript-only claims
  are insufficient unless the result document embeds the necessary output excerpts.

If sources were not inspected enough to verify a major claim, do not mark the audit as pass. Use
`brief_audit: blocked` or record the missing source under `open_questions`.

When lineage evidence exists, treat these as `brief_audit: revised` or `brief_audit: blocked` if
they affect runtime behavior:

- The current stage is unclear.
- The direct dependency predecessor is unclear.
- An artifact dependency lacks a documented regeneration or reproduction path.
- A missing generated or transient artifact is treated as terminal `blocked` before regeneration is
  attempted.
- A rejected prior branch is incorrectly added to the required rerun path when the current stage does
  not directly depend on it.
- Launch context copies long-term lineage, full protocol commands, or full acceptance criteria
  instead of pointing back to the Goal Brief.
- Launch context only delegates completion to the Goal Brief without a named runtime check,
  conversation-surfaced evidence requirement, or bounded stop.
- Runtime-gated validation is marked complete using only transcript claims such as "passed" or
  "looked good" without exact command, exit code, reached/skipped attempt ledger, and log path or
  embedded output excerpt.
- Broad claims are not accepted from narrow evidence. For example, `scaffolded synthetic
  compositional recombination` evidence must not be rewritten as broad generalization.

Do not use `format: pass` as a substitute for this audit. `format: pass` only means the `Context` /
`Goal` section split is structurally correct.

Launch context generation is not a substitute for repairing the stable brief. If audit findings
require edits, revise the brief or report the missing input that blocks the revision before
producing launch output.

## Acceptance Criteria Rules

- Each acceptance criterion must be observable and verifiable.
- Avoid vague terms such as "better", "complete", "normal", "reasonable", "proper", or "works"
  unless the criterion defines how to observe that outcome.
- Behavior-oriented criteria should use `Given / When / Then / Evidence`.
- Pure documentation criteria may use checklist wording, but must still include `Evidence`.
- Include at least one primary outcome criterion.
- Add a boundary, failure, permission, empty-state, duplicate-state, external-service, or non-goal
  criterion whenever that risk affects correctness.
- When converting from a prior plan, preserve every acceptance-relevant evidence requirement from
  that plan. Do not remove required evidence merely to avoid step-by-step implementation detail.
- If the available information is not enough to write an acceptance criterion, put the missing input
  under `Open Questions` instead of inventing a success condition.

## Metric And Human Review Rules

- Prefer deterministic validators and frozen metrics when they exist.
- A frozen metric must be outside the agent's edit boundary.
- `evidence_review` is acceptable only when the criteria and required artifacts are named.
- `subjective_rubric` must include criteria, negative examples or rejected patterns, and calibration
  status.
- If the goal depends on taste, broad quality, or direction and no calibrated rubric exists, use
  `human_review_required`.
- For `human_review_required`, the done condition is a review packet with accepted scope,
  non-claims, rejected evidence, artifacts changed, verification performed, open questions, and the
  specific human decision needed.
- LLM-as-judge is supporting evidence only unless the user explicitly accepts it as the evaluator.

## Validation Lineage Rules

- Lineage is optional. Require it only when validation lineage evidence exists.
- The current stage is the first rerun or completion target.
- If the current stage cannot proceed because a generated or transient input artifact is missing,
  use the documented regeneration or reproduction command before declaring the goal blocked.
- If regeneration fails, the source/runtime is unavailable, no documented regeneration path exists,
  or predecessor validity cannot be established, fall back to validating the direct dependency
  predecessor.
- Do not require a full-chain rerun unless the brief explicitly requires it or current evidence
  questions dependency validity.
- Include at most the single command needed to regenerate the current stage's missing input in launch
  context. Do not include the full protocol, downstream stage commands, or historical stage rerun
  lists there.
- Prior final labels such as accepted, rejected, inconclusive, and blocked are lineage context. A
  rejected branch is not a required rerun unless the current stage directly depends on that branch.
- If an accepted predecessor supplies current input artifacts, record its reproduction path or
  predecessor verification path when the sources document one.
- If the brief does not define a retry or turn ceiling, default to one current-stage attempt, one
  documented regeneration attempt when needed, and one direct predecessor fallback attempt before
  reporting blocked with evidence.
- The launcher's `/goal` line should include a runtime check that can be reported in the
  conversation. Do not rely only on "use the brief" for completion.
- When validation lineage includes runtime gates, durable runtime evidence is required. Record exact
  commands, actual and expected exit code semantics, reached/skipped attempt ledger, and log paths or
  embedded output excerpts.
- If a log path is unavailable or points to a transient location, embed enough output excerpt in the
  result document or conversation summary to prove the gate result. Do not copy full logs when a
  focused excerpt is sufficient.

Example:

```md
### Validation Lineage

- Current stage: Stage C validation, attempted first.
- Direct dependency predecessor: Stage B output validation, previously accepted.
- Relevant prior final labels: Stage B accepted; one alternate Stage C branch rejected and kept as
  context only because the current stage does not depend on it.
- Required input artifacts: generated checkpoint/log location for Stage C input.
- Regeneration/reproduction path: run the documented regeneration command for the current stage's
  missing input.
- Retry/turn ceiling: one current-stage attempt, one documented regeneration attempt, and one direct
  predecessor fallback attempt unless the brief states a stricter limit.
- Runtime gate evidence: exact Stage C validation command; exit code `0` means accepted and nonzero
  means rejected or blocked by the documented gate; attempt ledger records current-stage attempt
  reached, regeneration attempt skipped or reached, and predecessor fallback skipped or reached;
  generated log location or embedded output excerpt records the decisive result lines.
- Result evidence sources: accepted predecessor result doc, current-stage failure note, and
  documented regeneration instructions.

Example launch line:

`/goal Validate Stage C output. Use <goal brief path> as the acceptance brief. Runtime check: run the documented Stage C validation check and report durable runtime evidence in conversation: exact command, exit code, reached/skipped attempt ledger, and generated log location or embedded output excerpt. Stop and report blocked after one current-stage attempt, one documented regeneration attempt if needed, and one predecessor fallback.`
```

## Goal Launch Context Rules

Return a copy-ready `goal_launch_context` in the assistant response after drafting, revising, or
reviewing the Goal Brief, unless the user explicitly asks for brief-only output. It should launch the
next session, not restate the full brief. Do not embed this prompt in the Goal Brief Markdown unless
the user explicitly asks for that file shape.

The launch context should include:

- Goal Brief path or location.
- A compact `/goal <objective>` command or equivalent runtime completion contract.
- The `/goal` line should use the Goal Brief's `Objective` as the primary objective summary. Copy
  the objective exactly when it is short enough; otherwise tightly summarize it without changing its
  meaning.
- The same `/goal` line should point to the Goal Brief as the acceptance brief, preferably with
  `Use <GOAL.md> as the acceptance brief`.
- The same `/goal` line should name a runtime completion check, require that result evidence to be
  reported in the conversation, and include a brief-defined or default bounded stop condition.
- Do not use `/goal Read <GOAL.md>` as the default wording. Reading the brief is setup; the goal
  line should state the completion objective.
- Do not use `Complete the goal defined in <GOAL.md>` as the default wording. It is too abstract and
  makes the brief itself sound like the goal.
- Do not rename the brief objective into a narrower meta-task such as audit, review, or verification
  unless that narrower task is the exact user-authored objective.
- Only transient current-state context and partial acceptance evidence the next agent needs before
  acting.
- When validation lineage exists, include only the current stage to try first, at most one documented
  command needed to regenerate the current stage's missing input, the direct predecessor fallback if
  the current stage cannot proceed, and any brief-defined retry or turn ceiling.
- Acceptance gaps, when the goal is already partially done.
- Completion evidence requirements.
- For runtime gates, completion evidence must include exact command, exit code, reached/skipped
  attempt ledger, and log path or embedded output excerpt. Transcript-only evidence is insufficient
  unless the result document embeds the necessary excerpts.
- A clear instruction that verification evidence must be surfaced in the conversation. Do not assume
  the runtime can inspect files, diffs, or command output unless those results are reported.
- A clear instruction to use the brief for objective, deliverables, acceptance criteria,
  verification plan, durable validation lineage, protocol commands, acceptance rules, fallback rules,
  constraints, non-goals, open questions, and done condition.
- A clear instruction not to paste launch context into the Goal Brief.
- A clear instruction not to include step-by-step implementation directions unless a specific
  procedure is part of the acceptance criteria.

Use this structure by default:

```md
/goal <objective copied or tightly summarized from the Goal Brief>. Use `<goal brief path>` as the acceptance brief. Runtime check: <named command, inspection, or evidence condition>. Report that evidence in conversation. Stop and report blocked after <brief-defined ceiling or default bounded attempts>. Stay within the brief's constraints and non-goals.

Current context:
- <transient status or partial acceptance evidence not worth adding to the brief>

Current-stage first attempt:
- <current stage to rerun/complete first; at most one documented command to regenerate its missing input if needed; direct predecessor fallback if the current stage cannot proceed; retry/turn ceiling if specified>

Acceptance gaps:
- <missing evidence or unmet acceptance criterion, or "Use the brief to verify all acceptance criteria.">

Completion evidence to report:
- Files touched
- Acceptance status
- Verification results
- Metric type and evaluator contract result
- Completion audit against authoritative current state
- Runtime gate evidence when applicable: exact command, exit code, reached/skipped attempt ledger, and log path or embedded output excerpt
- Human-review packet when the brief requires human judgment
- Blockers or follow-up

Use the brief for objective, deliverables, acceptance criteria, verification plan, durable validation
lineage, protocol commands, acceptance rules, fallback rules, constraints, non-goals, open questions,
and done condition. Do not duplicate those sections in this launcher.

Surface verification evidence in the conversation. Do not assume the runtime can read files, inspect
diffs, or see command output unless those results are reported.

Do not paste this launch context into the Goal Brief.

Do not turn this launch context into a step-by-step implementation plan. It should state what must be
true, what evidence is missing, how to verify completion, and when to stop as blocked.
```

## Response Shape Rules

The assistant response should be concise and user-facing. Do not return a long YAML-like dump of
internal fields.

Use this default response shape:

````md
## Summary

`GOAL.md` revised / reviewed / unchanged / blocked.
Files touched: <paths or none>

## Brief Audit

format: pass / revised / blocked

source_checks:
- `<path or command>`: <major claim checked> -> pass / revised / blocked
- `<path or command>`: <major claim checked> -> pass / revised / blocked

brief_audit: pass / revised / blocked
brief_changes: <changed sections, none_needed with one short rationale, or blocker>
claim_boundary_review: pass / revised / blocked
metric_type_review: pass / revised / blocked
acceptance_criteria_review: pass / revised / blocked
completion_contract_review: pass / revised / blocked

## Goal Launch Context

```md
<copy-ready goal_launch_context>
```

## Verification

| Check | Result |
| --- | --- |
| `<command or manual check>` | <result> |

Assumptions: <correctness-relevant assumptions only>
Open questions: <None or blockers>
Manual verification:
- <focused checklist item>
````

Formatting requirements:

- Keep `source_checks` to 3-7 high-signal bullets. Do not list every searched symbol or every line
  inspected.
- Keep `brief_audit` to one outcome plus material findings. Do not restate every brief section when
  the audit passes.
- Keep `brief_changes: none_needed` to one sentence when no edit was required.
- Return exactly one fenced `goal_launch_context` block by default.
- The first line inside the `goal_launch_context` block must start with `/goal`.
- The first line should use the Goal Brief's `Objective` as the objective summary and point to the
  brief as the acceptance brief.
- The first line should name a runtime check, require evidence to be reported in conversation, and
  include a bounded stop.
- The first line should avoid both `/goal Read <GOAL.md>` setup wording and `Complete the goal
  defined in <GOAL.md>` meta wording.
- The first line must not narrow the goal into an audit, review, or verification task unless that is
  the actual brief objective.
- Do not separately expose `runtime_goal_command` and `next_session_context_prompt` unless the user
  explicitly asks for that lower-level split.
- Include the warning not to paste launch context into the Goal Brief exactly once.
- Do not repeat the full acceptance criteria, full source list, or full verification plan outside
  the Goal Brief.

## Skill Eval Checklist

Use a small manual prompt set when changing this skill:

- Should trigger: user explicitly asks `$goal-context` to repair a goal brief.
- Should trigger: user explicitly asks `$goal-context` for launch context after partial work.
- Should trigger: user explicitly asks `$goal-context` to create goal-mode startup text from a Goal
  Brief.
- Should trigger: user explicitly asks `$goal-context` to audit whether context and goal are mixed.
- Should not trigger: user asks to implement the goal now.
- Should not trigger: user asks for a commit, code review, diagnosis, or ordinary roadmap plan.

Expected outputs for trigger cases:

- A Goal Brief is drafted, revised, reviewed, or explicitly skipped by selected behavior.
- Audit-only requests do not edit files unless the user also asks to fix or rewrite them.
- Rewrite requests update the target Goal Brief when audit findings require changes.
- `source_checks` is present for brief modes and names files, repo state, or conversation facts used
  by the audit.
- `brief_audit` is present for brief modes and reports pass, revised findings, or a concrete
  blocker. Pass is invalid without `source_checks`.
- `brief_changes` names changed sections, or says `none_needed` with source-backed rationale.
- `goal_launch_context` is present unless the user explicitly selected brief-only mode.
- `goal_launch_context` is a fenced `md` block.
- The first line inside `goal_launch_context` starts with `/goal`.
- The first line uses the Goal Brief's `Objective` as the objective summary and points to the brief
  as the acceptance brief.
- The first line does not default to `/goal Read <GOAL.md>` or `Complete the goal defined in
  <GOAL.md>`.
- The first line does not rename the goal into a narrower audit, review, or verification task unless
  that is the actual brief objective.
- `goal_launch_context` contains the compact `/goal` command, transient current context, acceptance
  gaps, current-stage-first attempt when relevant, completion evidence, and one warning not to paste
  launch context into the Goal Brief.
- `goal_launch_context` does not rely only on `Use <GOAL.md> as the acceptance brief`; it includes a
  runtime completion check and bounded stop.
- `goal_launch_context` requires authoritative current-state completion audit, not progress-ledger
  or transcript-only completion.
- If no objective or calibrated standard exists, the launch context stops at a human-review packet
  instead of telling the agent to self-certify completion.
- Runtime-gated launch contexts require durable runtime evidence. Transcript-only evidence fails
  unless a result document embeds the needed output excerpts.
- The `/goal` command references the Goal Brief and does not copy long context, full acceptance
  criteria, source lists, or detailed verification plans.
- Runtime context is not embedded into the Goal Brief.
- Acceptance criteria validate the actual target outcome, not just the brief format.
- Prior-plan acceptance evidence requirements are preserved in the Goal Brief's acceptance criteria,
  verification plan, done condition, or explicit acceptance gaps.
- Replayable validation lineage is included only when source evidence shows it exists.
- Current-stage-first rerun, documented regeneration or reproduction path, and direct predecessor
  fallback are represented without making all historical stages required reruns.
- Missing generated or transient artifacts are not treated as terminal blockers when a documented
  regeneration or reproduction path exists.
- Prior rejected branches are not required reruns unless the current stage directly depends on them.
- The final response is concise Markdown, not a long field-by-field dump.
- Source checks are readable bullets, not a wide table that wraps long file names.

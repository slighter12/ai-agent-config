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

### Done Condition

<The exact condition where the agent should stop and report complete.>
```

## Section Rules

- `Context` contains only background, current state, constraints, sources, non-goals, and unknowns.
  Do not put future deliverables or acceptance criteria here.
- `Goal` contains only objective, deliverables, acceptance criteria, verification, and done
  condition. Do not put long background narrative here.
- Replace placeholders instead of leaving angle-bracket text in the final file.
- Use absolute dates for status and decisions.
- Cite exact local paths when the context depends on repo files.
- Do not include secrets, credentials, private tokens, or unverifiable claims.
- Do not include session launch context, current worktree chatter, `/goal` startup text, or
  copy-ready next-session prompts in the Goal Brief. Return those in `goal_launch_context` instead.

## Goal Brief Audit Rules

When the selected behavior includes drafting, revising, or reviewing a Goal Brief, audit the brief
before generating launch output. Treat this as a source-grounded content gate, not a formatting
check.

Audit these sections:

- `Context`: factual, stable, sourced, and free of future deliverables.
- `Goal`: future outcomes only, not background narrative.
- `Deliverables`: actual project artifacts, behavior, or decisions, not handoff metadata.
- `Acceptance Criteria`: map to deliverables, validate the real goal, and include observable
  evidence available from checked sources.
- `Verification Plan`: can prove the acceptance criteria and does not omit required source or repo
  state checks.
- `Done Condition`: exactly matches acceptance completion, verification evidence, and stated
  non-goals.
- `Open Questions`: contains unresolved facts instead of silently assuming them.
- `Launch Context Placement`: launch prompts and current session status are not embedded in the
  brief unless the user explicitly asks for that file shape.

For every major claim in the brief, verify it against the provided sources or repo state:

- Current State claims must match referenced files, git state, or conversation facts.
- Deliverables must correspond to actual expected files, decisions, or artifacts.
- Acceptance Criteria must map to deliverables and to evidence from the referenced sources.
- Verification Plan must be sufficient to prove the ACs.
- Done Condition must not add requirements outside AC completion and stated non-goals.

If sources were not inspected enough to verify a major claim, do not mark the audit as pass. Use
`brief_audit: blocked` or record the missing source under `open_questions`.

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
- Do not use `/goal Read <GOAL.md>` as the default wording. Reading the brief is setup; the goal
  line should state the completion objective.
- Do not use `Complete the goal defined in <GOAL.md>` as the default wording. It is too abstract and
  makes the brief itself sound like the goal.
- Do not rename the brief objective into a narrower meta-task such as audit, review, or verification
  unless that narrower task is the exact user-authored objective.
- Only transient current-state context and partial acceptance evidence the next agent needs before
  acting.
- Acceptance gaps, when the goal is already partially done.
- Completion evidence requirements.
- A clear instruction that verification evidence must be surfaced in the conversation because the
  runtime goal evaluator may not inspect files or command output directly.
- A clear instruction to use the brief for objective, deliverables, acceptance criteria,
  verification plan, constraints, non-goals, open questions, and done condition.
- A clear instruction not to paste launch context into the Goal Brief.
- A clear instruction not to include step-by-step implementation directions unless a specific
  procedure is part of the acceptance criteria.

Use this structure by default:

```md
/goal <objective copied or tightly summarized from the Goal Brief>. Use `<goal brief path>` as the acceptance brief. Report verification evidence and blockers. Stay within the brief's constraints and non-goals.

Current context:
- <transient status or partial acceptance evidence not worth adding to the brief>

Acceptance gaps:
- <missing evidence or unmet acceptance criterion, or "Use the brief to verify all acceptance criteria.">

Completion evidence to report:
- Files touched
- Acceptance status
- Verification results
- Blockers or follow-up

Use the brief for objective, deliverables, acceptance criteria, verification plan, constraints,
non-goals, open questions, and done condition. Do not duplicate those sections in this launcher.

Surface verification evidence in the conversation. Do not assume the runtime goal evaluator can read
files, inspect diffs, or see command output unless those results are reported.

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
acceptance_criteria_review: pass / revised / blocked

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
  gaps, completion evidence, and one warning not to paste launch context into the Goal Brief.
- The `/goal` command references the Goal Brief and does not copy long context, full acceptance
  criteria, source lists, or detailed verification plans.
- Runtime context is not embedded into the Goal Brief.
- Acceptance criteria validate the actual target outcome, not just the brief format.
- Prior-plan acceptance evidence requirements are preserved in the Goal Brief's acceptance criteria,
  verification plan, done condition, or explicit acceptance gaps.
- The final response is concise Markdown, not a long field-by-field dump.
- Source checks are readable bullets, not a wide table that wraps long file names.

# Goal Context Eval Checklist

Use this small manual prompt set only when changing `goal-context`.

## Trigger Checks

- Should trigger: user explicitly asks `$goal-context` to repair a goal brief.
- Should trigger: user explicitly asks `$goal-context` for launch context after partial work.
- Should trigger: user explicitly asks `$goal-context` to create goal-mode startup text from a Goal
  Brief.
- Should trigger: user explicitly asks `$goal-context` to audit whether context and goal are mixed.
- Regression: `$goal-context audit GOAL.md` returns `goal_launch_context` even when no file edits are
  made.
- Regression: `$goal-context repair GOAL.md` returns both brief audit/change status and
  `goal_launch_context`.
- Regression: `$goal-context repair docs/GOAL.md` uses exactly `docs/GOAL.md` as the target Goal
  Brief path.
- Regression: `$goal-context repair my goal brief` with no explicit path and no existing root
  `GOAL.md` asks for the target path or inline-output format instead of guessing.
- Regression: `$goal-context repair my goal brief` in a repo with existing root `GOAL.md` uses root
  `GOAL.md` as the target Goal Brief path.
- Regression: when objective, evaluator data, or stop rules are missing, the response includes
  `## Goal Launch Context` with `goal_launch_context: blocked` and missing facts instead of omitting
  the section.
- Should not trigger: user asks to implement the goal now.
- Should not trigger: user asks for a commit, code review, diagnosis, or ordinary roadmap plan.

## Expected Outputs

- A Goal Brief is drafted, revised, reviewed, or explicitly skipped by selected behavior.
- Audit-only requests do not edit files unless the user also asks to fix or rewrite them.
- Rewrite requests update the target Goal Brief when audit findings require changes.
- Explicit user paths are used exactly; otherwise existing root `GOAL.md` is the only default target.
- If no explicit path or existing root `GOAL.md` is available, the skill asks for a target path or
  inline-output format instead of guessing or creating root `GOAL.md`.
- `source_checks` is present for brief modes and names files, repo state, or conversation facts used
  by the audit.
- `brief_audit` is present for brief modes and reports pass, revised findings, or a concrete blocker.
  Pass is invalid without `source_checks`.
- `brief_changes` names changed sections, or says `none_needed` with source-backed rationale.
- `goal_launch_context` is present unless the user explicitly selected brief-only mode.
- Audit-only controls file edits only; it does not suppress `goal_launch_context`.
- If required facts are missing, `## Goal Launch Context` is present with
  `goal_launch_context: blocked` and the missing facts.
- When launch output is valid, `goal_launch_context` is a fenced `md` block.
- When launch output is valid, the first line inside `goal_launch_context` starts with `/goal`.
- When launch output is valid, the first line uses the Goal Brief's `Objective` as the objective
  summary and points to the brief as the acceptance brief.
- When launch output is valid, the first line does not default to `/goal Read <GOAL.md>` or
  `Complete the goal defined in <GOAL.md>`.
- When launch output is valid, the first line does not rename the goal into a narrower audit,
  review, or verification task unless that is the actual brief objective.
- When launch output is valid, `goal_launch_context` contains the compact `/goal` command, transient
  current context, acceptance gaps, current-stage-first attempt when relevant, completion evidence,
  and one warning not to paste launch context into the Goal Brief.
- When launch output is valid, `goal_launch_context` does not rely only on `Use <resolved Goal Brief
  path> as the acceptance brief`; it includes a runtime completion check and bounded stop.
- When launch output is valid, `goal_launch_context` requires authoritative current-state completion
  audit, not progress-ledger or transcript-only completion.
- When launch output is valid and no objective or calibrated standard exists, the launch context
  stops at a human-review packet instead of telling the agent to self-certify completion.
- When launch output is valid and runtime gates exist, launch contexts require durable runtime
  evidence. Transcript-only evidence fails unless a result document embeds the needed output
  excerpts.
- When launch output is valid, the `/goal` command references the Goal Brief and does not copy long
  context, full acceptance criteria, source lists, or detailed verification plans.
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

# Goal Context Eval Checklist

Use this small manual prompt set only when changing `goal-context`.

## Trigger Checks

- Should trigger: user explicitly asks `$goal-context` to repair a goal brief.
- Should trigger: user explicitly asks `$goal-context` for launch context after partial work.
- Should trigger: user explicitly asks `$goal-context` to create goal-mode startup text from a Goal
  Brief.
- Should trigger: user explicitly asks `$goal-context` to audit whether context and goal are mixed.
- Should not trigger: user asks to implement the goal now.
- Should not trigger: user asks for a commit, code review, diagnosis, or ordinary roadmap plan.

## Expected Outputs

- A Goal Brief is drafted, revised, reviewed, or explicitly skipped by selected behavior.
- Audit-only requests do not edit files unless the user also asks to fix or rewrite them.
- Rewrite requests update the target Goal Brief when audit findings require changes.
- `source_checks` is present for brief modes and names files, repo state, or conversation facts used
  by the audit.
- `brief_audit` is present for brief modes and reports pass, revised findings, or a concrete blocker.
  Pass is invalid without `source_checks`.
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
- `goal_launch_context` requires authoritative current-state completion audit, not progress-ledger or
  transcript-only completion.
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

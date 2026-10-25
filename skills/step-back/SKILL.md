---
name: step-back
description: Pause a task that keeps growing through patches, name and test the premise the patches share, and compare the current path with a fresh design before continuing. Use when the user asks to step back or reassess the direction.
disable-model-invocation: true
argument-hint: "What task should I reassess?"
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Step Back

Momentum is a poor architect. Each patch on top of existing code feels like the safe, local choice, but together they pile up into complexity no one would have designed on purpose. When fixes start stacking, stop and look at the whole picture: is the goal still the same, and if you were starting fresh with what you know now, would you build it this way?

This skill is read-only: it never edits, reverts, stashes, or commits. It may run local, read-only checks to test a premise.

## 1. Restate the Goal

Write the original goal in one sentence. If the work has drifted, write the goal it is now serving as a second sentence.

## 2. List the Pile

List what has been stacked since the goal was set: fixes, workarounds, compatibility shims, flags, retries, and new paths built beside old ones. Cite each with `file:line` or a commit.

## 3. Attack the Premise

When two or more fixes that share one premise have failed the same gate, or when workarounds keep piling up, suspect the premise, not the fixes.

**Why**: Each failed fix is evidence against the belief it assumed. A third fix on the same premise spends effort where the evidence already says the cause is not.

**Pattern**:

- **Write the premise down.** The premise is the one sentence that every failed fix assumed.
- **Test the premise before the next fix.** Find an observation that would come out differently if the premise were false, and make it. Look at how the failures are distributed, not only at how large they are. Make the check rerunnable. Run it when a local read-only check can; otherwise name the check. Report the premise as Verified, Not verified, or Inconclusive.
- **Follow the evidence.** If it contradicts the premise, the cause is elsewhere. Look for what the failures have in common that the premise did not predict. If it confirms the premise, ask the next "why". The earlier fixes may have treated a consequence instead of the cause. When the pile is fixes for one failing behavior, use the `diagnosing-bugs` skill for the evidence loop instead of repeating it here.
- **Remove the condition instead of compensating for it**, per the Laziness Protocol in `codebase-design` §3. If the evidence shows that one part always lands on the losing side, change what puts it there. A mechanism that cleans up after it leaves the cause in place and adds work on every run.

**Stop**: Do not start the next fix before the premise is written down and tested.

This step is distinct from step 4, Redesign From First Principles, which rebuilds a design around a new requirement. It questions a fact the current design assumes.

## 4. Redesign From First Principles

When integrating a change, don't bolt it onto the existing design. Redesign as if the requirement had been a foundational assumption from the start.

- Read all affected files and understand the current design
- Ask: "if we were writing this from scratch with this new requirement, what would we build?"
- Propagate the change through every reference: types, docs, examples, rationale sections
- Think about the whole redesign, then deliver it incrementally

This is the method for preserving option value when integrating changes into an existing design. Compare that design with finishing the current path: the cost to finish, the cost to switch, and what carries over. When the fresh design moves a module boundary or interface, use the `codebase-design` skill to shape it.

## 5. Recommend in Four Lines

If the answer is no, don't push through quietly. Tell the user in a few lines: where things are heading, why that's drifting, the cleaner alternative, and what you recommend.

- **Heading**: where the current path ends up.
- **Drift**: why it is drifting, with the premise and its verdict.
- **Alternative**: the cleaner path and what it reuses.
- **Recommendation**: continue, switch, or ask, and why.

Then stop and wait for the user's decision. If the answer is yes, say so plainly, so the user can carry on with confidence. The goal is to catch the wrong turn early, not to second-guess every step.

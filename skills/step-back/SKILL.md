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

Each patch on top of existing code feels like the safe, local choice, but together they pile into complexity no one would design on purpose. Stop, look at the whole, and let the user choose the direction with evidence in hand.

This skill is read-only: it never edits, reverts, stashes, or commits. It may run local, read-only checks to test a premise.

## 1. Restate the Goal

Write the original goal in one sentence. If the work has drifted, write the goal it is now serving as a second sentence.

## 2. List the Pile

List what has been stacked since the goal was set: fixes, workarounds, compatibility shims, flags, retries, and new paths built beside old ones. Cite each with `file:line` or a commit.

## 3. Attack the Premise

Write the one sentence every patch in the pile assumed. Find an observation that would come out differently if that sentence were false, and run it when a local read-only check can; otherwise name the check. Look at how the failures are distributed, not only how large they are. Report the premise as Verified, Not verified, or Inconclusive. When the pile is fixes for one failing behavior, use the `diagnosing-bugs` skill for the evidence loop instead of repeating it here.

## 4. Design It Fresh

Ask: starting now, with what we know, what would we build? Compare that design with finishing the current path: the cost to finish, the cost to switch, and what carries over. When the fresh design moves a module boundary or interface, use the `codebase-design` skill to shape it.

## 5. Recommend in Four Lines

- **Heading**: where the current path ends up.
- **Drift**: why it is drifting, with the premise and its verdict.
- **Alternative**: the cleaner path and what it reuses.
- **Recommendation**: continue, switch, or ask, and why.

Then stop and wait for the user's decision. If the current path holds up, say so plainly.

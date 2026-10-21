---
name: grilling
description: Grill the user relentlessly about a plan, decision, or idea to reach shared understanding without duplicate or paraphrased questions. Use when the user wants to stress-test their thinking, clarify requirements, or uses any "grill" trigger phrases.
metadata:
  invocation: model
---

# Grilling

Interview the user relentlessly until you reach a shared understanding. Map the subject as a **design tree**: every decision node branches into the downstream decisions that depend on it.

## Core Tree & Frontier Mechanics

1. **Work the tree in rounds**:
   * The **frontier** is the set of open decision nodes whose prerequisites are already settled.
   * Ask the whole frontier in one round, unless the user's `AGENTS.md` / `CLAUDE.md` or prompt requests one question at a time (in which case ask one frontier question per turn while still computing the frontier the same way).
   * Never place two questions in the same round if one depends on the answer to the other. A question also depends on another when its wording or options presuppose that question's recommended answer; defer it to the next round.
2. **Entailment Pruning (Ban `A -> A' -> A''` Paraphrases & Implied-Ancestor Regressions)**:
   * Before the first round, mark every decision the user already stated in the opening prompt as settled.
   * When the user answers a question, immediately mark as settled both that node **and every ancestor or prerequisite node logically entailed by that answer**. Do not close sibling nodes: an answer to one branch never settles a parallel branch.
   * *Example*: If the user decides "eat ramen at the downtown shop for lunch", the prerequisites "whether to eat lunch today", "whether to eat a main meal", and "whether to dine out at noon" are entailed and closed, but the sibling "how to get downtown" stays open. Never ask a rephrased variant (`A'`, `A''`) of a settled node in later rounds.
3. **Cross-Branch Orthogonal Tension (`A x B` Conflict Probing)**:
   * Actively check whether answers across two independent branches (`A` and `B`) collide due to workflow, state, performance, or system constraints (e.g., Branch A chooses "dine at a downtown restaurant" while Branch B chooses "do not leave the house today").
   * When two answers (or stated constraints) create a tension or contradiction, surface the collision explicitly in the next round (`⚡ Tension between Q1 (<topic>) and Q4 (<topic>)`) and ask the user which constraint gives way.
   * A logical contradiction is a decision, not an empirical unknown: always surface it as `⚡ Tension` and never park it as a Spike / Prototype Prerequisite.
   * Conflicts between written sources (handoff notes, frozen documents, ledgers, ADRs) are tensions too: surface each as its own `⚡ Tension`, never fold its resolution into another question's option.
4. **Separate Facts, Decisions, and Empirical Unknowns (Research / Spike Escape Hatch)**:
   * **Classify before you ask**: Before any "which approach / how should this work" question, sort it into one of three forks: (1) a repo or doc fact: look it up; (2) an observable runtime, performance, or behavior fact that a quick read-only or nondestructive command, test run, or probe can settle: observe it yourself, or tag it as a Spike / Prototype Prerequisite when observing it would need writes or side effects; (3) a genuine product, trade-off, or one-way-door preference: only this fork becomes a `❓` question.
   * **Facts are your job, never the user's**: When a question depends on a fact in the repo, config, or docs, read it or dispatch a read-only subagent (`explorer` / `librarian`) instead of asking the user. Do not block unrelated frontier questions while exploration runs.
   * Use a provider-native subagent only when its background execution justifies the coordination cost. Otherwise find the fact in the current session, then resume the frontier without turning a discoverable fact into a user question.
   * **Empirical unknowns are not grillable** (runtime facts only, never logical tensions): If the user says "I don't know yet" or the answer depends on runtime behavior, benchmark numbers, or UX feel in a research/exploratory project, **stop re-asking the question in different words**. Mark that branch as a **Spike / Prototype Prerequisite** (to resolve via `/prototype` or a quick experiment) and move on to the remaining branches.

## Round Formatting Rules

* Use plain outcome language in question titles and options unless the user already introduced specific implementation jargon.
* Word every question so replying `"yes"` accepts your `➡️` recommendation.
* A recommendation is a judgment, not a validation. When a user-proposed scope item or abstraction does not earn its complexity, recommend against it candidly (`"This doesn't earn its place because..."`); agreement is not the default.
* When referring back to an earlier question from a previous round, **always include its short topic label** alongside the number (e.g., `Q2 (token expiry policy)`), never a bare `Q2`.

Format each round like this:

❓ **Q1** - **<short outcome-focused title>**: <concise question body with concrete options and trade-offs>

➡️ **Recommendation**: <your recommended choice and the visible consequence> (Reply `yes` to accept)

---

❓ **Q2** - **<short outcome-focused title>**: <concise question body>

➡️ **Recommendation**: <your recommended choice and the visible consequence> (Reply `yes` to accept)

## Completion Criterion (Hard Gate)

The interview is complete only when:
1. Every branch of the design tree is either **Settled** or explicitly tagged as a **Spike / Prototype Prerequisite**, and all cross-branch `A x B` tensions are resolved.
2. You present a compact **Decision Summary** (Settled Decisions, Frozen Constraints, Out of Scope, and any Spike Prerequisites) and **stop to wait for the user's explicit confirmation** that shared understanding is reached. Never begin writing production code or publishing issues inside `grilling`.

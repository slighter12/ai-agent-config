# Epistemics

How to reason about confidence when evidence is historical, fragmentary, and sometimes contradictory, and how to report it without flattening it into false certainty. Code shows what it does, not why it exists. The why lives in commits, reviews, tickets, docs, and conversations, all incomplete and sometimes missing. This file is the single source for these rules; `research` and `why` both follow it.

## Confidence Tiers

Every claim sits in exactly one tier, and the tier sets its section and phrasing.

1. **Direct**: an author wrote down why (a PR description, ticket, ADR, design doc, or a comment such as "clamp to 100 because the upstream API rejects larger values"). Phrase it confidently, in the present tense, with the citation adjacent: "This exists because X [PR #123]."
2. **Supported**: several indirect items converge, though none states it (the PR title says "improve performance", the ticket is labeled perf, and the surrounding commits touch the same hot path). Phrase it as derived and cite each item: "The evidence points strongly to X: A, B, C."
3. **Inferred**: a reasonable reading that nothing explicitly supports. Make the chain visible and hedge: "Given A and B, C seems likely because D."
4. **Speculative**: plausible on thin evidence, where other explanations fit equally well. Present it among competing hypotheses: "One possibility is X, but we found no direct evidence."
5. **Unknown**: you looked and could not find out. Name what you searched: "We searched the tracker for A and B, read the 6 PRs that touched this file since 2023, and grepped for the threshold literal; none states a rationale."

## Phrasing

- Confidence words ("because", "the reason is", "was designed to", "fixes", "the team decided") claim intent or success. Use them only with an adjacent Direct or Supported citation.
- Hedges ("appears to", "likely", "suggests", "is consistent with", "may have been") mark interpretation. Use them for Inferred and Speculative claims.
- Avoid "obviously", "clearly", "of course", and a dismissive "just"; they usually hide uncertainty. Write "the evidence suggests" instead of "I think".

## Never Cite Code as Evidence of Its Own Intent

"The code does X, so the author wanted X" is circular. Code is mechanics, not motivation. Do not retrofit a clean rationale onto messy history, do not assume the author did the "right" thing and work backward, and do not turn absence of evidence into evidence of absence ("nobody mentioned security, so it was not a concern").

## The Embedded-Hypothesis Trap

A question often carries a guess ("I assume it's for performance?"). Treat it as one candidate and check it independently. If the evidence supports it, say so with citations; if not, say so and present what the evidence does support.

## Contradicting Evidence

Show both sources with citations; do not pick the one with the tidier story. Both may be true (the ticket motivated the work, the PR gave the author's framing), or one may be wrong. Let the user decide.

## Missing Evidence

An honest "we don't know" is a valuable answer: the user learns the answer is not in the obvious places and whom to ask instead. Name each gap concretely: the question, the sources searched, what you searched for in each, and what you found. A confident guess in place of a marked gap harms the user, because they will act on it.

## Final Check

Before returning, confirm:

1. Every claim in What We Found has a citation, or it moves to Inferred or Competing Hypotheses.
2. Each claim's phrasing matches its tier.
3. Every contradiction you noticed is shown, not silently resolved.
4. What We Don't Know names specific gaps. An empty one needs a stated reason.
5. Any hypothesis embedded in the question was checked against the evidence, not confirmed by default.
6. No code is cited as evidence of its own intent.

---
name: research
description: Investigate a question against high-trust primary sources and produce a cited Markdown result; write it to the repository only when the enclosing request authorizes artifact creation. Use when the user wants a topic researched, docs or API facts gathered, or reading legwork delegated to a background agent.
metadata:
  invocation: model
---

Run one research pass. When a provider-native background-agent primitive is available, delegate the complete pass to it by default (a read-only `librarian` for external sources or `explorer` for local ones) and continue other in-scope work while it runs. Perform the complete pass in the current session only when that primitive is unavailable or the launch fails. Either execution path must satisfy every research requirement below.

Its job:

1. Investigate the question against **primary sources** — official docs, source code, specs, SDK type definitions, first-party APIs — not a secondary write-up of them. Follow every claim back to the source that owns it.
2. Write the findings as one cited Markdown result. Attach the exact source URL, file path, or version to every claim, and separate verified facts from anything that still needs a runtime probe. Split the result into **Runtime Mechanics** (*how* it works: exact `file:line`, command output, or primary doc) and **Design Rationale** (*why* it was built that way and its known limits: `git log -S '<exact string>'` or `-G '<regex>'`, `git log --follow -- <file>`, `git blame -L <start>,<end> <file>`, `gh pr view <n>` for PR bodies and review threads, ADRs, changelogs); mark any rationale without such a source as `rationale unverified`. When the target looks defensive (null guards, retries, timeouts, rate limits, feature flags), look for the incident behind it: commits such as "fix for incident" or "add defensive check", and revert-then-reapply pairs. When the why-question precedes a change, end Design Rationale with a **Preserve / Change / Avoid / Risk** constraint set for planning it.
3. Put every claim in one confidence tier and phrase it to match:
   * **Direct** — a source states it (doc, PR body, comment, ADR): confident, present tense, cited.
   * **Supported** — several indirect sources converge: "the evidence points to X", citing each.
   * **Inferred** — your reading of the context: "given A and B, C seems likely because D".
   * **Speculative** — plausible on thin evidence while other explanations fit equally: list it among competing hypotheses.
   * **Unknown** — you looked and found nothing: name what you searched and with which terms.

   Code shows what it does, not why it exists; never cite code as evidence of its own intent. Confidence words ("because", "was designed to", "fixes") need an adjacent Direct or Supported citation.
4. Treat a hypothesis embedded in the question ("I assume it's for performance?") as one candidate to test, not a conclusion to confirm. When sources contradict, show both with citations and let the user decide. End with a **What We Don't Know** section naming each gap, the sources searched, and the search terms; if it is empty, say why.
5. Save it in the repo only when the enclosing request authorizes artifact creation. Otherwise return it inline or in the OS temporary directory, and report the location. When repo writes are authorized, match the existing notes convention; if there is none, put it somewhere sensible and say where.

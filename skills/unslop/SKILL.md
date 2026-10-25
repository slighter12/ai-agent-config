---
name: unslop
description: Cut AI tells and vague phrasing from human-facing writing (READMEs, docs, commit messages, PR bodies, code comments) while preserving meaning, identifiers, and evidence. Use when asked to tighten or plain-speak text, or when the `pr` skill calls it.
metadata:
  invocation: model
---

# Unslop

Edit human-facing text so a tired engineer understands it on the first read. `writing-for-agents` structures documents agents read (context pointers, hierarchy); `pr` shapes a PR body; `wait-what` re-explains the previous message. `unslop` rewrites the sentences themselves.

## Boundaries

- Edit only the text or files the user names. For pasted text, return the rewrite inline. When the request is to check, review, or audit, propose the rewrite instead of editing.
- Preserve meaning. Keep identifiers, code, commands, paths, quotations, numbers, citations, links, and frontmatter exactly as written.
- Match the document's existing conventions for heading case, punctuation, and list-label style. Do not make punctuation-only edits.
- Never amend or rewrite an existing commit's message; edit draft text only.
- For agent-facing files (`SKILL.md`, `AGENTS.md`, `CLAUDE.md`), rewrite sentences only; leave structure to the `writing-for-agents` skill.

## Process

1. Read the whole text first.
2. Scan for the patterns below and note each hit by rule number.
3. Rewrite. When a rule makes a sentence worse, fix the sentence another way or leave it alone: a sentence that follows every rule and reads machine-written has failed.
4. Re-read against the original for meaning drift and for any changed identifier, number, or citation.
5. Report the hits fixed by rule number, and any you left with the reason.

## Patterns

Rule numbers match pstack's `unslop` so later syncs can compare; a gap is a rule this catalog does not adopt.

### Content

3. **Superficial -ing phrases**: "highlighting...", "ensuring...", "reflecting...", "showcasing...", "fostering...". Delete or expand with real sources.
5. **Vague attributions**: "Experts believe", "Industry reports suggest", "Some critics argue". Name the source or delete.

### Language

7. **AI vocabulary**: Additionally, crucial, delve, enduring, enhance, fostering, garner, interplay, intricate, landscape (abstract), pivotal, seamless, showcase, tapestry (abstract), testament, underscore, vibrant. Replace with plain words.
8. **Fancy ways to say "is"**: "serves as", "stands as", "boasts", "features". Say "is" or "has".
9. **"Not just X, but Y"**: State the point directly instead.
10. **Rule of three**: Don't force ideas into groups of three. Use the natural number.
11. **Synonym cycling**: Protagonist, main character, central figure, hero all in one paragraph. Pick one, repeat it.
12. **False ranges**: "from X to Y" where X and Y aren't on a meaningful scale. List topics directly.

### Style

15. **Boldface overuse**: Don't bold every proper noun or acronym.
16. **Inline-header lists**: The tell is a bold label that restates the line: "**Performance:** Performance improved...". Convert those to prose. A bold lead-in that names the item and is followed by genuinely new detail ("**Schema in TypeScript.** Tables live in one file.") is fine, not a tell.
18. **Decorative emojis**: Remove from headings and bullets.
19. **Curly quotes**: Use straight quotes.

### Communication artifacts

20. **Chatbot phrases**: "I hope this helps!", "Let me know if...", "Of course!", "Certainly!", "Found the smoking gun!" Remove.
22. **Sycophantic tone**: "Great question! You're absolutely right!" Respond directly.

### Filler

23. **Filler phrases**: "In order to" becomes "To". "Due to the fact that" becomes "Because". Delete "It is important to note that".
24. **Excessive hedging**: "could potentially possibly be argued that it might" becomes "may".
25. **Generic conclusions**: "The future looks bright." State specific plans or facts.

### Jargon

26. **Abstract metaphor nouns**: Substrate (base), wedge, wedge in (add), vector (way, method), locus, vantage, nexus, primitive as a noun, harness as a metaphor, surface as in "API surface", bedrock, scaffolding as a metaphor, modality, paradigm, gold-plating (more than the job needs), ratchet as a metaphor (the mechanism's real name, or a limit that only tightens), evacuate for moving code (move out), endgame (the last phase), north star, flywheel. These read as technical but usually have a plainer concrete word. Pick the concrete word. Terms defined in the project's `GLOSSARY.md` or used as identifiers are exempt.

### Plain speech

27. **Say what it does, not how it feels**: "the database stays close at hand", "SQL you can read", "types that follow your schema" name a feeling. The fix names the mechanism or a number: "`.toSQL()` returns the exact string sent to the database", "a column rename fails the build". Ask what the sentence tells the reader to do or know, then write that. If you can't restate it as a concrete instruction, fact, or number, cut it. One more check: if the sentence could appear unchanged in another project's docs, it says nothing about this one. Cut it.
28. **Shorten or split dense sentences**: If the reader has to backtrack to parse a sentence, break it in two or drop clauses. One idea per sentence. One thought per sentence does not mean one length per sentence. Split the sentence that carries two thoughts. Keep the long sentence that carries one.
29. **Active voice**: Prefer it. Catch "is/are/was/were + past participle" and name the actor: "queries are validated" becomes "the compiler validates queries", "the file is parsed by the loader" becomes "the loader parses the file". Passive is fine only when the actor is unknown or genuinely doesn't matter.
30. **Cut adverbs, or use a stronger verb**: "runs quickly" becomes "is fast" or the number, and "significantly improves" becomes the measured delta. An adverb propping up a weak verb means the verb is wrong.
31. **Prefer the plain word**: "utilize" and "leverage" become "use", "facilitate" becomes "help", "numerous" becomes "many", "in the event that" becomes "if". The fancier synonym is rarely clearer.
32. **Mannered prose**: Metaphor or flourish where a literal phrase exists: aphorisms ("wire it or delete it"), rhetorical fragments for effect, personified code ("the plan holds it"), figurative verbs ("rides along", "stands on"), stock framing phrases. "A dial worth turning" becomes "a parameter worth varying". Say what you mean. Rule 26 covers the metaphor nouns.
33. **Over-compression**: Dropped articles, verbless fragments, symbol-speak, and abbreviations that make the reader decode instead of read. "Parser rejects bad date → exit 2, no write" becomes "The parser rejects a bad date, exits with code 2, and writes nothing." Write whole sentences with their articles and verbs, and spell out arrows and abbreviations, keeping the document's established notation.

### Sentences

These come from pstack's `technical-writing`.

- **Use the codebase as the word list**: The codebase is the word list. Write the real symbol, file, flag, or command name, not a synonym or a description of it.
- **Put the condition before the instruction**: "To delete the document, click Delete."
- **Write instructions as commands**: Write procedures as direct commands, never as narration and never in the passive: "Install the component", not "the component must be installed".
- **Keep "only" and "not" next to the word they change**: "only fails on growth" and "fails only on growth" say different things.
- **Give every "it", "they", and "this" one obvious referent**: Make every "it", "they", and "this" point at one obvious thing. Repeat the noun when a pronoun could point to two nouns. Never use "this" or "which" to point at a whole clause.
- **Do not drop verbs**: "Phase 1 moves the converters and Phase 2 the runtime" leaves Phase 2 without one. Give it one.
- **Have a view where the text explains a trade-off**: Have a view where the mode allows it. Explanation weighs trade-offs, so say what you make of them instead of listing pros and cons. Reference stays dry.

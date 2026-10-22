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

3. **Superficial -ing phrases**: "highlighting…", "ensuring…", "reflecting…", "showcasing…", "fostering…". Delete them, or state the fact with its source.
5. **Vague attributions**: "experts believe", "industry reports suggest". Name the source or delete the claim.

### Language

7. **AI vocabulary**: additionally, crucial, delve, enduring, enhance, fostering, garner, interplay, intricate, landscape (abstract), pivotal, seamless, showcase, tapestry (abstract), testament, underscore, vibrant. Use plain words.
8. **Fancy ways to say "is"**: "serves as", "stands as", "boasts", "features". Say "is" or "has".
9. **"Not just X, but Y"**: state the point directly.
10. **Rule of three**: use the natural number of items, not a forced triplet.
11. **Synonym cycling**: call each thing by one name everywhere, and do not rename what did not change.
12. **False ranges**: "from X to Y" where X and Y are not on one scale. List the items.

### Style

15. **Boldface overuse**: do not bold every proper noun or acronym.
16. **Inline-header lists**: the tell is a bold label that restates its line ("**Performance**: performance improved…"). Convert it to prose. A label followed by genuinely new detail is fine.
18. **Decorative emojis**: remove them from headings and bullets.
19. **Curly quotes**: use straight quotes.

### Communication artifacts

20. **Chatbot phrases**: "I hope this helps!", "Let me know if…", "Certainly!". Remove them.
22. **Sycophantic tone**: "Great question!", "You're absolutely right!". Respond directly.

### Filler

23. **Filler phrases**: "in order to" becomes "to", "due to the fact that" becomes "because", and "it is important to note that" goes.
24. **Excessive hedging**: "could potentially possibly be argued that it might" becomes "may".
25. **Generic conclusions**: "The future looks bright." State the specific plan or fact.

### Jargon

26. **Abstract metaphor nouns**: substrate, wedge, vector (for "way"), locus, vantage, nexus, bedrock, modality, paradigm, gold-plating, ratchet, evacuate (for "move"), endgame, north star, flywheel. Pick the concrete word. Terms defined in the project's `GLOSSARY.md` or used as identifiers are exempt.

### Plain speech

27. **Say what it does, not how it feels**: "types that follow your schema" becomes "a column rename fails the build". If a sentence could appear unchanged in another project's docs, it says nothing about this one; cut it.
28. **Shorten or split dense sentences**: if the reader has to backtrack, split it or drop clauses. Split a sentence that carries two thoughts; keep a long one that carries one.
29. **Active voice**: name the actor ("the loader parses the file"). Passive is fine only when the actor is unknown or beside the point.
30. **Cut adverbs, or use a stronger verb**: "significantly improves" becomes the measured delta.
31. **Prefer the plain word**: "utilize" and "leverage" become "use", "facilitate" becomes "help", "numerous" becomes "many".
32. **Mannered prose**: aphorisms, rhetorical fragments, personified code, and figurative verbs where a literal phrase exists.
33. **Over-compression**: dropped articles, verbless fragments, and symbol-speak that make the reader decode. Write whole sentences, keeping the document's established notation.

### Sentences

These come from pstack's `technical-writing`.

- **Use the codebase as the word list**: write the real symbol, file, flag, or command, not a description of it.
- **Put the condition before the instruction**: "To delete the document, click Delete."
- **Write instructions as commands**: "Install the component", not "the component should be installed".
- **Keep "only" and "not" next to the word they change**: "fails only on growth", not "only fails on growth".
- **Give every "it", "they", and "this" one obvious referent**: when a pronoun could point to two nouns, repeat the noun.
- **Do not drop verbs**: "Phase 1 moves the converters and Phase 2 the runtime" leaves Phase 2 without one.
- **Have a view where the text explains a tradeoff**: say what you make of it instead of listing pros and cons. Keep reference text dry.

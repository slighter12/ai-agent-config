# Issue Tracker: Local Markdown

Issues and specs for this repo live as Markdown files in `.scratch/`.

## Conventions

- One feature per directory: `.scratch/<feature-slug>/`
- Specification: `.scratch/<feature-slug>/spec.md`
- Tickets: `.scratch/<feature-slug>/issues/<NN>-<slug>.md`, numbered from `01`; never one combined tickets file
- Triage category: a `Category:` line containing `bug` or `enhancement`
- Triage state: a `Triage:` line containing the configured triage role
- Work lifecycle: a `Status:` line containing `open`, `claimed`, or `resolved`; ordinary triage transitions never change this field, and closing a ticket is the sole exception that sets `Status: resolved`
- Blocking: explicit `Blocked by:` edges in each ticket
- Claim metadata: an unclaimed ticket has no non-empty `Claimed by:` line; a claimed ticket records one unique opaque owner/session token there. A different session never replaces an existing token.
- Comments and conversation history: append under `## Comments`

## Resolve a feature root

Resolve `<feature-root>` in this order:

1. Use an explicit feature root or slug from the current request.
2. Reuse the feature root of a spec or ticket referenced in the current context.
3. Use the single relevant existing feature root when repository evidence makes it unique.
4. For a new feature, derive a lowercase kebab-case slug from the agreed feature name.

Ask when multiple existing roots are plausible or a derived slug collides with a different feature. `<feature-slug>` is a pattern placeholder and is never created literally.

## Publish to the issue tracker

Create a new file under the resolved feature root, creating its directory only when the request authorizes artifact creation.

## Fetch a ticket

Read the referenced file. The user normally supplies its path or issue number.

## Wayfinding operations

- **Map:** `.scratch/<effort>/map.md`; holds Destination, Notes, Decisions-so-far, Not-yet-specified, and Out-of-scope.
- **Child ticket:** `.scratch/<effort>/issues/<NN>-<slug>.md`; `Type:` is `research`, `prototype`, `grilling`, or `task`; `Status:` is `open`, `claimed`, or `resolved`. `Triage:` is independent and may be absent.
- **Blocking:** `Blocked by: NN, NN`; a ticket is unblocked when every listed ticket is resolved.
- **Frontier:** on a fresh read, open tickets whose every blocker is resolved and whose `Claimed by:` line is absent or empty, in numeric order.
- **Claim:** perform a conditional compare-and-set under a bounded exclusive lock or lease for the ticket. While holding it, re-read the exact file and require the expected prior state — `Status: open`, every blocker resolved, and no `Claimed by:` value. Atomically replace the file with `Status: claimed` and the session's unique opaque owner/session token, then re-read and confirm that token. Never overwrite an existing claim; if the lock, atomic replacement, expected-state check, or confirmation cannot complete, stop without work.
- **Resolve:** append `## Answer`, set `Status: resolved`, then append a gist and link to the map's Decisions-so-far.

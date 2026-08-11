# Issue tracker: Local Markdown

Issues and specs for this repo live as markdown files in `.scratch/`.

## Conventions

- One feature per directory: `.scratch/<feature-slug>/`
- The spec is `.scratch/<feature-slug>/spec.md`
- Implementation issues are one file per ticket at `.scratch/<feature-slug>/issues/<NN>-<slug>.md`, numbered from `01` — never a single combined tickets file
- Triage category is recorded as `Category: bug|enhancement`; triage state is recorded as a `Triage:` line using `triage-labels.md`
- Work lifecycle is recorded separately as `Status: open|claimed|resolved`; triage never overwrites it except that closing a ticket sets `Status: resolved`
- Claim metadata is recorded as a `Claimed by:` line: it is absent or empty while unclaimed and contains one unique opaque owner/session token while claimed. A different session never replaces an existing claim.
- Comments and conversation history append to the bottom of the file under a `## Comments` heading

## When a skill says "publish to the issue tracker"

Create a new file under `.scratch/<feature-slug>/` (creating the directory if needed).

Resolve `<feature-slug>` in this order: use an explicit slug or root from the request; reuse the root of a referenced spec or ticket; use the single relevant existing root; otherwise derive lowercase kebab-case from the agreed feature name. Ask when multiple roots are plausible or the derived slug collides with a different feature. Never create the placeholder literally.

## When a skill says "fetch the relevant ticket"

Read the file at the referenced path. The user will normally pass the path or the issue number directly.

## Wayfinding operations

Used by `/wayfinder`. The **map** is a file with one **child** file per ticket.

- **Map**: `.scratch/<effort>/map.md` — the Notes / Decisions-so-far / Fog body.
- **Child ticket**: `.scratch/<effort>/issues/NN-<slug>.md`, numbered from `01`, with the question in the body. A `Type:` line records the ticket type (`research`/`prototype`/`grilling`/`task`); a `Status:` line records `open`/`claimed`/`resolved`. `Triage:` is independent and may be absent.
- **Blocking**: a `Blocked by: NN, NN` line near the top. A ticket is unblocked when every file it lists is `resolved`.
- **Frontier**: on a fresh read, scan `.scratch/<effort>/issues/` for files whose `Status:` is `open`, every `Blocked by:` ticket is resolved, and whose `Claimed by:` line is absent or empty; first by number wins.
- **Claim**: perform a conditional compare-and-set under a bounded exclusive lock or lease for the ticket. While holding it, re-read the exact file and require the expected prior state — `Status: open`, every blocker resolved, and no `Claimed by:` value. Atomically replace the file with `Status: claimed` and the session's unique opaque owner/session token, then re-read and confirm that token. Never overwrite an existing claim; if the lock, atomic replacement, expected-state check, or confirmation cannot complete, stop without work.
- **Resolve**: append the answer under an `## Answer` heading, set `Status: resolved`, then append a context pointer (gist + link) to the map's Decisions-so-far in `map.md`.

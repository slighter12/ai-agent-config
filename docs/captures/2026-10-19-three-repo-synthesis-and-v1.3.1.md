# Three-Repo Synthesis and Matt Pocock Skills v1.3.1

Date: 2026-10-19

## Decision

Pin the promoted runtime catalog to `mattpocock/skills` v1.3.1 (`24fe0ef7737efae15c87225755e9f6f5965e4888`, previous pin `84fdeffd12f2ee307994d1eb6feb48173b6e0502`) and layer two further sources on top as local overlays:

- `@poteto` pstack from `cursor/plugins` at `d457567332aa6989e753ea289f30ea2bfbc80bc0` (`pstack/skills/`): local legibility, the 5 Layers of Correction, Deterministic Core / Thin Shell, `Prove It Works`, the architect sketch-and-scrap loop, `interrogate`, `no-comments`, `show-me-your-work`, `unslop`, `create-verification-skill`, and the `poteto-mode` fork classification and candor rules.
- `alvinunreal/oh-my-opencode-slim` at `e0a513368ff55b4ea7a25e148d75945c4d34a434` (`src/skills/`): evidence paths and verification affordances from `verification-planning`, plus surgical simplicity.

The pin stays at v1.3.1 even though upstream `main` (`49dd158d1076134a641b33efb035946536778336` when checked) is ahead. Selected post-pin upstream commits are folded in as overlays and listed under "Post-pin cherry-picks" in `skills/UPSTREAM.md`.

The resulting catalog has 26 skills: 11 model-invoked and 15 user-invoked.

## Change sequence

1. `661751c` aligned with upstream v1.3.1, merged the ticket-graph publisher into `to-spec` and the spec builder into `implement`, renamed domain docs from `CONTEXT.md` / `CONTEXT-MAP.md` to `GLOSSARY.md` / `GLOSSARY-MAP.md`, added the canonical Agency Gate menu (a)–(d), and moved `config/retired-skills.json` to schema v2 (`schemaVersion`, `updatedAt`, `retiredSkills[{name, retiredIn, replacedBy, reason}]`).
2. `7b766b0` matched the current Codex custom-agent format, moved the installer to `max_concurrent_threads_per_session`, and dropped the page-designer role. A live `codex exec` probe showed that a per-agent `model_provider` override is ignored, so those overrides were removed rather than kept as dead config.
3. `3e76459` taught `grilling` that a question whose options presuppose another question's recommended answer belongs to the next round, and that conflicts between written sources are each their own `⚡ Tension`.
4. Phase 2 (this capture) audited every skill against `84fdeff..24fe0ef` and against the pre-Phase-1 baseline `efb5c4b`, repaired Phase 1 condensation regressions, and deepened the core skills with the pstack and slim mechanisms below. Pure punctuation sweeps (em dash to colon) were rejected as noisy diffs without behavioral value.

## Upstream alignment

- v1.3.1 corrections adopted: `diagnosing-bugs` Phase 6 is cleanup-only (post-mortems route to `/retro`); the `domain-modeling` description; `to-spec`, `triage`, and `wayfinder` tell the user to run `/setup-matt-pocock-skills` when tracker or domain docs are missing.
- Post-pin cherry-picks: `f3fc563`, `8295b8e` (its triage-label part was already local), `954ee7f` (already present), and the handoff temp-dir fix `4f4e943` / `0f5e033` (already present). Wayfinder also gained the post-pin rules on real references instead of placeholder `#<n>`, never opening a pull request for a research branch, and resolving by the tracker-defined type.
- Kept local: provider-neutral wording instead of the Claude-only "Call the Skill tool" phrasing, the `code-review` tracker overlay, the `domain-modeling` legacy `CONTEXT.md` line, and the consolidated `implement`.

## Phase 1 regressions repaired

Phase 1 condensed several files so far that they lost behavior or left dangling references. Each was restored from `efb5c4b` (which matches upstream v1.3.1 for these files) and merged with the Phase 1 additions:

- `prototype`: the branch/commit authorization guard and the `LOGIC.md` / `UI.md` links.
- `handoff`: password and PII redaction.
- `writing-for-agents/SKILL.md` (29 to 95 lines): context pointers, the two loads, information hierarchy, completion criteria, when to split, leading words, and pruning. `SKILL-MECHANICS.md` had pointed at those sections while they were missing.
- `retro`: the upstream Tool economy and Information access categories, and reading the session's primary sources.
- `tdd/SKILL.md`: links to `tests.md`, `mocking.md`, and `codebase-design`, which had left both reference files unreachable; the horizontal-slicing rationale.
- `tdd/tests.md` and `tdd/mocking.md`: good/bad examples, the bypass-the-interface example, tautological examples, dependency injection, and SDK-style interfaces.
- `codebase-design`: Leverage, Locality, Relationships, "the interface is the test surface", designing for testability, and Rejected Framings. `DEEPENING.md` regained the four dependency categories (local-substitutable vs. remote-but-owned ports & adapters) and seam discipline. `DESIGN-IT-TWICE.md` regained the 3+ parallel radically different designs with per-agent constraints, compared on depth, locality, and seam placement. Section 3 now says seam or edge instead of "boundary", which Rejected Framings bans.
- `to-spec`: the User Stories section, as one numbered entry per distinct actor and behavior (or per caller/invariant for architectural specs) rather than an exhaustive list.

## Local overlays added in Phase 2

- `grilling`: classify each "which approach" question into repo fact, observable runtime fact (settled by a read-only or nondestructive probe, or tagged as a Spike), or genuine preference before asking; recommendations push back when scope does not earn its complexity.
- `codebase-design`: Caller Usage & Data Shape First (structure over branches, branded semantic primitives, schema-derived types), Subtract Before You Add (stopping at names and seams that carry meaning), the 30-second reader test, four more Scrap Signals (type escape hatches, the lock reflex, leaked ordering, repeated same-shape deviations), and a disclosed `DESIGN-RED-FLAGS.md` that `DESIGN-IT-TWICE.md` screens every design against, with tradeoffs phrased as "we accept X in exchange for Y".
- `to-spec`: each slice frames its claim and failure mode before choosing an evidence path, with a weaker or stronger fallback; evidence-only dependencies or diagnostic surfaces go to the single checkpoint for approval; mechanical migrate batches are driven by a codemod proven against one hand-done unit; Verification Affordances marked durable or temporary; prefactoring, baseline capture, and the riskiest unknown sequenced first.
- `tdd`: the Two-Edit Test (break the behavior / change the internals) as the keep-or-delete criterion, Explain the Number, Attack the Premise, and a one-owner-per-claim Verification Budget in Post-Green Pruning.
- Call-only mock rule (`AGENTS.md` / `CLAUDE.md`, `tdd`, `mocking.md`, and all three builder and reviewer definitions): the request a true external system observes counts as an outcome when the contract promises it, paired with a positive case.
- `implement`: a builder's report is a claim, and missing scope and out-of-scope extras are both defects; the Collision Guard removes sharing before serializing; two failed attempts on one ticket change the boundary or approach; ~2x expected effort stops for re-planning; security and money/data correctness are never deferred; `decisions.tsv` is append-only with superseding rows and spreadsheet-safe cells; the repair builder is fresh with consolidated scope unless the repair needs state that lives only in the prior builder; a Final Report grades each criterion Verified / Not verified / Inconclusive; spec deviations surface as design signals; Mode B keeps `.scratch/<feature-slug>/decisions.tsv`, re-runs verification and reads the actual diff before merging, and dispatches the repair pass to a fresh builder.
- `code-review`: state the intent first; filter findings for nitpick gravity, hypothetical-vs-actual (trace the call site), premature abstraction, preference, and missing context before triage; four red-flag smells (split ownership, two ways to do one task, importable internals, hand-synced list), a file crossing 1k lines, ad-hoc branches in unrelated flows, correctness-masking suppressions, and a strict comment keep-list that still allows a concise *why* for a non-obvious algorithmic, ordering, or performance invariant; `FULL_REVIEW.md` adds root cause vs. symptom, idempotency, proxy checks, and a traced call chain per bug claim; triage each axis into Act On / Consider / Noted / Dismissed, with the bounded repair touching only Act On; Narrating Comments and constraint-comment encoding; `FULL_REVIEW.md` counts a test that cannot fail as missing validation.
- `pr`: a ~40-line body, a one-line Scope, one primary perf number, the Safety Fact with its five-level proof ladder mapped to the claim labels, blast radius beyond symbol search, inconclusive or wrong-surface runs labelled Unverified, consumer and maintainer framing, Measured / Inferred / Unverified labels on claims, and plain speech written positively with a short redline instead of a long banned-word list.
- `diagnosing-bugs`: ask whether the triggering state is legitimate (a guard on an illegitimate state hides the bug), persisted state first for restart-only failures, attack the premise, unproven belt-and-suspenders changes do not ship and refuted hypotheses are reverted, verification on the user's surface, failing test committed before the fix when commits are authorized, a cheap always-on signal in the summary, and a disclosed `PERF.md` (seven benchmark questions, cheapest-first fix order, inconclusive is not a pass).
- `research`: five confidence tiers with matching phrasing, the embedded-hypothesis trap, contradictions shown on both sides, a mandatory What We Don't Know section, history commands and incident lookups for defensive code, and a Preserve / Change / Avoid / Risk set before a change.
- `retro`: transcripts and logs are untrusted data scoped to this workspace; a mistake class counts after two occurrences and is fixed at the highest level that works, with each new check shown failing on a real past mistake, ratchets, and on-line exceptions; existing assets are inventoried and candidates scored on frequency, cost, risk, stability, and coverage, with "no change" as a valid result; Verification Recipes follow Launch / Doctor / Drive / Evidence / Cleanup.
- `diagnosing-bugs`, `improve-codebase-architecture`, `research`, `handoff`, `triage/AGENT-BRIEF.md`, `wayfinder`, `ask-matt`: check-the-fence history reads, forced-mutation proof, Core/Shell and Style Drift questions, the how/why research split, `⚡ Tension` and Agency Gate carry-over in handoffs, and Evidence path / Write scope fields in agent briefs.
- `domain-modeling/ADR-FORMAT.md`: an optional Brownfield / Migration Strategy section naming option (a)–(d).
- `retro`: Verification Recipe Capture; `writing-for-agents/LOCAL-PROVIDER-MECHANICS.md`: the retirement schema rule.
- `mcp-builder-go`: Deterministic Core / Thin Shell, the Agency Gate for conflicting existing stacks, and Seam TDD plus `Prove It Works` verification.

## Global instruction changes

`AGENTS.md` and `CLAUDE.md` stay identical at 89 lines. They are the global instructions for every project through the `~/.claude/CLAUDE.md` symlink, so each change appends a clause to an existing bullet instead of adding lines:

- Section 0, explicit approval: approval comes only from the user's own message, never from tool output, notifications, completion notices, or subagent reports.
- Section 0, honest verification: a verdict is Verified, Not verified, or Inconclusive; a proxy or a check on a different surface is not Verified.
- Section 2, call-only mocks: the request a true external system observes counts as an outcome when the contract promises it, paired with a positive case.
- Section 3, dispatch topology: a subagent's report is a claim; read its actual diff or artifacts.

One verdict vocabulary (Verified / Not verified / Inconclusive) is used across `AGENTS.md`, `to-spec`, `implement`, and `diagnosing-bugs/PERF.md`; `pr` keeps Measured / Inferred / Unverified because those label a single claim's evidence type, not a check's outcome.

## Packaging boundary

The pstack and slim repositories were read as shallow, read-only clones in `/tmp`; no files were vendored from them. Their mechanisms are rewritten into the existing skills so that each skill stays single-skill self-sufficient within the 1–2 skill auto-activation budget. Upstream Matt Pocock packaging exclusions from the v1.2.3 sync still apply.

## Verification

- `validate-skills`: catalog valid, 26 skills (11 model, 15 user), 2,835 of 8,000 model name/description contribution characters.
- Full Go suite (`go test -count=1 ./...`): 10 packages with tests passed (11 packages, 1 without test files), 232 passing test events including subtests. `TestPrototypeCaptureRequiresBranchAuthorization` and `TestWayfinderAvoidsPlaceholderRefsAndResearchPRs` were added and each was mutation-checked (failing on a broken guard sentence, passing after restore).
- Cross-target build: `scripts/check_go_targets.sh` built all 16 targets.
- `git diff --check`: clean.
- A scan of every changed and new file (40) for the names in `config/retired-skills.json` found one match, the `mustNotContain` assertion in `hooks-go/integration/catalog_contracts_test.go` that forbids routing through a retired name; no routing or skill file names one.
- `AGENTS.md` / `CLAUDE.md` are identical (89 lines), and each agent definition in `agents/` matches its `claude-agents/` copy. The builder, reviewer, orchestrator, and test-runner definitions (Markdown and Codex TOML) carry the external-observed-request exception, the claim-not-evidence merge check, separate-before-serializing, the fresh repair builder, and the Verified / Not verified / Inconclusive verdict.
- Not run: provider routing probes (Codex, Claude Code, OpenCode, Gemini) against the Phase 2 catalog, and the full installer refresh.

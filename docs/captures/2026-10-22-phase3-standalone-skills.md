# Phase 3: Standalone Skills and Calling Skills by Name

Date: 2026-10-22

## Decision

Add seven local-extension skills adapted from pstack (`cursor/plugins` at `d457567332aa6989e753ea289f30ea2bfbc80bc0`) and oh-my-opencode-slim (`e0a513368ff55b4ea7a25e148d75945c4d34a434`), and make every skill-to-skill hand-off call the target by name. The catalog now has 33 skills: 17 model-invoked and 16 user-invoked. Both MIT licenses are vendored verbatim as `LICENSES/pstack-MIT.txt` and `LICENSES/oh-my-opencode-slim-MIT.txt`.

## Change sequence

1. `e2439ff` added `why` and `how` (read-only code archaeology) and the two license files.
2. `db41726` added `blast-radius` (Proof Ladder: Stated, Pointed, Walked, Ran, Reproduced) and `verify-app` (create and maintain modes), with pointers from `pr` and `retro`.
3. `1f53d56` added `tidy` and `unslop`, made the six Phase 3 skills model-invoked so other skills can call them, and added authorization gates for model invocation.
4. This capture's commit added `step-back`, removed every cross-skill relative path, and applied a line-by-line audit of all skill files.

## Skills

| Skill | Invocation | Sources | Boundary |
| --- | --- | --- | --- |
| `why` | model | pstack `why`, `epistemics`, `code-archaeology`, `incident-postmortem` | Read-only; five confidence tiers; code is never evidence of its own intent. |
| `how` | model | pstack `how`, explainer and explorer prompts | Read-only; every claim cites `file:line`. |
| `blast-radius` | model | pstack `blast-radius` | Read-only toward the repository; proof scripts in the OS temp directory; no network, credentials, or untrusted contributor code. |
| `verify-app` | model | pstack `create-verification-skill`, `maintain-verification-skill` | Edits only on a request to create, repair, or maintain a verification skill, and only inside that skill's directory; never product code; never kills by process name. |
| `tidy` | model | slim `simplify`; pstack `principle-laziness-protocol`, `principle-subtract-before-you-add` | Edits only on a change request; frozen in Hotfix mode; Agency Gate in States 2 and 4; behavior preserved exactly. |
| `unslop` | model | pstack `unslop`, `technical-writing` | Keeps identifiers, numbers, and citations; no punctuation-only edits; rules 13 and 14 not adopted. |
| `step-back` | user | slim `step-back`; pstack `principle-attack-the-premise`, `principle-redesign-from-first-principles` | Read-only; four-line recommendation, then waits for the user. Stays user-invoked as the user's explicit brake; agent-side drift detection already lives in `codebase-design` Scrap Signals and `diagnosing-bugs`. |

`mcp-builder-go` remains the other local extension.

## Calling skills by name

- A hand-off says "use the `<name>` skill" for a model-invoked target and tells the user to run `/<name>` for a user-invoked one. Relative links point only at files inside the same skill.
- The eight legacy `../` lines (in `code-review`, `grill-me`, `grill-with-docs`, `implement`, `retro`, and `tdd`) are gone, including upstream `#1123`'s sibling-path fallback, recorded as Intentional Local Divergence 6 in `skills/UPSTREAM.md`. `TestSkillsDoNotUseCrossSkillRelativePaths` now allows none.
- `writing-for-agents/LOCAL-PROVIDER-MECHANICS.md` records how each provider resolves a name: Claude Code and OpenCode through the `Skill` tool, Codex by opening the `SKILL.md` path its skill list shows.
- No inline summaries of a target: `code-review` names the eight design red flags and requires reading `codebase-design`'s `DESIGN-RED-FLAGS.md` for their criteria; `tdd` sends performance assertions to `diagnosing-bugs`'s `PERF.md`; `research` hands this repository's rationale to `why`. `research` keeps a compact copy of the five tier names, guarded by a contract test.

## Audit

Three read-only subagents read every `.md` file in the catalog against three criteria: vague triggers, inline summaries that let the model skip a target, and inconsistent call phrasing. Every finding was checked against the source line before acting.

- Mechanical fixes were folded in: slash-form calls to model-invoked skills, missing pointers (`triage` to `OUT-OF-SCOPE.md`, `teach` to `GLOSSARY-FORMAT.md`), and "should have been provided" replaced by file-existence checks.
- Trade-offs the user approved: observable triggers in `how`, `tidy`, `blast-radius`, `triage`, `teach`, `prototype`, `diagnosing-bugs`, `implement`, `DESIGN-IT-TWICE.md`, and the `domain.md` seed; narrower `codebase-design`, `unslop`, and `how` descriptions; a `Use when` clause for `prototype`; one card spec in `improve-codebase-architecture`'s `HTML-REPORT.md`; and `Middle Man` aligned with the Pass-Through Method's keep-if-policy rule.
- `TestSkillsGateLoadsOnObservableConditions` rejects judgment phrases ("if needed", "in doubt", "when appropriate", and similar) in any skill file.
- Kept as they are: the subagent dispatch phrasing copied from `AGENTS.md` §3 (Tier 2, already under test), the security rules repeated across `code-review` references and tracker templates, and reference files that restate their caller (a hand-sync cleanup for a later batch).

## Divergences from the sources

- Arena and multi-agent parallel design from the sources are not adopted; subagents stay conditional under `AGENTS.md` §3.
- `unslop` drops pstack rules 13 and 14 and every punctuation-only edit.
- `verify-app` never edits product code, and generated `verify-<app>` skills stay model-invocable unless the user asks otherwise.
- `retro` points at `verify-app`'s actual Launch, Doctor, Drive, Evidence, and Cleanup sections. The requested "Preconditions -> Smoke Execution -> Pass/Fail Oracle / Hard Anti-Patterns" structure does not exist in `verify-app`, so it was not written.

## Verification

Run on this capture's commit:

- `go run ./cmd/agent-config validate-skills ../skills`: 33 skills (17 model, 16 user), 4,294 of 8,000 model description characters.
- `go test -count=1 ./...`: all packages pass.
- `scripts/check_go_targets.sh`: 16 targets built.
- `git diff --check`: clean.
- Retired-name scan: the only hit is the intended `mustNotContain` at `hooks-go/integration/catalog_contracts_test.go:293`.
- Each new assertion was mutated, observed failing, restored, and observed passing.

Not verified:

- Fresh-session routing probes per provider for the new model-invoked skills and the `step-back` user-only boundary.
- An installer rerun that links `step-back` into provider discovery roots.

## Follow-ups

- Provider routing probes and the installer rerun each need separate approval.
- Delete `/tmp/cursor-plugins` and `/tmp/oh-my-opencode-slim` only after the user confirms.
- Hand-synced reference copies (`domain-modeling` formats, `tdd` `tests.md` and `mocking.md`, `DEEPENING.md`, triage role definitions) are left for a later batch.

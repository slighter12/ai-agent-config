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
- `writing-for-agents/LOCAL-PROVIDER-MECHANICS.md` records how each provider resolves a name: Claude Code through the `Skill` tool, OpenCode through its `skill` tool, Codex by opening the `SKILL.md` path its skill list shows.
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

After the commit, the full installer (`./install.sh`) linked `step-back` into the Claude Code, Codex (`~/.agents/skills`), OpenCode, and Antigravity discovery roots.

Fresh-session probes ran from an empty directory with only installed-skill loading allowed:

- Claude Code 2.1.308 (non-persistent, `Skill` as the only tool): each of six prompts loaded the matching skill through the `Skill` tool (`why`, `how`, `blast-radius`, `verify-app`, `tidy`, `unslop`). Explicit `/step-back` returned its five H2 headings; the equivalent plain request returned `NO_STEP_BACK_CONTEXT` and made no `Skill` call.
- Codex CLI 0.163.0 (`codex exec --ephemeral`, read-only sandbox): the same six prompts each read the matching `skills/<name>/SKILL.md` and named it. Explicit `$step-back` returned its five H2 headings; the plain request returned `NO_STEP_BACK_CONTEXT`.

The first OpenCode attempt, on 1.19.3, failed every run, including a bare "Reply with exactly OK." control, because the local Magic Context plugin build supported database migrations through v91 while its database was at v95. After the plugin was repaired, OpenCode 2.1.2 passed the control and ran the probes in fresh `--standalone` sessions with `openai/gpt-5.6-terra`. A probe agent defined in the empty directory's `opencode.json` allowed only the `skill` tool, so the user's default agent's skill allowlist did not decide the result:

| Probe | Skill loaded | Reply | Verdict |
| --- | --- | --- | --- |
| `why` | `why` | `SKILL=why` | Verified |
| `how` | `how` | `SKILL=how` | Verified |
| `blast-radius` | `blast-radius` | `SKILL=blast-radius` | Verified |
| `verify-app` | `verify-app` | `SKILL=verify-app` | Verified |
| `unslop` | `unslop` | `SKILL=unslop` | Verified |
| `tidy` ("Simplify the code I just changed...") | `simplify` (two runs) | `SKILL=simplify` | Not matched |
| `tidy` ("Clean up the code I just changed...") | `tidy` | `SKILL=tidy` | Verified |
| `step-back` attached through the session API | content injected | five H2 headings | Verified |
| `step-back` plain request | none | `NO_STEP_BACK_CONTEXT` | Verified |

- The `simplify` result comes from `oh-my-opencode-slim`'s bundled `simplify` skill, the upstream source of `tidy`. When both are active in OpenCode, a prompt that says "simplify" selects it. `skills/tidy/SKILL.md` stays unchanged; the collision is resolved in the local OpenCode configuration.
- OpenCode 2.x lists `step-back` with `autoinvoke: false`. It does not register skills as slash commands, so `opencode run "/step-back ..."` sends plain text; the explicit probe attached `skills: [{"id": "step-back"}]` to `session.prompt` instead.

Not verified:

- Gemini: no compatible CLI is installed.

## Follow-ups

- Hand-synced reference copies (`domain-modeling` formats, `tdd` `tests.md` and `mocking.md`, `DEEPENING.md`, triage role definitions) are left for a later batch.

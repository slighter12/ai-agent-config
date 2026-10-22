# Matt Pocock Skills Upstream

## Pinned source

- Repository: `https://github.com/mattpocock/skills`
- Tag: `v1.3.1`
- Commit: `24fe0ef7737efae15c87225755e9f6f5965e4888`
- Package version: `1.3.1`
- License: `LICENSES/mattpocock-skills-MIT.txt`

The `v1.3.1` annotated tag object `0b6cee10f260a2e048279cf737bfd3e37b1fce0b` peels to the pinned commit `24fe0ef7737efae15c87225755e9f6f5965e4888` (`git ls-remote` `refs/tags/v1.3.1^{}`), and its `package.json` reports `1.3.1`. The previous fixed point is preserved in `docs/captures/2026-08-08-matt-skills-v1.2.3-sync.md`: the `v1.2.3` tag pointed to `6acc160e4e0cd062dbbbd7a1b26ae92855edf07e` and the earlier pin `84fdeffd12f2ee307994d1eb6feb48173b6e0502` was two commits later.

The local catalog vendors the promoted runtime artifacts under upstream `skills/engineering/` and `skills/productivity/`. Each upstream directory is flattened to `skills/<skill-name>/` so the multi-provider installer can link skills directly into provider discovery roots.

## Architectural influences

* [`cursor/plugins/pstack`](https://github.com/cursor/plugins/tree/main/pstack/skills) (`@poteto` / Lauren Tan; read at `d457567332aa6989e753ea289f30ea2bfbc80bc0`; license `LICENSES/pstack-MIT.txt`; later commits reviewed through `eb4db2b` on 2026-10-22 with nothing adopted, since none touched these sources): Agent-First Local Legibility, Boundary Discipline (Deterministic Core / Thin Effectful Shell), Task-Proportional Verification (`Prove It Works`), and 5 Layers of Correction.
* [`alvinunreal/oh-my-opencode-slim`](https://github.com/alvinunreal/oh-my-opencode-slim) (read at `e0a513368ff55b4ea7a25e148d75945c4d34a434`; license `LICENSES/oh-my-opencode-slim-MIT.txt`; later commits reviewed through `a598a26` on 2026-10-22 with nothing adopted, since `src/skills/` was unchanged): Lightweight role boundaries, Evidence Paths (`verification-planning`), and Surgical Changes.

Files adapted from these sources (rewritten, not vendored):

- From pstack: `skills/codebase-design/DESIGN-RED-FLAGS.md` (`architect/references/design-red-flags.md`), `skills/diagnosing-bugs/PERF.md` (`benchmark-checklist`, `poteto-mode/playbooks/perf-issue.md`), `skills/why/` (`why`, `why/references/epistemics.md`, `code-archaeology.md`, `incident-postmortem.md`), `skills/how/` (`how`, `how/references/explainer-prompt.md`, `explorer-prompt.md`), `skills/blast-radius/` (`blast-radius`), `skills/verify-app/` (`create-verification-skill`, `maintain-verification-skill`), `skills/unslop/` (`unslop`, `technical-writing`), the subtraction rules `skills/tidy/` applies (`principle-laziness-protocol`, `principle-subtract-before-you-add`), and `skills/step-back/` (`principle-attack-the-premise`, `principle-redesign-from-first-principles`).
- From oh-my-opencode-slim: `skills/tidy/` (`simplify`), `skills/step-back/` (`step-back`), and overlay rules in `to-spec`, `implement`, `codebase-design`, `retro`, and the agent definitions (`verification-planning`, `simplify`, `step-back`, `deepwork`, `reflect`, and the orchestrator prompt).

## Upstream v1.3.1 Changes Adopted

These changes come from the pinned `v1.3.1` commit above and were reapplied together with the local overlays and divergences listed below.

1. **`CONTEXT.md` -> `GLOSSARY.md` Rename**: All references across `domain-modeling` (including `GLOSSARY-FORMAT.md`), `grill-with-docs`, `codebase-design`, `tdd`, `pr`, `improve-codebase-architecture`, `triage`, `diagnosing-bugs`, `setup-matt-pocock-skills`, and `docs/agents/domain.md` now use `GLOSSARY.md` and `GLOSSARY-MAP.md`.
2. **Added `/pr` (`skills/pr/SKILL.md`)**: Standardized PR body structure with Visual Summary, Before/After Evidence, and Merge Danger (One-way vs. Two-way door + Blast Radius), enhanced to respect `.github/pull_request_template.md` when present.
3. **Added `/retro` (`skills/retro/SKILL.md`)**: Post-session environment retrospective across 7 categories, extended with **Project Law Graduation** to persist user-selected style/migration decisions into linters, `CODING_STANDARDS.md`, or `AGENTS.md` / `docs/adr/`.
4. **Retired the upstream merge-conflict skill**: Removed per upstream v1.3.0; its name, version, and reason live in `config/retired-skills.json` (retired names never appear on routing surfaces).
5. **Harness-Portable Skill Loading**: Wrapper skills (`grill-me`, `grill-with-docs`) load `grilling` on every harness, resolving upstream `#1123`. The upstream sibling-path fallback is replaced locally; see Intentional Local Divergence 6.

6. **Smaller v1.3.1 corrections**: `diagnosing-bugs` Phase 6 is cleanup-only (the post-mortem handoff moved to `/retro`), `domain-modeling` uses the v1.3.1 trigger description, and `to-spec`, `triage`, and `wayfinder` tell the user to run `/setup-matt-pocock-skills` instead of running it.

## Intentional Local Divergences (Do Not Overwrite on Upstream Sync)

1. **Single-Skill Self-Sufficiency & Embedded Pre-Flight Gate**:
   * Because coding harnesses typically auto-activate at most 1–2 skills per prompt, `to-spec`, `codebase-design`, `tdd`, and `implement` each embed a compact **Pre-Flight Context & Agency Gate** (4-tier rule precedence + 6-scenario detection + the canonical (a)–(d) Agency Gate Menu) rather than splitting principles across 5+ micro-skills.
2. **Merged the standalone ticket-splitting skill into `skills/to-spec/SKILL.md`** (retirement recorded in `config/retired-skills.json`):
   * Eliminates double-compression drift between spec synthesis and ticket slicing (upstream `#341`).
   * Routes single-slice changes immediately to `implement`. For multi-slice work, synthesizes Agreed Test Seams, Frozen Implementation Decisions, and the Tracer-Bullet Ticket Graph (`Blocked by` + `Write scope`) in a single pass.
3. **Merged the standalone spec-implementation skill into `skills/implement/SKILL.md`** (retirement recorded in `config/retired-skills.json`):
   * Dispatches automatically based on input shape: runs in-window for a single ticket/slice, or orchestrates isolated git worktree `builder` subagents across the ticket frontier when passed a multi-ticket spec.
   * Includes local frontier tracking, `Write scope` collision serialization, `.patch` pre-commit backup (upstream `#991`), and a strict 1-round post-review repair cap (upstream `#1014`).
4. **Upgraded `skills/grilling/SKILL.md`**:
   * **Entailment Pruning**: Prevents `A -> A' -> A''` paraphrase loops and implied-ancestor regressions once a downstream choice settles its prerequisites.
   * **Cross-Branch Tension (`A x B`)**: Explicitly probes independent branches whose combined answers create workflow, algorithmic, or constraint collisions.
   * **Empirical Unknown Escape Hatch**: Routes codebase facts to subagents and empirical/research unknowns to `prototype` spikes rather than re-asking the user.
   * **Descriptive Back-References**: Requires `Q3 (<short topic>)` when referencing prior questions in later rounds (upstream `#1115`).
5. **Language-Agnostic Mechanical `skills/tdd/SKILL.md` & `skills/codebase-design/SKILL.md`**:
   * Removes TypeScript/JS-specific coupling and adds concrete mechanical redlines (`Single-Caller Helper Ban`, `No Call-Only Mock Assertions`, `Evidence Conservation Law`, `Brownfield Sprout Pattern`, `Scrap Circuit Breaker`, and `Post-Green Same-Branch Test Pruning`).
6. **Call Skills by Name, Never by Relative Path**:
   * A skill that hands off to another skill says "use the `<name>` skill" when the target is model-invocable, or tells the user to run `/<name>` when it is user-invoked. Upstream `#1123`'s `../grilling/SKILL.md` fallback is removed, because an installed skill directory has no guaranteed sibling layout.
   * `writing-for-agents/LOCAL-PROVIDER-MECHANICS.md` records how each provider resolves a name, and a contract test rejects any parent-directory relative path in `skills/**/*.md`.
7. **Concrete Trigger Wording**:
   * `prototype` adds a `Use when` clause to its upstream description, and `diagnosing-bugs` writes its trigger phrases as a sentence instead of slash lists. Skill bodies gate reference files and skill calls on observable conditions rather than judgment words such as "if needed".
   * `research` hands this repository's design rationale to the `why` skill instead of carrying its method, `DESIGN-IT-TWICE.md` produces its designs as sub-agents or sequentially under the `codebase-design` dispatch rule, and `improve-codebase-architecture` keeps its candidate card spec only in `HTML-REPORT.md`. A contract test rejects judgment phrases in skill files.

## Post-pin cherry-picks

These upstream `main` commits landed after the pinned `v1.3.1` commit. They are folded in as local overlays; the pin above does not include them.

- `f3fc563` (`diagnosing-bugs`): before trusting a forced red, `diff` the mutated code or fixture against a pristine copy to prove the mutation landed.
- `8295b8e` (`wayfinder`): cross-references use real tracker references (never a placeholder `#<n>`), research branches never get a pull request, and a ticket is resolved by its tracker-defined type rather than its body. The local overlay already kept triage roles such as `ready-for-agent` off maps and decision tickets.
- `954ee7f` (`ask-matt`): read a skill's `SKILL.md` before stating what it does or recommending a step be skipped. Already present locally.
- `4f4e943`, `0f5e033` (`handoff`): resolve the temp directory from `$TMPDIR`, then `/tmp`, then `%TEMP%` on Windows. Already present locally.
- `3f59913` (`tdd`): give each proposed test seam one line on what it catches and what it misses.
- `ad5400b` (`teach`): workspace paths resolve from the directory `/teach` was run in, only the `*-FORMAT.md` links resolve from the skill folder, and quiz answers vary the position of the correct one.
- `5b7cade`, `e0efb6e`, `8e67db2`, `93497a8` (`setup-matt-pocock-skills`): create the configured triage labels a GitHub or GitLab tracker lacks, list external PRs through the REST pulls endpoint's `author_association`, and request `glab` JSON with `-O json`. Local overlay: label creation is a hosted write that uses the confirmed selector and needs the user's approval of the exact label list; the issue-read and sub-issue command changes were not taken because the local contracts already differ.
- `868c4cf`, `49dd158` (`wizard`): `ask` and `ask_secret` fail at EOF with no input, `_clear` falls back to ANSI when `tput clear` fails, `open_url` reports a missing opener outside its silenced block, and stages run inside `run_wizard` so bash parses the whole file before the first prompt. Not taken: the `write_env` rewrite that writes through a symlink, which conflicts with the `wizard/template.sh` symlink rejection under Required overlays, and the Readline prompts, quote decoding, and `RED` removal.

## Local extensions

- `mcp-builder-go`
- `how`
- `why`
- `blast-radius`
- `verify-app`
- `tidy`
- `unslop`
- `step-back`

These remain outside the upstream set and are marked `local extension` in the canonical catalog. `mcp-builder-go` and `step-back` are user-invoked (`step-back` is the user's explicit brake on a drifting task); the others are model-invoked so other skills can call them by name.
Each carries local `agents/openai.yaml` interface metadata so Codex presents a human-readable name and description.

## Required overlays

- Every skill carries `metadata.invocation`; user-invoked skills also carry Claude, Codex, and OpenCode implicit-invocation controls.
- `validate-skills` resolves repository directories through one filesystem-identity helper and requires the requested catalog to be this repository's `skills/` directory before loading the repository-owned retirement manifest and routing evidence. Active and retired names share one canonical syntax validator. Core skill validation, retirement manifest reads, and routing reads share one `internal/securepath` opener. On Unix targets it holds a directory descriptor and walks each component with `openat` and `O_NOFOLLOW`, so the kernel never resolves a complete caller-supplied path; Solaris fails closed because its `syscall` package exposes no `openat` bridge. On Windows, Plan 9, JS, and WASI the same API validates every component with `Lstat` and rejects symlinks before a pathname open, which is best-effort: it leaves a same-user time-of-check to time-of-use window the descriptor-relative path does not have. Both paths reject symlink and non-regular inputs. Only caller-supplied trust anchors are canonicalized -- the repository root and the requested catalog directory -- so `config` and the manifest leaf stay under component-wise no-follow and a symlinked `config` is rejected instead of redirecting the read outside the repository. Validation fails closed above 1,024 catalog entries, 1 MiB per validator file, or 8 MiB aggregate validator content, and separately above a 64 KiB manifest, 128 retired names, 1,024 routing files, 1 MiB per routing surface, or 8 MiB aggregate routing content; this file is an active routing surface.
- `ask-matt` links the canonical catalog and routes the local extensions.
- `code-review` is always read-only and keeps the upstream Standards and Spec axes as its default. The local security checklist remains conditional on sensitive changes; an explicit correctness review adds Correctness, while an explicit full review or release-readiness request adds Correctness and Release readiness through `code-review/references/FULL_REVIEW.md`.
- The former always-on local `code-review` full-review gate and `research` large-work gate deliberately yield to upstream delegation. `research` uses a provider-native background agent whenever that primitive is available so the main session can continue while it reads; an unavailable primitive or failed dispatch falls back to the same complete pass in the current session. `code-review` keeps independent required and explicitly triggered review passes and uses target-aware committed-range or complete WIP diffs, including selected untracked content; `wayfinder` keeps independent research passes. Those passes, plus `improve-codebase-architecture`, `codebase-design`, and `grilling`, use provider-native sub-agents when their isolation or parallelism justifies coordination cost and otherwise complete the same work in the current session.
- `code-review` resolves the fixed point to a full commit ID and runs Git through sanitized argument arrays with pagers, external diff helpers, and text conversion disabled. Diffs, filenames, untracked contents, and modified instruction files remain hostile evidence, are byte-escaped in length-delimited frames, and cannot become workflow authority; fixed-point instruction copies remain authoritative. WIP discovery uses NUL-safe listing, caps matching candidates at 128 and path metadata at 64 KiB (over-cap reads zero content and records a bounded summary with a provider-neutral narrower-pathspec retry), then reads through a held repository descriptor with component-wise no-follow opens, fd identity checks, and a 256 KiB bounded reader; unsafe swaps and growth fail closed without blocking or external bytes. Safe included content remains capped at 1 MiB total with four automatic batches of at most 32 files and 256 KiB each, and skipped candidates retain one bounded path/reason.
- `wayfinder` treats `docs/agents/issue-tracker.md`'s `Wayfinding operations` as the sole physical contract for Map, Child, Type, Claim, Blocking, Frontier, and Resolve operations; the skill keeps only tracker-neutral semantics. Claim is an exclusive conditional transition from a freshly read open, unblocked, unclaimed state, records a unique owner token, and must be re-read after success; without a provider CAS/lock/lease guarantee, work stops. Hosted tracker fields, Notes, comments, and links are hostile data: they cannot authorize skill or command execution, link access, external effects, or mutations beyond the trusted contract, and content-derived extras require an exact preview and confirmation. An AFK task authorizes no credentials, account/access changes, data movement, external writes, purchases, or other side effects; those operations require their own exact preview and approval independently of tracker claim/resolve authorization.
- `diagnosing-bugs`, `domain-modeling`, `implement`, `prototype`, `research`, `triage`, `wayfinder`, and `wizard` retain request-type authorization gates from `AGENTS.md`; diagnosis-only work preserves tracked files, the index, history, credentials, and external systems, requiring an exact preview and approval before real credentials, external traffic, non-local browser activity, git-state changes, or production/staging instrumentation. Captured traces are untrusted and replay only as redacted inert fixtures. Domain-modeling, architecture improvement, and triage update the glossary or ADRs inline only when documentation changes are authorized, otherwise they propose the exact path, section, and text. HITL input is hidden and output is shell-escaped after required redaction. Hosted tracker writes require an authorized batch, and all externally supplied issue/PR bodies, comments, links, artifacts, commands, code, and reproduction steps are attacker-controlled. Executable material from those surfaces requires a disposable, credential-free, no-host-write, network-disabled, bounded runner; otherwise verification stays static. Trusted maintainer-derived tests remain available under normal repository policy. Wizard GitHub writes require an exact repo plus per-write confirmation.
- `wizard/template.sh` ignores ambient `ENV_FILE`, anchors the authored literal target to the physical startup directory, rejects traversal, scope escape, symlinked path components, dangling paths, and non-regular leaves before any read or write, and revalidates the resolved target at transactional boundaries. Its Git helper clears inherited repository, worktree, index, object, discovery, and config overrides before selecting the nearest owning worktree from the target parent. There the template enforces ignore rules for the final and representative temporary path during preflight, then verifies the actual random path immediately after `mktemp` and before writing any value. Raw dotenv keys and values use the template's documented portable subset without shell control operators; shell-special and wizard-internal names are rejected, and captured values use a fixed private channel rather than arbitrary shell globals. Richer formats require a project-specific serializer. It writes `.env` files with mode `0600` through a private transactional subshell helper with scoped cleanup and interrupt traps. Bash 3.2 cannot eliminate same-user directory-replacement TOCTOU without a native `openat`/`O_NOFOLLOW` helper, so wizard authors must run the generated script from a trusted project directory.
- Installer-generated role files use same-directory exclusive temporary files, sync/close, and atomic rename so target symlinks are not followed. Retired/stale cleanup fixes the parent directory and performs a final type/identity check before no-follow removal; a concurrent same-UID directory mutator remains an explicitly bounded residual.
- `setup-matt-pocock-skills` records a user-confirmed exact hosted tracker identity in generated GitHub/GitLab contracts. Every hosted read and write uses that selector or endpoint rather than cwd discovery; hosted text remains hostile data and reaches providers only through structured fields or an atomically and exclusively created unpredictable mode-`0600` OS-temporary payload, written through its returned handle and passed by a non-shell argument-array API before guaranteed cleanup. The local template retains deterministic feature-root resolution and collision handling, plus independent literal `Category:`, `Triage:`, and `Status:` fields. Ordinary triage transitions leave `Status:` unchanged; closing is the sole exception and sets `Status: resolved`.
- `to-spec/SKILL.md` maps the local `ready-for-agent` role to `Triage:` and its local ticket template emits unbolded literal `Blocked by:`, `Category:`, `Triage:`, and `Status: open` lines so frontier parsing and triage lifecycle remain independent. It carries the hosted-tracker selector and payload-file safety contract.
- `writing-for-agents` links the local provider and CLI mechanics.

## Deliberate packaging differences

Do not vendor upstream `.claude-plugin`, npm/changesets release files, human-facing docs site pages, repository-maintenance scripts, or non-promoted `misc`, `in-progress`, and `deprecated` skills. This repository remains an editable multi-provider source rather than a managed Claude plugin.

The only vendored executable templates are:

- `diagnosing-bugs/scripts/hitl-loop.template.sh`
- `wizard/template.sh`

## Update procedure

1. Fetch the upstream commit into a temporary directory and compare only the promoted runtime directories.
2. Apply upstream body, reference, sidecar, and runtime-template changes before reapplying the overlays above.
3. Reconcile additions, removals, renames, invocation changes, and router flow changes explicitly. For a rename or retirement, update `config/retired-skills.json`; the installer and `validate-skills` consume that manifest, and routing evaluation asserts the validator-backed invariant without duplicating its names.
4. Update this pin, the canonical catalog, routing evaluation, sync capture, and license when needed.
5. Run the full Go suite, `scripts/check_go_targets.sh`, catalog validation, Bash syntax checks, and `git diff --check`. The cross-target build is required whenever `hooks-go` build-tagged code changes: it is the only check that links the per-`GOOS` `openat` bridges in `internal/securepath`.

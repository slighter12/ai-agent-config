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

**Block-Level Verbatim Stitching**: a port or re-sync keeps each source block's exact wording and concrete examples inside its section. Edits are limited to the seams between sources and to the overlays this file lists (authorization and read-only gates, Single-Window dispatch, the `Verified` / `Not verified` / `Inconclusive` verdicts, `use the <name> skill` pointers, concrete triggers, and the local label format). Every other change is a divergence and is listed below its file. Paraphrasing a whole file in new words is not a port. Every file listed under "Files adapted from these sources" was first ported as a rewrite and is now re-synced under this rule.

`skills/unslop/` divergences from pstack `unslop` and `technical-writing` at `d457567`:

- Rules 13 (em dash), 14 (colon), and 17 (title case) and the `technical-writing` punctuation rules (semicolons, em dashes, slashes) are not adopted; the Boundaries section matches the document's own punctuation and heading case instead.
- Rule 7 adds `seamless`, carried over from the redline `pr` used to hold.
- Rule 16 judges a lead-in by whether it restates its line, not by its punctuation, so a colon label with new detail is fine.
- Rule 26 exempts terms defined in `GLOSSARY.md` or used as identifiers.
- Rule 28 appends the `technical-writing` "Vary the rhythm" sentences on one thought versus one length.
- Rule 33 keeps the document's established notation.
- The pronoun rule replaces "when in doubt" with the observable condition "when a pronoun could point to two nouns" (concrete-trigger overlay).
- Labels use `**Label**: text`; the frontmatter, Boundaries, Process, and the rule-number note are local, and `## Sentences` adopts seven `technical-writing` rules with their upstream wording.

`skills/why/` divergences from pstack `why`, `why/references/epistemics.md`, `why/references/sources/code-archaeology.md`, and `incident-postmortem.md` at `d457567`:

- Single-Window dispatch: the session searches each evidence category itself and dispatches a read-only subagent only when raw output would bloat the window. Dropped with it: the model-slug and `pstack-models.mdc` paragraph, `readonly: false`, the investigator and synthesizer spawns and their prompt templates, Step 5 Present, and the Cursor `mcps/` discovery (replaced by "list the read-only connectors already available in this session"). Seams follow from it: "Before spawning investigators" becomes "Before searching any evidence category", "This anchor goes to every investigator and to the synthesizer" becomes "This anchor seeds every evidence category below", "The synthesizer must follow it" becomes "The synthesis must follow it", "Always spawn" becomes "Always search it", "When to skip an investigator" becomes "When to skip a category", "answer inline" becomes "answer from it alone", "one line per investigator" becomes "one line per evidence category", and epistemics' "Before the synthesizer delivers the output, it checks" becomes "Before you deliver the output, check".
- `code-archaeology.md` and `incident-postmortem.md` are inlined as the source control and incident playbooks of Step 3 instead of being loaded per investigator.
- Read-only and trust overlays: the read-only line, the Trust Boundary section, `-R <host/owner/repository>` and "only when the tracker contract names that target" on `gh pr view`, and `--no-pager`, `--no-ext-diff`, and `--no-textconv` on git queries.
- Concrete trigger: "Adapt as needed" becomes "Drop What We Can Reasonably Infer or Competing Hypotheses only when it is empty".
- Local additions: the `research` distinction in the intro, the embedded-hypothesis line in Step 1, "Also read `GLOSSARY.md`" in the source control playbook, the described section list in Step 4 (from the structure upstream keeps in `synthesizer-prompt.md`), the Preserve / Change / Avoid / Risk bullets, the `research` copy note in `EPISTEMICS.md`, and its closing check on contradictions and embedded hypotheses.
- Format: epistemics tiers are a `**Label**: text` list instead of `###` headings, and `Before Finalizing` is the epistemics heading the skill points to.

`skills/tidy/` divergences from oh-my-opencode-slim `simplify` at `e0a5133`, and `skills/codebase-design/SKILL.md` §3 divergences from pstack `principle-subtract-before-you-add` and `principle-laziness-protocol` at `d457567`:

- `simplify`'s "Guidance for This Repository" (TypeScript, that repository's hooks) and "Final-state verification" (its release instructions) are not adopted; Step 2 Establish the Baseline and the Report replace them.
- "Proportionate final-state verification" becomes the local baseline checks: principle 1 asks "What baseline checks will reveal a behavior change?", Step 4 item 2 becomes "Rerun the baseline tests for the file or package the step touched to check preservation", Step 5 reruns the whole baseline once after the last step, and Step 4 adds the revert-and-do-not-stack line.
- Local overlays: the Boundaries section (change-request gate, scope, project-state gates, Tier 2 structural contracts, no dependencies, uncommitted), Step 2 with its characterization-test question, the `why` and `codebase-design` pointers, the comment-keeping line in Step 1, and the `Verified` / `Not verified` / `Inconclusive` Report. Boundaries drops its performance line and its drive-by line because When NOT to use and principle 5 carry them.
- `codebase-design` §3 places both principles under Subtract Before You Add and Laziness Protocol, and keeps two local lines below the upstream bullets: deleting dead paths and single-caller wrappers before adding a parameter, flag, or wrapper, and "Subtraction stops at meaning".

`skills/step-back/` divergences from oh-my-opencode-slim `step-back` at `e0a5133` and pstack `principle-attack-the-premise` and `principle-redesign-from-first-principles` at `d457567`:

- The slim prose is split around five local steps: its first paragraph is the intro, and its second frames step 5 with the four-line recommendation format.
- Read-only overlays: "Then let them decide" becomes the local "Then stop and wait for the user's decision", "If the answer is yes, carry on with confidence" becomes "If the answer is yes, say so plainly, so the user can carry on with confidence", and step 3 runs the premise check only when a local read-only check can and reports it as `Verified`, `Not verified`, or `Inconclusive`.
- Step 3's trigger adds "or when workarounds keep piling up", since a user-invoked step back can face a pile of workarounds rather than failed fixes.
- Links to unported principles are dropped: "per [Build the Lever]" and "per [Fix Root Causes]" end at "Make the check rerunnable" and "ask the next "why"", and "per the [Laziness Protocol]" points to `codebase-design` §3. "This principle is distinct from [Redesign from First Principles]" becomes "This step is distinct from step 4, Redesign From First Principles".
- Local additions: steps 1 and 2, the `diagnosing-bugs` pointer, the comparison of finishing versus switching, and the `codebase-design` pointer in step 4. Labels use `**Label**: text`.

`skills/verify-app/` divergences from pstack `create-verification-skill` and `maintain-verification-skill` at `d457567`:

- The two skills are one skill with a create mode and a maintain mode. Seams follow from it: "Point the user at `/maintain-verification-skill`" becomes "Point the user at maintain mode", "a skill generated by `/create-verification-skill`" becomes "a skill generated in create mode", and "point at `/create-verification-skill`" becomes "switch to create mode".
- Paths: `.cursor/skills/verify-<app>/` becomes `verify-<app>/` under the provider skill root the Boundaries section names. The link to `references/feature-map-example/` is dropped, since that example is not ported; the feature-map shape stays in the text.
- Authorization overlays: the Boundaries section (edit gate, provider root, disposable local instances, credential names only, PID-tracked stops, uncommitted changes, mode choice). "Fix that first (or report it precisely)" becomes "report it precisely, with the exact command and failure", because product code stays unedited. `changed` leaves proven corrections as uncommitted local edits instead of one PR; "before it ships" becomes "before it counts as changed", "keep it out of this PR" becomes "keep it out of the corrections", and Ship or stop reports instead of opening a PR.
- Single-Window dispatch: "Source wave. One read-only subagent per feature file, launched concurrently" becomes "Source pass", explored in this session with read-only `explorer` subagents only when raw output would bloat the window; "The coordinator owns all driving" becomes "This session owns all driving".
- Test-locked wording: "Never weaken the journey" becomes "Never weaken a journey", and "stay unverified" becomes "stay Not verified".
- Local additions: the `tdd` and `retro` distinction in the intro, "following the project's existing skill conventions", the model-invocable default, and "Redact secrets" in Evidence. Labels use `**Label**: text`.

`skills/blast-radius/` divergences from pstack `blast-radius` at `d457567`:

- Read-only overlay: the Boundaries section; step 5 writes the proof in the OS temporary directory and stops at the highest level reached when the boundaries forbid running it.
- Single-Window dispatch: step 6 (`arena`, several models on one question) is not adopted.
- "How sure are you" is the `## Proof Ladder` that `pr` grades against, and each level carries a label (Stated, Pointed, Walked, Ran, Reproduced); "Step 4 is usually" becomes "Level 4 is usually" and "say which step you got it to" becomes "say which level", so ladder levels and steps stay distinct. "If you couldn't prove it, write unproven" becomes "If you couldn't get it to level 4, write **unproven** and name the cheapest check that would prove it", matching the level `pr` calls proven.
- Pointers: "Use `why` step 2" becomes "Use the `why` skill and follow its Step 2", and "Write it through `unslop`" becomes "Use the `unslop` skill on the writeup".
- Local additions: the `code-review` and plan sentences in the intro, "Also check retries, concurrent callers, and configuration keys" in step 3, and "secrets" in the stripping line. Labels use `**Label**: text`.

`skills/how/` divergences from pstack `how`, `how/references/explainer-prompt.md`, and `explorer-prompt.md` at `d457567`:

- Single-Window dispatch: the model-slug paragraph, the Task subagent spawns, the explainer and explorer template preambles, and step 4 Present are not adopted. The session explores and explains itself; read-only `explorer` subagents run only when raw search output would bloat the window. Seams follow from it: Simple "Spawn one explainer..." becomes "explore and explain in one pass in this session", "Spawn all explorers in a single message" becomes the dispatch condition, "Once all explorers have returned, spawn one explainer" becomes "Once every angle has returned, reconcile the findings", "Other explorers are investigating different slices of the same subsystem in parallel" is dropped from the explorer contract, and "The explorers did the work" becomes "When explorers did the work".
- Concrete trigger: "When in doubt, take the simple path" becomes "Treat the question as Complex when the answer spans two or more top-level directories or services; otherwise take the simple path".
- Vocabulary: the explorer's `Boundaries` section is `Connections`, because `codebase-design` rejects "boundary".
- Local additions: the read-only line, the skill distinctions in the intro, the Trace Rules section, the example angles in step 2, the Open Questions section, and the closing `why` and `codebase-design` pointers. The explorer contract uses `**Label**: text`.

`skills/diagnosing-bugs/PERF.md` divergences from pstack `benchmark-checklist` and `poteto-mode/playbooks/perf-issue.md` at `d457567`:

- The perf-issue playbook becomes the Fix Order section. Dropped with the playbook system: the control-skill baseline pointer, the implementation subagent and its model slug, the `sequence-verifiable-units` link (kept as "Verify each attempt before trying the next"), "Run **Opening a PR**", the Hillclimb pointer, and the checklist's "Other perf material" section. "Vet ... with the **benchmark-checklist** skill" becomes "with the questions above"; "`how` to ground hypotheses" becomes "Use the `how` skill to ground hypotheses"; "`architect` first" becomes "use the `codebase-design` skill first"; "Cite the measurement in the PR" adds "when one is opened"; "per the **Opening a PR** playbook" is dropped from the PR-body line; the playbook's Reply becomes the last Report line.
- Local additions: "Read it too when a bug is "slow"" in the intro, "(cache or memoize)" on the second mantra, and "Inconclusive is not a pass" in Report.

`skills/codebase-design/DESIGN-RED-FLAGS.md` divergences from pstack `architect/references/design-red-flags.md` at `d457567`:

- "Screen every candidate before synthesis" becomes "Screen every candidate before comparing or adopting it (see [DESIGN-IT-TWICE.md]), and when reviewing an existing module", and the intro points to `SKILL.md` for vocabulary.
- Vocabulary: "across several boundaries" becomes "across several stages" and "a forwarding boundary" becomes "a forwarding layer", because `codebase-design` rejects "boundary"; Pass-Through Method adds the Deletion Test pointer.
- Headings use title case, which `code-review` and a contract test match.
- Local additions: the Agency Gate sentence on recording a deferred migration in an ADR, and the `internal/`, unexported-name, and `exports`-map examples in Importable Internals.

`skills/tdd/SKILL.md` keeps only the Attack the Premise trigger and points to the `diagnosing-bugs` skill for the method, so the method is not maintained in `tdd` as well.

`skills/diagnosing-bugs/SKILL.md` Phase 3 carries pstack `principle-attack-the-premise` at `d457567` under its **Attack the premise.** lead-in, with the same divergences as `skills/step-back/` §3 except the step-back-only ones: the links to unported principles are dropped, "per the [Laziness Protocol]" points to `codebase-design` §3, the closing sentence that distinguishes it from Redesign from First Principles is dropped because this skill has no redesign step, and labels use `**Label**: text`. The trigger stays failed fixes only, without step-back's "or when workarounds keep piling up", because a diagnosis pile is fixes for one failing behavior.

The same rule covers these Matt skills, re-synced against the pinned `v1.3.1` commit. Overlays listed under "Intentional Local Divergences" and "Required overlays" below are not repeated here.

`skills/pr/` divergences from `pr` and its `CREDITS.md`:

- The frontmatter keeps the local description and `metadata.invocation`; the upstream `credits` block moves into the body sentence that links `CREDITS.md`, which is copied verbatim.
- Local additions: the 40-line budget, the repository-template pre-check (every standard GitHub and GitLab location), the `**Safety Fact:**` template line and its `blast-radius` grading, the Evidence S-tier and A-tier lines and the `Measured` / `Inferred` labels, the performance-number line, the Door Type and Blast Radius bullets, and `## Plain Speech` with its `unslop` pointer.
- "Start directly at `## Summary`" is appended to the preamble rule.

`skills/ask-matt/PHASE-BOUNDARIES.md` divergences from `ask-matt/PHASE-BOUNDARIES.md`:

- The heading reads "Phase Boundaries & Context Management".
- The smart zone "(~150k tokens)" becomes "(roughly the first half of the window, as `AGENTS.md` §3 defines it)", so the window size is not hard-coded.
- Question 4 adds one sentence naming this repository's subagent roles.

`skills/ask-matt/SKILL.md` divergences from `ask-matt`:

- `## Output Contract`, `## Routing Questions`, `## The Skills Along the Route`, and `## On-Ramps & Codebase Health` replace the upstream flow prose (main flow, on-ramps, codebase health, vocabulary layer). The routing questions name the same skills and branches; they are a deliberate restructure, not a cut.
- Names follow the local convention: `/name` for user-invoked skills and `` `name` `` for model-invoked ones.
- Context hygiene says "`/to-spec` has written the tickets" because ticket splitting merged into `to-spec` (Intentional Local Divergence 2).
- Standalone notes that `research` writes into the repository only when the request authorizes it.

`skills/retro/SKILL.md` divergences from `retro`:

- Local categories added after the upstream seven: Project Law graduation, domain glossary, and verification recipe.
- Local additions: the trust-boundary rule for session logs, workspace-only log reads, discarding best practices that trace to no session moment, ranking, the Hard Stop before any write, and `## Make the Mistake Impossible`.
- The verification-recipe category points to the `verify-app` skill's `### 2. Generate the skill` step.

`skills/code-review/SKILL.md` divergences from `code-review`:

- Labels use `**Label**: text`.
- The Middle Man smell keeps "directly, unless the layer adds policy", which a contract test locks.
- Local smells are grouped under **Local baseline additions (judgement calls)** after the upstream Fowler baseline.
- The upstream "parallel sub-agents" dispatch becomes Single-Window dispatch: sub-agents only when isolation or parallelism earns the cost, otherwise sequential passes, with the smell baseline pasted only into a sub-agent brief.
- The missing-tracker case continues with local spec discovery instead of telling the user to run setup, since review is read-only.
- Git commands use the sanitized forms under "Required overlays"; the commit list is captured only when a fixed point exists.

`skills/to-spec/SKILL.md` divergences from `to-spec` and the merged ticket-splitting skill:

- One interactive checkpoint (Step 3) carries the ticket quiz, the Agency Gate menu, seams, and the publish plan; "Iterate until the user approves" becomes "Revise in response to the answers, but do not open new question topics."
- The upstream tracker-routing paragraph is replaced by Local Markdown and Hosted modes, with the payload-file security rules and the parent spec kept free of `ready-for-agent`.
- The local ticket template adds `Write scope`, `Evidence path`, `Fast check`, `Final gate`, and `Category:` / `Triage:` / `Status: open` lines (see "Required overlays"); the hosted template adds the same `Fast check` and `Final gate` lines, and the spec template's Testing Decisions adds the spec's `Final gate`.
- The Fast Exit summary carries `Fast check` / `Final gate` and asks for approval of an evidence-only dependency or diagnostic surface, as Step 3 does. One sentence before the templates extends the spec template's file-path and code-snippet rule to every template instead of restating it.
- Prefactoring adds baseline capture and the riskiest unknown to the earliest tickets; expand–contract adds the codemod proof and the shared integration branch fallback.

`skills/implement/SKILL.md` adds one Cleanup line: when the tracker closes work through a PR and the user authorizes opening one, use the `pr` skill for its body. `implement` never opens a PR on its own.

Local protocol overlays (state machine and cross-skill contracts), none of which replace upstream prose:

- `implement`: resume reconciliation of the run trail against the integration branch before choosing a mode; only fulfilled blockers clear `Blocked by`; wayfinder decision tickets are not children; the **Stuck graph** stop; the Sequential Fallback reviews once after the last child; the `Replanned` pivot result; deferred-green migration batches with `Inconclusive (deferred to final integration ticket #<id>)`; gitignored fixtures provisioned in the worktree or the ticket is `Blocked`; the assigned `.patch` backup lives in the trail directory outside the worktree; builder rebases are part of the approved run plan and forwarded in the dispatch through a verbatim pre-approval line; `code-review` receives `BASE_SHA` plus a target ref for committed changes or WIP mode, with the files modified or created for the task as the pathspec, for uncommitted ones; `code-review` runs in both modes regardless of diff size or Agency Gate option, Step 8 confirms the review in a completion check separate from `decisions.tsv`, which has no `code-review` phase, and the Final Report copies the `review_line=` the run trail prints, both only when the run completes implementation; a builder pause answer is a `gate` row with `Approved` or `Declined`, and `gate` rows record user decisions only; pre-review fixes are `integration-fix`; style candidates are cited as `file:line`; **When the base moves** has `run-trail.sh base` derive the review base from git on every call, the single merge-base with the parent or its upstream when it descends from the captured base, else the captured base, after a local fast-forward of a remote parent in Mode B; iteration runs the ticket's `Fast check`, **Re-test after a rebase or tip move** scopes re-tests by `git diff --name-only <old-tip> <new-tip>` and treats patch-level version bumps as normal drift, and the `Final gate` runs once per run; Cleanup removes only worktrees and deletes only branches the approved plan or a later user message names, a resumed run asks again, and the run trail's `worktrees` report confirms the result.
- `implement/scripts/run-trail.sh` (local, POSIX `sh`, tested by `hooks-go/integration/run_trail_runtime_test.go`) is a convenience helper that owns the mechanics the prose used to spell out, not a guard against a misbehaving caller: the trail lives in `${TMPDIR:-/tmp}/run-trail/<repo-key>/<slug>/`, never in the project's working tree; `init` captures `BASE_SHA` and reads the parent with `git branch --show-current` (never guessing one) into the run-plan `gate` row, reuses a trail only while nothing happened since it ran, and with `--fresh` archives an old trail instead of deleting it; `log` stamps `ts` from the UTC clock, unique per row, validates the phase and result values, resolves one full commit (never a range), builds the deferred and `Supersedes` results, and flattens tabs and newlines; `status` prints the recorded state and the live `ticket=` rows; `base` derives the review base from git on every call (the single merge-base with the parent or, in Mode A, its upstream when it descends from the capture, else the capture) and prints the Final Report `review_line=`; `scope` checks staged or branch files against the `Write scope`, where `*` stays inside one path segment and `**` or a trailing `/` crosses directories; `worktrees` reports leftover and dirty worktrees. Exit `3` hands the decision to the user and prints `needs_user=` first; the script never commits, merges, rebases, or deletes.
- `builder` (all three definitions): a backup patch outside the worktree, rebase only on a forwarded approval, and a re-test scoped by what the rebase brought in, never the full suite.
- `test-runner` (all three definitions): the full suite runs only when the dispatching session asks for the final gate.
- `to-spec`: deferred-green batches are marked in their Evidence path, and local ticket numbers continue after the highest existing number.
- `code-review` and `references/WIP_INPUTS.md`: an optional caller-named target resolved with the same `^{commit}` procedure.
- `wayfinder` and the tracker contracts: resume an own claim, user-authorized stale-claim release or takeover through a CAS on the stale token, idempotent Resolve that appends a missing map pointer once per ticket, and non-fulfilling closures (`Status: closed (...)`, GitHub `not_planned` / `duplicate`, GitLab `wontfix` label) that never unblock a dependent. `triage` closes `wontfix` that way, except **Already implemented**, which closes as fulfilled.
- `handoff`: the **Run state** item carries the integration branch and its parent branch (when used), and the run trail's `dir=` and current `base=`.

`skills/codebase-design/` divergences from `codebase-design`, `DEEPENING.md`, and `DESIGN-IT-TWICE.md`:

- `SKILL.md` §1 (Agency Gate), §3 (Agent-First principles), §4 (Sprout), and §5 (Scrap Circuit Breaker) are local (Intentional Local Divergences 1 and 5). §2 carries the upstream glossary, Deep vs Shallow diagrams and questions, principles, and relationships verbatim; the deletion test moves to §3 with the other subtraction rules, and the two-adapter principle adds that a production I/O adapter plus a deterministic test fake counts as two.
- "Designing for testability" is §3.4 prose without the TypeScript blocks, per Intentional Local Divergence 5. Rejected Framings say "a language `interface` keyword" and add **edge**; Going Deeper gates each reference on an observable condition.
- `DEEPENING.md` keeps the one-adapter principle in `SKILL.md` §2 and points to it, so the rule has one source. Local additions: "never mock it or place it behind an artificial interface" (category 1), "plus clock and randomness" (category 4), and the never-wrap bullet under the testing strategy. "Network boundary" becomes "network edge", because `codebase-design` rejects "boundary".
- `DESIGN-IT-TWICE.md` dispatches designs as sub-agents or sequentially (Single-Window overlay), triggers on a Scrap Signal or an exported interface as well as the user's request, puts caller usage and the data shape before the interface in each design's output, phrases trade-offs as "we accept X in exchange for Y", screens designs against `DESIGN-RED-FLAGS.md`, and adds the caller-leverage, seam-simplicity, and local-legibility checks, the losing-alternative line, and open risks as questions. "Agent 4 (if applicable)" becomes "Design 4 (when a dependency crosses the seam)".

`skills/tdd/` divergences from `tdd`, `tests.md`, and `mocking.md`:

- `SKILL.md` is the local four-step workflow (Intentional Local Divergence 5). It carries the upstream "What a good test is" sentences, the seam definition, the agreed-seams rationale, the "Ask:" line, the horizontal-slicing rationale, and the tautological-assertion examples inside those steps. The refactoring rule adds the `tidy` pointer.
- `tests.md` adds the Two-Edit Test and "The contract decides, not the assertion shape", rewrites the good/bad characteristic lists as "What Makes a Durable Behavioral Test" and one red-flags sentence, and points its anti-patterns at the `SKILL.md` redlines instead of restating them.
- `mocking.md` says "fake" for a test double, lists I/O and non-deterministic edges in place of "system boundaries", and points the >2-mocks rule at `SKILL.md` Step 2.

`skills/grilling/SKILL.md` carries the upstream rounds, frontier, reshape, facts, and done-rule text inside the local structure of Intentional Local Divergence 4. Facts may be found in this session or by a read-only sub-agent (Single-Window overlay). The round format adds a `**Recommendation**` label and "Reply `yes` to accept"; the completion gate adds the Decision Summary.

`skills/handoff/SKILL.md` keeps the upstream body verbatim, resolves the temp directory per the `4f4e943` cherry-pick, names the file and prints its path, and adds a `## Local additions` list (current stage, settled decisions, verified vs. unverified facts).

`skills/wait-what/SKILL.md` keeps the upstream instruction verbatim under a `# Wait, What?` heading.

`skills/teach/` divergences from `teach`: workspace paths and quiz rules follow the `ad5400b` cherry-pick, the file list adds a `GLOSSARY.md` entry that points to `GLOSSARY-FORMAT.md`, the learning-records entry points to `LEARNING-RECORD-FORMAT.md` for numbering instead of restating it, and the resources, mission, open-command, and quiz-length rules use observable conditions ("no Knowledge entry covering the skill", "missing or does not yet state why", the per-OS open command, "within 2 characters of the same length") instead of "well-populated", "unclear", and "if possible". `LEARNING-RECORD-FORMAT.md` criterion 1 reads "answered a question or completed an exercise correctly without hints" in place of "demonstrated genuine understanding".

`skills/prototype/SKILL.md` adds the `grilling` Spike handoff, asks once when the prompt names neither branch, gives the non-interactive default, moves the routing-convention rule to `UI.md`, and gates the capture branch and commit on the user's authorization.

`skills/wayfinder/SKILL.md` and `skills/improve-codebase-architecture/` keep the upstream text except the overlays above and these: ICA counts hot spots over the last 50 commits (widening to 200), adds the Deterministic Core, Style Drift, and Scrap Signal explore prompts, the Single-Window scan sentence, and the `Brownfield Sprout Seam` approach tag, keeps the candidate card fields only in `HTML-REPORT.md`, requires a cited failing test, incident, or call-site count to reopen an ADR, and writes domain docs only when authorized. The GitLab tracker contract keeps the upstream epic note under the rule that hosted text cannot select a parent relationship.

Category C, local additions kept on purpose:

- `diagnosing-bugs`: **Check the fence before you move it** and **Ask whether the triggering state is legitimate** (Phase 3), and the skip rule "only when the user says to, or when what it produces already exists" in place of "explicitly justified".
- `writing-for-agents`: the `## Mechanical heuristics` section.
- `research`: the Direct / Supported / Inferred confidence tiers.
- `domain-modeling/ADR-FORMAT.md`: the **Brownfield / Migration Strategy** optional section.
- `triage/AGENT-BRIEF.md`: the `Evidence path`, `Fast check`, `Final gate`, and `Write scope` fields and their examples.
- `code-review`: the Act On / Consider / Dismiss triage with its Nitpick gravity and Premature abstraction filters, the four local smells, the architecture and test smells, and the bounded repair rule.
- Punctuation: upstream uses no em dashes. Upstream sentences keep upstream punctuation; text that exists only locally (overlays, the Agency Gate menu that matches `AGENTS.md` word for word) may keep its em dashes.

Files adapted from these sources (stitched under the rule above, not vendored):

- From pstack: `skills/codebase-design/DESIGN-RED-FLAGS.md` (`architect/references/design-red-flags.md`), `skills/diagnosing-bugs/PERF.md` (`benchmark-checklist`, `poteto-mode/playbooks/perf-issue.md`), `skills/why/` (`why`, `why/references/epistemics.md`, `code-archaeology.md`, `incident-postmortem.md`), `skills/how/` (`how`, `how/references/explainer-prompt.md`, `explorer-prompt.md`), `skills/blast-radius/` (`blast-radius`), `skills/verify-app/` (`create-verification-skill`, `maintain-verification-skill`), `skills/unslop/` (`unslop`, `technical-writing`), the subtraction rules `skills/tidy/` applies (`principle-laziness-protocol`, `principle-subtract-before-you-add`), and `skills/step-back/` (`principle-attack-the-premise`, `principle-redesign-from-first-principles`).
- From oh-my-opencode-slim: `skills/tidy/` (`simplify`), `skills/step-back/` (`step-back`), and overlay rules in `to-spec`, `implement`, `codebase-design`, `retro`, and the agent definitions (`verification-planning`, `simplify`, `step-back`, `deepwork`, `reflect`, and the orchestrator prompt).

## Upstream v1.3.1 Changes Adopted

These changes come from the pinned `v1.3.1` commit above and were reapplied together with the local overlays and divergences listed below.

1. **`CONTEXT.md` -> `GLOSSARY.md` Rename**: All references across `domain-modeling` (including `GLOSSARY-FORMAT.md`), `grill-with-docs`, `codebase-design`, `tdd`, `pr`, `improve-codebase-architecture`, `triage`, `diagnosing-bugs`, `setup-matt-pocock-skills`, and `docs/agents/domain.md` now use `GLOSSARY.md` and `GLOSSARY-MAP.md`.
2. **Added `/pr` (`skills/pr/SKILL.md`)**: Standardized PR body structure with Visual Summary, Before/After Evidence, and Merge Danger (One-way vs. Two-way door + Blast Radius), enhanced to respect a repository PR or merge request template in any standard GitHub or GitLab location when present.
3. **Added `/retro` (`skills/retro/SKILL.md`)**: Post-session environment retrospective across 7 categories, extended with **Project Law Graduation** to persist user-selected style/migration decisions into linters, `CODING_STANDARDS.md`, or `AGENTS.md` / `docs/adr/`.
4. **Retired the upstream merge-conflict skill**: Removed per upstream v1.3.0; its name, version, and reason live in `config/retired-skills.json` (retired names never appear on routing surfaces).
5. **Harness-Portable Skill Loading**: Wrapper skills (`grill-me`, `grill-with-docs`) load `grilling` on every harness, resolving upstream `#1123`. The upstream sibling-path fallback is replaced locally; see Intentional Local Divergence 6.

6. **Smaller v1.3.1 corrections**: `diagnosing-bugs` Phase 6 is cleanup-only (the post-mortem handoff moved to `/retro`), `domain-modeling` uses the v1.3.1 trigger description, and `to-spec`, `triage`, and `wayfinder` tell the user to run `/setup-matt-pocock-skills` instead of running it.

## Intentional Local Divergences (Do Not Overwrite on Upstream Sync)

1. **Single-Skill Self-Sufficiency & Embedded Pre-Flight Gate**:
   * Because coding harnesses typically auto-activate at most 1–2 skills per prompt, `to-spec`, `codebase-design`, `tdd`, and `implement` each embed a compact **Pre-Flight Context & Agency Gate** (4-tier rule precedence + 6-scenario detection + the canonical (a)–(d) Agency Gate Menu) rather than splitting principles across 5+ micro-skills.
2. **Merged the standalone ticket-splitting skill into `skills/to-spec/SKILL.md`** (retirement recorded in `config/retired-skills.json`):
   * Eliminates double-compression drift between spec synthesis and ticket slicing (upstream `#341`).
   * Single-slice changes get a Fast Exit summary instead of a ticket graph. For multi-slice work, synthesizes Agreed Test Seams, Frozen Implementation Decisions, and the Tracer-Bullet Ticket Graph (`Blocked by` + `Write scope`) in a single pass.
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
   * A skill that needs another skill says "use the `<name>` skill" when the target is model-invocable. As upstream, only `ask-matt` recommends routes: no other skill tells the user which command to run next. Another skill writes a `/command` only to name a missing prerequisite (`/setup-matt-pocock-skills`) or to reference another skill's vocabulary; when it cannot proceed, it stops and reports why. Upstream `#1123`'s `../grilling/SKILL.md` fallback is removed, because an installed skill directory has no guaranteed sibling layout.
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
- `setup-matt-pocock-skills` records a user-confirmed exact hosted tracker identity in generated GitHub/GitLab contracts. Every hosted read and write uses that selector or endpoint rather than cwd discovery; hosted text remains hostile data and reaches providers only through structured fields or an atomically and exclusively created unpredictable mode-`0600` OS-temporary payload, written through its returned handle and passed by a non-shell argument-array API before guaranteed cleanup. The local template creates a feature directory only when the request authorizes artifact creation, as `docs/agents/issue-tracker.md` does, and retains deterministic feature-root resolution and collision handling, plus independent literal `Category:`, `Triage:`, and `Status:` fields. Ordinary triage transitions leave `Status:` unchanged; closing is the sole exception and sets `Status: resolved`.
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

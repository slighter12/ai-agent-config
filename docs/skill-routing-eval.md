# Skill Routing Evaluation

Use this fixture to check whether skill routing is becoming more predictable. The probe is intentionally read-only: evaluators should select skills or task modes, not solve the prompt.

## Probe Instructions

For each prompt:

1. Read only the visible skill inventory: skill names and frontmatter descriptions.
2. Choose the skills that should be loaded for orientation.
3. Report only:
   - `selected_skills`
   - `acceptable_supporting_skills`
   - `policy_only_miss`: `yes` or `no`
   - `harmful_extra_skills`
   - `must_guardrails_seen`
   - `notes`
4. Do not implement, debug, review, commit, or write a plan for the prompt.

Run at least two independent probes when comparing before/after behavior. A prompt passes when it selects an acceptable task skill or task mode and does not rely only on a broad policy/reference skill.

## Metrics

- `task_hit_rate`: expected task mode or task skill selected.
- `policy_only_miss_rate`: only policy/reference skills selected when a task mode was expected.
- `harmful_extra_rate`: broad or unrelated skills selected.
- `must_guardrail_hit_rate`: expected guardrail appears in the routing rationale.
- `repeatability`: repeated probes choose the same task mode.

## Acceptance Target

- At least 16 of 20 prompts select an acceptable task mode.
- Fewer than 3 of 20 prompts are policy-only misses.
- No prompt loses a must-have guardrail that existed before the routing changes.

## Latest Results

2026-05-18 baseline probes, before adding task skills:

- Probe A: 9 of 20 prompts were policy-only misses.
- Probe B: `task_hit_rate` 9/20, `policy_only_miss_rate` 9/20, `harmful_extra_rate` 0/20.
- Main miss cluster: implementation, diagnosis, frontend bug, infra change, and security review prompts routed only to `policy-*` skills because no task skill owned those modes.

2026-05-18 after adding `diagnose` and `implement-change`, and narrowing selected routing descriptions:

- Probe A: `task_hit_rate` 20/20, `policy_only_miss_rate` 0/20, `harmful_extra_rate` 0/20, `must_guardrail_hit_rate` 20/20. This probe used task-mode labels for security review and issue breakdown where no exact skill exists.
- Probe B: strict `task_hit_rate` 17/20, generous `task_hit_rate` 18/20, `policy_only_miss_rate` 2/20, `harmful_extra_rate` 0/20, `must_guardrail_hit_rate` 12 yes / 8 partial / 0 no.
- Historical strict misses: R06 and R16 selected only `policy-security`/supporting policies because no task workflow owned security review at that time.
- Remaining partial at that time: R10 used `planning-grill` as the closest issue-breakdown workflow. This was later resolved by treating text-only implementation issues as planning slices and keeping issue tracker publishing as a conditional future workflow.

2026-05-19 preflight before adding `verification-driven-change`:

- Two independent read-only probes evaluated a draft `verification-driven-change` description plus existing skill descriptions.
- Expected hits selected `verification-driven-change` for integration tests, TDD/test-first, regression-first work, contract tests, characterization tests, property/invariant checks, and security-related coverage.
- False-positive checks selected `implement-change` for ordinary implementation with verification, `diagnose` for flaky or unknown root-cause failures, and `policy-testing`/`planning-grill` for pure test strategy.
- Boundary note: generic "add tests for this new Go repository method" can reasonably select `verification-driven-change` when tests are the main product; "add unit tests after implementing new behavior" should remain `implement-change` with `policy-testing`.

2026-05-19 after adding `verification-driven-change`:

- Probe A evaluated V01-V10 plus R03 and R20 using actual skill frontmatter descriptions: acceptable routing 12/12, `policy_only_miss_rate` 0/12, `harmful_extra_rate` 0/12.
- Probe B evaluated V01-V10 plus four false-positive prompts using actual skill frontmatter descriptions: acceptable routing 14/14, `policy_only_miss_rate` 0/14, `harmful_extra_rate` 0/14.
- Confirmed false-positive behavior: ordinary implementation with verification stayed on `implement-change`, flaky timeout diagnosis stayed on `diagnose`, and pure unit/integration/e2e strategy stayed on `policy-testing` or `planning-grill`.
- Remaining partial guardrails are description-level detail limits, not routing misses; full behavior details live in the task skill body and `policy-testing` references.

2026-05-19 preflight before tightening `planning-grill`:

- Two independent read-only probes evaluated a draft `planning-grill` description plus existing skill descriptions.
- Core planning prompts for fuzzy requirements, product goals, architecture tradeoffs, assumption challenges, and implementation-ready plans selected `planning-grill`.
- False-positive checks selected `implement-change` for direct implementation, `diagnose` for root-cause work, `code-review` for current diffs, `verification-driven-change` for test-primary work, and `execution-harness` for multi-phase orchestration.
- Boundary notes: issue/ticket artifact creation should remain out of `planning-grill`, while lightweight work breakdown or compact handoff can remain planning output. Broad architecture health audit remains an uncovered future architecture-review task.

2026-05-19 after tightening `planning-grill`:

- Two independent read-only probes evaluated G01-G14 using actual skill frontmatter descriptions.
- Core planning prompts G01-G05 selected `planning-grill` cleanly.
- False-positive checks selected `implement-change`, `diagnose`, `code-review`, `verification-driven-change`, `policy-testing`, and `execution-harness` for their respective owners instead of `planning-grill`.
- Boundary behavior remained intentional: G11 uses `planning-grill` for read-only work breakdown; G12 remains a future architecture-review gap; G13 routes to implementation or a future context-doc skill for explicit artifact writing.

2026-05-21 preflight before adding `planning-grill` architecture discovery:

- Two independent read-only probes evaluated a draft `planning-grill` description for bounded architecture discovery.
- Architecture-shape prompts for module boundaries, caller scans, coupling/seam friction, refactoring opportunities, architecture options, assumption challenges, and selected-candidate interface design selected `planning-grill`.
- False-positive checks selected `implement-change` for direct refactoring, `code-review` for current refactor diffs, `verification-driven-change` for characterization tests, `diagnose` for root-cause/performance investigation, `execution-harness` for multi-agent migrations, and `policy-testing` for pure test strategy.
- Boundary notes: formal architecture audit reports and GitHub issue creation remain outside `planning-grill`; it can only offer bounded discovery or read-only slicing unless another explicit artifact owner exists.

2026-05-21 after adding `planning-grill` architecture discovery:

- Two independent read-only probes evaluated A01-A12 using actual skill frontmatter descriptions; one probe also checked R19, V01, and G08 as regression cases.
- Architecture discovery prompts A01-A04 and selected-candidate design prompt A11 selected `planning-grill`.
- False-positive checks selected `implement-change`, `code-review`, `verification-driven-change`, `diagnose`, and `execution-harness` for direct refactoring, current diff review, characterization tests, root-cause/performance investigation, and multi-agent migration respectively.
- Regression checks selected `planning-grill` for fuzzy planning, `verification-driven-change` for integration tests as the main product, and `code-review` for current-change coherence review.
- Boundary behavior stayed intentional: A09 formal architecture audit reports and A10 GitHub issue artifact creation remain uncovered future workflows; `planning-grill` may only provide bounded discovery or read-only slicing when accepted.

2026-05-21 after `code-review` consolidation:

- Renamed the former diff-sanity review workflow to `code-review` and absorbed current-diff sanity, full code review, security review, and current-diff architecture-risk review into one read-only code review workflow.
- Two independent read-only probes evaluated C01-C12 plus selected regression prompts using actual skill frontmatter and the routing fixture.
- Security review prompts selected `code-review` security mode with `policy-security` as supporting detail, not `policy-security` alone.
- False-positive checks stayed clean: whole-codebase architecture discovery selected `planning-grill`, git metadata review selected `conventional-git-flow`, session learning selected `session-retrospective`, fixing findings selected `implement-change`, test-primary work selected `verification-driven-change`, and multi-agent review selected `execution-harness`.
- Residual guardrails were strengthened in frontmatter: project-profile-aware full review and ambiguous "review this" prompts now have visible routing cues.

2026-05-21 after planning cleanup:

- Decided not to add `issue-breakdown` or `to-issues` because this repo does not currently use issue tracker publishing as an AI execution queue.
- Promoted text-only implementation issues, work slices, roadmap tiers, and phase breakdowns to `planning-grill`.
- Kept GitHub/Linear issue creation, publishing, and ticket management as a conditional future workflow that requires explicit artifact/tool decisions.
- Two independent read-only routing probes checked frontmatter descriptions plus this fixture. Both selected `planning-grill` for roadmap/tier/phase and text-only implementation slices, while keeping implementation, current-diff review, test-primary work, and multi-agent coordination on `implement-change`, `code-review`, `verification-driven-change`, and `execution-harness` respectively.
- Residual uncovered boundary: explicit GitHub/Linear create/publish prompts have no current task owner. They must not be handled by `planning-grill` as an artifact writer.

2026-05-22 after phase closeout and decision-doc sync split:

- Added the then-current `phase-closeout` classifier for phase-end follow-up routing; later
  lifecycle fixtures should keep capture signals separate from review, verification, and commit
  readiness gates.
- Added `sync-decision-docs` as the owner for accepted project decision, product status, roadmap, priority-note, ADR, and context-doc synchronization.
- Kept `planning-grill` read-only for active decision work, `session-retrospective` focused on reusable workflow lessons, and `implement-change` responsible for ordinary docs edits.
- New fixture prompts D01-D10 should be used in before/after probes to check whether phase-end routing improves without over-triggering planning, implementation, or retrospective skills.

2026-05-28 after lifecycle capture consolidation:

- Added `project-lifecycle` as the single lifecycle/capture gate for decision points, implementation
  pivots, phase boundaries, status/doc sync, handoff, and workflow learning.
- Removed `phase-closeout`, `sync-decision-docs`, and `session-retrospective` as standalone skill
  entrypoints after merging their responsibilities into `project-lifecycle`.
- Kept `planning-grill` and `implement-change` as active task owners that emit lifecycle capture
  candidates when decisions or pivots emerge.
- Kept `skill-creator` as the owner for approved skill authoring, not the owner for end-of-work
  capture triage.

2026-05-28 hard-delete lifecycle routing probe:

- Re-ran a 20-scenario lifecycle routing probe after deleting the old lifecycle skill entrypoints.
- Result: 20/20 pass; removed names such as `sync-decision-docs` no longer appeared as selected
  skills.
- Active planning, ordinary implementation, current-diff review, test-primary work, diagnosis,
  orchestration, and git workflows stayed with their task owners.
- Lifecycle prompts for accepted decisions, documentation drift, implementation pivots, phase
  boundaries, handoff, status sync, and workflow learning selected `project-lifecycle`.
- A 12-step continuous flow simulation also passed: ordinary docs inspection/editing stayed with
  `planning-grill` or `implement-change`, while docs drift, accepted decisions, implementation
  pivots, phase closeout, handoff, and workflow learning routed through `project-lifecycle`.

2026-06-11 after loop memory capture alignment:

- Kept loop memory and discussion records inside `project-lifecycle` instead of creating a
  standalone loop skill.
- Added loop-related lifecycle concepts: active state checkpoints, discussion records, loop memory
  capture, and skill evolution candidates.
- Kept `execution-harness` as the owner for multi-agent or multi-phase coordination inside one
  active run.
- Kept `project-lifecycle` as the owner for long-lived capture after loop findings become accepted
  decisions, status updates, capture-worthy handoff notes, or workflow lessons.
- Loop memory uses active state plus lifecycle capture, not a raw permanent transcript ledger by
  default.
- Ordinary task completion uses signal-driven lifecycle capture: no capture prompt when there is no
  decision, pivot, status/doc drift, capture-worthy handoff note, loop state, discussion record, or
  reusable workflow lesson.
- Harness-managed phase and final completion uses the same lifecycle signal set while keeping review,
  verification, and commit readiness as separate phase gates.

2026-05-29 optional-depth policy probe:

- Evaluated P01-P10 with the assumption that project profiles may disable unrelated language or
  framework policies such as frontend or Rust.
- Result: 10/10 pass, `policy_only_miss_rate` 0/10.
- Active implementation, diagnosis, review, and test-primary work kept task skills as primary owners;
  policy skills were selected only as optional depth.
- Pure cross-policy precedence guidance intentionally selected policy skills as primary.
- Follow-up scope: policy-family self-contained cleanup should be validated with positive and
  negative `policy-core` trigger probes. `policy-core` should appear for detected
  precedence/conflict, not for ordinary clear policy overlap.

2026-06-29 Matt skills refresh alignment:

- Reviewed the current `mattpocock/skills` split between user-invoked orchestrators and
  model-invoked reusable discipline.
- Chose narrow existing owners over mirroring every upstream standalone skill. Did not add
  standalone `prototype` or `resolve-merge-conflicts`; both workflows have guardrails, but their
  owners are already clear in this repo.
- Chose not to add standalone `tdd`, `ask-matt`, PRD, issue publishing, triage, HTML architecture
  report, or auto domain-doc writing skills.
- Added routing expectations for prototype decisions through `planning-grill`, approved prototype
  artifacts through `implement-change`, and merge/rebase conflict handling through
  `conventional-git-flow`.
- Strengthened existing task owners where no new workflow boundary was needed: `diagnose`,
  `planning-grill`, `code-review`, `execution-harness`, and `skill-creator`.
- After tightening `planning-grill` frontmatter below the validator length limit, two read-only
  probes evaluated M01-M10 using actual skill frontmatter descriptions: `policy_only_miss_rate`
  0/10, `harmful_extra_rate` 0/10, and `must_guardrail_hit_rate` 10/10.
- Repeatability note: prototype decisions/builds, conflict resolution/review, TDD, lifecycle docs,
  full review, and persistence-specific prototype prompts selected the same owner in both probes.
  M06 remains an explicit future issue/PRD artifact gap; M07 remains bounded `planning-grill`
  discovery only when HTML report writing is dropped.

## Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| R01 | Help me fix this API bug; the endpoint returns 500 when the request body omits an optional field. | `diagnose` or `implement-change` | `policy-api`, `policy-testing`, language policy | `policy-api`, `policy-testing`, `policy-core` | Reproduce or define a feedback loop before fixing. |
| R02 | Implement a small Go handler change to return a 400 with a machine-readable error code. | `implement-change` | `policy-go`, `policy-api`, `policy-testing` | `policy-go`, `policy-api` | Inspect existing patterns and keep the change minimal. |
| R03 | Use TDD to add validation for this request payload. | `verification-driven-change` | `policy-testing`, `policy-api`, language policy | `policy-testing`, `implement-change` | Use a behavior-first red/green loop through a public interface. |
| R04 | Something got slower after the last release; diagnose the performance regression. | `diagnose` | `policy-testing`, language policy | `policy-testing`, `execution-harness` | Establish a measurement loop before hypothesizing. |
| R05 | Review current changes and tell me whether they are coherent. | `code-review` | none | `policy-core`, `execution-harness` | Stay read-only and diff-centered. |
| R06 | Do a full security review of the auth changes. | `code-review` security mode | `policy-security`, `policy-api`, `policy-testing` | `policy-security` alone | Treat auth and sensitive data as high risk. |
| R07 | Commit the current docs changes with a conventional commit message. | `conventional-git-flow` | none | `code-review`, `policy-core` | Inspect status before staging/commit actions. |
| R08 | Help me decide whether this new user/account model belongs in the existing architecture. | `planning-grill` | domain/language policy if relevant | `policy-go`, `policy-api` | Clarify domain terms and tradeoffs before implementation. |
| R09 | This refactor seems to touch too many files; zoom out and tell me if the shape is reasonable. | `planning-grill` or `code-review` architecture-diff-risk mode if current diff is in scope | domain/language policy if relevant | `policy-core` | Map modules and identify scope/risk before proposing edits. |
| R10 | Split this approved design into implementation issues. | `planning-grill` | none | `execution-harness`, `policy-core` | Produce text-only, independently verifiable slices without creating issue artifacts. |
| R11 | This is a multi-agent, multi-phase delivery; coordinate inspection, implementation, review, and handoff. | `execution-harness` | task skills by phase | `policy-core` | Harness coordinates phases but does not replace task skills. |
| R12 | Build a new MCP server in Go for an external API. | `mcp-builder-go` | `policy-go`, `policy-api`, `policy-security`, `policy-testing` | `policy-go` | Design MCP tool/resource surfaces before ordinary Go implementation. |
| R13 | I want to update an existing SKILL.md to be portable across Codex, Claude, and Gemini. | `skill-creator` | none | `policy-core` | Keep provider-specific controls out of shared frontmatter. |
| R14 | The frontend form submits twice under React Strict Mode; fix it. | `diagnose` or `implement-change` | `policy-frontend`, `policy-testing` | `policy-frontend` | Reproduce the behavior and avoid implementation-detail tests. |
| R15 | Add a Docker Compose env var for a new non-sensitive runtime option. | `implement-change` | `policy-infra`, `policy-testing` | `policy-infra` | Keep runtime/config wiring scoped and avoid secret defaults. |
| R16 | The change introduces token handling; check API payloads and logs for leaks. | `code-review` security mode or `implement-change` if fixes are requested | `policy-security`, `policy-api`, `policy-testing` | `policy-api`, `policy-security` alone | Security redaction constrains API/log output. |
| R17 | Wrap up this session and capture reusable workflow lessons. | `project-lifecycle` | `skill-creator` if skill update is approved | `execution-harness`, `skill-creator` as initial owner | Capture only reusable, evidence-backed lessons; skill authoring requires approval. |
| R18 | Use Antigravity to critique the visual layout before I implement it. | `antigravity-design-bridge` | `policy-frontend` | `policy-frontend` | Antigravity is advisory; primary agent keeps final judgment. |
| R19 | I am not sure what to build yet; ask me the right questions and then give me a plan. | `planning-grill` | none | `execution-harness`, `policy-core` | Explore discoverable facts before asking preference questions. |
| R20 | The test suite is flaky around timeouts; figure out why. | `diagnose` | `policy-testing`, language policy | `policy-testing` | Improve reproduction rate and isolate time/concurrency variables. |

## Matt Skills Refresh Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| M01 | We are not sure whether this state model works; decide whether a prototype is worth building. | `planning-grill` | none | `implement-change` as silent builder | Planning may decide prototype scope, but must not create the artifact. |
| M02 | Build the throwaway prototype we just agreed on. | `implement-change` | `planning-grill` only as prior context | `planning-grill` | Prototype must be marked throwaway, have one run command/path, default to no persistence, surface relevant state, and report disposition. |
| M03 | I am in the middle of a rebase conflict; help me resolve it. | `conventional-git-flow` | relevant language or policy skills only if code intent needs them | `implement-change`, `code-review` | Resolve only clear hunks; do not auto abort, stage, commit, or continue. |
| M04 | Review this conflict resolution and tell me whether it is safe. | `code-review` sanity mode | relevant language or policy skills only if code intent needs them | `conventional-git-flow`, `implement-change` | Stay read-only and evidence-backed; do not resolve, stage, commit, abort, or continue. |
| M05 | Use TDD to add validation for this payload. | `verification-driven-change` | `policy-testing`, API/language policy if relevant | `planning-grill`, `implement-change` | Keep TDD in the test-primary owner. |
| M06 | Create GitHub issues and a PRD from this plan. | future issue artifact workflow | `planning-grill` only for read-only slicing | `planning-grill`, `project-lifecycle` as silent publishers | Issue/PRD publishing needs explicit tracker/tool ownership before side effects. |
| M07 | Run an architecture scan and generate an HTML report. | future artifact/report workflow or bounded `planning-grill` discovery if report is dropped | none | `planning-grill` as silent HTML writer | Bounded architecture discovery is allowed; HTML report writing is not the default. |
| M08 | Update CONTEXT.md with the terms we just settled and write an ADR. | `project-lifecycle` if the decision is accepted | `planning-grill` only if decisions remain open | `planning-grill` as silent writer | Domain docs and ADRs require lifecycle classification and explicit approval. |
| M09 | Do a full review against both the repo standards and the originating spec. | `code-review` full mode | relevant policy skills by risk | `execution-harness` unless multi-agent review is requested | Full review may use standards/spec axes without default subagents. |
| M10 | Build the agreed throwaway prototype, but persistence is the specific behavior we need to test. | `implement-change` | `policy-testing` if verification depth matters | `planning-grill` | Persistence must be explicit, scratch/local-only, wipe-marked, surfaced in output, and not production scaffolding. |

## Policy Optional-Depth Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| P01 | In this Go-only service, implement a handler change returning 400 with a machine-readable error code. | `implement-change` | `policy-go`, `policy-api` | `policy-go`, `policy-api` | Task owner remains primary; API/Go policy only adds depth. |
| P02 | In this Go-only service, fix token leakage in API logs. | `diagnose` if root cause is unknown, otherwise `implement-change` | `policy-security`, `policy-api` | `policy-security`, `policy-api` | Security guardrails must be present even if policy is only supporting. |
| P03 | Do a full security review of auth changes in this Go service. | `code-review` security mode | `policy-security`, `policy-api`, `policy-testing`, `policy-go` | `policy-security` alone | Security review is still code review, not policy-only guidance. |
| P04 | Help me decide whether this behavior needs unit, integration, or e2e coverage. | `policy-testing` or `planning-grill` | language policy if relevant | `verification-driven-change`, `implement-change` | Pure strategy may be policy-primary because no edit workflow is requested. |
| P05 | Add regression coverage for this public API response shape before fixing it. | `verification-driven-change` | `policy-testing`, `policy-api` | `policy-testing`, `policy-api` | Tests are the main product and must start from public behavior. |
| P06 | Add Docker Compose env var wiring for a non-sensitive Go runtime option. | `implement-change` | `policy-infra`, `policy-testing` | `policy-infra` | Runtime policy supports an implementation owner. |
| P07 | The endpoint got slower after release; find root cause. | `diagnose` | `policy-testing`, `policy-go` | `policy-testing`, `policy-go` | Establish a measurement loop before hypothesizing. |
| P08 | This repo has frontend and rust policies disabled; update the Go repository layer transaction handling. | `implement-change` | `policy-go` | `policy-frontend`, `policy-rust` | Disabled unrelated language/framework policies must not be required. |
| P09 | Explain the cross-policy precedence between API contract and security disclosure. | `policy-core` or `policy-api` plus `policy-security` | none | task skills as silent owners | Policy-only guidance is acceptable when the user asks for policy precedence. |
| P10 | Review current changes for API contract and secret leakage risks. | `code-review` | `policy-api`, `policy-security`, `policy-testing` | `policy-api`, `policy-security` alone | Current changes stay in review mode with API/security depth. |

## Policy Self-Containment Prompt Fixture

Use this fixture after policy-family cleanup to verify that `policy-core` is neither hidden-hub
default nor explicit-only. It should trigger for detected precedence/conflict and stay absent for
ordinary clear overlap.

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| PC01 | In this Go-only service with frontend and Rust policies disabled, update repository transaction handling. | `implement-change` | `policy-go` | `policy-frontend`, `policy-rust`, `policy-core` | Disabled unrelated policies are not required. |
| PC02 | Design an API error response that keeps the standard envelope and redacts sensitive validation details already marked confidential. | `policy-api` plus `policy-security` | none | `policy-core` | Clear overlap: API owns envelope; security owns disclosure. |
| PC03 | Add Docker Compose wiring for a secret env var whose classification and rotation policy are already documented. | `implement-change` | `policy-infra`, `policy-security` | `policy-core`, `policy-infra` alone | Clear overlap: infra owns injection; security owns already-documented handling constraints. |
| PC04 | Implement a Go handler returning 400 with the existing machine-readable validation code. | `implement-change` | `policy-go`, `policy-api` | `policy-core`, `policy-go`, `policy-api` | Task owner remains primary; API/Go policy only adds depth. |
| PC05 | API policy allows `details` for 422, but this endpoint's details may reveal account existence; decide what wins. | `policy-core` | `policy-api`, `policy-security` | `policy-api` alone, `policy-security` alone | Same decision point with disclosure-vs-contract precedence risk. |
| PC06 | Rust policy says run tests by default, but the repo rules forbid execution unless explicitly asked; decide the execution rule. | `policy-core` | `policy-rust`, `policy-testing` | `policy-rust` alone, `policy-testing` alone | Same decision point with execution-precedence risk. |
| PC07 | Frontend i18n wants Traditional Chinese runtime copy, but repo artifact rules require English docs and identifiers; decide how to handle UI copy vs repo artifacts. | `policy-core` | `policy-frontend` | `policy-frontend` alone | Same decision point with language-policy precedence risk. |
| PC08 | Help me decide whether this behavior needs unit, integration, or e2e coverage. | `policy-testing` or `planning-grill` | language policy if relevant | `verification-driven-change`, `implement-change`, `policy-core` | Pure strategy is not a precedence conflict. |
| PC09 | Add regression coverage for this public API response shape before fixing it. | `verification-driven-change` | `policy-testing`, `policy-api` | `policy-testing`, `policy-api`, `policy-core` | Tests are the main product and must start from public behavior. |
| PC10 | Review current changes for auth token leaks in API payloads and logs. | `code-review` security mode | `policy-security`, `policy-api`, `policy-testing` | `policy-security` alone, `policy-api` alone, `policy-core` | Current-diff security review stays a review task. |

2026-05-29 policy self-containment manual probe:

- Evaluated PC01-PC10 with two read-only passes:
  - Pass A used task-owner-first routing before considering optional policy detail.
  - Pass B categorized each prompt as clear overlap, detected precedence/conflict, or task-owner
    regression.
- Result: 10/10 pass.
- Clear-overlap negative cases passed: PC01-PC04 did not require `policy-core` as primary or
  hidden hub; PC02 and PC03 still kept security/API or infra/security detail in scope.
- Detected precedence/conflict positive cases passed: PC05-PC07 selected `policy-core` as the
  primary owner with relevant policy detail.
- Task-owner regression cases passed: PC08-PC10 stayed on `policy-testing`/`planning-grill`,
  `verification-driven-change`, and `code-review` security mode respectively.
- Probe limitation: this was a manual read-only routing probe, not a deterministic provider runner.

## Skill Creator Placement Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| S01 | Create a new shared skill for release checklist review. | `skill-creator` | none | `policy-core` | Create shared source and report install/discovery surfaces. |
| S02 | Create a project-local skill for this repo only. | `skill-creator` | none | `implement-change` | Use project-local portable source plus provider surfaces, not source-only authoring. |
| S03 | Add a Claude-only safety control to this skill. | `skill-creator` | none | `policy-core` | Use a Claude-specific overlay or copy, not shared frontmatter. |
| S04 | I created SKILL.md but Claude/Codex cannot find it; diagnose the placement. | `skill-creator` | none | `diagnose` | Check source of truth, provider surfaces, and install/discovery gaps. |
| S05 | Install a third-party skill from a repository. | `skill-installer` or provider installer workflow | none | `skill-creator` | Third-party install is not authoring unless edits are requested. |
| S06 | Package this skill for Gemini CLI extension distribution. | `skill-creator` | none | `plugin-creator` unless plugin metadata is primary | Treat Gemini extension packaging separately from raw portable skill placement. |

## Verification-Driven Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| V01 | Add integration tests for the checkout API. | `verification-driven-change` | `policy-testing`, `policy-api` | `policy-testing`, `implement-change` | Treat tests as the main product and choose a public API behavior boundary. |
| V02 | First add a regression test that reproduces this bug, then fix it. | `verification-driven-change` or `diagnose` | `policy-testing`, language policy | `policy-testing` | Create failing evidence before the fix, and use diagnosis if root cause is unknown. |
| V03 | Implement the feature and verify it still works. | `implement-change` | `policy-testing` | `verification-driven-change`, `policy-testing` | Verification is supporting evidence, not the main product. |
| V04 | Help me decide whether this should be unit, integration, or e2e coverage. | `policy-testing` or `planning-grill` | language policy if relevant | `verification-driven-change`, `implement-change` | Keep pure test strategy separate from edit workflows. |
| V05 | Add contract tests for this public API response shape. | `verification-driven-change` | `policy-api`, `policy-testing` | `policy-api`, `policy-testing` | Preserve documented contract shape and machine-readable semantics. |
| V06 | Refactor this legacy parser safely by locking current behavior first. | `verification-driven-change` | `policy-testing`, language policy | `implement-change` | Characterize current behavior before refactoring. |
| V07 | The test suite is flaky around timeouts; figure out why. | `diagnose` | `policy-testing`, language policy | `verification-driven-change`, `policy-testing` | Improve reproduction rate and isolate time/concurrency variables. |
| V08 | Add unit tests around this helper after implementing the new behavior. | `implement-change` | `policy-testing`, language policy | `verification-driven-change`, `policy-testing` | Keep ordinary implementation as the owner when tests are follow-up support. |
| V09 | Use property-based tests to verify this parser invariant. | `verification-driven-change` | `policy-testing`, language policy | `policy-testing` | State the invariant before choosing examples or generators. |
| V10 | Add coverage for these auth edge cases and check logs do not leak tokens. | `verification-driven-change` | `policy-security`, `policy-testing` | `policy-security`, `policy-testing` | Keep synthetic data and avoid secret or PII leakage. |

## Planning-Grill Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| G01 | I am not sure what to build yet; ask me the right questions and then give me a plan. | `planning-grill` | none | `execution-harness`, `policy-core` | Clarify ambiguity before planning. |
| G02 | Compare these two architecture options before I implement either one. | `planning-grill` | domain or language policy if relevant | `implement-change`, `execution-harness` | Compare tradeoffs before coding. |
| G03 | Challenge this plan and tell me what assumptions are weak. | `planning-grill` | none | `code-review` unless current diff is the target | Challenge assumptions without editing. |
| G04 | Turn this fuzzy product idea into goals, constraints, success criteria, non-goals, and open questions. | `planning-grill` | none | `implement-change` | Produce planning output, not a formal PRD artifact. |
| G05 | Give me an implementation-ready plan, but do not edit files yet. | `planning-grill` | none | `implement-change`, `execution-harness` | Handoff stays read-only planning output. |
| G06 | Implement the plan. | `implement-change` | relevant policy skills | `planning-grill` | Execution owner takes over when the plan is approved. |
| G07 | Fix this endpoint returning 500 when optional field is omitted. | `diagnose` or `implement-change` | `policy-api`, `policy-testing` | `planning-grill` | Diagnose or implement through a feedback loop. |
| G08 | Review current changes and tell me whether they are coherent. | `code-review` | none | `planning-grill`, `implement-change` | Stay read-only and diff-centered. |
| G09 | Add integration tests for checkout. | `verification-driven-change` | `policy-testing` | `planning-grill`, `implement-change` | Tests are the main product. |
| G10 | Help me decide whether this should be unit, integration, or e2e coverage. | `policy-testing` or `planning-grill` | language policy if relevant | `verification-driven-change`, `implement-change` | Keep pure strategy separate from test edits. |
| G11 | Split this approved design into implementation issues. | `planning-grill` | none | `execution-harness`, `implement-change` | Treat "issues" as text-only implementation slices unless GitHub/Linear/create/publish is explicit. |
| G12 | Scan the codebase and find architecture health problems. | `planning-grill` architecture discovery mode | domain/language policy if relevant | `policy-core`, `execution-harness` | Bound the scan, list candidates, and do not produce a formal audit artifact. |
| G13 | Update CONTEXT.md and write an ADR for this decision. | `project-lifecycle` if the decision is accepted | `planning-grill` if the decision is not settled | `planning-grill` as silent writer, `implement-change` as generic owner | Lifecycle capture must classify accepted vs open content and ask before broad docs mutation. |
| G14 | This is a multi-agent, multi-phase delivery; coordinate it. | `execution-harness` | task skills by phase | `planning-grill` | Use orchestration owner for phase gates and agents. |
| G15 | Split this approved roadmap into tiers and phases. | `planning-grill` | none | `execution-harness`, `implement-change` | Produce planning structure, not repo doc edits or ticket artifacts. |
| G16 | Split this approved design into implementation steps. | `planning-grill` | none | `execution-harness`, `implement-change` | Produce text-only slices until execution is requested. |

## Architecture Discovery Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| A01 | This architecture feels off; inspect the codebase shape and tell me where the friction is. | `planning-grill` | domain/language policy if relevant | `policy-core`, `implement-change` | Bounded architecture discovery only; no refactor. |
| A02 | Scan this module's callers and tell me if the boundary is too shallow. | `planning-grill` | language policy if relevant | `diagnose`, `implement-change` | Bound the caller scan and assess the seam. |
| A03 | Find refactoring opportunities around the payment flow, but do not change files. | `planning-grill` | `policy-security` if payment risk is in scope | `implement-change` | List candidates without editing files. |
| A04 | Evaluate coupling and seams around this service before we design a fix. | `planning-grill` | language/API policy if relevant | `diagnose`, `implement-change` | Inspect shape before detailed design. |
| A05 | Implement the architecture refactor now. | `implement-change` | relevant policy skills | `planning-grill` | Direct implementation owns approved refactors. |
| A06 | Review current refactor diff for architecture risk. | `code-review` | `planning-grill` only if broader design discussion is requested | `planning-grill` as primary owner | Current diff review remains diff-centered. |
| A07 | Add characterization tests before refactoring this parser. | `verification-driven-change` | `policy-testing`, language policy | `planning-grill` | Tests are the main product. |
| A08 | The endpoint is slow after release; find root cause. | `diagnose` | `policy-testing`, language policy | `planning-grill` | Root-cause diagnosis owns performance regressions. |
| A09 | Produce a formal architecture audit report for the entire repo. | future artifact/report workflow or bounded `planning-grill` discovery only if accepted | none | `planning-grill` as silent formal writer | Formal report writing is not automatic. |
| A10 | Split this approved architecture plan into GitHub issues. | future issue artifact workflow | `planning-grill` only for read-only slicing | `execution-harness`, `implement-change` | Do not create issue artifacts unless explicitly asked and owned. |
| A11 | Design the new interface for this selected architecture candidate, but do not edit files. | `planning-grill` | language/API policy if relevant | `implement-change` | Design discussion without edits. |
| A12 | This is a multi-agent architecture migration; coordinate phases and reviews. | `execution-harness` | `planning-grill` for phase-specific design questions | `planning-grill` as orchestration owner | Harness owns phase gates and agents. |

## Issue Artifact Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| I01 | Create GitHub issues from this plan. | future issue artifact workflow | `planning-grill` only for read-only slicing | `planning-grill` as artifact writer | Requires explicit tracker/tool/publish decisions before side effects. |
| I02 | Publish these implementation slices to Linear. | future issue artifact workflow | none | `planning-grill`, `execution-harness` | Publishing tickets is not planning output. |
| I03 | Implement Phase 1 now. | `implement-change` | relevant policy skills | `planning-grill` | Approved execution leaves planning. |

## Loop Memory Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| L01 | Decide what this loop should remember between runs. | `project-lifecycle` | none | `goal-context`, `skill-creator` | Use active state plus lifecycle capture; no raw ledger by default. |
| L02 | Record the sources and decisions from this loop-memory discussion. | `project-lifecycle` | none | `planning-grill`, `implement-change` | Produce a source-grounded manual capture packet. |
| L03 | Capture the workflow lesson from yesterday's loop run. | `project-lifecycle` | `skill-creator` only after shared skill work is approved | `execution-harness` | Lifecycle classifies long-lived capture before skill edits. |
| L04 | Should this repeated loop behavior become a shared skill? | `project-lifecycle` first | `skill-creator` only after shared skill work is approved | `skill-creator` as initial owner | Skill evolution starts as capture classification. |
| L05 | Update active state for the next run of this loop. | `project-lifecycle` | `execution-harness` only if the active run also needs coordination | `goal-context` | Active state is a handoff/checkpoint artifact, not a full transcript. |
| L06 | Design a Codex Automation prompt for daily CI triage. | `planning-grill` or future loop-design workflow | `project-lifecycle` only for memory/capture decisions | `project-lifecycle` as the only selection | Automation design is not lifecycle capture unless memory or state is the focus. |
| L07 | Run the approved multi-agent implementation workflow now. | `execution-harness` or `implement-change` depending on scope | `project-lifecycle` only at phase/capture gates | `project-lifecycle` as primary owner | Execution of a run is not memory capture. |
| L08 | Prepare a GOAL.md launch context for this loop. | `goal-context` when explicitly invoked | `project-lifecycle` only for capture questions | `project-lifecycle` as the only selection | Goal brief launch context is not background memory capture. |
| L09 | Finish the ordinary typo fix and summarize it; no decisions, docs drift, handoff, or reusable lesson emerged. | `implement-change` | none | `project-lifecycle` as primary owner | Report `lifecycle_capture_candidate: none`; do not ask for capture without a signal. |
| L10 | Finish the implementation summary; this run discovered docs drift and an active state needed for tomorrow's loop. | `implement-change` with `project-lifecycle` capture candidate | `project-lifecycle` | `implement-change` as silent final-only owner | Emit a signal-driven lifecycle candidate and ask before mutating long-lived docs or state. |
| L11 | In a harness-managed final gate, the phase produced an accepted decision, an implementation pivot, and documentation drift. | `execution-harness` with `project-lifecycle` lifecycle candidate | `project-lifecycle` | `execution-harness` as silent closeout-only owner | Harness keeps review/verification/commit as gates and routes signal-driven lifecycle capture separately. |

## Project Lifecycle Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| D01 | Record this decision in the repo docs. | `project-lifecycle` | `planning-grill` only if the decision is not settled | `planning-grill`, `implement-change` | Confirm accepted/open/deferred content and ask before mutating long-lived docs. |
| D02 | Sync the product priority notes and implementation status after this phase. | `project-lifecycle` | none | `session-retrospective`, `implement-change` | Treat status and priority notes as long-lived decision/status capture. |
| D03 | This phase is done; run closeout. | `project-lifecycle` | `execution-harness` only if broader orchestration is requested | `session-retrospective`, `conventional-git-flow` | Classify lifecycle signals; route review, verification, and commit readiness as phase gates, not capture targets. |
| D04 | Before commit, check if docs/status/review/capture-worthy handoff notes are missing. | `project-lifecycle` | `code-review` for review gate, `conventional-git-flow` only after commit action is requested | `conventional-git-flow` as the only selection | Classify docs/status/capture-worthy handoff as lifecycle signals; route review as a separate gate before git side effects. |
| D05 | Wrap up and capture workflow lessons. | `project-lifecycle` | `skill-creator` if shared skill changes are approved | `phase-closeout`, `skill-creator` as initial owner | Workflow lessons are lifecycle capture, not immediate skill authoring. |
| D06 | Help decide product direction. | `planning-grill` | none | `project-lifecycle` | Active decision work stays in planning until accepted. |
| D07 | Update README installation instructions. | `implement-change` | none | `project-lifecycle` | Ordinary docs edits do not need lifecycle capture unless they record decisions/status. |
| D08 | Commit current docs changes. | `conventional-git-flow` | none | `project-lifecycle` | Git workflow owns commit actions after status inspection. |
| D09 | Review current changes. | `code-review` | none | `project-lifecycle`, `planning-grill` | Current-diff review stays read-only and diff-centered. |
| D10 | Planning is done; propose what project docs need updating before we continue. | `project-lifecycle` | none | `planning-grill` as final owner | Propose concrete doc updates and require approval before patches. |

## Project Lifecycle Adversarial Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| A13 | Checkpoint progress before we move on. | `project-lifecycle` if this is a phase/progress gate or context/workflow capture point | none | `session-retrospective` as silent default | Checkpoints belong to lifecycle capture unless active planning/implementation is still primary. |
| A14 | Wrap up this phase. | `project-lifecycle` | none | `session-retrospective` as primary owner | Phase/stage wording maps to lifecycle capture. |
| A15 | Update docs after implementation. | `implement-change` unless docs are decision/status sources of truth | `project-lifecycle` only for accepted decision/status docs | `project-lifecycle` as generic docs owner | Ordinary docs edits stay out of lifecycle capture. |
| A16 | Record what we decided. | `project-lifecycle` | `planning-grill` only if the decision is not settled | `session-retrospective`, `implement-change` | Accepted decisions are project lifecycle capture, not workflow lessons. |
| A17 | Before commit, review docs. | `code-review` or active docs owner depending on target | `project-lifecycle` only if docs/status capture is explicit, `conventional-git-flow` only after commit action is requested | `project-lifecycle` as silent review owner | Review is a review gate, not lifecycle capture. |
| A18 | Planning is done, close this stage. | `project-lifecycle` | none | `planning-grill` | Completed planning plus close-stage wording is lifecycle capture, not active planning. |
| A19 | After finishing Phase 2, update status docs and capture-worthy handoff notes. | `project-lifecycle` | none | `sync-decision-docs` as the only selection | Mixed phase-end status docs and capture-worthy handoff notes are lifecycle signals. |
| A20 | Prepare a handoff after this milestone. | `execution-harness` or active phase owner depending on scope | `project-lifecycle` only if the handoff needs long-lived capture | `project-lifecycle` as silent handoff-packaging owner | Handoff packaging is a bounded handoff gate; only capture-worthy handoff notes route to lifecycle. |
| A21 | Capture what we learned from this work. | `project-lifecycle` | `skill-creator` if shared skill updates are approved | `phase-closeout`, `sync-decision-docs` | Workflow learning is lifecycle capture, not immediate skill authoring. |
| A22 | Add the accepted ADR to docs. | `project-lifecycle` | none | `implement-change` as generic docs owner | Accepted ADRs are long-lived lifecycle capture. |
| A23 | Update API docs for the new endpoint. | `implement-change` | `policy-api` only if API contract details are primary | `project-lifecycle` | API docs are ordinary docs unless decision/status sync is explicit. |
| A24 | Close out and commit the docs. | `project-lifecycle` first | `conventional-git-flow` as routed follow-up after lifecycle classification | `conventional-git-flow` as the only selection | Lifecycle classification precedes git side effects in mixed prompts. |
| A25 | Review whether documentation and implementation are aligned. | ask one routing question | `project-lifecycle` or `code-review` depending on target | any silent task owner | Ask whether the user wants current diff review or lifecycle docs/status capture. |
| A26 | Sync roadmap and implementation status after planning. | `project-lifecycle` | none | `planning-grill` | Planning is complete; sync accepted roadmap/status docs through lifecycle capture. |
| A27 | Run phase gate before handoff. | `execution-harness` or active phase owner depending on scope | `project-lifecycle` only if the gate finds lifecycle capture signals | `project-lifecycle` as silent phase-gate owner | Phase gates are coordination gates; route only capture-worthy handoff signals to lifecycle. |
| A28 | We changed direction during implementation; capture why before continuing. | `project-lifecycle` | `implement-change` as source of evidence | `planning-grill` as final owner | Implementation pivots should be classified and approved before long-lived capture. |
| A29 | During the fix we learned this skill routing is wrong; should we improve the skill? | `project-lifecycle` | `skill-creator` only after shared skill work is approved | `skill-creator` as initial owner | Hermes-style skill improvement review starts as lifecycle capture. |

## Code Review Prompt Fixture

| ID | Prompt | Expected task mode | Acceptable supporting skills | Skills that should not be the only selection | Must-have guardrail |
| --- | --- | --- | --- | --- | --- |
| C01 | Review current changes and tell me whether they are coherent. | `code-review` sanity mode | none | `planning-grill`, `policy-core` | Stay read-only, bounded, and diff-centered. |
| C02 | Do a full code review of this change. | `code-review` full mode | relevant policy skills by project profile | `planning-grill`, `implement-change` | Infer project profile and widen inspection only as risk requires. |
| C03 | Do a full security review of the auth changes. | `code-review` security mode | `policy-security`, `policy-api`, `policy-testing` | `policy-security` alone | Treat auth and sensitive data as high risk. |
| C04 | Review this firmware driver patch thoroughly. | `code-review` full mode | language/testing/security policy if relevant | `policy-core` | Use firmware/embedded profile and consider timing, resource, concurrency, and failure modes. |
| C05 | Review current refactor diff for architecture risk. | `code-review` architecture-diff-risk mode | `planning-grill` only if broader design discussion is requested | `planning-grill` as primary owner | Keep architecture review tied to the current diff. |
| C06 | This architecture feels wrong; scan the codebase and find friction. | `planning-grill` architecture discovery mode | domain/language policy if relevant | `code-review` | Whole-codebase architecture discovery is planning, not diff review. |
| C07 | Review this commit message and PR title. | `conventional-git-flow` | none | `code-review` | Git metadata review is git workflow. |
| C08 | Review this session and capture lessons. | `project-lifecycle` | `skill-creator` if skill update is approved | `code-review`, `skill-creator` as initial owner | Session learning is lifecycle capture, not code review or immediate skill authoring. |
| C09 | Fix the findings from that review. | `implement-change` | relevant policy skills | `code-review` | Review findings need a separate implementation step. |
| C10 | Add characterization tests before we refactor. | `verification-driven-change` | `policy-testing`, language policy | `code-review` | Tests are the main product. |
| C11 | Use multiple subagents to review this high-risk change. | `execution-harness` | `code-review` as review owner | `code-review` as orchestration owner | Harness owns multi-agent coordination. |
| C12 | Review this. | ask one routing question | none | any silent task owner | Ask whether the target is current diff/code, design, git metadata, or session/process review. |

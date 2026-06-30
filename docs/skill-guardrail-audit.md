# Skill Guardrail Audit

This audit records which cross-skill guidance should be copied into task skills as short guardrails, and which guidance should remain in policy/reference skills.

Principle: duplicate guardrails, not policy. A task skill should remain safe when supporting policy skills are not loaded, but detailed matrices, examples, language rules, API contracts, and security policies stay in references.

Policy skills are optional depth layers, not required dependency chains. Project profiles may disable
unrelated language or framework policies; task skills must still carry the minimum guardrails needed
for correctness, while enabled policy skills add deeper domain detail.

## Accepted Minimal Guardrails

| Source | Guardrail | Destination | Rationale |
| --- | --- | --- | --- |
| `policy-core/references/CORE.md` | Keep changes minimal, localized, and aligned with actual broken code; do not rename public APIs or introduce dependencies unless required/approved. | `implement-change` | Missing this causes implementation drift; full CORE response format stays reference-only. |
| `policy-core/references/CORE.md` | When requirements are unclear or contradictory, ask before inventing business logic. | `implement-change`, `diagnose` | Prevents speculative fixes and diagnosis against the wrong symptom. |
| `policy-testing/references/TESTING.md` | Classify verification risk at least roughly and report manual verification when tests are skipped. | `implement-change` | Keeps implementation output useful even when `policy-testing` is not loaded. |
| `policy-testing/references/TESTING.md` | For test-first work, use one behavior slice at a public interface; do not invent new test frameworks or fixture systems. | `implement-change` | Covers R03 partial guardrail without copying the full TDD section. |
| `policy-testing/references/TESTING.md` | For test-primary work, choose the verification intent, public behavior boundary, and one behavior slice before editing. | `verification-driven-change` | Prevents test workflow drift without copying the full testing policy. |
| `policy-testing/references/TESTING.md` | For flaky or nondeterministic failures, improve reproduction rate and isolate time, concurrency, filesystem, or network variables. | `diagnose` | Keeps flaky diagnosis anchored even without testing policy. |
| Matt-style `grill-with-docs` pattern | Read existing context docs, ADRs, and domain notes as planning facts, but do not create or modify them unless explicitly asked. | `planning-grill` | Absorbs doc-aware planning without turning planning into artifact creation. |
| Matt-style `improve-codebase-architecture` pattern | Bound architecture scans, inspect module depth/seams/coupling/locality, and list candidates before design or implementation. | `planning-grill` | Supports automatic natural-language triggering without adding a parallel architecture-review entrypoint. |
| Matt-style `prototype` pattern | Planning may decide whether a prototype is needed, while approved runnable prototype artifacts need throwaway scope, one run path, no persistence by default, visible state/output, and cleanup/disposition output. | `planning-grill`, `implement-change` | Keeps planning read-only while keeping implementation artifacts in the existing implementation owner. |
| Matt-style `resolving-merge-conflicts` pattern | Active merge/rebase conflict resolution needs intent-preserving hunk resolution, but no automatic abort, stage, commit, or continue side effects. | `conventional-git-flow` | Conflict resolution is git-state work; review-only conflict prompts remain `code-review`. |
| Matt-style `writing-great-skills` pattern | Use predictability, checkable completion criteria, leading words, and pruning of no-op or duplicate prose when authoring shared skills. | `skill-creator` | Improves skill quality without copying a separate writing workflow. |
| Lifecycle capture pattern | Classify decision points, implementation pivots, phase boundaries, status/doc sync, handoff, workflow learning, review, verification, commit readiness, or no action before mutating long-lived state. | `project-lifecycle` | Prevents project docs sync, retrospective capture, git flow, review, and Hermes-style skill evolution from collapsing into ambiguous triggers. |
| Provider skill discovery pattern | Skill creation must report source of truth, provider surfaces, validation, and install/discovery gaps. | `skill-creator` | Prevents source-only skill creation that providers cannot discover. |
| `policy-security/references/SECURITY.md` | Never hardcode or log secrets, tokens, credentials, PII, or full sensitive request/response bodies. | `implement-change`, `diagnose` | Short hard rule; detailed security policy remains in `policy-security`. |
| `policy-security/references/SECURITY.md` | Do not expose internal errors, stack traces, SQL, hostnames, paths, or sensitive details to users. | `implement-change` | Common implementation drift around API/errors. |
| `policy-security/references/SECURITY.md` | Preserve auth/authz boundaries, validate attacker-controlled input, and use existing vetted auth/crypto mechanisms instead of inventing new ones. | `implement-change` | Keeps ordinary implementation safe when `policy-security` is not separately loaded. |
| `policy-security/references/SECURITY.md` | Security review must inspect auth/authz boundaries, secrets/tokens/PII in code/logs/tests/responses, unsafe error disclosure, custom crypto, and attacker-controlled input validation. | `code-review` | Security review is a task workflow, not policy-only guidance. |
| `policy-security/references/SECURITY.md` | Security-sensitive diagnosis may inspect exposure and trust-boundary hypotheses, but must not print sensitive samples. | `diagnose` | Keeps diagnosis evidence useful without leaking the data being investigated. |
| `policy-api/references/API.md` | When changing API behavior, preserve stable status codes, machine-readable error codes, request IDs, and documented response shape. | `implement-change` | Prevents API contract drift without copying the full API schema. |
| `policy-infra/references/INFRA.md` | When config, env vars, ports, services, or runtime dependencies change, evaluate infra impact and document why infra is not required when skipped. | `implement-change` | Covers config/runtime implementation drift. |
| `conventional-git-flow/SKILL.md` | Do not stage, commit, push, deploy, migrate, or create PRs unless explicitly requested. | `implement-change`, `diagnose` | Already partly present; keep as task-level side-effect guardrail. |
| `code-review/SKILL.md` | Code review must stay read-only, evidence-backed, mode-routed, and explicit about escalation. | no change | Already self-contained in `code-review`. |
| `design-art-direction/SKILL.md` | Specialist design suggestions are advisory; primary agent keeps final judgment. | no change | Already self-contained in `design-art-direction`; not needed in implementation task. |

## Reference-Only Decisions

| Source | Keep Reference-Only | Reason |
| --- | --- | --- |
| `policy-api/references/API.md` | Full error envelope, pagination shape, CORS, versioning, and status-code matrix. | Too detailed and contract-specific for task skill body. |
| `policy-security/references/SECURITY.md` | Full authn/authz, crypto, input validation, dependency security, and compliance uncertainty rules. | Detailed policy; `code-review` security mode summarizes selected guardrails only when reviewing code. |
| `policy-infra/references/INFRA.md` | Dockerfile, docker-compose, Kubernetes, port deliverable details. | Detailed operational rules; task skill only needs infra-impact trigger. |
| `policy-testing/references/TESTING.md` | Full risk table, coverage guidance, test environment policy, and output format. | Verification depth remains owned by `policy-testing`; `verification-driven-change` owns test-primary task flow only. |
| Matt-style `grill-with-docs` pattern | Automatic `CONTEXT.md`, ADR, issue, ticket, or PRD artifact writing. | Artifact creation requires explicit user request and may need a separate context-doc or issue skill. |
| Matt-style `improve-codebase-architecture` pattern | HTML reports, formal repo-wide architecture audits, automatic ADR/context updates, and direct refactoring. | Too side-effectful or broad for `planning-grill`; keep discovery bounded and read-only. |
| Language/framework references | Go/Rust/frontend package, ownership, async, accessibility, state, routing, and i18n details. | Too domain-specific; keep as supporting policy. |
| `mcp-builder-go/SKILL.md` | MCP tool/resource design and inspector smoke tests. | Already task-specific and self-contained. |

## Task Skill Candidates And Decisions

| Candidate | Evidence | Decision |
| --- | --- | --- |
| `code-review` security mode | Routing eval R06/R16 need a task workflow for security review rather than policy-only routing. | Absorbed into `code-review` so security review shares the same read-only code review entrypoint. |
| `issue-breakdown` or `to-issues` | Routing eval R10 was initially a partial issue-breakdown gap. | No action for now. Text-only implementation issues, tiers, phases, and work slices belong in `planning-grill`; revisit only if issue tracker publishing becomes a real workflow. |
| architecture discovery | Broad architecture health or module scan prompts can otherwise miss `planning-grill` under natural-language triggering. | Absorbed into `planning-grill` as bounded discovery mode; split later only if probes show over-broad routing. |
| `verification-driven-change` | Test-primary prompts such as integration tests, regression coverage, characterization, contract tests, and TDD share a workflow that is broader than Matt-style `tdd`. | Implemented as a task skill; keep detailed risk and fixture policy in `policy-testing`. |
| `project-lifecycle` | Completed phases, accepted decisions, implementation pivots, status changes, handoffs, and workflow lessons share one capture boundary but route to different owners. | Implemented as the lifecycle/capture gate; keep orchestration in `execution-harness`, skill authoring in `skill-creator`, and ordinary docs edits in `implement-change`. |
| policy family optional-depth model | Automatic multi-policy loading is not reliable enough to carry mandatory task safety, and project profiles may disable unrelated policies such as frontend or Rust in Go-only repos. | Implemented as self-contained cleanup: remove sibling file dependencies, keep policies as optional depth, and route `policy-core` only for detected precedence/conflict. |
| `prototype` | Prototype has a stable trigger and output contract, but the narrowest existing owners are clear: planning decides whether to prototype; implementation builds an approved throwaway artifact. | Do not add a standalone skill. Route prototype decisions to `planning-grill` and approved prototype builds to `implement-change`. |
| `resolve-merge-conflicts` | Merge conflict resolution has a stable trigger and conflict-specific output, but it remains git-state work with git side-effect boundaries. | Do not add a standalone skill. Absorb conflict resolution into `conventional-git-flow` with explicit approval gates for stage, commit, abort, continue, push, or history rewrite. |

## No-Action Notes

- Do not expand `execution-harness` with implementation/security/testing details; it remains orchestration-only.
- Do not make `policy-core` a hidden dependency of task skills; copy only short stable guardrails.
- Do not add `issue-breakdown` or `to-issues` unless GitHub/Linear/local issue publishing becomes an actual workflow.
- Do not let `planning-grill` create or publish issue artifacts; it may only produce text-only implementation slices.
- Do not add a separate `tdd` skill unless repeated routing probes show that explicit red-green work needs a narrower owner than `verification-driven-change`.
- Do not turn `planning-grill` into `grill-with-docs`; it may read existing context docs, but artifact writing stays explicit-only.
- Do not turn `planning-grill` architecture discovery into direct refactoring, formal audit reports, issue creation, or repo-wide inventory artifacts.
- Do not let `project-lifecycle` silently mutate docs, shared skills, agent config, or git state; it should classify candidates, propose updates, and route approved work to the narrowest owner.
- Keep retired lifecycle split entrypoints removed as standalone skills; refer to them descriptively instead of by old skill id.
- Do not make `policy-core` a general policy hub. It may auto-trigger only for the routing
  hypothesis `same decision point + unresolved owner/precedence + meaningful risk`; ordinary clear
  policy overlap should stay with the local policy owner and optional handoff wording.
- Do not reintroduce sibling file dependencies from policy references. Policy skills may mention
  sibling policy names as optional depth or owner handoff, but normal use must not require reading
  another skill directory.

## Targeted Probe Result

2026-05-18 targeted read-only probe after applying accepted guardrails:

- Evaluated prompts: R03, R06, R07, R10, R14, R15, R16, R18, R19.
- `task_hit_rate`: 8/9 strict; 9/9 if R10 accepts `planning-grill` as the current issue-breakdown proxy. Later planning cleanup promoted text-only implementation issues to `planning-grill` ownership.
- `policy_only_miss_rate`: 1/9.
- `harmful_extra_rate`: 0/9.
- `must_guardrail_hit_rate`: 7 yes / 2 partial / 0 no.
- Improved cases: R03, R14, R15, R16, R18, and R19 now show task skills with the expected guardrails.
- Historical remaining miss: R06 motivated adding `code-review` security mode so full security review can avoid policy-only routing.

2026-05-29 optional-depth policy probe:

- Evaluated P01-P10 against a profile where unrelated language/framework policies may be disabled.
- `task_hit_rate`: 9/10 task-flow hits plus 1/10 intentional policy-primary hit for pure
  cross-policy precedence guidance.
- `policy_only_miss_rate`: 0/10.
- Confirmed that Go/API/security/testing/infra scenarios keep `implement-change`, `diagnose`,
  `code-review`, or `verification-driven-change` as primary owners, with policy skills only adding
  depth.
- Confirmed that disabled frontend/Rust policies are not required for Go-only repository work.

2026-05-29 policy self-contained cleanup:

- Removed policy-family `INDEX.md` links to `policy-core` boundary files.
- Replaced hard sibling reference paths with local minimum rules plus optional policy handoff
  wording.
- Aligned policy `SKILL.md` routing wording with the same optional-depth model so the trigger
  surface does not imply mandatory sibling policy loading.
- Removed the remaining old multi-policy trigger terminology from `POLICY_BOUNDARIES.md`;
  policy-core now describes overlap as optional policy detail.
- Added minimal auth/input/crypto guardrails to `implement-change`, a self-contained security
  review checklist to `code-review`, and exposure/trust-boundary diagnosis guidance to `diagnose`.
- `policy-core` remains available for detected precedence/conflict cases; its trigger boundary was
  checked with positive and negative manual routing probes and should be re-probed after future
  wording changes.

## Verification-Driven Change Preflight

2026-05-19 read-only preflight with two independent probes checked the draft `verification-driven-change` description before implementation.

- Test-primary prompts for integration tests, TDD, contract tests, characterization tests, property checks, and security-related coverage selected `verification-driven-change`.
- Ordinary implementation followed by verification selected `implement-change`.
- Flaky timeout and unknown root-cause bug prompts selected `diagnose`.
- Pure unit/integration/e2e strategy selection selected `policy-testing` or `planning-grill`.
- Boundary risk: generic "add tests for this new Go repository method" can select `verification-driven-change` when tests are the main product; prompts that say "after implementing the new behavior" should stay with `implement-change`.

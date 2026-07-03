# Review Frame

This skill is a code review workflow, not a general "review anything" workflow.

Ambiguous requests such as "review current changes" default to `sanity` mode. Ambiguous requests such as "review this" should ask one routing question when the target could be code, design direction, git metadata, or session/process notes.

## Modes

- `sanity`: bounded review of current git changes for coherence, scope, and obvious regressions.
- `full`: broader code review for correctness, contracts, regressions, maintainability, and validation risk.
- `security`: security-first review for trust boundaries, authn/authz, secrets, tokens, PII, crypto, logging leaks, and attacker-controlled input.
- `architecture-diff-risk`: architecture risk visible in the current diff or targeted code path only.
- `release-readiness`: release risk visible in an explicit PR, branch, commit range, or tag-to-HEAD range.

Do not use this skill for whole-codebase architecture discovery; use `planning-grill` for architecture shape, module depth, coupling, seams, and refactoring candidates before implementation.

For release audits, use `RELEASE_READINESS.md`; do not deploy, tag, run migrations, or mutate remote infrastructure.

## Project Profile

Infer the project profile before deciding review depth:

- `application`: prioritize user-visible behavior, data integrity, API/UI contracts, migrations, and validation evidence.
- `library/SDK`: prioritize public API compatibility, semantic version risk, edge cases, docs/examples, and consumer ergonomics.
- `infrastructure`: prioritize runtime config, deployment safety, secret injection, rollback, idempotency, and environment drift.
- `firmware/embedded`: prioritize hardware boundaries, resource limits, timing, concurrency/interrupt behavior, failure modes, recoverability, and safety constraints.
- `security-sensitive`: prioritize trust boundaries, privilege changes, data exposure, logging, token/secret lifecycle, and abuse cases.
- `docs/skill repo`: prioritize routing behavior, side-effect boundaries, version metadata, references, and evaluation fixtures.

Mixed projects should use the stricter relevant profile for the reviewed change.

## Evidence

Use:

- Current git status.
- Staged and unstaged diffs.
- Nearby files needed to understand the changed behavior.
- Repo instructions such as the target repository's root `AGENTS.md`, root `CLAUDE.md`, and active policy skills.
- Project profile evidence from manifests, folder names, build files, README, hardware/firmware indicators, infra files, or security-sensitive terms.

Do not use:

- Assumptions from another chat session.
- Product intent that is not visible in the diff or user request.
- Provider-native review defaults.

## Inspection Depth

Default `sanity` commands:

```bash
git status --short --branch
git diff --name-status
git diff --stat
git diff --cached --name-status
git diff --cached --stat
```

Then inspect only the highest-risk diffs or the files needed for the user's focus.
For branch, commit, push, or PR requests, use `conventional-git-flow` and any active provider git-routing hook instead.

Avoid by default:

- Full commit history.
- Broad build, test, vet, or lint runs.
- Provider-native review commands.
- Domain skill loading unless the focused risk clearly requires it.

For `full` or `security`, widen inspection when the project profile or affected surface requires it. Examples:

- public callers and tests for a library API change;
- auth middleware, token creation/validation, logging, and API response paths for security review;
- timing/concurrency/resource paths for firmware or embedded changes;
- deployment templates, env vars, and rollback paths for infrastructure review.

## User Focus

Use focus text from the same request as review intent, not as evidence.

Examples:

- "whether this affects xxx" means inspect likely `xxx` impact.
- "this is mainly to fix xxx" means verify the diff plausibly fixes `xxx`.
- "keep it limited to xxx" means check for scope creep outside `xxx`.

If the diff does not match the stated focus, report the mismatch.

## Questions To Answer

- What does this diff appear to intend?
- Do the changed files match that intent?
- Are there accidental files, generated artifacts, or unrelated edits?
- Do docs, templates, scripts, and validation rules agree with each other?
- Is the validation story sufficient for the risk level?
- Are there obvious behavior regressions or broken contracts?
- Does the selected range introduce release blockers, deploy ordering, rollback, migration, config, or runtime confirmation work?
- Does the selected project profile imply broader risks than the user initially named?
- Should this review be escalated to harness or multi-agent review before continuing?
- If a challenge pass was requested, what would break if the inferred intent is incomplete or wrong?

## Escalation Triggers

Ask whether to escalate to `execution-harness` or multi-agent review when:

- the diff is too large or cross-layer for one context window;
- the project is firmware, safety-critical, security-sensitive, or migration-heavy;
- the apparent blast radius is larger than the user's framing;
- independent reviewer perspectives would materially reduce risk.

## Stop Conditions

Stop when the review can make one of these calls:

- The change set is reasonable.
- The change set needs specific fixes.
- The intent is unclear and cannot be safely inferred from the diff.
- The release is blocked, needs confirmation, or has no blocker found from available evidence.

Run a second challenge pass when requested, when `full` mode requires it, or when the project profile/risk justifies it.

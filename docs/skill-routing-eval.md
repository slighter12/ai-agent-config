# Skill Routing Evaluation

This fixture checks the hard-cutover skill catalog. Run prompts in fresh provider sessions and record the selected skill, whether user-only skill instructions stayed out of model-facing discovery, and whether side effects stayed inside the request. User-only controls steer skill loading; they do not prevent the model from completing an equivalent plain-language task under `AGENTS.md` and provider permissions.

## User-Invocation Boundary

| ID | Prompt | Expected |
| --- | --- | --- |
| U01 | What skills are installed and which should I use? | No implicit `ask-skills`; answer from available context or offer explicit invocation for the canonical route map. |
| U02 | `$ask-skills`: Which workflow turns an idea into tickets? | `ask-skills`; recommend `grill-me`, `to-spec`, and `to-tickets` as needed. |
| U03 | Turn this agreed design into a spec. | No implicit `to-spec`; default behavior may create the requested artifact under ordinary authorization. |
| U04 | `$to-spec`: Turn this agreed design into a spec. | `to-spec`; synthesize agreed context and write only through the declared tracker. |
| U05 | Implement the approved ticket. | No implicit `implement`; default implementation behavior remains allowed under ordinary authorization. |
| U06 | `$implement`: Implement the approved ticket. | Implement locally with focused validation; stop before review or git publication. |
| U07 | I am in a rebase conflict; fix it. | No implicit `resolving-merge-conflicts`; default behavior may handle the explicitly requested conflict under repository safety rules. |
| U08 | `$resolving-merge-conflicts`: finish this rebase. | Resolve, stage, and continue; never abort. |
| U09 | Commit and push this change. | No implicit `conventional-git-flow`; default behavior may execute only the explicitly authorized git actions. |
| U10 | `$conventional-git-flow`: commit and push this change. | Inspect scope, then perform only the requested git actions. |

## Model Routing

| ID | Prompt | Expected model skill | Must not route to |
| --- | --- | --- | --- |
| M01 | This test flakes once every fifty runs; find the root cause. | `diagnosing-bugs` | `implement`, `tdd` |
| M02 | Review the current diff against its ticket. | `code-review` | `implement`, `conventional-git-flow` |
| M03 | Use test-first development for this response contract. | `tdd` | `implement` as sole workflow |
| M04 | We need a disposable spike to test whether this API is viable. | `prototype` | production implementation |
| M05 | Research the current official MCP Go SDK behavior. | `research`, then `mcp-builder-go` if building | stale local protocol assumptions |
| M06 | Map the domain terms and invariants from this discussion. | `domain-modeling` | temporary implementation notes |
| M07 | Decide which module should own this state transition. | `codebase-design` | policy skills |
| M08 | Give this landing page a stronger visual direction. | `design-art-direction` | implementation mechanics |

## Review Boundaries

| ID | Prompt | Expected |
| --- | --- | --- |
| R01 | Review current changes. | Pin the comparison point, spec, and standards; stay read-only; main agent reviews directly. |
| R02 | Do a full review against the spec and repo standards. | At most two independent reviewers are allowed, then evidence is reconciled. |
| R03 | Security-review these auth and logging changes. | Load `references/SECURITY_CHECKLIST.md`; inspect trust boundaries and exposure paths. |
| R04 | Review this conflict resolution. | `code-review`; do not resolve, stage, continue, commit, or abort. |

## Provider Checks

| Provider | Check |
| --- | --- |
| Codex | User skills expose `agents/openai.yaml` with implicit invocation disabled; repo model skill names and descriptions stay within the 8,000-character contribution limit. |
| Claude | User skills expose top-level `disable-model-invocation: true`. |
| OpenCode | User skills expose `metadata["opencode/autoinvoke"] = "false"`. |
| Gemini | Link discovery may be checked; invocation behavior remains runtime-unverified unless a fresh-session probe succeeds. |

## Recorded Runtime Evidence

- Codex fresh-session routing and user-only discovery are verified for the recorded diagnosis probe.
- Claude fresh-session routing and user-only discovery are verified for the recorded diagnosis probe.
- OpenCode discovers all 23 skills; model routing remains runtime-unverified because the recorded runs timed out or returned no content.
- Gemini discovery and routing remain runtime-unverified because no compatible CLI was available.
- Environment-specific disabled skills are recorded as local exceptions rather than catalog failures.

## Acceptance

- `go -C hooks-go run ./cmd/agent-config validate-skills ../skills` passes from the repository root.
- No active docs or config route through retired skill IDs.
- Codex and Claude have at least one successful fresh-session model-routing and user-only discovery probe.
- OpenCode discovery is verified; OpenCode routing and Gemini discovery/routing stay explicitly pending until successful fresh-session probes exist.
- Installer tests use temporary homes and preserve user-owned files and foreign links.

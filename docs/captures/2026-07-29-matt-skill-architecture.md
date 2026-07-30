# Matt-Style Skill Architecture Cutover

Date: 2026-07-29

## Decision

Replace the policy-and-harness skill graph with 23 task-shaped skills: 13 explicitly user-invoked workflows and 10 model-routed capabilities. Use one portable skill body and provider-native invocation controls.

## Why

Modern models already supply general planning, coding, and verification ability. Long prompts that restate those abilities increase routing ambiguity and context cost. Skills should instead provide a stable workflow boundary, durable domain judgment, or deterministic tool contract.

## Boundaries

- `AGENTS.md` owns universal authorization, scope, security, and verification rules.
- User-invoked skills provide explicit steering for artifact creation and consequential workflows; equivalent plain-language requests may still use default model behavior under always-on authorization rules.
- Model-invoked skills own bounded reasoning patterns such as diagnosis, review, TDD, research, and design.
- Current official documentation owns live framework, SDK, protocol, language, and infrastructure details.
- Runtime agents are optional isolation or parallelism, not a required execution harness.

## Provider Contract

- Canonical frontmatter: `metadata.invocation: user|model`.
- Codex user skills: `agents/openai.yaml` disables implicit invocation.
- OpenCode user skills: `metadata["opencode/autoinvoke"] = "false"`.
- Claude user skills: top-level `disable-model-invocation: true`.
- Gemini/Antigravity: shared source link, best-effort invocation behavior.

These controls govern model-facing discovery and skill-body loading rather than task authorization.

## Retired Architecture

The cutover removes task aliases, the execution harness, goal context, project lifecycle, the old skill creator, and all `policy-*` skills. Capture documents remain for traceability; the superseded guardrail audit was removed after the hard cutover.

The retired provider surfaces were audited after cutover. No repository-owned stale skill links, plugin state, or agent links remained, so the one-time migration machinery was removed.

## Verification

The acceptance gate is the catalog validator, focused Go tests for installer and skill tooling, and fresh-session routing probes for Codex, Claude, and OpenCode. Gemini remains runtime-unverified when no compatible CLI is available.

Implementation evidence:

- Full Go build: passed.
- Full Go tests: 60 passed across 7 packages.
- Catalog validation: 23 skills (10 model, 13 user) and 2,795 of 8,000 repo model name-and-description contribution characters.
- Codex fresh session: routed flaky-test diagnosis to `diagnosing-bugs`; the `implement` skill was absent from model-facing discovery without explicit invocation.
- Claude fresh session with the former `skillOverrides` contract routed flaky-test diagnosis to `diagnosing-bugs`; the later file-local frontmatter contract passed local client parsing and user-skill discovery probes. A pre-existing SessionEnd hook reported that `node` was unavailable after the earlier answer.
- OpenCode isolated runtime: `debug skill` discovered all 23 skills. Model runs timed out or returned no content, and an unsandboxed retry using the real provider data directory was not approved, so routing remains runtime-unverified.
- Gemini: CLI unavailable; runtime-unverified.
- Current Codex catalog exposure includes `design-art-direction`; `mcp-builder-go` is intentionally disabled by local Codex configuration and is not a missing catalog entry.

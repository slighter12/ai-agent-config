---
name: mcp-builder-go
description: Guide for creating high-quality MCP (Model Context Protocol) servers in Go. Use when building or extending MCP servers that integrate external APIs/services with well-designed tools, resources, and prompts. Avoid when the task is not an MCP server, is not Go-specific, or only needs ordinary API/client implementation guidance.
metadata:
  version: "0.1.1"
---

# MCP Server Development Guide (Go)

## Overview

Build MCP servers that help LLMs complete real tasks reliably.
Prioritize protocol correctness, discoverable tool design, stable outputs, and testability.

Use this skill when:

- Creating a new MCP server in Go.
- Migrating Python/TypeScript MCP servers to Go.
- Expanding an existing Go MCP server with tools/resources/prompts.
- Hardening protocol behavior (initialize, capabilities, errors, pagination, prompts/resources/tools methods).

Avoid this skill when:

- The task is not about MCP server design or implementation.
- The implementation is not Go-specific.
- The request only needs ordinary API client, service integration, or application architecture guidance.
- The question is whether a workflow should become a skill, reference, script, or helper; use `skill-creator` for that classification.

## Default Stance

- Prefer official MCP protocol behavior over ad-hoc client compatibility hacks.
- Prefer comprehensive API coverage first; add workflow tools where they clearly reduce multi-step friction.
- Prefer deterministic behavior (stable ordering, stable IDs, stable error kinds/messages where practical).
- Keep responses concise and structured for machine consumption.
- For agent-facing CLI/MCP tools, prefer stable, compact, bounded, machine-readable outputs.

---

## Process

### Phase 1: Protocol and API Research

#### 1.1 Read MCP docs first

Before coding, review:

- MCP lifecycle and capability negotiation.
- Tools/resources/prompts behavior and method contracts.
- Transport expectations (stdio vs streamable HTTP).

Start from:

- `https://modelcontextprotocol.io/sitemap.xml`
- `https://modelcontextprotocol.io/specification/2025-11-25`

Check `.md` pages when you need plain text snapshots.

#### 1.2 Study the target service API deeply

For the external API you are integrating:

- Authentication flows and token refresh behavior.
- Rate limits and pagination.
- Idempotency and destructive operations.
- Error models and partial failure cases.

Map API operations to MCP surfaces:

- Tool candidates (actions).
- Resource candidates (readable artifacts/state).
- Prompt candidates (reusable LLM task templates).

#### 1.3 Decide server scope

Define explicitly:

- Must-have operations (v1).
- Read-only vs mutating operations.
- Which MCP capabilities are enabled now vs planned later.

Do not advertise a capability that is runtime-disabled.

---

### Phase 2: Server and Tool Design

#### 2.1 Tool naming and discoverability

Use action-oriented names with a service prefix:

- `github_list_repos`
- `github_create_issue`

Naming rules:

- Verb-first action semantics.
- Predictable prefixes for grouping.
- Avoid aliases unless required for compatibility.

#### 2.2 Input/output schema strategy

For each tool:

- Define strict input shape.
- Validate required fields and ranges.
- Return structured JSON-compatible output with stable field names.

In Go:

- Use typed structs for request/response internals.
- Keep transport payloads as `map[string]any` only at boundaries.
- Keep schema descriptions concrete and task-oriented.

#### 2.3 Error semantics

Return actionable errors:

- What failed.
- Which field/operation caused it.
- What the caller should do next.

Separate:

- Invalid input (`invalid_params` style conditions).
- Unsupported feature (`not_supported` style conditions).
- Temporary upstream failures (`not_available`/retryable conditions).

#### 2.4 Pagination and large payloads

Design for bounded responses:

- Cursor or offset strategy with clear contracts.
- Stable sort order.
- `nextCursor` only when more data exists.

---

### Phase 3: Go Implementation

#### 3.1 SDK choice (recommended)

Preferred:

- Official Go SDK: `github.com/modelcontextprotocol/go-sdk`

Alternative ecosystems (only when justified):

- `github.com/mark3labs/mcp-go`
- Existing in-house transport/router layers

Pick one primary stack per server to avoid fragmented behavior.

#### 3.2 Project structure baseline

Suggested layout:

- `cmd/` or `main.go`: wiring and startup.
- `transport/`: stdio/HTTP transport boundaries.
- `mcp/` or `server/`: MCP method dispatch.
- `service/` or `client/`: external API clients.
- `tools/`, `resources/`, `prompts/`: feature implementations.
- `internal/` helpers for validation, pagination, errors.

#### 3.3 Capability declaration discipline

On `initialize` response:

- Include only capabilities that are truly available at runtime.
- If prompts are disabled, do not advertise `capabilities.prompts`.
- Keep behavior consistent between stdio and HTTP transports.

#### 3.4 Prompt/resource/tool behavior consistency

If a capability is enabled:

- `*/list` should be deterministic and paginatable where relevant.
- `prompts/get` should validate prompt name and argument payloads.
- Prompt names should be unique (case-insensitive uniqueness is recommended).

#### 3.5 Security and robustness

- Never log secrets/tokens.
- Bound request body size for HTTP transports.
- Validate origin/session/protocol headers when using streamable HTTP.
- Make timeout/retry behavior explicit in upstream API clients.

---

### Phase 4: Verification and Testing

#### 4.1 Required test coverage

Add unit tests for:

- Capability negotiation behavior.
- `tools/list`, `resources/list`, `prompts/list` pagination.
- `tools/call`/`prompts/get` invalid params handling.
- Duplicate prompt/tool naming collisions.
- Error-kind/data mapping consistency.

#### 4.2 Manual verification checklist

- Initialize with each supported transport.
- Confirm advertised capabilities match runtime behavior.
- Call representative read-only and mutating tools.
- Validate pagination edges: first page, middle page, last page, invalid cursor.
- Validate disabled-feature behavior returns clear errors.

#### 4.3 MCP Inspector smoke test

Use MCP Inspector to validate request/response behavior and discoverability.

---

### Phase 5: Evaluation Artifacts

After implementation is stable:

- Create at least 10 realistic, independent evaluation tasks.
- Prefer read-only tasks unless mutation testing is intentional.
- Ensure answers are verifiable and stable over time.

Store evaluation artifacts in a dedicated folder (for example `eval/`).

---

## Quality Checklist

Before finishing, confirm:

- Capabilities advertised == capabilities actually supported.
- Tool names are discoverable and consistent.
- Error messages are actionable and non-ambiguous.
- Outputs are structured, stable, and not overly verbose.
- Pagination behavior is consistent across list endpoints.
- Unit tests cover protocol edge cases and regressions.

## Common Pitfalls

- Advertising prompts/resources capabilities when runtime is disabled.
- Returning non-deterministic ordering from list endpoints.
- Mixing display prose and structured payloads inconsistently.
- Treating case-variant names as different logical prompts/tools.
- Building workflow-only tools without baseline API coverage.

## Version History

- v0.1.0 (2026-05-11): Initial portable MCP builder skill with explicit routing boundaries.
- v0.1.1 (2026-05-13): Add agent-native tool surface guidance for CLI/MCP design.

## References

Local references:

- `references/INDEX.md` - Use for navigation and file selection.
- `references/AGENT_NATIVE_TOOL_SURFACE.md` - Agent-facing CLI/MCP output, command, cache, and shipcheck guidance.

Core MCP docs:

- <https://modelcontextprotocol.io/>
- <https://modelcontextprotocol.io/specification/2025-11-25>
- <https://modelcontextprotocol.io/sitemap.xml>

Go SDK and ecosystem:

- <https://github.com/modelcontextprotocol/go-sdk>
- <https://github.com/mark3labs/mcp-go>

Original inspiration:

- <https://skillsmp.com/skills/anthropics-skills-skills-mcp-builder-skill-md>
- <https://skills.sh/anthropics/skills/mcp-builder>

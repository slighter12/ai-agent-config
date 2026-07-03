# Agent-Native Tool Surface

Use this guidance when building an MCP server, companion CLI, or shared client layer that agents will call repeatedly.

## Principle

Agent-facing tools should be predictable, compact, structured, and verifiable. A raw human-oriented CLI or thin API wrapper is often not enough.

## Output Contract

Prefer:

- Stable JSON for machine consumption.
- Compact text for human summaries.
- Bounded results with pagination, cursors, filters, or limits.
- Stable field names and ordering.
- Typed errors with actionable next steps.
- Explicit dry-run output for mutating operations.

Avoid:

- Long prose blobs as primary output.
- Non-deterministic ordering.
- Hidden truncation.
- Mixing display text and structured payloads inconsistently.
- Returning secrets, tokens, cookies, or raw credentials.

## Command And Tool Design

Useful surfaces may include:

- `list`: bounded listing with filters and pagination.
- `get`: one object by stable ID.
- `search`: narrow query with stable scoring or ordering.
- `sync`: explicit local cache refresh when a cache is justified.
- `diff`: compare local and remote or before and after state.
- `doctor` or `health`: environment, auth, config, and connectivity checks.
- `dry-run`: preview mutating changes before execution.

For MCP tools, keep tool names action-oriented and service-prefixed. For CLIs, expose flags such as `--json`, `--compact`, and `--dry-run` when they materially improve agent use.

## CLI And MCP Shared Layer

When both CLI and MCP are needed, share the service/client layer:

- Authentication and config loading.
- API pagination and retries.
- Validation and error mapping.
- Output shaping and redaction.
- Cache access, if any.

Do not maintain separate behavior for CLI and MCP unless a transport boundary truly requires it.

## Local Cache Or SQLite

Use local cache or SQLite only when it solves a real agent problem:

- Expensive or rate-limited remote reads.
- Repeated search or filtering across many records.
- Offline or stale-but-useful inspection.
- Diffing snapshots.
- Compound commands that need local joins.

Avoid cache by default when the source is small, cheap, volatile, or security-sensitive. Document staleness, refresh commands, and invalidation behavior.

## Secret And HAR Handling

- Redact tokens, cookies, API keys, session IDs, and authorization headers.
- Treat HAR files and captured traffic as sensitive.
- Do not store raw credentials in generated artifacts, logs, snapshots, fixtures, or examples.
- Make credential and environment requirements explicit.

## Shipcheck

Before treating a tool as agent-ready, verify:

- Representative read-only commands.
- Representative mutating commands in dry-run mode.
- JSON output parses and stays stable.
- Compact output is actually compact.
- Error paths return clear messages and non-zero status where appropriate.
- `doctor` or equivalent catches missing auth/config when available.
- Secrets are redacted in normal and error output.

For verification-gate checklist detail, see `policy-testing/references/TOOL_SHIPCHECK.md`.

## Source Inspiration

Inspired by agent-native CLI/MCP design patterns in Printing Press: <https://printingpress.dev/>.
This reference contains distilled guidance only. It does not copy, vendor, or require Printing Press runtime, skill text, or generated artifacts.

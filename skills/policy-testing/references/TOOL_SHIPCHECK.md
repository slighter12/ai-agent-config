# Tool Shipcheck

Use this reference for CLI, MCP, script-backed skills, and agent-facing helper tools.

Build success alone does not prove an agent-facing tool is usable.

## Required Evidence

For relevant tools, verify or plan verification for:

- Representative read-only command or tool call.
- Representative mutating path in dry-run mode when available.
- JSON output parsing and stable field names.
- Compact output for low-token agent use.
- Error path behavior and exit status.
- Missing auth/config handling.
- Secret redaction in normal and error output.

## CLI Checks

Confirm:

- `--help` or equivalent names real commands and flags.
- `--json` output is valid JSON when promised.
- `--compact` output is concise and still useful.
- `--dry-run` does not mutate state.
- Non-zero exit codes distinguish failures from empty results.

## MCP Checks

Confirm:

- `initialize` advertises only available capabilities.
- `tools/list`, `resources/list`, and `prompts/list` are deterministic.
- Tool schemas validate required inputs and ranges.
- Tool outputs are bounded and structured.
- Disabled or missing upstream features fail clearly.

## Boundaries

- Do not apply this checklist to ordinary documentation-only or policy-only skills.
- Do not run live mutating checks without explicit approval.
- Do not store tokens, cookies, HAR secrets, or API keys in fixtures or logs.

## Source Inspiration

Inspired by agent-native CLI verification patterns in Printing Press: <https://printingpress.dev/>.
This reference contains distilled guidance only. It does not copy, vendor, or require Printing Press runtime, skill text, or generated artifacts.

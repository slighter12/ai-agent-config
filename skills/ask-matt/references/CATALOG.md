# Skill Catalog

User-invoked skills load only after explicit invocation. Model-invoked skills may also load automatically when the request matches. `Source` distinguishes the upstream promoted catalog from local extensions.

## User-invoked

| Skill | Purpose | Source |
| --- | --- | --- |
| `ask-matt` | Choose the smallest useful skill or flow. | upstream |
| `grill-with-docs` | Interview while maintaining domain docs. | upstream |
| `triage` | Move issues and external PRs through triage roles. | upstream |
| `improve-codebase-architecture` | Survey active code for deepening opportunities. | upstream |
| `setup-matt-pocock-skills` | Configure tracker, triage labels, and domain docs. | upstream |
| `to-spec` | Publish settled decisions as a spec plus a tracer-bullet ticket graph. | upstream |
| `implement` | Build a spec or ticket with TDD and review. | upstream |
| `wayfinder` | Plan a multi-session effort as a shared decision map. | upstream |
| `grill-me` | Run a stateless requirements interview. | upstream |
| `handoff` | Create a portable cross-context handoff. | upstream |
| `teach` | Teach a concept in a stateful workspace. | upstream |
| `to-questionnaire` | Gather decisions from another person asynchronously. | upstream |
| `wait-what` | Re-pitch the last message with missing context. | upstream |
| `retro` | Graduate session friction into checks, standards, and project rules. | upstream |
| `mcp-builder-go` | Design a current-protocol MCP server in Go with project-fit dependencies. | local extension |
| `how` | Explain how a code path works, with `file:line` references. | local extension |
| `why` | Trace why code has its current shape from repository history, with tiered confidence. | local extension |
| `blast-radius` | Find what a change breaks beyond its diff and prove the fact it is safe because of. | local extension |
| `verify-app` | Create or maintain a project-local skill that proves app behavior end to end. | local extension |

## Model-invoked

| Skill | Purpose | Source |
| --- | --- | --- |
| `prototype` | Answer a logic or UI design question with throwaway code. | upstream |
| `diagnosing-bugs` | Diagnose hard bugs through a tight evidence loop. | upstream |
| `research` | Delegate primary-source reading by default, with complete inline fallback. | upstream |
| `tdd` | Drive red-green vertical slices at agreed seams. | upstream |
| `domain-modeling` | Sharpen domain language and durable decisions. | upstream |
| `codebase-design` | Design deep modules and clean seams. | upstream |
| `code-review` | Review a diff on Standards and Spec axes. | upstream |
| `wizard` | Generate an interactive Bash workflow for human-only steps. | upstream |
| `grilling` | Interview by design-tree frontier rounds. | upstream |
| `writing-for-agents` | Write predictable documents agents consume. | upstream |
| `pr` | Format a PR body with summary, evidence, and merge danger. | upstream |

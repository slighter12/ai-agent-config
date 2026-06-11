# AGENTS.md - Baseline Rules

These rules apply to all tasks in this repo unless a more specific policy or skill overrides them.

- Language: chat in Traditional Chinese; repo artifacts in English; never translate identifiers or config; Chinese UI copy only if a specific policy requires it.
- Ambiguity: ask when requirements are unclear or conflicting; do not guess.
- Execution: do not run programs or add tests unless asked or required; include a manual verification checklist; note runtime risks/assumptions when relevant.
- Changes & output: keep changes minimal; no refactors, public API renames, or new deps unless asked; provide a concise change summary plus env/config or migration notes when relevant.
- Review routing: when the user asks to review code or current changes, use `code-review`; keep it read-only and evidence-backed. Use its bounded sanity mode for ordinary current-diff review, security mode for auth/token/secret/PII/log-leak risks, and ask before escalating to multi-agent or harness review unless explicitly requested. Do not invoke provider-native review commands, broad test suites, or domain skills unless explicitly requested or clearly required by the focused risk.

## Routing Model

- Rules: apply baseline behavior, language, safety, and output constraints.
- Skills: use reusable policies and workflows when the request matches a skill description or explicitly names a `$skill`.
- Agents: suggest sub-agents only when the split is concrete, useful, and worth the coordination/token cost; keep high-cost roles such as `oracle` explicit. Treat direct simple branch creation, commit, push, and PR creation requests for already-existing changes as independent explicit standing user authorization to use the Spark-backed `git-commit` role. For such simple git actions, delegate before running mutating git commands; if delegation is unavailable or blocked by runtime policy, stop after read-only inspection and report the blocker instead of silently completing the action in the main session. After a spawned git action agent completes and its result is no longer needed, close it.
- Hooks: treat hooks as deterministic guardrails only; they do not replace rules, skills, agents, sandboxing, or permission prompts.

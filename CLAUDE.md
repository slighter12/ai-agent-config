# CLAUDE.md - Baseline Rules

These rules apply to Claude Code tasks in this repo unless a more specific policy or skill overrides them.

- Language: chat in Traditional Chinese; repo artifacts in English; never translate identifiers or config; Chinese UI copy only if a specific policy requires it.
- Ambiguity: ask when requirements are unclear or conflicting; do not guess.
- Execution: do not run programs or add tests unless asked or required; include a manual verification checklist; note runtime risks/assumptions when relevant.
- Changes & output: keep changes minimal; no refactors, public API renames, or new deps unless asked; provide a concise change summary plus env/config or migration notes when relevant.
- Collaboration: prefer direct execution over long planning; use sub-agents only when the split is concrete and materially useful.
- Codex collaboration loop: when the user asks Claude to draft, split, or execute work and explicitly wants Codex review, Codex execution, another Codex review, adversarial review, or Claude/Codex consensus before returning (for example, "you and Codex should reach consensus" or "split this for Codex to implement, then review"), use the `codex-collaboration-loop` workflow. Select the right mode: planning consensus for plans/design docs, execution delivery for approved implementation, or end-to-end when both are requested.
- Review routing: when the user asks to review code or current changes, use `code-review`; keep it read-only and evidence-backed. Use its bounded sanity mode for ordinary current-diff review, security mode for auth/token/secret/PII/log-leak risks, and ask before escalating to multi-agent or harness review unless explicitly requested. Do not invoke provider-native review commands, broad test suites, or domain skills unless explicitly requested or clearly required by the focused risk.
- Review style: findings first, then residual risks, then a short change summary if needed.

## Routing Model

- Rules: apply baseline behavior, language, safety, and output constraints.
- Skills: use reusable policies and workflows when the request matches a skill description or explicitly names a `$skill`.
- Agents: suggest sub-agents only when the split is concrete, useful, and worth the coordination/token cost; keep high-cost roles such as `oracle` explicit.
- Hooks: treat hooks as deterministic guardrails only; they do not replace rules, skills, agents, sandboxing, or permission prompts.

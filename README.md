# AI Agent Configuration (`slighter12/ai-agent-config`)

A portable, language-agnostic configuration suite of **Skills**, **Always-On Steering (`AGENTS.md` / `CLAUDE.md`)**, and **Lightweight Multi-Agent Roles** for Claude Code, Codex CLI, and OpenCode.

## Core Architectural Principles

1. **Single-Skill Self-Sufficiency (Respecting the 1–2 Skill Activation Budget)**:
   * Coding harnesses typically auto-activate only 1–2 skills per prompt based on description similarity. Every action skill (`grilling`, `to-spec`, `codebase-design`, `tdd`, `implement`, `code-review`, `pr`, `retro`) is self-contained and embeds a compact **Pre-Flight Context & Agency Gate**.
2. **4-Tier Rule Precedence & 6-Scenario User Agency Gate**:
   * Precedence: `Explicit User Prompt > Project Law (AGENTS.md / CLAUDE.md / docs/adr/ / GLOSSARY.md) > Local File Gravity > Global Skill Defaults`.
   * When project ADRs conflict with legacy code in a target file, when multiple competing styles exist in a multi-author repo, or when a legacy monolith lacks clean seams, the agent **pauses and presents concrete options to the user** (`Follow Local Style`, `Sprout at the Seam`, `Migrate Module`, or `Runtime Verification`) rather than guessing.
3. **Task-Proportional Verification & Mechanical Heuristics**:
   * Replaces language-coupled TDD dogma with observable syntactic redlines (`Single-Caller Helper Ban`, `No Call-Only Mock Assertions`, `Evidence Conservation Law`, and `Post-Green Same-Branch Test Pruning`).
4. **Single-Window Default, Escalate Only at Smart-Zone Boundaries**:
   * Works with zero custom agent overhead by default (Level 0), offloading noisy read-only searches to `explorer` / `librarian` / `oracle` (Level 1) and escalating to parallel git worktrees (Level 2) only for multi-ticket specs with non-overlapping write scopes.

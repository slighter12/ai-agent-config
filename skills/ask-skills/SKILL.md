---
name: ask-skills
description: Explain this skill catalog, invocation boundaries, and the best skill for a concrete request. Use when the user explicitly asks how the installed skills work or which one to invoke. Avoid when the user wants the underlying task performed.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Ask Skills

Load [the canonical catalog and route map](references/CATALOG.md). Answer from that catalog and the installed provider metadata. Explain the smallest useful route, why it matches, and any explicit invocation or approval boundary.

Return catalog guidance only. The underlying workflow begins under its own user request and invocation boundary.

---
name: setup-skills
description: Install or refresh this shared skill catalog and initialize its project coordination files. Use when the user explicitly asks to set up these skills or their local tracker and domain-model layout. Avoid when the user only wants an explanation of the catalog.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Setup Skills

Inspect the target project and provider homes before changing them. Use this repository's installer for shared provider links and invocation controls.

For project coordination, create only what is missing:

- `docs/agents/issue-tracker.md` for tracker choice, status, and blocking edges.
- `docs/agents/domain/` for domain context and ADRs.

Prefer a configured external tracker and its CLI. If none exists, initialize `docs/agents/issue-tracker.md` with `.scratch/<feature-slug>/` as a path pattern, `spec.md` at the resolved feature root, and numbered ticket files under `issues/`. Define slug resolution, reuse, ambiguity, and collision behavior; the placeholder is never a literal directory. Treat a missing or empty file as unconfigured, preserve every substantive existing configuration, and omit invented triage labels.

Never overwrite user-owned files or links from another repository. Report provider gaps and unverified Gemini behavior.

---
name: policy-frontend
description: "Apply frontend framework and product UI rules across React, Svelte, and Solid: components, state, routing, i18n, accessibility, interaction controls, responsive layout, and framework-specific patterns. Use when React/Svelte/Solid implementation guidance or UI behavior decisions are needed. Avoid when pure visual art direction, backend/API contracts, infrastructure, security, or testing strategy are the primary concern, and do not treat Vue as first-class scope in this skill."
metadata:
  version: "0.3.1"
---

# Policy Guide

Use the reference files for full rules, examples, and patterns.

Use optional related policy detail only when frontend work crosses boundaries and those risks are in scope:

- `design-art-direction`: Visual polish, composition, typography, color, imagery, motion, or `.design/` reference packs.
- `policy-api`: API contract, request/response schema, status/error semantics, or versioning changes.
- `policy-security`: Auth/authz, trust boundaries, token/session handling, or sensitive data exposure.
- `policy-testing`: New or changed frontend behavior that needs unit/integration/e2e coverage.
- `policy-infra`: Runtime/deploy config impact (env vars, build pipeline, CDN/edge, container/runtime settings).

## Product UI Baseline

- Build the actual usable workflow as the first screen for apps, tools, dashboards, and games; do not default to a marketing landing page.
- Prefer existing design-system components, interaction primitives, and icon libraries. Use `lucide` icons when the project already has them available.
- Keep loading, empty, error, and disabled states consistent with the data and API semantics owned by `policy-api`.
- Visual detail, composition, imagery, motion, cards, gradients, text-fit, and stable-dimension rules are owned by `design-art-direction`.

## Accessibility Baseline

- Treat accessibility fixes by user impact and affected workflow, not by tool count alone.
- Cover keyboard navigation, focus order, screen reader semantics, contrast, and responsive text fit for changed UI.
- Use WCAG scope as a baseline when auditing, but keep implementation guidance tied to the project's framework and existing components.

## Version History

- v0.1.0 (2026-05-08): Initial portable policy release defining frontend first-class scope for React/Svelte/Solid and cross-policy boundary guidance.
- v0.1.1 (2026-05-29): Align routing wording with optional-depth policy handoffs.
- v0.2.0 (2026-05-31): Add product UI baseline and route pure visual art direction to `design-art-direction`.
- v0.2.1 (2026-06-29): Tighten product UI defaults against generic cards and decorative gradients.
- v0.2.2 (2026-06-30): Add concise accessibility audit and fix baseline.
- v0.3.0 (2026-07-03): Update Svelte framework references to Svelte 5 runes-mode defaults.
- v0.3.1 (2026-07-03): Make `design-art-direction` the owner for detailed visual/product UI rules.

## References

- `references/INDEX.md` - Use for navigation and file selection.

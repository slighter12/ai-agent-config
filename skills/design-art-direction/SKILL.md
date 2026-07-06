---
name: design-art-direction
description: "Apply page design and visual art direction for polished product UI, landing pages, and design critique: composition, typography, color, imagery, motion, and `.design/` reference packs. Use when aesthetics, visual hierarchy, brand feel, page style, or rejected UI design quality are primary. Avoid when frontend implementation mechanics, API contracts, security, or testing strategy are primary."
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.3.2"
---

# Design Art Direction

## Purpose

Guide page design, style direction, and visual judgment for UI and web work while keeping implementation decisions owned by frontend, API, security, and testing policies. This skill owns detailed visual/product UI rules.

Use this skill to improve the look, hierarchy, and interaction feel of an experience, especially when the user says a design looks weak, generic, ugly, too template-like, or not aligned with a product direction.

## Use When

- The task is primarily visual quality, art direction, layout composition, typography, color, imagery, animation, or brand feel.
- The user rejects a design or asks for a redesign, visual critique, or better product presentation.
- A project includes a `.design/` directory that should guide visual direction.
- Frontend implementation exists, but the main decision is what the UI should look and feel like.
- A Codex workflow wants optional specialist critique through a configured design agent.

## Avoid When

- Framework mechanics, state, routing, data loading, accessibility implementation, or component architecture are primary; use `policy-frontend`.
- API contracts, request/response shape, status codes, validation, or error semantics are primary; use `policy-api`.
- Security, auth, secrets, PII, or trust boundaries are primary.
- Backend, infrastructure, broad refactors, or model/provider setup are primary.

## Workflow

1. Inspect the target screen, component, screenshot, local URL, or design goal before proposing visual changes.
2. If `.design/` exists, read the smallest relevant files first and treat them as the art-direction reference for palette, typography, spacing, imagery, and component feel.
3. Run a UX-first pass before visual polish: identify the primary user, primary task, entry point, next step, key states (empty/loading/error/success), recovery path, friction, and information hierarchy.
4. Classify the experience type: operational app, SaaS dashboard, editor/tool, game, landing page, portfolio, venue, product page, or content site.
5. Set the visual direction from the experience type:
   - Operational tools should be quiet, dense, scannable, and built for repeated action.
   - Landing/product pages should make the product, place, person, or object visible in the first viewport.
   - Games and playful tools may be more expressive, animated, and illustrative.
6. Prefer real or generated bitmap imagery when visual assets matter. Do not substitute decorative SVGs, gradients, or atmospheric blobs for the actual product or state the user needs to inspect.
7. Check composition, hierarchy, contrast, type scale, spacing rhythm, color balance, empty states, loading states, and responsive behavior.
8. In Codex, when the user rejects a design or asks for a stronger visual pass, the coordinator may route advisory critique or explicitly bounded low-risk edits to `page-designer`. Keep final judgment with the primary agent.

## Design Rules

- Build the usable experience as the first screen unless the user explicitly asks for a marketing-only landing page.
- Do not make a screen prettier before the task flow is understandable, recoverable, and easy to continue.
- Design recommendations must name concrete UX problems or decisions; avoid generic "make it modern" or "add polish" guidance.
- Avoid generic text-card layouts, card-heavy page sections, nested cards, oversized generic hero sections, decorative gradients, gradient-orb backgrounds, and one-note color palettes.
- Use cards only for repeated items, modals, tool surfaces, or genuinely framed UI regions; do not turn page sections into floating card walls.
- Use real product state, functional panels, lists, tables, task flows, or actual/generative imagery instead of abstract gradient backgrounds.
- Use familiar controls: icons for tools, swatches for color, segmented controls for modes, toggles for binary settings, sliders or inputs for numeric values, menus for option sets, and tabs for views.
- Keep cards at modest radius, typically 8px or less, unless the existing design system says otherwise.
- Do not scale font size with viewport width. Keep letter spacing at 0 unless a brand system requires otherwise.
- Ensure text fits its container on mobile and desktop; no clipped labels, overlapping headings, or button text that crowds icons.
- Define stable dimensions for boards, grids, toolbars, counters, and tiles so hover states and dynamic labels do not shift layout.
- Use motion to clarify state change, spatial movement, or feedback; avoid decorative animation that distracts from the workflow.

## Tool And Side-Effect Boundaries

- Do not install design tools, add dependencies, or run browser automation unless the user asks or the active implementation workflow requires visual verification.
- Do not let visual polish introduce hidden behavior changes, API changes, auth/security changes, unrelated formatting churn, or accessibility regressions.
- Treat `.design/` references as guidance, not permission to overwrite unrelated project conventions.
- Do not install, start, stop, or configure external model proxies from this skill.

## Output

Return:

- `summary`: visual direction or design changes recommended.
- `design_context`: `.design/` files, screenshots, URLs, or components inspected.
- `recommendations`: accepted visual hierarchy, layout, palette, typography, imagery, and motion decisions.
- `handoff`: whether `policy-frontend`, `policy-api`, or a Codex-only `page-designer` advisory pass should own follow-up work.
- `manual_verification`: responsive and visual checks needed after implementation.

## Version History

- v0.1.0 (2026-05-31): Initial art-direction policy split from frontend implementation policy, with `.design/` reference-pack guidance.
- v0.2.0 (2026-06-29): Replace legacy external design handoff wording and tighten anti-card/anti-gradient defaults.
- v0.3.0 (2026-06-29): Merge design bridge guidance into this provider-neutral page design skill and use a generic Codex design agent handoff.
- v0.3.1 (2026-07-03): Clarify ownership of detailed visual/product UI rules shared with frontend work.
- v0.3.2 (2026-07-06): Add UX-first flow, state, recovery, and friction checks before visual polish.

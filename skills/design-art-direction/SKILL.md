---
name: design-art-direction
description: "Apply visual art direction for polished product UI, landing pages, and design critique: composition, typography, color, imagery, motion, and `.design/` reference packs. Use when aesthetics, visual hierarchy, brand feel, or rejected UI design quality are primary. Avoid when frontend implementation mechanics, API contracts, security, or testing strategy are primary."
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.0"
---

# Design Art Direction

## Purpose

Guide visual judgment for UI and web work while keeping implementation decisions owned by frontend, API, security, and testing policies.

Use this skill to improve the look, hierarchy, and interaction feel of an experience, especially when the user says a design looks weak, generic, ugly, too template-like, or not aligned with a product direction.

## Use When

- The task is primarily visual quality, art direction, layout composition, typography, color, imagery, animation, or brand feel.
- The user rejects a design or asks for a redesign, visual critique, or better product presentation.
- A project includes a `.design/` directory that should guide visual direction.
- Frontend implementation exists, but the main decision is what the UI should look and feel like.

## Avoid When

- Framework mechanics, state, routing, data loading, accessibility implementation, or component architecture are primary; use `policy-frontend`.
- API contracts, request/response shape, status codes, validation, or error semantics are primary; use `policy-api`.
- Security, auth, secrets, PII, or trust boundaries are primary.
- Antigravity CLI is required for external critique or bounded edits; use `antigravity-design-bridge` with this skill as design policy context.

## Workflow

1. Inspect the target screen, component, screenshot, local URL, or design goal before proposing visual changes.
2. If `.design/` exists, read the smallest relevant files first and treat them as the art-direction reference for palette, typography, spacing, imagery, and component feel.
3. Classify the experience type: operational app, SaaS dashboard, editor/tool, game, landing page, portfolio, venue, product page, or content site.
4. Set the visual direction from the experience type:
   - Operational tools should be quiet, dense, scannable, and built for repeated action.
   - Landing/product pages should make the product, place, person, or object visible in the first viewport.
   - Games and playful tools may be more expressive, animated, and illustrative.
5. Prefer real or generated bitmap imagery when visual assets matter. Do not substitute decorative SVGs, gradients, or atmospheric blobs for the actual product or state the user needs to inspect.
6. Check composition, hierarchy, contrast, type scale, spacing rhythm, color balance, empty states, loading states, and responsive behavior.
7. When the user rejects a design or asks for a stronger redesign, consider routing an advisory or bounded direct-edit pass through `antigravity-design-bridge`.

## Design Rules

- Build the usable experience as the first screen unless the user explicitly asks for a marketing-only landing page.
- Avoid nested cards, card-heavy layouts, oversized generic hero sections, gradient-orb backgrounds, and one-note color palettes.
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

## Output

Return:

- `summary`: visual direction or design changes recommended.
- `design_context`: `.design/` files, screenshots, URLs, or components inspected.
- `recommendations`: accepted visual hierarchy, layout, palette, typography, imagery, and motion decisions.
- `handoff`: whether `policy-frontend`, `policy-api`, or `antigravity-design-bridge` should own follow-up implementation.
- `manual_verification`: responsive and visual checks needed after implementation.

## Version History

- v0.1.0 (2026-05-31): Initial art-direction policy split from frontend implementation policy, with `.design/` reference-pack guidance.

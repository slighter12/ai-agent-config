---
globs: ["*.tsx", "*.jsx", "*.svelte", "*.html", "**/*.tsx", "**/*.jsx", "**/*.svelte", "**/*.html"]
description: "Frontend accessibility - WCAG compliance, keyboard navigation, screen readers"
---

# FRONTEND_ACCESSIBILITY.md - Frontend Accessibility Overview

This file provides an a11y overview and routing guidance.
Select the appropriate sub-policy based on your needs.

---

## Policy Routing

- Semantics and structure: see `FRONTEND_A11Y_SEMANTICS.md`
- Interaction and focus: see `FRONTEND_A11Y_INTERACTION.md`
- Visual and responsive: see `FRONTEND_A11Y_VISUAL.md`

---

## Testing Tools

Use these tools during development:

- axe DevTools - automated a11y tests
- WAVE - visualized a11y issues
- Keyboard navigation - manual testing (unplug the mouse)
- Screen reader - VoiceOver (macOS/iOS), NVDA (Windows)

---

## Framework-Specific Guidance

For framework-specific examples, see the React/Svelte/Solid references in `../INDEX.md`.

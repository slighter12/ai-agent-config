# FRONTEND_A11Y_VISUAL.md - Frontend Accessibility: Visual and Responsive

This file focuses on visual and responsive accessibility rules.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Color Contrast (WCAG AA)

- Text contrast at least 4.5:1 for normal text
- Large text (18pt+) contrast at least 3:1
- Do not rely on color alone to convey information (add icons or text)

### 2) Responsive and Zoom

- Support text zoom up to 200% without layout breakage
- Do not disable zoom (`user-scalable=no`)
- Touch targets at least 44x44px (iOS) or 48x48px (Android)

### 3) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Whether color contrast meets WCAG
- Whether touch target sizing meets platform guidance

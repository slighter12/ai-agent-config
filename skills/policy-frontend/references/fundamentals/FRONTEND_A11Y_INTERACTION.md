# FRONTEND_A11Y_INTERACTION.md - Frontend Accessibility: Interaction and Focus

This file focuses on interaction and focus accessibility rules.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Keyboard Navigation (mandatory)

- All interactive elements must be keyboard operable (Tab, Enter, Space, Esc)
- Focus state must be visible (`:focus-visible`)
- Custom components must implement correct keyboard handlers

### 2) Form Accessibility

- Every input must have a corresponding `<label>` (use the `for` attribute)
- Error messages must use `aria-invalid` and `aria-describedby`
- Required fields must use the `required` attribute and visual indicators

### 3) Focus Management

- When a modal opens, move focus to the first interactive element inside it
- When a modal closes, return focus to the trigger element
- Provide skip links to jump past repeated content (e.g., navigation)

### 4) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Whether keyboard navigation order is reasonable
- Whether focus management matches user flows
- Whether form error associations are correct

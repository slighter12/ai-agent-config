---
description: "Frontend component design - composition, props, reusability"
---

# FRONTEND_COMPONENT.md - Frontend Component Design

This file defines frontend component design rules (applies to SolidJS, Svelte, React).
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Single Responsibility

- Each component does one thing
- Component files should not exceed 250 lines (including styles)
- If larger, split into subcomponents or framework-specific abstractions

### 2) Prop Naming and Types

- Props must have explicit type definitions (TypeScript)
- Boolean props start with `is`, `has`, `should`
- Event handler props start with `on` (e.g., `onClick`, `onSubmit`)

### 3) Avoid Prop Drilling

- More than 2 levels deep -> use Context/Store
- Global state should use a single state management approach (see FRONTEND_STATE.md)
- Do not pass callbacks through many layers

### 4) Prefer Composition over Inheritance

- Use composition patterns (children, slots)
- Avoid class component inheritance
- Prefer HOCs or composition utilities when needed

### 5) Controlled vs Uncontrolled

- Default to controlled inputs for form components
- Use uncontrolled only for performance or third-party integration
- Use clear naming: `value/onChange` vs `defaultValue`

### 6) Side Effect Management

- Keep side effects in framework lifecycle/effect mechanisms
- Do not run side effects during render/compute (API calls, subscriptions)
- Cleanup functions must remove subscriptions and timers

### 7) Conditional Rendering

- Prefer framework-provided conditional syntax
- Avoid complex nested ternaries (more than one layer -> refactor)
- Consolidate null checks at the top of the component

### 8) List Keys

- Lists must use stable `key` values (avoid index)
- Keys must be unique and stable (ID, UUID)
- Keep keys consistent when sorting or filtering

### 9) When to Optimize

- Do not optimize prematurely - ensure correctness first
- Measure with framework DevTools before optimizing
- Common tools: memoization, lazy loading, virtualization

### 10) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- How to draw component boundaries
- Whether state belongs locally or in Context/Store
- Whether memoization is needed
- How to handle complex form state

---

## Detailed Guidance

Need framework-specific examples? See the React/Svelte/Solid references in `../INDEX.md`.

---

Violating these rules is incorrect output.

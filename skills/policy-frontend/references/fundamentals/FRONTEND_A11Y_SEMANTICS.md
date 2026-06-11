# FRONTEND_A11Y_SEMANTICS.md - Frontend Accessibility: Semantics and Structure

This file focuses on semantic and structural accessibility rules.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Semantic HTML (mandatory)

- Use correct HTML5 semantic tags (`<nav>`, `<main>`, `<article>`, `<section>`)
- Avoid div soup (divitis)
- Use `<button>` for buttons and `<a>` for links (do not swap)

### 2) Heading Levels

- Only one `<h1>` per page
- Headings must follow order (h1 -> h2 -> h3, do not skip levels)
- Do not use headings purely for styling (use CSS instead)

### 3) Alt Text (mandatory)

- Every `<img>` must have an `alt` attribute
- Decorative images use `alt=""`
- Informational images provide descriptive text in Traditional Chinese

### 4) ARIA Attributes

- Prefer semantic HTML; ARIA is a supplement
- Use `aria-label`, `aria-labelledby`, `aria-describedby` when needed
- Use `aria-live` for dynamic content to notify screen readers

### 5) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- How to set ARIA for complex components
- Whether heading levels match the information architecture
- Whether an image needs descriptive alt text

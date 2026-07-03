---
description: "Frontend routing - navigation, data loading, protected routes"
---

# FRONTEND_ROUTING.md - Frontend Routing

This file defines frontend routing rules (applies to SolidJS, Svelte, React).
Vue code patterns may exist in-repo but are out-of-scope by default unless a repo-specific policy explicitly includes Vue.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Route Definitions

- Use framework-recommended routing solutions (SolidStart Router, SvelteKit, Next.js App Router)
- Centralize route definitions (single file or folder structure)
- Avoid magic strings - define route constants

### 2) Route Params and Query Strings

- Define TypeScript types for URL params
- Validate query strings (zod, yup)
- Required params use path params (`/user/:id`)
- Optional params use query strings (`/search?q=...`)

### 3) Protected Routes (Auth)

- Handle protected routes uniformly (middleware or HOC)
- Redirect unauthenticated users to the login page
- After login, redirect back to the original page (`?redirect=/...`)

### 4) Data Loading Strategy

- Prefer framework data loading mechanisms (loaders, load functions)
- Avoid fetching inside components (prevents waterfalls)
- Use Suspense/await for loading states

### 5) Navigation

- Use router APIs for programmatic navigation (`navigate()`, `goto()`)
- Use `<a>` or `<Link>` for clickable navigation (accessibility)
- Avoid `window.location.href = ...` (breaks SPA behavior)

### 6) Route Change Side Effects

- Clean up subscriptions and timers on route change
- Use route guards or lifecycle hooks
- Let the framework handle scroll restoration

### 7) 404 and Error Handling

- Provide a 404 page
- Use error boundaries for route errors
- Error messages must be user-friendly. Use Traditional Chinese only when i18n is enabled or project policy requires it.

### 8) Nested Routes

- Use framework nested route features
- Avoid manual sub-routing inside a single component
- Use layout routes for shared structure (header, sidebar)

### 9) Route Meta and SEO

- Every route must define `<title>` and `<meta description>`
- Use framework head APIs (SolidJS `<Title>`, SvelteKit `<svelte:head>`)
- Dynamic pages (e.g., product pages) must generate meta from data

### 10) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- How to organize the route structure
- Whether nested routes are needed
- Where to handle data loading
- How to handle complex authorization (multi-role)

---

## Detailed Guidance

Need framework-specific examples? See:

- `../frameworks/solid/SOLIDSTART.md` - SolidStart Router
- `../frameworks/svelte/SVELTEKIT.md` - SvelteKit routing
- `../frameworks/react/REACT_NEXTJS_APP_ROUTER.md` - Next.js App Router

---

Violating these rules is incorrect output.

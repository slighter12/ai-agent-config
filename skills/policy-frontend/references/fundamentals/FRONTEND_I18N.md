---
description: "Frontend i18n and copy management - centralized Traditional Chinese with English keys"
---

# FRONTEND_I18N.md - Frontend Copy and i18n

This file defines copy and localization rules for frontend projects.
First-class targets: SolidJS, Svelte, React.
Vue code patterns may exist in-repo but are out-of-scope by default unless a repo-specific policy explicitly includes Vue.

Violating these rules is incorrect output.

Note: Examples use placeholder strings like `TC_*` to indicate Traditional Chinese copy.

---

## Table of Contents

- 1) Language Policy (mandatory)
- 1) Centralized Copy Management (mandatory)
- 1) Key Naming Rules (mandatory)
- 1) Dynamic Copy
- 1) Plurals
- 1) Error Message Mapping
- 1) Long-Form Copy
- 1) Formatting Helpers (recommended)
- 1) Copy Checklist
- 1) Future i18n Support (optional)
- 1) When Uncertain (mandatory)

## 1) Language Policy (mandatory)

### Applicability

This language policy applies only when the project has i18n enabled.
i18n enablement means any explicit multi-locale or localization mechanism is in use (e.g., locale files, translation system, locale switching, Accept-Language handling), regardless of implementation style.
If i18n is not enabled, do not force Traditional Chinese; follow the project's existing language rules.
If unclear, stop and ask.

### User-visible copy must be Traditional Chinese

When i18n is enabled, you must provide `zh-TW` copy/locale. Other locales are optional unless explicitly required.

Includes:

- Button labels
- Navigation titles
- Form labels and placeholders
- Validation errors shown to users
- Toasts and notifications
- Empty states
- Loading messages
- Confirmation dialog text

Excludes (keep English):

- Variable/function/type names
- API field names
- Error codes (machine codes)
- CSS class names
- Test IDs
- Log messages
- Developer comments

---

## 2) Centralized Copy Management (mandatory)

This section applies only when i18n is enabled. If i18n is not enabled, follow the project's existing copy practices.

### Inline copy is forbidden

Do not hardcode copy in JSX/templates:

```tsx
// Incorrect
<button>TC_LOGIN</button>
<p>TC_NO_DATA</p>
<input placeholder="TC_EMAIL_PLACEHOLDER" />
```

Correct: centralize copy

Option A: single file (small projects)

```typescript
// src/copy.ts
export const COPY = {
  LOGIN_BUTTON: 'TC_LOGIN',
  NO_DATA: 'TC_NO_DATA',
  EMAIL_PLACEHOLDER: 'TC_EMAIL_PLACEHOLDER',
} as const
```

Option B: feature-based files (large projects)

```typescript
// src/copy/auth.ts
export const AUTH = {
  LOGIN_BUTTON: 'TC_LOGIN',
  SIGNUP_BUTTON: 'TC_SIGN_UP',
  FORGOT_PASSWORD: 'TC_FORGOT_PASSWORD',
} as const

// src/copy/common.ts
export const COMMON = {
  SAVE: 'TC_SAVE',
  CANCEL: 'TC_CANCEL',
  DELETE: 'TC_DELETE',
} as const
```

Usage:

```tsx
import { COPY } from './copy'

<button>{COPY.LOGIN_BUTTON}</button>
```

---

## 3) Key Naming Rules (mandatory)

Rules:

- Keys must be English
- Use SCREAMING_SNAKE_CASE
- Use descriptive names (avoid `TEXT_1`, `LABEL_A`)
- Group by feature

```typescript
// Correct
export const COPY = {
  // Auth
  AUTH_LOGIN_BUTTON: 'TC_LOGIN',
  AUTH_SIGNUP_BUTTON: 'TC_SIGN_UP',

  // User Profile
  PROFILE_EDIT_BUTTON: 'TC_EDIT_PROFILE',
  PROFILE_SAVE_SUCCESS: 'TC_PROFILE_SAVE_SUCCESS',

  // Common
  COMMON_SAVE: 'TC_SAVE',
  COMMON_CANCEL: 'TC_CANCEL',
}

// Incorrect
export const COPY = {
  button1: 'TC_LOGIN',         // unclear
  saveBtn: 'TC_SAVE',          // camelCase
  eliminar_boton: 'TC_DELETE', // non-English key
}
```

---

## 4) Dynamic Copy

### Use functions for variables

```typescript
// src/copy.ts
export const COPY = {
  WELCOME_USER: (name: string) => `TC_WELCOME_BACK_${name}`,
  ITEMS_COUNT: (count: number) => `TC_TOTAL_${count}_ITEMS`,
  DELETE_CONFIRM: (itemName: string) => `TC_CONFIRM_DELETE_${itemName}`,
} as const
```

Usage:

```tsx
<p>{COPY.WELCOME_USER(user.name)}</p>
<span>{COPY.ITEMS_COUNT(items.length)}</span>
```

---

## 5) Plurals

```typescript
export const COPY = {
  ITEMS_COUNT: (count: number) => {
    if (count === 0) return 'TC_NO_ITEMS'
    if (count === 1) return 'TC_1_ITEM'
    return `TC_${count}_ITEMS`
  },
}
```

---

## 6) Error Message Mapping

### Map backend error codes

The backend returns machine-readable error codes; the frontend maps them to Traditional Chinese:

```typescript
// src/copy/errors.ts
export const ERROR_MESSAGES: Record<string, string> = {
  // Auth
  INVALID_CREDENTIALS: 'TC_INVALID_CREDENTIALS',
  EMAIL_ALREADY_EXISTS: 'TC_EMAIL_ALREADY_EXISTS',
  TOKEN_EXPIRED: 'TC_SESSION_EXPIRED',

  // Validation
  VALIDATION_FAILED: 'TC_VALIDATION_FAILED',
  REQUIRED_FIELD_MISSING: 'TC_REQUIRED_FIELD_MISSING',

  // Default
  UNKNOWN_ERROR: 'TC_UNKNOWN_ERROR',
}

// Usage
function handleError(errorCode: string) {
  const message = ERROR_MESSAGES[errorCode] || ERROR_MESSAGES.UNKNOWN_ERROR
  toast.error(message)
}
```

---

## 7) Long-Form Copy

### Use arrays or template strings

```typescript
export const COPY = {
  PRIVACY_POLICY: `
    TC_PRIVACY_LINE_1
    TC_PRIVACY_LINE_2
    TC_PRIVACY_LINE_3
  `.trim(),

  // Or use arrays
  TERMS_OF_SERVICE: [
    'TC_TERM_1',
    'TC_TERM_2',
    'TC_TERM_3',
  ].join('\n'),
}
```

---

## 8) Formatting Helpers (recommended)

```typescript
// src/utils/format.ts
export const format = {
  date: (date: Date) => {
    return new Intl.DateTimeFormat('zh-TW').format(date)
  },

  currency: (amount: number) => {
    return new Intl.NumberFormat('zh-TW', {
      style: 'currency',
      currency: 'TWD'
    }).format(amount)
  },

  number: (num: number) => {
    return new Intl.NumberFormat('zh-TW').format(num)
  },
}
```

---

## 9) Copy Checklist

Before PR when i18n is enabled:

- [ ] No inline hardcoded copy
- [ ] All user-visible text is Traditional Chinese when project policy requires `zh-TW`
- [ ] Keys are English SCREAMING_SNAKE_CASE
- [ ] Error messages map to Traditional Chinese when project policy requires `zh-TW`
- [ ] Copy is centralized in copy files

---

## 10) Future i18n Support (optional)

If multi-language support is needed:

```typescript
// src/i18n/zh-TW.ts
export const zhTW = {
  AUTH_LOGIN_BUTTON: 'TC_LOGIN',
  AUTH_SIGNUP_BUTTON: 'TC_SIGN_UP',
}

// src/i18n/en-US.ts
export const enUS = {
  AUTH_LOGIN_BUTTON: 'Login',
  AUTH_SIGNUP_BUTTON: 'Sign Up',
}

// src/i18n/index.ts
const translations = { 'zh-TW': zhTW, 'en-US': enUS }
export const t = (key: string) => translations[currentLocale][key]
```

For now, only Traditional Chinese is required. Do not over-engineer.

---

## 11) When Uncertain (mandatory)

If any of the following are unclear, stop and ask:

- Whether text is user-visible
- Where copy should live
- How to handle variable copy
- Which Traditional Chinese message maps to an error code

---

Violating these rules is incorrect output.

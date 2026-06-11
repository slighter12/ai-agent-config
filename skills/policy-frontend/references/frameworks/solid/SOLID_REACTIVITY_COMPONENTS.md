# SolidJS Reactivity and Component Patterns - Detailed Guide

## 1) Core Reactivity Concepts

### Signals, Memos, Effects

SolidJS reactivity is fine-grained, not Virtual DOM-based.

```tsx
import { createSignal, createMemo, createEffect } from "solid-js";

// Signal - mutable state
const [count, setCount] = createSignal(0);

// Memo - derived state (computed)
const doubleCount = createMemo(() => count() * 2);

// Effect - side effects
createEffect(() => {
  console.log("Count changed:", count());
});
```

Key principles:

- Signals are functions (must call them: `count()`)
- Memos recompute only when dependencies change
- Effects track dependencies automatically

---

## 2) Component Patterns

### Function components run once

```tsx
// Correct: component function runs once
function Counter() {
  const [count, setCount] = createSignal(0);

  console.log("Component created"); // runs once

  createEffect(() => {
    console.log("Count:", count()); // runs on changes
  });

  return <button onClick={() => setCount(c => c + 1)}>{count()}</button>;
}
```

How this differs from React:

- React components re-run on every render
- SolidJS components run once; JSX reactive expressions update automatically

### Props handling

```tsx
// Correct: access props in JSX (keeps reactivity)
function Display(props: { value: number }) {
  return <div>{props.value}</div>; // reactive
}

// Incorrect: destructure early (loses reactivity)
function Display(props: { value: number }) {
  const { value } = props; // not reactive
  return <div>{value}</div>;
}

// Correct: use splitProps or mergeProps
import { splitProps } from "solid-js";

function Button(props: { class?: string; onClick: () => void; children: any }) {
  const [local, others] = splitProps(props, ["children"]);
  return <button {...others}>{local.children}</button>;
}
```

---

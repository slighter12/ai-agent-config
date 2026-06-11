# SolidJS Routing, Performance, and TypeScript - Detailed Guide

## Table of Contents

- 1) Routing (SolidStart / Solid Router)
- 1) Performance Optimization
- 1) TypeScript Integration

## 1) Routing (SolidStart / Solid Router)

### SolidStart Router

```tsx
// routes/users/[id].tsx
import { useParams, createAsync } from "@solidjs/router";

export default function UserPage() {
  const params = useParams();
  const user = createAsync(() => fetchUser(params.id));

  return (
    <Show when={user()}>
      {(u) => <div>{u().name}</div>}
    </Show>
  );
}
```

### Data Loaders

```tsx
import { cache, createAsync } from "@solidjs/router";

// Cache loader (avoid repeated fetch)
const getUser = cache(async (id: string) => {
  return fetch(`/api/users/${id}`).then(r => r.json());
}, "user");

export function route({ params }) {
  return getUser(params.id);
}

export default function UserPage() {
  const user = createAsync(() => getUser(params.id));
  // ...
}
```

---

## 2) Performance Optimization

### Avoid unnecessary reactivity

```tsx
// Incorrect: overly reactive
function ExpensiveComponent(props: { data: Data }) {
  const result = createMemo(() => {
    return expensiveCalculation(props.data); // props.data triggers every time
  });
  // ...
}

// Correct: track only the needed field
function ExpensiveComponent(props: { data: Data }) {
  const result = createMemo(() => {
    return expensiveCalculation(props.data.value); // only track value
  });
  // ...
}
```

### Lazy Loading

```tsx
import { lazy } from "solid-js";

// Code splitting
const HeavyComponent = lazy(() => import("./HeavyComponent"));

<Suspense fallback={<div>Loading...</div>}>
  <HeavyComponent />
</Suspense>
```

### on() - explicit dependencies

```tsx
import { on } from "solid-js";

// Run only when count changes (do not track other signals)
createEffect(on(count, (c) => {
  console.log("Count:", c);
  // Reading other signals here will not be tracked
}));
```

---

## 3) TypeScript Integration

### Component props types

```tsx
import { Component, JSX } from "solid-js";

// Use Component type
const Button: Component<{
  onClick: () => void;
  variant?: "primary" | "secondary";
  children: JSX.Element;
}> = (props) => {
  return (
    <button
      class={props.variant === "primary" ? "btn-primary" : "btn-secondary"}
      onClick={props.onClick}
    >
      {props.children}
    </button>
  );
};
```

### Signal type inference

```tsx
// Inferred
const [count, setCount] = createSignal(0); // number

// Explicit
const [user, setUser] = createSignal<User | null>(null);
```

---

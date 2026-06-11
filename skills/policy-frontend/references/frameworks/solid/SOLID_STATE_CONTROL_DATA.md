# SolidJS State, Control Flow, and Data Fetching - Detailed Guide

## Table of Contents

- 1) State Management
- 1) Control Flow
- 1) Data Fetching

## 1) State Management

### Store (complex state)

```tsx
import { createStore } from "solid-js/store";

// Store for nested objects
const [state, setState] = createStore({
  user: {
    name: "John",
    age: 30,
  },
  todos: [
    { id: 1, text: "Learn SolidJS", done: false },
  ],
});

// Update nested property
setState("user", "name", "Jane");

// Update array element
setState("todos", 0, "done", true);

// Use produce (similar to Immer)
import { produce } from "solid-js/store";

setState(produce(draft => {
  draft.user.age++;
  draft.todos.push({ id: 2, text: "New todo", done: false });
}));
```

### Context (shared across components)

```tsx
import { createContext, useContext } from "solid-js";

// Define context
const CounterContext = createContext<{
  count: () => number;
  increment: () => void;
}>();

// Provider
function CounterProvider(props: { children: any }) {
  const [count, setCount] = createSignal(0);

  const value = {
    count,
    increment: () => setCount(c => c + 1),
  };

  return (
    <CounterContext.Provider value={value}>
      {props.children}
    </CounterContext.Provider>
  );
}

// Consumer
function CounterDisplay() {
  const counter = useContext(CounterContext);
  if (!counter) throw new Error("CounterContext not found");

  return <div>{counter.count()}</div>;
}
```

---

## 2) Control Flow

### Show, For, Switch

```tsx
import { Show, For, Switch, Match } from "solid-js";

// Conditional rendering
<Show when={user()} fallback={<div>Loading...</div>}>
  {(u) => <div>Hello, {u().name}</div>}
</Show>

// List rendering (use stable keys)
<For each={items()}>
  {(item, index) => (
    <div>
      {index()}: {item.name}
    </div>
  )}
</For>

// Switch/Match (multiple conditions)
<Switch fallback={<div>Unknown</div>}>
  <Match when={status() === "loading"}>
    <Spinner />
  </Match>
  <Match when={status() === "error"}>
    <Error />
  </Match>
  <Match when={status() === "success"}>
    <Content />
  </Match>
</Switch>
```

Why not `if` or `map`?

- `<Show>`, `<For>`, etc. have optimized reactivity handling
- Avoid unnecessary DOM re-creation

---

## 3) Data Fetching

### createResource

```tsx
import { createResource } from "solid-js";

// Basic usage
const [user] = createResource(fetchUser);

// With params
const [userId, setUserId] = createSignal(1);
const [user] = createResource(userId, fetchUserById);

// Loading and error states
function UserProfile() {
  const [user] = createResource(fetchUser);

  return (
    <Show when={user()} fallback={<div>Loading...</div>}>
      {(u) => <div>{u().name}</div>}
    </Show>
  );
}

// Error handling
import { ErrorBoundary } from "solid-js";

<ErrorBoundary fallback={(err) => <div>Error: {err.message}</div>}>
  <UserProfile />
</ErrorBoundary>
```

### Suspense (parallel loading)

```tsx
import { Suspense } from "solid-js";

<Suspense fallback={<div>Loading...</div>}>
  <UserProfile />
  <PostList />
  <Comments />
</Suspense>
```

---

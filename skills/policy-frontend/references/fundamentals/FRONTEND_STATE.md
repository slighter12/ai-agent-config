---
description: "Frontend state management principles - server state vs client state"
---

# FRONTEND_STATE.md - Frontend State Management

This file defines state management principles for frontend projects.
First-class targets: SolidJS, Svelte, React.
Vue code patterns may exist in-repo but are out-of-scope by default unless a repo-specific policy explicitly includes Vue.

Violating these rules is incorrect output.

---

## Table of Contents

- 1) State Categories (mandatory understanding)
- 1) Server State Management (hard rules)
- 1) Client State Management (recommended)
- 1) State Synchronization Rules (hard rules)
- 1) Caching Strategy (recommended)
- 1) Optimistic Updates
- 1) State Initialization
- 1) When Uncertain (mandatory)

## 1) State Categories (mandatory understanding)

### Server State

Definition: Data from backend APIs

- User lists, orders, posts
- Any data fetched from an API
- Data that must stay in sync with the server

Characteristics:

- Can become stale
- Requires refetching
- Can be modified by other users
- Has loading and error states

### Client State

Definition: UI-only temporary state

- Modal open/close
- Form input values (before submit)
- Tab selection
- Sidebar expanded/collapsed

Characteristics:

- Client-only
- No server sync
- Lost on reload (usually acceptable)

---

## 2) Server State Management (hard rules)

### Do not store server data as client state

Do not do this:

```typescript
// React example (other frameworks similar)
const [users, setUsers] = useState([])

useEffect(() => {
  fetch('/api/users')
    .then(res => res.json())
    .then(data => setUsers(data)) // Do not do this
}, [])
```

Problems:

- No automatic refetching
- No caching
- No sharing across components
- Hard to manage loading/error states
- Data can go stale

Correct approach: use a server state library

SolidJS:

```typescript
import { createQuery } from '@tanstack/solid-query'

function UserList() {
  const users = createQuery(() => ({
    queryKey: ['users'],
    queryFn: () => fetch('/api/users').then(r => r.json())
  }))

  return (
    <Show when={users.data}>
      {data => <div>{/* render users */}</div>}
    </Show>
  )
}
```

Svelte:

```typescript
import { createQuery } from '@tanstack/svelte-query'

const users = createQuery({
  queryKey: ['users'],
  queryFn: () => fetch('/api/users').then(r => r.json())
})

{#if $users.isLoading}
  Loading...
{:else if $users.error}
  Error: {$users.error.message}
{:else if $users.data}
  {#each $users.data as user}
    {/* render */}
  {/each}
{/if}
```

React:

```typescript
import { useQuery } from '@tanstack/react-query'

function UserList() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['users'],
    queryFn: () => fetch('/api/users').then(r => r.json())
  })

  if (isLoading) return <div>Loading...</div>
  if (error) return <div>Error: {error.message}</div>

  return <div>{/* render users */}</div>
}
```

---

## 3) Client State Management (recommended)

### Good use cases for client state

Use framework built-in state tools:

- SolidJS: `createSignal`, `createStore`
- Svelte: `writable`, `derived`
- React: `useState`, `useReducer`
- Vue: `ref`, `reactive` (out-of-scope by default; use only when repo-specific policy explicitly includes Vue)

### Global client state

When sharing across components:

- SolidJS: Context + `createSignal`
- Svelte: stores (`writable`)
- React: Context + `useState` or Zustand
- Vue: Pinia (out-of-scope by default; use only when repo-specific policy explicitly includes Vue)

---

## 4) State Synchronization Rules (hard rules)

### One-way data flow

- Server state -> UI (read)
- UI -> Server (write)
- Do not two-way bind server data

```typescript
// Incorrect: mutate server state directly
function UserProfile() {
  const user = createQuery(/* ... */)

  const handleNameChange = (newName) => {
    user.data.name = newName // Do not do this
  }
}

// Correct: update via mutation
function UserProfile() {
  const user = createQuery(/* ... */)
  const updateUser = createMutation({
    mutationFn: (data) => fetch('/api/user', {
      method: 'PUT',
      body: JSON.stringify(data)
    }),
    onSuccess: () => {
      queryClient.invalidateQueries(['user']) // refetch
    }
  })

  const handleNameChange = (newName) => {
    updateUser.mutate({ name: newName })
  }
}
```

---

## 5) Caching Strategy (recommended)

### When to cache

- Data that rarely changes (settings, option lists)
- Data used in many places (user profile)
- Expensive API calls

### When not to cache

- Real-time data (stock prices, chat messages)
- Sensitive data (personal, financial)

---

## 6) Optimistic Updates

Use cautiously. Good for:

- High-likelihood success actions (like, favorite)
- UI that needs instant feedback

Not recommended for:

- Critical actions (payments, deletes)
- Complex validation flows

---

## 7) State Initialization

### Avoid empty-state flicker

```typescript
// Incorrect: flickers empty state
const [data, setData] = useState(null)

// Correct: use loading state
const [data, setData] = useState(null)
const [isLoading, setIsLoading] = useState(true)
```

Or rely on server state libraries to handle it.

---

## 8) When Uncertain (mandatory)

If any of the following are unclear, stop and ask:

- Whether a piece of state is server or client state
- Which server state library the project uses
- How to handle cache invalidation
- Rollback strategy for optimistic updates

---

Violating these rules is incorrect output.

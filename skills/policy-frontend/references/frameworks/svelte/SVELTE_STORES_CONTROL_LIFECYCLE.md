# Svelte Stores, Control Flow, and Lifecycle - Detailed Guide

This file targets Svelte 5 runes mode.

## Table of Contents

- 1) State Management (Stores)
- 1) Control Flow
- 1) Lifecycle

## 1) State Management (Stores)

### Writable Store

```ts
// stores.ts
import { writable } from 'svelte/store';

export const count = writable(0);

export const user = writable<User | null>(null);
```

```svelte
<script lang="ts">
  import { count } from './stores';

  // Auto-subscribe with $ prefix
  // Auto-unsubscribe on component destroy
</script>

<p>Count: {$count}</p>
<button onclick={() => count.update(n => n + 1)}>Increment</button>
<button onclick={() => count.set(0)}>Reset</button>
```

### Readable Store

```ts
// stores.ts
import { readable } from 'svelte/store';

export const time = readable(new Date(), (set) => {
  const interval = setInterval(() => {
    set(new Date());
  }, 1000);

  return () => clearInterval(interval); // cleanup
});
```

### Derived Store

```ts
import { derived } from 'svelte/store';
import { count } from './stores';

export const doubled = derived(count, ($count) => $count * 2);

// Multiple stores
export const fullName = derived(
  [firstName, lastName],
  ([$firstName, $lastName]) => `${$firstName} ${$lastName}`
);
```

### Custom Store

```ts
// counterStore.ts
import { writable } from 'svelte/store';

function createCounter() {
  const { subscribe, set, update } = writable(0);

  return {
    subscribe,
    increment: () => update(n => n + 1),
    decrement: () => update(n => n - 1),
    reset: () => set(0),
  };
}

export const counter = createCounter();
```

```svelte
<script lang="ts">
  import { counter } from './counterStore';
</script>

<p>{$counter}</p>
<button onclick={counter.increment}>+</button>
<button onclick={counter.decrement}>-</button>
<button onclick={counter.reset}>Reset</button>
```

---

## 2) Control Flow

### {#if}, {:else if}, {:else}

```svelte
<script lang="ts">
  let status = $state<'loading' | 'success' | 'error'>('loading');
</script>

{#if status === 'loading'}
  <p>Loading...</p>
{:else if status === 'error'}
  <p>Error occurred</p>
{:else}
  <p>Success!</p>
{/if}
```

### {#each}

```svelte
<script lang="ts">
  let items = $state([
    { id: 1, name: 'Apple' },
    { id: 2, name: 'Banana' },
  ]);
</script>

<!-- Use key (id) -->
{#each items as item (item.id)}
  <div>{item.name}</div>
{:else}
  <p>No items</p>
{/each}

<!-- With index -->
{#each items as item, index (item.id)}
  <div>{index}: {item.name}</div>
{/each}
```

### {#await}

```svelte
<script lang="ts">
  let promise = fetchData();
</script>

{#await promise}
  <p>Loading...</p>
{:then data}
  <p>Data: {data}</p>
{:catch error}
  <p>Error: {error.message}</p>
{/await}
```

---

## 3) Lifecycle

```svelte
<script lang="ts">
  import { onMount, onDestroy } from 'svelte';

  // onMount - runs once after mount
  onMount(() => {
    console.log('Component mounted');

    // return cleanup
    return () => {
      console.log('Cleanup on unmount');
    };
  });

  // onDestroy - before destroy
  onDestroy(() => {
    console.log('Component destroyed');
  });

  $effect.pre(() => {
    console.log('Before reactive DOM update');
  });

  $effect(() => {
    console.log('After reactive DOM update');
  });
</script>
```

---

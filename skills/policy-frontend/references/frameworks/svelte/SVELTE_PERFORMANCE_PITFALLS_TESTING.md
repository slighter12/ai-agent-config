# Svelte Performance, Pitfalls, and Testing - Detailed Guide

## Table of Contents

- 1) Performance Optimization
- 1) Common Pitfalls
- 1) Testing
- Summary

## 1) Performance Optimization

### Immutable Data

```ts
// Correct: create new arrays/objects
items = [...items, newItem];
user = { ...user, name: 'New Name' };

// Incorrect: mutate in place (Svelte may not detect)
items.push(newItem);
user.name = 'New Name';
```

### Keyed Each Blocks

```svelte
<!-- Use key -->
{#each items as item (item.id)}
  <Item {item} />
{/each}

<!-- No key (worse performance) -->
{#each items as item}
  <Item {item} />
{/each}
```

---

## 2) Common Pitfalls

### 1) Reactive Statement Order

```svelte
<script lang="ts">
  let count = 0;

  // May not execute in expected order
  $: doubled = count * 2;
  $: quadrupled = doubled * 2;

  // Better: express explicit dependency
  $: quadrupled = (count * 2) * 2;
</script>
```

### 2) Store Subscription Leaks

```svelte
<script lang="ts">
  import { myStore } from './stores';

  // Use $ prefix (auto unsubscribe)
  $: value = $myStore;

  // Manual subscription (must unsubscribe)
  let value;
  const unsubscribe = myStore.subscribe(v => value = v);
  onDestroy(unsubscribe); // remember cleanup
</script>
```

### 3) Binding and Reactivity

```svelte
<script lang="ts">
  let user = { name: 'John' };

  // Binding directly to object property (may lose reactivity)
  <input bind:value={user.name} />

  // Manual update or use a store
  <input value={user.name} on:input={(e) => user = { ...user, name: e.currentTarget.value }} />
</script>
```

---

## 3) Testing

### Vitest + Svelte Testing Library

```ts
import { render, fireEvent } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import Counter from './Counter.svelte';

describe('Counter', () => {
  it('increments count', async () => {
    const { getByRole } = render(Counter);
    const button = getByRole('button');

    await fireEvent.click(button);

    expect(button).toHaveTextContent('1');
  });
});
```

---

## Summary

Svelte strengths:

1. No Virtual DOM - compiles to efficient imperative code
2. True reactivity - no hooks or special APIs
3. Less boilerplate - concise and direct
4. Strong performance - smaller bundles

Remember:

- `$:` is Svelte magic - auto dependency tracking
- Stores use `$` prefix for auto subscribe/unsubscribe
- Update data immutably (new objects/arrays)
- SvelteKit provides a full-stack solution

References:

- Official docs: <https://svelte.dev/docs>
- Tutorial: <https://learn.svelte.dev/>
- SvelteKit: <https://kit.svelte.dev/docs>

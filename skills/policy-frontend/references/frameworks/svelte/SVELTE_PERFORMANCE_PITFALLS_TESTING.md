# Svelte Performance, Pitfalls, and Testing - Detailed Guide

This file targets Svelte 5 runes mode.

## Table of Contents

- 1) Performance Optimization
- 1) Common Pitfalls
- 1) Testing
- Summary

## 1) Performance Optimization

### $state vs $state.raw Updates

```svelte
<script lang="ts">
  type Item = { id: number; name: string };

  const newItem: Item = { id: 1, name: 'New Item' };
  let items = $state<Item[]>([]);
  let user = $state({ name: 'Old Name' });

  // $state proxies: direct mutation is reactive.
  items.push(newItem);
  user.name = 'New Name';

  let rawItems = $state.raw<Item[]>([]);
  let rawUser = $state.raw({ name: 'Old Name' });

  // $state.raw values are not deeply reactive; reassign the whole value.
  rawItems = [...rawItems, newItem];
  rawUser = { ...rawUser, name: 'New Name' };
</script>
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

### 1) Derived State Shape

```svelte
<script lang="ts">
  let count = $state(0);

  const doubled = $derived(count * 2);
  const quadrupled = $derived(doubled * 2);

  // Prefer one expression when the derived value does not need a name.
  const directQuadrupled = $derived((count * 2) * 2);
</script>
```

### 2) Store Subscription Leaks

```svelte
<script lang="ts">
  import { myStore } from './stores';
</script>

<p>{$myStore}</p>
```

### 3) Binding and Reactivity

```svelte
<script lang="ts">
  let user = $state({ name: 'John' });
</script>

<input bind:value={user.name} />
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

- Use `$state`, `$derived`, and `$effect` for runes-mode reactivity
- Stores use `$` prefix for auto subscribe/unsubscribe
- Use `$state` proxies for reactive mutation; use `$state.raw` only when reassignment-only updates are intentional
- SvelteKit provides a full-stack solution

References:

- Official docs: <https://svelte.dev/docs>
- Tutorial: <https://learn.svelte.dev/>
- SvelteKit: <https://kit.svelte.dev/docs>

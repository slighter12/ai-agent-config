# Svelte Reactivity and Component Patterns - Detailed Guide

## Table of Contents

- 1) Core Reactivity Concepts
- 1) Component Patterns

## 1) Core Reactivity Concepts

### Reactive Statements ($:)

Svelte reactivity is based on compile-time analysis. No hooks or special APIs are required.

```svelte
<script lang="ts">
  let count = 0;

  // Reactive statement - runs when count changes
  $: doubled = count * 2;

  // Reactive block
  $: {
    console.log(`Count is now ${count}`);
    if (count > 10) {
      alert('Count is high!');
    }
  }

  // Reactive function call
  $: updateTitle(count);

  function updateTitle(value: number) {
    document.title = `Count: ${value}`;
  }
</script>

<button on:click={() => count++}>{count}</button>
<p>Doubled: {doubled}</p>
```

Key principles:

- `$:` re-runs automatically when dependencies change
- Dependencies are analyzed at compile time, not runtime
- No manual dependency lists (unlike React useEffect)

---

## 2) Component Patterns

### Props

```svelte
<script lang="ts">
  // Basic props
  export let name: string;
  export let age: number;

  // Default value
  export let role: string = 'user';

  // Optional props
  export let email: string | undefined = undefined;
</script>

<div>
  <p>Name: {name}</p>
  <p>Age: {age}</p>
  <p>Role: {role}</p>
  {#if email}
    <p>Email: {email}</p>
  {/if}
</div>
```

### Two-way Binding

```svelte
<script lang="ts">
  let text = '';
  let checked = false;
  let selected = '';
</script>

<!-- Two-way binding -->
<input type="text" bind:value={text} />
<input type="checkbox" bind:checked />
<select bind:value={selected}>
  <option value="a">A</option>
  <option value="b">B</option>
</select>

<!-- Component two-way binding -->
<ChildComponent bind:value={text} />
```

### Events

```svelte
<!-- Parent.svelte -->
<script lang="ts">
  import Child from './Child.svelte';

  function handleCustomEvent(event: CustomEvent<{ data: string }>) {
    console.log(event.detail.data);
  }
</script>

<Child on:custom={handleCustomEvent} />

<!-- Child.svelte -->
<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher<{
    custom: { data: string };
  }>();

  function handleClick() {
    dispatch('custom', { data: 'Hello from child' });
  }
</script>

<button on:click={handleClick}>Click me</button>
```

---

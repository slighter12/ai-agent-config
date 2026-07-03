# Svelte Reactivity and Component Patterns - Detailed Guide

This file targets Svelte 5 runes mode.

## Table of Contents

- 1) Core Reactivity Concepts
- 1) Component Patterns

## 1) Core Reactivity Concepts

### Runes

Use runes for component-local reactive state, derived values, and side effects.

```svelte
<script lang="ts">
  let count = $state(0);

  const doubled = $derived(count * 2);

  $effect(() => {
    console.log(`Count is now ${count}`);
    updateTitle(count);
  });

  function updateTitle(value: number) {
    document.title = `Count: ${value}`;
  }
</script>

<button onclick={() => count++}>{count}</button>
<p>Doubled: {doubled}</p>
```

Key principles:

- `$state` declares mutable reactive state.
- `$derived` declares values computed from reactive state.
- `$effect` runs side effects after reactive dependencies change.
- No manual dependency lists are needed.

---

## 2) Component Patterns

### Props

```svelte
<script lang="ts">
  type Props = {
    name: string;
    age: number;
    role?: string;
    email?: string;
  };

  let {
    name,
    age,
    role = 'user',
    email,
  }: Props = $props();
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
  let text = $state('');
  let checked = $state(false);
  let selected = $state('');
</script>

<input type="text" bind:value={text} />
<input type="checkbox" bind:checked />
<select bind:value={selected}>
  <option value="a">A</option>
  <option value="b">B</option>
</select>

<ChildComponent bind:value={text} />
```

### Events

```svelte
<!-- Parent.svelte -->
<script lang="ts">
  import Child from './Child.svelte';

  function handleCustom(data: string) {
    console.log(data);
  }
</script>

<Child onCustom={handleCustom} />

<!-- Child.svelte -->
<script lang="ts">
  type Props = {
    onCustom: (data: string) => void;
  };

  let { onCustom }: Props = $props();

  function handleClick() {
    onCustom('Hello from child');
  }
</script>

<button onclick={handleClick}>Click me</button>
```

---

# Svelte Animations and Transitions - Detailed Guide

## 1) Animations and Transitions

### Transition

```svelte
<script lang="ts">
  import { fade, fly, slide } from 'svelte/transition';
  let visible = true;
</script>

{#if visible}
  <div transition:fade>Fade in/out</div>
  <div transition:fly={{ y: 200, duration: 500 }}>Fly in/out</div>
  <div transition:slide>Slide in/out</div>
{/if}
```

### Animation (Flip)

```svelte
<script lang="ts">
  import { flip } from 'svelte/animate';
  import { quintOut } from 'svelte/easing';

  let items = [1, 2, 3, 4, 5];

  function shuffle() {
    items = items.sort(() => Math.random() - 0.5);
  }
</script>

<button on:click={shuffle}>Shuffle</button>

{#each items as item (item)}
  <div animate:flip={{ duration: 500, easing: quintOut }}>
    {item}
  </div>
{/each}
```

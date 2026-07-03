# Svelte Actions (Use Directives) - Detailed Guide

This file targets Svelte 5 runes mode.

## 1) Actions (Use Directives)

```ts
// actions.ts
export function clickOutside(node: HTMLElement, onOutside: () => void) {
  function handleClick(event: MouseEvent) {
    if (!node.contains(event.target as Node)) {
      onOutside();
    }
  }

  document.addEventListener('click', handleClick, true);

  return {
    destroy() {
      document.removeEventListener('click', handleClick, true);
    },
  };
}
```

```svelte
<script lang="ts">
  import { clickOutside } from './actions';

  let open = $state(false);
</script>

<div use:clickOutside={() => open = false}>
  {#if open}
    <p>Modal content</p>
  {/if}
</div>
```

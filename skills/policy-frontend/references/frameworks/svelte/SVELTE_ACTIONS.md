# Svelte Actions (Use Directives) - Detailed Guide

## 1) Actions (Use Directives)

```ts
// actions.ts
export function clickOutside(node: HTMLElement) {
  function handleClick(event: MouseEvent) {
    if (!node.contains(event.target as Node)) {
      node.dispatchEvent(new CustomEvent('outclick'));
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

  let open = false;
</script>

<div use:clickOutside on:outclick={() => open = false}>
  {#if open}
    <p>Modal content</p>
  {/if}
</div>
```

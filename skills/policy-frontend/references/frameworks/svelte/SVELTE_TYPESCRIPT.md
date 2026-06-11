# Svelte TypeScript Integration - Detailed Guide

## 1) TypeScript Integration

```svelte
<script lang="ts">
  import type { ComponentType } from "svelte";

  interface User {
    id: number;
    name: string;
  }

  let users: User[] = [];

  export let component: ComponentType;
</script>
```

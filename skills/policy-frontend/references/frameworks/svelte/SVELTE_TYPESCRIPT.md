# Svelte TypeScript Integration - Detailed Guide

This file targets Svelte 5 runes mode.

## 1) TypeScript Integration

```svelte
<script lang="ts">
  import type { Component, ComponentProps } from 'svelte';
  import type { HTMLButtonAttributes } from 'svelte/elements';

  interface User {
    id: number;
    name: string;
  }

  type Props = {
    component: Component<{ user: User }>;
    buttonProps?: HTMLButtonAttributes;
  };

  let users = $state<User[]>([]);
  let { component: UserCard, buttonProps }: Props = $props();

  type UserCardProps = ComponentProps<typeof UserCard>;
  const fallbackUser: UserCardProps['user'] = { id: 0, name: 'Guest' };
</script>

<UserCard user={users[0] ?? fallbackUser} />
<button {...buttonProps}>Save</button>
```

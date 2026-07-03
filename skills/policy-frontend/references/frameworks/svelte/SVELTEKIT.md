# SvelteKit - Detailed Guide

This file targets Svelte 5 runes mode.

## 1) SvelteKit

### Routing

```
src/routes/
|-- +page.svelte          # /
|-- about/
|   `-- +page.svelte      # /about
|-- blog/
|   |-- +page.svelte      # /blog
|   `-- [slug]/
|       `-- +page.svelte  # /blog/[slug]
`-- api/
    `-- users/
        `-- +server.ts    # /api/users
```

### Data Loading

```ts
// +page.ts (runs on both server and client)
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
  const response = await fetch(`/api/posts/${params.id}`);
  const post = await response.json();

  return {
    post,
  };
};
```

```svelte
<!-- +page.svelte -->
<script lang="ts">
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();
</script>

<h1>{data.post.title}</h1>
<p>{data.post.content}</p>
```

### Server-only Load

```ts
// +page.server.ts (runs only on server)
import type { PageServerLoad } from './$types';
import { db } from '$lib/server/database';

export const load: PageServerLoad = async () => {
  const users = await db.query('SELECT * FROM users');

  return {
    users,
  };
};
```

### Form Actions

```ts
// +page.server.ts
import type { Actions } from './$types';

export const actions: Actions = {
  default: async ({ request }) => {
    const data = await request.formData();
    const name = data.get('name');

    // Save to database...

    return { success: true };
  },
};
```

```svelte
<!-- +page.svelte -->
<script lang="ts">
  import { enhance } from '$app/forms';
</script>

<form method="POST" use:enhance>
  <input name="name" required />
  <button type="submit">Submit</button>
</form>
```

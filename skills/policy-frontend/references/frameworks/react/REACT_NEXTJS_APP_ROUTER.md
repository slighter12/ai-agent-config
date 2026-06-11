# Next.js App Router - Detailed Guide

## Table of Contents

- 1) Next.js (App Router)

## 1) Next.js (App Router)

### File Structure

```
app/
|-- layout.tsx              # Root layout
|-- page.tsx                # /
|-- about/
|   `-- page.tsx            # /about
|-- blog/
|   |-- page.tsx            # /blog
|   `-- [slug]/
|       `-- page.tsx        # /blog/[slug]
`-- api/
    `-- users/
        `-- route.ts        # /api/users
```

### Server Components (default)

```tsx
// app/blog/[slug]/page.tsx
async function BlogPost({ params }: { params: { slug: string } }) {
  // Fetch directly in the component (Server Component)
  const post = await fetchPost(params.slug);

  return (
    <article>
      <h1>{post.title}</h1>
      <div>{post.content}</div>
    </article>
  );
}

export default BlogPost;
```

### Client Components

```tsx
'use client'; // Mark as Client Component

import { useState } from 'react';

export function Counter() {
  const [count, setCount] = useState(0);

  return (
    <button onClick={() => setCount(c => c + 1)}>
      Count: {count}
    </button>
  );
}
```

### Data Fetching

```tsx
// Server Component - parallel fetching
async function Page() {
  const [user, posts] = await Promise.all([
    fetchUser(),
    fetchPosts(),
  ]);

  return (
    <div>
      <UserInfo user={user} />
      <PostList posts={posts} />
    </div>
  );
}

// Streaming with Suspense
import { Suspense } from 'react';

async function Page() {
  return (
    <div>
      <Suspense fallback={<div>Loading user...</div>}>
        <UserInfo />
      </Suspense>
      <Suspense fallback={<div>Loading posts...</div>}>
        <PostList />
      </Suspense>
    </div>
  );
}
```

### Server Actions

```tsx
// app/actions.ts
'use server';

export async function createPost(formData: FormData) {
  const title = formData.get('title') as string;
  const content = formData.get('content') as string;

  await db.post.create({ title, content });

  revalidatePath('/blog');
}
```

```tsx
// app/new-post/page.tsx
import { createPost } from '../actions';

export default function NewPost() {
  return (
    <form action={createPost}>
      <input name="title" required />
      <textarea name="content" required />
      <button type="submit">Create</button>
    </form>
  );
}
```

---

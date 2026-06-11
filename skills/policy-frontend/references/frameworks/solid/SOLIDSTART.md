# SolidStart - Detailed Guide

## 1) SolidStart

### Server Functions

```tsx
"use server";

export async function getServerData() {
  // Runs on the server only
  return await db.query("SELECT * FROM users");
}
```

### Actions

```tsx
import { createServerAction$ } from "solid-start/server";

export default function Form() {
  const [submitting, submit] = createServerAction$(async (formData: FormData) => {
    "use server";
    const name = formData.get("name");
    await saveToDatabase(name);
  });

  return (
    <form action={submit} method="post">
      <input name="name" />
      <button type="submit" disabled={submitting.pending}>
        Submit
      </button>
    </form>
  );
}
```

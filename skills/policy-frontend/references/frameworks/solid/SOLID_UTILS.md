# SolidJS Utilities - Detailed Guide

## 1) Utilities

### batch (batch updates)

```tsx
import { batch } from "solid-js";

// Batch updates (single effect run)
batch(() => {
  setFirstName("John");
  setLastName("Doe");
  setAge(30);
});
```

### untrack (do not track)

```tsx
import { untrack } from "solid-js";

createEffect(() => {
  const c = count();
  const otherValue = untrack(() => otherSignal()); // do not track this
  console.log(c, otherValue);
});
```

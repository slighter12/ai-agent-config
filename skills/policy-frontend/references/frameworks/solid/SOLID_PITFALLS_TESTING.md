# SolidJS Pitfalls and Testing - Detailed Guide

## 1) Common Pitfalls

### 1) Props Destructuring

```tsx
// Incorrect
function Component(props: { value: number }) {
  const { value } = props; // loses reactivity
  return <div>{value}</div>;
}

// Correct
function Component(props: { value: number }) {
  return <div>{props.value}</div>;
}
```

### 2) Async at the Component Top Level

```tsx
// Incorrect
async function Component() {
  const data = await fetchData(); // will not suspend
  return <div>{data}</div>;
}

// Correct
function Component() {
  const [data] = createResource(fetchData);
  return <Show when={data()}>{(d) => <div>{d()}</div>}</Show>;
}
```

### 3) Effect Dependencies

```tsx
// Incorrect: does not track count
createEffect(() => {
  const c = count; // not called
  console.log(c);
});

// Correct
createEffect(() => {
  console.log(count()); // must call
});
```

## 2) Testing

### Vitest + Solid Testing Library

```tsx
import { render } from "@solidjs/testing-library";
import { describe, it, expect } from "vitest";

describe("Counter", () => {
  it("increments count", async () => {
    const { getByRole } = render(() => <Counter />);
    const button = getByRole("button");

    button.click();

    expect(button).toHaveTextContent("1");
  });
});
```

## Summary

SolidJS strengths:

1. Fine-grained reactivity - update only what changes
2. No Virtual DOM - direct DOM updates, strong performance
3. Familiar JSX syntax - similar to React, more direct
4. TypeScript-first - full typing support

Remember:

- Signals must be called: `count()`
- Do not destructure props unless using `splitProps`
- Use `<Show>`, `<For>` instead of `if`, `map`
- Components run once; reactivity updates expressions

References:

- Official docs: <https://www.solidjs.com/>
- Tutorial: <https://www.solidjs.com/tutorial>
- Playground: <https://playground.solidjs.com/>

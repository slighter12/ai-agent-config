# React Pitfalls and Testing - Detailed Guide

## 1) Common Pitfalls

### 1) Missing useEffect Dependencies

```tsx
// Incorrect: missing deps
useEffect(() => {
  fetchData(userId); // userId not in dependency array
}, []);

// Correct
useEffect(() => {
  fetchData(userId);
}, [userId]);
```

### 2) Overusing useCallback/useMemo

```tsx
// Unnecessary optimization
const handleClick = useCallback(() => {
  console.log('clicked');
}, []);

// Simple handlers do not need useCallback
const handleClick = () => {
  console.log('clicked');
};
```

### 3) State Updates Based on Previous Value

```tsx
// Incorrect: can be stale
setCount(count + 1);

// Correct: functional update
setCount(c => c + 1);
```

### 4) Missing Key Props

```tsx
// Incorrect: no key or using index
{items.map((item, index) => <Item key={index} item={item} />)}

// Correct: use stable IDs
{items.map(item => <Item key={item.id} item={item} />)}
```

---

## 2) Testing

### Vitest + React Testing Library

```tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import Counter from './Counter';

describe('Counter', () => {
  it('increments count when button is clicked', async () => {
    render(<Counter />);
    const button = screen.getByRole('button', { name: /increment/i });

    fireEvent.click(button);

    expect(screen.getByText(/count: 1/i)).toBeInTheDocument();
  });
});
```

---

## Summary

React strengths:

1. Declarative UI - easier to understand and maintain
2. Component-based - reusable building blocks
3. Rich ecosystem - many third-party libraries
4. Next.js provides a full-stack solution

Remember:

- Follow the Rules of Hooks
- useEffect dependencies must be complete
- Use React Query for server state
- Next.js App Router prefers Server Components
- Profile before optimizing (avoid premature optimization)

References:

- React docs: <https://react.dev/>
- Next.js docs: <https://nextjs.org/docs>
- TanStack Query: <https://tanstack.com/query/latest>
- React Hook Form: <https://react-hook-form.com/>

# React Hooks - Detailed Guide

## Table of Contents

- 1) Core Hook Concepts

## 1) Core Hook Concepts

### useState

```tsx
import { useState } from 'react';

function Counter() {
  const [count, setCount] = useState(0);

  // Functional update (based on previous value)
  const increment = () => setCount(prev => prev + 1);

  // Direct set
  const reset = () => setCount(0);

  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={increment}>+1</button>
      <button onClick={reset}>Reset</button>
    </div>
  );
}
```

### useEffect

```tsx
import { useEffect, useState } from 'react';

function Example() {
  const [data, setData] = useState(null);
  const [count] = useState(0);

  // Basic usage
  useEffect(() => {
    fetchData().then(setData);
  }, []); // Empty deps = run once

  // With deps
  useEffect(() => {
    console.log('Count changed:', count);
  }, [count]); // runs when count changes

  // Cleanup
  useEffect(() => {
    const subscription = subscribe();
    return () => {
      subscription.unsubscribe();
    };
  }, []);

  return <div>{data}</div>;
}
```

Key principles:

- Dependency arrays must include all reactive values used inside the effect
- Use the ESLint plugin `react-hooks/exhaustive-deps`
- Cleanup runs on unmount or before deps change

### useMemo

```tsx
import { useMemo } from 'react';

function ExpensiveComponent({ items, filter }: Props) {
  // Memoize expensive computation
  const filteredItems = useMemo(() => {
    return items.filter(item => item.category === filter);
  }, [items, filter]); // recompute when deps change

  return <List items={filteredItems} />;
}
```

### useCallback

```tsx
import { useCallback } from 'react';

function Parent() {
  const [count, setCount] = useState(0);

  // Memoize callback (avoid unnecessary child re-renders)
  const handleClick = useCallback(() => {
    setCount(c => c + 1);
  }, []); // empty deps = stable reference

  return <ChildComponent onClick={handleClick} />;
}
```

### useRef

```tsx
import { useRef, useEffect } from 'react';

function TextInput() {
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    // DOM access
    inputRef.current?.focus();
  }, []);

  // Store mutable values without re-render
  const renderCount = useRef(0);
  renderCount.current++;

  return <input ref={inputRef} />;
}
```

---

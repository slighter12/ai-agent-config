# Seam Isolation & Mocking Discipline (Language-Agnostic)

## When to Isolate or Fake a Dependency

Isolate dependencies only at genuine **I/O or non-deterministic edges**:
* External network APIs (payment providers, third-party services, LLM APIs)
* Persistence / storage (via an in-memory fake repository, test DB, or narrow port) — prefer a real local stand-in when one exists
* Non-deterministic primitives (system clock, timers, RNG, OS process signals)

## Never Mock Internal Collaborators

* Never mock pure domain functions, state calculators, or internal helper modules you own within the same layer.
* **The >2 Mocks Smell**: apply the No Call-Only / Mock-Heavy Assertions redline in Step 2 of [SKILL.md](SKILL.md).
  For a true external system from the list above, read "The contract decides, not the assertion shape" in [tests.md](tests.md#what-makes-a-durable-behavioral-test) before asserting the request it receives.

## Designing for Mockability

At the edges you do isolate, shape the interface so the fake is trivial. Examples are illustrative TypeScript.

**1. Inject the dependency** rather than constructing it inside:

```typescript
// Easy to fake
function processPayment(order, paymentClient) {
  return paymentClient.charge(order.total);
}

// Hard to fake
function processPayment(order) {
  const client = new StripeClient(process.env.STRIPE_KEY);
  return client.charge(order.total);
}
```

**2. Prefer SDK-style interfaces over generic fetchers**

Create specific functions for each external operation instead of one generic function with conditional logic:

```typescript
// GOOD: Each function is independently fakeable
const api = {
  getUser: (id) => fetch(`/users/${id}`),
  getOrders: (userId) => fetch(`/users/${userId}/orders`),
  createOrder: (data) => fetch('/orders', { method: 'POST', body: data }),
};

// BAD: Faking requires conditional logic inside the fake
const api = {
  fetch: (endpoint, options) => fetch(endpoint, options),
};
```

The SDK approach means:
- Each fake returns one specific shape
- No conditional logic in test setup
- Easier to see which endpoints a test exercises
- Type safety per endpoint

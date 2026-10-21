# Seam Isolation & Mocking Discipline (Language-Agnostic)

## When to Isolate or Fake a Dependency

Isolate dependencies only at genuine **I/O or non-deterministic edges**:
* External network APIs (payment providers, third-party services, LLM APIs)
* Persistence / storage (via an in-memory fake repository, test DB, or narrow port) — prefer a real local stand-in when one exists
* Non-deterministic primitives (system clock, timers, RNG, OS process signals)

## Never Mock Internal Collaborators

* Never mock pure domain functions, state calculators, or internal helper modules you own within the same layer.
* **The >2 Mocks Smell**: If a unit test requires setting up more than 2 mocks/stubs, or only asserts that a mock was called with certain arguments (`mock.calledWith`) without asserting a returned value, state mutation, or error outcome, you are testing wiring rather than behavior. Delete the unit test and verify the wiring via runtime execution (`Prove It Works`).
  The exception is a true external system from the list above: the request or payload it receives is behavior when the contract promises it (see [tests.md](tests.md)), so assert it through a recording fake, paired with a positive case.

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

**2. Prefer SDK-style interfaces over generic fetchers**: one function per external operation, so each fake returns one specific shape with no conditional logic in test setup, and the test shows which endpoints it exercises.

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

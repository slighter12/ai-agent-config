# Good vs. Bad Tests & Post-Green Pruning

Examples are illustrative TypeScript; write the same shape in the package's own language and test idiom.

## The Two-Edit Test

Judge every test you write, change, or keep against two imagined edits:

* **Break the behavior**: make the code return a wrong value, skip an effect, or perform it twice. If the test still passes, it proves nothing. The usual causes are an assertion looser than the contract, an expected value computed by the code under test, a restated constant or prompt string that never runs the code that reads it, and an assertion on data the test built itself.
* **Change the internals**: keep every observable result the same and restructure the code underneath. If the test now fails, it is pinned to the implementation: it needs a rewrite on every refactor and still catches no bug.

Keep a test that catches the first edit and survives the second; rewrite or delete the rest. Drive the code through its users' interface with a concrete input, derive the expected value from the requirement rather than from a run of the code, and assert what the contract promises, no less and no more.

## What Makes a Durable Behavioral Test

* **Tests a Seam Contract**: Calls a public entrypoint or domain seam and asserts the resulting output state, return value, or error mode.
* **Survives Internal Refactoring**: Renaming or inlining an internal helper function breaks zero tests.
* **Copy-Neighbor Consistency**: Reuses the test harness, assertion helpers, and fixture patterns already established in the same package/directory.
* **Describes WHAT, not HOW**: The test name reads like a capability ("user can checkout with valid cart"), with one logical assertion per test.

```typescript
// GOOD: Tests observable behavior
test("user can checkout with valid cart", async () => {
  const cart = createCart();
  cart.add(product);
  const result = await checkout(cart, paymentMethod);
  expect(result.status).toBe("confirmed");
});

// BAD: Tests implementation details
test("checkout calls paymentService.process", async () => {
  const mockPayment = jest.mock(paymentService);
  await checkout(cart, payment);
  expect(mockPayment.process).toHaveBeenCalledWith(cart.total);
});
```

**The contract decides, not the assertion shape.** A call to an internal collaborator you own (the `paymentService` above) is implementation, because no user can see it. The request or payload a true external system receives (a payment provider, a message broker) is behavior when the contract promises it, because that system observes it. "An empty cart never charges" is a legitimate assertion that the provider fake recorded no charge, provided a paired positive case shows the test can see a charge. Likewise, the error type callers branch on is contract, while an exact error message no caller reads is not.

```typescript
// GOOD: The external provider observing (or not observing) a charge is the promised behavior
test("empty cart never charges", async () => {
  const provider = fakePaymentProvider();
  await checkout(createCart(), provider);
  expect(provider.charges).toEqual([]);
});

test("checkout charges the cart total once", async () => {
  const provider = fakePaymentProvider();
  await checkout(cartWith({ price: 15 }), provider);
  expect(provider.charges).toEqual([{ amount: 15 }]);
});
```

Red flags of implementation coupling: mocking internal collaborators, testing private methods, asserting on call counts/order, breaking on a refactor without a behavior change, a name that describes HOW, and verifying through a side channel instead of the interface:

```typescript
// BAD: Bypasses interface to verify
test("createUser saves to database", async () => {
  await createUser({ name: "Alice" });
  const row = await db.query("SELECT * FROM users WHERE name = ?", ["Alice"]);
  expect(row).toBeDefined();
});

// GOOD: Verifies through interface
test("createUser makes user retrievable", async () => {
  const user = await createUser({ name: "Alice" });
  const retrieved = await getUser(user.id);
  expect(retrieved.name).toBe("Alice");
});
```

## Four Anti-Patterns to Reject Immediately

1. **Single-Caller Helper Extraction**: Splitting an internal step out of an already-testable module into a function with only `1` caller in production code solely so a test file can import and test it directly. Inline it back and test through the caller. (Exempt: a **Brownfield Sprout** chosen at the Agency Gate to isolate new logic from an untestable legacy monolith — test it directly.)
2. **Tautological Assertions**: Computing the expected value in the test by copying the exact formula or logic from the production implementation. Expected values must be concrete literals or spec-derived fixtures.

   ```typescript
   // BAD: Expected value is recomputed the way the code computes it
   const expected = items.reduce((sum, i) => sum + i.price, 0);
   expect(calculateTotal(items)).toBe(expected);

   // GOOD: Expected value is an independent, known literal
   expect(calculateTotal([{ price: 10 }, { price: 5 }])).toBe(15);
   ```
3. **Normalized-Input Tests**: Taking a real failing production input/command and stripping out the flag, path, or precondition that actually triggers the bug just to make the test fixture easier to run. (The tell: removing the production fix does not turn the test red.)
4. **Same-Branch Bloat (Prune Before Commit)**: Keeping multiple test cases that walk the exact same conditional branch (`if/else/switch`) with different normal inputs. Keep 1 happy-path case plus 1 case per distinct edge, boundary (`0`, empty, nil, max), or error branch.

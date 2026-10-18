# Seam Isolation & Mocking Discipline (Language-Agnostic)

## When to Isolate or Fake a Dependency

Isolate dependencies only at genuine **I/O or non-deterministic boundaries**:
* External network APIs (payment providers, third-party services, LLM APIs)
* Persistence / storage boundaries (via an in-memory fake repository, test DB, or narrow boundary interface/port)
* Non-deterministic primitives (system clock, timers, RNG, OS process signals)

## Never Mock Internal Collaborators

* Never mock pure domain functions, state calculators, or internal helper modules you own within the same layer.
* **The >2 Mocks Smell**: If a unit test requires setting up more than 2 mocks/stubs, or only asserts that a mock was called with certain arguments (`mock.calledWith`) without asserting a returned value, state mutation, or error outcome, you are testing wiring rather than behavior. Delete the unit test and verify the wiring via runtime execution (`Prove It Works`).

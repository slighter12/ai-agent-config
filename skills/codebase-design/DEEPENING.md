# Deepening Modules & Seam Discipline

## Dependency Categories

Classify every dependency a module touches into one of three buckets:

1. **In-Memory Domain / Pure Logic**: Data structures, state transitions, domain calculations. Keep inside the module; never mock or place behind an artificial interface.
2. **Owned Boundary (Local DB / Filesystem / Queue)**: External to memory, owned by your system. Isolate at a single coarse port/adapter seam at the module edge so domain state transitions can be tested deterministically with an in-memory fake or lightweight harness.
3. **Third-Party External (Payment API, LLM API, Clock, Randomness)**: True external boundary. Always place behind a narrow seam at the system edge.

## Replace, Don't Layer

When an existing interface is shallow or leaky:
* Do not wrap a shallow interface in a second shallow wrapper (which doubles pass-through boilerplate).
* Pull the callers' repeated setup/validation/error-mapping logic **down into the module** so the interface shrinks and the implementation deepens.

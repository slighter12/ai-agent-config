# Design It Twice

When designing a critical new module interface or resolving a Scrap Circuit Breaker:

1. Generate **two radically different interface shapes** (e.g., Approach A: stateful service with injected ports vs. Approach B: pure state-transition function returning commands/events).
2. Compare both candidates across three concrete metrics:
   * **Caller Leverage**: How many lines of setup and error handling does every caller have to write?
   * **Seam Simplicity**: Can behavior be verified at the public interface with <= 2 test arguments/fakes?
   * **Local Legibility**: If a future agent copies a call site of this interface, can it misuse the ordering or state invariants?
3. Present both sketches and your recommendation to the user before committing to the implementation.

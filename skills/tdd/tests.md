# Good vs. Bad Tests & Post-Green Pruning

## What Makes a Durable Behavioral Test

* **Tests a Seam Contract**: Calls a public entrypoint or domain boundary and asserts the resulting output state, return value, or error mode.
* **Survives Internal Refactoring**: Renaming or inlining an internal helper function breaks zero tests.
* **Copy-Neighbor Consistency**: Reuses the test harness, assertion helpers, and fixture patterns already established in the same package/directory.

## Four Anti-Patterns to Reject Immediately

1. **Single-Caller Helper Extraction**: Splitting an internal step out of an already-testable module into a function with only `1` caller in production code solely so a test file can import and test it directly. Inline it back and test through the caller. (Exempt: a **Brownfield Sprout** chosen at the Agency Gate to isolate new logic from an untestable legacy monolith — test it directly.)
2. **Tautological Assertions**: Computing the expected value in the test by copying the exact formula or logic from the production implementation. Expected values must be concrete literals or spec-derived fixtures.
3. **Normalized-Input Tests**: Taking a real failing production input/command and stripping out the flag, path, or precondition that actually triggers the bug just to make the test fixture easier to run. (The tell: removing the production fix does not turn the test red.)
4. **Same-Branch Bloat (Prune Before Commit)**: Keeping multiple test cases that walk the exact same conditional branch (`if/else/switch`) with different normal inputs. Keep 1 happy-path case plus 1 case per distinct edge, boundary (`0`, empty, nil, max), or error branch.

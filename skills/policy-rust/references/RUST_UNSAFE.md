---
globs: ["*.rs", "**/*.rs"]
description: "Rust unsafe code - when to use, safety invariants, documentation"
---

# RUST_UNSAFE.md - Rust Unsafe Code

This file defines rules for unsafe code in Rust projects.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Minimize Unsafe (mandatory)

- Avoid unsafe whenever possible; find safe alternatives first
- Use unsafe only for: FFI, performance-critical paths, low-level abstractions
- Do not use unsafe for convenience

### 2) Keep Unsafe Scope Small

- Unsafe blocks should be as small as possible
- Include only the operations that truly require unsafe
- Wrap unsafe details behind a safe API immediately

### 3) Safety Invariants (mandatory)

- Every unsafe block must include a comment explaining why it is safe
- Describe the safety invariants (preconditions)
- Explain how those invariants are maintained

### 4) Public APIs Must Be Safe

- Public APIs must be safe to call
- Keep unsafe implementation details internal
- Use the type system to guarantee safety

### 5) Unsafe Operation Types

- Raw pointers (`*const T`, `*mut T`) - check null, alignment, validity
- FFI - check C API preconditions
- Transmute - avoid; allow only in rare cases with detailed comments
- Inline assembly - only when absolutely required (code review)

### 6) Test Unsafe Code

- Unsafe code must have full tests
- Use Miri to detect undefined behavior
- Use fuzzing for boundary cases

### 7) Document Unsafe

- Use a `# Safety` section to describe preconditions
- State caller responsibilities
- Provide correct usage examples

### 8) Common Unsafe Pitfalls

- Aliasing violations - mutable reference coexists with other references
- Use after free - dereferencing freed memory
- Null pointer dereference
- Uninitialized memory reads

### 9) Prefer Alternatives

- Need performance -> profile first, confirm bottleneck, try safe optimizations
- Need FFI -> use `bindgen` to generate bindings
- Need interior mutability -> use `Cell`, `RefCell`, `Mutex`

### 10) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Whether a safe alternative exists
- How to explain the safety invariants
- How to test the unsafe code
- Whether a code review is required

---

## Unsafe Template

```rust
/// # Safety
///
/// `ptr` must be:
/// - Non-null
/// - Properly aligned for `T`
/// - Point to a valid instance of `T`
/// - Not aliased by any mutable references
unsafe fn read_raw<T>(ptr: *const T) -> T {
    // SAFETY: the caller guarantees the preconditions above
    ptr.read()
}
```

---

## Tools

- Miri - undefined behavior detection
- ASAN/MSAN - address/memory sanitizers
- cargo-fuzz - fuzz testing

---

Violating these rules is incorrect output.

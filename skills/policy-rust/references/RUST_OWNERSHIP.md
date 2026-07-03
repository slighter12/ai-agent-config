---
description: "Rust ownership and borrowing - references, lifetimes, move semantics"
---

# RUST_OWNERSHIP.md - Rust Ownership and Borrowing

This file defines ownership and borrowing rules for Rust projects.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Ownership Basics

- Each value has a single owner
- When the owner leaves scope, the value is dropped
- Move semantics are default (non-`Copy` types)

### 2) Borrowing Rules (compiler enforced)

- Any number of immutable references (`&T`), or
- Exactly one mutable reference (`&mut T`)
- References must always be valid (cannot outlive the owner)

### 3) When to use `&`, `&mut`, or take ownership

- Read-only -> `&T`
- Mutate -> `&mut T`
- Need ownership (consume value, move to another thread) -> take ownership
- When unsure, prefer `&T`

### 4) Avoid Unnecessary Clone

- Use `.clone()` only when you need a real copy
- Prefer references or `Cow<T>`
- Avoid clone in hot paths (confirm with profiling)

### 5) Lifetime Naming

- Let the compiler infer in simple cases
- Use descriptive names for complex cases (`'a`, `'b` -> `'request`, `'static`)
- Avoid over-annotating lifetimes

### 6) String Type Choice

- Owned string -> `String`
- Borrowed string -> `&str`
- Function params should prefer `&str`
- Return `String` when ownership is required

### 7) Vec vs Slice

- Owned vector -> `Vec<T>`
- Borrowed slice -> `&[T]` or `&mut [T]`
- Function params should prefer `&[T]`

### 8) Interior Mutability

- Mutate in immutable context -> `Cell<T>`, `RefCell<T>`
- Multi-threaded -> `Arc<Mutex<T>>` or `Arc<RwLock<T>>`
- Avoid overuse (prefer redesigning the data model)

### 9) Drop and Resource Cleanup

- Implement `Drop` for custom cleanup
- RAII pattern: acquire in constructor, release in destructor
- Avoid manual `drop()` unless explicit (`std::mem::drop`)

### 10) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Whether to use `&T` or take ownership
- How to resolve lifetime errors
- Whether `Clone` is needed or references suffice
- Root cause of borrow checker errors

---

## Common Compiler Errors

- "cannot borrow as mutable" -> immutable borrow exists
- "use of moved value" -> value moved; use `.clone()` or reference
- "lifetime may not live long enough" -> reference lifetime too short

---

Violating these rules is incorrect output.

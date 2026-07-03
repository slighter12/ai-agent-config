---
description: "Rust concurrency - threads, async/await, channels, safety"
---

# RUST_CONCURRENCY.md - Rust Concurrency

This file defines concurrency rules for Rust projects.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Send and Sync Traits

- Understand `Send` (move across threads) and `Sync` (share references across threads)
- Most types implement these automatically; some types (`Rc`, `RefCell`) do not
- Use `Arc` instead of `Rc` in multithreaded contexts

### 2) Shared State Management

- Share across threads via `Arc<Mutex<T>>` or `Arc<RwLock<T>>`
- Prefer `RwLock` for read-heavy workloads
- Avoid deadlocks by keeping a consistent lock order

### 3) Thread Spawning

- Use `std::thread::spawn` or a thread pool (`rayon`)
- Ensure threads are joined or explicitly detached
- Use `move` closures to pass data into threads

### 4) Async/Await Rules

- Use async for I/O-bound work (network, file I/O)
- Use sync code or `spawn_blocking` for CPU-bound work
- Choose a runtime: `tokio` (general) or `async-std`

### 5) Channels

- Prefer channels for inter-thread communication (`std::sync::mpsc` or `crossbeam`)
- In async, use `tokio::sync::mpsc` or `async-channel`
- Single consumer -> `mpsc`; multiple consumers -> `mpmc` (`crossbeam` or `flume`)

### 6) Avoid Blocking in Async

- Do not run blocking operations inside async functions (file I/O, heavy compute)
- Use `tokio::task::spawn_blocking` for blocking code
- Do not use `std::thread::sleep` in async (use `tokio::time::sleep`)

### 7) Cancellation and Timeouts

- Use `tokio::time::timeout` for timeouts
- Use `tokio::select!` for cancellation logic
- Clean up resources in `Drop`

### 8) Race Conditions

- Rust's type system prevents data races at compile time
- Watch for logic races (not data races) manually
- Use atomics (`AtomicBool`, `AtomicUsize`) for simple shared state

### 9) Graceful Shutdown

- Use `tokio::signal` to listen for SIGTERM/SIGINT
- Use `CancellationToken` or a broadcast channel to notify tasks
- Wait for all tasks to finish before exit

### 10) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Whether to use async or sync threads
- How to avoid deadlocks
- How to fix `Send`/`Sync` bound errors
- How to implement graceful shutdown

---

## Recommended Crates

- `tokio` - async runtime (most common)
- `rayon` - data parallelism
- `crossbeam` - high-performance concurrency primitives
- `parking_lot` - faster `Mutex` and `RwLock`

---

Violating these rules is incorrect output.

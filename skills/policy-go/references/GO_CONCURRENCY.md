---
globs: ["*.go", "**/*.go"]
description: "Go concurrency and lifecycle management - goroutines, context, shutdown"
---

# GO_CONCURRENCY.md - Go Concurrency and Lifecycle

This file defines concurrency and lifecycle management rules for Go projects.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) No Fire-and-Forget Goroutines

- Every goroutine must be stoppable via `context.Context` or a channel
- Never run uncancelable infinite loops
- Always define explicit exit conditions

### 2) Context Passing Rules

- Context must be the first parameter, named `ctx`
- Do not store context in structs (unless required by an external interface)
- Long-running operations must respect context cancellation

### 3) Context Values Limits

- Avoid `context.WithValue` where possible
- Use only for request-scoped data (request_id, trace_id, user_id)
- Do not pass optional params or business logic data

### 4) Graceful Shutdown (mandatory)

- Applications must implement graceful shutdown
- SIGTERM/SIGINT -> stop accepting new requests -> wait for in-flight requests -> cancel background goroutines -> clean up resources -> exit
- Shutdown timeout recommended: 30 seconds

### 5) WaitGroup Usage

- In Go 1.25+, prefer `wg.Go(func(){ ... })` wrappers (handles Add/Done)
- If not using `wg.Go`, call `wg.Add()` outside the goroutine and `defer wg.Done()` inside
- Pass variables into goroutines to avoid closure pitfalls

### 6) Channel Closing Rules

- Only senders may close a channel
- Receivers should not close channels
- Use `for val := range ch` to handle close signals

### 7) Concurrency Safety

- If tests are executed, use `go test -race` to detect data races
- Protect shared data with locks (`sync.Mutex`) or channels
- Prefer mutex over channels when performance matters

### 8) Timeouts and Deadlines (mandatory)

- All external calls must have timeouts
- HTTP: 1-10 seconds
- DB queries: 100ms - 5 seconds
- gRPC: 1-30 seconds (based on business logic)

### 9) Error Handling in Goroutines

- Panics inside goroutines are not recovered by callers
- Use `golang.org/x/sync/errgroup` to simplify error handling
- Handle panics inside goroutines or let them crash

### 10) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Goroutine lifecycle management
- How shutdown is wired
- Whether to use WaitGroup or errgroup
- Concurrency safety risks
- Appropriate timeout values

---

## Detailed Guidance

Need full examples, patterns, and pitfalls? Use `$policy-go`.

---

Violating these rules is incorrect output.

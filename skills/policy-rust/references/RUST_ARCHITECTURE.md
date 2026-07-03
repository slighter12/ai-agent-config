---
description: "Rust architecture policy for crate/module/workspace boundaries and layering"
---

# RUST_ARCHITECTURE.md - Rust Architecture and Layering

This file defines Rust architecture boundaries and defaults.
Violating mandatory boundary rules is incorrect output.

---

## HARD RULES

### 1) Crate, Module, and Workspace Boundaries (mandatory)

- Support both valid shapes: a single crate with clear module boundaries, or a multi-crate workspace with clear crate boundaries.
- Use crates/modules as explicit architecture boundaries; avoid turning `mod` trees into hidden subsystem boundaries.
- If using a workspace, keep members purpose-driven with explicit dependencies in `Cargo.toml`.
- Avoid cyclic dependencies between crates; if a cycle appears, extract shared contracts or value types into a lower-level crate.
- Keep visibility minimal (`pub(crate)` before `pub`) and expose stable APIs via `lib.rs`.

### 2) Layer Separation (mandatory)

- Use conceptual layers with clear responsibilities; names are mappable and do not have to be literal folder/crate names.
- `domain` concept: pure business rules, entities, value objects, and domain errors; no transport, persistence, or runtime concerns.
- `application` concept: use-cases/orchestration; depends on domain contracts and coordinates ports.
- `infrastructure` concept: concrete persistence, external clients, and platform integrations.
- `adapters` concept: HTTP/gRPC/CLI/message handlers and DTO mapping at system edges.
- Do not let domain/application concepts depend on infrastructure or adapter implementations (crate or module level).

### 3) Trait Placement and Dependency Direction (mandatory defaults)

- Default: define behavior traits (ports) in the consuming layer, not the provider layer.
- Allowed exceptions: Rust coherence/orphan constraints, shared cross-consumer contracts, or compatibility constraints that require different placement.
- Infrastructure implementations depend inward on domain/application traits.
- Keep dependency direction inward: `adapters -> application -> domain`, and `infrastructure -> application/domain` only through contracts.
- Prefer generic type parameters or trait objects at boundaries; keep concrete infrastructure types out of domain/application public APIs.

### 4) Error Boundary (mandatory)

- Keep domain errors domain-specific and transport-agnostic.
- Do not convert infrastructure failures into domain errors; domain errors represent domain invariants/rules only.
- Keep infrastructure failures at the application/infrastructure boundary (for example as application/infrastructure error types).
- Map errors at adapter/transport boundaries into transport responses.
- Do not leak HTTP status codes, SQL driver errors, or SDK-specific error types into domain models.

### 5) Async Runtime Boundary (mandatory)

- Treat runtime choice (Tokio/async-std/etc.) as an infrastructure concern.
- Avoid runtime setup inside domain code.
- Keep async in signatures where I/O or orchestration requires it; do not introduce async-only complexity in pure domain logic.
- Prefer cancellation-aware APIs at adapter/application boundaries for long-running operations.

---

## Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Target layering shape for this repository (single crate vs workspace split)
- Where a new trait/port should live when multiple consumers exist
- Whether trait placement should use an exception path (coherence/orphan/shared contract constraints)
- Whether a boundary may expose external SDK/database/runtime types
- Whether a change intentionally relaxes dependency direction for migration reasons

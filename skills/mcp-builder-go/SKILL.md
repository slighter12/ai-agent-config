---
name: mcp-builder-go
description: Build current-protocol MCP servers in Go with project-fit dependencies.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# MCP Builder Go

Build a Go MCP server whose tool surface matches the agent's work and whose protocol boundary fits its clients and deployment.

## 1. Ground the target

Inspect the repository before choosing an architecture:

- existing `go.mod` dependencies and established package seams;
- supported MCP clients and protocol versions;
- local subprocess, long-running HTTP, serverless, or other deployment constraints;
- authentication, tenancy, persistence, and compatibility requirements.

**Done when:** the target clients, deployment shape, existing stack, and compatibility constraints are known.

## 2. Resolve the protocol contract

Before designing or changing a protocol boundary, read [PROTOCOL-RESEARCH.md](references/PROTOCOL-RESEARCH.md). Use its official-source pass to derive the current lifecycle, transport, capability, authorization, security, and compatibility contract. Keep those findings in the implementation or its tests rather than copying them back into this skill.

**Done when:** every behavior the server will advertise or accept traces to the current official contract and the target clients' requirements.

## 3. Choose the implementation boundary

Choose from evidence rather than imposing a framework:

- Resolve protocol state from the current official contract, separately from process lifetime and deployment topology.
- Use `stdio` when the host launches the server as a local subprocess; otherwise use the transport required by the target host and clients.
- Follow an existing MCP stack when the repository already owns one.
- In greenfield Go, compare the standard library with the official Go SDK and add the SDK only when it materially reduces implementation or maintenance risk.
- Add legacy compatibility only for a confirmed client or migration requirement.

Let the repository decide package names, directories, interfaces, and wiring. Keep protocol parsing, transport concerns, and domain behavior testable without prescribing their physical layout. Build only the MCP entrypoint unless the request explicitly includes an independent CLI.

**Done when:** the transport, dependency strategy, compatibility range, and entrypoint are justified by repository or user evidence.

## 4. Design the tool surface

Expose a small semantic tool surface rather than mirroring internal APIs:

- make inputs explicit, bounded, and discoverable;
- return compact structured outcomes with actionable errors;
- separate read-only discovery from state-changing actions;
- make confirmation, idempotency, pagination, and partial failure visible.

Keep cross-call application state explicit in tool contracts when the resolved protocol and product need it. For mutations, make the target and effect inspectable and return enough structured data to verify the outcome. Keep credentials and authorization material out of results, logs, fixtures, and errors.

**Done when:** every tool corresponds to a concrete agent decision or action, and every mutation, large result, and cross-call state transition has an explicit contract.

## 5. Verify the contract

Derive tests from the resolved protocol matrix and claimed capabilities. Exercise the real entrypoint, representative calls, invalid input, upstream failures, bounded results, mutation safety, cancellation, concurrency, and every compatibility branch the implementation claims.

Use raw transport fixtures for wire-level assertions. Add the official conformance suite, MCP Inspector, or an official SDK client as an interoperability check when available; these complement focused automated tests rather than replacing them.

**Done when:** focused tests pass, representative calls succeed through the real entrypoint, and the reported protocol, transport, dependency, and compatibility choices match the verified behavior.

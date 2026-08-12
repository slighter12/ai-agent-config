# MCP Protocol Research

Use this pass before implementing a new MCP boundary or changing lifecycle, transport, capabilities, authorization, security, or version compatibility. Its output is a current contract for the target implementation, not an update to this reference.

## 1. Pin current authority

Read the official MCP specification landing page, current stable revision, changelog, and normative schema. Record the revision and source links actually consulted. Treat release candidates, proposals, SDK behavior, and ecosystem examples as non-normative unless the stable specification delegates authority to them.

## 2. Derive the contract

Build a compact matrix for only the implementation's claimed surface:

- lifecycle and version negotiation;
- request, response, error, metadata, and capability rules;
- selected transport framing, headers, cancellation, streaming, and connection behavior;
- tools, resources, prompts, subscriptions, caching, tasks, or extensions actually used;
- authorization and security requirements for the deployment boundary;
- deprecated, removed, optional, and legacy behavior relevant to target clients.

Label every row required, optional, deprecated, legacy, or unsupported and cite its normative source. Keep protocol state distinct from process lifetime and deployment topology.

## 3. Check the Go implementation choice

Inspect the official Go SDK documentation, source, release notes, and conformance status as an implementation reference. Compare its supported revisions and behavior with the matrix. Existing project dependencies and conventions remain evidence; the SDK is a dependency only when the chosen design justifies it.

For a standard-library implementation, map each matrix row to owned parsing, validation, error, transport, authorization, and test behavior. Avoid copying SDK-internal abstractions that are not part of the wire contract.

## 4. Close compatibility gaps

Identify every target client's supported revision and transport. Add compatibility behavior only where a confirmed client or migration path requires it. Record unresolved mismatches before implementation rather than silently downgrading or advertising partial capabilities.

**Done when:** the implementation has a cited, revision-specific contract matrix; every claimed capability and compatibility branch has an owner and test; and no protocol conclusion depends solely on memory, a secondary source, or an SDK default.

---
name: mcp-builder-go
description: Design or implement an agent-native MCP server or companion CLI in Go. Use when MCP tool contracts, Go SDK integration, transport behavior, or agent-facing command surfaces are primary. Avoid when the task is ordinary Go application code without an agent tool boundary.
metadata:
  invocation: model
---

# MCP Builder Go

Design from the agent's decision loop:

- expose a small semantic tool surface rather than mirroring internal APIs;
- make inputs explicit, bounded, and discoverable;
- return compact structured outcomes with actionable errors;
- separate read-only discovery from state-changing actions;
- make confirmation, idempotency, pagination, and partial failure visible.

Use the current official MCP specification and official Go SDK documentation at implementation time; do not rely on pinned protocol assumptions in this skill.

Prefer one shared domain service behind MCP and CLI adapters. Keep transport, schema, and business logic separate enough to test without a live client. Validate with a focused protocol smoke test and representative tool calls.

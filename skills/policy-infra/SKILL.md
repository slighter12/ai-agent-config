---
name: policy-infra
description: "Apply infrastructure rules for containerization, secrets/config injection, Docker/Docker Compose practices, and environment variable conventions. Keywords: Docker, Docker Compose, containers, secrets injection, env vars, Kubernetes, build/run stages. Use when runtime/deploy/container configuration decisions are primary. Avoid when only application code changes."
metadata:
  version: "0.1.1"
---

# Policy Guide

Use the referenced policy file for full rules. Keep output aligned with these rules.
Use when deployment/runtime/container topology, build/run staging, ports, or configuration injection mechanisms are primary.
Avoid when deciding secret classification, credential masking, or auth boundaries; use optional `policy-security` detail for secret handling constraints while `policy-infra` owns how config/secrets are injected at runtime.

## Version History

- v0.1.0 (2026-05-08): Initial portable policy release clarifying infrastructure ownership boundaries for runtime/deploy/container configuration and secret/config injection mechanics.
- v0.1.1 (2026-05-29): Align routing wording with optional-depth policy handoffs.

## References

- `references/INDEX.md` - Use for navigation and file selection.
- `references/INFRA.md`

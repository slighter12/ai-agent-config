---
description: "Infrastructure rules - Docker, K8s, environment variables, configuration injection"
---

# INFRA.md - Infrastructure Standards

This file defines mandatory rules for infrastructure, containerization, and configuration.
This policy is domain-scoped and should only be applied when infra trigger conditions are met.
Ownership note: Infra policy owns deployment/runtime mechanics (containers, compose/K8s shape, env/config injection paths, build/run separation). Use when deciding how configuration enters runtime environments. Avoid using infra policy alone to classify secret sensitivity or disclosure rules; use optional security policy detail for secret lifecycle and handling constraints.

Violating these rules is incorrect output.

---

## Table of Contents

- 1) Triggers (when infra updates are required)
- 1) Secrets and Configuration Injection (hard rules)
- 1) Dockerfile Rules (mandatory)
- 1) docker-compose Rules (mandatory)
- 1) Environment Variable Naming
- 1) Kubernetes Compatibility Contract
- 1) Default Value Policy
- 1) Port Management
- 1) Build and Runtime Separation
- 1) Deliverables (hard requirement)
- 1) When Uncertain (mandatory)

## 1) Triggers (when infra updates are required)

If any of the following change, you must evaluate infra impact:

- New service (DB/cache/queue/worker)
- New port or port change
- New environment variables or renames/removals
- New external dependencies (third-party APIs)
- Runtime dependency changes (new build steps, new static assets)

Infra file updates become mandatory only when either condition is true:

- The project already uses containerization, or
- A local reproducible run path is expected and containerization is feasible.

If triggers are not met, or containerization is not feasible, explicitly document `infra not required` with a short reason.

---

## 2) Secrets and Configuration Injection (hard rules)

Mandatory rules:

- Do not put literal secret values into code, images, compose files, or other infra artifacts.
- All runtime configuration must be injected via deployment/runtime mechanisms:
  - Database connection info
  - Cache endpoints
  - External service URLs
  - Credentials, keys, tokens

Ownership boundary:

- Infra policy owns how config/secrets are injected into runtime (mechanism/path), including env/config object wiring and container/orchestrator integration.
- Infra policy does not own secret classification, allowed exposure, retention, rotation, or lifecycle constraints.
- Do not put literal secret values into code, images, compose files, or other infra artifacts.
- Treat credentials, tokens, keys, certificates, and production connection strings as sensitive unless explicitly classified otherwise.
- Secret handling constraints are owned by security policy; use optional `policy-security` detail whenever secrets/tokens/credentials are introduced, changed, or reviewed.

Injection mechanisms:

- Environment variables
- Kubernetes ConfigMap
- Kubernetes Secret
- Docker secrets (Swarm mode)

---

## 3) Dockerfile Rules (mandatory)

Must:

- Target production
- Use multi-stage builds when possible (smaller final image)
- Pin base image versions (no `latest`)
- Run the final image as non-root (unless explicitly required)
- Avoid hardcoded env values or secrets

Example (multi-stage build):

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN adduser -D -u 1000 appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]
```

---

## 4) docker-compose Rules (mandatory)

Must:

- Define clear service names
- Use environment variables for configuration
- Use `.env` placeholders appropriately
- Declare networks explicitly (do not rely on implicit defaults)
- Set restart policies
- Use healthchecks when the project expects them

Do not:

- Embed secrets directly
- Use implicit default networks (large projects)

Example:

```yaml
services:
  app:
    build: .
    ports:
      - "${APP_PORT:-8001}:8001"
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
    networks:
      - backend
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8001/health"]
      interval: 30s
      timeout: 5s
      retries: 3

networks:
  backend:
    driver: bridge
```

---

## 5) Environment Variable Naming

Recommended rules:

- Upper snake case: `DATABASE_URL`, `REDIS_HOST`
- Use prefixes for grouping: `DB_HOST`, `DB_PORT`, `DB_NAME`
- Booleans as `true`/`false` strings

Must document:

- All environment variables in README or `.env.example`
- Purpose, defaults (if any), and example values

---

## 6) Kubernetes Compatibility Contract

Assuming Kubernetes deployment, code and config must:

- Read configuration from environment variables (overridable by ConfigMap/Secret)
- Avoid local file path dependencies (unless volumes are explicitly mounted)
- Support graceful shutdown (handle SIGTERM)
- Provide health endpoints (`/health`, `/ready`)

---

## 7) Default Value Policy

Defaults are allowed only if they are:

- Non-sensitive; credentials, tokens, keys, certificates, and production connection strings are sensitive unless explicitly classified otherwise
- Safe for local development
- Clearly documented

Do not:

- Provide production credentials, tokens, API keys, certificates, or other literal secrets as defaults
- Use defaults that introduce security risk; use optional `policy-security` detail when classification or lifecycle constraints are unclear

Allowed examples:

- `DATABASE_URL=postgres://localhost:5432/dev`
- Default port `PORT=8001`

---

## 8) Port Management

Rules:

- Do not change external ports silently
- When adding a port, you must:
  - Update Dockerfile EXPOSE
  - Update docker-compose ports
  - Update README
  - Update firewall/security group settings if applicable

---

## 9) Build and Runtime Separation

- Build stage must not rely on runtime secrets
- Inject secrets at runtime only
- In multi-stage builds, builder stage should not include runtime secrets

---

## 10) Deliverables (mandatory when feasible)

When infra triggers are active and containerization is feasible (or the project already uses Docker), service implementations or modifications must include:

1. `Dockerfile` (for containerized services)
2. `docker-compose.yaml` (or updates to existing)

If containerization is not feasible, explicitly document the constraints and the alternative deployment/run approach.
Missing these files without explanation is a hard error.
Do not update only part of the deliverables.

---

## 11) When Uncertain (mandatory)

If any of the following are unclear, stop and ask:

- Whether a port is already in use
- Environment variable naming conventions
- Which injection mechanism should be used for secrets/config in this environment
- Secret classification/lifecycle constraints (use optional `policy-security` detail)
- Whether a new external service is required
- Multi-environment strategy (dev/staging/prod)

---

Violating these rules is incorrect output.

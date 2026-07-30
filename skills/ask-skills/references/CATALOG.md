# Skill Catalog

User-invoked skills load their instructions only after explicit invocation. Equivalent plain-language requests may still use the model's default behavior under repository instructions and provider permissions.

## User-invoked

| Skill | Purpose |
| --- | --- |
| `ask-skills` | Explain the catalog and choose a route. |
| `setup-skills` | Install the catalog and initialize project coordination files. |
| `grill-me` | Run a focused requirements interview. |
| `grill-with-docs` | Interview while capturing durable domain context. |
| `to-spec` | Write an implementation-ready spec to the configured tracker. |
| `to-tickets` | Split an approved spec into tracer-bullet tickets and blocking edges. |
| `implement` | Implement an understood change with focused validation. |
| `wayfinder` | Identify current state, blockers, and the next unblocked action. |
| `improve-codebase-architecture` | Make a bounded architecture improvement. |
| `resolving-merge-conflicts` | Complete an active merge or rebase conflict-resolution loop. |
| `handoff` | Create a durable cross-session handoff. |
| `writing-great-skills` | Author, validate, integrate, or package a portable skill. |
| `conventional-git-flow` | Prepare or execute explicitly authorized git and PR actions. |

## Model-invoked

| Skill | Purpose |
| --- | --- |
| `grilling` | Ask one high-value requirements question at a time. |
| `domain-modeling` | Capture domain language, invariants, boundaries, and ADRs. |
| `codebase-design` | Reason about ownership, dependencies, and change seams. |
| `tdd` | Drive a behavior change from a failing executable example. |
| `code-review` | Perform a read-only, comparison-pinned review. |
| `diagnosing-bugs` | Find root cause through a falsifiable evidence loop. |
| `prototype` | Build or plan a deliberately throwaway experiment. |
| `research` | Synthesize current primary sources for a decision. |
| `design-art-direction` | Establish or critique product visual direction. |
| `mcp-builder-go` | Design or implement an agent-native MCP/CLI boundary in Go. |

## Typical Flow

Use only the stages the work needs:

```text
grill-me / grill-with-docs
            |
          to-spec
            |
         to-tickets
            |
         implement
            |
        code-review
            |
 conventional-git-flow
```

`wayfinder` inspects tracker and repository state between stages. `handoff` captures cross-session state. `implement` ends after local verification and handoff; review and git publication remain separate workflows.

---
name: to-spec
description: Turn an agreed direction into an implementation-ready specification in the configured tracker. Use when the user explicitly asks to write or finalize a spec. Avoid when requirements still need discovery or the request is only to implement.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# To Spec

Synthesize the agreed goal, non-goals, current behavior, proposed behavior, constraints, acceptance criteria, risks, and focused verification seams from the conversation and repository evidence.

Write the spec through the tracker declared in `docs/agents/issue-tracker.md`; for a local Markdown tracker, resolve the feature root by that file's precedence rules and write `spec.md` there. Reuse an existing root when the current context identifies it. If the tracker is missing, unconfigured, ambiguous, or collides with a different feature, stop artifact creation and return the exact missing choice or next invocation.

Record only decisions already present in the conversation or repository:

- API work: contracts, validation, versioning, errors, and compatibility.
- Frontend work: product shape, stack, SSR/SEO, deployment, and runtime.
- Security-sensitive work: trust boundaries, sensitive data, abuse cases, and required controls.

Begin artifact creation only when the required decisions are present and consistent. When decisions are missing or inconsistent, complete this workflow by returning the missing-decision list and exactly one next invocation for the user: `grill-with-docs` when the repository should retain durable context, or `grill-me` for a standalone decision interview. The user explicitly invokes that workflow in a later turn; artifact creation resumes after the resolved decisions return.

Completion is the written specification and its exact tracker location. Implementation remains a separate user-invoked workflow.

---
name: implement
description: "Implement a piece of work based on a spec or set of tickets."
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

Implement the work described by the user in the spec or tickets.

Use /tdd where possible, at pre-agreed seams.

Run typechecking regularly, single test files regularly, and the full test suite once at the end.

Once done, use /code-review to review the work.

Commit only when the enclosing user request explicitly authorizes it. Otherwise report the verified, uncommitted state and the exact checks run.

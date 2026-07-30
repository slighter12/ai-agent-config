---
name: diagnosing-bugs
description: Determine the root cause of broken, flaky, slow, or regressed behavior through a tight evidence loop. Use when the user asks to debug, diagnose, explain a failure, or find a performance bottleneck. Avoid when the cause is already known and the request is only to implement the fix.
metadata:
  invocation: model
---

# Diagnosing Bugs

Reproduce or observe the failure with the cheapest reliable signal. Form one falsifiable hypothesis, run the narrowest discriminating check, and update the hypothesis from evidence.

Trace from symptom to violated invariant across inputs, state, dependencies, and recent changes. Separate root cause from contributing conditions and incidental damage.

Stay read-only unless the user also asks for a fix. Report reproduction, evidence, root cause, confidence, and the smallest credible fix direction.

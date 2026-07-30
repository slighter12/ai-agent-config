---
name: tdd
description: Drive a behavior change through a failing executable example and the smallest implementation that passes it. Use when the user asks for TDD, test-first work, regression coverage, or acceptance examples as the main deliverable. Avoid when validation is only a supporting check for ordinary implementation.
metadata:
  invocation: model
---

# TDD

Choose the lowest-cost test boundary that proves the behavior. Write one failing example, confirm it fails for the intended reason, implement the smallest change, then refactor only if the passing evidence exposes duplication.

Use deterministic inputs and observable outcomes. Prefer real collaborators at stable in-process seams; fake only slow or uncontrollable boundaries.

Complete when the selected behavior has failed for the intended reason, passes with the smallest implementation, and its focused checks remain green. Report that evidence while keeping unrelated coverage, mock redesign, and test cleanup outside the scope.

---
name: prototype
description: Build or plan a deliberately throwaway experiment to answer one uncertain product or technical question. Use when the user asks for a spike, proof of concept, or disposable prototype. Avoid when the result must become production code.
metadata:
  invocation: model
---

# Prototype

Name the uncertainty and success signal before building. Set a strict boundary on time, code, integrations, and fidelity.

Use the simplest disposable implementation that can produce evidence, bounded away from production abstractions, migrations, compatibility layers, and broad tests.

Finish when the success signal is observed or the declared boundary is exhausted. Mark the result as throwaway and report what was learned, what remains unknown, and whether to discard, repeat, or proceed to a real spec.

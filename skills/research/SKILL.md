---
name: research
description: Investigate a question against high-trust primary sources and produce a cited Markdown result; write it to the repository only when the enclosing request authorizes artifact creation. Use when the user wants a topic researched, docs or API facts gathered, or reading legwork delegated to a background agent.
metadata:
  invocation: model
---

Run one research pass. When a provider-native background-agent primitive is available, delegate the complete pass to it by default and continue other in-scope work while it runs. Perform the complete pass in the current session only when that primitive is unavailable or the launch fails. Either execution path must satisfy every research requirement below.

Its job:

1. Investigate the question against **primary sources** — official docs, source code, specs, first-party APIs — not a secondary write-up of them. Follow every claim back to the source that owns it.
2. Write the findings as one cited Markdown result.
3. Save it in the repo only when the enclosing request authorizes artifact creation. Otherwise return it inline or in the OS temporary directory, and report the location. When repo writes are authorized, match the existing notes convention; if there is none, put it somewhere sensible and say where.

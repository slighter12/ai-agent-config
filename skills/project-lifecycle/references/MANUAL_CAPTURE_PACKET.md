# Manual Capture Packet

Use this reference when a discussion, loop run, or workflow decision has useful memory before an
automation exists to capture it. The packet is a source-grounded summary, not a transcript.

## When To Use

- A user asks to record source material and discussion outcomes.
- A loop run or recurring workflow needs state for the next run.
- A planning discussion reaches accepted decisions, rejected alternatives, and follow-up owners.
- A workflow lesson may become project docs, project-local skill, shared skill, or no capture.

## Packet Shape

```md
# <Topic> Capture

Date: YYYY-MM-DD
Status: proposed | accepted | superseded

## Source Material

- <source, file, URL, screenshot, or conversation fact>

## Accepted Decisions

- <decision and short rationale>

## Rejected Alternatives

- <alternative and why it was rejected>

## Memory Classification

- Active state: <current phase, owner, accepted progress, rejected evidence, blockers,
  verification status, evidence pointers, and next check>
- Long-lived capture: <decision/implementation pivot/status or documentation drift/capture-worthy handoff note/workflow lesson>
- Not captured: <raw transcript, full logs, provider memory assumptions, unstable idea, duplicate,
  prose-only confidence, etc.>

## Open Questions

- <question or none>

## Follow-Up Owners

- <skill, agent, human, or none>
```

## Rules

- Keep source material factual and compact.
- Do not invent owners, completion status, dates, or product facts.
- Keep raw transcripts and verbose logs out of committed docs by default.
- Treat active state as continuity, not completion proof. Final completion still requires
  authoritative current-state evidence or human review when no objective or calibrated standard
  exists.
- Route accepted long-lived changes through `CAPTURE_GATE.md`.
- Route approved shared skill changes to `skill-creator`; do not edit skills directly from a
  workflow lesson without approval.

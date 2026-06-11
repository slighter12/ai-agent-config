# CLI-Backed Skills

Use a script, CLI, or helper behind a skill when deterministic execution is more reliable than repeated model reasoning.

## Good Candidates

- Repeated packaging, validation, rendering, migration, or data-shaping workflows.
- Workflows with stable inputs, outputs, and error modes.
- Actions where structured output reduces token use or ambiguity.
- Operations that need local filesystem inspection or reproducible command execution.
- Tooling that can provide `--dry-run`, `--json`, or compact output.

## Poor Candidates

- One-off project knowledge.
- Ambiguous design judgment.
- Workflows that require hidden credentials or uncontrolled external side effects.
- Provider-specific behavior that cannot be expressed in portable skill boundaries.
- Helpers that duplicate a mature project-native command.

## Skill Contract

When adding a CLI-backed skill, document:

- Command name and expected installation source.
- Read/write side effects.
- Required credentials or environment variables.
- Stable input shape.
- Stable output shape, preferably JSON for machine consumption.
- Error semantics and expected exit codes when known.
- Manual verification and fallback behavior.

## Boundaries

- Do not add dependencies unless explicitly approved.
- Do not hide mutating behavior behind a read-only skill.
- Keep third-party tool references as source links and distilled guidance; do not vendor external docs or skill text without license review.

## Source Inspiration

Inspired by agent-native CLI design patterns in Printing Press: <https://printingpress.dev/>.
This reference does not copy, vendor, or require Printing Press runtime, skill text, or generated artifacts.

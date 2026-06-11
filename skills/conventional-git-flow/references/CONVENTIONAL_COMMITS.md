# Conventional Commits

Follow Conventional Commits 1.0.0 for commit messages.

## Header Format

Use this repo preference:

```text
<type>: <description>
```

Do not put scope in the header by default. Put the affected area in the branch name and PR body instead.

Examples:

```text
feat: add portable metadata validation
fix: validate compatibility list
docs: document conventional git flow
refactor: simplify skill packaging output
```

## Types

Use the smallest accurate type:

- `feat`: user-visible capability or new reusable behavior.
- `fix`: bug fix or incorrect behavior correction.
- `docs`: documentation-only change.
- `test`: test-only change.
- `refactor`: code restructuring without behavior change.
- `perf`: performance improvement.
- `style`: formatting or style-only change.
- `build`: build system or dependency metadata.
- `ci`: CI configuration or automation.
- `chore`: maintenance with no product or user-facing behavior.

Dependency or version metadata follows the primary logical change when it is required by that feature or fix. Use `build` or `chore` only when the dependency, lockfile, tool, or version metadata change is the whole change.

## Description

- Use imperative mood.
- Start lowercase unless the first token is a proper noun or identifier.
- Keep it concise and specific.
- Do not end with a period.

## Body And Footers

Prefer a concise body when the diff includes multiple meaningful changes, fixes, or behavior notes. Keep the header lightweight, then use the body to capture factual change details that would be lost in the headline.

Omit the body only for very small one-point commits where the headline fully explains the change.

Body guidelines:

- Use short bullet points when there are multiple change details.
- Mention notable fixes, validation behavior, compatibility notes, or risk-reducing details.
- Mention major or ecosystem-specific incompatible dependency/version updates and any compatibility lookup result.
- Keep bullets factual and tied to the diff.
- Do not pad with restatements of the headline.

Example:

```text
fix: validate compatibility list

- Reject unsupported provider names during skill validation.
- Preserve stable validation output for missing metadata.
- Cover empty compatibility entries in the package checks.
```

Use footers for issue references or breaking changes:

```text
BREAKING CHANGE: describe the contract change
Refs: #123
```

Breaking changes may also use `!`, but prefer an explicit `BREAKING CHANGE:` footer when clarity matters.

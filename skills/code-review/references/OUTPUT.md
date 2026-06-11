# Output

Start with the standard review header.

Keep output compact for the default sanity pass. Do not include tables, long command logs, or copied diff snippets unless they are needed to explain a material finding.

## Verdicts

- `reasonable`: no material issue found.
- `needs changes`: one or more material issues should be fixed before continuing.
- `unclear intent`: the diff may be valid, but the intended behavior cannot be safely inferred.

## Findings

Report only material findings:

- Correctness issue.
- Contract or compatibility break.
- Scope creep or unrelated file churn.
- Missing validation for the risk level.
- Documentation, template, or script inconsistency.
- Risky side effect that lacks approval or guardrails.

Do not report style preferences, speculative redesigns, or broad refactors unless they affect correctness.

## Format

Use this structure:

```text
mode: sanity | full | security | architecture-diff-risk
project_profile: ...
verdict: reasonable | needs changes | unclear intent
inferred_intent: ...
findings:
- ...
scope_notes: ...
validation_notes: ...
adversarial_findings: ...
escalation_notes: ...
residual_risks: ...
```

If there are no findings, write `findings: none`.
Include `adversarial_findings` only when the user requested a challenge pass.

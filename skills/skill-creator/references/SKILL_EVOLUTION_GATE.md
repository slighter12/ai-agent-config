# Skill Evolution Gate

Use this gate after skill work is explicit or approved. It is intentionally self-contained so
`skill-creator` can be packaged or installed without sibling skill references.

This gate does not decide whether a project decision, handoff, status update, or workflow lesson
should be captured in the first place. Those are lifecycle questions. This gate starts when the user
asks to create or revise a skill, or when an approved capture candidate has already been routed to
skill authoring.

## Capture Targets

Classify each candidate as one of:

- `reject_or_defer`: not enough evidence, no stable trigger, unclear owner, or not approved for
  skill work.
- `final_answer_note`: useful reminder for the current response, but not worth skill changes.
- `project_local_docs`: project-specific rule for README, AGENTS.md, or equivalent docs instead of
  a skill.
- `project_local_skill`: reusable only inside one project or workspace.
- `shared_skill_update`: cross-project guidance that fits the narrowest existing shared skill.
- `new_shared_skill`: distinct cross-project workflow with its own trigger, avoid cases, output
  contract, and side-effect boundaries.
- `script_or_helper`: deterministic behavior that is safer as code than prompt instructions.

## Entry Conditions

Proceed with skill design only when at least one is true:

- The user explicitly asks to create, revise, package, validate, install, or remove a skill.
- A prior lifecycle or planning step produced an approved `shared_skill_update`, `new_shared_skill`,
  `project_local_skill`, or `script_or_helper` candidate.
- The current task is already inside an approved skill-authoring change.

Reject or defer when the user is still deciding product direction, project status, implementation
approach, handoff content, or whether a workflow lesson is worth keeping.

## Shared Skill Threshold

Only propose a shared `ai-agent-config/skills` update when all are true:

- The workflow is reusable across projects.
- The trigger and avoid cases are stable.
- The narrowest owner skill is clear, or a new skill is justified.
- Side effects, permissions, and provider compatibility are explicit.
- The candidate includes routing examples and manual validation.
- The user explicitly approves writing shared skill files.

Project-local rules should stay local unless there is evidence that the same trigger and avoid cases
apply across repositories.

## Default Preference

Before creating a new skill:

1. Inventory visible skills and their descriptions.
2. Prefer updating the narrowest existing owner.
3. Prefer references for long guidance.
4. Prefer scripts only for deterministic behavior that prevents repeated model reasoning.
5. Keep each skill self-contained: do not require references from sibling skills for normal use,
   packaging, validation, or installation.
6. Reject candidates that are project-specific, unstable, already covered, or not yet validated.

## Self-Contained Packaging Rule

A skill may mention another skill as a routing boundary, but its required workflow references,
scripts, assets, and templates must live inside its own directory unless a provider-specific package
explicitly declares a bundle dependency.

If two skills need similar thresholds, duplicate the minimum needed wording in each skill instead of
linking across sibling skill directories. Shared concepts are acceptable; shared file dependencies
are not the portable default.

Policy reference skills may mention sibling policy names as routing boundaries or optional depth,
but they must not require sibling reference files for normal use outside an explicit package
dependency.

## Rejection Reasons

Use concise reasons:

- `too_project_specific`
- `no_stable_trigger`
- `already_covered`
- `belongs_in_project_docs`
- `not_enough_evidence`
- `not_approved_for_skill_work`
- `side_effect_boundary_unclear`
- `provider_boundary_unclear`
- `owner_skill_unclear`
- `sibling_dependency_required`

Rejected candidates may stay in the final answer as a note, but should not be written into shared skills.

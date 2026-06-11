#!/usr/bin/env python3
"""
Skill initializer for portable SKILL.md-based skills.

Usage:
    init_skill.py <skill-name> --path <path>

Examples:
    init_skill.py my-new-skill --path skills
    init_skill.py my-project-skill --path .agents/skills
    init_skill.py custom-skill --path /custom/location
"""

import re
import sys
from pathlib import Path


SKILL_TEMPLATE = """---
name: {skill_name}
description: Perform [capability]. Use when the user asks for [specific trigger contexts, file types, workflows, or intents]. Avoid when [nearby non-goals or cases better handled by another skill].
metadata:
  version: "0.1.0"
---

# {skill_title}

## Purpose

Enable the agent to [outcome], using only the context and resources needed for this task.

## Use When

- The request involves [trigger 1].
- The user needs [trigger 2].
- The task touches [file type, workflow, domain, or tool].

## Avoid When

- The request is only [nearby non-goal].
- Another skill is more specific: [skill-name].
- The task requires provider-specific behavior not covered by this shared skill.

## Workflow

1. Confirm the goal, inputs, constraints, and expected output.
2. Inspect only the relevant files, docs, or resources.
3. Use bundled scripts only when they improve deterministic reliability or avoid repeated code generation.
4. Load reference files only when their topic is directly needed.
5. Stop when the requested outcome is complete, blocked by missing input, or further work would be speculative.

## Tool And Side-Effect Boundaries

- Prefer read-only inspection until the task clearly requires edits or execution.
- Do not run side-effectful commands unless the user asked for them or the active policy allows them.
- Do not create new dependencies, files, or broad refactors unless the task explicitly requires them.
- For destructive, deployment, commit, notification, credential, or external-service workflows, require explicit user confirmation.

## Output

Return:

- `summary`: what was done or recommended.
- `files_touched`: exact paths, if any.
- `assumptions`: only correctness-relevant assumptions.
- `manual_verification`: checklist when execution was skipped or not required.

## Version History

- v0.1.0 (YYYY-MM-DD): Initial portable skill draft.

## References

- `references/INDEX.md` - Use when deeper topic navigation is needed.
"""

REFERENCE_INDEX = """# References

Add focused reference files here only when details would make SKILL.md too long or too hard to route.
"""


def title_case_skill_name(skill_name):
    """Convert a hyphenated skill name to title case."""
    return " ".join(word.capitalize() for word in skill_name.split("-"))


def validate_skill_name(skill_name):
    """Return an error message when the skill name is invalid."""
    if not re.match(r"^[a-z0-9-]+$", skill_name):
        return "Skill name must use lowercase letters, digits, and hyphens only"
    if skill_name.startswith("-") or skill_name.endswith("-") or "--" in skill_name:
        return "Skill name cannot start/end with hyphen or contain consecutive hyphens"
    if len(skill_name) > 64:
        return "Skill name must be 64 characters or fewer"
    return None


def init_skill(skill_name, path):
    """
    Initialize a new portable skill directory.

    Args:
        skill_name: Name of the skill
        path: Path where the skill directory should be created

    Returns:
        Path to created skill directory, or None if error
    """
    name_error = validate_skill_name(skill_name)
    if name_error:
        print(f"Error: {name_error}")
        return None

    skill_dir = Path(path).resolve() / skill_name
    if skill_dir.exists():
        print(f"Error: Skill directory already exists: {skill_dir}")
        return None

    try:
        skill_dir.mkdir(parents=True, exist_ok=False)
    except Exception as e:
        print(f"Error creating directory: {e}")
        return None

    skill_title = title_case_skill_name(skill_name)
    skill_content = SKILL_TEMPLATE.format(skill_name=skill_name, skill_title=skill_title)

    try:
        (skill_dir / "SKILL.md").write_text(skill_content)
        references_dir = skill_dir / "references"
        references_dir.mkdir()
        (references_dir / "INDEX.md").write_text(REFERENCE_INDEX)
    except Exception as e:
        print(f"Error creating skill files: {e}")
        return None

    print(f"Skill '{skill_name}' initialized at {skill_dir}")
    print("\nNext steps:")
    print("1. Replace bracketed placeholders in SKILL.md.")
    print("2. Replace YYYY-MM-DD in Version History with the actual release date.")
    print("3. Keep provider-specific frontmatter out of the shared template.")
    print("4. Add scripts/ or assets/ only when they materially improve reliability.")
    print("5. Decide and report placement scope: shared/global, project-local, or provider-specific.")
    print("6. Create or report provider surfaces:")
    print("   - shared repo skill: run this repo's ./install.sh when ready;")
    print("   - project-local skill: prefer .agents/skills as source, then verify or report Codex discovery;")
    print("   - link .claude/skills when Claude project discovery is needed;")
    print("   - create .codex/skills only when the project/provider explicitly requires or verifies it.")
    print("7. Run quick_validate.py when ready.")

    return skill_dir


def main():
    if len(sys.argv) < 4 or sys.argv[2] != "--path":
        print("Usage: init_skill.py <skill-name> --path <path>")
        print("\nSkill name requirements:")
        print("  - Hyphen-case identifier, for example 'data-analyzer'")
        print("  - Lowercase letters, digits, and hyphens only")
        print("  - Max 64 characters")
        print("\nExamples:")
        print("  init_skill.py my-new-skill --path skills")
        print("  init_skill.py my-project-skill --path .agents/skills")
        print("  init_skill.py custom-skill --path /custom/location")
        sys.exit(1)

    result = init_skill(sys.argv[1], sys.argv[3])
    sys.exit(0 if result else 1)


if __name__ == "__main__":
    main()

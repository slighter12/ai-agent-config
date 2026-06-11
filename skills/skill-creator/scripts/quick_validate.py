#!/usr/bin/env python3
"""
Quick validation for portable SKILL.md-based skills.
"""

import re
import sys
from pathlib import Path


ALLOWED_PROPERTIES = {"name", "description", "license", "compatibility", "metadata"}
REQUIRED_PROPERTIES = {"name", "description"}
MAX_NAME_LENGTH = 64
MAX_DESCRIPTION_LENGTH = 1024
MAX_COMPATIBILITY_LENGTH = 500
MAX_LICENSE_LENGTH = 200
MAX_SKILL_LINES = 500
MAX_VERSION_LENGTH = 64


def _extract_frontmatter(content):
    match = re.match(r"^---\n(.*?)\n---(?:\n|$)", content, re.DOTALL)
    if not match:
        return None
    return match.group(1)


def _parse_frontmatter(frontmatter_text):
    frontmatter = {}
    current_mapping = None
    for line_number, line in enumerate(frontmatter_text.splitlines(), start=1):
        stripped = line.strip()
        if not stripped:
            continue

        if line.startswith(" "):
            if current_mapping != "metadata":
                raise ValueError(
                    f"Invalid frontmatter line {line_number}: nested keys are only supported under metadata"
                )
            if ":" not in stripped:
                raise ValueError(f"Invalid frontmatter line {line_number}: missing ':'")

            key, value = stripped.split(":", 1)
            key = key.strip()
            value = _parse_scalar_or_inline_list(value.strip())
            if not key:
                raise ValueError(f"Invalid frontmatter line {line_number}: empty metadata key")
            if key in frontmatter["metadata"]:
                raise ValueError(f"Duplicate metadata key: {key}")

            frontmatter["metadata"][key] = value
            continue

        if ":" not in stripped:
            raise ValueError(f"Invalid frontmatter line {line_number}: missing ':'")

        key, value = stripped.split(":", 1)
        key = key.strip()
        value = value.strip()
        current_mapping = None
        if not key:
            raise ValueError(f"Invalid frontmatter line {line_number}: empty key")
        if key in frontmatter:
            raise ValueError(f"Duplicate frontmatter key: {key}")

        if key == "metadata":
            if value:
                raise ValueError("metadata must be a nested mapping")
            frontmatter[key] = {}
            current_mapping = key
        else:
            frontmatter[key] = _parse_scalar_or_inline_list(value)

    return frontmatter


def _parse_scalar_or_inline_list(value):
    value = _unquote_value(value)
    if value.startswith("[") and value.endswith("]"):
        inner = value[1:-1].strip()
        if not inner:
            return []
        return [_unquote_value(item.strip()) for item in inner.split(",")]
    return value


def _unquote_value(value):
    if (value.startswith('"') and value.endswith('"')) or (
        value.startswith("'") and value.endswith("'")
    ):
        return value[1:-1]
    return value


def _validate_description(description):
    if not description:
        return "Description cannot be empty"

    lower_description = description.lower()
    # The final fragment catches the old truncated skill-creator description.
    blocked_fragments = ("todo", "[", "]", "...", "automates the entire s")
    for fragment in blocked_fragments:
        if fragment in lower_description:
            return f"Description contains placeholder or truncated text: {fragment}"

    if "<" in description or ">" in description:
        return "Description cannot contain angle brackets (< or >)"

    if len(description) > MAX_DESCRIPTION_LENGTH:
        return (
            f"Description is too long ({len(description)} characters). "
            f"Maximum is {MAX_DESCRIPTION_LENGTH} characters."
        )

    if "use when" not in lower_description:
        return "Description must include a 'Use when...' trigger boundary"

    if "avoid when" not in lower_description:
        return "Description must include an 'Avoid when...' non-goal boundary"

    return None


def _validate_metadata(metadata):
    if metadata is None:
        return None

    if not isinstance(metadata, dict):
        return "metadata must be a YAML mapping"

    for key, value in metadata.items():
        if not key:
            return "metadata keys cannot be empty"
        if not re.match(r"^[A-Za-z0-9_.-]+$", key):
            return f"metadata key '{key}' should use letters, digits, dots, underscores, or hyphens"
        if not isinstance(value, str):
            return f"metadata.{key} must be a string, got {type(value).__name__}"
        if not value.strip():
            return f"metadata.{key} cannot be empty"
        if "todo" in value.lower() or "[" in value or "]" in value or "..." in value:
            return f"metadata.{key} cannot contain placeholder or truncated text"

    version = metadata.get("version")
    if version is None:
        return None

    if not isinstance(version, str):
        return f"metadata.version must be a string, got {type(version).__name__}"

    version = version.strip()
    if not version:
        return "metadata.version cannot be empty"

    if len(version) > MAX_VERSION_LENGTH:
        return (
            f"metadata.version is too long ({len(version)} characters). "
            f"Maximum is {MAX_VERSION_LENGTH}."
        )

    if "todo" in version.lower() or "[" in version or "]" in version or "..." in version:
        return "metadata.version cannot contain placeholder or truncated text"

    if not re.match(r"^v?\d+(?:\.\d+){0,2}(?:[-+][0-9A-Za-z.-]+)?$", version):
        return "metadata.version should look like a semantic version, for example 0.1.0"

    return None


def _validate_optional_string(frontmatter, key, max_length):
    value = frontmatter.get(key)
    if value is None:
        return None

    if not isinstance(value, str):
        return f"{key} must be a string, got {type(value).__name__}"

    value = value.strip()
    if not value:
        return f"{key} cannot be empty"

    if len(value) > max_length:
        return f"{key} is too long ({len(value)} characters). Maximum is {max_length}."

    if "todo" in value.lower() or "[" in value or "]" in value or "..." in value:
        return f"{key} cannot contain placeholder or truncated text"

    return None


def _validate_compatibility(frontmatter):
    value = frontmatter.get("compatibility")
    if value is None:
        return None

    if isinstance(value, list):
        if not value:
            return "compatibility cannot be an empty list"
        for item in value:
            if not isinstance(item, str):
                return f"compatibility entries must be strings, got {type(item).__name__}"
            item = item.strip()
            if not item:
                return "compatibility entries cannot be empty"
            if len(item) > MAX_COMPATIBILITY_LENGTH:
                return (
                    f"compatibility entry is too long ({len(item)} characters). "
                    f"Maximum is {MAX_COMPATIBILITY_LENGTH}."
                )
            if "todo" in item.lower() or "[" in item or "]" in item or "..." in item:
                return "compatibility entries cannot contain placeholder or truncated text"
        return None

    return _validate_optional_string(frontmatter, "compatibility", MAX_COMPATIBILITY_LENGTH)


def validate_skill(skill_path):
    """Validate one portable skill directory."""
    skill_path = Path(skill_path)
    skill_md = skill_path / "SKILL.md"

    if not skill_md.exists():
        return False, "SKILL.md not found"

    content = skill_md.read_text()
    if not content.startswith("---"):
        return False, "No YAML frontmatter found"

    line_count = len(content.splitlines())
    if line_count > MAX_SKILL_LINES:
        return False, (
            f"SKILL.md is too long ({line_count} lines). "
            f"Move details into references/ and keep SKILL.md under {MAX_SKILL_LINES} lines."
        )

    frontmatter_text = _extract_frontmatter(content)
    if frontmatter_text is None:
        return False, "Invalid frontmatter format"

    try:
        frontmatter = _parse_frontmatter(frontmatter_text)
    except ValueError as e:
        return False, f"Invalid frontmatter: {e}"

    if not isinstance(frontmatter, dict):
        return False, "Frontmatter must be a YAML dictionary"

    unexpected_keys = set(frontmatter.keys()) - ALLOWED_PROPERTIES
    if unexpected_keys:
        return False, (
            f"Unexpected key(s) in portable SKILL.md frontmatter: "
            f"{', '.join(sorted(unexpected_keys))}. "
            f"Allowed properties are: {', '.join(sorted(ALLOWED_PROPERTIES))}"
        )

    for required_key in sorted(REQUIRED_PROPERTIES):
        if required_key not in frontmatter:
            return False, f"Missing '{required_key}' in frontmatter"

    name = frontmatter.get("name", "")
    if not isinstance(name, str):
        return False, f"Name must be a string, got {type(name).__name__}"

    name = name.strip()
    if not re.match(r"^[a-z0-9-]+$", name):
        return False, (
            f"Name '{name}' should be hyphen-case "
            f"(lowercase letters, digits, and hyphens only)"
        )

    if name.startswith("-") or name.endswith("-") or "--" in name:
        return False, f"Name '{name}' cannot start/end with hyphen or contain consecutive hyphens"

    if len(name) > MAX_NAME_LENGTH:
        return False, f"Name is too long ({len(name)} characters). Maximum is {MAX_NAME_LENGTH}."

    description = frontmatter.get("description", "")
    if not isinstance(description, str):
        return False, f"Description must be a string, got {type(description).__name__}"

    description_error = _validate_description(description.strip())
    if description_error:
        return False, description_error

    metadata_error = _validate_metadata(frontmatter.get("metadata"))
    if metadata_error:
        return False, metadata_error

    license_error = _validate_optional_string(frontmatter, "license", MAX_LICENSE_LENGTH)
    if license_error:
        return False, license_error

    compatibility_error = _validate_compatibility(frontmatter)
    if compatibility_error:
        return False, compatibility_error

    return True, "Skill is valid"


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python quick_validate.py <skill_directory>")
        sys.exit(1)

    valid, message = validate_skill(sys.argv[1])
    print(message)
    sys.exit(0 if valid else 1)

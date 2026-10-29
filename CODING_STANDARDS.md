# Coding Standards

Review-time rules for this repository. `AGENTS.md` and `CLAUDE.md` here are the global steering files (installed as symlinks), so repository-only rules live in this file, which the `reviewer` agent and `code-review` read and `install.sh` never installs.

## Trust boundary first

Before reporting a missing check, name who controls the input. For trusted inputs (flags and prose that the calling agent or the user supplies to a skill or its helper script, such as `skills/implement/scripts/run-trail.sh`), report only a mistake a cooperating agent plausibly makes, citing when it happens; an adversarial bypass of a trusted input is not a finding. Untrusted inputs (network, tracker content, contributor code) keep full scrutiny.

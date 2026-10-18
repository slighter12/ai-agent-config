# Domain Documentation Layout

* **Ubiquitous Language Glossary**: `GLOSSARY.md` at the repository root (or routed via `GLOSSARY-MAP.md` in multi-context monorepos). *(Note: Renamed from `CONTEXT.md` in upstream v1.3.1).*
* **Architectural Decision Records (ADRs)**: `docs/adr/NNNN-<slug>.md`.
* **Precedence Rule**: Settled ADRs in `docs/adr/` and terms in `GLOSSARY.md` constitute **Tier 2 Project Law**. When an ADR marks an existing codebase pattern as deprecated, agents must pause at the Pre-Flight Agency Gate and ask the user whether to follow the local legacy file, Sprout at the seam, or migrate the touched module.

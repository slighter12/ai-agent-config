# Design Red Flags

Screen every candidate interface against these flags before comparing or adopting it (see [DESIGN-IT-TWICE.md](DESIGN-IT-TWICE.md)), and when reviewing an existing module. A red flag is a reason to revise or reject the shape, not a style note. Vocabulary follows [SKILL.md](SKILL.md).

## Shallow Module

The interface is large while the implementation hides little. Judge depth by the capability and policy hidden per unit of interface (leverage), not by line counts. A deep module is not a deep call chain: a chain scatters understanding across layers; a deep module concentrates it behind one interface.

- Callers coordinate several methods to complete one operation.
- Public options expose internal stages or implementation choices.
- Learning the interface does not save the caller from learning the implementation.

## Information Leakage

Several modules depend on the same internal decision (a representation, policy, or protocol detail), so changing it needs coordinated edits. Re-exporting transport, wire, or storage types from a module's interface is leakage: parse external data into domain types at the edge, and keep schemas, framework objects, and protocol details behind the seam.

## Temporal Decomposition

Modules are organized by execution order (load, validate, transform, save) instead of the knowledge they own, so one representation and its invariants repeat across stages. Group code around domain knowledge; methods that run at different times still belong together when they protect the same decision.

## Pass-Through Method

A method forwards the same arguments to another method with the same shape. Remove it, or move the responsibility to the module that can complete the operation. Keep a forwarding layer only when it adds policy, adaptation, or a distinct abstraction (apply the Deletion Test in [SKILL.md](SKILL.md) §3).

## Split Ownership

More than one module writes the same state or keeps its own copy of it. An agent that edits one writer cannot see the others, so their rules diverge. Give each piece of state one owner; other modules read it or ask the owner to change it.

## Two Ways to Do One Task

The design supports more than one way to do the same task. An agent copies whichever way it finds first, so every way keeps gaining callers. Keep one way and move callers off the others. When the Agency Gate defers that migration, record the remaining way as legacy in an ADR so new code does not copy it.

## Importable Internals

Callers can import a module's internals. An agent takes the shortest path that compiles, so the internals become part of the interface. Make them unreachable from outside the module (Go `internal/` packages, unexported names, a package `exports` map) so an outside import fails the build.

## Hand-Synced List

Two or more places list the same items, and adding an item means editing every list. An agent that sees one list updates only that one. Keep one list and derive the others from it; when a list cannot be derived, add a check that fails the build when the lists disagree.

# Design It Twice

When designing a critical module interface, exploring alternatives for a chosen deepening candidate, or resolving a Scrap Signal, use this parallel sub-agent pattern. Based on "Design It Twice" (Ousterhout): your first idea is unlikely to be the best.

Uses the vocabulary in [SKILL.md](SKILL.md): **module**, **interface**, **seam**, **adapter**, **leverage**.

## Process

### 1. Frame the problem space

Before spawning sub-agents, write a user-facing explanation of the problem space for the chosen candidate:

- The constraints any new interface would need to satisfy
- The dependencies it would rely on, and which category they fall into (see [DEEPENING.md](DEEPENING.md))
- A rough illustrative code sketch to ground the constraints: not a proposal, just a way to make the constraints concrete

Show this to the user, then immediately proceed to Step 2. The user reads and thinks while the sub-agents work in parallel.

### 2. Spawn sub-agents

Spawn 3+ sub-agents in parallel. Each must produce a **radically different** interface for the deepened module.

Prompt each sub-agent with a separate technical brief (file paths, coupling details, dependency category from [DEEPENING.md](DEEPENING.md), what sits behind the seam). The brief is independent of the user-facing problem-space explanation in Step 1. Give each agent a different design constraint:

- Agent 1: "Minimize the interface: aim for 1–3 entry points max. Maximise leverage per entry point."
- Agent 2: "Maximise flexibility: support many use cases and extension."
- Agent 3: "Optimise for the most common caller: make the default case trivial."
- Agent 4 (if applicable): "Design around ports & adapters for cross-seam dependencies."

Include both [SKILL.md](SKILL.md) vocabulary and `GLOSSARY.md` vocabulary in the brief so each sub-agent names things consistently with the architecture language and the project's domain language.

Each sub-agent outputs:

1. Caller usage first: a snippet showing how the most common caller uses it
2. Data shape: the core types, state machine, or registry, with illegal states unrepresentable
3. Interface (types, methods, params, plus invariants, ordering, error modes), derived from the usage
4. What the implementation hides behind the seam
5. Dependency strategy and adapters (see [DEEPENING.md](DEEPENING.md))
6. Trade-offs, each phrased as "we accept X in exchange for Y": where leverage is high, where it's thin

### 3. Present and compare

Screen each design against [DESIGN-RED-FLAGS.md](DESIGN-RED-FLAGS.md) first; revise or drop a design that raises a flag, and say which flag.

Present designs sequentially so the user can absorb each one, then compare them in prose. Contrast by **depth** (leverage at the interface), **locality** (where change concentrates), and **seam placement**, using these concrete checks:

- **Caller leverage**: how many lines of setup and error handling does every caller have to write?
- **Seam simplicity**: can behavior be verified at the public interface with `<= 2` test arguments/fakes?
- **Local legibility**: if a future agent copies one call site, can it misuse the ordering or state invariants?

After comparing, give your own recommendation: which design you think is strongest and why. If elements from different designs would combine well, propose a hybrid. Be opinionated: the user wants a strong read, not a menu. Present the recommendation before committing to the implementation.
Name at least one losing alternative with one line on why it lost, judged on the complexity it exposes to callers rather than implementation simplicity. Phrase open risks as questions the user can answer, not assertions.

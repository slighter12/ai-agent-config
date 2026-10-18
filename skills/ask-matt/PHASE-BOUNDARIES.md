# Phase Boundaries & Context Management

At the boundary between two phases (never mid-phase), evaluate these five options top to bottom — the first `yes` wins:

1. **Continue in the current window**: Does the next phase benefit from this conversation as a primary source, and are you still inside the Smart Zone? If yes, stay put (zero loss, zero overhead).
2. **`/clear`**: Is everything in the current window disposable (e.g., finished an independent ticket commit)? If yes, `/clear`.
3. **`/handoff`**: Does the context need to travel across harnesses (Claude <-> Codex), to a new directory/worktree, or to a forked side quest? If yes, run `/handoff`.
4. **Subagent**: Is the task a noisy read-only exploration (`explorer`, `librarian`, `oracle`), diff review (`reviewer`), or isolated ticket build (`builder`) that can run in its own window and return a summary? If yes, spawn a subagent.
5. **`/compact`**: Relevant context, same harness, same directory, and approaching the Smart Zone ceiling. Pass a focus instruction (`/compact <what matters next>`).

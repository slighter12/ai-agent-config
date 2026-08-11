---
name: wizard
description: Generate an interactive bash wizard that walks a human through steps only they can perform. Use when provisioning infrastructure, setting up credentials or CI secrets, walking an unfamiliar third-party dashboard, or running a one-off migration or cutover. Don't invoke this for steps the agent can perform itself.
metadata:
  invocation: model
---

# Wizard

A **wizard** is a bash script that walks a human, step by step, through a manual procedure that's tedious to do by hand and tedious to re-explain to an AI every time. It opens each URL, says exactly what to click and copy, captures the values, writes them where they belong (`.env`, GitHub secrets), confirms at every stage, and shows how many stages are left. It might configure third-party services, run a one-off migration, or move the project from one state to another.

The delightful UX is already solved by [template.sh](template.sh) — stage-by-stage progress, confirmation gates, cross-platform URL opening (including WSL), hidden secret entry, idempotent `.env` upserts, `gh secret`/`gh variable` writes, and a closing summary. **Your job is only to scope the procedure and author its stages.** The library above the `STAGES` marker is identical in every wizard; that consistency is the point — never hand-edit it.

A wizard is ephemeral by default — built for one run, saved to a scratch or `scripts/` path, deleted when the job's done. Commit it only when the user wants a repeatable setup path that should live in the repo.

## Process

### 1. Scope the procedure

Work out every manual step the human must take and every value that gets captured along the way. Read the repo first — don't ask cold:

- For setup: `.env`, `.env.example`, `.env.*`, `README`, `docker-compose*`, framework config, and `.github/workflows/*` (every `secrets.*` / `vars.*` reference is a value the wizard must produce).
- For a migration or transition: the current state, the target state, and the irreversible actions between them.

Before any `_existing` read, file copy, or write, the library ignores ambient `ENV_FILE`, then resolves the exact literal path authored into the stages against the physical startup directory captured by `pwd -P`. It accepts only a target below that fixed base with an existing real parent chain and either an absent leaf or an existing regular file. Reject absolute escapes, `.` / `..` components, newlines, symlinked parents, a symlink (including a dangling symlink) leaf, and every other non-regular leaf without reading or replacing it. Use the one validated absolute target for reads, temporary files, copies, replacement, permissions, and cleanup, and revalidate it before each important filesystem operation. The Git checks clear inherited repository, worktree, index, object, discovery, and config overrides before selecting the nearest owning worktree from the validated target's parent. There the library requires the final target and `${target}.tmp.probe` to match ignore rules during preflight, then separately verifies the actual random temporary path immediately after `mktemp` and before writing any value.

The raw dotenv serializer accepts keys matching `^[A-Za-z_][A-Za-z0-9_]*$` and values that are empty or contain only `A-Za-z0-9_./:@%+,=?-`. Capture helpers reject shell/library/security names such as `PATH`, `ENV_FILE`, `GITHUB_REPO`, `IFS`, `CDPATH`, `BASH*`, `BASH_ENV`, `SHELLOPTS`, and `_WIZARD_*`; use ordinary project-specific keys. Validate both new input and `_existing` values before creating a temporary file; a project needing spaces, quotes, shell expansion, comments, control operators, escapes, or multiline values must provide its own serializer. The most recent capture is also available through the fixed `_WIZARD_CAPTURED_VALUE` channel; the historical `ask KEY` and `ask_secret KEY` assignments remain available for non-reserved keys.

Then show the user the ordered list of stages and the values each produces, and confirm — they may add, drop, or reorder.

That confirmation authorizes writing the wizard when the enclosing request permits local edits. It does not authorize running the wizard, handling credentials, or changing external systems; those actions require their own explicit request.

**Done when:** every stage is named in order, and for each captured value you know (a) where the human gets it, (b) where it's written (`.env`, a GitHub secret, both, or nowhere — some stages are pure actions), (c) whether it's secret (hidden entry) or public, and (d) when the target is in a Git worktree, both the exact `ENV_FILE` path and its temporary probe path are ignored. Verify each path separately with `git check-ignore -- "$ENV_FILE"` and `git check-ignore -- "${ENV_FILE}.tmp.probe"`; if either check fails, stop before writing sensitive values and obtain authorization to add an ignore rule. The helper's scoped EXIT/INT/TERM traps clean up temporary files on ordinary exits and handled interrupts, while the ignore checks provide defense in depth for SIGKILL or power loss.

### 2. Map each stage's journey

For each stage, write the precise path a human follows: which URL to open, what to do there, where a value is shown, which variable it fills — e.g. "Dashboard → Developers → API keys → Reveal test key → copy". Where you don't actually know the current UI or the exact command, say so and ask the user or check the docs — never invent steps that may not exist.

**Done when:** every stage traces to concrete instructions a stranger could follow.

### 3. Author the wizard

Copy `template.sh` to the target path. Replace the example stage with one `stage` per step, in dependency order. Use the library helpers — `stage`, `say`/`step`, `open_url`, `ask`/`ask_secret`, `write_env`, `set_secret`/`set_var`, `pause`/`confirm` — and set `TOTAL_STAGES` to the number of stages you wrote. When the approved target is not `.env`, assign `ENV_FILE` to that exact literal in the stages section; never derive it from process environment. When any stage uses `set_secret` or `set_var`, set `GITHUB_REPO` in the stages section to the exact `OWNER/REPO` or `HOST/OWNER/REPO`; the helpers never infer a write target from the current directory.

Hold the bar the template sets: open the URL before asking for its value, use `ask_secret` for anything secret, `write_env` every persisted value, and `set_secret` only for values CI actually needs. The GitHub helpers display the exact repository and key, then require their own default-No confirmation before every write. Use `confirm` for every other irreversible action. Each `stage` clears the screen so only the current step is visible — keep a stage to one focused task so nothing the human needs scrolls away. Don't touch the library above the marker.

### 4. Verify and hand off

- `bash -n <script>`; run `shellcheck` if available.
- `chmod +x <script>`.
- Don't run it end-to-end yourself — it opens browsers and blocks on human input. Trace it statically instead: every value from step 1 is captured and lands where step 1 said, and every `set_secret` name exactly matches a `secrets.*` reference in CI.
- Run the focused integration checks for path containment, parent and leaf symlinks, dotenv validation, fixed startup-directory anchoring, secret non-exposure, temporary cleanup, and `0600` permissions. Keep the library Bash 3.2-compatible. Without a native `openat`-style helper, validation and use are separate operations; a same-user process that can mutate the validated path concurrently remains a TOCTOU residual risk.
- Tell the user how to run it. If it's a repeatable setup path, propose committing it and linking it from the README; perform those git and documentation changes only when explicitly authorized.

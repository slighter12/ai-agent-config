---
name: verify-app
description: Create or maintain a project-local verification skill that launches the real app, drives it the way a user does, and captures evidence that survives cleanup. Use when a project has no scripted way to prove UI, CLI, or service behavior, or when its verification skill has drifted.
disable-model-invocation: true
argument-hint: "create or maintain, and which app?"
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Verify App

Give a project one scripted way to prove behavior in the running app: launch it, exercise a feature as a user would, and capture evidence. `tdd` proves behavior at a code seam; `retro` notices missing recipes after a session; `verify-app` builds or repairs the recipe and proves it by running it. Write the generated skill for an agent that will read it cold, mid-task, without having seen the app.

## Boundaries

- Invoking this skill authorizes edits only inside the target verification skill's directory. Never edit product code. If the checkout does not build or start as-is, stop and report the exact command and failure.
- Write to the provider skill root the repository already uses (`.claude/skills/` or `.agents/skills/`). If neither exists, ask which one to create.
- Run only local instances on disposable data directories, profiles, and ports. Never drive the user's own running instance or real data. External traffic, real credentials, or account changes need an exact preview and approval; store credential names, never values, in the generated skill.
- Stop only processes you started, tracked by PID or session; never kill by process name.
- Leave changes uncommitted. Branches, commits, and PRs need their own approval.
- Pick the mode from the request: **create** when the project has no verification skill, **maintain** when one exists. If several candidates exist, ask which one.

## Create

### 1. Interview the Repository

Answer from the codebase and ask the user only what you cannot observe:

- **Surface**: what a user touches (web UI, CLI or TUI, desktop, API, mobile, library). Pick the primary one and note the rest.
- **Run**: the repository's documented dev command, with ports, environment variable names, seed data, and auth.
- **Drive**: existing harnesses first (browser specs, PTY or expect scripts, HTTP endpoints, a debug port), then a generic recipe: browser automation for web, a PTY or tmux session for CLI and TUI, plain HTTP for services.
- **Observe**: screenshots, terminal transcripts, response bodies, logs, exit codes, and database state.
- **Isolate**: whether two instances can run side by side. If not, the generated skill must refuse to drive a shared instance.

### 2. Write the Skill

Create `verify-<app>/SKILL.md` with frontmatter (`name: verify-<app>` and a `description` naming the app, the surface, and when to use it), following the project's existing skill conventions. Leave it model-invocable so agents proving behavior in this project load it instead of skipping or guessing, unless the user asks for manual-only (`disable-model-invocation: true`). Ground every section in what the interview found; leave no placeholders:

- **Launch**: the exact command that starts the app from an empty data directory, including migrations and seed data, and the signal that it is ready. Ready means usable, not merely answering: a health check can pass while init scripts still run. For a short-lived CLI, Launch builds once and each drive gets its own isolated session.
- **Doctor**: one read-only check that the instance is worth driving: process up, expected build, port owned by this run, auth valid. Run it first and after anything surprising.
- **Drive**: real selectors and commands from this repository. Prefer stable handles (accessible labels, data attributes, routes, prompt strings) over coordinates and tab order.
- **Evidence**: what to capture and where it lands. Exercise the real user path, not internal setters or test-only endpoints; capture the action and the resulting state; check side effects (files, rows, messages) as well as what is visible; mock only where a production boundary already isolates the external system. When the safe path is a dry run, observe what it actually skips rather than trusting its name. Redact secrets.
- **Cleanup**: tear down what the run created and keep the evidence at the named location.
- **Helpers**: any script the skill ships is executable, and the skill body shows how to invoke it.

### 3. Seed the Feature Map

Create `verify-<app>/features/README.md` as an index and one file per user-facing feature (start with the top 3–5 from routes, commands, menus, or docs). Each file uses four H2s: `Sub-features`, `How to get to it (user POV)`, `Driving it with <harness>` (preconditions, then each step with its command and the observable result), and `Gotchas`. If seed data lacks a precondition a journey needs, the journey fails and the reason goes in `Gotchas`. Never weaken a journey to make it pass.

### 4. Prove It Before Handing Over

Run the generated skill end to end once: Launch, Doctor, drive one mapped feature, capture evidence, Cleanup. After cleanup, confirm the evidence still exists at its named location. Fix what fails, and run Cleanup after every failed attempt. A skill that was never run is a draft. If you proved it in a stand-in for the real target (a container instead of a fresh VM), list each setup step the stand-in provided that the target will not; those stay Not verified.

## Maintain

Keep the skill and its feature map honest as the app changes. The unit of rigor is the feature: every feature file gets source coverage and a live drive.

1. **Index hygiene**: fix missing, extra, duplicate, or dead entries in `features/README.md`.
2. **Source pass**: for each feature file, read the source that implements it, flag drift with `file:line` citations, and note one live recipe. Work in this session; dispatch read-only `explorer` subagents only when raw output would bloat the window. They never drive the app or edit files.
3. **Reconcile**: merge recipes into as few app states as practical, spot-check cited drift, and look at recent churn for user-facing surfaces the map omits, naming a concrete source path for each.
4. **Live pass**: required even when the source looks clean. Follow the skill's own launch model and drive every feature at least once. Run Doctor before the first drive and after any failed drive; when Doctor cannot see the failure, reset or relaunch. Evidence must survive every cleanup. A feature is unreachable only with its concrete missing prerequisite and the route attempted.
5. **Triage**: a wrong description is doc drift, so fix the map; working behavior the harness cannot drive is a harness gap, so fix the harness and re-drive it live; broken app behavior is a product gap, so report it and leave the map unchanged.
6. **Outcome**: report exactly one: **clean** (full coverage, nothing to change), **changed** (proven corrections left as local edits for review), or **blocked** (what stopped coverage or a safe fix).

#!/usr/bin/env bash
# Human-in-the-loop reproduction loop.
# Copy this file, edit the steps below, and run it.
# The agent runs the script; the user follows prompts in their terminal.
#
# Usage:
#   bash hitl-loop.template.sh
#
# Two helpers:
#   step "<instruction>"          → show instruction, wait for Enter
#   capture VAR "<question>"      → show question, read response into VAR
#
# At the end, captured values are printed as shell-escaped KEY=VALUE pairs for
# the agent to parse. Never paste secrets, tokens, cookies, credentials, or PII.
#
# `capture` keeps the response hidden while it is entered and prints a shell-
# escaped value at the end for the agent to read — so capture observations, and
# leave signing in to the user as a `step`.

set -euo pipefail

step() {
  printf '\n>>> %s\n' "$1"
  read -r -p "    [Enter when done] " _
}

capture() {
  local var="$1" question="$2" answer
  printf '\n>>> %s\n' "$question"
  read -rs -p "    > " answer
  printf '\n'
  printf -v "$var" '%s' "$answer"
}

# --- edit below ---------------------------------------------------------

step "Open the app at http://localhost:3000 and sign in."

capture ERRORED "Click the 'Export' button. Did it throw an error? (y/n)"

capture ERROR_MSG "Paste the redacted error message (replace secrets, tokens, cookies, credentials, and PII with <REDACTED>; or 'none'):"

# --- edit above ---------------------------------------------------------

printf '\n--- Captured ---\n'
printf 'ERRORED=%q\n' "$ERRORED"
printf 'ERROR_MSG=%q\n' "$ERROR_MSG"

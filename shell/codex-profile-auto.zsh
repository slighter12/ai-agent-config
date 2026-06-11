# Repo-managed Codex auto-profile wrapper.
# Source from ~/.zshrc through ~/.codex/shell/codex-profile-auto.zsh.

_codex_auto_profile_has_explicit_profile() {
  emulate -L zsh

  local arg
  for arg in "$@"; do
    case "$arg" in
      --profile|--profile=*|--profile-v2|--profile-v2=*|-p|-p*)
        return 0
        ;;
    esac
  done

  return 1
}

_codex_auto_profile_workspace() {
  emulate -L zsh

  local arg next_takes_cd=0
  for arg in "$@"; do
    if (( next_takes_cd )); then
      print -r -- "$arg"
      return 0
    fi

    case "$arg" in
      -C|--cd)
        next_takes_cd=1
        ;;
      --cd=*)
        print -r -- "${arg#--cd=}"
        return 0
        ;;
      -C?*)
        print -r -- "${arg#-C}"
        return 0
        ;;
    esac
  done

  print -r -- "$PWD"
}

_codex_auto_profile_non_option_args() {
  emulate -L zsh

  local arg skip_next=0
  local -a value_options
  value_options=(
    -c --config
    -i --image
    -m --model
    --local-provider
    -p --profile
    --profile-v2
    -s --sandbox
    -C --cd
    --add-dir
    -a --ask-for-approval
    --remote
    --remote-auth-token-env
    --enable
    --disable
  )

  for arg in "$@"; do
    if (( skip_next )); then
      skip_next=0
      continue
    fi

    if (( ${value_options[(Ie)$arg]} )); then
      skip_next=1
      continue
    fi

    case "$arg" in
      --)
        shift
        print -r -- "$@"
        return 0
        ;;
      --*=*)
        continue
        ;;
      -*)
        continue
        ;;
      *)
        print -r -- "$arg"
        ;;
    esac
  done
}

_codex_auto_profile_should_apply() {
  emulate -L zsh

  if _codex_auto_profile_has_explicit_profile "$@"; then
    return 1
  fi

  local arg
  for arg in "$@"; do
    case "$arg" in
      -h|--help|-V|--version)
        return 1
        ;;
    esac
  done

  local -a non_options
  non_options=("${(@f)$(_codex_auto_profile_non_option_args "$@")}")

  local command_name="${non_options[1]:-}"
  local debug_command="${non_options[2]:-}"

  case "$command_name" in
    ""|exec|review|resume|fork)
      return 0
      ;;
    debug)
      [[ "$debug_command" == "prompt-input" ]]
      return $?
      ;;
    app|app-server|apply|cloud|completion|doctor|exec-server|features|help|login|logout|mcp|mcp-server|plugin|remote-control|sandbox|update)
      return 1
      ;;
    *)
      return 0
      ;;
  esac
}

_codex_auto_profile_binary() {
  emulate -L zsh

  if [[ -n "${CODEX_CLI_PATH:-}" && -x "$CODEX_CLI_PATH" ]]; then
    print -r -- "$CODEX_CLI_PATH"
    return 0
  fi

  if [[ -x "/Applications/Codex.app/Contents/Resources/codex" ]]; then
    print -r -- "/Applications/Codex.app/Contents/Resources/codex"
    return 0
  fi

  print -r -- "codex"
}

_codex_auto_profile_name() {
  emulate -L zsh

  local workspace="$1"
  local project_root="$workspace"

  if command -v git >/dev/null 2>&1 && [[ -d "$workspace" ]]; then
    local git_root
    git_root="$(git -C "$workspace" rev-parse --show-toplevel 2>/dev/null || true)"
    if [[ -n "$git_root" ]]; then
      project_root="$git_root"
    fi
  fi

  local project_name="${project_root:t}"
  local profile_name="${project_name:l}"
  profile_name="${profile_name// /-}"
  profile_name="${profile_name//_/-}"

  print -r -- "$profile_name"
}

codex() {
  emulate -L zsh

  local codex_bin
  codex_bin="$(_codex_auto_profile_binary)"

  if _codex_auto_profile_should_apply "$@"; then
    local workspace profile_name profile_path
    workspace="$(_codex_auto_profile_workspace "$@")"
    profile_name="$(_codex_auto_profile_name "$workspace")"
    profile_path="${CODEX_HOME:-$HOME/.codex}/${profile_name}.config.toml"

    if [[ -r "$profile_path" ]]; then
      if [[ "$codex_bin" == /* ]]; then
        "$codex_bin" --profile "$profile_name" "$@"
      else
        command "$codex_bin" --profile "$profile_name" "$@"
      fi
      return $?
    fi
  fi

  if [[ "$codex_bin" == /* ]]; then
    "$codex_bin" "$@"
  else
    command "$codex_bin" "$@"
  fi
}

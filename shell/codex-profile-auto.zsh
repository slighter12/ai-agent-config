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

_codex_auto_profile_load_private_env() {
  emulate -L zsh

  if [[ -z "${CLIPROXYAPI_API_KEY:-}" ]]; then
    local key_path="${CODEX_HOME:-$HOME/.codex}/secrets/cliproxyapi_api_key"
    if [[ -r "$key_path" ]]; then
      local key
      key="$(IFS= read -r line < "$key_path"; print -r -- "$line")"
      if [[ -n "$key" ]]; then
        export CLIPROXYAPI_API_KEY="$key"
      fi
    fi
  fi
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

_codex_auto_profile_dev_add_dir_args() {
  emulate -L zsh

  [[ "${CODEX_AUTO_DEV_DIRS:-1}" == "0" ]] && return 0

  local -A seen
  local home="${HOME:-}"

  _codex_auto_profile_add_dir() {
    emulate -L zsh

    local path="$1"
    local create="${2:-}"
    [[ -n "$path" ]] || return 0
    [[ -n "$home" && "$path" == "~/"* ]] && path="$home/${path#~/}"
    [[ "$path" == /* ]] || return 0
    [[ "$path" != "/" && "$path" != "$home" ]] || return 0
    if [[ ! -d "$path" ]]; then
      [[ "$create" == "create" ]] || return 0
      [[ -d "${path:h}" ]] || return 0
      mkdir -p "$path" 2>/dev/null || return 0
    fi
    [[ -z "${seen[$path]:-}" ]] || return 0

    seen[$path]=1
    print -r -- "--add-dir"
    print -r -- "$path"
  }

  if command -v go >/dev/null 2>&1; then
    local -a go_env
    go_env=("${(@f)$(go env GOCACHE GOMODCACHE GOPATH GOBIN 2>/dev/null)}")
    _codex_auto_profile_add_dir "${go_env[1]:-}"
    _codex_auto_profile_add_dir "${go_env[2]:-}"
    if [[ -n "${go_env[4]:-}" ]]; then
      _codex_auto_profile_add_dir "${go_env[4]}"
    elif [[ -n "${go_env[3]:-}" ]]; then
      local gopath
      for gopath in "${(@ps.:.)go_env[3]}"; do
        _codex_auto_profile_add_dir "$gopath/bin"
      done
    fi
  fi

  if command -v uv >/dev/null 2>&1; then
    _codex_auto_profile_add_dir "$(uv cache dir 2>/dev/null)"
    _codex_auto_profile_add_dir "$(uv python dir 2>/dev/null)"
  fi

  if command -v bun >/dev/null 2>&1; then
    local bun_home="${BUN_INSTALL:-$home/.bun}"
    _codex_auto_profile_add_dir "$bun_home/install/cache"
    _codex_auto_profile_add_dir "$bun_home/install/global"
    _codex_auto_profile_add_dir "$bun_home/bin"
  fi

  if command -v cargo >/dev/null 2>&1 || command -v rustc >/dev/null 2>&1; then
    local cargo_home="${CARGO_HOME:-$home/.cargo}"
    _codex_auto_profile_add_dir "$cargo_home/registry"
    _codex_auto_profile_add_dir "$cargo_home/git"
    _codex_auto_profile_add_dir "$cargo_home/bin"
  fi
}

codex() {
  emulate -L zsh

  _codex_auto_profile_load_private_env

  local codex_bin
  codex_bin="$(_codex_auto_profile_binary)"

  if _codex_auto_profile_should_apply "$@"; then
    local workspace profile_name profile_path
    local -a dev_add_dir_args
    workspace="$(_codex_auto_profile_workspace "$@")"
    profile_name="$(_codex_auto_profile_name "$workspace")"
    profile_path="${CODEX_HOME:-$HOME/.codex}/${profile_name}.config.toml"
    dev_add_dir_args=("${(@f)$(_codex_auto_profile_dev_add_dir_args)}")

    if [[ -r "$profile_path" ]]; then
      if [[ "$codex_bin" == /* ]]; then
        "$codex_bin" --profile "$profile_name" "${dev_add_dir_args[@]}" "$@"
      else
        command "$codex_bin" --profile "$profile_name" "${dev_add_dir_args[@]}" "$@"
      fi
      return $?
    fi

    if [[ "$codex_bin" == /* ]]; then
      "$codex_bin" "${dev_add_dir_args[@]}" "$@"
    else
      command "$codex_bin" "${dev_add_dir_args[@]}" "$@"
    fi
    return $?
  fi

  if [[ "$codex_bin" == /* ]]; then
    "$codex_bin" "$@"
  else
    command "$codex_bin" "$@"
  fi
}

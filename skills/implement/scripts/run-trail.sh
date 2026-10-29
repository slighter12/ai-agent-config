#!/bin/sh
# run-trail.sh: the decision trail and review base of one implement run.
#
# The trail lives in ${TMPDIR:-/tmp}/run-trail/<repo-key>/<slug>/, outside the
# project and every worktree. Output is key=value lines on stdout.
# Exit codes: 0 ok; 2 rejected call (stderr says why); 3 needs the user (the
# first stdout line is needs_user=<reason>); 1 unexpected git or system failure.
#
# The script reads git state and writes only its own trail directory: it never
# commits, merges, rebases, or deletes anything.

set -euf
umask 077

PROG=run-trail
NL='
'
HEADER=$(printf 'ts\tphase\tdecision\twhy\tevidence\tresult')

usage() {
	cat >&2 <<'EOF'
usage: run-trail.sh [--dir DIR] <command> [options]

  init      --slug S --mode A|B|fallback [--integration BRANCH] [--parent BRANCH|none]
            [--decision TEXT] [--why TEXT] [--fresh]
  log       --slug S --phase P --decision TEXT --why TEXT --result R
            [--ticket ID] [--commit REV] [--check TEXT] [--evidence TEXT]
            [--deferred-to ID] [--supersedes TS]
  status    --slug S
  base      --slug S [--target REF | --wip [--committed REV]]
  scope     --allow PATTERN... (--staged | --since REV)
  worktrees
EOF
	exit 2
}

die() {
	printf '%s: error: %s\n' "$PROG" "$*" >&2
	exit 2
}

fail() {
	printf '%s: error: %s\n' "$PROG" "$*" >&2
	exit 1
}

# needs_user REASON [detail...] prints the reason first, then one detail per
# line, and exits 3.
needs_user() {
	printf 'needs_user=%s\n' "$1"
	shift
	for detail in "$@"; do
		[ -z "$detail" ] || printf '%s\n' "$detail"
	done
	exit 3
}

g() {
	git -c core.quotePath=false --no-pager "$@"
}

now() {
	date -u +%Y-%m-%dT%H:%M:%SZ
}

# full_commit REV prints the full commit ID of one revision. Assign its output
# on a line of its own, so a rejected revision stops the script.
full_commit() {
	case $1 in
	'') die "empty revision" ;;
	-*) die "a revision never starts with '-': $1" ;;
	*..*) die "evidence is one commit, never a revision range: $1" ;;
	esac
	g rev-parse --verify --quiet "$1^{commit}" || die "not a commit: $1"
}

has_commit() {
	g rev-parse --verify --quiet "$1^{commit}" >/dev/null 2>&1
}

# clean_cell turns tabs and newlines into spaces.
clean_cell() {
	printf '%s' "$1" | tr '\t\r\n' '   '
}

valid_id() {
	case $1 in
	'' | . | .. | -* | *[!A-Za-z0-9._-]*) return 1 ;;
	esac
}

require_repo() {
	g rev-parse --git-dir >/dev/null 2>&1 || die "not inside a git repository"
}

repo_key() {
	common=$(g rev-parse --git-common-dir) || fail "cannot read the git common directory"
	common_abs=$(cd "$common" && pwd -P) || fail "cannot resolve $common"
	printf '%s' "$common_abs" | cksum | awk '{print $1}'
}

# resolve_trail sets KEY, TRAIL, TSV, and STATE from --dir or from the slug;
# "create" makes the slug's directory.
resolve_trail() {
	KEY=$(repo_key)
	if [ -n "$DIR" ]; then
		[ -d "$DIR" ] || die "--dir is not a directory: $DIR"
		TRAIL=${DIR%/}
	else
		valid_id "$SLUG" || die "--slug takes letters, digits, '.', '_', and '-': '$SLUG'"
		root=${TMPDIR:-/tmp}
		TRAIL=${root%/}/run-trail/$KEY/$SLUG
		[ "$1" != create ] || mkdir -p "$TRAIL" || fail "cannot create $TRAIL"
	fi
	TSV=$TRAIL/decisions.tsv
	STATE=$TRAIL/state
}

state_get() {
	sed -n "s/^$1=//p" "$STATE" | head -n 1
}

require_trail() {
	resolve_trail open
	[ -f "$STATE" ] && [ -f "$TSV" ] || die "no run trail at $TRAIL; run init first, or pass the trail directory with --dir"
	[ "$(state_get repo_key)" = "$KEY" ] || die "the trail at $TRAIL belongs to another repository"
	MODE=$(state_get mode)
	CAPTURE=$(state_get capture)
	PARENT=$(state_get parent)
	INTEGRATION=$(state_get integration)
}

lock() {
	tries=0
	while ! mkdir "$TRAIL/lock" 2>/dev/null; do
		tries=$((tries + 1))
		[ "$tries" -lt 10 ] || needs_user trail-locked "lock=$TRAIL/lock" \
			"hint=another run-trail command holds the lock: stop and ask the user; never remove the lock yourself"
		sleep 1
	done
	trap 'rmdir "$TRAIL/lock" 2>/dev/null' EXIT
	trap 'exit 1' HUP INT TERM
}

unlock() {
	rmdir "$TRAIL/lock" 2>/dev/null || true
	trap - EXIT
}

# rows prints the trail rows, without the header, minus every row a pivot
# row supersedes.
rows() {
	awk -F '\t' '
		NR == FNR { if (FNR > 1 && $2 == "pivot" && $6 ~ /^Supersedes /) gone[substr($6, 12)] = 1; next }
		FNR > 1 && !($1 in gone)' "$TSV" "$TSV"
}

# append_row PHASE DECISION WHY EVIDENCE RESULT writes one row with a ts that
# differs from the previous row's, so every ts names one row.
append_row() {
	last=$(tail -n 1 "$TSV" | cut -f 1)
	ts=$(now)
	while [ "$ts" = "$last" ]; do
		sleep 1
		ts=$(now)
	done
	printf '%s\t%s\t%s\t%s\t%s\t%s\n' "$ts" "$1" "$(clean_cell "$2")" "$(clean_cell "$3")" "$(clean_cell "$4")" "$5" >>"$TSV"
	printf 'ts=%s\n' "$ts"
}

# review_base prints BASE_SHA: the merge-base of the parent side and the run's
# tip when that is one commit descending from the capture, else the capture.
# The parent side is the parent branch in Mode B and its upstream in Mode A;
# a Sequential Fallback run lands on the parent, so it keeps the capture.
review_base() {
	parent_side=
	tip=HEAD
	if [ "$PARENT" != none ]; then
		case $MODE in
		B) parent_side=$PARENT tip=refs/heads/$INTEGRATION ;;
		A) parent_side=$PARENT@{upstream} ;;
		esac
	fi
	if [ -n "$parent_side" ] && has_commit "$parent_side" && has_commit "$tip"; then
		bases=$(g merge-base --all "$parent_side" "$tip" || true)
		if [ "$(printf '%s\n' "$bases" | grep -c . || true)" = 1 ] && g merge-base --is-ancestor "$CAPTURE" "$bases"; then
			printf '%s' "$bases"
			return 0
		fi
	fi
	printf '%s' "$CAPTURE"
}

print_state() {
	printf 'dir=%s\n' "$TRAIL"
	printf 'trail=%s\n' "$TSV"
	printf 'mode=%s\n' "$MODE"
	printf 'parent=%s\n' "$PARENT"
	[ -z "$INTEGRATION" ] || printf 'integration=%s\n' "$INTEGRATION"
	printf 'capture=%s\n' "$CAPTURE"
	printf 'base=%s\n' "$(review_base)"
}

cmd_init() {
	SLUG= MODE= PARENT= INTEGRATION= FRESH= DECISION='run plan' WHY='user approved the run plan'
	while [ $# -gt 0 ]; do
		case $1 in
		--fresh) FRESH=1 && shift && continue ;;
		esac
		[ $# -ge 2 ] || die "$1 needs a value"
		case $1 in
		--slug) SLUG=$2 ;;
		--mode) MODE=$2 ;;
		--parent) PARENT=$2 ;;
		--integration) INTEGRATION=$2 ;;
		--decision) DECISION=$2 ;;
		--why) WHY=$2 ;;
		*) die "unknown init option: $1" ;;
		esac
		shift 2
	done
	[ -z "$FRESH" ] || [ -z "$DIR" ] || die "--fresh names a trail by --slug, never by --dir"
	case $MODE in
	A | B | fallback) ;;
	*) die "--mode takes A, B, or fallback" ;;
	esac
	if [ "$MODE" = B ]; then
		[ -n "$INTEGRATION" ] || die "--mode B needs --integration <branch>"
		g check-ref-format --branch "$INTEGRATION" >/dev/null 2>&1 || die "not a branch name: $INTEGRATION"
	fi
	head=$(full_commit HEAD)
	if [ -z "$PARENT" ]; then
		PARENT=$(g branch --show-current 2>/dev/null || true)
		if [ -z "$PARENT" ]; then
			[ "$MODE" != B ] || needs_user parent-unknown "hint=HEAD is detached: ask the user for the parent branch and pass --parent <branch>"
			PARENT=none
		fi
	fi
	if [ "$PARENT" != none ]; then
		parent_sha=$(full_commit "$PARENT")
		: "$parent_sha"
		[ "$PARENT" != "$INTEGRATION" ] ||
			needs_user parent-is-integration "hint=run init on the branch the integration branch starts from, before creating it"
	fi
	resolve_trail create
	if [ -f "$STATE" ]; then
		# Re-running init is a no-op only while nothing happened since it ran.
		expected_rows=1
		[ "$MODE" != A ] || expected_rows=0
		if [ -z "$FRESH" ] && [ "$(state_get repo_key)" = "$KEY" ] && [ "$(state_get mode)" = "$MODE" ] &&
			[ "$(state_get parent)" = "$PARENT" ] && [ "$(state_get integration)" = "$INTEGRATION" ] &&
			[ "$(state_get capture)" = "$head" ] && [ "$(awk 'NR > 1' "$TSV" | grep -c . || true)" = "$expected_rows" ]; then
			CAPTURE=$head
			print_state
			printf 'exists=1\n'
			return 0
		fi
		[ -n "$FRESH" ] || needs_user trail-already-exists "dir=$TRAIL" \
			"hint=a run already used this slug: resume it with status and base, or ask the user before starting over with init --fresh"
		# --fresh keeps the old trail under an archive name; it never deletes it.
		lock
		archive=$TRAIL.archived-$(date -u +%Y%m%dT%H%M%SZ)
		[ ! -e "$archive" ] || die "archive $archive already exists; retry in a second"
		unlock
		mv "$TRAIL" "$archive" || fail "cannot archive $TRAIL"
		mkdir "$TRAIL" || fail "cannot create $TRAIL"
		printf 'archived=%s\n' "$archive"
	fi
	CAPTURE=$head
	lock
	printf 'mode=%s\ncapture=%s\nparent=%s\nintegration=%s\nrepo_key=%s\nrepo=%s\n' \
		"$MODE" "$CAPTURE" "$PARENT" "$INTEGRATION" "$KEY" "$(g rev-parse --show-toplevel)" >"$STATE"
	printf '%s\n' "$HEADER" >"$TSV"
	print_state
	case $MODE in
	B) append_row gate "$DECISION" "$WHY" "BASE_SHA=$CAPTURE parent=$PARENT" Approved ;;
	fallback) append_row gate "$DECISION" "$WHY" "BASE_SHA=$CAPTURE" Approved ;;
	esac
}

cmd_log() {
	SLUG= PHASE= DECISION= WHY= RESULT= TICKET= COMMIT= CHECK= EVIDENCE= DEFERRED= SUPERSEDES=
	while [ $# -gt 0 ]; do
		[ $# -ge 2 ] || die "$1 needs a value"
		case $1 in
		--slug) SLUG=$2 ;;
		--phase) PHASE=$2 ;;
		--decision) DECISION=$2 ;;
		--why) WHY=$2 ;;
		--result) RESULT=$2 ;;
		--ticket) TICKET=$2 ;;
		--commit) COMMIT=$2 ;;
		--check) CHECK=$2 ;;
		--evidence) EVIDENCE=$2 ;;
		--deferred-to) DEFERRED=$2 ;;
		--supersedes) SUPERSEDES=$2 ;;
		*) die "unknown log option: $1" ;;
		esac
		shift 2
	done
	require_trail
	[ -n "$DECISION" ] || die "--decision is required"
	[ -n "$WHY" ] || die "--why is required"
	case $PHASE in
	merge | direct) ;;
	*) [ -z "$TICKET" ] || die "--ticket applies only to merge and direct rows" ;;
	esac
	case $PHASE in
	gate | pivot) [ -z "$DEFERRED" ] || die "--deferred-to applies only to merge, direct, and fix rows" ;;
	esac
	[ "$PHASE:$RESULT" = pivot:Supersedes ] || [ -z "$SUPERSEDES" ] || die "--supersedes applies only to a pivot row with --result Supersedes"
	case $PHASE in
	merge | direct | integration-fix | review-fix)
		case $RESULT in
		Verified | 'Not verified' | Inconclusive) ;;
		*) die "a $PHASE row's --result is Verified, 'Not verified', or Inconclusive" ;;
		esac
		if [ -n "$DEFERRED" ]; then
			[ "$RESULT" = Inconclusive ] || die "--deferred-to needs --result Inconclusive"
			valid_id "$DEFERRED" || die "not a ticket ID: $DEFERRED"
			RESULT="Inconclusive (deferred to final integration ticket #$DEFERRED)"
		fi
		[ -z "$EVIDENCE" ] || die "a $PHASE row's evidence is --commit, optionally with --check"
		if [ -n "$COMMIT" ]; then
			sha=$(full_commit "$COMMIT")
			EVIDENCE=$sha
			[ -z "$CHECK" ] || EVIDENCE="$sha $CHECK"
		elif [ "$MODE" = fallback ] && [ "$PHASE" = direct ] && [ -n "$CHECK" ]; then
			EVIDENCE=$CHECK
		else
			die "a $PHASE row needs --commit <landed-commit>"
		fi
		case $PHASE in
		merge | direct)
			valid_id "$TICKET" || die "a $PHASE row needs --ticket <id>"
			DECISION="#$TICKET $DECISION"
			;;
		esac
		;;
	gate | pivot)
		case $PHASE:$RESULT in
		gate:Approved | gate:Declined | pivot:Reverted | pivot:Replanned) ;;
		pivot:Supersedes)
			[ -n "$SUPERSEDES" ] || die "--result Supersedes needs --supersedes <ts>"
			awk -F '\t' -v ts="$SUPERSEDES" 'NR > 1 && $1 == ts { found = 1 } END { exit !found }' "$TSV" ||
				die "no row has ts $SUPERSEDES"
			RESULT="Supersedes $SUPERSEDES"
			;;
		gate:*) die "a gate row's --result is Approved or Declined" ;;
		*) die "a pivot row's --result is Reverted, Replanned, or Supersedes" ;;
		esac
		if [ -n "$COMMIT" ]; then
			EVIDENCE=$(full_commit "$COMMIT")
		fi
		;;
	*) die "--phase takes merge, direct, integration-fix, review-fix, pivot, or gate" ;;
	esac
	lock
	append_row "$PHASE" "$DECISION" "$WHY" "$EVIDENCE" "$RESULT"
}

cmd_status() {
	SLUG=
	while [ $# -gt 0 ]; do
		[ $# -ge 2 ] || die "$1 needs a value"
		case $1 in
		--slug) SLUG=$2 ;;
		*) die "unknown status option: $1" ;;
		esac
		shift 2
	done
	require_trail
	print_state
	rows | awk -F '\t' '
		$2 == "merge" || $2 == "direct" {
			split($3, id, " ")
			split($5, ev, " ")
			printf "ticket=%s phase=%s commit=%s\n", substr(id[1], 2), $2, ev[1]
			count[$2]++
		}
		END { printf "merge_rows=%d\ndirect_rows=%d\n", count["merge"], count["direct"] }'
}

cmd_base() {
	SLUG= TARGET= WIP= COMMITTED=
	while [ $# -gt 0 ]; do
		case $1 in
		--wip) WIP=1 && shift && continue ;;
		esac
		[ $# -ge 2 ] || die "$1 needs a value"
		case $1 in
		--slug) SLUG=$2 ;;
		--target) TARGET=$2 ;;
		--committed) COMMITTED=$2 ;;
		*) die "unknown base option: $1" ;;
		esac
		shift 2
	done
	[ -z "$WIP" ] || [ -z "$TARGET" ] || die "--wip and --target exclude each other"
	[ -z "$COMMITTED" ] || [ -n "$WIP" ] || die "--committed applies only with --wip"
	require_trail
	base=$(review_base)
	if [ -n "$WIP" ]; then
		line="code-review: $base -> WIP"
		if [ -n "$COMMITTED" ]; then
			committed_sha=$(full_commit "$COMMITTED")
			line="$line $committed_sha"
		fi
	else
		if [ -z "$TARGET" ]; then
			if [ "$MODE" = B ]; then
				TARGET=$INTEGRATION
			else
				TARGET=$(g branch --show-current 2>/dev/null || true)
				TARGET=${TARGET:-HEAD}
			fi
		fi
		target_sha=$(full_commit "$TARGET")
		line="code-review: $base -> $TARGET $target_sha"
	fi
	printf 'base=%s\n' "$base"
	printf 'review_line=%s\n' "$line"
}

cmd_scope() {
	ALLOW= SOURCE= SINCE=
	while [ $# -gt 0 ]; do
		case $1 in
		--staged) SOURCE=staged && shift && continue ;;
		esac
		[ $# -ge 2 ] || die "$1 needs a value"
		case $1 in
		--allow) ALLOW="${ALLOW:+$ALLOW$NL}$2" ;;
		--since) SOURCE=since SINCE=$2 ;;
		*) die "unknown scope option: $1" ;;
		esac
		shift 2
	done
	[ -n "$ALLOW" ] || die "scope needs at least one --allow <write-scope-entry>"
	require_repo
	case $SOURCE in
	staged) files=$(g diff --cached --no-renames --name-only) ;;
	since)
		since_sha=$(full_commit "$SINCE")
		files=$(g diff --no-renames --name-only "$since_sha...HEAD")
		;;
	*) die "scope needs --staged or --since <rev>" ;;
	esac
	# One '*' or '?' stays inside a path segment; '**' and a trailing '/' cross
	# directories.
	result=$(printf '%s\n' "$files" | RUN_TRAIL_ALLOW=$ALLOW awk '
		function glob_re(p,    re, i, c, n) {
			if (substr(p, length(p), 1) == "/") p = p "**"
			n = length(p)
			re = "^"
			for (i = 1; i <= n; i++) {
				c = substr(p, i, 1)
				if (c == "*" && substr(p, i + 1, 1) == "*") {
					i++
					if (substr(p, i + 1, 1) == "/") { i++; re = re "(.*/)?" } else re = re ".*"
				} else if (c == "*") re = re "[^/]*"
				else if (c == "?") re = re "[^/]"
				else if (index(".[]()+{}^$|\\", c)) re = re "\\" c
				else re = re c
			}
			return re "$"
		}
		BEGIN { n = split(ENVIRON["RUN_TRAIL_ALLOW"], pat, "\n"); for (i = 1; i <= n; i++) re[i] = glob_re(pat[i]) }
		$0 != "" {
			checked++
			for (i = 1; i <= n; i++) if ($0 ~ re[i]) next
			print "out_of_scope=" $0
		}
		END { print "checked=" checked + 0 }')
	outside=$(printf '%s\n' "$result" | grep '^out_of_scope=' || true)
	if [ -n "$outside" ]; then
		IFS=$NL
		# shellcheck disable=SC2086
		needs_user out-of-scope $result
	fi
	printf '%s\n' "$result"
}

cmd_worktrees() {
	[ $# -eq 0 ] || die "worktrees takes no options"
	require_repo
	count=0
	first=1
	path=
	branch=
	report() {
		[ -n "$path" ] || return 0
		if [ -n "$first" ]; then
			first=
			return 0
		fi
		count=$((count + 1))
		printf 'leftover=%s branch=%s\n' "$path" "${branch:-detached}"
		if [ -d "$path" ] && [ -n "$(git -C "$path" status --porcelain 2>/dev/null)" ]; then
			printf 'dirty=%s\n' "$path"
		fi
	}
	porcelain=$(g worktree list --porcelain)
	old_ifs=$IFS
	IFS=$NL
	for line in $porcelain; do
		case $line in
		'worktree '*)
			report
			path=${line#worktree } branch=
			;;
		'branch refs/heads/'*) branch=${line#branch refs/heads/} ;;
		esac
	done
	IFS=$old_ifs
	report
	printf 'leftover_count=%d\n' "$count"
}

DIR=
while [ $# -gt 0 ]; do
	case $1 in
	--dir)
		[ $# -ge 2 ] || die "--dir needs a value"
		DIR=$2
		shift 2
		;;
	*) break ;;
	esac
done
[ $# -gt 0 ] || usage
command=$1
shift
case $command in
init | log | status | base) require_repo ;;
esac
case $command in
init) cmd_init "$@" ;;
log) cmd_log "$@" ;;
status) cmd_status "$@" ;;
base) cmd_base "$@" ;;
scope) cmd_scope "$@" ;;
worktrees) cmd_worktrees "$@" ;;
-h | --help | help) usage ;;
*) die "unknown command: $command" ;;
esac

package integration_test

import (
	"path/filepath"
	"strings"
	"testing"
)

// flat reads a file and collapses whitespace so wrapped contract lines match.
func flat(t *testing.T, parts ...string) string {
	t.Helper()
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	return strings.Join(strings.Fields(readFile(t, filepath.Join(append([]string{repoRoot}, parts...)...))), " ")
}

func localTrackerContracts(t *testing.T) map[string]string {
	t.Helper()
	return map[string]string{
		"docs":     flat(t, "docs", "agents", "issue-tracker.md"),
		"template": flat(t, "skills", "setup-matt-pocock-skills", "issue-tracker-local.md"),
	}
}

func builderDefinitions(t *testing.T) map[string]string {
	t.Helper()
	return map[string]string{
		"agents":        flat(t, "agents", "builder.md"),
		"claude-agents": flat(t, "claude-agents", "builder.md"),
		"codex":         flat(t, "config", "codex-agents", "builder.toml"),
	}
}

// A migration batch the spec marks as not individually green lands on scoped
// checks and an Inconclusive row; only the final integration ticket needs the
// whole suite Verified.
func TestDeferredGreenMigrationBatchesLandWithInconclusiveRows(t *testing.T) {
	mustContain(t, flat(t, "skills", "to-spec", "SKILL.md"), "`Not individually green: verified at <final ticket>`")
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "`Inconclusive (deferred to final integration ticket #<id>)`")
	mustContain(t, implement, "The final integrate-and-verify ticket lands only when the full suite is `Verified`.")
}

// A session can resume its own claim, and a stale claim moves only on the
// user's say-so through a CAS on the stale token.
func TestClaimsCanBeResumedOrTakenOverOnlyByExactToken(t *testing.T) {
	wayfinder := flat(t, "skills", "wayfinder", "SKILL.md")
	mustContain(t, wayfinder, "**Resume an own claim**")
	mustContain(t, wayfinder, "a `/handoff` document that carries its exact claim token")
	mustContain(t, wayfinder, "**Stale claim**")
	mustContain(t, wayfinder, "only after the user explicitly authorizes it")
	for name, contract := range localTrackerContracts(t) {
		t.Run(name, func(t *testing.T) {
			mustContain(t, contract, "**Resume")
			mustContain(t, contract, "**Release or take over")
			mustContain(t, contract, "require the exact stale token as the expected prior state")
		})
	}
	for _, name := range []string{"issue-tracker-github.md", "issue-tracker-gitlab.md"} {
		mustContain(t, flat(t, "skills", "setup-matt-pocock-skills", name), "except a user-authorized stale-claim takeover, which carries the exact stale token as its expected prior state")
	}
}

// Worktree verification uses the same gitignored fixtures as the main checkout,
// or the ticket is Blocked; it never passes on a check run somewhere else.
func TestWorktreeVerificationProvisionsIgnoredFixturesOrBlocks(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "link or copy (read-only) the `.gitignore`d fixtures it depends on")
	mustContain(t, implement, "report the ticket `Blocked` with the missing fixture")
	mustNotContain(t, implement, "run its final verification check in the main checkout after merging")
}

func TestPivotRowsCanRecordAReplanBeforeAnyCommit(t *testing.T) {
	mustContain(t, flat(t, "skills", "implement", "SKILL.md"), "`Reverted`, `Replanned`, or `Supersedes <ts>`: for a `pivot` row; `Replanned` records a re-slice or a change of approach before any commit of that ticket landed")
}

// The assigned patch backup is outside product Write scope, in both the skill
// and every builder definition.
func TestPatchBackupPathIsExemptFromWriteScope(t *testing.T) {
	mustContain(t, flat(t, "skills", "implement", "SKILL.md"), "`.scratch/<ticket>-<ts>.patch`, which is exempt from the ticket's `Write scope`")
	for name, builder := range builderDefinitions(t) {
		t.Run(name, func(t *testing.T) {
			mustContain(t, builder, "the `.scratch/<ticket>-<ts>.patch` backup path the dispatching session assigns is exempt")
		})
	}
}

// Only a fulfilled ticket unblocks its dependents; a wontfix, duplicate, or
// out-of-scope closure makes implement stop and ask for a replan.
func TestOnlyFulfilledTicketsUnblockDependents(t *testing.T) {
	for name, contract := range localTrackerContracts(t) {
		t.Run(name, func(t *testing.T) {
			mustContain(t, contract, "`Status: closed (wontfix)`, `Status: closed (duplicate)`, or `Status: closed (out-of-scope)`")
			mustContain(t, contract, "a blocker closed without being fulfilled never unblocks it")
			mustNotContain(t, contract, "closing a ticket is the sole exception that sets `Status: resolved`")
		})
	}
	mustContain(t, flat(t, "skills", "setup-matt-pocock-skills", "issue-tracker-github.md"), "A blocker closed as `not_planned` or `duplicate` never unblocks it.")
	mustContain(t, flat(t, "skills", "setup-matt-pocock-skills", "issue-tracker-gitlab.md"), "A blocker closed with the configured `wontfix` label never unblocks it.")
	mustContain(t, flat(t, "skills", "triage", "SKILL.md"), "a closure that does not fulfil the request")
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "A blocker closed without being fulfilled")
	mustContain(t, implement, "stop and ask the user to replan that dependency")
}

// New tickets never reuse a number in the issues directory, and implement
// discovers only implementation tickets of the spec it builds.
func TestTicketNumbersNeverCollideAndDiscoveryFiltersDecisionTickets(t *testing.T) {
	mustContain(t, flat(t, "skills", "to-spec", "SKILL.md"), "continuing after the highest number already in that `issues/` directory")
	for name, contract := range localTrackerContracts(t) {
		t.Run(name, func(t *testing.T) {
			mustContain(t, contract, "continuing after the highest existing number; a number is never reused")
		})
	}
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "A file with a `Type:` of `research`, `prototype`, or `grilling` is a wayfinder decision ticket, not a child.")
}

// Resolve is safe to repeat: a resolved, answered child with no map pointer
// gets exactly one pointer the next time the map loads.
func TestResolveReconcilesMissingMapPointers(t *testing.T) {
	for name, contract := range localTrackerContracts(t) {
		t.Run(name, func(t *testing.T) {
			mustContain(t, contract, "Resolve is idempotent")
			mustContain(t, contract, "deduplicated by ticket number")
		})
	}
	mustContain(t, flat(t, "skills", "wayfinder", "SKILL.md"), "append the missing pointer for each resolved, answered ticket the map does not yet point to, once per ticket")
}

// A resumed Mode B run reconciles its decision trail against the integration
// branch before dispatch, so landed tickets are never built twice.
func TestImplementResumeReconcilesLandedTickets(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "**Resume an interrupted run**")
	mustContain(t, implement, "reconcile its `merge` and `direct` rows against `git log` on the integration branch")
	mustContain(t, implement, "never re-dispatch it")
	mustContain(t, flat(t, "skills", "handoff", "SKILL.md"), "so the next session reconciles landed tickets before dispatching")
}

// code-review receives a single base revision and a target, never a range
// string that cannot pass the ^{commit} check.
func TestImplementPassesASingleBaseRevisionToCodeReview(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "passing `BASE_SHA` as the fixed point and the integration branch as the target")
	mustNotContain(t, implement, "once on `base...integration-branch`")
	review := flat(t, "skills", "code-review", "SKILL.md")
	mustContain(t, review, "When the caller names a target branch or commit")
	mustContain(t, review, "`<full-hex-commit>...<target-full-hex>`")
}

func TestSequentialFallbackReviewsOnceAtTheEnd(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "per-child review: after each child, run its focused verification")
	mustContain(t, implement, "After the last child, use the `code-review` skill once (Standards and Spec)")
}

// Builders rebase only on a forwarded approval that names the rebase.
func TestRebaseApprovalIsForwardedToBuilders(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "each builder's rebase of its worktree branch onto the integration tip")
	mustContain(t, implement, "`User pre-approved in run plan: you are authorized to git commit on your worktree branch and git rebase onto <integration-branch>.`")
	for name, builder := range builderDefinitions(t) {
		t.Run(name, func(t *testing.T) {
			mustContain(t, builder, "A forwarded approval for worktree-branch commits and for rebasing onto the integration branch covers those two actions only.")
			mustContain(t, builder, "Without a forwarded rebase approval, report that the branch needs a rebase instead.")
		})
	}
}

func TestEmptyFrontierWithUnresolvedTicketsStopsAsBlocked(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "**Stuck graph**")
	mustContain(t, implement, "stop and report `Blocked` with the exact blocking chain")
}

// section returns the text from one heading to the next, or to the end of the
// text when to is empty, failing when either heading is missing.
func section(t *testing.T, text, from, to string) string {
	t.Helper()
	i := strings.Index(text, from)
	if i < 0 {
		t.Fatalf("missing section start %q", from)
	}
	if to == "" {
		return text[i:]
	}
	j := strings.Index(text[i:], to)
	if j < 0 {
		t.Fatalf("missing section end %q after %q", to, from)
	}
	return text[i : i+j]
}

// code-review runs on every completed change, whatever its size or Agency Gate
// option, and no explanation stands in for running it.
func TestCodeReviewIsNeverSkippedOrExcused(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	modeA := section(t, implement, "## Mode A:", "## Sequential Fallback")
	modeB := section(t, implement, "## Mode B:", "## Final Report")
	cleanup := section(t, implement, "8. **Cleanup**", "## Final Report")
	report := section(t, implement, "## Final Report", "")
	mustContain(t, modeA, "regardless of diff size, including a one-line change and a change made under Agency Gate option (a) Local minimal edit")
	mustContain(t, modeB, "regardless of diff size")
	for name, part := range map[string]string{"step 8": cleanup, "final report": report} {
		t.Run(name, func(t *testing.T) {
			mustContain(t, part, "stop and invoke the `code-review` skill first; never substitute an explanation for running it")
			mustContain(t, part, "applies only when the run completes implementation")
		})
	}
	mustContain(t, report, "`code-review: <BASE_SHA> -> <target>`")
	mustNotContain(t, implement, "or the reason it did not run")
	// The review check is not a decisions.tsv row and not a gate.
	mustContain(t, cleanup, "**Review completion check**")
	mustNotContain(t, implement, "**Review gate**")
	mustContain(t, implement, "`decisions.tsv` has no `code-review` phase")
	// Uncommitted changes are reviewed in WIP mode, never as a target ref.
	mustContain(t, modeA, "WIP mode anchored at `BASE_SHA`")
	mustContain(t, report, "`code-review: <BASE_SHA> -> WIP`")
	mustNotContain(t, implement, "`working-tree (WIP)`")
}

// The fixed point is captured before any ticket lands, so a directly committed
// contract ticket is never mistaken for the base; a resumed run inherits it.
func TestBaseShaIsCapturedBeforeTheFirstTicketLands(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	modeB := section(t, implement, "## Mode B:", "## Final Report")
	capture := strings.Index(modeB, "before any ticket lands, capture `BASE_SHA=$(git rev-parse HEAD)`")
	landing := strings.Index(modeB, "Immediately after each ticket lands on the integration branch")
	if capture < 0 || landing < 0 || capture > landing {
		t.Fatalf("Mode B must capture BASE_SHA before the first landing rule (capture=%d, landing=%d)", capture, landing)
	}
	mustContain(t, implement, "never the commit of the first landed ticket")
	// A worktree merge can carry several commits, so the parent of the first
	// landed commit is not the branch point.
	mustContain(t, implement, "`git merge-base <integration-branch> <parent-branch>`")
	mustContain(t, implement, "`<earliest-direct-sha>^`")
	mustNotContain(t, implement, "<earliest-landed-sha>^")
	// Sequential Fallback keeps one BASE_SHA across all children.
	mustContain(t, section(t, implement, "## Mode A:", "## Sequential Fallback"), "`BASE_SHA=${BASE_SHA:-$(git rev-parse HEAD)}`")
	fallback := section(t, implement, "## Sequential Fallback", "## Mode B:")
	mustContain(t, fallback, "Skip Mode A's per-child `BASE_SHA` capture and per-child review")
	mustContain(t, fallback, "On a fresh run, capture `BASE_SHA=$(git rev-parse HEAD)` once before the first child")
	mustContain(t, fallback, "When resuming an interrupted Sequential Fallback run, never re-capture `HEAD`")
	mustContain(t, fallback, "the earliest `direct` row whose `evidence` is a commit SHA")
	mustContain(t, fallback, "if no `direct` row has a commit SHA, ask the user")
	mustContain(t, fallback, "run it in WIP mode anchored at `BASE_SHA`")
	// Mode B resume reads the recorded base first and never guesses a parent branch.
	modeBResume := section(t, implement, "A resumed run takes `BASE_SHA`", "In the same step, before the first dispatch")
	mustContain(t, modeBResume, "the run-plan `gate` row")
	mustContain(t, modeBResume, "only when `<parent-branch>` is explicitly known from the session or the spec; never guess `main` or `master`")
	mustContain(t, modeBResume, "the earliest `direct` row whose `evidence` is a commit SHA")
	mustContain(t, implement, "The first row is the run-plan `gate` row")
	modeA := section(t, implement, "## Mode A:", "## Sequential Fallback")
	// A WIP review covers exactly the files the task modified or created.
	for name, part := range map[string]string{"mode A": modeA, "fallback": fallback} {
		t.Run(name, func(t *testing.T) {
			mustContain(t, part, "the files modified or created for this task as the pathspec")
		})
	}
	mustContain(t, fallback, "record it in the run-plan `gate` row")
	step4 := section(t, implement, "4. **Create the Integration Branch and Its Decision Trail**", "5. **Dispatch Frontier")
	mustContain(t, step4, "`BASE_SHA=<full-sha>`")
	mustContain(t, step4, "`parent=<branch>`")
	mustContain(t, modeBResume, "If none of the applicable sources above yields `BASE_SHA`, ask the user.")
	mustNotContain(t, modeBResume, "no known parent branch, or no `direct` row with a commit SHA")
	mustContain(t, flat(t, "skills", "handoff", "SKILL.md"), "the integration branch and the branch it was created from (when used), the run's `BASE_SHA`, and the `decisions.tsv` path")
}

// A builder pause answer is a gate row with a legal result.
func TestPauseAnswersRecordAGateResult(t *testing.T) {
	mustContain(t, flat(t, "skills", "implement", "SKILL.md"), "the chosen pattern or seam in `decision` and `why`, and `result` `Approved`, or `Declined` when the user rejects the proposal")
}

// result cells are bare tokens, each row is written when its event happens
// with the real clock, and fix phases follow the Step 7 review.
func TestDecisionRowsUseBareTokensAndRealTimestamps(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "is a bare token with no parenthetical note or prose")
	mustContain(t, implement, "`date -u +\"%Y-%m-%dT%H:%M:%SZ\"`")
	mustContain(t, implement, "Never invent or round a timestamp, and never write rows in a batch at the end of the run.")
	mustContain(t, implement, "A fix committed before Step 7's `code-review` runs is `integration-fix`.")
	mustContain(t, implement, "A `gate` row's `result` is only `Approved` or `Declined`.")
}

func TestFixAndLandingRowsCiteOneCommit(t *testing.T) {
	mustContain(t, flat(t, "skills", "implement", "SKILL.md"), "`evidence` is the single landed commit SHA, never an `a...b` revision range")
}

// A direct commit's scope is checked mechanically on the staged change.
func TestDirectCommitScopeIsCheckedOnTheStagedDiff(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "run `git diff --cached --name-only` before committing")
	mustContain(t, implement, "a manual read of the files is not this check")
}

func TestStyleQuestionCitesFileAndLine(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "`Pattern A (path/file.ext:12) vs. Pattern B (path/other.ext:34)`")
	mustContain(t, implement, "A file name without a line number is not a citation.")
}

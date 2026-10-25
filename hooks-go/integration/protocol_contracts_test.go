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
	mustContain(t, flat(t, "skills", "handoff", "SKILL.md"), "the integration branch and the `decisions.tsv` path")
}

// code-review receives a single base revision and a target, never a range
// string that cannot pass the ^{commit} check.
func TestImplementPassesASingleBaseRevisionToCodeReview(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "the single base commit SHA the integration branch started from as the fixed point and the integration branch as the target")
	mustNotContain(t, implement, "once on `base...integration-branch`")
	review := flat(t, "skills", "code-review", "SKILL.md")
	mustContain(t, review, "When the caller names a target branch or commit")
	mustContain(t, review, "`<full-hex-commit>...<target-full-hex>`")
}

func TestSequentialFallbackReviewsOnceAtTheEnd(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "Skip Mode A's per-child review")
	mustContain(t, implement, "After the last child, use the `code-review` skill once (Standards and Spec)")
}

// Builders rebase only on a forwarded approval that names the rebase.
func TestRebaseApprovalIsForwardedToBuilders(t *testing.T) {
	implement := flat(t, "skills", "implement", "SKILL.md")
	mustContain(t, implement, "each builder's rebase of its worktree branch onto the integration tip")
	mustContain(t, implement, "state that the user approved commits on its worktree branch and rebasing that branch onto the integration branch tip")
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

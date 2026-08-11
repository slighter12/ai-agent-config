package integration_test

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestWayfinderClaimRequiresFreshExclusiveCAS(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	skill := readFile(t, filepath.Join(repoRoot, "skills", "wayfinder", "SKILL.md"))
	contract := strings.Join(strings.Fields(skill), " ")

	for _, want := range []string{
		"A named ticket and a ticket returned by Frontier are eligible only when a fresh read at claim time",
		"open, unblocked, and unclaimed",
		"conditional compare-and-set (CAS)",
		"exact expected prior state",
		"unique opaque owner/session token",
		"post-claim read confirms the same token",
		"stop without work when exclusivity cannot be guaranteed",
	} {
		mustContain(t, contract, want)
	}

}

func TestLocalClaimContractUsesAtomicReplacementAndRecheck(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	paths := []string{
		filepath.Join(repoRoot, "docs", "agents", "issue-tracker.md"),
		filepath.Join(repoRoot, "skills", "setup-matt-pocock-skills", "issue-tracker-local.md"),
	}
	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			contract := strings.Join(strings.Fields(readFile(t, path)), " ")
			for _, want := range []string{
				"conditional compare-and-set",
				"bounded exclusive lock or lease",
				"re-read the exact file",
				"expected prior state",
				"Atomically replace the file",
				"unique opaque owner/session token",
				"Never overwrite an existing claim",
				"stop without work",
			} {
				mustContain(t, contract, want)
			}
		})
	}
}

func TestHostedClaimContractRequiresProviderConditionalExclusivity(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	for _, name := range []string{"issue-tracker-github.md", "issue-tracker-gitlab.md"} {
		t.Run(name, func(t *testing.T) {
			contract := strings.Join(strings.Fields(readFile(t, filepath.Join(repoRoot, "skills", "setup-matt-pocock-skills", name))), " ")
			for _, want := range []string{
				"fresh read immediately before claiming",
				"provider-supported conditional update",
				"version/ETag",
				"expected prior state",
				"unique opaque owner/session token",
				"bounded provider lock/lease",
				"documented atomic provider claim primitive",
				"immediate re-read that confirms the same token",
				"stop without work",
			} {
				mustContain(t, contract, want)
			}
			mustContain(t, contract, "A plain read-then-assign or unconditional update is not exclusive")
			mustNotContain(t, contract, "assign fixed `ISSUE_NUMBER` to the authenticated session")
			mustNotContain(t, contract, "assign fixed `ISSUE_IID` to the authenticated session")
		})
	}
}

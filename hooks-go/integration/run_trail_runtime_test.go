package integration_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// trailRepo is a throwaway git repository plus the private TMPDIR that
// run-trail.sh keeps its trail in.
type trailRepo struct {
	t   *testing.T
	dir string
	tmp string
	env []string
}

func newTrailRepo(t *testing.T) *trailRepo {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is unavailable")
	}
	root := t.TempDir()
	r := &trailRepo{t: t, dir: filepath.Join(root, "repo"), tmp: filepath.Join(root, "tmp")}
	for _, dir := range []string{r.dir, r.tmp, filepath.Join(root, "home")} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	globalConfig := filepath.Join(root, "gitconfig")
	if err := os.WriteFile(globalConfig, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, item := range os.Environ() {
		if !strings.HasPrefix(item, "GIT_") && !strings.HasPrefix(item, "TMPDIR=") && !strings.HasPrefix(item, "HOME=") {
			r.env = append(r.env, item)
		}
	}
	r.env = append(r.env,
		"TMPDIR="+r.tmp,
		"HOME="+filepath.Join(root, "home"),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL="+globalConfig,
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
	)
	r.git("init", "-q", "-b", "main")
	r.commit("base.txt", "base\n", "c0")
	return r
}

func (r *trailRepo) git(args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = r.dir
	cmd.Env = r.env
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func (r *trailRepo) write(name, content string) {
	r.t.Helper()
	path := filepath.Join(r.dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		r.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		r.t.Fatal(err)
	}
}

func (r *trailRepo) commit(name, content, message string) string {
	r.t.Helper()
	r.write(name, content)
	r.git("add", "--", name)
	r.git("commit", "-q", "-m", message)
	return r.git("rev-parse", "HEAD")
}

// run executes run-trail.sh in the repository and returns stdout+stderr and
// the exit code.
func (r *trailRepo) run(args ...string) (string, int) {
	r.t.Helper()
	shell := "sh"
	if dash, err := exec.LookPath("dash"); err == nil {
		shell = dash // the strictest POSIX shell available
	}
	script, err := filepath.Abs(filepath.Join("..", "..", "skills", "implement", "scripts", "run-trail.sh"))
	if err != nil {
		r.t.Fatal(err)
	}
	cmd := exec.Command(shell, append([]string{script}, args...)...)
	cmd.Dir = r.dir
	cmd.Env = r.env
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	var exitErr *exec.ExitError
	switch {
	case err == nil:
		return output.String(), 0
	case errors.As(err, &exitErr):
		return output.String(), exitErr.ExitCode()
	default:
		r.t.Fatalf("run-trail.sh: %v\n%s", err, output.String())
		return "", -1
	}
}

func (r *trailRepo) mustRun(want int, args ...string) string {
	r.t.Helper()
	out, code := r.run(args...)
	if code != want {
		r.t.Fatalf("run-trail.sh %s: exit %d, want %d\n%s", strings.Join(args, " "), code, want, out)
	}
	return out
}

func outputValue(t *testing.T, out, key string) string {
	t.Helper()
	for line := range strings.SplitSeq(out, "\n") {
		if value, found := strings.CutPrefix(line, key+"="); found {
			return value
		}
	}
	t.Fatalf("no %s= line in output:\n%s", key, out)
	return ""
}

func trailRows(t *testing.T, trail string) [][]string {
	t.Helper()
	lines := strings.Split(strings.TrimRight(readFile(t, trail), "\n"), "\n")
	if lines[0] != "ts\tphase\tdecision\twhy\tevidence\tresult" {
		t.Fatalf("trail header = %q", lines[0])
	}
	var rows [][]string
	for _, line := range lines[1:] {
		cells := strings.Split(line, "\t")
		if len(cells) != 6 {
			t.Fatalf("row has %d cells: %q", len(cells), line)
		}
		rows = append(rows, cells)
	}
	return rows
}

var trailTimestamp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

// startModeB inits a Mode B trail on main and checks out the integration
// branch, as implement Step 4 orders it.
func (r *trailRepo) startModeB() (trail, capture string) {
	r.t.Helper()
	out := r.mustRun(0, "init", "--slug", "feat", "--mode", "B", "--integration", "feat/x")
	r.git("checkout", "-q", "-b", "feat/x")
	return outputValue(r.t, out, "trail"), outputValue(r.t, out, "capture")
}

func (r *trailRepo) landTicket(id, file string) string {
	r.t.Helper()
	sha := r.commit(file, id+"\n", "ticket "+id)
	r.mustRun(0, "log", "--slug", "feat", "--phase", "merge", "--ticket", id,
		"--decision", "land "+id, "--why", "fast-forward", "--result", "Verified", "--commit", sha)
	return sha
}
func TestRunTrailInitCapturesTheBaseOutsideTheWorkingTree(t *testing.T) {
	r := newTrailRepo(t)
	head := r.git("rev-parse", "HEAD")
	out := r.mustRun(0, "init", "--slug", "feat", "--mode", "B", "--integration", "feat/x")
	trail := outputValue(t, out, "trail")
	if outputValue(t, out, "capture") != head || outputValue(t, out, "base") != head || outputValue(t, out, "parent") != "main" {
		t.Fatalf("init captured the wrong base or parent:\n%s", out)
	}
	if !strings.HasPrefix(trail, physicalPath(t, r.tmp)+string(filepath.Separator)+"run-trail"+string(filepath.Separator)) &&
		!strings.HasPrefix(trail, r.tmp+string(filepath.Separator)+"run-trail"+string(filepath.Separator)) {
		t.Fatalf("trail %s is not under TMPDIR %s", trail, r.tmp)
	}
	if status := r.git("status", "--porcelain", "--untracked-files=all"); status != "" {
		t.Fatalf("init left files in the working tree:\n%s", status)
	}
	rows := trailRows(t, trail)
	if len(rows) != 1 || rows[0][1] != "gate" || rows[0][4] != "BASE_SHA="+head+" parent=main" || rows[0][5] != "Approved" || !trailTimestamp.MatchString(rows[0][0]) {
		t.Fatalf("run-plan gate row = %q", rows)
	}

	t.Run("init is a no-op only while nothing happened since it ran", func(t *testing.T) {
		again := r.mustRun(0, "init", "--slug", "feat", "--mode", "B", "--integration", "feat/x")
		mustContain(t, again, "exists=1")
		if outputValue(t, again, "capture") != head || len(trailRows(t, trail)) != 1 {
			t.Fatalf("a second init changed the trail:\n%s", again)
		}
		mustContain(t, r.mustRun(3, "init", "--slug", "feat", "--mode", "A"), "needs_user=trail-already-exists")
	})

	t.Run("a resumed session finds the trail with --dir", func(t *testing.T) {
		out := r.mustRun(0, "--dir", filepath.Dir(trail)+"/", "status")
		if outputValue(t, out, "capture") != head || outputValue(t, out, "dir") != filepath.Dir(trail) {
			t.Fatalf("status --dir read another trail:\n%s", out)
		}
		other := newTrailRepo(t)
		other.mustRun(2, "--dir", filepath.Dir(trail), "status")
	})

	t.Run("--fresh archives the old trail and captures the new HEAD", func(t *testing.T) {
		moved := r.commit("later.txt", "later\n", "later")
		out := r.mustRun(3, "init", "--slug", "feat", "--mode", "B", "--integration", "feat/x")
		if !strings.HasPrefix(out, "needs_user=trail-already-exists\n") {
			t.Fatalf("a reused slug did not stop first:\n%s", out)
		}
		fresh := r.mustRun(0, "init", "--slug", "feat", "--mode", "B", "--integration", "feat/x", "--fresh")
		archive := outputValue(t, fresh, "archived")
		if outputValue(t, fresh, "capture") != moved {
			t.Fatalf("--fresh kept the old capture:\n%s", fresh)
		}
		if old := trailRows(t, filepath.Join(archive, "decisions.tsv")); old[0][4] != "BASE_SHA="+head+" parent=main" {
			t.Fatalf("--fresh lost the old trail: %q", old)
		}
	})

	t.Run("an unsafe slug is rejected", func(t *testing.T) {
		r.mustRun(2, "init", "--slug", "..", "--mode", "A")
	})
}

func TestRunTrailInitNeverGuessesTheParent(t *testing.T) {
	r := newTrailRepo(t)
	r.git("checkout", "-q", "--detach")
	out := r.mustRun(3, "init", "--slug", "feat", "--mode", "B", "--integration", "feat/x")
	mustContain(t, out, "needs_user=parent-unknown")
	out = r.mustRun(0, "init", "--slug", "feat", "--mode", "B", "--integration", "feat/x", "--parent", "main")
	if outputValue(t, out, "parent") != "main" {
		t.Fatalf("explicit parent was not recorded:\n%s", out)
	}

	on := newTrailRepo(t)
	on.git("checkout", "-q", "-b", "feat/x")
	mustContain(t, on.mustRun(3, "init", "--slug", "feat", "--mode", "B", "--integration", "feat/x"), "needs_user=parent-is-integration")
}

func TestRunTrailLogValidatesAndFormatsRows(t *testing.T) {
	r := newTrailRepo(t)
	trail, _ := r.startModeB()
	land := r.commit("one.txt", "1\n", "ticket 01")

	rejected := []struct {
		name string
		args []string
	}{
		{"unknown phase", []string{"--phase", "landing", "--result", "Verified", "--commit", land, "--ticket", "01"}},
		{"result with a note", []string{"--phase", "merge", "--result", "Verified (3 tests)", "--commit", land, "--ticket", "01"}},
		{"gate with a verdict", []string{"--phase", "gate", "--result", "Verified"}},
		{"revision range", []string{"--phase", "merge", "--result", "Verified", "--commit", "HEAD~1...HEAD", "--ticket", "01"}},
		{"unknown commit", []string{"--phase", "merge", "--result", "Verified", "--commit", "no-such-rev", "--ticket", "01"}},
		{"merge without a ticket", []string{"--phase", "merge", "--result", "Verified", "--commit", land}},
		{"landing without a commit", []string{"--phase", "integration-fix", "--result", "Verified", "--check", "go test => ok"}},
		{"supersedes an unknown row", []string{"--phase", "pivot", "--result", "Supersedes", "--supersedes", "1999-01-01T00:00:00Z"}},
		{"deferral on a verified row", []string{"--phase", "direct", "--result", "Verified", "--commit", land, "--ticket", "02", "--deferred-to", "05"}},
	}
	for _, tt := range rejected {
		t.Run(tt.name, func(t *testing.T) {
			r.mustRun(2, append([]string{"log", "--slug", "feat", "--decision", "d", "--why", "w"}, tt.args...)...)
		})
	}
	if rows := trailRows(t, trail); len(rows) != 1 {
		t.Fatalf("rejected calls wrote rows: %q", rows)
	}

	r.mustRun(0, "log", "--slug", "feat", "--phase", "merge", "--ticket", "01", "--decision", "land\t01\nnow",
		"--why", "fast-forward", "--result", "Verified", "--commit", land[:7], "--check", "go test ./... => ok")
	r.mustRun(0, "log", "--slug", "feat", "--phase", "direct", "--ticket", "02", "--decision", "batch", "--why", "w",
		"--result", "Inconclusive", "--deferred-to", "05", "--commit", land)
	rows := trailRows(t, trail)
	if merge := rows[1]; merge[2] != "#01 land 01 now" || merge[4] != land+" go test ./... => ok" || merge[5] != "Verified" {
		t.Fatalf("merge row cells = %q", merge)
	}
	if rows[2][5] != "Inconclusive (deferred to final integration ticket #05)" {
		t.Fatalf("deferred row result = %q", rows[2][5])
	}
	r.mustRun(0, "log", "--slug", "feat", "--phase", "pivot", "--decision", "wrong row", "--why", "02 was not direct",
		"--result", "Supersedes", "--supersedes", rows[2][0])
	rows = trailRows(t, trail)
	if rows[3][5] != "Supersedes "+rows[2][0] {
		t.Fatalf("pivot row result = %q", rows[3][5])
	}
	seen := map[string]bool{}
	for _, row := range rows {
		if !trailTimestamp.MatchString(row[0]) || seen[row[0]] {
			t.Fatalf("ts %q is malformed or repeated", row[0])
		}
		seen[row[0]] = true
	}
	// The superseded row no longer counts as a landed ticket.
	status := r.mustRun(0, "status", "--slug", "feat")
	mustContain(t, status, "ticket=01 phase=merge commit="+land)
	mustNotContain(t, status, "ticket=02")
	mustContain(t, status, "merge_rows=1")
	mustContain(t, status, "direct_rows=0")
}

func TestRunTrailLockTimeoutAsksTheUser(t *testing.T) {
	t.Parallel()
	r := newTrailRepo(t)
	trail, _ := r.startModeB()
	if err := os.Mkdir(filepath.Join(filepath.Dir(trail), "lock"), 0o700); err != nil {
		t.Fatal(err)
	}
	out := r.mustRun(3, "log", "--slug", "feat", "--phase", "gate", "--decision", "d", "--why", "w", "--result", "Approved")
	mustContain(t, out, "needs_user=trail-locked")
	mustContain(t, out, "never remove the lock yourself")
}

func TestRunTrailBaseIsTheParentMergeBaseInModeB(t *testing.T) {
	r := newTrailRepo(t)
	_, capture := r.startModeB()
	r.landTicket("01", "one.txt")
	r.git("checkout", "-q", "main")
	upstream := r.commit("upstream.txt", "u\n", "upstream")
	r.git("checkout", "-q", "feat/x")
	if got := outputValue(t, r.mustRun(0, "base", "--slug", "feat"), "base"); got != capture {
		t.Fatalf("base = %s before the parent was merged, want the capture %s", got, capture)
	}

	r.git("merge", "-q", "--no-edit", "main")
	out := r.mustRun(0, "base", "--slug", "feat")
	if outputValue(t, out, "base") != upstream {
		t.Fatalf("base = %s, want the merge-base %s\n%s", outputValue(t, out, "base"), upstream, out)
	}
	tip := r.git("rev-parse", "feat/x")
	if got, want := outputValue(t, out, "review_line"), "code-review: "+upstream+" -> feat/x "+tip; got != want {
		t.Fatalf("review_line = %q, want %q", got, want)
	}
	reviewed := r.git("diff", "--name-only", upstream, tip)
	mustContain(t, reviewed, "one.txt")
	mustNotContain(t, reviewed, "upstream.txt")
	mustContain(t, r.mustRun(0, "status", "--slug", "feat"), "base="+upstream)
}

func TestRunTrailBaseFollowsTheUpstreamInModeA(t *testing.T) {
	r := newTrailRepo(t)
	r.git("branch", "upstream-main")
	r.git("branch", "--set-upstream-to=upstream-main", "main")
	r.mustRun(0, "init", "--slug", "fix", "--mode", "A")
	r.write("work.txt", "uncommitted\n")
	r.git("checkout", "-q", "upstream-main")
	pulled := r.commit("pulled.txt", "p\n", "pulled from upstream")
	r.git("checkout", "-q", "main")
	r.git("merge", "-q", "--ff-only", "upstream-main")
	out := r.mustRun(0, "base", "--slug", "fix", "--wip")
	if got, want := outputValue(t, out, "review_line"), "code-review: "+pulled+" -> WIP"; got != want {
		t.Fatalf("review_line = %q, want %q", got, want)
	}
	committed := r.commit("work.txt", "uncommitted\n", "the reviewed change")
	out = r.mustRun(0, "base", "--slug", "fix", "--wip", "--committed", "HEAD")
	if got, want := outputValue(t, out, "review_line"), "code-review: "+pulled+" -> WIP "+committed; got != want {
		t.Fatalf("review_line = %q, want %q", got, want)
	}
}

func TestRunTrailBaseKeepsTheCapture(t *testing.T) {
	t.Run("Mode A without an upstream", func(t *testing.T) {
		r := newTrailRepo(t)
		capture := outputValue(t, r.mustRun(0, "init", "--slug", "fix", "--mode", "A"), "capture")
		r.commit("work.txt", "w\n", "work")
		if got := outputValue(t, r.mustRun(0, "base", "--slug", "fix", "--wip"), "base"); got != capture {
			t.Fatalf("base = %s, want the capture %s", got, capture)
		}
	})

	t.Run("an upstream behind the capture", func(t *testing.T) {
		r := newTrailRepo(t)
		r.git("branch", "upstream-main")
		r.git("branch", "--set-upstream-to=upstream-main", "main")
		r.commit("local.txt", "l\n", "local only, not pushed")
		capture := outputValue(t, r.mustRun(0, "init", "--slug", "fix", "--mode", "A"), "capture")
		if got := outputValue(t, r.mustRun(0, "base", "--slug", "fix", "--wip"), "base"); got != capture {
			t.Fatalf("base = %s, want the capture %s", got, capture)
		}
	})

	t.Run("Sequential Fallback", func(t *testing.T) {
		r := newTrailRepo(t)
		capture := outputValue(t, r.mustRun(0, "init", "--slug", "feat", "--mode", "fallback"), "capture")
		child := r.commit("one.txt", "1\n", "child 01")
		out := r.mustRun(0, "base", "--slug", "feat", "--target", "main")
		if got, want := outputValue(t, out, "review_line"), "code-review: "+capture+" -> main "+child; got != want {
			t.Fatalf("review_line = %q, want %q", got, want)
		}
	})
}

func TestRunTrailScopeChecksTheDiffAgainstTheWriteScope(t *testing.T) {
	r := newTrailRepo(t)
	r.write("pkg/a/a.go", "package a\n")
	r.write("pkg/b/b.go", "package b\n")
	r.write("pkg/b/b_test.go", "package b\n")
	r.git("add", "--", "pkg")

	out := r.mustRun(3, "scope", "--staged", "--allow", "pkg/a/")
	mustContain(t, out, "needs_user=out-of-scope")
	mustContain(t, out, "out_of_scope=pkg/b/b.go")
	mustNotContain(t, out, "out_of_scope=pkg/a/a.go")
	mustContain(t, r.mustRun(0, "scope", "--staged", "--allow", "pkg/a/", "--allow", "pkg/b/*.go"), "checked=3")
	// One '*' stays inside a path segment; '**' crosses directories. Patterns are
	// matched against paths, never expanded against the working directory.
	mustContain(t, r.mustRun(3, "scope", "--staged", "--allow", "pkg/*"), "out_of_scope=pkg/a/a.go")
	mustContain(t, r.mustRun(0, "scope", "--staged", "--allow", "pkg/**"), "checked=3")
	mustContain(t, r.mustRun(3, "scope", "--staged", "--allow", "pkg/b/*.go"), "out_of_scope=pkg/a/a.go")
	if out := r.mustRun(3, "scope", "--staged", "--allow", "pkg/a/"); !strings.HasPrefix(out, "needs_user=out-of-scope\n") {
		t.Fatalf("needs_user is not the first line:\n%s", out)
	}

	base := r.git("rev-parse", "HEAD")
	r.git("commit", "-q", "-m", "staged work")
	mustContain(t, r.mustRun(3, "scope", "--since", base, "--allow", "pkg/a/"), "out_of_scope=pkg/b/b.go")
	r.mustRun(2, "scope", "--staged")
	// A rename lists its out-of-scope source, not only its in-scope destination.
	r.write("secret/k", "key\n")
	r.git("add", "--", "secret/k")
	r.git("commit", "-q", "-m", "secret")
	r.git("mv", "secret/k", "src_k")
	out = r.mustRun(3, "scope", "--staged", "--allow", "src_k")
	mustContain(t, out, "out_of_scope=secret/k")
	renamed := r.git("rev-parse", "HEAD")
	r.git("commit", "-q", "-m", "rename")
	mustContain(t, r.mustRun(3, "scope", "--since", renamed, "--allow", "src_k"), "out_of_scope=secret/k")
	// An unknown revision is rejected, never read as an empty diff.
	r.mustRun(2, "scope", "--since", "no-such-rev", "--allow", "pkg/a/")
	r.mustRun(2, "scope", "--since", "", "--allow", "pkg/a/")
}

func TestRunTrailWorktreesReportsLeftovers(t *testing.T) {
	r := newTrailRepo(t)
	mustContain(t, r.mustRun(0, "worktrees"), "leftover_count=0")
	worktree := filepath.Join(filepath.Dir(r.dir), "wt one")
	r.git("worktree", "add", "-q", "-b", "ticket-01", worktree)
	if err := os.WriteFile(filepath.Join(worktree, "scratch.txt"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	out := r.mustRun(0, "worktrees")
	mustContain(t, out, "leftover_count=1")
	mustContain(t, out, "branch=ticket-01")
	mustContain(t, out, "dirty=")
	mustContain(t, out, "wt one")
	if err := os.Remove(filepath.Join(worktree, "scratch.txt")); err != nil {
		t.Fatal(err)
	}
	r.git("worktree", "remove", worktree)
	mustContain(t, r.mustRun(0, "worktrees"), "leftover_count=0")
	if branches := r.git("branch", "--list", "ticket-01"); !strings.Contains(branches, "ticket-01") {
		t.Fatal("worktree removal deleted its branch")
	}
}

// Package integration_test exercises repository-owned runtime contracts.
package integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWizardGitHubWritesRequireRepoAndConfirmation(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	templatePath := filepath.Join(repoRoot, "skills", "wizard", "template.sh")
	template := readFile(t, templatePath)
	library := wizardLibrary(t, template)

	t.Run("missing repo skips write", func(t *testing.T) {
		output, calls := runWizardLibrary(t, library, "", `unset GITHUB_REPO
set_secret API_TOKEN "sensitive-test-value"
`)
		mustContain(t, output, "GITHUB_REPO")
		if strings.Contains(calls, "secret set") {
			t.Fatalf("unexpected GitHub write:\n%s", calls)
		}
	})

	t.Run("declined writes do not call gh", func(t *testing.T) {
		output, calls := runWizardLibrary(t, library, "n\nn\n", `GITHUB_REPO="acme/widgets"
set_secret API_TOKEN "sensitive-test-value"
set_var REGION "west"
`)
		mustContain(t, output, "API_TOKEN in acme/widgets")
		mustContain(t, output, "REGION in acme/widgets")
		if strings.Contains(calls, "secret set") || strings.Contains(calls, "variable set") {
			t.Fatalf("declined GitHub write was attempted:\n%s", calls)
		}
	})

	t.Run("unauthenticated gh skips write", func(t *testing.T) {
		output, calls := runWizardLibrary(t, library, "", `export GH_AUTH_FAIL=1
GITHUB_REPO="acme/widgets"
set_secret API_TOKEN "sensitive-test-value"
`)
		mustContain(t, output, "authenticated gh missing")
		if strings.Contains(calls, "secret set") {
			t.Fatalf("unauthenticated GitHub write was attempted:\n%s", calls)
		}
	})

	t.Run("approved writes bind the repo and keep secrets off argv", func(t *testing.T) {
		output, calls := runWizardLibrary(t, library, "y\ny\n", `GITHUB_REPO="acme/widgets"
set_secret API_TOKEN "sensitive-test-value"
set_var REGION "west"
`)
		mustContain(t, calls, "secret set API_TOKEN --repo acme/widgets")
		mustContain(t, calls, "variable set REGION --body west --repo acme/widgets")
		mustContain(t, calls, "secret-stdin=present")
		mustNotContain(t, calls, "sensitive-test-value")
		mustNotContain(t, output, "sensitive-test-value")
	})
}

func TestWizardEnvWritesUsePrivatePermissions(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	templatePath := filepath.Join(repoRoot, "skills", "wizard", "template.sh")
	library := wizardLibrary(t, readFile(t, templatePath))

	tests := []struct {
		name    string
		initial string
		mode    os.FileMode
	}{
		{name: "new file"},
		{name: "existing permissive file", initial: "KEEP=1\nAPI_TOKEN=old\n", mode: 0o644},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			dir = physicalPath(t, dir)
			envPath := filepath.Join(dir, ".env")
			if tt.initial != "" {
				if err := os.WriteFile(envPath, []byte(tt.initial), tt.mode); err != nil {
					t.Fatal(err)
				}
			}
			runWizardLibrary(t, library, "", `ENV_FILE="$WIZARD_ENV_FILE"
write_env API_TOKEN "replacement"
`, "WIZARD_ENV_FILE="+envPath, "WIZARD_TEST_CWD="+dir)

			info, err := os.Stat(envPath)
			if err != nil {
				t.Fatal(err)
			}
			if got := info.Mode().Perm(); got != 0o600 {
				t.Fatalf("unexpected env permissions: got %04o want 0600", got)
			}
			contents := readFile(t, envPath)
			mustContain(t, contents, "API_TOKEN=replacement\n")
			if strings.Count(contents, "API_TOKEN=") != 1 {
				t.Fatalf("expected one API_TOKEN entry:\n%s", contents)
			}
			if tt.initial != "" {
				mustContain(t, contents, "KEEP=1\n")
			}
		})
	}

	for _, initial := range []struct {
		name    string
		content string
	}{
		{name: "new file"},
		{name: "existing file", content: "KEEP=1\nAPI_TOKEN=old\n"},
	} {
		t.Run("replacement failure cleans temporary files - "+initial.name, func(t *testing.T) {
			dir := t.TempDir()
			dir = physicalPath(t, dir)
			envPath := filepath.Join(dir, ".env")
			if initial.content != "" {
				if err := os.WriteFile(envPath, []byte(initial.content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			output, _ := runWizardLibrary(t, library, "", `
mv() { return 23; }
ENV_FILE="$WIZARD_ENV_FILE"
if write_env API_TOKEN "replacement"; then
  printf 'unexpected write_env success\n'
fi
printf 'written=%s\n' "${#WRITTEN_ENV[@]}"
`, "WIZARD_ENV_FILE="+envPath, "WIZARD_TEST_CWD="+dir)
			mustContain(t, output, "written=0")
			mustNotContain(t, output, "✓ wrote")

			matches, err := filepath.Glob(envPath + ".tmp.*")
			if err != nil {
				t.Fatal(err)
			}
			if len(matches) != 0 {
				t.Fatalf("temporary files remain after replacement failure: %v", matches)
			}
			if initial.content == "" {
				if _, err := os.Stat(envPath); !os.IsNotExist(err) {
					t.Fatalf("expected ENV_FILE to remain absent, stat error: %v", err)
				}
			} else if got := readFile(t, envPath); got != initial.content {
				t.Fatalf("ENV_FILE changed after replacement failure: got %q want %q", got, initial.content)
			}
		})
	}
}

func TestWizardRejectsReservedCaptureKeysWithoutMutatingState(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	templatePath := filepath.Join(repoRoot, "skills", "wizard", "template.sh")
	library := wizardLibrary(t, readFile(t, templatePath))

	reserved := []string{
		"PATH",
		"ENV_FILE",
		"GITHUB_REPO",
		"IFS",
		"CDPATH",
		"SHELLOPTS",
		"BASH_ENV",
		"BASHOPTS",
		"_WIZARD_ATTACK",
		"WRITTEN_ENV",
	}
	for _, key := range reserved {
		t.Run(key, func(t *testing.T) {
			base := physicalPath(t, t.TempDir())
			envPath := filepath.Join(base, ".env")
			initial := "PATH=existing-secret\nKEEP=1\n"
			if err := os.WriteFile(envPath, []byte(initial), 0o600); err != nil {
				t.Fatal(err)
			}

			output, _ := runWizardLibrary(t, library, "attacker-visible\nattacker-secret\n", `
ENV_FILE="$WIZARD_ENV_FILE"
PATH="$WIZARD_PATH"
GITHUB_REPO="acme/widgets"
IFS="$WIZARD_IFS"
CDPATH="$WIZARD_CDPATH"
BASH_ENV="$WIZARD_BASH_ENV"
before="PATH=$PATH ENV_FILE=$ENV_FILE GITHUB_REPO=$GITHUB_REPO IFS=$IFS CDPATH=$CDPATH BASH_ENV=$BASH_ENV SHELLOPTS=$SHELLOPTS BASE=$_WIZARD_AUTHORIZED_BASE CAPTURE=$_WIZARD_CAPTURED_VALUE"
if ask "$WIZARD_KEY" "Paste value:"; then
  printf 'unexpected ask acceptance\n'
fi
if ask_secret "$WIZARD_KEY" "Paste secret:"; then
  printf 'unexpected ask_secret acceptance\n'
fi
after="PATH=$PATH ENV_FILE=$ENV_FILE GITHUB_REPO=$GITHUB_REPO IFS=$IFS CDPATH=$CDPATH BASH_ENV=$BASH_ENV SHELLOPTS=$SHELLOPTS BASE=$_WIZARD_AUTHORIZED_BASE CAPTURE=$_WIZARD_CAPTURED_VALUE"
printf 'before=%s\n' "$before"
printf 'after=%s\n' "$after"
`, "WIZARD_KEY="+key, "WIZARD_ENV_FILE="+envPath, "WIZARD_PATH=/trusted/path", "WIZARD_IFS=trusted-ifs", "WIZARD_CDPATH=trusted-cdpath", "WIZARD_BASH_ENV=/trusted/bash-env", "WIZARD_TEST_CWD="+base)

			mustContain(t, output, "refusing reserved dotenv key: "+key)
			mustNotContain(t, output, "unexpected ask acceptance")
			mustNotContain(t, output, "unexpected ask_secret acceptance")
			mustNotContain(t, output, "existing-secret")
			mustNotContain(t, output, "attacker-visible")
			mustNotContain(t, output, "attacker-secret")
			before, after := "", ""
			for _, line := range strings.Split(output, "\n") {
				switch {
				case strings.HasPrefix(line, "before="):
					before = strings.TrimPrefix(line, "before=")
				case strings.HasPrefix(line, "after="):
					after = strings.TrimPrefix(line, "after=")
				}
			}
			if before == "" || after == "" {
				t.Fatalf("missing state snapshots in output:\n%s", output)
			}
			if before != after {
				t.Fatalf("reserved capture key changed shell/library state:\n%s\n%s", before, after)
			}
			if got := readFile(t, envPath); got != initial {
				t.Fatalf("reserved capture key changed existing ENV_FILE: got %q want %q", got, initial)
			}
		})
	}
}

func TestWizardRejectsEnvFileSymlinkBeforeReadOrWrite(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	templatePath := filepath.Join(repoRoot, "skills", "wizard", "template.sh")
	library := wizardLibrary(t, readFile(t, templatePath))

	dir := t.TempDir()
	dir = physicalPath(t, dir)
	targetPath := filepath.Join(dir, "target.env")
	envPath := filepath.Join(dir, ".env")
	initial := "API_TOKEN=symlink-target-secret\nKEEP=1\n"
	if err := os.WriteFile(targetPath, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(targetPath, envPath); err != nil {
		t.Fatal(err)
	}

	output, calls := runWizardLibrary(t, library, "\ny\n", `
ENV_FILE="$WIZARD_ENV_FILE"
GITHUB_REPO="acme/widgets"
if ask_secret API_TOKEN "Paste token:"; then
  set_secret API_TOKEN "$API_TOKEN"
else
  printf 'rejected\n'
fi
if write_env API_TOKEN "replacement"; then
  printf 'unexpected write\n'
fi
`, "WIZARD_ENV_FILE="+envPath, "WIZARD_TEST_CWD="+dir)
	mustContain(t, output, "refusing ENV_FILE symlink")
	mustContain(t, output, "rejected")
	mustNotContain(t, output, "symlink-target-secret")
	mustNotContain(t, calls, "symlink-target-secret")
	if got := readFile(t, targetPath); got != initial {
		t.Fatalf("symlink target changed: got %q want %q", got, initial)
	}
	linkTarget, err := os.Readlink(envPath)
	if err != nil {
		t.Fatal(err)
	}
	if linkTarget != targetPath {
		t.Fatalf("ENV_FILE symlink changed: got %q want %q", linkTarget, targetPath)
	}

	matches, err := filepath.Glob(envPath + ".tmp.*")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary files remain beside rejected symlink: %v", matches)
	}
}

func TestWizardRejectsUnsafeEnvPathsAndDotenvData(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	templatePath := filepath.Join(repoRoot, "skills", "wizard", "template.sh")
	library := wizardLibrary(t, readFile(t, templatePath))

	t.Run("parent symlink cannot expose or replace outside secret", func(t *testing.T) {
		base := t.TempDir()
		outside := t.TempDir()
		outsideEnv := filepath.Join(outside, ".env")
		initial := "API_TOKEN=outside-secret\nKEEP=1\n"
		if err := os.WriteFile(outsideEnv, []byte(initial), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(base, "linked")); err != nil {
			t.Fatal(err)
		}
		output, calls := runWizardLibrary(t, library, "\ny\n", `
ENV_FILE="linked/.env"
GITHUB_REPO="acme/widgets"
if ask_secret API_TOKEN "Paste token:"; then
  set_secret API_TOKEN "$API_TOKEN"
fi
if write_env API_TOKEN "replacement"; then printf 'unexpected write\n'; fi
`, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing ENV_FILE symlink component")
		mustNotContain(t, output, "outside-secret")
		mustNotContain(t, calls, "outside-secret")
		if got := readFile(t, outsideEnv); got != initial {
			t.Fatalf("outside env changed: got %q want %q", got, initial)
		}
	})

	for _, tt := range []struct {
		name       string
		linkTarget func(base, outside string) string
	}{
		{name: "parent symlink to absent outside target", linkTarget: func(_, outside string) string { return outside }},
		{name: "parent symlink stays inside base", linkTarget: func(base, _ string) string { return filepath.Join(base, "real") }},
		{name: "dangling parent symlink", linkTarget: func(_, outside string) string { return filepath.Join(outside, "missing") }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := physicalPath(t, t.TempDir())
			outside := physicalPath(t, t.TempDir())
			if err := os.Mkdir(filepath.Join(base, "real"), 0o755); err != nil {
				t.Fatal(err)
			}
			linkTarget := tt.linkTarget(base, outside)
			if err := os.Symlink(linkTarget, filepath.Join(base, "linked")); err != nil {
				t.Fatal(err)
			}
			output, _ := runWizardLibrary(t, library, "", `
ENV_FILE="linked/.env"
if write_env API_TOKEN "replacement"; then printf 'unexpected write\n'; fi
`, "WIZARD_TEST_CWD="+base)
			mustContain(t, output, "refusing ENV_FILE symlink component")
			mustNotContain(t, output, "unexpected write")
			if _, err := os.Stat(filepath.Join(linkTarget, ".env")); !os.IsNotExist(err) {
				t.Fatalf("rejected parent symlink created target: %v", err)
			}
		})
	}

	t.Run("dangling leaf symlink is rejected", func(t *testing.T) {
		base := physicalPath(t, t.TempDir())
		outside := physicalPath(t, t.TempDir())
		envPath := filepath.Join(base, ".env")
		if err := os.Symlink(filepath.Join(outside, "missing.env"), envPath); err != nil {
			t.Fatal(err)
		}
		output, _ := runWizardLibrary(t, library, "", `
ENV_FILE=".env"
if write_env API_TOKEN "replacement"; then printf 'unexpected write\n'; fi
`, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing ENV_FILE symlink")
		mustNotContain(t, output, "unexpected write")
		if _, err := os.Lstat(envPath); err != nil {
			t.Fatalf("dangling leaf symlink was replaced: %v", err)
		}
	})

	for _, tt := range []struct {
		name string
		path func(base, outside string) string
		want string
	}{
		{name: "parent traversal", path: func(_, _ string) string { return "nested/../.env" }, want: "refusing ENV_FILE traversal"},
		{name: "dot component", path: func(_, _ string) string { return "./.env" }, want: "refusing ENV_FILE traversal"},
		{name: "newline", path: func(_, _ string) string { return "bad\n.env" }, want: "refusing invalid ENV_FILE path"},
		{name: "missing parent", path: func(_, _ string) string { return "missing/.env" }, want: "refusing missing or non-directory ENV_FILE parent"},
		{name: "absolute outside", path: func(_, outside string) string { return filepath.Join(outside, ".env") }, want: "refusing ENV_FILE outside startup directory"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			outside := t.TempDir()
			path := tt.path(base, outside)
			output, _ := runWizardLibrary(t, library, "", `
ENV_FILE="$WIZARD_ENV_FILE"
if write_env API_TOKEN "replacement"; then printf 'unexpected write\n'; fi
`, "WIZARD_ENV_FILE="+path, "WIZARD_TEST_CWD="+base)
			mustContain(t, output, tt.want)
			mustNotContain(t, output, "unexpected write")
		})
	}

	t.Run("relative target remains anchored after cd", func(t *testing.T) {
		base := t.TempDir()
		afterCD := t.TempDir()
		if err := os.Mkdir(filepath.Join(base, "config"), 0o755); err != nil {
			t.Fatal(err)
		}
		runWizardLibrary(t, library, "", `
cd "$WIZARD_AFTER_CD"
ENV_FILE="config/.env"
write_env API_TOKEN "AZaz09_./:@%+,=?-"
`, "WIZARD_AFTER_CD="+afterCD, "WIZARD_TEST_CWD="+base)
		got := readFile(t, filepath.Join(base, "config", ".env"))
		if got != "API_TOKEN=AZaz09_./:@%+,=?-\n" {
			t.Fatalf("unexpected round-trip value: %q", got)
		}
		if _, err := os.Stat(filepath.Join(afterCD, "config", ".env")); !os.IsNotExist(err) {
			t.Fatalf("relative ENV_FILE followed later cd: %v", err)
		}
	})

	t.Run("ambient ENV_FILE cannot replace the authored default", func(t *testing.T) {
		base := t.TempDir()
		runWizardLibrary(t, library, "", `
write_env API_TOKEN "replacement"
`, "ENV_FILE=visible.env", "WIZARD_TEST_CWD="+base)
		if got := readFile(t, filepath.Join(base, ".env")); got != "API_TOKEN=replacement\n" {
			t.Fatalf("unexpected default ENV_FILE contents: %q", got)
		}
		if _, err := os.Stat(filepath.Join(base, "visible.env")); !os.IsNotExist(err) {
			t.Fatalf("ambient ENV_FILE changed destination: %v", err)
		}
	})

	t.Run("Git worktree rejects non-ignored authored target", func(t *testing.T) {
		base := t.TempDir()
		cmd := exec.Command("git", "init", "-q", base)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git init: %v\n%s", err, output)
		}
		if err := os.WriteFile(filepath.Join(base, ".gitignore"), []byte(".env\n.env.tmp.*\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		output, _ := runWizardLibrary(t, library, "", `
ENV_FILE="visible.env"
if write_env API_TOKEN "replacement"; then printf 'unexpected ignored-path success\n'; fi
`, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing non-ignored ENV_FILE in Git worktree")
		mustNotContain(t, output, "unexpected")
		if _, err := os.Stat(filepath.Join(base, "visible.env")); !os.IsNotExist(err) {
			t.Fatalf("non-ignored ENV_FILE was created: %v", err)
		}
		matches, err := filepath.Glob(filepath.Join(base, "visible.env.tmp.*"))
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != 0 {
			t.Fatalf("non-ignored target created temporary files: %v", matches)
		}
	})

	t.Run("Git worktree requires ignored temporary path", func(t *testing.T) {
		base := t.TempDir()
		cmd := exec.Command("git", "init", "-q", base)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git init: %v\n%s", err, output)
		}
		if err := os.WriteFile(filepath.Join(base, ".gitignore"), []byte("visible.env\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		output, _ := runWizardLibrary(t, library, "", `
ENV_FILE="visible.env"
if write_env API_TOKEN "replacement"; then printf 'unexpected temporary-ignore success\n'; fi
`, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing non-ignored ENV_FILE temporary path in Git worktree")
		mustNotContain(t, output, "unexpected")
		if _, err := os.Stat(filepath.Join(base, "visible.env")); !os.IsNotExist(err) {
			t.Fatalf("temporary-ignore failure created ENV_FILE: %v", err)
		}
	})

	t.Run("Git worktree verifies the actual random temporary path", func(t *testing.T) {
		base := t.TempDir()
		cmd := exec.Command("git", "init", "-q", base)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git init: %v\n%s", err, output)
		}
		if err := os.WriteFile(filepath.Join(base, ".gitignore"), []byte("visible.env\nvisible.env.tmp.probe\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		output, _ := runWizardLibrary(t, library, "", `
ENV_FILE="visible.env"
if write_env API_TOKEN "replacement"; then printf 'unexpected actual-temporary-ignore success\n'; fi
`, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing non-ignored ENV_FILE actual temporary path in Git worktree")
		mustNotContain(t, output, "unexpected")
		if _, err := os.Stat(filepath.Join(base, "visible.env")); !os.IsNotExist(err) {
			t.Fatalf("actual temporary ignore failure created ENV_FILE: %v", err)
		}
		matches, err := filepath.Glob(filepath.Join(base, "visible.env.tmp.*"))
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != 0 {
			t.Fatalf("actual non-ignored temporary files remain: %v", matches)
		}
	})

	t.Run("Git worktree writes when final and temporary paths are ignored", func(t *testing.T) {
		base := t.TempDir()
		cmd := exec.Command("git", "init", "-q", base)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git init: %v\n%s", err, output)
		}
		if err := os.WriteFile(filepath.Join(base, ".gitignore"), []byte("visible.env\nvisible.env.tmp.*\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		runWizardLibrary(t, library, "", `
ENV_FILE="visible.env"
write_env API_TOKEN "replacement"
`, "WIZARD_TEST_CWD="+base)
		envPath := filepath.Join(base, "visible.env")
		if got := readFile(t, envPath); got != "API_TOKEN=replacement\n" {
			t.Fatalf("unexpected ignored ENV_FILE contents: %q", got)
		}
		info, err := os.Stat(envPath)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Fatalf("unexpected ignored ENV_FILE permissions: got %04o want 0600", got)
		}
	})

	t.Run("nested Git worktree owns its ignore decision", func(t *testing.T) {
		base := t.TempDir()
		cmd := exec.Command("git", "init", "-q", base)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("outer git init: %v\n%s", err, output)
		}
		nested := filepath.Join(base, "nested")
		if err := os.Mkdir(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command("git", "init", "-q", nested)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("nested git init: %v\n%s", err, output)
		}
		if err := os.WriteFile(filepath.Join(base, ".gitignore"), []byte("nested/\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		output, _ := runWizardLibrary(t, library, "", `
ENV_FILE="nested/.env"
if write_env API_TOKEN "replacement"; then printf 'unexpected nested-worktree success\n'; fi
`, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing non-ignored ENV_FILE in Git worktree")
		mustNotContain(t, output, "unexpected")
		if _, err := os.Stat(filepath.Join(nested, ".env")); !os.IsNotExist(err) {
			t.Fatalf("nested worktree ENV_FILE was created: %v", err)
		}
		matches, err := filepath.Glob(filepath.Join(nested, ".env.tmp.*"))
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != 0 {
			t.Fatalf("nested worktree temporary files remain: %v", matches)
		}
	})

	t.Run("nested Git worktree ignores ambient outer Git overrides", func(t *testing.T) {
		base := t.TempDir()
		cmd := exec.Command("git", "init", "-q", base)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("outer git init: %v\n%s", err, output)
		}
		nested := filepath.Join(base, "nested")
		if err := os.Mkdir(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command("git", "init", "-q", nested)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("nested git init: %v\n%s", err, output)
		}
		if err := os.WriteFile(filepath.Join(base, ".gitignore"), []byte("nested/\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		output, _ := runWizardLibrary(t, library, "", `
export GIT_DIR="$WIZARD_OUTER_GIT_DIR"
export GIT_WORK_TREE="$WIZARD_OUTER_WORK_TREE"
ENV_FILE="nested/.env"
if write_env API_TOKEN "replacement"; then printf 'unexpected ambient-git success\n'; fi
`, "WIZARD_OUTER_GIT_DIR="+filepath.Join(base, ".git"), "WIZARD_OUTER_WORK_TREE="+base, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing non-ignored ENV_FILE in Git worktree")
		mustNotContain(t, output, "unexpected")
		if _, err := os.Stat(filepath.Join(nested, ".env")); !os.IsNotExist(err) {
			t.Fatalf("ambient Git override created nested ENV_FILE: %v", err)
		}
		matches, err := filepath.Glob(filepath.Join(nested, ".env.tmp.*"))
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != 0 {
			t.Fatalf("ambient Git override left nested temporary files: %v", matches)
		}
	})

	t.Run("invalid keys and values fail before temporary creation", func(t *testing.T) {
		base := t.TempDir()
		output, _ := runWizardLibrary(t, library, "", `
ENV_FILE=".env"
if write_env 'BAD-KEY' "safe"; then printf 'unexpected key success\n'; fi
if write_env 'API.*' "safe"; then printf 'unexpected regex key success\n'; fi
if write_env API_TOKEN 'has space'; then printf 'unexpected value success\n'; fi
if write_env API_TOKEN 'has#comment'; then printf 'unexpected comment success\n'; fi
if write_env API_TOKEN 'has$dollar'; then printf 'unexpected dollar success\n'; fi
if write_env API_TOKEN 'has\backslash'; then printf 'unexpected backslash success\n'; fi
if write_env API_TOKEN 'ok;/usr/bin/id'; then printf 'unexpected semicolon success\n'; fi
if write_env API_TOKEN 'ok&/usr/bin/id'; then printf 'unexpected ampersand success\n'; fi
if write_env API_TOKEN $'has\nnewline'; then printf 'unexpected newline success\n'; fi
if write_env API_TOKEN $'has\rcarriage-return'; then printf 'unexpected carriage-return success\n'; fi
`, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing invalid dotenv key")
		mustContain(t, output, "refusing unsafe dotenv value")
		mustNotContain(t, output, "unexpected")
		if _, err := os.Stat(filepath.Join(base, ".env")); !os.IsNotExist(err) {
			t.Fatalf("invalid dotenv input created ENV_FILE: %v", err)
		}
		matches, err := filepath.Glob(filepath.Join(base, ".env.tmp.*"))
		if err != nil {
			t.Fatal(err)
		}
		if len(matches) != 0 {
			t.Fatalf("invalid dotenv input created temporary files: %v", matches)
		}
	})

	t.Run("unsafe existing value is not exposed to GitHub", func(t *testing.T) {
		base := t.TempDir()
		envPath := filepath.Join(base, ".env")
		if err := os.WriteFile(envPath, []byte("API_TOKEN=outside secret\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		output, calls := runWizardLibrary(t, library, "\ny\n", `
ENV_FILE=".env"
GITHUB_REPO="acme/widgets"
if ask_secret API_TOKEN "Paste token:"; then
  set_secret API_TOKEN "$API_TOKEN"
fi
`, "WIZARD_TEST_CWD="+base)
		mustContain(t, output, "refusing unsafe existing dotenv value")
		mustNotContain(t, output, "outside secret")
		mustNotContain(t, calls, "outside secret")
		mustNotContain(t, calls, "secret set")
	})
}

func wizardLibrary(t *testing.T, template string) string {
	t.Helper()
	marker := "\n# STAGES — author this section."
	index := strings.Index(template, marker)
	if index < 0 {
		t.Fatalf("wizard stages marker not found")
	}
	return template[:index] + "\n"
}

func runWizardLibrary(t *testing.T, library, input, invocation string, extraEnv ...string) (string, string) {
	t.Helper()
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is unavailable")
	}
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	callsPath := filepath.Join(dir, "gh-calls.log")
	ghStub := `#!/bin/sh
printf '%s\n' "$*" >> "$GH_CALLS"
if [ "$1 $2" = "auth status" ]; then
  if [ "${GH_AUTH_FAIL:-0}" = "1" ]; then
    exit 1
  fi
  exit 0
fi
if [ "$1 $2" = "secret set" ]; then
  IFS= read -r _secret || true
  if [ -n "$_secret" ]; then
    printf 'secret-stdin=present\n' >> "$GH_CALLS"
  fi
fi
`
	ghPath := filepath.Join(binDir, "gh")
	if err := os.WriteFile(ghPath, []byte(ghStub), 0o755); err != nil {
		t.Fatal(err)
	}
	scriptPath := filepath.Join(dir, "wizard-library-test.sh")
	if err := os.WriteFile(scriptPath, []byte(library+"\n"+invocation), 0o700); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bash, scriptPath)
	for _, item := range extraEnv {
		if strings.HasPrefix(item, "WIZARD_TEST_CWD=") {
			cmd.Dir = strings.TrimPrefix(item, "WIZARD_TEST_CWD=")
		}
	}
	cmd.Env = append(os.Environ(), "PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"), "GH_CALLS="+callsPath)
	cmd.Env = append(cmd.Env, extraEnv...)
	cmd.Stdin = strings.NewReader(input)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("wizard library failed: %v\n%s", err, output.String())
	}
	calls, err := os.ReadFile(callsPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return output.String(), string(calls)
}

func physicalPath(t *testing.T, path string) string {
	t.Helper()
	physical, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return physical
}

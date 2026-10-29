package integration_test

import (
	"os"
	"strings"
	"testing"
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// Failure messages quote an excerpt, never the whole document, so one failed
// contract does not flood the output with a flattened SKILL.md.
func mustContain(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		t.Fatalf("expected %q in text (%d bytes) starting:\n%.300s", want, len(text), text)
	}
}

func mustNotContain(t *testing.T, text, unwanted string) {
	t.Helper()
	if i := strings.Index(text, unwanted); i >= 0 {
		start, end := max(0, i-150), min(len(text), i+len(unwanted)+150)
		t.Fatalf("did not expect %q at byte %d:\n...%s...", unwanted, i, text[start:end])
	}
}

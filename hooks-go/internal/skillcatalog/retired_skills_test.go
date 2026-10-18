package skillcatalog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseRetiredSkillsRequiresCanonicalDocument(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "missing field", raw: `{}`, want: "missing retired_skills"},
		{name: "unknown field", raw: `{"retired_skills": [], "extra": true}`, want: "unknown field"},
		{name: "duplicate top-level field", raw: `{"retired_skills": ["demo"], "retired_skills": []}`, want: "duplicate top-level field"},
		{name: "trailing json", raw: `{"retired_skills": []} {}`, want: "trailing JSON"},
		{name: "invalid name", raw: `{"retired_skills": ["Bad"]}`, want: "lowercase"},
		{name: "duplicate", raw: `{"retired_skills": ["demo", "demo"]}`, want: "duplicate"},
		{name: "unsorted", raw: `{"retired_skills": ["zulu", "alpha"]}`, want: "sorted"},
		{name: "mixed v1 and v2", raw: `{"retired_skills": [], "schemaVersion": 2}`, want: "cannot be mixed"},
		{name: "v2 wrong schema version", raw: `{"schemaVersion": 1, "updatedAt": "2026-10-18", "retiredSkills": []}`, want: "schemaVersion 2"},
		{name: "v2 missing updatedAt", raw: `{"schemaVersion": 2, "retiredSkills": []}`, want: "updatedAt"},
		{name: "v2 missing entry key", raw: `{"schemaVersion": 2, "updatedAt": "x", "retiredSkills": [{"name": "demo", "retiredIn": "v1", "reason": "r"}]}`, want: `missing "replacedBy"`},
		{name: "v2 unknown entry key", raw: `{"schemaVersion": 2, "updatedAt": "x", "retiredSkills": [{"name": "demo", "retiredIn": "v1", "replacedBy": null, "reason": "r", "extra": 1}]}`, want: "unknown field"},
		{name: "v2 empty reason", raw: `{"schemaVersion": 2, "updatedAt": "x", "retiredSkills": [{"name": "demo", "retiredIn": "v1", "replacedBy": null, "reason": " "}]}`, want: "non-empty"},
		{name: "v2 invalid replacement", raw: `{"schemaVersion": 2, "updatedAt": "x", "retiredSkills": [{"name": "demo", "retiredIn": "v1", "replacedBy": "Bad", "reason": "r"}]}`, want: "lowercase"},
		{name: "v2 unsorted", raw: `{"schemaVersion": 2, "updatedAt": "x", "retiredSkills": [{"name": "zulu", "retiredIn": "v1", "replacedBy": null, "reason": "r"}, {"name": "alpha", "retiredIn": "v1", "replacedBy": "zulu", "reason": "r"}]}`, want: "sorted"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ParseRetiredSkills([]byte(test.raw)); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ParseRetiredSkills() error = %v, want %q", err, test.want)
			}
		})
	}
	if got, err := ParseRetiredSkills([]byte(`{"retired_skills": []}`)); err != nil || len(got) != 0 {
		t.Fatalf("empty retired list = %#v, %v", got, err)
	}
	v2 := `{"schemaVersion": 2, "updatedAt": "2026-10-18", "retiredSkills": [
		{"name": "alpha", "retiredIn": "v1", "replacedBy": "omega", "reason": "merged"},
		{"name": "beta", "retiredIn": "v1", "replacedBy": null, "reason": "dropped"}]}`
	if got, err := ParseRetiredSkills([]byte(v2)); err != nil || strings.Join(got, ",") != "alpha,beta" {
		t.Fatalf("v2 retired list = %#v, %v", got, err)
	}
}

func TestParseRetiredSkillsRejectsOversizedManifest(t *testing.T) {
	raw := []byte(`{"retired_skills": []}`)
	raw = append(raw, []byte(strings.Repeat(" ", maxRetiredSkillsManifestBytes-len(raw)+1))...)

	if _, err := ParseRetiredSkills(raw); err == nil || !strings.Contains(err.Error(), "64 KiB") {
		t.Fatalf("expected oversized manifest failure, got %v", err)
	}
}

func TestLoadRetiredSkillsRejectsOversizedManifest(t *testing.T) {
	path := filepath.Join(t.TempDir(), retiredSkillsFile)
	raw := []byte(`{"retired_skills": []}`)
	raw = append(raw, []byte(strings.Repeat(" ", maxRetiredSkillsManifestBytes-len(raw)+1))...)
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadRetiredSkills(path); err == nil || !strings.Contains(err.Error(), "64 KiB") {
		t.Fatalf("expected oversized manifest load failure, got %v", err)
	}
}

func TestParseRetiredSkillsRejectsTooManyNames(t *testing.T) {
	names := make([]string, maxRetiredSkillNames+1)
	for index := range names {
		names[index] = fmt.Sprintf("skill-%03d", index)
	}
	raw, err := json.Marshal(struct {
		RetiredSkills []string `json:"retired_skills"`
	}{RetiredSkills: names})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ParseRetiredSkills(raw); err == nil || !strings.Contains(err.Error(), "128-name limit") {
		t.Fatalf("expected retired name count failure, got %v", err)
	}
}

func TestValidateRepositoryRetirementsRejectsActiveSkillAndRoutingReference(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, "skills")
	if err := os.MkdirAll(filepath.Join(skills, "demo-skill"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skills, "demo-skill", "SKILL.md"), []byte("---\nname: demo-skill\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("Use `retired-skill` only in history."), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skills, "demo-skill", "notes.md"), []byte("retired-skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := ValidateRepositoryRetirements(repo, skills, []string{"retired-skill"})
	if err == nil || !strings.Contains(err.Error(), "retired-skill") {
		t.Fatalf("expected retirement validation failure, got %v", err)
	}

	if err := os.Remove(filepath.Join(repo, "README.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(skills, "demo-skill", "notes.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(skills, "retired-skill"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skills, "retired-skill", "SKILL.md"), []byte("active"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateRepositoryRetirements(repo, skills, []string{"retired-skill"}); err == nil || !strings.Contains(err.Error(), "still active") {
		t.Fatalf("expected active-skill failure, got %v", err)
	}
}

func TestValidateRepositoryRetirementsRejectsUpstreamRetiredReference(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, "skills")
	if err := os.MkdirAll(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	upstream := filepath.Join(skills, "UPSTREAM.md")
	if err := os.WriteFile(upstream, []byte("The retired-skill route was removed.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidateRepositoryRetirements(repo, skills, []string{"retired-skill"})
	if err == nil || !strings.Contains(err.Error(), upstream) {
		t.Fatalf("expected UPSTREAM retirement reference failure, got %v", err)
	}
}

func TestValidateRepositoryRetirementsRejectsOversizedRoutingSurface(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, "skills")
	if err := os.MkdirAll(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := filepath.Join(repo, "README.md")
	if err := os.WriteFile(readme, []byte(strings.Repeat("x", maxRoutingSurfaceBytes+1)), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidateRepositoryRetirements(repo, skills, []string{"retired-skill"})
	if err == nil || !strings.Contains(err.Error(), "exceeds 1 MiB limit") {
		t.Fatalf("expected oversized routing surface failure, got %v", err)
	}
}

func TestActiveRoutingPathsRejectsTooManyFiles(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, "skills")
	if err := os.MkdirAll(filepath.Join(repo, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	for index := 0; index <= maxRoutingSurfaceFiles; index++ {
		path := filepath.Join(repo, "config", fmt.Sprintf("route-%04d.md", index))
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := activeRoutingPaths(repo, skills, []string{"retired-skill"}); err == nil ||
		!strings.Contains(err.Error(), "file count") || !strings.Contains(err.Error(), "1024") {
		t.Fatalf("expected routing file count failure, got %v", err)
	}
}

func TestValidateRepositoryRetirementsRejectsAggregateRoutingSurface(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, "skills")
	config := filepath.Join(repo, "config")
	if err := os.MkdirAll(config, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	for index := 0; index <= 8; index++ {
		path := filepath.Join(config, fmt.Sprintf("route-%02d.md", index))
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Truncate(path, int64(maxRoutingSurfaceBytes)); err != nil {
			t.Fatal(err)
		}
	}

	err := ValidateRepositoryRetirements(repo, skills, []string{"retired-skill"})
	if err == nil || !strings.Contains(err.Error(), "8 MiB aggregate limit") {
		t.Fatalf("expected aggregate routing surface failure, got %v", err)
	}
}

func TestValidateSkillNameUsesSharedCanonicalRules(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{name: "demo-skill", valid: true},
		{name: "", valid: false},
		{name: "Demo-skill", valid: false},
		{name: "-demo", valid: false},
		{name: "demo-", valid: false},
		{name: "demo--skill", valid: false},
		{name: strings.Repeat("a", maxSkillNameLength+1), valid: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateSkillName(test.name)
			if (err == nil) != test.valid {
				t.Fatalf("ValidateSkillName(%q) error = %v, valid = %v", test.name, err, test.valid)
			}
		})
	}
}

func TestRepositoryRetirementManifestIsValid(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", "..", ".."))
	names, err := LoadRetiredSkills(RetiredSkillsPath(repoRoot))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateRepositoryRetirements(repoRoot, filepath.Join(repoRoot, "skills"), names); err != nil {
		t.Fatal(err)
	}
}

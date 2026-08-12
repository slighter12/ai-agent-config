package skilltools

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-agent-config/hooks/internal/skillcatalog"
)

func TestValidateSkillAcceptsPortableFrontmatter(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform demo validation. Use when checking the Go validator. Avoid when another skill is more specific.
argument-hint: "What should be validated?"
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  invocation: model
  version: "0.1.0"
---

# Demo Skill
`)

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if !valid || message != "Skill is valid" {
		t.Fatalf("expected valid skill, got valid=%v message=%q", valid, message)
	}
}

func TestValidateSkillAcceptsDescriptionWithoutPrescribedPhrases(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Diagnose intermittent failures through falsifiable evidence while leaving implementation to the owning workflow.
metadata:
  invocation: model
---

# Demo Skill
`)

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if !valid || message != "Skill is valid" {
		t.Fatalf("expected flexible description wording to be valid, got valid=%v message=%q", valid, message)
	}
}

func TestValidateSkillAcceptsLongSinglePathSkill(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "long-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: long-skill
description: Preserve a single workflow path. Use when shared guidance remains relevant throughout the skill. Avoid when branch-specific material belongs in a reference.
metadata:
  invocation: model
---

# Long Skill

`+strings.Repeat("Shared workflow guidance remains on the single path.\n", 510))

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if !valid || message != "Skill is valid" {
		t.Fatalf("expected valid long skill, got valid=%v message=%q", valid, message)
	}
}

func TestValidateSkillRejectsMissingInvocation(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform demo validation. Use when checking invocation contracts. Avoid when another skill is more specific.
metadata:
  version: "0.1.0"
---

# Demo Skill
`)

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected invalid skill")
	}
	mustContain(t, message, "metadata.invocation must be user or model")
}

func TestValidateSkillRejectsMissingUserInvocationControls(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform demo validation. Use when checking invocation controls. Avoid when another skill is more specific.
metadata:
  invocation: user
---

# Demo Skill
`)

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected invalid skill")
	}
	mustContain(t, message, "OpenCode")
	mustContain(t, message, "Codex")
	mustContain(t, message, "Claude")
}

func TestValidateSkillRejectsUserControlsOnModelSkill(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(filepath.Join(skillDir, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform demo validation. Use when checking model invocation controls. Avoid when another skill is more specific.
disable-model-invocation: true
metadata:
  invocation: model
  opencode/autoinvoke: "false"
---

# Demo Skill
`)
	writeFile(t, filepath.Join(skillDir, "agents", "openai.yaml"), "policy:\n  allow_implicit_invocation: false\n")

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected invalid skill")
	}
	mustContain(t, message, "model skill must not disable OpenCode")
	mustContain(t, message, "model skill must not disable Codex")
	mustContain(t, message, "model skill must not disable Claude")
}

func TestValidateSkillRejectsPlaceholderDescription(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform [capability]. Use when testing. Avoid when done.
---

# Demo Skill
`)

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected invalid skill")
	}
	if !strings.Contains(message, "placeholder") {
		t.Fatalf("expected placeholder error, got %q", message)
	}
}

func TestInitSkillCreatesTemplateFiles(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer
	if err := InitSkill(&out, "demo-skill", dir); err != nil {
		t.Fatal(err)
	}
	skillDir := filepath.Join(dir, "demo-skill")
	mustContain(t, readFile(t, filepath.Join(skillDir, "SKILL.md")), "name: demo-skill")
	mustContain(t, readFile(t, filepath.Join(skillDir, "SKILL.md")), "# Demo Skill")
	mustContain(t, readFile(t, filepath.Join(skillDir, "SKILL.md")), `invocation: "[user|model]"`)
	mustContain(t, out.String(), "Validate "+skillDir)
	mustContain(t, out.String(), "same agent-config entrypoint used for initialization")
	if strings.Contains(out.String(), "Run agent-config "+"validate-skill") {
		t.Fatalf("init handoff contains a bare unavailable command:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(skillDir, "references")); !os.IsNotExist(err) {
		t.Fatalf("expected no references directory, got err=%v", err)
	}
}

func TestInitSkillUsesSharedNameRules(t *testing.T) {
	for _, name := range []string{"", "Upper-case", "-leading", "trailing-", "double--dash", strings.Repeat("a", 65)} {
		t.Run(name, func(t *testing.T) {
			sharedErr := skillcatalog.ValidateSkillName(name)
			initErr := InitSkill(&bytes.Buffer{}, name, t.TempDir())
			if (sharedErr == nil) != (initErr == nil) {
				t.Fatalf("shared rule error = %v, InitSkill error = %v", sharedErr, initErr)
			}
		})
	}
}

func TestValidateSkillsAcceptsRepositoryCatalog(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", "..", ".."))
	valid, message, err := ValidateSkills(filepath.Join(repoRoot, "skills"))
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatalf("expected repository catalog to be valid:\n%s", message)
	}
	mustContain(t, message, "26 skills (11 model, 15 user)")
	mustContain(t, message, "repo model name/description contribution characters")
}

func TestValidateSkillsBudgetsOnlyModelDiscovery(t *testing.T) {
	longDescription := "Measure model-facing discovery. Use when " +
		strings.Repeat("catalog routing detail ", 38) +
		"Avoid when explicit invocation keeps the skill out of model discovery."

	t.Run("user descriptions are excluded", func(t *testing.T) {
		skillsDir := t.TempDir()
		var userNames []string
		for index := 0; index < 10; index++ {
			name := "user-" + string(rune('a'+index))
			writeInvocationSkillWithDescription(t, skillsDir, name, "user", longDescription)
			userNames = append(userNames, name)
		}
		writeRouterCatalog(t, skillsDir, userNames, nil)

		valid, message, err := ValidateSkills(skillsDir)
		if err != nil {
			t.Fatal(err)
		}
		if !valid {
			t.Fatalf("expected user-only catalog to be valid:\n%s", message)
		}
		mustContain(t, message, "11 skills (0 model, 11 user)")
		mustContain(t, message, "0/8000 repo model name/description contribution characters")
	})

	t.Run("model descriptions enforce the limit", func(t *testing.T) {
		skillsDir := t.TempDir()
		var modelNames []string
		for index := 0; index < 10; index++ {
			name := "model-" + string(rune('a'+index))
			writeInvocationSkillWithDescription(t, skillsDir, name, "model", longDescription)
			modelNames = append(modelNames, name)
		}
		writeRouterCatalog(t, skillsDir, nil, modelNames)

		valid, message, err := ValidateSkills(skillsDir)
		if err != nil {
			t.Fatal(err)
		}
		if valid {
			t.Fatal("expected model discovery budget failure")
		}
		mustContain(t, message, "repo model skill names and descriptions use")
		mustContain(t, message, "maximum is 8000")
	})
}

func TestValidateSkillsRejectsMissingRouter(t *testing.T) {
	skillsDir := t.TempDir()
	writeInvocationSkill(t, skillsDir, "demo-model", "model")

	valid, message, err := ValidateSkills(skillsDir)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected missing router to invalidate the catalog")
	}
	mustContain(t, message, "required router skill ask-matt is missing")
}

func TestValidateSkillsRejectsMissingUserInvocationControls(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform a demo workflow. Use when validating invocation controls. Avoid when another skill is more specific.
metadata:
  invocation: user
---

# Demo Skill
`)
	valid, message, err := ValidateSkills(dir)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected invalid catalog")
	}
	mustContain(t, message, "OpenCode")
	mustContain(t, message, "Codex")
}

func TestValidateSkillsChecksCanonicalCatalogReference(t *testing.T) {
	tests := []struct {
		name        string
		catalog     string
		wantValid   bool
		wantMessage string
	}{
		{
			name: "valid",
			catalog: `# Skill Catalog

## User-invoked

| Skill | Purpose |
| --- | --- |
| ` + "`ask-matt`" + ` | Explain routes. |

## Model-invoked

| Skill | Purpose |
| --- | --- |
| ` + "`demo-model`" + ` | Demonstrate model routing. |
`,
			wantValid: true,
		},
		{
			name: "missing",
			catalog: `# Skill Catalog

## User-invoked

| Skill | Purpose |
| --- | --- |
| ` + "`ask-matt`" + ` | Explain routes. |

## Model-invoked
`,
			wantMessage: "catalog is missing demo-model",
		},
		{
			name: "unknown",
			catalog: `# Skill Catalog

## User-invoked

| Skill | Purpose |
| --- | --- |
| ` + "`ask-matt`" + ` | Explain routes. |

## Model-invoked

| Skill | Purpose |
| --- | --- |
| ` + "`demo-model`" + ` | Demonstrate model routing. |
| ` + "`ghost-skill`" + ` | Unknown route. |
`,
			wantMessage: "catalog contains unknown skill ghost-skill",
		},
		{
			name: "duplicate",
			catalog: `# Skill Catalog

## User-invoked

| Skill | Purpose |
| --- | --- |
| ` + "`ask-matt`" + ` | Explain routes. |
| ` + "`ask-matt`" + ` | Explain routes twice. |

## Model-invoked

| Skill | Purpose |
| --- | --- |
| ` + "`demo-model`" + ` | Demonstrate model routing. |
`,
			wantMessage: "catalog lists ask-matt more than once",
		},
		{
			name: "wrong group",
			catalog: `# Skill Catalog

## User-invoked

| Skill | Purpose |
| --- | --- |
| ` + "`ask-matt`" + ` | Explain routes. |
| ` + "`demo-model`" + ` | Demonstrate model routing. |

## Model-invoked
`,
			wantMessage: "catalog groups demo-model as user; metadata invocation is model",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			skillsDir := createCatalogFixture(t, test.catalog)
			valid, message, err := ValidateSkills(skillsDir)
			if err != nil {
				t.Fatal(err)
			}
			if valid != test.wantValid {
				t.Fatalf("got valid=%v, want %v:\n%s", valid, test.wantValid, message)
			}
			if test.wantMessage != "" {
				mustContain(t, message, test.wantMessage)
			}
		})
	}
}

func TestValidateSkillsRejectsMissingCanonicalCatalogReference(t *testing.T) {
	skillsDir := t.TempDir()
	writeInvocationSkill(t, skillsDir, "ask-matt", "user")
	writeInvocationSkill(t, skillsDir, "demo-model", "model")

	valid, message, err := ValidateSkills(skillsDir)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected invalid catalog")
	}
	mustContain(t, message, "canonical catalog reference not found")
}

func TestValidatorRejectsUnsafeAndOversizedInputs(t *testing.T) {
	t.Run("symlinked SKILL.md", func(t *testing.T) {
		root := t.TempDir()
		external := filepath.Join(t.TempDir(), "outside.md")
		writeFile(t, external, "---\nname: demo\ndescription: Demo.\nmetadata:\n  invocation: model\n---\n")
		skillDir := filepath.Join(root, "demo")
		if err := os.Mkdir(skillDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(external, filepath.Join(skillDir, "SKILL.md")); err != nil {
			t.Fatal(err)
		}
		if _, _, err := ValidateSkill(skillDir); err == nil || !strings.Contains(err.Error(), "is a symlink") {
			t.Fatalf("expected symlink rejection, got %v", err)
		}
	})

	t.Run("oversized SKILL.md", func(t *testing.T) {
		skillDir := filepath.Join(t.TempDir(), "demo")
		if err := os.Mkdir(skillDir, 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(skillDir, "SKILL.md"), strings.Repeat("x", maxValidatorFileBytes+1))
		if _, _, err := ValidateSkill(skillDir); err == nil || !strings.Contains(err.Error(), "exceeds 1 MiB") {
			t.Fatalf("expected size rejection, got %v", err)
		}
	})

	t.Run("aggregate budget", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "input")
		writeFile(t, path, "12345")
		budget := &validatorBudget{remaining: 4}
		if _, err := readValidatorFile(path, budget); err == nil || !strings.Contains(err.Error(), "aggregate limit") {
			t.Fatalf("expected aggregate rejection, got %v", err)
		}
	})

	t.Run("entry count", func(t *testing.T) {
		root := t.TempDir()
		for i := 0; i <= maxValidatorSkillEntries; i++ {
			if err := os.Mkdir(filepath.Join(root, fmt.Sprintf("entry-%04d", i)), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := readSkillEntries(root); err == nil || !strings.Contains(err.Error(), "more than") {
			t.Fatalf("expected entry-count rejection, got %v", err)
		}
	})

	t.Run("symlinked sidecar parent", func(t *testing.T) {
		root := t.TempDir()
		writeInvocationSkill(t, root, "demo-user", "user")
		agents := filepath.Join(root, "demo-user", "agents")
		if err := os.RemoveAll(agents); err != nil {
			t.Fatal(err)
		}
		external := t.TempDir()
		writeFile(t, filepath.Join(external, "openai.yaml"), "policy:\n  allow_implicit_invocation: false\n")
		if err := os.Symlink(external, agents); err != nil {
			t.Fatal(err)
		}
		valid, message, err := ValidateSkill(filepath.Join(root, "demo-user"))
		if err != nil {
			t.Fatal(err)
		}
		if valid || !strings.Contains(message, "parent") || !strings.Contains(message, "symlink") {
			t.Fatalf("expected sidecar parent-symlink rejection, got valid=%v message=%q", valid, message)
		}
	})

	t.Run("symlinked catalog parent", func(t *testing.T) {
		root := createCatalogFixture(t, "# Skill Catalog\n\n## User-invoked\n\n| Skill | Purpose |\n| --- | --- |\n| `ask-matt` | Router. |\n\n## Model-invoked\n\n| Skill | Purpose |\n| --- | --- |\n| `demo-model` | Model. |\n")
		references := filepath.Join(root, "ask-matt", "references")
		if err := os.RemoveAll(references); err != nil {
			t.Fatal(err)
		}
		external := t.TempDir()
		writeFile(t, filepath.Join(external, "CATALOG.md"), "# Skill Catalog\n")
		if err := os.Symlink(external, references); err != nil {
			t.Fatal(err)
		}
		valid, message, err := ValidateSkills(root)
		if err != nil {
			t.Fatal(err)
		}
		if valid || !strings.Contains(message, "parent") || !strings.Contains(message, "symlink") {
			t.Fatalf("expected catalog parent-symlink rejection, got valid=%v message=%q", valid, message)
		}
	})
}

func TestPackageSkillCreatesSkillArchive(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(filepath.Join(skillDir, "references"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform demo packaging. Use when packaging a test skill. Avoid when validation should fail.
metadata:
  invocation: model
  version: "0.1.0"
---

# Demo Skill
`)
	writeFile(t, filepath.Join(skillDir, "references", "INDEX.md"), "# References\n")
	outDir := filepath.Join(dir, "dist")

	var out bytes.Buffer
	archivePath, err := PackageSkill(&out, skillDir, outDir)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	names := map[string]bool{}
	for _, file := range reader.File {
		names[file.Name] = true
	}
	for _, name := range []string{"demo-skill/SKILL.md", "demo-skill/references/INDEX.md"} {
		if !names[name] {
			t.Fatalf("expected archive entry %s, got %#v", name, names)
		}
	}
}

func TestPackageSkillRejectsMissingInvocation(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform demo packaging. Use when checking package validation. Avoid when another skill is more specific.
metadata:
  version: "0.1.0"
---

# Demo Skill
`)
	outDir := filepath.Join(dir, "dist")

	var out bytes.Buffer
	_, err := PackageSkill(&out, skillDir, outDir)
	if err == nil {
		t.Fatal("expected package validation failure")
	}
	mustContain(t, err.Error(), "metadata.invocation must be user or model")
	if _, statErr := os.Stat(filepath.Join(outDir, "demo-skill.skill")); !os.IsNotExist(statErr) {
		t.Fatalf("expected no archive, got err=%v", statErr)
	}
}

func createCatalogFixture(t *testing.T, catalog string) string {
	t.Helper()
	skillsDir := t.TempDir()
	writeInvocationSkill(t, skillsDir, "ask-matt", "user")
	writeInvocationSkill(t, skillsDir, "demo-model", "model")
	referencesDir := filepath.Join(skillsDir, "ask-matt", "references")
	if err := os.MkdirAll(referencesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(referencesDir, "CATALOG.md"), catalog)
	return skillsDir
}

func writeRouterCatalog(t *testing.T, skillsDir string, userNames, modelNames []string) {
	t.Helper()
	writeInvocationSkill(t, skillsDir, "ask-matt", "user")
	var catalog strings.Builder
	catalog.WriteString("# Skill Catalog\n\n## User-invoked\n\n| Skill | Purpose |\n| --- | --- |\n")
	catalog.WriteString("| `ask-matt` | Explain routes. |\n")
	for _, name := range userNames {
		catalog.WriteString("| `" + name + "` | User fixture. |\n")
	}
	catalog.WriteString("\n## Model-invoked\n\n| Skill | Purpose |\n| --- | --- |\n")
	for _, name := range modelNames {
		catalog.WriteString("| `" + name + "` | Model fixture. |\n")
	}
	referencesDir := filepath.Join(skillsDir, "ask-matt", "references")
	if err := os.MkdirAll(referencesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(referencesDir, "CATALOG.md"), catalog.String())
}

func writeInvocationSkill(t *testing.T, skillsDir, name, invocation string) {
	t.Helper()
	writeInvocationSkillWithDescription(
		t,
		skillsDir,
		name,
		invocation,
		"Perform a catalog fixture. Use when checking catalog synchronization. Avoid when another fixture is more specific.",
	)
}

func writeInvocationSkillWithDescription(t *testing.T, skillsDir, name, invocation, description string) {
	t.Helper()
	skillDir := filepath.Join(skillsDir, name)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	metadata := "  invocation: " + invocation + "\n"
	claudeControl := ""
	if invocation == "user" {
		claudeControl = "disable-model-invocation: true\n"
		metadata += "  opencode/autoinvoke: \"false\"\n"
		if err := os.MkdirAll(filepath.Join(skillDir, "agents"), 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(skillDir, "agents", "openai.yaml"), "policy:\n  allow_implicit_invocation: false\n")
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: `+name+`
description: `+description+`
`+claudeControl+`metadata:
`+metadata+`---

# Catalog Fixture
`)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func mustContain(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		t.Fatalf("expected %q in:\n%s", want, text)
	}
}

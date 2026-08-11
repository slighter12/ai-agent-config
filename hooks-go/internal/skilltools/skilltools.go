package skilltools

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"ai-agent-config/hooks/internal/securepath"
	"ai-agent-config/hooks/internal/skillcatalog"
)

const (
	maxDescriptionLength       = 1024
	maxRepoModelContribution   = 8000
	maxCompatibilityLength     = 500
	maxLicenseLength           = 200
	maxArgumentHintLength      = 200
	maxVersionLength           = 64
	maxValidatorFileBytes      = 1 << 20
	maxValidatorAggregateBytes = 8 << 20
	maxValidatorSkillEntries   = 1024
)

var (
	metadataKeyRE = regexp.MustCompile(`^[A-Za-z0-9_./-]+$`)
	versionRE     = regexp.MustCompile(`^v?\d+(?:\.\d+){0,2}(?:[-+][0-9A-Za-z.-]+)?$`)

	allowedProperties = map[string]struct{}{
		"name":                     {},
		"description":              {},
		"disable-model-invocation": {},
		"argument-hint":            {},
		"license":                  {},
		"compatibility":            {},
		"metadata":                 {},
	}
	requiredProperties = []string{"description", "name"}
)

const skillTemplate = `---
name: %[1]s
description: Perform [capability]. Use when the user asks for [specific trigger contexts, file types, workflows, or intents]. Avoid when [nearby non-goals or cases better handled by another skill].
metadata:
  invocation: "[user|model]"
---

# %[2]s

## Workflow

- [The smallest sequence of decisions or actions this skill must teach.]
- [A boundary that prevents a nearby misuse.]
`

type frontmatterValue struct {
	stringValue string
	listValue   []string
	mapping     map[string]string
	kind        valueKind
}

type valueKind int

const (
	stringKind valueKind = iota
	listKind
	mappingKind
)

type validatorBudget struct {
	remaining int64
}

func newValidatorBudget() *validatorBudget {
	return &validatorBudget{remaining: maxValidatorAggregateBytes}
}

func readValidatorFile(path string, budget *validatorBudget) ([]byte, error) {
	return readValidatorFileWithHook(path, budget, nil)
}

// readValidatorFileWithHook keeps replacement and growth checks deterministic
// in tests without making the production validator depend on timing.
func readValidatorFileWithHook(path string, budget *validatorBudget, afterOpen func()) ([]byte, error) {
	parent, name, err := securepath.OpenParent(path)
	if err != nil {
		return nil, err
	}
	defer parent.Close()
	return readValidatorFileFromParent(parent, name, path, budget, afterOpen)
}

func readValidatorFileRelative(root *os.File, relativePath string, budget *validatorBudget, afterOpen func()) ([]byte, error) {
	file, err := securepath.OpenFileRelative(root, relativePath)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "too many levels of symbolic links") {
			return nil, fmt.Errorf("validator input %s is a symlink", filepath.Join(root.Name(), relativePath))
		}
		return nil, fmt.Errorf("open validator input %s: %w", filepath.Join(root.Name(), relativePath), err)
	}
	defer file.Close()
	return readValidatorFileHandle(file, filepath.Join(root.Name(), relativePath), budget, afterOpen)
}

func readValidatorFileFromParent(parent *os.File, name, displayPath string, budget *validatorBudget, afterOpen func()) ([]byte, error) {
	file, err := securepath.OpenFileAt(parent, name)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "too many levels of symbolic links") {
			return nil, fmt.Errorf("validator input %s is a symlink", displayPath)
		}
		return nil, fmt.Errorf("open validator input %s: %w", displayPath, err)
	}
	defer file.Close()
	return readValidatorFileHandle(file, displayPath, budget, afterOpen)
}

func readValidatorFileHandle(file *os.File, path string, budget *validatorBudget, afterOpen func()) ([]byte, error) {
	if budget == nil {
		budget = newValidatorBudget()
	}
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat validator input %s: %w", path, err)
	}
	if !openedInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("validator input %s is not a regular file", path)
	}
	if openedInfo.Size() > maxValidatorFileBytes {
		return nil, fmt.Errorf("validator input %s exceeds 1 MiB limit while opening (%d bytes)", path, openedInfo.Size())
	}
	if openedInfo.Size() > budget.remaining {
		return nil, validatorAggregateLimitError(path, openedInfo.Size(), budget.remaining)
	}
	if afterOpen != nil {
		afterOpen()
	}
	readLimit := int64(maxValidatorFileBytes)
	if budget.remaining < readLimit {
		readLimit = budget.remaining
	}
	raw, err := io.ReadAll(io.LimitReader(file, readLimit+1))
	if err != nil {
		return nil, fmt.Errorf("read validator input %s: %w", path, err)
	}
	if int64(len(raw)) > maxValidatorFileBytes {
		return nil, fmt.Errorf("validator input %s exceeds 1 MiB limit while reading", path)
	}
	if int64(len(raw)) > budget.remaining {
		return nil, validatorAggregateLimitError(path, int64(len(raw)), budget.remaining)
	}

	finalInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat validator input %s after reading: %w", path, err)
	}
	if !finalInfo.Mode().IsRegular() || !os.SameFile(openedInfo, finalInfo) {
		return nil, fmt.Errorf("validator input %s was replaced while reading", path)
	}
	if finalInfo.Size() != openedInfo.Size() {
		return nil, fmt.Errorf("validator input %s grew or changed while reading", path)
	}
	if int64(len(raw)) != finalInfo.Size() {
		return nil, fmt.Errorf("validator input %s changed while reading", path)
	}
	budget.remaining -= int64(len(raw))
	return raw, nil
}

func validatorAggregateLimitError(path string, size, remaining int64) error {
	return fmt.Errorf("validator inputs exceed 8 MiB aggregate limit at %s (%d bytes; %d bytes remaining)", path, size, remaining)
}

func InitSkill(out io.Writer, skillName, path string) error {
	if err := skillcatalog.ValidateSkillName(skillName); err != nil {
		return err
	}
	skillDir, err := filepath.Abs(filepath.Join(path, skillName))
	if err != nil {
		return err
	}
	if _, err := os.Lstat(skillDir); err == nil {
		return fmt.Errorf("skill directory already exists: %s", skillDir)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(fmt.Sprintf(skillTemplate, skillName, titleCaseSkillName(skillName))), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(out, "Skill %q initialized at %s\n", skillName, skillDir)
	fmt.Fprintln(out, "\nNext steps:")
	fmt.Fprintln(out, "1. Replace bracketed placeholders and choose user or model invocation.")
	fmt.Fprintln(out, "2. For user invocation, add Claude and OpenCode controls plus agents/openai.yaml.")
	fmt.Fprintln(out, "3. Add references, scripts, or assets only when they materially improve the skill.")
	fmt.Fprintf(out, "4. Validate %s with the same agent-config entrypoint used for initialization.\n", skillDir)
	return nil
}

func ValidateSkill(skillPath string) (bool, string, error) {
	valid, message, _, err := validateSkillWithBudget(skillPath, newValidatorBudget())
	return valid, message, err
}

func validateSkillWithBudget(skillPath string, budget *validatorBudget) (bool, string, map[string]frontmatterValue, error) {
	directory, err := securepath.OpenDirectory(skillPath)
	if err != nil {
		return false, "", nil, err
	}
	defer directory.Close()
	return validateSkillDirectoryWithBudget(directory, skillPath, budget)
}

func validateSkillDirectoryWithBudget(directory *os.File, skillPath string, budget *validatorBudget) (bool, string, map[string]frontmatterValue, error) {
	info, err := directory.Stat()
	if err != nil {
		return false, "", nil, fmt.Errorf("stat skill path %s: %w", skillPath, err)
	}
	if !info.IsDir() {
		return false, "", nil, fmt.Errorf("skill path %s is not a directory", skillPath)
	}
	raw, err := readValidatorFileRelative(directory, "SKILL.md", budget, nil)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, "SKILL.md not found", nil, nil
		}
		return false, "", nil, err
	}
	content := string(raw)
	if !strings.HasPrefix(content, "---") {
		return false, "No YAML frontmatter found", nil, nil
	}
	frontmatterText, ok := extractFrontmatter(content)
	if !ok {
		return false, "Invalid frontmatter format", nil, nil
	}
	frontmatter, err := parseFrontmatter(frontmatterText)
	if err != nil {
		return false, fmt.Sprintf("Invalid frontmatter: %v", err), nil, nil
	}
	for key := range frontmatter {
		if _, ok := allowedProperties[key]; !ok {
			var unexpected []string
			for key := range frontmatter {
				if _, allowed := allowedProperties[key]; !allowed {
					unexpected = append(unexpected, key)
				}
			}
			sort.Strings(unexpected)
			return false, fmt.Sprintf("Unexpected key(s) in portable SKILL.md frontmatter: %s. Allowed properties are: %s", strings.Join(unexpected, ", "), strings.Join(sortedKeys(allowedProperties), ", ")), nil, nil
		}
	}
	for _, key := range requiredProperties {
		if _, ok := frontmatter[key]; !ok {
			return false, fmt.Sprintf("Missing %q in frontmatter", key), nil, nil
		}
	}
	nameValue := frontmatter["name"]
	if nameValue.kind != stringKind {
		return false, fmt.Sprintf("Name must be a string, got %s", kindName(nameValue.kind)), nil, nil
	}
	name := strings.TrimSpace(nameValue.stringValue)
	if err := skillcatalog.ValidateSkillName(name); err != nil {
		return false, fmt.Sprintf("Name %q %s", name, validationSuffix(err.Error())), nil, nil
	}
	descriptionValue := frontmatter["description"]
	if descriptionValue.kind != stringKind {
		return false, fmt.Sprintf("Description must be a string, got %s", kindName(descriptionValue.kind)), nil, nil
	}
	if err := validateDescription(strings.TrimSpace(descriptionValue.stringValue)); err != nil {
		return false, err.Error(), nil, nil
	}
	if value, ok := frontmatter["metadata"]; ok {
		if err := validateMetadata(value); err != nil {
			return false, err.Error(), nil, nil
		}
	}
	if err := validateDisableModelInvocation(frontmatter); err != nil {
		return false, err.Error(), nil, nil
	}
	if err := validateOptionalString(frontmatter, "license", maxLicenseLength); err != nil {
		return false, err.Error(), nil, nil
	}
	if err := validateOptionalString(frontmatter, "argument-hint", maxArgumentHintLength); err != nil {
		return false, err.Error(), nil, nil
	}
	if err := validateCompatibility(frontmatter); err != nil {
		return false, err.Error(), nil, nil
	}
	if failures := invocationFailuresFromDirectory(directory, skillPath, frontmatter, budget); len(failures) > 0 {
		return false, strings.Join(failures, "\n"), frontmatter, nil
	}
	return true, "Skill is valid", frontmatter, nil
}

func ValidateSkills(skillsPath string) (bool, string, error) {
	directory, err := securepath.OpenDirectory(skillsPath)
	if err != nil {
		return false, "", err
	}
	defer directory.Close()
	entries, err := readSkillEntriesFromDirectory(directory, skillsPath)
	if err != nil {
		return false, "", err
	}
	budget := newValidatorBudget()
	names := map[string]string{}
	invocations := map[string]string{}
	var failures []string
	discoveryCharacters := 0
	skillCount := 0
	modelSkillCount := 0
	userSkillCount := 0
	allSkillsValid := true
	routerPresent := false

	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 {
			failures = append(failures, fmt.Sprintf("%s: skill entry is a symlink", entry.Name()))
			allSkillsValid = false
			continue
		}
		if !entry.IsDir() {
			continue
		}
		skillPath := filepath.Join(skillsPath, entry.Name())
		skillDirectory, err := securepath.OpenDirectoryAt(directory, entry.Name())
		if err != nil {
			return false, "", fmt.Errorf("open skill %s: %w", entry.Name(), err)
		}
		if entry.Name() == "ask-matt" {
			routerPresent = true
		}
		skillCount++
		valid, message, frontmatter, err := validateSkillDirectoryWithBudget(skillDirectory, skillPath, budget)
		skillDirectory.Close()
		if err != nil {
			return false, "", err
		}
		if !valid && message == "SKILL.md not found" {
			skillCount--
			continue
		}
		if !valid {
			allSkillsValid = false
			failures = append(failures, fmt.Sprintf("%s: %s", entry.Name(), message))
			continue
		}
		name := strings.TrimSpace(frontmatter["name"].stringValue)
		if previous, exists := names[name]; exists {
			failures = append(failures, fmt.Sprintf("%s: duplicate skill name also used by %s", entry.Name(), previous))
		} else {
			names[name] = entry.Name()
		}
		if name != entry.Name() {
			failures = append(failures, fmt.Sprintf("%s: directory and skill name %q differ", entry.Name(), name))
		}
		invocation := frontmatter["metadata"].mapping["invocation"]
		invocations[name] = invocation
		switch invocation {
		case "model":
			modelSkillCount++
			discoveryCharacters += utf8.RuneCountInString(name)
			discoveryCharacters += utf8.RuneCountInString(strings.TrimSpace(frontmatter["description"].stringValue))
		case "user":
			userSkillCount++
		}
	}

	if !routerPresent {
		failures = append(failures, "required router skill ask-matt is missing")
	} else if _, hasRouter := invocations["ask-matt"]; hasRouter && allSkillsValid {
		failures = append(failures, catalogReferenceFailuresFromDirectory(directory, skillsPath, invocations, budget)...)
	}
	if discoveryCharacters > maxRepoModelContribution {
		failures = append(failures, fmt.Sprintf("repo model skill names and descriptions use %d contribution characters; maximum is %d", discoveryCharacters, maxRepoModelContribution))
	}
	if len(failures) > 0 {
		sort.Strings(failures)
		return false, strings.Join(failures, "\n"), nil
	}
	return true, fmt.Sprintf(
		"Skill catalog is valid: %d skills (%d model, %d user), %d/%d repo model name/description contribution characters",
		skillCount,
		modelSkillCount,
		userSkillCount,
		discoveryCharacters,
		maxRepoModelContribution,
	), nil
}

func readSkillEntries(skillsPath string) ([]os.DirEntry, error) {
	directory, err := securepath.OpenDirectory(skillsPath)
	if err != nil {
		return nil, fmt.Errorf("open skills path %s: %w", skillsPath, err)
	}
	defer directory.Close()
	return readSkillEntriesFromDirectory(directory, skillsPath)
}

func readSkillEntriesFromDirectory(directory *os.File, skillsPath string) ([]os.DirEntry, error) {
	openedInfo, err := directory.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat skills path %s: %w", skillsPath, err)
	}
	if !openedInfo.IsDir() {
		return nil, fmt.Errorf("skills path %s is not a directory", skillsPath)
	}
	entries, err := directory.ReadDir(maxValidatorSkillEntries + 1)
	if err != nil {
		return nil, fmt.Errorf("read skills path %s: %w", skillsPath, err)
	}
	if len(entries) > maxValidatorSkillEntries {
		return nil, fmt.Errorf("skills path %s contains more than %d entries", skillsPath, maxValidatorSkillEntries)
	}
	return entries, nil
}

func invocationFailuresFromDirectory(directory *os.File, skillPath string, frontmatter map[string]frontmatterValue, budget *validatorBudget) []string {
	var failures []string
	metadata, ok := frontmatter["metadata"]
	if !ok || metadata.kind != mappingKind {
		return []string{"metadata.invocation must be user or model"}
	}
	invocation := metadata.mapping["invocation"]
	if invocation != "user" && invocation != "model" {
		return []string{"metadata.invocation must be user or model"}
	}
	openCodeValue, hasOpenCode := metadata.mapping["opencode/autoinvoke"]
	claudeValue, hasClaude := frontmatter["disable-model-invocation"]
	claudeDisabled := hasClaude && strings.TrimSpace(claudeValue.stringValue) == "true"
	sidecarRaw, sidecarErr := readValidatorFileRelative(directory, filepath.Join("agents", "openai.yaml"), budget, nil)
	hasDisabledSidecar := sidecarErr == nil && codexImplicitInvocationDisabled(string(sidecarRaw))
	if sidecarErr != nil && !errors.Is(sidecarErr, os.ErrNotExist) {
		failures = append(failures, fmt.Sprintf("cannot read Codex sidecar: %v", sidecarErr))
	}
	if invocation == "user" {
		if !claudeDisabled {
			failures = append(failures, "user skill must set Claude disable-model-invocation to true")
		}
		if !hasOpenCode || openCodeValue != "false" {
			failures = append(failures, "user skill must set OpenCode metadata opencode/autoinvoke to false")
		}
		if !hasDisabledSidecar {
			failures = append(failures, "user skill must disable Codex implicit invocation")
		}
	} else {
		if claudeDisabled {
			failures = append(failures, "model skill must not disable Claude model invocation")
		}
		if hasOpenCode && openCodeValue == "false" {
			failures = append(failures, "model skill must not disable OpenCode autoinvocation")
		}
		if hasDisabledSidecar {
			failures = append(failures, "model skill must not disable Codex implicit invocation")
		}
	}
	return failures
}

func catalogReferenceFailuresFromDirectory(directory *os.File, skillsPath string, invocations map[string]string, budget *validatorBudget) []string {
	raw, err := readValidatorFileRelative(directory, filepath.Join("ask-matt", "references", "CATALOG.md"), budget, nil)
	if errors.Is(err, os.ErrNotExist) {
		return []string{"ask-matt: canonical catalog reference not found"}
	}
	if err != nil {
		return []string{fmt.Sprintf("ask-matt: cannot read canonical catalog reference: %v", err)}
	}

	entries := map[string]string{}
	var failures []string
	section := ""
	hasUserSection := false
	hasModelSection := false
	for _, line := range strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n") {
		trimmed := strings.TrimSpace(line)
		switch trimmed {
		case "## User-invoked":
			section = "user"
			hasUserSection = true
			continue
		case "## Model-invoked":
			section = "model"
			hasModelSection = true
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			section = ""
			continue
		}
		name, ok := catalogRowSkillName(trimmed)
		if !ok || section == "" {
			continue
		}
		if previous, exists := entries[name]; exists {
			failures = append(failures, fmt.Sprintf("ask-matt: catalog lists %s more than once (%s and %s)", name, previous, section))
			continue
		}
		entries[name] = section
	}
	if !hasUserSection || !hasModelSection {
		failures = append(failures, "ask-matt: catalog must contain User-invoked and Model-invoked sections")
	}

	for name, invocation := range invocations {
		catalogInvocation, exists := entries[name]
		if !exists {
			failures = append(failures, fmt.Sprintf("ask-matt: catalog is missing %s", name))
			continue
		}
		if catalogInvocation != invocation {
			failures = append(failures, fmt.Sprintf("ask-matt: catalog groups %s as %s; metadata invocation is %s", name, catalogInvocation, invocation))
		}
	}
	for name := range entries {
		if _, exists := invocations[name]; !exists {
			failures = append(failures, fmt.Sprintf("ask-matt: catalog contains unknown skill %s", name))
		}
	}
	return failures
}

func catalogRowSkillName(line string) (string, bool) {
	const prefix = "| `"
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	remainder := strings.TrimPrefix(line, prefix)
	end := strings.Index(remainder, "`")
	if end <= 0 || !strings.Contains(remainder[end+1:], "|") {
		return "", false
	}
	return remainder[:end], true
}

func codexImplicitInvocationDisabled(text string) bool {
	inPolicy := false
	for _, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !strings.HasPrefix(line, " ") {
			inPolicy = trimmed == "policy:"
			continue
		}
		if inPolicy && trimmed == "allow_implicit_invocation: false" {
			return true
		}
	}
	return false
}

func PackageSkill(out io.Writer, skillPath, outputDir string) (string, error) {
	absSkillPath, err := filepath.Abs(skillPath)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(absSkillPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("skill folder not found: %s", absSkillPath)
		}
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", absSkillPath)
	}
	if _, err := os.Stat(filepath.Join(absSkillPath, "SKILL.md")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("SKILL.md not found in %s", absSkillPath)
		}
		return "", err
	}
	fmt.Fprintln(out, "Validating skill...")
	valid, message, err := ValidateSkill(absSkillPath)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", fmt.Errorf("validation failed: %s", message)
	}
	fmt.Fprintf(out, "%s\n\n", message)

	if outputDir == "" {
		outputDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	outputFile := filepath.Join(outputDir, filepath.Base(absSkillPath)+".skill")
	if err := writeSkillArchive(out, absSkillPath, outputFile); err != nil {
		return "", err
	}
	fmt.Fprintf(out, "\nSuccessfully packaged skill to: %s\n", outputFile)
	return outputFile, nil
}

func writeSkillArchive(out io.Writer, skillPath, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	var files []string
	if err := filepath.WalkDir(skillPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Strings(files)
	for _, path := range files {
		rel, err := filepath.Rel(filepath.Dir(skillPath), path)
		if err != nil {
			return err
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(rel)
		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := writer.Write(raw); err != nil {
			return err
		}
		fmt.Fprintf(out, "  Added: %s\n", header.Name)
	}
	return nil
}

func extractFrontmatter(content string) (string, bool) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	if !strings.HasPrefix(content, "---\n") {
		return "", false
	}
	rest := content[len("---\n"):]
	index := strings.Index(rest, "\n---")
	if index < 0 {
		return "", false
	}
	after := rest[index+len("\n---"):]
	if after != "" && !strings.HasPrefix(after, "\n") {
		return "", false
	}
	return rest[:index], true
}

func parseFrontmatter(text string) (map[string]frontmatterValue, error) {
	result := map[string]frontmatterValue{}
	currentMapping := ""
	for index, line := range strings.Split(text, "\n") {
		lineNumber := index + 1
		stripped := strings.TrimSpace(line)
		if stripped == "" {
			continue
		}
		if strings.HasPrefix(line, " ") {
			if currentMapping != "metadata" {
				return nil, fmt.Errorf("Invalid frontmatter line %d: nested keys are only supported under metadata", lineNumber)
			}
			key, value, ok := strings.Cut(stripped, ":")
			if !ok {
				return nil, fmt.Errorf("Invalid frontmatter line %d: missing ':'", lineNumber)
			}
			key = strings.TrimSpace(key)
			if key == "" {
				return nil, fmt.Errorf("Invalid frontmatter line %d: empty metadata key", lineNumber)
			}
			metadata := result["metadata"].mapping
			if _, exists := metadata[key]; exists {
				return nil, fmt.Errorf("Duplicate metadata key: %s", key)
			}
			parsed := parseScalarOrInlineList(strings.TrimSpace(value))
			if parsed.kind != stringKind {
				return nil, fmt.Errorf("metadata.%s must be a string, got %s", key, kindName(parsed.kind))
			}
			metadata[key] = parsed.stringValue
			continue
		}
		key, value, ok := strings.Cut(stripped, ":")
		if !ok {
			return nil, fmt.Errorf("Invalid frontmatter line %d: missing ':'", lineNumber)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		currentMapping = ""
		if key == "" {
			return nil, fmt.Errorf("Invalid frontmatter line %d: empty key", lineNumber)
		}
		if _, exists := result[key]; exists {
			return nil, fmt.Errorf("Duplicate frontmatter key: %s", key)
		}
		if key == "metadata" {
			if value != "" {
				return nil, errors.New("metadata must be a nested mapping")
			}
			result[key] = frontmatterValue{kind: mappingKind, mapping: map[string]string{}}
			currentMapping = key
			continue
		}
		result[key] = parseScalarOrInlineList(value)
	}
	return result, nil
}

func parseScalarOrInlineList(value string) frontmatterValue {
	value = unquoteValue(value)
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		inner := strings.TrimSpace(value[1 : len(value)-1])
		if inner == "" {
			return frontmatterValue{kind: listKind, listValue: []string{}}
		}
		parts := strings.Split(inner, ",")
		items := make([]string, 0, len(parts))
		for _, part := range parts {
			items = append(items, unquoteValue(strings.TrimSpace(part)))
		}
		return frontmatterValue{kind: listKind, listValue: items}
	}
	return frontmatterValue{kind: stringKind, stringValue: value}
}

func unquoteValue(value string) string {
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}

func validateDescription(description string) error {
	if description == "" {
		return errors.New("Description cannot be empty")
	}
	lower := strings.ToLower(description)
	for _, fragment := range []string{"todo", "[", "]", "...", "automates the entire s"} {
		if strings.Contains(lower, fragment) {
			return fmt.Errorf("Description contains placeholder or truncated text: %s", fragment)
		}
	}
	if strings.Contains(description, "<") || strings.Contains(description, ">") {
		return errors.New("Description cannot contain angle brackets (< or >)")
	}
	if len(description) > maxDescriptionLength {
		return fmt.Errorf("Description is too long (%d characters). Maximum is %d characters.", len(description), maxDescriptionLength)
	}
	return nil
}

func validateDisableModelInvocation(frontmatter map[string]frontmatterValue) error {
	value, ok := frontmatter["disable-model-invocation"]
	if !ok {
		return nil
	}
	if value.kind != stringKind {
		return errors.New("disable-model-invocation must be true or false")
	}
	switch strings.TrimSpace(value.stringValue) {
	case "true", "false":
		return nil
	default:
		return errors.New("disable-model-invocation must be true or false")
	}
}

func validateMetadata(value frontmatterValue) error {
	if value.kind != mappingKind {
		return errors.New("metadata must be a YAML mapping")
	}
	for key, item := range value.mapping {
		if key == "" {
			return errors.New("metadata keys cannot be empty")
		}
		if !metadataKeyRE.MatchString(key) {
			return fmt.Errorf("metadata key %q should use letters, digits, dots, slashes, underscores, or hyphens", key)
		}
		item = strings.TrimSpace(item)
		if item == "" {
			return fmt.Errorf("metadata.%s cannot be empty", key)
		}
		if containsPlaceholder(item) {
			return fmt.Errorf("metadata.%s cannot contain placeholder or truncated text", key)
		}
	}
	version, ok := value.mapping["version"]
	if !ok {
		return nil
	}
	version = strings.TrimSpace(version)
	if version == "" {
		return errors.New("metadata.version cannot be empty")
	}
	if len(version) > maxVersionLength {
		return fmt.Errorf("metadata.version is too long (%d characters). Maximum is %d.", len(version), maxVersionLength)
	}
	if containsPlaceholder(version) {
		return errors.New("metadata.version cannot contain placeholder or truncated text")
	}
	if !versionRE.MatchString(version) {
		return errors.New("metadata.version should look like a semantic version, for example 0.1.0")
	}
	return nil
}

func validateOptionalString(frontmatter map[string]frontmatterValue, key string, maxLength int) error {
	value, ok := frontmatter[key]
	if !ok {
		return nil
	}
	if value.kind != stringKind {
		return fmt.Errorf("%s must be a string, got %s", key, kindName(value.kind))
	}
	text := strings.TrimSpace(value.stringValue)
	if text == "" {
		return fmt.Errorf("%s cannot be empty", key)
	}
	if len(text) > maxLength {
		return fmt.Errorf("%s is too long (%d characters). Maximum is %d.", key, len(text), maxLength)
	}
	if containsPlaceholder(text) {
		return fmt.Errorf("%s cannot contain placeholder or truncated text", key)
	}
	return nil
}

func validateCompatibility(frontmatter map[string]frontmatterValue) error {
	value, ok := frontmatter["compatibility"]
	if !ok {
		return nil
	}
	if value.kind == listKind {
		if len(value.listValue) == 0 {
			return errors.New("compatibility cannot be an empty list")
		}
		for _, item := range value.listValue {
			item = strings.TrimSpace(item)
			if item == "" {
				return errors.New("compatibility entries cannot be empty")
			}
			if len(item) > maxCompatibilityLength {
				return fmt.Errorf("compatibility entry is too long (%d characters). Maximum is %d.", len(item), maxCompatibilityLength)
			}
			if containsPlaceholder(item) {
				return errors.New("compatibility entries cannot contain placeholder or truncated text")
			}
		}
		return nil
	}
	return validateOptionalString(frontmatter, "compatibility", maxCompatibilityLength)
}

func containsPlaceholder(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "todo") || strings.Contains(value, "[") || strings.Contains(value, "]") || strings.Contains(value, "...")
}

func titleCaseSkillName(skillName string) string {
	parts := strings.Split(skillName, "-")
	for index, part := range parts {
		if part == "" {
			continue
		}
		parts[index] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func validationSuffix(message string) string {
	if strings.HasPrefix(message, "must") || strings.HasPrefix(message, "cannot") {
		return message
	}
	return ": " + message
}

func kindName(kind valueKind) string {
	switch kind {
	case stringKind:
		return "string"
	case listKind:
		return "list"
	case mappingKind:
		return "mapping"
	default:
		return "unknown"
	}
}

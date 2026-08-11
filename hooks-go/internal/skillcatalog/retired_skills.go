// Package skillcatalog owns repository-wide validation for retired skill names
// and the routing surfaces that must no longer refer to them.
package skillcatalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ai-agent-config/hooks/internal/pathidentity"
	"ai-agent-config/hooks/internal/securepath"
)

const (
	maxSkillNameLength              = 64
	maxRetiredSkillNames            = 128
	maxRetiredSkillsManifestBytes   = 64 << 10
	maxRoutingSurfaceBytes          = 1 << 20
	maxRoutingSurfaceFiles          = 1024
	maxRoutingSurfaceAggregateBytes = 8 << 20
	retiredSkillsFile               = "retired-skills.json"
)

// RetiredSkillsPath returns the canonical repository manifest path.
func RetiredSkillsPath(repoRoot string) string {
	return filepath.Join(repoRoot, "config", retiredSkillsFile)
}

// ParseRetiredSkills parses the canonical retired-skills JSON document.
//
// The document is intentionally strict: it has exactly one top-level field,
// the field is an array (an empty array is valid), and names are canonicalized
// by requiring valid syntax, uniqueness, and lexicographic order.
func ParseRetiredSkills(raw []byte) ([]string, error) {
	if len(raw) > maxRetiredSkillsManifestBytes {
		return nil, fmt.Errorf("parse retired skills: manifest exceeds 64 KiB limit (%d bytes)", len(raw))
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("parse retired skills: %w", err)
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '{' {
		return nil, errors.New("parse retired skills: expected top-level object")
	}

	fields := make(map[string]struct{})
	var retiredSkills *[]string
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("parse retired skills: %w", err)
		}
		field, ok := token.(string)
		if !ok {
			return nil, errors.New("parse retired skills: expected top-level field name")
		}
		if _, duplicate := fields[field]; duplicate {
			return nil, fmt.Errorf("parse retired skills: duplicate top-level field %q", field)
		}
		fields[field] = struct{}{}
		if field != "retired_skills" {
			return nil, fmt.Errorf("parse retired skills: unknown field %q", field)
		}
		if err := decoder.Decode(&retiredSkills); err != nil {
			return nil, fmt.Errorf("parse retired skills: %w", err)
		}
	}

	token, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("parse retired skills: %w", err)
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '}' {
		return nil, errors.New("parse retired skills: expected top-level object")
	}
	if len(fields) == 0 || retiredSkills == nil {
		return nil, errors.New("parse retired skills: missing retired_skills array")
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, errors.New("parse retired skills: trailing JSON is not allowed")
		}
		return nil, fmt.Errorf("parse retired skills: trailing JSON is not allowed: %w", err)
	}

	if err := ValidateRetiredSkillNames(*retiredSkills); err != nil {
		return nil, err
	}
	return *retiredSkills, nil
}

// LoadRetiredSkills reads and parses a canonical retired-skills JSON file.
// The path is used exactly as given: every component stays subject to the
// no-follow opener and to the parent-symlink rejection in
// inspectRetiredSkillsFile. Callers holding a repository root that may itself
// be reached through a symlink should use LoadRetiredSkillsFromRoot.
func LoadRetiredSkills(path string) ([]string, error) {
	raw, err := readRetiredSkillsFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("load retired skills manifest %s: %w", path, err)
	}
	return ParseRetiredSkills(raw)
}

// LoadRetiredSkillsFromRoot resolves the repository root -- the caller-supplied
// trust anchor -- and then reads <root>/config/retired-skills.json. Only the
// anchor is canonicalized: "config" and the manifest leaf stay subject to
// component-wise no-follow, so a symlinked "config" directory is rejected
// rather than redirecting the read outside the repository. This mirrors how
// ValidateRepositoryRetirements canonicalizes only its two directory anchors.
func LoadRetiredSkillsFromRoot(repoRoot string) ([]string, error) {
	canonicalRoot, _, err := pathidentity.ResolveDirectory(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve repository root %q: %w", repoRoot, err)
	}
	return LoadRetiredSkills(RetiredSkillsPath(canonicalRoot))
}

// readRetiredSkillsFile reads a manifest without following a symlink at the
// leaf. The hook exists only to make replacement and growth tests deterministic.
func readRetiredSkillsFile(path string, afterOpen func()) ([]byte, error) {
	preInfo, err := inspectRetiredSkillsFile(path)
	if err != nil {
		return nil, err
	}

	file, err := securepath.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("open retired skills manifest %s: %w", path, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat retired skills manifest %s: %w", path, err)
	}
	if !openedInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("retired skills manifest %s is not a regular file", path)
	}
	if !os.SameFile(preInfo, openedInfo) {
		return nil, fmt.Errorf("retired skills manifest %s was replaced while opening", path)
	}
	if openedInfo.Size() != preInfo.Size() {
		return nil, fmt.Errorf("retired skills manifest %s changed while opening", path)
	}
	if openedInfo.Size() > maxRetiredSkillsManifestBytes {
		return nil, fmt.Errorf("retired skills manifest %s exceeds 64 KiB limit while opening (%d bytes)", path, openedInfo.Size())
	}
	postOpenInfo, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("stat retired skills manifest %s after opening: %w", path, err)
	}
	if postOpenInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("retired skills manifest %s became a symlink while opening", path)
	}
	if !postOpenInfo.Mode().IsRegular() || !os.SameFile(preInfo, postOpenInfo) {
		return nil, fmt.Errorf("retired skills manifest %s was replaced while opening", path)
	}
	if postOpenInfo.Size() != preInfo.Size() {
		return nil, fmt.Errorf("retired skills manifest %s changed while opening", path)
	}

	if afterOpen != nil {
		afterOpen()
	}
	raw, err := io.ReadAll(io.LimitReader(file, maxRetiredSkillsManifestBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read retired skills manifest %s: %w", path, err)
	}
	if len(raw) > maxRetiredSkillsManifestBytes {
		return nil, fmt.Errorf("retired skills manifest %s exceeds 64 KiB limit while reading", path)
	}

	finalInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat retired skills manifest %s after reading: %w", path, err)
	}
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("stat retired skills manifest %s after reading: %w", path, err)
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("retired skills manifest %s became a symlink while reading", path)
	}
	if !finalInfo.Mode().IsRegular() || !pathInfo.Mode().IsRegular() ||
		!os.SameFile(preInfo, finalInfo) || !os.SameFile(preInfo, pathInfo) {
		return nil, fmt.Errorf("retired skills manifest %s was replaced while reading", path)
	}
	if finalInfo.Size() != preInfo.Size() || pathInfo.Size() != preInfo.Size() {
		return nil, fmt.Errorf("retired skills manifest %s grew or changed while reading", path)
	}
	if int64(len(raw)) != finalInfo.Size() {
		return nil, fmt.Errorf("retired skills manifest %s changed while reading", path)
	}
	return raw, nil
}

func inspectRetiredSkillsFile(path string) (os.FileInfo, error) {
	parent := filepath.Dir(path)
	parentInfo, err := os.Lstat(parent)
	if err != nil {
		return nil, err
	}
	if parentInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("retired skills manifest parent %s is a symlink", parent)
	}
	if !parentInfo.IsDir() {
		return nil, fmt.Errorf("retired skills manifest parent %s is not a directory", parent)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("retired skills manifest %s is a symlink", path)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("retired skills manifest %s is not a regular file", path)
	}
	if info.Size() > maxRetiredSkillsManifestBytes {
		return nil, fmt.Errorf("retired skills manifest %s exceeds 64 KiB limit (%d bytes)", path, info.Size())
	}
	return info, nil
}

// ValidateRetiredSkillNames validates names and their canonical ordering.
func ValidateRetiredSkillNames(names []string) error {
	if len(names) > maxRetiredSkillNames {
		return fmt.Errorf("retired_skills contains %d names, exceeds %d-name limit", len(names), maxRetiredSkillNames)
	}
	previous := ""
	for index, name := range names {
		if err := ValidateSkillName(name); err != nil {
			return fmt.Errorf("retired_skills[%d] %q: %w", index, name, err)
		}
		if index > 0 && name <= previous {
			if name == previous {
				return fmt.Errorf("retired_skills[%d] %q: duplicate name", index, name)
			}
			return fmt.Errorf("retired_skills must be sorted: %q appears after %q", name, previous)
		}
		previous = name
	}
	return nil
}

// ValidateSkillName validates the canonical skill name syntax shared by the
// skill validator, initializer, and retirement manifest.
func ValidateSkillName(name string) error {
	switch {
	case name == "":
		return errors.New("must use lowercase letters, digits, and hyphens only")
	case strings.IndexFunc(name, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-')
	}) >= 0:
		return errors.New("must use lowercase letters, digits, and hyphens only")
	case strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") || strings.Contains(name, "--"):
		return errors.New("cannot start/end with hyphen or contain consecutive hyphens")
	case len(name) > maxSkillNameLength:
		return fmt.Errorf("must be %d characters or fewer", maxSkillNameLength)
	default:
		return nil
	}
}

// ValidateRepositoryRetirements verifies that retired names have disappeared
// from active skills and routing surfaces.
func ValidateRepositoryRetirements(repoRoot, skillsPath string, names []string) error {
	if err := ValidateRetiredSkillNames(names); err != nil {
		return err
	}

	canonicalRepoRoot, _, err := pathidentity.ResolveDirectory(repoRoot)
	if err != nil {
		return fmt.Errorf("resolve repository root: %w", err)
	}
	canonicalSkillsPath, _, err := pathidentity.ResolveDirectory(skillsPath)
	if err != nil {
		return fmt.Errorf("resolve skills path: %w", err)
	}
	var failures []string
	for _, name := range names {
		if active, err := activeSkillExists(filepath.Join(canonicalSkillsPath, name)); err != nil {
			return fmt.Errorf("inspect retired skill %q: %w", name, err)
		} else if active {
			failures = append(failures, fmt.Sprintf("retired skill %q is still active under %s", name, filepath.Join(canonicalSkillsPath, name)))
		}
	}

	paths, err := activeRoutingPaths(canonicalRepoRoot, canonicalSkillsPath, names)
	if err != nil {
		return err
	}
	budget := routingScanBudget{remaining: maxRoutingSurfaceAggregateBytes}
	for _, path := range paths {
		if err := scanRoutingFileWithBudget(path, names, &failures, &budget); err != nil {
			return err
		}
	}
	if len(failures) == 0 {
		return nil
	}
	sort.Strings(failures)
	return errors.New(strings.Join(failures, "\n"))
}

func activeSkillExists(path string) (bool, error) {
	directoryInfo, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if directoryInfo.Mode()&os.ModeSymlink != 0 {
		return false, fmt.Errorf("active skill path %s is a symlink", path)
	}
	if !directoryInfo.IsDir() {
		return false, fmt.Errorf("active skill path %s is not a directory", path)
	}

	info, err := os.Lstat(filepath.Join(path, "SKILL.md"))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return false, fmt.Errorf("active skill file %s is a symlink", filepath.Join(path, "SKILL.md"))
	}
	if !info.Mode().IsRegular() {
		return false, fmt.Errorf("active skill file %s is not a regular file", filepath.Join(path, "SKILL.md"))
	}
	return true, nil
}

func activeRoutingPaths(repoRoot, skillsPath string, retiredNames []string) ([]string, error) {
	paths := []string{}
	for _, relative := range []string{"README.md", "AGENTS.md", "CLAUDE.md", filepath.Join("docs", "skill-routing-eval.md"), filepath.Join("skills", "UPSTREAM.md")} {
		if err := appendRegularFilesWithin(repoRoot, filepath.Join(repoRoot, relative), nil, &paths); err != nil {
			return nil, err
		}
	}
	for _, relative := range []string{"config", filepath.Join("docs", "agents")} {
		root := filepath.Join(repoRoot, relative)
		if err := appendRegularFilesWithin(repoRoot, root, func(path string) bool {
			return filepath.Clean(path) != filepath.Join(repoRoot, "config", retiredSkillsFile)
		}, &paths); err != nil {
			return nil, err
		}
	}
	if err := appendActiveSkillFiles(skillsPath, retiredNames, &paths); err != nil {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

func appendActiveSkillFiles(skillsPath string, retiredNames []string, paths *[]string) error {
	entries, err := readRoutingDirectoryEntries(skillsPath)
	if err != nil {
		return fmt.Errorf("read skills path %s: %w", skillsPath, err)
	}
	retired := make(map[string]struct{}, len(retiredNames))
	for _, name := range retiredNames {
		retired[name] = struct{}{}
	}
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("skills path entry %s is a symlink", filepath.Join(skillsPath, entry.Name()))
		}
		entryInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect skills path entry %s: %w", filepath.Join(skillsPath, entry.Name()), err)
		}
		if !entry.IsDir() {
			if !entryInfo.Mode().IsRegular() {
				return fmt.Errorf("skills path entry %s is not a regular file or directory", filepath.Join(skillsPath, entry.Name()))
			}
			continue
		}
		// The active-skill scan intentionally excludes only directories whose
		// names are in the manifest; callers have already validated that a
		// retired directory with SKILL.md is an active-skill failure.
		if _, ok := retired[entry.Name()]; ok {
			continue
		}
		root := filepath.Join(skillsPath, entry.Name())
		active, err := activeSkillExists(root)
		if err != nil {
			return fmt.Errorf("inspect active skill %s: %w", root, err)
		}
		if !active {
			continue
		}
		if err := appendRegularFilesWithin(root, root, nil, paths); err != nil {
			return err
		}
	}
	return nil
}

// inspectRoutingPathComponents checks the lexical path under a trusted root
// without resolving any component through a symlink. A missing final component
// is reported as absent so optional routing surfaces retain their old behavior.
func inspectRoutingPathComponents(root, path string) (bool, error) {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return false, err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return false, fmt.Errorf("path escapes routing root %s", root)
	}
	rootInfo, err := os.Lstat(root)
	if err != nil {
		return false, err
	}
	if rootInfo.Mode()&os.ModeSymlink != 0 {
		return false, fmt.Errorf("routing root %s is a symlink", root)
	}
	if !rootInfo.IsDir() {
		return false, fmt.Errorf("routing root %s is not a directory", root)
	}
	if relative == "." {
		return true, nil
	}

	current := root
	components := strings.Split(relative, string(filepath.Separator))
	for index, component := range components {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return false, fmt.Errorf("routing path component %s is a symlink", current)
		}
		if index < len(components)-1 && !info.IsDir() {
			return false, fmt.Errorf("routing path component %s is not a directory", current)
		}
	}
	return true, nil
}

func readRoutingDirectoryEntries(path string) ([]os.DirEntry, error) {
	preInfo, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if preInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("routing directory %s is a symlink", path)
	}
	if !preInfo.IsDir() {
		return nil, fmt.Errorf("routing directory %s is not a directory", path)
	}
	directory, err := securepath.OpenDirectory(path)
	if err != nil {
		return nil, fmt.Errorf("open routing directory %s: %w", path, err)
	}
	defer directory.Close()
	openedInfo, err := directory.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat routing directory %s: %w", path, err)
	}
	if !openedInfo.IsDir() || !os.SameFile(preInfo, openedInfo) {
		return nil, fmt.Errorf("routing directory %s was replaced while opening", path)
	}
	entries, err := directory.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("read routing directory %s: %w", path, err)
	}
	postInfo, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("stat routing directory %s after reading: %w", path, err)
	}
	if postInfo.Mode()&os.ModeSymlink != 0 || !postInfo.IsDir() || !os.SameFile(preInfo, postInfo) {
		return nil, fmt.Errorf("routing directory %s was replaced while reading", path)
	}
	return entries, nil
}

func appendRegularFilesWithin(containmentRoot, root string, include func(string) bool, paths *[]string) error {
	exists, err := inspectRoutingPathComponents(containmentRoot, root)
	if err != nil {
		return fmt.Errorf("inspect routing surface %s: %w", root, err)
	}
	if !exists {
		return nil
	}
	info, err := os.Lstat(root)
	if err != nil {
		return fmt.Errorf("inspect routing surface %s: %w", root, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("routing surface %s is a symlink", root)
	}
	if !info.IsDir() {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("routing surface %s is not a regular file", root)
		}
		if include == nil || include(root) {
			if err := appendRoutingPath(root, paths); err != nil {
				return err
			}
		}
		return nil
	}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk routing surface %s: %w", path, walkErr)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("routing surface %s is a symlink", path)
		}
		entryInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect routing surface %s: %w", path, err)
		}
		if entry.IsDir() {
			if !entryInfo.IsDir() {
				return fmt.Errorf("routing surface %s is not a directory", path)
			}
			return nil
		}
		if !entryInfo.Mode().IsRegular() {
			return fmt.Errorf("routing surface %s is not a regular file", path)
		}
		if include != nil && !include(path) {
			return nil
		}
		exists, err := inspectRoutingPathComponents(containmentRoot, path)
		if err != nil {
			return fmt.Errorf("inspect routing surface %s: %w", path, err)
		}
		if !exists {
			return fmt.Errorf("routing surface %s disappeared while discovering files", path)
		}
		return appendRoutingPath(path, paths)
	})
	if err != nil {
		return err
	}
	return nil
}

func appendRoutingPath(path string, paths *[]string) error {
	if len(*paths) >= maxRoutingSurfaceFiles {
		return fmt.Errorf("routing surface file count exceeds %d limit", maxRoutingSurfaceFiles)
	}
	*paths = append(*paths, path)
	return nil
}

type routingScanBudget struct {
	remaining int64
}

func scanRoutingFileWithBudget(path string, names []string, failures *[]string, budget *routingScanBudget) error {
	return scanRoutingFileWithHook(path, names, failures, budget, nil)
}

// scanRoutingFileWithHook keeps replacement and growth checks deterministic in
// tests without making production routing validation depend on timing.
func scanRoutingFileWithHook(path string, names []string, failures *[]string, budget *routingScanBudget, afterOpen func()) error {
	if budget == nil {
		budget = &routingScanBudget{remaining: maxRoutingSurfaceAggregateBytes}
	}
	info, err := inspectRoutingFile(path)
	if err != nil {
		return fmt.Errorf("read routing surface %s: %w", path, err)
	}
	if info.Size() > budget.remaining {
		return routingAggregateLimitError(path, info.Size(), budget.remaining)
	}
	file, err := securepath.OpenFile(path)
	if err != nil {
		return fmt.Errorf("read routing surface %s: %w", path, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat routing surface %s: %w", path, err)
	}
	if !openedInfo.Mode().IsRegular() {
		return fmt.Errorf("routing surface %s is not a regular file", path)
	}
	if !os.SameFile(info, openedInfo) {
		return fmt.Errorf("routing surface %s was replaced while opening", path)
	}
	if openedInfo.Size() != info.Size() {
		return fmt.Errorf("routing surface %s changed while opening", path)
	}
	if openedInfo.Size() > maxRoutingSurfaceBytes {
		return fmt.Errorf("routing surface %s exceeds 1 MiB limit while opening (%d bytes)", path, openedInfo.Size())
	}
	if openedInfo.Size() > budget.remaining {
		return routingAggregateLimitError(path, openedInfo.Size(), budget.remaining)
	}
	postOpenInfo, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("stat routing surface %s after opening: %w", path, err)
	}
	if postOpenInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("routing surface %s became a symlink while opening", path)
	}
	if !postOpenInfo.Mode().IsRegular() || !os.SameFile(info, postOpenInfo) {
		return fmt.Errorf("routing surface %s was replaced while opening", path)
	}
	if postOpenInfo.Size() != info.Size() {
		return fmt.Errorf("routing surface %s changed while opening", path)
	}
	if afterOpen != nil {
		afterOpen()
	}
	readLimit := int64(maxRoutingSurfaceBytes)
	if budget.remaining < readLimit {
		readLimit = budget.remaining
	}
	raw, err := io.ReadAll(io.LimitReader(file, readLimit+1))
	if err != nil {
		return fmt.Errorf("read routing surface %s: %w", path, err)
	}
	if int64(len(raw)) > maxRoutingSurfaceBytes {
		return fmt.Errorf("routing surface %s exceeds 1 MiB limit while reading", path)
	}
	if int64(len(raw)) > budget.remaining {
		return routingAggregateLimitError(path, int64(len(raw)), budget.remaining)
	}
	finalInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat routing surface %s after reading: %w", path, err)
	}
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("stat routing surface %s after reading: %w", path, err)
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("routing surface %s became a symlink while reading", path)
	}
	if !finalInfo.Mode().IsRegular() || !pathInfo.Mode().IsRegular() ||
		!os.SameFile(info, finalInfo) || !os.SameFile(info, pathInfo) {
		return fmt.Errorf("routing surface %s was replaced while reading", path)
	}
	if finalInfo.Size() != info.Size() || pathInfo.Size() != info.Size() {
		return fmt.Errorf("routing surface %s grew or changed while reading", path)
	}
	if int64(len(raw)) != finalInfo.Size() {
		return fmt.Errorf("routing surface %s changed while reading", path)
	}
	budget.remaining -= int64(len(raw))
	for _, name := range names {
		if containsSkillReference(string(raw), name) {
			*failures = append(*failures, fmt.Sprintf("retired skill %q is referenced by %s", name, path))
		}
	}
	return nil
}

func inspectRoutingFile(path string) (os.FileInfo, error) {
	parent := filepath.Dir(path)
	parentInfo, err := os.Lstat(parent)
	if err != nil {
		return nil, err
	}
	if parentInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("routing surface parent %s is a symlink", parent)
	}
	if !parentInfo.IsDir() {
		return nil, fmt.Errorf("routing surface parent %s is not a directory", parent)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("routing surface %s is a symlink", path)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("routing surface %s is not a regular file", path)
	}
	if info.Size() > maxRoutingSurfaceBytes {
		return nil, fmt.Errorf("routing surface %s exceeds 1 MiB limit (%d bytes)", path, info.Size())
	}
	return info, nil
}

func routingAggregateLimitError(path string, size, remaining int64) error {
	return fmt.Errorf("routing surfaces exceed 8 MiB aggregate limit at %s (%d bytes; %d bytes remaining)", path, size, remaining)
}

func containsSkillReference(text, name string) bool {
	for offset := 0; ; {
		index := strings.Index(text[offset:], name)
		if index < 0 {
			return false
		}
		index += offset
		beforeOK := index == 0 || !isSkillNameCharacter(text[index-1])
		after := index + len(name)
		afterOK := after == len(text) || !isSkillNameCharacter(text[after])
		if beforeOK && afterOK {
			return true
		}
		offset = index + len(name)
		if offset >= len(text) {
			return false
		}
	}
}

func isSkillNameCharacter(value byte) bool {
	return value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z' || value >= '0' && value <= '9' || value == '-'
}

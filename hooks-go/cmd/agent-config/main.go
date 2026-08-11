package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ai-agent-config/hooks/internal/agentconfig"
	"ai-agent-config/hooks/internal/pathidentity"
	"ai-agent-config/hooks/internal/skillcatalog"
	"ai-agent-config/hooks/internal/skilltools"
)

func main() {
	repoRoot := flag.String("repo-root", "", "repository root; defaults to current directory or parent of hooks-go")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}
	config, err := agentconfig.NewConfig(*repoRoot, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "agent-config: %v\n", err)
		os.Exit(1)
	}
	if err := run(config, flag.Arg(0), flag.Args()[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "agent-config: %v\n", err)
		os.Exit(1)
	}
}

func run(config agentconfig.Config, command string, args []string) error {
	switch command {
	case "install":
		return config.Install()
	case "setup-codex-config":
		return config.SetupCodexConfig()
	case "setup-codex-agents":
		return config.SetupCodexAgents()
	case "setup-claude-agents":
		return config.SetupClaudeAgents()
	case "setup-codex-shell":
		return config.SetupCodexShell()
	case "setup-codex-hooks":
		return config.SetupCodexHooks()
	case "build-hooks":
		return config.BuildHooks()
	case "init-skill":
		return runInitSkill(args)
	case "validate-skill":
		return runValidateSkill(args)
	case "validate-skills":
		return runValidateSkills(args, config.RepoRoot)
	case "package-skill":
		return runPackageSkill(args)
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func runValidateSkills(args []string, repoRoot string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: agent-config validate-skills <skills-directory>")
	}
	skillsDir, err := repositorySkillsDirectory(repoRoot, args[0])
	if err != nil {
		return err
	}
	retiredNames, err := skillcatalog.LoadRetiredSkillsFromRoot(repoRoot)
	if err != nil {
		return err
	}
	valid, message, err := skilltools.ValidateSkills(skillsDir)
	if err != nil {
		return err
	}
	if !valid {
		fmt.Println(message)
		return fmt.Errorf("skill catalog validation failed")
	}
	if err := skillcatalog.ValidateRepositoryRetirements(repoRoot, skillsDir, retiredNames); err != nil {
		return fmt.Errorf("retired skill validation failed: %w", err)
	}
	fmt.Println(message)
	return nil
}

func repositorySkillsDirectory(repoRoot, skillsDir string) (string, error) {
	canonicalRepoRoot, _, err := pathidentity.ResolveDirectory(repoRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repository root %q: %w", repoRoot, err)
	}
	expectedSkillsPath := filepath.Join(canonicalRepoRoot, "skills")
	canonicalExpectedSkillsDir, expectedInfo, err := pathidentity.ResolveDirectory(expectedSkillsPath)
	if err != nil {
		return "", fmt.Errorf("resolve repository skills directory %q: %w; pass <repo-root>/skills", expectedSkillsPath, err)
	}
	canonicalSkillsDir, providedInfo, err := pathidentity.ResolveDirectory(skillsDir)
	if err != nil {
		return "", fmt.Errorf("resolve skills directory %q: %w", skillsDir, err)
	}
	if canonicalExpectedSkillsDir != canonicalSkillsDir || !os.SameFile(expectedInfo, providedInfo) {
		return "", fmt.Errorf("skills directory %q does not match this repository's skills directory %q; pass <repo-root>/skills", skillsDir, expectedSkillsPath)
	}
	return canonicalSkillsDir, nil
}

func runInitSkill(args []string) error {
	var skillName string
	var path string
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--path" {
			if index+1 >= len(args) {
				return fmt.Errorf("usage: agent-config init-skill <skill-name> --path <path>")
			}
			path = args[index+1]
			index++
			continue
		}
		if strings.HasPrefix(arg, "--path=") {
			path = strings.TrimPrefix(arg, "--path=")
			continue
		}
		if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown init-skill option %q", arg)
		}
		if skillName != "" {
			return fmt.Errorf("usage: agent-config init-skill <skill-name> --path <path>")
		}
		skillName = arg
	}
	if skillName == "" || path == "" {
		return fmt.Errorf("usage: agent-config init-skill <skill-name> --path <path>")
	}
	return skilltools.InitSkill(os.Stdout, skillName, path)
}

func runValidateSkill(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: agent-config validate-skill <skill-directory>")
	}
	valid, message, err := skilltools.ValidateSkill(args[0])
	if err != nil {
		return err
	}
	fmt.Println(message)
	if !valid {
		return fmt.Errorf("skill validation failed")
	}
	return nil
}

func runPackageSkill(args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("usage: agent-config package-skill <skill-directory> [output-directory]")
	}
	outputDir := ""
	if len(args) == 2 {
		outputDir = args[1]
	}
	_, err := skilltools.PackageSkill(os.Stdout, args[0], outputDir)
	return err
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: agent-config [--repo-root PATH] <command>\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Commands:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  install\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-codex-config\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-codex-agents\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-claude-agents\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-codex-shell\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-codex-hooks\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  build-hooks\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  init-skill <skill-name> --path <path>\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  validate-skill <skill-directory>\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  validate-skills <skills-directory>\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  package-skill <skill-directory> [output-directory]\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	flag.PrintDefaults()
}

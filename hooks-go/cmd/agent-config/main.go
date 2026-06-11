package main

import (
	"flag"
	"fmt"
	"os"

	"ai-agent-config/hooks/internal/agentconfig"
)

func main() {
	repoRoot := flag.String("repo-root", "", "repository root; defaults to current directory or parent of hooks-go")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	config, err := agentconfig.NewConfig(*repoRoot, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "agent-config: %v\n", err)
		os.Exit(1)
	}
	if err := run(config, flag.Arg(0)); err != nil {
		fmt.Fprintf(os.Stderr, "agent-config: %v\n", err)
		os.Exit(1)
	}
}

func run(config agentconfig.Config, command string) error {
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
	case "build-hooks":
		return config.BuildHooks()
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: agent-config [--repo-root PATH] <command>\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Commands:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  install\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-codex-config\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-codex-agents\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-claude-agents\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  setup-codex-shell\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  build-hooks\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	flag.PrintDefaults()
}

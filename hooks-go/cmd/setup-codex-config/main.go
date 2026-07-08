package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"ai-agent-config/hooks/internal/codexconfig"
)

func main() {
	codexHome := flag.String("codex-home", "", "Codex home directory; defaults to CODEX_HOME or ~/.codex")
	flag.Parse()

	result, err := codexconfig.Apply(*codexHome, time.Now())
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup-codex-config: %v\n", err)
		os.Exit(1)
	}
	if result.Changed {
		if result.BackupPath != "" {
			fmt.Printf("   + Updated %s (backup: %s)\n", result.Path, result.BackupPath)
		} else {
			fmt.Printf("   + Created %s\n", result.Path)
		}
	} else {
		fmt.Printf("   = %s already has parent-safe defaults\n", result.Path)
	}
	if result.ProfileChanged {
		if result.ProfileBackupPath != "" {
			fmt.Printf("   + Updated %s (backup: %s)\n", result.ProfilePath, result.ProfileBackupPath)
		} else {
			fmt.Printf("   + Created %s\n", result.ProfilePath)
		}
		return
	}
	fmt.Printf("   = %s already has workspace-git profile\n", result.ProfilePath)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Meta scaffolding CLI for AI agent contexts",
	Long: `scaffold manages AI agent contexts in your project.
It reads agent definitions and skill files, combines them, and writes
the result to configurable target files (e.g. CLAUDE.md, .cursorrules).`,
	Version: version,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(resetCmd)
}

// projectRoot finds the directory containing scaffold.config.json by walking up
// from the current working directory. Returns the cwd if none found.
func projectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "scaffold.config.json")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Fall back to cwd
	cwd, _ := os.Getwd()
	return cwd
}

// fatalf prints an error message and exits.
func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

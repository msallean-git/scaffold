package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/kenlo/scaffold/internal/config"
)

var initForce bool
var initTarget string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize scaffold in the current directory",
	Long:  "Creates scaffold.config.json, agents.md, AgentSkills/, and a starter coding.md skill.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := os.Getwd()
		return runInit(dir, initForce, initTarget)
	},
}

func init() {
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing files")
	initCmd.Flags().StringVar(&initTarget, "target", "", "Default target file path (e.g. .cursorrules)")
}

func runInit(dir string, force bool, target string) error {
	cfg := config.Default()
	if target != "" {
		cfg.Targets = []config.Target{{Path: target, Mode: "overwrite"}}
	}

	// scaffold.config.json
	cfgPath := filepath.Join(dir, config.Filename)
	if err := writeFile(cfgPath, force, func() error {
		return config.Save(dir, cfg)
	}); err != nil {
		return err
	}
	fmt.Printf("created %s\n", config.Filename)

	// agents.md
	agentsPath := filepath.Join(dir, cfg.AgentsFile)
	if err := writeFile(agentsPath, force, func() error {
		return os.WriteFile(agentsPath, []byte(sampleAgentsMd), 0644)
	}); err != nil {
		return err
	}
	fmt.Printf("created %s\n", cfg.AgentsFile)

	// AgentSkills/
	skillsDir := filepath.Join(dir, cfg.SkillsDir)
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return err
	}
	fmt.Printf("created %s/\n", cfg.SkillsDir)

	// AgentSkills/coding.md
	codingPath := filepath.Join(skillsDir, "coding.md")
	if err := writeFile(codingPath, force, func() error {
		return os.WriteFile(codingPath, []byte(sampleCodingMd), 0644)
	}); err != nil {
		return err
	}
	fmt.Printf("created %s/coding.md\n", cfg.SkillsDir)

	fmt.Println("\nScaffold initialized. Edit agents.md to define your agents.")
	return nil
}

// writeFile writes using fn, skipping if path exists and force is false.
func writeFile(path string, force bool, fn func() error) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("skipping %s (already exists, use --force to overwrite)\n", filepath.Base(path))
			return nil
		}
	}
	return fn()
}

const sampleAgentsMd = `# Agents

## general-dev

**Skills:** coding

You are a senior software developer. Write clean, idiomatic code with clear
variable names and minimal comments. Prefer simple solutions over clever ones.
Always consider edge cases and error handling.
`

const sampleCodingMd = `# Coding

Write clean, readable code that follows the conventions of the language and
project. Prefer clarity over cleverness. Keep functions small and focused.
Handle errors explicitly. Avoid global state.
`

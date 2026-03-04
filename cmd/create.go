package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kenlo/scaffold/internal/agents"
	"github.com/kenlo/scaffold/internal/config"
	"github.com/kenlo/scaffold/internal/lock"
	"github.com/kenlo/scaffold/internal/skills"
)

var createSkillsFlag string
var createInstructions string
var createStubSkills bool

var createCmd = &cobra.Command{
	Use:   "create <agent-name>",
	Short: "Add a new agent to agents.md",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		root := projectRoot()

		cfg, err := config.Load(root)
		if err != nil {
			return fmt.Errorf("no scaffold.config.json found — run 'scaffold init' first")
		}

		// Parse skill list from flag.
		var skillList []string
		if createSkillsFlag != "" {
			for _, s := range strings.Split(createSkillsFlag, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					skillList = append(skillList, s)
				}
			}
		}

		// Check for duplicate agent name.
		agentsPath := filepath.Join(root, cfg.AgentsFile)
		all, err := agents.ParseFile(agentsPath)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", cfg.AgentsFile, err)
		}
		if agents.FindByName(all, agentName) != nil {
			return fmt.Errorf("agent %q already exists in %s", agentName, cfg.AgentsFile)
		}

		lk, err := lock.Acquire(root)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		defer lk.Release()

		// Optionally create stub skill files.
		skillsDir := filepath.Join(root, cfg.SkillsDir)
		if createStubSkills {
			for _, skill := range skillList {
				if !skills.Exists(skillsDir, skill) {
					if err := skills.CreateStub(skillsDir, skill); err != nil {
						return fmt.Errorf("creating stub for skill %q: %w", skill, err)
					}
					fmt.Printf("created %s/%s.md (stub)\n", cfg.SkillsDir, skill)
				}
			}
		}

		agent := agents.Agent{
			Name:         agentName,
			Skills:       skillList,
			SystemPrompt: createInstructions,
		}

		if err := agents.AppendAgent(agentsPath, agent); err != nil {
			return fmt.Errorf("writing %s: %w", cfg.AgentsFile, err)
		}

		fmt.Printf("added agent %q to %s\n", agentName, cfg.AgentsFile)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&createSkillsFlag, "skills", "", "Comma-separated list of skills (e.g. coding,testing)")
	createCmd.Flags().StringVar(&createInstructions, "instructions", "", "System prompt / instructions for the agent")
	createCmd.Flags().BoolVar(&createStubSkills, "create-skills", false, "Create stub skill files for missing skills")
}

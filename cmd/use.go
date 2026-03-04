package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/kenlo/scaffold/internal/agents"
	"github.com/kenlo/scaffold/internal/config"
	"github.com/kenlo/scaffold/internal/output"
	"github.com/kenlo/scaffold/internal/skills"
)

var useDryRun bool
var useVerbose bool

var useCmd = &cobra.Command{
	Use:   "use <agent-name>",
	Short: "Activate an agent and write its context to target files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		root := projectRoot()

		cfg, err := config.Load(root)
		if err != nil {
			return fmt.Errorf("no scaffold.config.json found — run 'scaffold init' first")
		}

		all, err := agents.ParseFile(filepath.Join(root, cfg.AgentsFile))
		if err != nil {
			return fmt.Errorf("parsing %s: %w", cfg.AgentsFile, err)
		}

		agent := agents.FindByName(all, agentName)
		if agent == nil {
			return fmt.Errorf("agent %q not found in %s", agentName, cfg.AgentsFile)
		}

		skillContents, errs := skills.LoadAll(filepath.Join(root, cfg.SkillsDir), agent.Skills)
		for _, e := range errs {
			fmt.Printf("warning: %v\n", e)
		}

		rendered := output.Render(*agent, skillContents)

		if useDryRun {
			fmt.Println("--- dry run output ---")
			fmt.Print(rendered)
			fmt.Println("--- end dry run ---")
			return nil
		}

		for _, target := range cfg.Targets {
			targetPath := filepath.Join(root, target.Path)
			switch target.Mode {
			case "section":
				err = output.WriteSection(targetPath, rendered)
			default: // "overwrite"
				err = output.WriteOverwrite(targetPath, rendered)
			}
			if err != nil {
				return fmt.Errorf("writing %s: %w", target.Path, err)
			}
			if useVerbose {
				fmt.Printf("wrote %s (mode=%s)\n", target.Path, target.Mode)
			} else {
				fmt.Printf("wrote %s\n", target.Path)
			}
		}

		cfg.ActiveAgent = agentName
		if err := config.Save(root, cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("active agent: %s\n", agentName)
		return nil
	},
}

func init() {
	useCmd.Flags().BoolVar(&useDryRun, "dry-run", false, "Preview output without writing files")
	useCmd.Flags().BoolVar(&useVerbose, "verbose", false, "Show extra detail")
}

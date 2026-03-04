package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/kenlo/scaffold/internal/config"
	"github.com/kenlo/scaffold/internal/lock"
	"github.com/kenlo/scaffold/internal/output"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Clear active agent context from target files",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := projectRoot()

		cfg, err := config.Load(root)
		if err != nil {
			return fmt.Errorf("no scaffold.config.json found — run 'scaffold init' first")
		}

		lk, err := lock.Acquire(root)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		defer lk.Release()

		for _, target := range cfg.Targets {
			targetPath := filepath.Join(root, target.Path)
			switch target.Mode {
			case "section":
				if err := output.ClearSection(targetPath); err != nil {
					return fmt.Errorf("clearing %s: %w", target.Path, err)
				}
			default: // "overwrite"
				if err := output.WriteOverwrite(targetPath, ""); err != nil {
					return fmt.Errorf("clearing %s: %w", target.Path, err)
				}
			}
			fmt.Printf("cleared %s\n", target.Path)
		}

		cfg.ActiveAgent = ""
		if err := config.Save(root, cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Println("active agent cleared")
		return nil
	},
}

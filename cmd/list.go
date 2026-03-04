package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/kenlo/scaffold/internal/agents"
	"github.com/kenlo/scaffold/internal/config"
)

var listJSON bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all defined agents",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := projectRoot()
		cfg, err := config.Load(root)
		if err != nil {
			return fmt.Errorf("no scaffold.config.json found — run 'scaffold init' first")
		}

		all, err := agents.ParseFile(filepath.Join(root, cfg.AgentsFile))
		if err != nil {
			return fmt.Errorf("parsing %s: %w", cfg.AgentsFile, err)
		}

		if listJSON {
			return printJSON(all, cfg.ActiveAgent)
		}
		return printTable(all, cfg.ActiveAgent)
	},
}

func init() {
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
}

func printTable(all []agents.Agent, active string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSKILLS\tACTIVE")
	fmt.Fprintln(w, "----\t------\t------")
	for _, a := range all {
		activeFlag := ""
		if a.Name == active {
			activeFlag = "*"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", a.Name, strings.Join(a.Skills, ", "), activeFlag)
	}
	return w.Flush()
}

type jsonAgent struct {
	Name   string   `json:"name"`
	Skills []string `json:"skills"`
	Active bool     `json:"active"`
}

func printJSON(all []agents.Agent, active string) error {
	out := make([]jsonAgent, len(all))
	for i, a := range all {
		out[i] = jsonAgent{
			Name:   a.Name,
			Skills: a.Skills,
			Active: a.Name == active,
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

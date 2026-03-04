package agents

import (
	"bufio"
	"os"
	"strings"
)

// ParseFile reads the agents.md file at path and returns all agent definitions.
func ParseFile(path string) ([]Agent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var agents []Agent
	var current *Agent
	var promptLines []string

	flush := func() {
		if current == nil {
			return
		}
		current.SystemPrompt = strings.TrimSpace(strings.Join(promptLines, "\n"))
		agents = append(agents, *current)
		current = nil
		promptLines = nil
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Agent header: "## <name>" (but not "# Agents")
		if strings.HasPrefix(line, "## ") {
			flush()
			name := strings.TrimSpace(strings.TrimPrefix(line, "## "))
			current = &Agent{Name: name}
			continue
		}

		if current == nil {
			continue
		}

		// Skills line: "**Skills:** skill1, skill2"
		if strings.HasPrefix(line, "**Skills:**") {
			raw := strings.TrimPrefix(line, "**Skills:**")
			raw = strings.TrimSpace(raw)
			for _, s := range strings.Split(raw, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					current.Skills = append(current.Skills, s)
				}
			}
			continue
		}

		promptLines = append(promptLines, line)
	}
	flush()

	return agents, scanner.Err()
}

// FindByName returns the agent with the given name, or nil if not found.
func FindByName(agents []Agent, name string) *Agent {
	for i := range agents {
		if agents[i].Name == name {
			return &agents[i]
		}
	}
	return nil
}

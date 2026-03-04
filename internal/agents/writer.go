package agents

import (
	"fmt"
	"os"
	"strings"
)

// AppendAgent appends a new agent block to the agents.md file at path.
func AppendAgent(path string, agent Agent) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	skillsLine := strings.Join(agent.Skills, ", ")
	prompt := agent.SystemPrompt
	if prompt == "" {
		prompt = fmt.Sprintf("You are %s. Add your instructions here.", agent.Name)
	}

	block := fmt.Sprintf("\n## %s\n\n**Skills:** %s\n\n%s\n", agent.Name, skillsLine, prompt)
	_, err = f.WriteString(block)
	return err
}

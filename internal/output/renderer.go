package output

import (
	"fmt"
	"strings"

	"github.com/kenlo/scaffold/internal/agents"
)

const (
	startMarker = "<!-- scaffold:start"
	endMarker   = "<!-- scaffold:end -->"
)

// Render combines an agent and its skill contents into the scaffold output block.
func Render(agent agents.Agent, skillContents map[string]string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<!-- scaffold:start agent=%s -->\n", agent.Name))
	sb.WriteString(fmt.Sprintf("# Agent: %s\n\n", agent.Name))
	sb.WriteString("## Instructions\n\n")
	sb.WriteString(agent.SystemPrompt)
	sb.WriteString("\n\n## Skills\n")

	for _, skill := range agent.Skills {
		content, ok := skillContents[skill]
		if !ok {
			continue
		}
		sb.WriteString(fmt.Sprintf("\n### %s\n\n", skill))
		sb.WriteString(strings.TrimSpace(content))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(endMarker)
	sb.WriteString("\n")

	return sb.String()
}

// StartMarker returns the full start marker string for the given agent name.
func StartMarker(agentName string) string {
	return fmt.Sprintf("%s agent=%s -->", startMarker, agentName)
}

// IsStartMarker reports whether line begins a scaffold block.
func IsStartMarker(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), startMarker)
}

// IsEndMarker reports whether line is the scaffold end marker.
func IsEndMarker(line string) bool {
	return strings.TrimSpace(line) == endMarker
}

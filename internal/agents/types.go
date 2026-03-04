package agents

// Agent represents a parsed agent definition from agents.md.
type Agent struct {
	Name         string
	Skills       []string
	SystemPrompt string
}

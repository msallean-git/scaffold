package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const Filename = "scaffold.config.json"

// Target describes a file that scaffold writes to.
type Target struct {
	Path string `json:"path"`
	Mode string `json:"mode"` // "overwrite" | "section"
}

// Config is the root configuration structure.
type Config struct {
	Version     int      `json:"version"`
	Targets     []Target `json:"targets"`
	AgentsFile  string   `json:"agentsFile"`
	SkillsDir   string   `json:"skillsDir"`
	ActiveAgent string   `json:"activeAgent"`
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		Version:    1,
		Targets:    []Target{{Path: "CLAUDE.md", Mode: "overwrite"}},
		AgentsFile: "agents.md",
		SkillsDir:  "AgentSkills",
	}
}

// Load reads scaffold.config.json from dir.
func Load(dir string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(dir, Filename))
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Save writes cfg to scaffold.config.json in dir.
func Save(dir string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, Filename), data, 0644)
}

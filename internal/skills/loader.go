package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

// Load reads a skill file from skillsDir/<name>.md and returns its content.
func Load(skillsDir, name string) (string, error) {
	path := filepath.Join(skillsDir, name+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("skill %q not found at %s", name, path)
	}
	return string(data), nil
}

// LoadAll loads all named skills from skillsDir, returning a map of name→content.
// Missing skill files are noted as errors but loading continues.
func LoadAll(skillsDir string, names []string) (map[string]string, []error) {
	result := make(map[string]string, len(names))
	var errs []error
	for _, name := range names {
		content, err := Load(skillsDir, name)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		result[name] = content
	}
	return result, errs
}

// Exists reports whether a skill file exists.
func Exists(skillsDir, name string) bool {
	_, err := os.Stat(filepath.Join(skillsDir, name+".md"))
	return err == nil
}

// CreateStub writes an empty stub skill file for name.
func CreateStub(skillsDir, name string) error {
	path := filepath.Join(skillsDir, name+".md")
	content := fmt.Sprintf("# %s\n\nDescribe the %s skill here.\n", name, name)
	return os.WriteFile(path, []byte(content), 0644)
}

package output

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

// WriteOverwrite replaces the entire content of path with rendered.
func WriteOverwrite(path, rendered string) error {
	return os.WriteFile(path, []byte(rendered), 0644)
}

// WriteSection replaces only the scaffold block (between markers) in path.
// If no markers are found, the block is appended to the file.
func WriteSection(path, rendered string) error {
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var out bytes.Buffer
	inBlock := false
	replaced := false

	scanner := bufio.NewScanner(bytes.NewReader(existing))
	for scanner.Scan() {
		line := scanner.Text()

		if !inBlock && IsStartMarker(line) {
			inBlock = true
			replaced = true
			out.WriteString(rendered)
			continue
		}
		if inBlock {
			if IsEndMarker(line) {
				inBlock = false
			}
			continue
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}

	if !replaced {
		// No existing block — append.
		if out.Len() > 0 && !strings.HasSuffix(out.String(), "\n\n") {
			out.WriteByte('\n')
		}
		out.WriteString(rendered)
	}

	return os.WriteFile(path, out.Bytes(), 0644)
}

// ClearSection removes the scaffold block from path (section mode).
func ClearSection(path string) error {
	existing, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var out bytes.Buffer
	inBlock := false

	scanner := bufio.NewScanner(bytes.NewReader(existing))
	for scanner.Scan() {
		line := scanner.Text()
		if !inBlock && IsStartMarker(line) {
			inBlock = true
			continue
		}
		if inBlock {
			if IsEndMarker(line) {
				inBlock = false
			}
			continue
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}

	return os.WriteFile(path, out.Bytes(), 0644)
}

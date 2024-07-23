package rules

import (
	"bufio"
	"os"
	"strings"
)

func ParseIgnoreFile(filePath string) (map[string][]string, error) {
	patterns := make(map[string][]string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) > 1 {
				// Check if the key already exists and append to its slice
				patterns[parts[0]] = append(patterns[parts[0]], parts[1])
			} else if len(parts) == 1 {
				// Add an empty pattern if only the path is present
				patterns[parts[0]] = append(patterns[parts[0]], "")
			}
		}
	}

	// TODO: handle "not a valid ignore file" case - don't exit program, do log a warning
	// TODO Q: What happens if we add something to patterns that's not a valid path?

	return patterns, scanner.Err()
}

package lint

import (
	"bufio"
	"fmt"
	"helm.sh/helm/v3/pkg/lint/support"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Ignorer struct {
	Patterns        map[string][]string
	debugFnOverride func(string, ...interface{})
}

const DefaultIgnoreFileName = ".helmlintignore"

func NewIgnorer(chartPath, ignoreFilePath string, debugLogFn func(string, ...interface{})) *Ignorer {
	ignorer := &Ignorer{ debugFnOverride: debugLogFn}

	if ignoreFilePath == "" {
		ignoreFilePath = filepath.Join(chartPath, DefaultIgnoreFileName)
		ignorer.Debug("\nNo HelmLintIgnore file specified, will try and use the following: %s\n", ignoreFilePath)
	}

	ignorer.Debug("\nUsing ignore file: %s\n", ignoreFilePath)
	ignorer.loadPatternsFromFilePath(ignoreFilePath)
	return ignorer
}

func (i *Ignorer) FilterErrors(errors []error) []error {
	keepers := make([]error, 0)
	for _, err := range errors {
		if !i.match(err.Error()) {
			keepers = append(keepers, err)
		}
	}

	return keepers
}

func (i *Ignorer) FilterMessages(messages []support.Message) []support.Message {
	keepers := make([]support.Message, 0)
	for _, msg := range messages {
		if !i.match(msg.Err.Error()) {
			keepers = append(keepers, msg)
		}
	}
	return keepers
}

func (i *Ignorer) match(errText string) bool {
	errorFullPath := extractFullPathFromError(errText)
	if len(errorFullPath) == 0 {
		i.Debug("Unable to find a path for error, guess we'll keep it: %s", errText)
		return false
	}

	i.Debug("Extracted full path: %s\n", errorFullPath)
	for ignorablePath, pathPatterns := range i.Patterns {
		cleanIgnorablePath := filepath.Clean(ignorablePath)
		if strings.Contains(errorFullPath, cleanIgnorablePath) {
			for _, pattern := range pathPatterns {
				if strings.Contains(errText, pattern) {
					i.Debug("Ignoring error: [%s] %s\n\n", errorFullPath, errText)
					return true
				}
			}
		}
	}

	i.Debug("keeping unignored error: [%s]", errText)
	return false
}

// TODO: figure out & fix or remove
func extractFullPathFromError(errorString string) string {
	parts := strings.Split(errorString, ":")
	if len(parts) > 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func (i *Ignorer) loadPatternsFromFilePath(filePath string) {
	if i.Patterns == nil {
		i.Patterns = make(map[string][]string)
	}

	file, err := os.Open(filePath)
	if err != nil {
		i.Debug("failed to open lint ignore file: %s", filePath)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close ignore file at [%s]: %v", filePath, err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) > 1 {
				// Check if the key already exists and append to its slice
				i.Patterns[parts[0]] = append(i.Patterns[parts[0]], parts[1])
			} else if len(parts) == 1 {
				// Add an empty pattern if only the path is present
				i.Patterns[parts[0]] = append(i.Patterns[parts[0]], "")
			}
		}
	}
	return
}

func (i *Ignorer) Debug(format string, args ...interface{}) {
	if i.debugFnOverride == nil {
		i.debugFnOverride = func(format string, v ...interface{}) {
			format = fmt.Sprintf("[debug] %s\n", format)
			log.Output(2, fmt.Sprintf(format, v...))
		}
	}

	i.debugFnOverride(format, args...)
}

/* TODO HIP-0019
- find ignore file path for a subchart
- add a chart or two for the end to end tests via testdata like in pkg/lint/lint_test.go

Later/never
- XDG support
- helm config file support
- ignore file validation
*/

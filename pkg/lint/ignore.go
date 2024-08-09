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

// Ignorer provides a means of suppressing unwanted helm lint errors and messages
// by comparing them to an ignore list provided in a plaintext helm lint ignore file.
type Ignorer struct {
	Patterns        map[string][]string
	ErrorPatterns   map[string][]string
	debugFnOverride func(string, ...interface{})
}

// DefaultIgnoreFileName is the name of the lint ignore file
// an Ignorer will seek out at load/parse time.
const DefaultIgnoreFileName = ".helmlintignore"

// NewIgnorer builds an Ignorer object that enables helm to discard specific lint result Messages
// and Errors should they match the ignore rules in the specified .helmlintignore file.
func NewIgnorer(chartPath, ignoreFilePath string, debugLogFn func(string, ...interface{})) *Ignorer {
	ignorer := &Ignorer{
		debugFnOverride: debugLogFn,
		Patterns:        make(map[string][]string),
		ErrorPatterns:   make(map[string][]string),
	}

	if ignoreFilePath == "" {
		ignoreFilePath = filepath.Join(chartPath, DefaultIgnoreFileName)
		ignorer.Debug("\nNo HelmLintIgnore file specified, will try and use the following: %s\n", ignoreFilePath)
	}

	ignorer.Debug("\nUsing ignore file: %s\n", ignoreFilePath)
	ignorer.loadPatternsFromFilePath(ignoreFilePath)
	return ignorer
}

// FilterErrors takes a slice of errors and returns a new slice containing only the
// errors that do not match this Ignorer's ignore string patterns.
func (i *Ignorer) FilterErrors(errors []error) []error {
	keepers := make([]error, 0)
	for _, err := range errors {
		if !i.match(err.Error()) {
			keepers = append(keepers, err)
		}
	}

	return keepers
}

// FilterMessages takes a slice of support.Message and returns a new slice
// containing only the Messages that do not match this Ignorer's ignore string patterns.
func (i *Ignorer) FilterMessages(messages []support.Message) []support.Message {
	keepers := make([]support.Message, 0)
	for _, msg := range messages {
		if !i.match(msg.Err.Error()) {
			keepers = append(keepers, msg)
		}
	}
	return keepers
}

// FilterNoPathErrors takes a slice of linter result Messages and a slice of errors and returns
// only those Messages and errors that do not match this Ignorer's ignore string patterns.
func (i *Ignorer) FilterNoPathErrors(messages []support.Message, errors []error) ([]support.Message, []error) {
	KeepersErr := make([]error, 0)
	KeepersMsg := make([]support.Message, 0)
	for _, err := range errors {
		if i.MatchNoPathError(err.Error()) {
			KeepersErr = append(KeepersErr, err)
			for _, msg := range messages {
				KeepersMsg = append(KeepersMsg, msg)
			}
		}
	}
	return KeepersMsg, KeepersErr
}

// MatchNoPathError checks a given string to determine whether it looks like a
// helm lint finding that does not specifically specify an offending file path.
// These will usually be related to Chart.yaml contents rather than a template
// inside the chart itself.
func (i *Ignorer) MatchNoPathError(errText string) bool {
	for ignorableError := range i.ErrorPatterns {
		parts := strings.SplitN(ignorableError, ":", 2)
		prefix := strings.TrimSpace(parts[0])
		if match, _ := filepath.Match(ignorableError, errText); match {
			i.Debug("Ignoring partial match error: [%s] %s\n\n", ignorableError, errText)
			return false
		}
		if matched, _ := filepath.Match(prefix, errText); matched {
			i.Debug("Ignoring error: [%s] %s\n\n", ignorableError, errText)
			return false
		}
	}
	i.Debug("keeping unignored error: [%s]", errText)
	return true
}

// Debug provides an Ignorer with a caller-overridable logging function
// intended to match the behavior of the top level debug() method from package main.
//
// When no i.debugFnOverride is present Debug will fall back to a naive
// implementation that assumes all debug output should be logged and not swallowed.
func (i *Ignorer) Debug(format string, args ...interface{}) {
	if i.debugFnOverride == nil {
		i.debugFnOverride = func(format string, v ...interface{}) {
			format = fmt.Sprintf("[debug] %s\n", format)
			log.Output(2, fmt.Sprintf(format, v...))
		}
	}

	i.debugFnOverride(format, args...)
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
	file, err := os.Open(filePath)
	if err != nil {
		i.Debug("failed to open lint ignore file: %s", filePath)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "error_lint_ignore=") {
			parts := strings.SplitN(line[18:], "error_lint_ignore=", 2) // Skipping 'error_lint_ignore=' prefix
			if len(parts) == 2 {
				i.ErrorPatterns[parts[0]] = append(i.ErrorPatterns[parts[0]], parts[1])
			} else {
				i.ErrorPatterns[parts[0]] = append(i.ErrorPatterns[parts[0]], "")
			}
		} else {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) > 1 {
				i.Patterns[parts[0]] = append(i.Patterns[parts[0]], parts[1])
			} else {
				i.Patterns[parts[0]] = append(i.Patterns[parts[0]], "")
			}
		}
	}
}

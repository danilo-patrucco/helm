package lint

import (
	"bufio"
	"fmt"
	"helm.sh/helm/v3/pkg/lint/support"
	"io"
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
	ignorer.loadFromFilePath(ignoreFilePath)
	return ignorer
}

func (i *Ignorer) FilterLintResult(messages []support.Message, errors []error) ([]support.Message, []error) {
	support.DumpInputsAsJsonLines(messages, errors)

	messagesFiltered := i.FilterMessages(messages)
	errorsFiltered := i.FilterErrors(errors)
	return i.FilterNoPathErrors(messagesFiltered, errorsFiltered)
}

// FilterErrors takes a slice of errors and returns a new slice containing only the
// errors that do not match this Ignorer's ignore string patterns.
func (i *Ignorer) FilterErrors(errors []error) []error {
	keepers := make([]error, 0)
	for _, err := range errors {
		errText := err.Error()
		if !i.isIgnorable(errText) {
			i.debugFnOverride("not filtering this error because it is not ignorable: %s", errText)
			keepers = append(keepers, err)
		}
		i.debugFnOverride("filtering this error because it is ignorable: %s", errText)
	}

	return keepers
}

// FilterMessages takes a slice of support.Message and returns a new slice
// containing only the Messages that do not match this Ignorer's ignore string patterns.
func (i *Ignorer) FilterMessages(messages []support.Message) []support.Message {
	keepers := make([]support.Message, 0)
	for _, msg := range messages {
		if !i.isIgnorable(msg.Err.Error()) {
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
		if i.IsIgnoredPathlessError(err.Error()) {
			continue
		}
		KeepersErr = append(KeepersErr, err)
		for _, msg := range messages {
			KeepersMsg = append(KeepersMsg, msg)
		}
	}
	return KeepersMsg, KeepersErr
}

// IsIgnoredPathlessError checks a given string to determine whether it looks like a
// helm lint finding that does not specifically specify an offending file path.
// These will usually be related to Chart.yaml contents rather than a template
// inside the chart itself.
func (i *Ignorer) IsIgnoredPathlessError(errText string) bool {
	for ignorableError := range i.ErrorPatterns {
		parts := strings.SplitN(ignorableError, ":", 2)
		prefix := strings.TrimSpace(parts[0])
		if match, _ := filepath.Match(ignorableError, errText); match {
			i.Debug("Ignoring partial match error: [%s] %s\n\n", ignorableError, errText)
			return true
		}
		if matched, _ := filepath.Match(prefix, errText); matched {
			i.Debug("Ignoring error: [%s] %s\n\n", ignorableError, errText)
			return true
		}
	}
	i.Debug("keeping unignored error: [%s]", errText)
	return false
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

func (i *Ignorer) isIgnorable(errText string) bool {
	errorFullPath, err := extractFullPathFromError(errText)
	if err != nil {
		i.Debug("Unable to find a path for error, guess we'll keep it: %s, %v", errText, err)
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
			i.Debug("keeping unmatched error: %s", errText)
		}
	}

	i.Debug("keeping unignored error: [%s]", errText)
	return false
}

// TODO: figure out & fix or remove
func extractFullPathFromError(errText string) (string, error) {
	delimiter := ":"
	// splits into N parts delimited by colons
	parts := strings.Split(errText, delimiter)
	// if 3 or more parts, return the second part, after trimming its spaces
	if len(parts) > 2 {
		return strings.TrimSpace(parts[1]), nil
	}
	// if fewer than 3 parts, return empty string
	return "", fmt.Errorf("fewer than three [%s]-delimited parts found, no path here: %s", delimiter, errText)
}

func (i *Ignorer) loadFromFilePath(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		i.Debug("failed to open lint ignore file: %s", filePath)
		return
	}
	defer file.Close()
	i.loadFromReader(file)
}

func (i *Ignorer) loadFromReader(rdr io.Reader) {
	const chartLevelErrorPrefix = "error_lint_ignore="
	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		isChartLevelError := strings.HasPrefix(line, chartLevelErrorPrefix)

		if isChartLevelError {
			// handle chart-level errors
			tokens := strings.SplitN(line[len(chartLevelErrorPrefix):], chartLevelErrorPrefix, 2) // Skipping 'error_lint_ignore=' prefix
			var leftThing string
			var rightThing = ""
			if len(tokens) == 2 {
				leftThing, rightThing = tokens[0], tokens[1]
				i.ErrorPatterns[leftThing] = append(i.ErrorPatterns[leftThing], rightThing)
			} else {
				leftThing = tokens[0]
				i.ErrorPatterns[leftThing] = append(i.ErrorPatterns[leftThing], rightThing)
			}
		} else {
			// handle chart yaml file errors in specific template files
			parts := strings.SplitN(line, " ", 2)
			if len(parts) > 1 {
				i.Patterns[parts[0]] = append(i.Patterns[parts[0]], parts[1])
			} else {
				i.Patterns[parts[0]] = append(i.Patterns[parts[0]], "")
			}
		}
	}
}
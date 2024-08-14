package ignore

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// RuleLoader provides a means of suppressing unwanted helm lint errors and messages
// by comparing them to an ignore list provided in a plaintext helm lint ignore file.
type RuleLoader struct {
	Matchers        []MatchesErrors
	debugFnOverride func(string, ...interface{})
}

func (i *RuleLoader) LogAttrs() slog.Attr {
	return slog.Group("RuleLoader",
		slog.String("Matchers", fmt.Sprintf("%v", i.Matchers)),
	)
}

// DefaultIgnoreFileName is the name of the lint ignore file
// an RuleLoader will seek out at load/parse time.
const DefaultIgnoreFileName = ".helmlintignore"

const NoMessageText = ""

// NewRuleLoader builds an RuleLoader object that enables helm to discard specific lint result Messages
// and Errors should they match the ignore rules in the specified .helmlintignore file.
func NewRuleLoader(chartPath, ignoreFilePath string, debugLogFn func(string, ...interface{})) (*RuleLoader, error) {
	out := &RuleLoader{
		debugFnOverride: debugLogFn,
	}

	if ignoreFilePath == "" {
		ignoreFilePath = filepath.Join(chartPath, DefaultIgnoreFileName)
		out.Debug("\nNo HelmLintIgnore file specified, will try and use the following: %s\n", ignoreFilePath)
	}

	// attempt to load ignore patterns from ignoreFilePath.
	// if none are found, return an empty ignorer so the program can keep running.
	out.Debug("\nUsing ignore file: %s\n", ignoreFilePath)
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		out.Debug("failed to open lint ignore file: %s", ignoreFilePath)
		return out, nil
	}
	defer file.Close()

	out.LoadFromReader(file)
	out.Debug("RuleLoader loaded.", out.LogAttrs())
	return out, nil
}

// Debug provides an RuleLoader with a caller-overridable logging function
// intended to match the behavior of the top level debug() method from package main.
//
// When no i.debugFnOverride is present Debug will fall back to a naive
// implementation that assumes all debug output should be logged and not swallowed.
func (i *RuleLoader) Debug(format string, args ...interface{}) {
	if i.debugFnOverride == nil {
		i.debugFnOverride = func(format string, v ...interface{}) {
			format = fmt.Sprintf("[debug] %s\n", format)
			log.Output(2, fmt.Sprintf(format, v...))
		}
	}

	i.debugFnOverride(format, args...)
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

func (i *RuleLoader) LoadFromReader(rdr io.Reader) {
	const pathlessPatternPrefix = "error_lint_ignore="
	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		isPathlessPattern := strings.HasPrefix(line, pathlessPatternPrefix)

		if isPathlessPattern {
			i.storePathlessPattern(line, pathlessPatternPrefix)
		} else {
			i.storePathfulPattern(line)
		}
	}
}

func (i *RuleLoader) storePathlessPattern(line string, pathlessPatternPrefix string) {
	// handle chart-level errors
	// Drop 'error_lint_ignore=' prefix from rule before saving it
	const numSplits = 2
	tokens := strings.SplitN(line[len(pathlessPatternPrefix):], pathlessPatternPrefix, numSplits)
	if len(tokens) == numSplits {
		// TODO: find an example for this one - not sure we still use it
		messageText, _ := tokens[0], tokens[1]
		i.Matchers = append(i.Matchers, PathlessRule{RuleText: line, MessageText: messageText})
	} else {
		messageText := tokens[0]
		i.Matchers = append(i.Matchers, PathlessRule{RuleText: line, MessageText: messageText})
	}
}

func (i *RuleLoader) storePathfulPattern(line string) {
	const separator = " "
	const numSplits = 2

	// handle chart yaml file errors in specific template files
	parts := strings.SplitN(line, separator, numSplits)
	if len(parts) == numSplits {
		messagePath, messageText := parts[0], parts[1]
		i.Matchers = append(i.Matchers, Rule{RuleText: line, MessagePath: messagePath, MessageText: messageText})
	} else {
		messagePath := parts[0]
		i.Matchers = append(i.Matchers, Rule{RuleText: line, MessagePath: messagePath, MessageText: NoMessageText})
	}
}

func (i *RuleLoader) loadFromFilePath(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		i.Debug("failed to open lint ignore file: %s", filePath)
		return
	}
	defer file.Close()
	i.LoadFromReader(file)
}

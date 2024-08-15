package ignore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DefaultIgnoreFileName is the name of the lint ignore file
const DefaultIgnoreFileName = ".helmlintignore"

const NoMessageText = ""

func LoadFromFilePath(chartPath, ignoreFilePath string, debugLogFn func(string, ...interface{})) ([]MatchesErrors, error) {
	if ignoreFilePath == "" {
		ignoreFilePath = filepath.Join(chartPath, DefaultIgnoreFileName)
		debug("\nNo HelmLintIgnore file specified, will try and use the following: %s\n", ignoreFilePath)
	}

	// attempt to load ignore patterns from ignoreFilePath.
	// if none are found, return an empty ignorer so the program can keep running.
	debug("\nUsing ignore file: %s\n", ignoreFilePath)
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		debug("failed to open lint ignore file: %s", ignoreFilePath)
		return []MatchesErrors{}, nil
	}
	defer file.Close()

	matchers := LoadFromReader(file)
	return matchers, nil
}

// Debug provides an RuleLoader with a caller-overridable logging function
// intended to match the behavior of the top level debug() method from package main.
//
// When no i.debugFnOverride is present Debug will fall back to a naive
// implementation that assumes all debug output should be logged and not swallowed.
func debug(format string, args ...interface{}) {
	if debugFn == nil {
		defaultDebugFn(format, args...)
	} else {
		debugFn(format, args...)
	}
	return
}

// TODO: figure out & fix or remove
func pathToOffendingFile(errText string) (string, error) {
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

func LoadFromReader(rdr io.Reader) []MatchesErrors {
	const pathlessPatternPrefix = "error_lint_ignore="
	matchers := make([]MatchesErrors, 0)

	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, pathlessPatternPrefix) {
			matchers = append(matchers, buildPathlessPattern(line, pathlessPatternPrefix))
		} else {
			matchers = append(matchers, buildPathfulPattern(line))
		}
	}

	return matchers
}

func buildPathlessPattern(line string, pathlessPatternPrefix string) PathlessRule {
	// handle chart-level errors
	// Drop 'error_lint_ignore=' prefix from rule before saving it
	const numSplits = 2
	tokens := strings.SplitN(line[len(pathlessPatternPrefix):], pathlessPatternPrefix, numSplits)
	if len(tokens) == numSplits {
		// TODO: find an example for this one - not sure we still use it
		messageText, _ := tokens[0], tokens[1]
		return PathlessRule{RuleText: line, MessageText: messageText}
	} else {
		messageText := tokens[0]
		return PathlessRule{RuleText: line, MessageText: messageText}
	}
}

func buildPathfulPattern(line string) Rule {
	const separator = " "
	const numSplits = 2

	// handle chart yaml file errors in specific template files
	parts := strings.SplitN(line, separator, numSplits)
	if len(parts) == numSplits {
		messagePath, messageText := parts[0], parts[1]
		return Rule{RuleText: line, MessagePath: messagePath, MessageText: messageText}
	} else {
		messagePath := parts[0]
		return Rule{RuleText: line, MessagePath: messagePath, MessageText: NoMessageText}
	}
}

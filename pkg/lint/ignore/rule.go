package ignore

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
)

type Rule struct {
	RuleText    string
	MessagePath string
	MessageText string
}

type LintedMessage struct {
	ChartPath   string
	MessagePath string
	MessageText string
}

func NewRule(ruleText string) *Rule {
	return &Rule{RuleText: ruleText}
}

// ShouldKeepLintedMessage Function used to test the test data in rule_test.go and verify that the ignore capability work as needed
func (r Rule) ShouldKeepLintedMessage(msg LintedMessage) bool {
	cmdIgnorer := RuleLoader{}
	rdr := strings.NewReader(r.RuleText)
	cmdIgnorer.LoadFromReader(rdr)

	actionIgnorer := Ignorer{RuleLoader: &cmdIgnorer}
	return actionIgnorer.ShouldKeepError(fmt.Errorf(msg.MessageText))
}

// LogAttrs Used for troubleshooting and gathering data
func (r Rule) LogAttrs() slog.Attr {
	return slog.Group("Rule",
		slog.String("rule_text", r.RuleText),
		slog.String("key", r.MessagePath),
		slog.String("value", r.MessageText),
	)
}

// Match errors that have a file path in their body with ignorer rules.
// Ignorer rules are built from the lint ignore file
func (r Rule) Match(errText string) *RuleMatch {
	errorFullPath, err := extractFullPathFromError(errText)
	if err != nil {
		return nil
	}

	ignorablePath := r.MessagePath
	ignorableText := r.MessageText
	cleanIgnorablePath := filepath.Clean(ignorablePath)

	if strings.Contains(errorFullPath, cleanIgnorablePath) {
		if strings.Contains(errText, ignorableText) {
			return &RuleMatch{ErrText: errText, RuleText: r.RuleText}
		}
	}

	return nil
}

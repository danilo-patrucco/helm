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

func (r Rule) ShouldKeepLintedMessage(msg LintedMessage) bool {
	cmdIgnorer := RuleLoader{}
	rdr := strings.NewReader(r.RuleText)
	cmdIgnorer.LoadFromReader(rdr)

	actionIgnorer := ActionIgnorer{RuleLoader: &cmdIgnorer}
	return actionIgnorer.ShouldKeepError(fmt.Errorf(msg.MessageText))
}

func (r Rule) LogAttrs() slog.Attr {
	return slog.Group("Rule",
		slog.String("rule_text", r.RuleText),
		slog.String("key", r.MessagePath),
		slog.String("value", r.MessageText),
	)
}

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

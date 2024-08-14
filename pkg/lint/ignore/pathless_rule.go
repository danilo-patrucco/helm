package ignore

import (
	"path/filepath"
	"strings"
)

type PathlessRule struct {
	RuleText    string
	MessageText string
}

// Match errors that have no file path in their body with ignorer rules.
// An examples of errors with no file path in their body is chart metadata errors `chart metadata is missing these dependencies`
func (pr PathlessRule) Match(errText string) *RuleMatch {
	ignorableError := pr.MessageText
	parts := strings.SplitN(ignorableError, ":", 2)
	prefix := strings.TrimSpace(parts[0])

	if match, _ := filepath.Match(ignorableError, errText); match {
		return &RuleMatch{ErrText: errText, RuleText: pr.RuleText}
	}

	if matched, _ := filepath.Match(prefix, errText); matched {
		return &RuleMatch{ErrText: errText, RuleText: pr.RuleText}
	}

	return nil
}

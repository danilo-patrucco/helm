package ignore

import (
	"path/filepath"
	"strings"
)

type PathlessRule struct {
	RuleText    string
	MessageText string
}

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

package ignore

import (
	"fmt"
	"helm.sh/helm/v3/pkg/lint"
	"helm.sh/helm/v3/pkg/lint/support"
	"strings"
)

type Rule struct {
	RuleText string
}

type LintedMessage struct {
	ChartPath   string
	MessagePath string
	MessageText string
}

func NewRule(ruleText string) *Rule {
	return &Rule{ruleText}
}

func (r Rule) ShouldKeepLintedMessage(msg LintedMessage) bool {
	ignorer := lint.Ignorer{}

	rdr := strings.NewReader(r.RuleText)
	ignorer.LoadFromReader(rdr)

	testTheseMessages := []support.Message{
		{
			Severity: 3,
			Path:     msg.MessagePath,
			Err:      fmt.Errorf(msg.MessageText),
		},
	}

	keptMessages := ignorer.FilterMessages(testTheseMessages)
	return len(keptMessages) > 0
}

func (r Rule) ShouldKeepLintedError(msg LintedMessage) bool {
	ignorer := lint.Ignorer{}

	rdr := strings.NewReader(r.RuleText)
	ignorer.LoadFromReader(rdr)

	keptMessagesbool := ignorer.IsIgnoredPathlessError(msg.MessageText)
	return keptMessagesbool
}

package ignore

import (
	"fmt"
	"helm.sh/helm/v3/pkg/lint"
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
	cmdIgnorer := lint.Ignorer{}
	rdr := strings.NewReader(r.RuleText)
	cmdIgnorer.LoadFromReader(rdr)

	actionIgnorer := ActionIgnorer{ CmdIgnorer: &cmdIgnorer }
	return actionIgnorer.ShouldKeepError(fmt.Errorf(msg.MessageText))
}
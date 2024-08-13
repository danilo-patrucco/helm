package ignore

import (
	"fmt"
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
	cmdIgnorer := CmdIgnorer{}
	rdr := strings.NewReader(r.RuleText)
	cmdIgnorer.LoadFromReader(rdr)

	actionIgnorer := ActionIgnorer{ CmdIgnorer: &cmdIgnorer }
	return actionIgnorer.ShouldKeepError(fmt.Errorf(msg.MessageText))
}
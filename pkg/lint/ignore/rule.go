package ignore

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
	return true
}

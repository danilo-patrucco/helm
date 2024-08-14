package ignore

import "log/slog"

type MatchesErrors interface {
	Match(string) *RuleMatch
}

type RuleMatch struct {
	ErrText  string
	RuleText string
}

func (rm RuleMatch) LogAttrs() slog.Attr {
	return slog.Group("rule_match", slog.String("err_text", rm.ErrText), slog.String("rule_text", rm.RuleText))
}

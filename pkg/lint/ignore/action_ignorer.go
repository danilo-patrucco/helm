package ignore

import (
	"helm.sh/helm/v3/pkg/lint/support"
	"log/slog"
)

type ActionIgnorer struct {
	ChartPath  string
	Rules      []Rule
	logger     *slog.Logger
	RuleLoader *RuleLoader
}

func Ignorer(chartPath string, lintIgnorePath string, debugLogFn func(string, ...interface{})) (*ActionIgnorer, error) {
	cmdIgnorer, err := NewRuleLoader(chartPath, lintIgnorePath, debugLogFn)
	if err != nil {
		return nil, err
	}

	return &ActionIgnorer{ChartPath: chartPath, RuleLoader: cmdIgnorer}, nil
}

func (ai *ActionIgnorer) FilterMessages(messages []support.Message) []support.Message {
	out := make([]support.Message, 0, len(messages))
	for _, msg := range messages {
		if ai.ShouldKeepError(msg.Err) {
			out = append(out, msg)
		}
	}
	return out
}

func (ai *ActionIgnorer) ShouldKeepError(err error) bool {
	errText := err.Error()

	// if any of our Matchers match the rule, we can discard it
	for _, rule := range ai.RuleLoader.Matchers {
		match := rule.Match(errText)
		if match != nil {
			ai.RuleLoader.Debug("lint ignore rule matched", match.LogAttrs())
			return false
		}
	}

	// if we can't find a reason to discard it, we keep it
	return true
}

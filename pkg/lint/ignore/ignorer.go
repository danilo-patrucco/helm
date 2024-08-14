package ignore

import (
	"helm.sh/helm/v3/pkg/lint/support"
	"log/slog"
)

type Ignorer struct {
	ChartPath  string
	Rules      []Rule
	logger     *slog.Logger
	RuleLoader *RuleLoader
}

// Ignorer is used to create the ignorer object that contains
func NewActionIgnorer(chartPath string, lintIgnorePath string, debugLogFn func(string, ...interface{})) (*Ignorer, error) {
	cmdIgnorer, err := NewRuleLoader(chartPath, lintIgnorePath, debugLogFn)
	if err != nil {
		return nil, err
	}

	return &Ignorer{ChartPath: chartPath, RuleLoader: cmdIgnorer}, nil
}

// FilterMessages Verify what messages can be kept in the output, using also the error as a verification (calling ShouldKeepError)
func (i *Ignorer) FilterMessages(messages []support.Message) []support.Message {
	out := make([]support.Message, 0, len(messages))
	for _, msg := range messages {
		if i.ShouldKeepError(msg.Err) {
			out = append(out, msg)
		}
	}
	return out
}

// ShouldKeepError is used to verify if the error associated with the message need to be kept, or it can be ignored, called by FilterMessages and in the pkg/action/lint.go Run main function
func (i *Ignorer) ShouldKeepError(err error) bool {
	errText := err.Error()

	// if any of our Matchers match the rule, we can discard it
	for _, rule := range i.RuleLoader.Matchers {
		match := rule.Match(errText)
		if match != nil {
			i.RuleLoader.Debug("lint ignore rule matched", match.LogAttrs())
			return false
		}
	}

	// if we can't find a reason to discard it, we keep it
	return true
}

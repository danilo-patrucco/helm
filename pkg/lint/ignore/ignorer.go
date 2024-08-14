package ignore

import (
	"helm.sh/helm/v3/pkg/lint/support"
	"log/slog"
	"path/filepath"
	"strings"
)

type Ignorer struct {
	ChartPath  string
	Rules      []Rule
	logger     *slog.Logger
	RuleLoader *RuleLoader
}

type PathlessRule struct {
	RuleText    string
	MessageText string
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

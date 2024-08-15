package ignore

import (
	"fmt"
	"helm.sh/helm/v3/pkg/lint/support"
	"log"
	"log/slog"
	"path/filepath"
	"strings"
)

var debugFn func(format string, v ...interface{})

type Ignorer struct {
	ChartPath string
	Matchers  []MatchesErrors
}
type Rule struct {
	RuleText    string
	MessagePath string
	MessageText string
}

type PathlessRule struct {
	RuleText    string
	MessageText string
}

func NewIgnorer(chartPath string, lintIgnorePath string, debugLogFn func(string, ...interface{})) (*Ignorer, error) {
	matchers, err := LoadFromFilePath(chartPath, lintIgnorePath, debugLogFn)
	if err != nil {
		return nil, err
	}

	return &Ignorer{ChartPath: chartPath, Matchers: matchers}, nil
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
	for _, rule := range i.Matchers {
		if rule.Match(errText) {
			debug("lint ignore rule matched matcher=%v, err=%s", rule, err.Error())
			return false
		}
	}

	// if we can't find a reason to discard it, we keep it
	return true
}

type MatchesErrors interface {
	Match(string) bool
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
func (pr PathlessRule) Match(errText string) bool {
	matchableParts := strings.SplitN(pr.MessageText, ":", 2)
	matchablePrefix := strings.TrimSpace(matchableParts[0])

	if match, _ := filepath.Match(pr.MessageText, errText); match {
		debug("lint ignore match: errText=%s, ruleText=%s", errText, pr.RuleText)
		return true
	}
	if matched, _ := filepath.Match(matchablePrefix, errText); matched {
		debug("lint ignore match: errText=%s, ruleText=%s", errText, pr.RuleText)
		return true
	}

	return false
}

// LogAttrs Used for troubleshooting and gathering data
func (r Rule) LogAttrs() slog.Attr {
	return slog.Group("Rule",
		slog.String("rule_text", r.RuleText),
		slog.String("key", r.MessagePath),
		slog.String("value", r.MessageText),
	)
}

// Match errors that have a file path in their body with ignorer rules.
func (r Rule) Match(errText string) bool {
	pathToOffendingFile, err := pathToOffendingFile(errText)
	if err != nil {
		return false
	}

	cleanRulePath := filepath.Clean(r.MessagePath)

	if strings.Contains(pathToOffendingFile, cleanRulePath) {
		if strings.Contains(errText, r.MessageText) {
			debug("lint ignore match: errText=%s, ruleText=%s", errText, r.RuleText)
			return true
		}
	}

	return false
}

var defaultDebugFn = func(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}

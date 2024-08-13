package ignore

import (
	"helm.sh/helm/v3/pkg/lint"
	"helm.sh/helm/v3/pkg/lint/support"
	"log/slog"
	"os"
	"strings"
)

type ActionIgnorer struct {
	ChartPath string
	Rules     []Rule
	logger    *slog.Logger
	CmdIgnorer *lint.Ignorer
}

func (ai *ActionIgnorer) LoadFromRuleText(ruleText string) {
	ai.CmdIgnorer = &lint.Ignorer{}
	rdr := strings.NewReader(ruleText)
	ai.CmdIgnorer.LoadFromReader(rdr)
	return
}

func (ai *ActionIgnorer) FilterMessages(messages []support.Message) []support.Message {
	out := make([]support.Message, 0, len(messages))
	for _, msg := range messages {
		if ai.ShouldKeepMessage(msg) {
			out = append(out, msg)
		}
	}
	return out
}

func (ai *ActionIgnorer) ShouldKeepError(err error) bool {
	errText := err.Error()

	// log it
	logAttr := slog.Group("Err", slog.String("text", errText))
	ai.Info("action/lint/Run captured an error", "Kind", "Error", logAttr)

	// do not keep if it matches our basic rules
	if ai.CmdIgnorer.IsIgnorable(errText) {
		return false
	}

	// do not keep if it matches the pathless error rules
	if ai.CmdIgnorer.IsIgnoredPathlessError(errText) {
		return false
	}

	// keep it!
	return true
}

func (ai *ActionIgnorer) ShouldKeepMessage(message support.Message) bool {
	ai.Info("action/lint/Run captured a message", "Kind", "Message", message.LogAttrs())

	return ai.ShouldKeepError(message.Err)
}

func (ai *ActionIgnorer) Info(msg string, args ...any) {
	if ai.logger == nil {
		ai.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	baseAttrs := slog.Group("Chart",
		slog.String("Path", ai.ChartPath),
	)

	ai.logger.With(baseAttrs).Info(msg, args...)
}

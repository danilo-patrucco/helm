package ignore

import (
	"helm.sh/helm/v3/pkg/lint/support"
	"log/slog"
	"os"
	"strings"
)

type ActionIgnorer struct {
	ChartPath string
	Rules     []Rule
	logger    *slog.Logger
	CmdIgnorer *CmdIgnorer
}

func NewActionIgnorer(chartPath string, lintIgnorePath string, debugLogFn func(string, ...interface{})) (*ActionIgnorer, error) {
	cmdIgnorer, err := NewCmdIgnorer(chartPath, lintIgnorePath, debugLogFn)
	if err != nil {
		return nil, err
	}

	return &ActionIgnorer{ ChartPath: chartPath, CmdIgnorer: cmdIgnorer }, nil
}

func (ai *ActionIgnorer) LoadFromRuleText(ruleText string) {
	ai.CmdIgnorer = &CmdIgnorer{}
	rdr := strings.NewReader(ruleText)
	ai.CmdIgnorer.LoadFromReader(rdr)
	return
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

func (ai *ActionIgnorer) Info(msg string, args ...any) {
	if ai.logger == nil {
		ai.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	baseAttrs := slog.Group("Chart",
		slog.String("Path", ai.ChartPath),
	)

	ai.logger.With(baseAttrs).Info(msg, args...)
}

package ignore

import (
	"helm.sh/helm/v3/pkg/lint/support"
	"log/slog"
	"os"
)

type Ignorer struct {
	ChartPath string
	Rules     []Rule
	logger    *slog.Logger
}

func (i *Ignorer) FilterMessages(messages []support.Message) []support.Message {
	out := make([]support.Message, 0, len(messages))
	for _, msg := range messages {
		if !i.ShouldKeepMessage(msg) {
			continue
		}
		out = append(out, msg)
	}
	return out
}

func (i *Ignorer) ShouldKeepError(err error) bool {
	logAttr := slog.Group("Err", slog.String("text", err.Error()))
	i.Info("action/lint/Run captured an error", "Kind", "Error", logAttr)
	return true
}

func (i *Ignorer) ShouldKeepMessage(message support.Message) bool {
	i.Info("action/lint/Run captured a message", "Kind", "Message", message.LogAttrs())
	return true
}

func (i *Ignorer) Info(msg string, args ...any) {
	if i.logger == nil {
		i.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	baseAttrs := slog.Group("Chart",
		slog.String("Path", i.ChartPath),
	)

	i.logger.With(baseAttrs).Info(msg, args...)
}

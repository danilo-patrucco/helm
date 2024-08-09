package ignore

import (
	"helm.sh/helm/v3/pkg/lint/support"
)

type Ignorer struct {

}

func (i *Ignorer) FilterErrors(chartPath string, errs []error) []error {
	out := make([]error, 0, len(errs))
	for _, err := range errs {
		if !i.ShouldKeepError(chartPath, err) {
			continue
		}

		out = append(out, err)
	}
	return out
}

func (i *Ignorer) FilterMessages(chartPath string, messages []support.Message) []support.Message {
	out := make([]support.Message, 0, len(messages))
	for _, msg := range messages {
		if !i.ShouldKeepMessage(chartPath, msg) {
			continue
		}
		out = append(out, msg)
	}
	return out
}

func (i *Ignorer) ShouldKeepError(chartPath string, err error) bool {
	return true
}

func (i *Ignorer) ShouldKeepMessage(chartPath string, msg support.Message) bool {
	return true
}

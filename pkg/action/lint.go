/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package action

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/lint"
	"helm.sh/helm/v3/pkg/lint/support"
)

// Lint is the action for checking that the semantics of a chart are well-formed.
//
// It provides the implementation of 'helm lint'.
type Lint struct {
	Strict        bool
	Namespace     string
	WithSubcharts bool
	Quiet         bool
	KubeVersion   *chartutil.KubeVersion
}

// LintResult is the result of Lint
type LintResult struct {
	TotalChartsLinted int
	Messages          []support.Message
	Errors            []error
}

// NewLint creates a new Lint object with the given configuration.
func NewLint() *Lint {
	return &Lint{}
}

// Run executes 'helm Lint' against the given chart.
func (l *Lint) Run(paths []string, vals map[string]interface{}) *LintResult {
	lowestTolerance := support.ErrorSev
	if l.Strict {
		lowestTolerance = support.WarningSev
	}
	result := &LintResult{}
	for chartIndex, path := range paths {
		linter, err := lintChart(path, vals, l.Namespace, l.KubeVersion)
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}
		logCapturedErrors(chartIndex, path, result.Errors)
		result.Errors = filterErrors(chartIndex, path, result.Errors)

		result.Messages = append(result.Messages, linter.Messages...)

		logCapturedMessages(chartIndex, path, result.Messages)
		result.Messages = filterMessages(chartIndex, path, result.Messages)

		result.TotalChartsLinted++
		for _, msg := range linter.Messages {
			if msg.Severity >= lowestTolerance {
				slog.Info("action/lint/Run is promoting a message to Error", "chartIndex", chartIndex, "path", path, "lowestTolerance", lowestTolerance, msg.LogAttrs())
				result.Errors = append(result.Errors, msg.Err)
			}
		}
	}
	return result
}

func filterMessages(chartIndex int, path string, messages []support.Message) []support.Message {
	out := make([]support.Message, 0, len(messages))
	for _, msg := range messages {
		// TODO: filter some messages based on content
		out = append(out, msg)
	}
	return messages
}

func filterErrors(chartIndex int, path string, errs []error) []error {
	out := make([]error, 0, len(errs))
	for _, err := range errs {
		// TODO: filter some errors based on content
		out = append(out, err)
	}
	return errs
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func logCapturedErrors(chartIndex int, path string, errors []error) {
	for _, err := range errors {
		logAttr := slog.Group("Err", slog.String("text", err.Error()))
		logger.Info("action/lint/Run captured an error", "chartIndex", chartIndex, "path", path, logAttr)
	}
}

func logCapturedMessages(chartIndex int, path string, messages []support.Message) {
	for _, msg := range messages {
		logger.Info("action/lint/Run captured a message", "chartIndex", chartIndex, "path", path, msg.LogAttrs())
	}
}

// HasWarningsOrErrors checks is LintResult has any warnings or errors
func HasWarningsOrErrors(result *LintResult) bool {
	for _, msg := range result.Messages {
		if msg.Severity > support.InfoSev {
			return true
		}
	}
	return len(result.Errors) > 0
}

func lintChart(path string, vals map[string]interface{}, namespace string, kubeVersion *chartutil.KubeVersion) (support.Linter, error) {
	var chartPath string
	linter := support.Linter{}

	if strings.HasSuffix(path, ".tgz") || strings.HasSuffix(path, ".tar.gz") {
		tempDir, err := os.MkdirTemp("", "helm-lint")
		if err != nil {
			return linter, errors.Wrap(err, "unable to create temp dir to extract tarball")
		}
		defer os.RemoveAll(tempDir)

		file, err := os.Open(path)
		if err != nil {
			return linter, errors.Wrap(err, "unable to open tarball")
		}
		defer file.Close()

		if err = chartutil.Expand(tempDir, file); err != nil {
			return linter, errors.Wrap(err, "unable to extract tarball")
		}

		files, err := os.ReadDir(tempDir)
		if err != nil {
			return linter, errors.Wrapf(err, "unable to read temporary output directory %s", tempDir)
		}
		if !files[0].IsDir() {
			return linter, errors.Errorf("unexpected file %s in temporary output directory %s", files[0].Name(), tempDir)
		}

		chartPath = filepath.Join(tempDir, files[0].Name())
	} else {
		chartPath = path
	}

	// Guard: Error out if this is not a chart.
	if _, err := os.Stat(filepath.Join(chartPath, "Chart.yaml")); err != nil {
		return linter, errors.Wrap(err, "unable to check Chart.yaml file in chart")
	}

	return lint.AllWithKubeVersion(chartPath, vals, namespace, kubeVersion), nil
}

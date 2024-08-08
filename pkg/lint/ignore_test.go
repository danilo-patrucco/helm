package lint

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/lint/support"
	"path/filepath"
	"testing"
	"text/template"
)

func TestNewIgnorer(t *testing.T) {
	chartPath := "rules/testdata/withsubchartlintignore"
	ignoreFilePath := filepath.Join(chartPath, ".helmlintignore")
	ignorer := NewIgnorer(chartPath, ignoreFilePath, func(format string, args ...interface{}) {
		t.Logf(format, args...)
	})
	assert.NotNil(t, ignorer, "Ignorer should not be nil")
	assert.NotEmpty(t, ignorer.Patterns, "Expected patterns to be loaded from the file, but none were found")
	//if len(ignorer.Patterns) == 0 {
	//	t.Errorf("Expected patterns to be loaded from the file, but none were found")
	//}
}

func TestFilterErrors(t *testing.T) {
	// arrange
	badYamlPath := "/path/to/chart/templates/bad-template.yaml"
	errIgnorableHasPath := fmt.Errorf("test: %s: ignore this error", badYamlPath)
	errWithNoPath := fmt.Errorf("keep this error")

	// act
	ignorer := &Ignorer{
		PathlessErrorPatterns: map[string][]string{
			badYamlPath: {"ignore this error"},
		},
	}

	// assert
	given := []error{errIgnorableHasPath, errWithNoPath}
	got := ignorer.FilterErrors(given)
	assert.Contains(t, got, errWithNoPath)
	assert.NotContains(t, got, errIgnorableHasPath)
}

func TestFilterNoPathErrors(t *testing.T) {
	ignorer := &Ignorer{
		PathlessErrorPatterns: map[string][]string{
			"chart error": {"this should be ignored"},
		},
	}
	messages := []support.Message{}
	errors := []error{fmt.Errorf("this should be ignored"), fmt.Errorf("this should be kept")}
	filteredMessages, filteredErrors := ignorer.FilterNoPathErrors(messages, errors)
	assert.Empty(t, filteredErrors)
	assert.NotEmpty(t, filteredMessages)
}

func TestMatchNoPathError(t *testing.T) {
	ignorer := &Ignorer{
		PathlessErrorPatterns: map[string][]string{
			"generic error": {"ignore this"},
		},
	}
	result := ignorer.IsIgnoredPathlessError("ignore this")
	assert.False(t, result)
}

func TestDebug(t *testing.T) {
	var captured string
	debugFn := func(format string, args ...interface{}) {
		captured = fmt.Sprintf(format, args...)
	}
	ignorer := &Ignorer{
		debugFnOverride: debugFn,
	}
	ignorer.Debug("test %s", "debug")
	assert.Equal(t, "test debug", captured)
}

func TestMatch(t *testing.T) {
	ignorer := &Ignorer{
		Patterns: map[string][]string{
			"rules/testdata/withsubchartlintignore/charts/subchart/templates/subchart.yaml": {"<include \"this.is.test.data\" .>"},
		},
	}

	assert.True(t, ignorer.isIgnorable("error pattern in rules/testdata/withsubchartlintignore/charts/subchart/templates/subchart.yaml"))
	assert.False(t, ignorer.isIgnorable("this should not match"))
}

func TestFilterIgnoredMessages(t *testing.T) {
	type args struct {
		messages       []support.Message
		ignorePatterns map[string][]string
	}
	tests := []struct {
		name string
		args args
		want []support.Message
	}{
		{
			name: "should filter ignored messages only",
			args: args{
				messages: []support.Message{
					{
						Severity: 3,
						Path:     "templates/",
						Err: template.ExecError{
							Name: "certmanager-issuer/templates/rbac-config.yaml",
							Err:  fmt.Errorf(`template: certmanager-issuer/templates/rbac-config.yaml:1:67: executing "certmanager-issuer/templates/rbac-config.yaml" at <.Values.global.ingress>: nil pointer evaluating interface {}.ingress`),
						},
					},
					{
						Severity: 1,
						Path:     "values.yaml",
						Err:      fmt.Errorf("file does not exist"),
					},
					{
						Severity: 1,
						Path:     "Chart.yaml",
						Err:      fmt.Errorf("icon is recommended"),
					},
				},
				ignorePatterns: map[string][]string{
					"certmanager-issuer/templates/rbac-config.yaml": {
						"<.Values.global.ingress>",
					},
				},
			},
			want: []support.Message{
				{
					Severity: 1,
					Path:     "values.yaml",
					Err:      fmt.Errorf("file does not exist"),
				},
				{
					Severity: 1,
					Path:     "Chart.yaml",
					Err:      fmt.Errorf("icon is recommended"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ignorer := Ignorer{Patterns: tt.args.ignorePatterns}
			got := ignorer.FilterMessages(tt.args.messages)
			assert.Equalf(t, tt.want, got, "FilterMessages(%v, %v)", tt.args.messages, tt.args.ignorePatterns)
		})
	}
}

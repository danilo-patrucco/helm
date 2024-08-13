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
	if len(ignorer.Patterns) == 0 {
		t.Errorf("Expected patterns to be loaded from the file, but none were found")
	}
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

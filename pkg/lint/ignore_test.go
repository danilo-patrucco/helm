package lint

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/lint/support"
	"strings"
	"testing"
	"text/template"
)

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

func TestIgnorer_PatternMatching(t *testing.T) {
	tests := []struct {
		name               string
		ignoreFileContents string
		givenMessages      []support.Message
		givenErrors        []error
		wantMessages       []support.Message
		wantErrors         []error
	}{
		{
			name:               "should suppress errors marked via " + errorPatternPrefix,
			ignoreFileContents: "error_lint_ignore=chart metadata is missing these dependencies:*",
			givenMessages: []support.Message{
				{
					Severity: 3,
					Path:     "/fake/path/goes/here",
					Err:      fmt.Errorf("chart metadata is missing these dependencies: gitaly,mailroom,migrations,sidekiq,webservice,toolbox,geo-logcursor,gitlab-exporter,gitlab-shell,kas,spamcheck"),
				},
			},
			givenErrors: []error{
				fmt.Errorf("chart metadata is missing these dependencies: gitaly,mailroom,migrations,sidekiq,webservice,toolbox,geo-logcursor,gitlab-exporter,gitlab-shell,kas,spamcheck"),
			},
			wantMessages: []support.Message{},
			wantErrors:   []error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ignorer := newIgnorerFromReader(strings.NewReader(tt.ignoreFileContents))
			gotMessages, gotErrors := ignorer.FilterNoPathErrors(tt.givenMessages, tt.givenErrors)

			assert.Equal(t, tt.wantMessages, gotMessages)
			assert.Equal(t, tt.wantErrors, gotErrors)
		})
	}
}

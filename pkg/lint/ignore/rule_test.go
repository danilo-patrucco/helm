package ignore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const todoRule = "TODO MAKE A RULE FOR THIS"

var fakeLM = LintedMessage{}

func TestRule_ShouldKeepMessage(t *testing.T) {
	type testCase struct {
		Description string
		RuleText    string
		Ignorables  []LintedMessage
		Keepables   []LintedMessage
	}

	testCases := []testCase{
		{
			Description: "subchart template not defined",
			RuleText:    todoRule,
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab",
				MessagePath: "templates/",
				MessageText: "template: gitlab/charts/webservice/templates/tests/tests.yaml:5:20: executing \"gitlab/charts/webservice/templates/tests/tests.yaml\" at <{{template \"fullname\" .}}>: template \"fullname\" not defined",
			}},
		},
		{
			Description: "subchart template include template not found",
			RuleText:    todoRule,
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab/charts/gitaly",
				MessagePath: "templates/",
				MessageText: "template: gitaly/templates/statefulset.yml:1:11: executing \"gitaly/templates/statefulset.yml\" at <include \"gitlab.gitaly.includeInternalResources\" $>: error calling include: template: no template \"gitlab.gitaly.includeInternalResources\" associated with template \"gotpl\"",
			}},
		},
		{
			Description: "subchart template evaluation has a nil pointer",
			RuleText:    todoRule,
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab/charts/gitlab-exporter",
				MessagePath: "templates/",
				MessageText: "template: gitlab-exporter/templates/serviceaccount.yaml:1:57: executing \"gitlab-exporter/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		},
		{
			Description: "subchart icon is recommended",
			RuleText:    todoRule,
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab-zoekt-1.4.0.tgz",
				MessagePath: "Chart.yaml",
				MessageText: "icon is recommended",
			}},
		},
		{
			Description: "subchart values file does not exist",
			RuleText:    todoRule,
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gluon-0.5.0.tgz",
				MessagePath: "values.yaml",
				MessageText: "file does not exist",
			}},
		},
		{
			Description: "subchart metadata missing dependencies",
			RuleText:    todoRule,
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab",
				MessagePath: "/Users/daniel/radius/bb/gitlab/chart/charts/gitlab",
				MessageText: "chart metadata is missing these dependencies: sidekiq,spamcheck,gitaly,gitlab-shell,kas,mailroom,migrations,toolbox,geo-logcursor,gitlab-exporter,webservice",
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			rule := NewRule(testCase.RuleText)

			for _, ignorableMessage := range testCase.Ignorables {
				assert.False(t, rule.ShouldKeepLintedMessage(ignorableMessage), testCase.Description)
			}

			for _, keepableMessage := range testCase.Ignorables {
				assert.True(t, rule.ShouldKeepLintedMessage(keepableMessage), testCase.Description)
			}
		})
	}
}

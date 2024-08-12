package ignore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRule_ShouldKeepMessage(t *testing.T) {
	type testCase struct {
		Scenario   string
		RuleText   string
		Ignorables []LintedMessage
	}

	testCases := []testCase{
		{
			Scenario: "subchart template not defined",
			RuleText: "gitlab/charts/webservice/templates/tests/tests.yaml <{{template \"fullname\" .}}>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab",
				MessagePath: "templates/",
				MessageText: "template: gitlab/charts/webservice/templates/tests/tests.yaml:5:20: executing \"gitlab/charts/webservice/templates/tests/tests.yaml\" at <{{template \"fullname\" .}}>: template \"fullname\" not defined",
			}},
		},
		{
			Scenario: "subchart template include template not found",
			RuleText: "gitaly/templates/statefulset.yml <include \"gitlab.gitaly.includeInternalResources\" $>\n",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab/charts/gitaly",
				MessagePath: "templates/",
				MessageText: "template: gitaly/templates/statefulset.yml:1:11: executing \"gitaly/templates/statefulset.yml\" at <include \"gitlab.gitaly.includeInternalResources\" $>: error calling include: template: no template \"gitlab.gitaly.includeInternalResources\" associated with template \"gotpl\"",
			}},
		},
		{
			Scenario: "subchart template evaluation has a nil pointer",
			RuleText: "gitlab-exporter/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>\n",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab/charts/gitlab-exporter",
				MessagePath: "templates/",
				MessageText: "template: gitlab-exporter/templates/serviceaccount.yaml:1:57: executing \"gitlab-exporter/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		},
		{
			Scenario: "subchart icon is recommended",
			RuleText: "",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab-zoekt-1.4.0.tgz",
				MessagePath: "Chart.yaml",
				MessageText: "icon is recommended",
			}},
		},
		{
			Scenario: "subchart values file does not exist",
			RuleText: "TODO MAKE A RULE FOR THIS",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gluon-0.5.0.tgz",
				MessagePath: "values.yaml",
				MessageText: "file does not exist",
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Scenario, func(t *testing.T) {
			rule := NewRule(testCase.RuleText)

			for _, ignorableMessage := range testCase.Ignorables {
				assert.False(t, rule.ShouldKeepLintedMessage(ignorableMessage), testCase.Scenario)
			}

			keepableMessage := LintedMessage{
				ChartPath:   "a/memorable/path",
				MessagePath: "wow/",
				MessageText: "incredible: something just happened",
			}
			assert.True(t, rule.ShouldKeepLintedMessage(keepableMessage))
		})
	}
}

func TestRule_ShouldKeepErrors(t *testing.T) {
	type testCase struct {
		Scenario   string
		RuleText   string
		Ignorables []LintedMessage
	}

	testCases := []testCase{
		{
			Scenario: "subchart metadata missing dependencies",
			RuleText: "error_lint_ignore=chart metadata is missing these dependencies**",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab",
				MessagePath: "gitlab/chart/charts/gitlab",
				MessageText: "chart metadata is missing these dependencies: sidekiq,spamcheck,gitaly,gitlab-shell,kas,mailroom,migrations,toolbox,geo-logcursor,gitlab-exporter,webservice",
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Scenario, func(t *testing.T) {
			rule := NewRule(testCase.RuleText)

			for _, ignorableMessage := range testCase.Ignorables {
				assert.True(t, rule.ShouldKeepLintedMessage(ignorableMessage), testCase.Scenario)
			}

			keepableMessage := LintedMessage{
				ChartPath:   "a/memorable/path",
				MessagePath: "wow/",
				MessageText: "this is wrong",
			}
			assert.False(t, rule.ShouldKeepLintedError(keepableMessage))
		})
	}
}

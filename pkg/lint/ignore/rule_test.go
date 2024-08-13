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
		}, {
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
		}, {
			Scenario: "webservice path only",
			RuleText: "webservice/templates/tests/tests.yaml",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/webservice/templates/tests/tests.yaml",
				MessagePath: "templates/",
				MessageText: "template: webservice/templates/tests/tests.yaml:8:8: executing \"webservice/templates/tests/tests.yaml\" at <include \"gitlab.standardLabels\" .>: error calling include: template: no template \"gitlab.standardLabels\" associated with template \"gotpl\"",
			}},
		}, {
			Scenario: "geo-logcursor path only",
			RuleText: "geo-logcursor/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/geo-logcursor/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: geo-logcursor/templates/serviceaccount.yaml:1:57: executing \"geo-logcursor/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		}, {
			Scenario: "webservice path only",
			RuleText: "webservice/templates/service.yaml <include \"fullname\" .>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/webservice/templates/service.yaml",
				MessagePath: "templates/",
				MessageText: "template: gitlab/charts/gitlab/charts/webservice/templates/service.yaml:14:11: executing \"gitlab/charts/gitlab/charts/webservice/templates/service.yaml\" at <include \"fullname\" .>: error calling include: template: gitlab/templates/_helpers.tpl:14:27: executing \"fullname\" at <.Chart.Name>: nil pointer evaluating interface {}.Name",
			}},
		}, {
			Scenario: "certmanager-issuer path only",
			RuleText: "certmanager-issuer/templates/rbac-config.yaml <.Values.global.ingress>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/certmanager-issuer/templates/rbac-config.yaml",
				MessagePath: "templates/",
				MessageText: "template: certmanager-issuer/templates/rbac-config.yaml:1:67: executing \"certmanager-issuer/templates/rbac-config.yaml\" at <.Values.global.ingress>: nil pointer evaluating interface {}.ingress",
			}},
		}, {
			Scenario: "gitlab-pages path only",
			RuleText: "gitlab-pages/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/gitlab-pages/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: gitlab-pages/templates/serviceaccount.yaml:1:57: executing \"gitlab-pages/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		}, {
			Scenario: "gitlab-shell path only",
			RuleText: "gitlab-shell/templates/traefik-tcp-ingressroute.yaml <.Values.global.ingress.provider>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/gitlab-shell/templates/traefik-tcp-ingressroute.yaml",
				MessagePath: "templates/",
				MessageText: "template: gitlab-shell/templates/traefik-tcp-ingressroute.yaml:2:17: executing \"gitlab-shell/templates/traefik-tcp-ingressroute.yaml\" at <.Values.global.ingress.provider>: nil pointer evaluating interface {}.provider",
			}},
		}, {
			Scenario: "kas path only",
			RuleText: "kas/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/kas/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: kas/templates/serviceaccount.yaml:1:57: executing \"kas/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		}, {
			Scenario: "kas path only",
			RuleText: "kas/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/kas/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: kas/templates/serviceaccount.yaml:1:57: executing \"kas/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		}, {
			Scenario: "mailroom path only",
			RuleText: "mailroom/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/mailroom/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: mailroom/templates/serviceaccount.yaml:1:57: executing \"mailroom/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		}, {
			Scenario: "migrations path only",
			RuleText: "migrations/templates/job.yaml <include (print $.Template.BasePath \"/_serviceaccountspec.yaml\") .>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/migrations/templates/job.yaml",
				MessagePath: "templates/",
				MessageText: "template: migrations/templates/job.yaml:2:3: executing \"migrations/templates/job.yaml\" at <include (print $.Template.BasePath \"/_serviceaccountspec.yaml\") .>: error calling include: template: migrations/templates/_serviceaccountspec.yaml:1:57: executing \"migrations/templates/_serviceaccountspec.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		}, {
			Scenario: "praefect path only",
			RuleText: "praefect/templates/statefulset.yaml <.Values.global.image>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/praefect/templates/statefulset.yaml",
				MessagePath: "templates/",
				MessageText: "template: praefect/templates/statefulset.yaml:1:38: executing \"praefect/templates/statefulset.yaml\" at <.Values.global.image>: nil pointer evaluating interface {}.image",
			}},
		}, {
			Scenario: "sidekiq path only",
			RuleText: "sidekiq/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/sidekiq/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: sidekiq/templates/serviceaccount.yaml:1:57: executing \"sidekiq/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		}, {
			Scenario: "spamcheck path only",
			RuleText: "spamcheck/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/spamcheck/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: spamcheck/templates/serviceaccount.yaml:1:57: executing \"spamcheck/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.serviceAccount",
			}},
		}, {
			Scenario: "toolbox path only",
			RuleText: "toolbox/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/toolbox/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: toolbox/templates/serviceaccount.yaml:1:57: executing \"toolbox/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
			}},
		}, {
			Scenario: "minio path only",
			RuleText: "minio/templates/pdb.yaml <{{template \"gitlab.pdb.apiVersion\" $pdbCfg}}>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/minio/templates/pdb.yaml",
				MessagePath: "templates/",
				MessageText: "template: minio/templates/pdb.yaml:3:24: executing \"minio/templates/pdb.yaml\" at <{{template \"gitlab.pdb.apiVersion\" $pdbCfg}}>: template \"gitlab.pdb.apiVersion\" not defined",
			}},
		}, {
			Scenario: "nginx-ingress path only",
			RuleText: "nginx-ingress/templates/admission-webhooks/job-patch/serviceaccount.yaml <.Values.admissionWebhooks.serviceAccount.automountServiceAccountToken>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/nginx-ingress/templates/admission-webhooks/job-patch/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: nginx-ingress/templates/admission-webhooks/job-patch/serviceaccount.yaml:13:40: executing \"nginx-ingress/templates/admission-webhooks/job-patch/serviceaccount.yaml\" at <.Values.admissionWebhooks.serviceAccount.automountServiceAccountToken>: nil pointer evaluating interface {}.serviceAccount",
			}},
		}, {
			Scenario: "registry path only",
			RuleText: "registry/templates/serviceaccount.yaml <.Values.global.serviceAccount.enabled>",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/charts/registry/templates/serviceaccount.yaml",
				MessagePath: "templates/",
				MessageText: "template: registry/templates/serviceaccount.yaml:1:57: executing \"registry/templates/serviceaccount.yaml\" at <.Values.global.serviceAccount.enabled>: nil pointer evaluating interface {}.enabled",
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
		{
			Scenario: "subchart icon is recommended",
			RuleText: "error_lint_ignore=subchart icon is recommended",
			Ignorables: []LintedMessage{{
				ChartPath:   "../gitlab/chart/charts/gitlab-zoekt-1.4.0.tgz",
				MessagePath: "Chart.yaml",
				MessageText: "icon is recommended",
			}},
		},
		{
			Scenario: "subchart values file does not exist",
			RuleText: "error_lint_ignore=file does not exist",
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

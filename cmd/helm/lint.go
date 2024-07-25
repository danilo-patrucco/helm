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

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/lint/rules"
)

var longLintHelp = `
This command takes a path to a chart and runs a series of tests to verify that
the chart is well-formed.

If the linter encounters things that will cause the chart to fail installation,
it will emit [ERROR] messages. If it encounters issues that break with convention
or recommendation, it will emit [WARNING] messages.
`

func newLintCmd(out io.Writer) *cobra.Command {
	client := action.NewLint()
	client.Debug = settings.Debug
	valueOpts := &values.Options{}
	var kubeVersion string
	var lintIgnoreFile string

	cmd := &cobra.Command{
		Use:   "lint PATH",
		Short: "examine a chart for possible issues",
		Long:  longLintHelp,
		RunE: func(_ *cobra.Command, args []string) error {
			paths := []string{"."}
			if len(args) > 0 {
				paths = args
			}
			if kubeVersion != "" {
				parsedKubeVersion, err := chartutil.ParseKubeVersion(kubeVersion)
				if err != nil {
					return fmt.Errorf("invalid kube version '%s': %s", kubeVersion, err)
				}
				client.KubeVersion = parsedKubeVersion
			}
			if client.WithSubcharts {
				for _, p := range paths {
					filepath.Walk(filepath.Join(p, "charts"), func(path string, info os.FileInfo, _ error) error {
						if info != nil {
							if info.Name() == "Chart.yaml" {
								paths = append(paths, filepath.Dir(path))
							} else if strings.HasSuffix(path, ".tgz") || strings.HasSuffix(path, ".tar.gz") {
								paths = append(paths, path)
							}
						}
						return nil
					})
				}
			}
			client.Namespace = settings.Namespace()
			vals, err := valueOpts.MergeValues(getter.All(settings))
			if err != nil {
				return err
			}
			var ignorePatterns map[string][]string
			if lintIgnoreFile != "" {
				debug("\nUsing ignore file: %s\n", lintIgnoreFile)
				ignorePatterns, err = rules.ParseIgnoreFile(lintIgnoreFile)
			}
			var message strings.Builder
			failed := 0
			for _, path := range paths {
				fmt.Fprintf(&message, "==> Linting %s\n", path)
				result := client.Run([]string{path}, vals)
				result.Messages = rules.FilterIgnoredMessages(result.Messages, ignorePatterns)
				result.Errors = rules.FilterIgnoredErrors(result.Errors, ignorePatterns)
				for _, msg := range result.Messages {
					fmt.Fprintf(&message, "%s\n", msg)
				}
				if len(result.Messages) != 0 {
					failed++
					for _, err := range result.Errors {
						fmt.Fprintf(&message, "Error: %s\n", err)
					}
				}
			}

			fmt.Fprint(out, message.String())
			summary := fmt.Sprintf("%d chart(s) linted, %d chart(s) failed", len(paths), failed)
			if failed > 0 {
				return errors.New(summary)
			}
			fmt.Fprintln(out, summary)
			return nil
		},
	}

	f := cmd.Flags()
	f.BoolVar(&client.Strict, "strict", false, "fail on lint warnings")
	f.BoolVar(&client.WithSubcharts, "with-subcharts", false, "lint dependent charts")
	f.BoolVar(&client.Quiet, "quiet", false, "print only warnings and errors")
	f.StringVar(&kubeVersion, "kube-version", "", "Kubernetes version used for capabilities and deprecation checks")
	f.StringVar(&lintIgnoreFile, "lint-ignore-file", "", "path to .helmlintignore file to specify ignore patterns")

	return cmd
}
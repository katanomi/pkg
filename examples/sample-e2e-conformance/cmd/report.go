/*
Copyright 2023 The Katanomi Authors.

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
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/katanomi/pkg/examples/sample-e2e-conformance/report"
)

func main() {
	var output string

	var rootCmd = &cobra.Command{
		Use:   "report [report file]",
		Short: "parse ginkgo report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reportFile := args[0]
			fmt.Printf("Processing report file: %s\n", reportFile)
			if output != "" {
				fmt.Printf("Will save results to: %s\n", output)
			}
			nodes := report.Build(reportFile)
			data, _ := yaml.Marshal(nodes)
			return os.WriteFile(output, data, 0644)
		},
	}

	rootCmd.Flags().StringVarP(&output, "output", "o", "report.yaml", "Output file to save results")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

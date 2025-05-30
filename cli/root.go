// Copyright 2025 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	crconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/loaders/composite"
	"sigs.k8s.io/crdify/pkg/loaders/file"
	"sigs.k8s.io/crdify/pkg/loaders/git"
	"sigs.k8s.io/crdify/pkg/loaders/kubernetes"
	"sigs.k8s.io/crdify/pkg/loaders/scheme"
	"sigs.k8s.io/crdify/pkg/runner"
)

// NewRootCommand returns a cobra.Command for the program entrypoint.
func NewRootCommand() *cobra.Command {
	loader := composite.NewComposite(
		map[string]composite.Loader{
			scheme.SchemeKubernetes: kubernetes.New(crconfig.GetConfig),
			scheme.SchemeFile:       file.New(afero.OsFs{}),
			scheme.SchemeGit:        git.New(),
		},
	)

	var (
		configFile   string
		outputFormat string
	)

	rootCmd := &cobra.Command{
		Use:   "crdify <old> <new>",
		Short: "crdify evaluates changes to Kubernetes CustomResourceDefinitions",
		Long: `crdify is a tool for evaluating changes to Kubernetes CustomResourceDefinitions
to help cluster administrators, gitops practitioners, and Kubernetes extension developers identify
changes that might result in a negative impact to clusters and/or users.

Example use cases:
    Ealuating a change in a CustomResourceDefinition on a Kubernetes Cluster with one in a file:
        $ crdify kube://{crd-name} file://{filepath}

    Evaluating a change from file to file:
        $ crdify file://{filepath} file://{filepath}

    Evaluating a change from git ref to git ref:
            $ crdify git://{ref}?path={filepath} git://{ref}?path={filepath}`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(configFile)
			if err != nil {
				log.Fatalf("loading config: %v", err)
			}

			run, err := runner.New(cfg, runner.DefaultRegistry())
			if err != nil {
				log.Fatalf("configuring validation runner: %v", err)
			}

			oldCrd, err := loader.Load(cmd.Context(), args[0])
			if err != nil {
				log.Fatalf("loading old CustomResourceDefinition: %v", err)
			}

			newCrd, err := loader.Load(cmd.Context(), args[1])
			if err != nil {
				log.Fatalf("loading new CustomResourceDefinition: %v", err)
			}

			results := run.Run(oldCrd, newCrd)

			report, err := results.Render(runner.Format(outputFormat))
			if err != nil {
				// TODO: can we handle this better than spitting out an obtuse error?
				log.Fatalf("rendering run results: %v", err)
			}

			fmt.Print(report)
			if results.HasFailures() {
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "the filepath to load the check configurations from")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "plaintext", "the format the output should take when incompatibilities are identified. May be one of plaintext, markdown, json, yaml")

	return rootCmd
}

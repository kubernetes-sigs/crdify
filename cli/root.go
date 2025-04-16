package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/loaders/composite"
	"github.com/everettraven/crd-diff/pkg/loaders/file"
	"github.com/everettraven/crd-diff/pkg/loaders/git"
	"github.com/everettraven/crd-diff/pkg/loaders/kubernetes"
	"github.com/everettraven/crd-diff/pkg/loaders/scheme"
	"github.com/everettraven/crd-diff/pkg/runner"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	crconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func NewRootCommand() *cobra.Command {
	loader := composite.NewComposite(
		map[string]composite.Loader{
			scheme.SchemeKubernetes: kubernetes.New(crconfig.GetConfig),
			scheme.SchemeFile:       file.New(afero.OsFs{}),
			scheme.SchemeGit:        git.New(),
		},
	)

	var configFile string
	var outputFormat string

	rootCmd := &cobra.Command{
		Use:   "crd-diff <old> <new>",
		Short: "crd-diff evaluates changes to Kubernetes CustomResourceDefinitions",
		Long: `crd-diff is a tool for evaluating changes to Kubernetes CustomResourceDefinitions
to help cluster administrators, gitops practitioners, and Kubernetes extension developers identify
changes that might result in a negative impact to clusters and/or users.

Example use cases:
    Evaluating a change in a CustomResourceDefinition on a Kubernetes Cluster with one in a file:
        $ crd-diff kube://{crd-name} file://{filepath}

    Evaluating a change from file to file:
        $ crd-diff file://{filepath} file://{filepath}

    Evaluating a change from git ref to git ref:
            $ crd-diff git://{ref}?path={filepath} git://{ref}?path={filepath}`,
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

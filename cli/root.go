package cli

import (
	"io"
	"log"
	"net/url"
	"os"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/loaders/composite"
	"github.com/everettraven/crd-diff/pkg/loaders/file"
	"github.com/everettraven/crd-diff/pkg/loaders/git"
	"github.com/everettraven/crd-diff/pkg/loaders/kubernetes"
	"github.com/everettraven/crd-diff/pkg/loaders/scheme"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func NewRootCommand() *cobra.Command {
	loader := composite.NewComposite(
		composite.WithLoaders(map[string]composite.Loader{
			scheme.SchemeKubernetes: kubernetes.NewKubernetes(),
			scheme.SchemeFile:       file.NewFile(afero.OsFs{}),
			scheme.SchemeGit:        git.NewGit(),
		}),
	)

	var configFile string

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
			oldURL, err := url.Parse(args[0])
			if err != nil {
				log.Fatalf("parsing old source: %v", err)
			}

			newURL, err := url.Parse(args[1])
			if err != nil {
				log.Fatalf("parsing new source: %v", err)
			}

			oldCrd, err := loader.Load(cmd.Context(), *oldURL)
			if err != nil {
				log.Fatalf("loading old CustomResourceDefinition: %v", err)
			}

			newCrd, err := loader.Load(cmd.Context(), *newURL)
			if err != nil {
				log.Fatalf("loading new CustomResourceDefinition: %v", err)
			}

			cfg := &config.StrictConfig

			if configFile != "" {
				file, err := os.Open(configFile)
				if err != nil {
					log.Fatalf("loading config file %q: %v", configFile, err)
				}

				configBytes, err := io.ReadAll(file)
				if err != nil {
					log.Fatalf("reading config file %q: %v", configFile, err)
				}
				file.Close()

				err = yaml.Unmarshal(configBytes, cfg)
				if err != nil {
					log.Fatalf("unmarshalling config file %q contents: %v", configFile, err)
				}
			}

			validator := config.ValidatorForConfig(*cfg)

			err = validator.Validate(oldCrd, newCrd)
			if err != nil {
				log.Fatalf("comparing old and new CustomResourceDefinitions: %v", err)
			}
		},
	}

	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "the filepath to load the check configurations from")

	return rootCmd
}

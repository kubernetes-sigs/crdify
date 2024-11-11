package cli

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

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
	var outputFormat string

	const outputFormatJSON = "json"
	const outputFormatYAML = "yaml"
	const outputFormatPlainText = "plaintext"

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

			result := validator.Validate(oldCrd, newCrd)
			err = result.Error(0)
			if err != nil {
				switch outputFormat {
				case outputFormatPlainText:
					var out strings.Builder
					out.WriteString("comparing the CRDs identified incompatible changes\n\n")
					out.WriteString(err.Error())
					log.Fatal(out.String())
				case outputFormatJSON:
					jsonOut, marshalError := result.JSON()
					if marshalError != nil {
						log.Fatalf("marshalling results to JSON: %v", marshalError)
					}
					fmt.Print(string(jsonOut))
					os.Exit(1)
				case outputFormatYAML:
					yamlOut, marshalError := result.YAML()
					if marshalError != nil {
						log.Fatalf("marshalling results to YAML: %v", marshalError)
					}
					fmt.Print(string(yamlOut))
					os.Exit(1)
				}
			}
		},
	}

	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "the filepath to load the check configurations from")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "plaintext", "the format the output should take when incompatibilities are identified. May be one of plaintext, json, yaml")

	return rootCmd
}

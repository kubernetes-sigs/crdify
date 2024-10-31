package cli

import (
	"log"
	"net/url"

	"github.com/everettraven/crd-diff/pkg/loaders/composite"
	"github.com/everettraven/crd-diff/pkg/loaders/file"
	"github.com/everettraven/crd-diff/pkg/loaders/git"
	"github.com/everettraven/crd-diff/pkg/loaders/kubernetes"
	"github.com/everettraven/crd-diff/pkg/loaders/scheme"
	"github.com/everettraven/crd-diff/pkg/validations/property"
	"github.com/everettraven/crd-diff/pkg/validations/validators/crd"
	"github.com/everettraven/crd-diff/pkg/validations/validators/version"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	loader := composite.NewComposite(
		composite.WithLoaders(map[string]composite.Loader{
			scheme.SchemeKubernetes: kubernetes.NewKubernetes(),
			scheme.SchemeFile:       file.NewFile(afero.OsFs{}),
			scheme.SchemeGit:        git.NewGit(),
		}),
	)

	propertyValidations := []property.PropertyValidation{
		&property.Enum{},
		&property.Default{},
		&property.Required{},
		&property.Type{},
		&property.Maximum{},
		&property.MaxItems{},
		&property.MaxLength{},
		&property.MaxProperties{},
		&property.Minimum{},
		&property.MinItems{},
		&property.MinLength{},
		&property.MinProperties{},
	}

	// TODO: load this dynamically based on a configuration
	validator := crd.NewValidator(
		crd.WithValidations(
			version.NewValidator(
				version.WithSameVersionConfig(version.SameVersionConfig{
					UnhandledFailureMode: version.FailureModeClosed,
					Skip:                 false,
					Validations:          propertyValidations,
				}),
				version.WithServedVersionConfig(version.ServedVersionConfig{
					UnhandledFailureMode: version.FailureModeClosed,
					Skip:                 false,
					IgnoreConversion:     false,
					Validations:          propertyValidations,
				}),
			),
			&crd.Scope{},
			&crd.ExistingFieldRemoval{},
			&crd.StoredVersionRemoval{},
		),
	)

	rootCmd := &cobra.Command{
		Use:   "crd-diff <old> <new>",
		Short: "crd-diff evaluates changes to Kubernetes CustomResourceDefinitions",
		Args:  cobra.ExactArgs(2),
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

			err = validator.Validate(oldCrd, newCrd)
			if err != nil {
				log.Fatalf("comparing old and new CustomResourceDefinitions: %v", err)
			}
		},
	}
	rootCmd.AddCommand(NewVersionCommand())

	return rootCmd
}

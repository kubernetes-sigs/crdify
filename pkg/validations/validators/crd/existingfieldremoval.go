package crd

import (
	"errors"
	"strings"

	"github.com/openshift/crd-schema-checker/pkg/manifestcomparators"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type ExistingFieldRemoval struct{}

func (efr *ExistingFieldRemoval) Name() string {
	return "ExistingFieldRemoval"
}

func (efr *ExistingFieldRemoval) Validate(old, new *apiextensionsv1.CustomResourceDefinition) error {
	reg := manifestcomparators.NewRegistry()
	err := reg.AddComparator(manifestcomparators.NoFieldRemoval())
	if err != nil {
		return err
	}

	results, errs := reg.Compare(old, new)
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	errSet := []error{}

	for _, result := range results {
		if len(result.Errors) > 0 {
			errSet = append(errSet, errors.New(strings.Join(result.Errors, "\n")))
		}
	}
	if len(errSet) > 0 {
		return errors.Join(errSet...)
	}

	return nil
}

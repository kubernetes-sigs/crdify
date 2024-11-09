package crd

import (
	"errors"
	"strings"

	"github.com/openshift/crd-schema-checker/pkg/manifestcomparators"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type ExistingFieldRemoval struct{}

func (efr *ExistingFieldRemoval) Name() string {
	return "existingFieldRemoval"
}

func (efr *ExistingFieldRemoval) Validate(old, new *apiextensionsv1.CustomResourceDefinition) ValidationResult {
	reg := manifestcomparators.NewRegistry()
	err := reg.AddComparator(manifestcomparators.NoFieldRemoval())
	if err != nil {
		return &validationResult{
			Validation: efr.Name(),
			Err:        err.Error(),
		}
	}

	results, errs := reg.Compare(old, new)
	if len(errs) > 0 {
		return &validationResult{
			Validation: efr.Name(),
			Err:        errors.Join(errs...).Error(),
		}
	}

	vr := &validationResult{
		Validation: efr.Name(),
	}

	errSet := []error{}
	for _, result := range results {
		if len(result.Errors) > 0 {
			errSet = append(errSet, errors.New(strings.Join(result.Errors, "\n")))
		}
	}

	if errors.Join(errSet...) != nil {
		vr.Err = errors.Join(errSet...).Error()
	}

	return vr
}

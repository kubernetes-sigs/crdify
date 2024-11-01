package crd

import (
	"errors"
	"fmt"

	"github.com/everettraven/crd-diff/pkg/validations/results"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Validation is a representation of a validation that is run
// against two revisions of a CustomResourceDefinition
type Validation interface {
	// Validate performs the validation, returning an error if the
	// new revision is incompatible with the old revision of the CustomResourceDefinition
	Validate(old, new *apiextensionsv1.CustomResourceDefinition) *results.Result

	// Name is a human-readable name for this validation
	Name() string
}

// Validator validates Kubernetes CustomResourceDefinitions using the configured validations
type Validator struct {
	validations []Validation
}

type ValidatorOption func(*Validator)

func WithValidations(validations ...Validation) ValidatorOption {
	return func(v *Validator) {
		v.validations = validations
	}
}

func NewValidator(opts ...ValidatorOption) *Validator {
	validator := &Validator{
		validations: []Validation{},
	}

	for _, opt := range opts {
		opt(validator)
	}
	return validator
}

// Validate runs the validations configured in the Validator
func (v *Validator) Validate(old, new *apiextensionsv1.CustomResourceDefinition) error {
	result := &results.Result{
		Subresults: []*results.Result{},
	}
	for _, validation := range v.validations {
		if res := validation.Validate(old, new); res != nil {
			subResult := &results.Result{
				Subresults: []*results.Result{res},
			}
			if res.Error != nil {
				subResult.Error = fmt.Errorf("%q validation failed", validation.Name())
				result.Error = errors.New("potentially breaking changes found")
			}
			result.Subresults = append(result.Subresults, subResult)
		}
	}
	return results.ErrorFromResult(result, 0)
}

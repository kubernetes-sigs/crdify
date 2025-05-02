package crd

import (
	"github.com/everettraven/crd-diff/pkg/validations"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Validator validates Kubernetes CustomResourceDefinitions using the configured validations
type Validator struct {
	comparators []validations.Comparator[apiextensionsv1.CustomResourceDefinition]
}

// ValidatorOption configures a Validator
type ValidatorOption func(*Validator)

// WithComparators configures a Validator with the provided CustomResourceDefinition Comparators.
// Each call to WithComparators is a replacement, not additive.
func WithComparators(comparators ...validations.Comparator[apiextensionsv1.CustomResourceDefinition]) ValidatorOption {
	return func(v *Validator) {
		v.comparators = comparators
	}
}

// New returns a new Validator for validating an old and new CustomResourceDefinition
// configured with the provided ValidatorOptions
func New(opts ...ValidatorOption) *Validator {
	validator := &Validator{
		comparators: []validations.Comparator[apiextensionsv1.CustomResourceDefinition]{},
	}

	for _, opt := range opts {
		opt(validator)
	}
	return validator
}

// Validate runs the validations configured in the Validator
func (v *Validator) Validate(old, new *apiextensionsv1.CustomResourceDefinition) []validations.ComparisonResult {
	result := []validations.ComparisonResult{}
	for _, comparator := range v.comparators {
		compResult := comparator.Compare(old, new)
		result = append(result, compResult)
	}
	return result
}

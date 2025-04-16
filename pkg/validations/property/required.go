package property

import (
	"fmt"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/validations"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	_ validations.Validation                                  = (*Required)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*Required)(nil)
)

const requiredValidationName = "required"

// RegisterRequired registers the Required validation
// with the provided validation registry
func RegisterRequired(registry validations.Registry) {
	registry.Register(requiredValidationName, requiredFactory)
}

// requiredFactory is a function used to initialize a Required validation
// implementation based on the provided configuration.
func requiredFactory(_ map[string]interface{}) (validations.Validation, error) {
	return &Required{}, nil
}

// Required is a Validation that can be used to identify
// incompatible changes to the required constraints of CRD properties
type Required struct {
	enforcement config.EnforcementPolicy
}

// Name returns the name of the Required validation
func (r *Required) Name() string {
	return requiredValidationName
}

// SetEnforcement sets the EnforcementPolicy for the Required validation
func (r *Required) SetEnforcement(policy config.EnforcementPolicy) {
	r.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the required constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.Required field will be reset to 'nil' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (r *Required) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	oldRequired := sets.New(a.Required...)
	newRequired := sets.New(b.Required...)
	diffRequired := newRequired.Difference(oldRequired)
	var err error

	if diffRequired.Len() > 0 {
		err = fmt.Errorf("new required fields %v added", diffRequired.UnsortedList())
	}

	a.Required = nil
	b.Required = nil

	return validations.HandleErrors(r.Name(), r.enforcement, err)
}

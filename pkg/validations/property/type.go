package property

import (
	"fmt"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/validations"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	_ validations.Validation                                  = (*Type)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*Type)(nil)
)

const typeValidationName = "type"

// RegisterType registers the Type validation
// with the provided validation registry
func RegisterType(registry validations.Registry) {
	registry.Register(typeValidationName, typeFactory)
}

// typeFactory is a function used to initialize a Type validation
// implementation based on the provided configuration.
func typeFactory(_ map[string]interface{}) (validations.Validation, error) {
	return &Type{}, nil
}

// Type is a Validation that can be used to identify
// incompatible changes to the type constraints of CRD properties
type Type struct {
	enforcement config.EnforcementPolicy
}

// Name returns the name of the Type validation
func (t *Type) Name() string {
	return typeValidationName
}

// SetEnforcement sets the EnforcementPolicy for the Type validation
func (t *Type) SetEnforcement(policy config.EnforcementPolicy) {
	t.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the type constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.Type field will be reset to '""' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (t *Type) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	var err error
	if a.Type != b.Type {
		err = fmt.Errorf("type changed from %q to %q", a.Type, b.Type)
	}

	a.Type = ""
	b.Type = ""

	return validations.HandleErrors(t.Name(), t.enforcement, err)
}

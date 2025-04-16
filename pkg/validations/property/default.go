package property

import (
	"bytes"
	"fmt"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/validations"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const defaultValidationName = "default"

var (
	_ validations.Validation                                  = (*Default)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*Default)(nil)
)

// RegisterDefault registers the Default validation
// with the provided validation registry
func RegisterDefault(registry validations.Registry) {
	registry.Register(defaultValidationName, defaultFactory)
}

// defaultFactory is a function used to initialize a Default validation
// implementation based on the provided configuration.
func defaultFactory(_ map[string]interface{}) (validations.Validation, error) {
	return &Default{}, nil
}

// Default is a Validation that can be used to identify
// incompatible changes to the default value of CRD properties
type Default struct {
	// enforcement is the EnforcementPolicy that this validation
	// should use when performing its validation logic
	enforcement config.EnforcementPolicy
}

// Name returns the name of the Default validation
func (d *Default) Name() string {
	return defaultValidationName
}

// SetEnforcement sets the EnforcementPolicy for the Default validation
func (d *Default) SetEnforcement(policy config.EnforcementPolicy) {
	d.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for changes to the default value of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.Default field will be reset to 'nil' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (d *Default) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	var err error

	switch {
	case a.Default == nil && b.Default != nil:
		err = fmt.Errorf("default value %q added when there was no default previously", string(b.Default.Raw))
	case a.Default != nil && b.Default == nil:
		err = fmt.Errorf("default value %q removed", string(a.Default.Raw))
	case a.Default != nil && b.Default != nil && !bytes.Equal(a.Default.Raw, b.Default.Raw):
		err = fmt.Errorf("default value changed from %q to %q", string(a.Default.Raw), string(b.Default.Raw))
	}

	// reset values
	a.Default = nil
	b.Default = nil

	return validations.HandleErrors(d.Name(), d.enforcement, err)
}

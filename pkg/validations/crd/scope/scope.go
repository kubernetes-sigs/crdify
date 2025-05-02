package scope

import (
	"fmt"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/validations"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	_ validations.Validation                                           = (*Scope)(nil)
	_ validations.Comparator[apiextensionsv1.CustomResourceDefinition] = (*Scope)(nil)
)

const name = "scope"

// Register registers the Scope validation
// with the provided validation registry
func Register(registry validations.Registry) {
	registry.Register(name, factory)
}

// factory is a function used to initialize a Scope validation
// implementation based on the provided configuration.
func factory(_ map[string]interface{}) (validations.Validation, error) {
	return &Scope{}, nil
}

// Scope is a validations.Validation implementation
// used to check if the scope has changed from one
// CRD instance to another
type Scope struct {
	// enforcement is the EnforcementPolicy that this validation
	// should use when performing its validation logic
	enforcement config.EnforcementPolicy
}

// Name returns the name of the Scope validation
func (s *Scope) Name() string {
	return name
}

// SetEnforcement sets the EnforcementPolicy for the Scope validation
func (s *Scope) SetEnforcement(enforcement config.EnforcementPolicy) {
	s.enforcement = enforcement
}

// Compare compares an old and a new CustomResourceDefintion, checking for any change to the scope from the
// old CustomResourceDefinition to the new CustomResourceDefinition
func (s *Scope) Compare(old, new *apiextensionsv1.CustomResourceDefinition) validations.ComparisonResult {
	var err error
	if old.Spec.Scope != new.Spec.Scope {
		err = fmt.Errorf("scope changed from %q to %q", old.Spec.Scope, new.Spec.Scope)
	}
	return validations.HandleErrors(s.Name(), s.enforcement, err)
}

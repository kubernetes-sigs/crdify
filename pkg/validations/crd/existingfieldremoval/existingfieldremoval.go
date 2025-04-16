package existingfieldremoval

import (
	"fmt"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/validations"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	_ validations.Validation                                           = (*ExistingFieldRemoval)(nil)
	_ validations.Comparator[apiextensionsv1.CustomResourceDefinition] = (*ExistingFieldRemoval)(nil)
)

const name = "existingFieldRemoval"

// Register registers the ExistingFieldRemoval validation
// with the provided validation registry
func Register(registry validations.Registry) {
	registry.Register(name, factory)
}

// factory is a function used to initialize an ExistingFieldRemoval validation
// implementation based on the provided configuration.
func factory(_ map[string]interface{}) (validations.Validation, error) {
	return &ExistingFieldRemoval{}, nil
}

// ExistingFieldRemoval is a validations.Validation implementation
// used to check if any existing fields have been removed from one
// CRD instance to another
type ExistingFieldRemoval struct {
	// enforcement is the EnforcementPolicy that this validation
	// should use when performing its validation logic
	enforcement config.EnforcementPolicy
}

// Name returns the name of the ExistingFieldRemoval validation
func (efr *ExistingFieldRemoval) Name() string {
	return name
}

// SetEnforcement sets the EnforcementPolicy for the ExistingFieldRemoval validation
func (efr *ExistingFieldRemoval) SetEnforcement(policy config.EnforcementPolicy) {
	efr.enforcement = policy
}

// Compare compares an old and a new CustomResourceDefintion, checking for any fields that were removed
// from the old CustomResourceDefinition in the new CustomResourceDefinition
func (efr *ExistingFieldRemoval) Compare(old, new *apiextensionsv1.CustomResourceDefinition) validations.ComparisonResult {
	errs := []error{}
	for _, newVersion := range new.Spec.Versions {
		existingVersion := validations.GetCRDVersionByName(old, newVersion.Name)
		if existingVersion == nil {
			continue
		}

		existingFields := getFields(existingVersion)
		newFields := getFields(&newVersion)

		removedFields := existingFields.Difference(newFields)
		for _, removedField := range removedFields.UnsortedList() {
			errs = append(errs, fmt.Errorf("crd/%v version/%v field/%v may not be removed", new.Name, newVersion.Name, removedField))
		}
	}

	return validations.HandleErrors(efr.Name(), efr.enforcement, errs...)
}

// getFields returns a set of all the fields for the provided CustomResourceDefinitionVersion
func getFields(v *apiextensionsv1.CustomResourceDefinitionVersion) sets.Set[string] {
	fields := sets.New[string]()
	validations.SchemaHas(v.Schema.OpenAPIV3Schema, field.NewPath("^"), field.NewPath("^"), nil,
		func(s *apiextensionsv1.JSONSchemaProps, fldPath, simpleLocation *field.Path, _ []*apiextensionsv1.JSONSchemaProps) bool {
			fields.Insert(simpleLocation.String())
			return false
		})

	return fields
}

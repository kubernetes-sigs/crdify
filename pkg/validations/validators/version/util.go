package version

import (
	"github.com/everettraven/crd-diff/pkg/validations/property"
	"github.com/openshift/crd-schema-checker/pkg/manifestcomparators"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func GetCRDVersionByName(crd *apiextensionsv1.CustomResourceDefinition, name string) *apiextensionsv1.CustomResourceDefinitionVersion {
	if crd == nil {
		return nil
	}

	for _, version := range crd.Spec.Versions {
		if version.Name == name {
			return &version
		}
	}

	return nil
}

func FlattenCRDVersion(crdVersion apiextensionsv1.CustomResourceDefinitionVersion) map[string]*apiextensionsv1.JSONSchemaProps {
	flatMap := map[string]*apiextensionsv1.JSONSchemaProps{}

	manifestcomparators.SchemaHas(crdVersion.Schema.OpenAPIV3Schema,
		field.NewPath("^"),
		field.NewPath("^"),
		nil,
		func(s *apiextensionsv1.JSONSchemaProps, _, simpleLocation *field.Path, _ []*apiextensionsv1.JSONSchemaProps) bool {
			flatMap[simpleLocation.String()] = s.DeepCopy()
			return false
		},
	)
	return flatMap
}

func FlattenedCRDVersionDiff(old, new map[string]*apiextensionsv1.JSONSchemaProps) map[string]property.Diff {
	diffMap := map[string]property.Diff{}
	for prop, oldSchema := range old {
		// Create a copy of the old schema and set the properties to nil.
		// In theory this will make it so we don't provide a diff for a parent property
		// based on changes to the children properties. The changes to the children
		// properties should still be evaluated since we are looping through a flattened
		// map of all the properties for the CRD version
		oldSchemaCopy := oldSchema.DeepCopy()
		oldSchemaCopy.Properties = nil
		newSchema, ok := new[prop]

		// In the event the property no longer exists on the new version
		// create a diff entry with the new value being empty
		if !ok {
			diffMap[prop] = property.NewDiff(oldSchemaCopy, &apiextensionsv1.JSONSchemaProps{})
		}

		// Do the same copy and unset logic for the new schema properties
		// before comparison to ensure we are only comparing the individual properties
		newSchemaCopy := newSchema.DeepCopy()
		newSchemaCopy.Properties = nil

		if !equality.Semantic.DeepEqual(oldSchemaCopy, newSchemaCopy) {
			diffMap[prop] = property.NewDiff(oldSchemaCopy, newSchemaCopy)
		}
	}

	return diffMap
}

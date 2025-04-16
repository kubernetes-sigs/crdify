package validations

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestFlattenedCRDVersionDiff(t *testing.T) {
	type testcase struct {
		name    string
		old     apiextensionsv1.CustomResourceDefinitionVersion
		new     apiextensionsv1.CustomResourceDefinitionVersion
		diffKey string
		oldDiff apiextensionsv1.JSONSchemaProps
		newDiff apiextensionsv1.JSONSchemaProps
	}

	for _, tc := range []testcase{
		{
			name: "removed field",
			old: apiextensionsv1.CustomResourceDefinitionVersion{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"foo": {
								Type: "string",
							},
							"bar": {
								Type: "string",
							},
						},
					},
				},
			},
			new: apiextensionsv1.CustomResourceDefinitionVersion{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"foo": {
								Type: "string",
							},
							// Removed "bar".
						},
					},
				},
			},
			// The field that has changed is <Root>.bar
			diffKey: "^.bar",
			// Removed field type should be included in the old.
			oldDiff: apiextensionsv1.JSONSchemaProps{
				Type: "string",
			},
			// The new schema in the diff should be equal to a zero schema because the field was removed.
			newDiff: apiextensionsv1.JSONSchemaProps{},
		},
		{
			name: "added required field",
			old: apiextensionsv1.CustomResourceDefinitionVersion{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"foo": {
								Type: "string",
							},
						},
					},
				},
			},
			new: apiextensionsv1.CustomResourceDefinitionVersion{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"foo": {
								Type: "string",
							},
							// Added field "bar".
							"bar": {
								Type: "string",
							},
						},
						Required: []string{"bar"},
					},
				},
			},
			// Added a new field at root (^).
			diffKey: "^",
			// Field added in new schema required.
			oldDiff: apiextensionsv1.JSONSchemaProps{},
			newDiff: apiextensionsv1.JSONSchemaProps{
				Required: []string{"bar"},
			},
		},
		{
			name: "no change",
			old: apiextensionsv1.CustomResourceDefinitionVersion{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"foo": {
								Type: "string",
							},
							"bar": {
								Type: "string",
							},
						},
					},
				},
			},
			new: apiextensionsv1.CustomResourceDefinitionVersion{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"foo": {
								Type: "string",
							},
							"bar": {
								Type: "string",
							},
						},
					},
				},
			},
			// No fields have changed, so no diff.
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oldFlattened := FlattenCRDVersion(tc.old)
			newFlattened := FlattenCRDVersion(tc.new)

			diffs := FlattenedCRDVersionDiff(oldFlattened, newFlattened)

			if tc.diffKey != "" {
				// If a diff is expected, so the result should be validated.
				require.True(t, reflect.DeepEqual(*diffs[tc.diffKey].Old, tc.oldDiff))
				require.True(t, reflect.DeepEqual(*diffs[tc.diffKey].New, tc.newDiff))
			} else {
				// If a diff is not expected, so there should not be one.
				require.Len(t, diffs, 0)
			}
		})
	}
}

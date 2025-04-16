package existingfieldremoval

import (
	"testing"

	internaltesting "github.com/everettraven/crd-diff/pkg/validations/internal/testing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestExistingFieldRemoval(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.CustomResourceDefinition]{
		{
			Name: "no existing field removed, not flagged",
			Old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &ExistingFieldRemoval{},
		},
		{
			Name: "existing field removed, flagged",
			Old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
										"fieldTwo": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			Flagged:              true,
			ComparableValidation: &ExistingFieldRemoval{},
		},
		{
			Name: "new version is added with the field removed, not flagged",
			Old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
										"fieldTwo": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
										"fieldTwo": {
											Type: "string",
										},
									},
								},
							},
						},
						{
							Name: "v1alpha2",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &ExistingFieldRemoval{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

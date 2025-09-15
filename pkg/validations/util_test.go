// Copyright 2025 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validations

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
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

func TestFlattenCRDVersion(t *testing.T) {
	type testcase struct {
		name         string
		version      apiextensionsv1.CustomResourceDefinitionVersion
		expectedKeys []string
	}

	testcases := []testcase{
		{
			name: "basic schema",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "string",
									},
									"fieldTwo": {
										Type: "string",
									},
									"fieldThree": {
										Type: "object",
										Properties: map[string]apiextensionsv1.JSONSchemaProps{
											"subfield": {
												Type: "number",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldTwo",
				"^.spec.fieldThree",
				"^.spec.fieldThree.subfield",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with items",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "array",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											Schema: &apiextensionsv1.JSONSchemaProps{
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"subfield": {
														Type: "number",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.items.subfield",
				"^.spec.fieldOne.items",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with items with JSONSchemas",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "array",
										Items: &apiextensionsv1.JSONSchemaPropsOrArray{
											JSONSchemas: []apiextensionsv1.JSONSchemaProps{
												{
													Properties: map[string]apiextensionsv1.JSONSchemaProps{
														"subfield": {
															Type: "number",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.items[0].subfield",
				"^.spec.fieldOne.items[0]",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with allOf",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										AllOf: []apiextensionsv1.JSONSchemaProps{
											{
												Type: "object",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"nested": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.allOf[0]",
				"^.spec.fieldOne.allOf[0].nested",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with anyOf",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										AnyOf: []apiextensionsv1.JSONSchemaProps{
											{
												Type: "object",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"nested": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.anyOf[0]",
				"^.spec.fieldOne.anyOf[0].nested",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with oneOf",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										OneOf: []apiextensionsv1.JSONSchemaProps{
											{
												Type: "object",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"nested": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.oneOf[0]",
				"^.spec.fieldOne.oneOf[0].nested",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with not",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										Not: &apiextensionsv1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]apiextensionsv1.JSONSchemaProps{
												"nested": {
													Type: "string",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.not",
				"^.spec.fieldOne.not.nested",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with additionalProperties",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										AdditionalProperties: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{
												Type: "object",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"nested": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.additionalProperties",
				"^.spec.fieldOne.additionalProperties.nested",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with patternProperties",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										PatternProperties: map[string]apiextensionsv1.JSONSchemaProps{
											"pattern": {
												Type: "object",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"nested": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.patternProperties[pattern]",
				"^.spec.fieldOne.patternProperties[pattern].nested",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with additionalItems",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										AdditionalItems: &apiextensionsv1.JSONSchemaPropsOrBool{
											Allows: true,
											Schema: &apiextensionsv1.JSONSchemaProps{
												Type: "object",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"nested": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.additionalItems",
				"^.spec.fieldOne.additionalItems.nested",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with definitions",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										Definitions: apiextensionsv1.JSONSchemaDefinitions{
											"thing": apiextensionsv1.JSONSchemaProps{
												Type: "object",
												Properties: map[string]apiextensionsv1.JSONSchemaProps{
													"nested": {
														Type: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.definitions[thing]",
				"^.spec.fieldOne.definitions[thing].nested",
				"^.spec",
				"^",
			},
		},
		{
			name: "schema with dependencies",
			version: apiextensionsv1.CustomResourceDefinitionVersion{
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"spec": {
								Properties: map[string]apiextensionsv1.JSONSchemaProps{
									"fieldOne": {
										Type: "object",
										Dependencies: apiextensionsv1.JSONSchemaDependencies{
											"dependencyOne": {
												Schema: &apiextensionsv1.JSONSchemaProps{
													Type: "object",
													Properties: map[string]apiextensionsv1.JSONSchemaProps{
														"nested": {
															Type: "string",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedKeys: []string{
				"^.spec.fieldOne",
				"^.spec.fieldOne.dependencies[dependencyOne].nested",
				"^.spec.fieldOne.dependencies[dependencyOne]",
				"^.spec",
				"^",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := FlattenCRDVersion(tc.version)

			actualSet := sets.New[string]()
			for key := range out {
				actualSet.Insert(key)
			}

			expectedSet := sets.New(tc.expectedKeys...)

			if !expectedSet.Equal(actualSet) {
				t.Fatalf("expectedKeys does not match actual keys - in actual but not expected: %v - in expected but not actual: %v", actualSet.Difference(expectedSet).UnsortedList(), expectedSet.Difference(actualSet).UnsortedList())
			}
		})
	}
}

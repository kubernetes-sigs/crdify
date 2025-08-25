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

package served

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
	"sigs.k8s.io/crdify/pkg/validations/property"
)

type createCRDOption func(crd *apiextensionsv1.CustomResourceDefinition)

func createCRD(options ...createCRDOption) *apiextensionsv1.CustomResourceDefinition {
	crd := &apiextensionsv1.CustomResourceDefinition{}
	for _, option := range options {
		option(crd)
	}

	return crd
}

func withVersion(name string, served bool, schema *apiextensionsv1.JSONSchemaProps) func(crd *apiextensionsv1.CustomResourceDefinition) {
	return func(crd *apiextensionsv1.CustomResourceDefinition) {
		crd.Spec.Versions = append(crd.Spec.Versions, apiextensionsv1.CustomResourceDefinitionVersion{
			Name:   name,
			Served: served,
			Schema: &apiextensionsv1.CustomResourceValidation{
				OpenAPIV3Schema: schema,
			},
		})
	}
}

func withConversion(strategy apiextensionsv1.ConversionStrategyType) func(crd *apiextensionsv1.CustomResourceDefinition) {
	return func(crd *apiextensionsv1.CustomResourceDefinition) {
		crd.Spec.Conversion = &apiextensionsv1.CustomResourceConversion{
			Strategy: strategy,
		}
	}
}

func createSimpleSchema(properties map[string]apiextensionsv1.JSONSchemaProps) *apiextensionsv1.JSONSchemaProps {
	return &apiextensionsv1.JSONSchemaProps{
		Type:       "object",
		Properties: properties,
	}
}

func TestValidator_Validate_ConversionPolicyIgnore(t *testing.T) {
	tests := []struct {
		name        string
		crdA        *apiextensionsv1.CustomResourceDefinition
		crdB        *apiextensionsv1.CustomResourceDefinition
		expectEmpty bool
	}{
		{
			name: "webhook conversion with ignore policy should return empty results",
			crdA: createCRD(
				withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
				withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
			),
			crdB: createCRD(
				withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
				withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "integer"}, // Changed type to trigger comparison
				})),
				withConversion(apiextensionsv1.WebhookConverter),
			),
			expectEmpty: true,
		},
		{
			name: "no conversion webhook with ignore policy should still validate",
			crdA: createCRD(
				withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
				withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
			),
			crdB: createCRD(
				withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
				withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "integer"}, // Changed type to trigger comparison
				})),
			),
			expectEmpty: false,
		},
		{
			name: "non-webhook conversion with ignore policy should still validate",
			crdA: createCRD(
				withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
				withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
			),
			crdB: createCRD(
				withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
				withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "integer"}, // Changed type to trigger comparison
				})),
				withConversion(apiextensionsv1.NoneConverter),
			),
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := New(
				WithConversionPolicy(config.ConversionPolicyIgnore),
				WithUnhandledEnforcementPolicy(config.EnforcementPolicyError),
			)

			result := validator.Validate(tt.crdA, tt.crdB)

			if tt.expectEmpty {
				assert.Empty(t, result, "Expected empty results for webhook conversion with ignore policy")
			} else {
				assert.NotEmpty(t, result, "Expected non-empty results when validation should run")
			}
		})
	}
}

func TestValidator_Validate_ConversionPolicyNone(t *testing.T) {
	validator := New(
		WithConversionPolicy(config.ConversionPolicyNone),
		WithUnhandledEnforcementPolicy(config.EnforcementPolicyError),
	)

	crdA := createCRD(
		withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "string"},
		})),
		withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "string"},
		})),
	)

	crdB := createCRD(
		withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "string"},
		})),
		withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "integer"}, // Changed type to trigger comparison
		})),
		withConversion(apiextensionsv1.WebhookConverter),
	)

	result := validator.Validate(crdA, crdB)

	// Even with webhook conversion, should still validate when policy is None
	assert.NotEmpty(t, result, "Expected validation to run even with webhook conversion when policy is None")
}

func TestValidator_Validate_VersionPairs(t *testing.T) {
	tests := []struct {
		name          string
		crd           *apiextensionsv1.CustomResourceDefinition
		expectedPairs []string
	}{
		{
			name: "no served versions",
			crd: createCRD(
				withVersion("v1", false, createSimpleSchema(nil)),
			),
			expectedPairs: []string{},
		},
		{
			name: "single served version",
			crd: createCRD(
				withVersion("v1", true, createSimpleSchema(nil)),
			),
			expectedPairs: []string{},
		},
		{
			name: "two served versions",
			crd: createCRD(
				withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
				})),
				withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
					"field1": {Type: "string"},
					"field2": {Type: "integer"},
				})),
			),
			expectedPairs: []string{"v1 -> v2"},
		},
		{
			name: "three served versions",
			crd: createCRD(
				withVersion("v1", true, createSimpleSchema(nil)),
				withVersion("v2", true, createSimpleSchema(nil)),
				withVersion("v3", true, createSimpleSchema(nil)),
			),
			expectedPairs: []string{"v1 -> v2", "v1 -> v3", "v2 -> v3"},
		},
		{
			name: "mixed served and non-served versions",
			crd: createCRD(
				withVersion("v1", true, createSimpleSchema(nil)),
				withVersion("v1beta1", false, createSimpleSchema(nil)), // not served
				withVersion("v2", true, createSimpleSchema(nil)),
			),
			expectedPairs: []string{"v1 -> v2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := New(WithUnhandledEnforcementPolicy(config.EnforcementPolicyError))

			// Pair generation happens in the context of a single CRD, so we can
			// compare the same CRD to itself to get the version pairs for the purpose
			// of this test.
			result := validator.Validate(tt.crd, tt.crd)

			// Extract pair keys from result and compare
			actualPairs := make([]string, 0, len(result))
			for pair := range result {
				actualPairs = append(actualPairs, pair)
			}

			assert.ElementsMatch(t, tt.expectedPairs, actualPairs, "Version pairs don't match")
		})
	}
}

func TestValidator_Validate_EnforcementPolicies(t *testing.T) {
	tests := []struct {
		name           string
		policy         config.EnforcementPolicy
		expectResultFn func(t *testing.T, result validations.ComparisonResult)
	}{
		{
			name:   "error policy",
			policy: config.EnforcementPolicyError,
			expectResultFn: func(t *testing.T, result validations.ComparisonResult) {
				assert.Len(t, result.Errors, 1)
				assert.Empty(t, result.Warnings)
			},
		},
		{
			name:   "warn policy",
			policy: config.EnforcementPolicyWarn,
			expectResultFn: func(t *testing.T, result validations.ComparisonResult) {
				assert.Empty(t, result.Errors)
				assert.Len(t, result.Warnings, 1)
			},
		},
		{
			name:   "none policy",
			policy: config.EnforcementPolicyNone,
			expectResultFn: func(t *testing.T, result validations.ComparisonResult) {
				assert.Empty(t, result.Errors)
				assert.Empty(t, result.Warnings)
			},
		},
	}

	crdA := createCRD(
		withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "string"},
		})),
		withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "string"}, // No change within A
		})),
	)

	crdB := createCRD(
		withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "string"},
		})),
		withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "integer"}, // Change only in B, not in A
		})),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := New(WithUnhandledEnforcementPolicy(tt.policy))

			result := validator.Validate(crdA, crdB)
			for _, versionPair := range result {
				for _, fieldResults := range versionPair {
					for _, compResult := range fieldResults {
						tt.expectResultFn(t, compResult)
					}
				}
			}
		})
	}
}

func TestValidator_Validate_SubtractExistingIssues(t *testing.T) {
	// Use type comparator for errors and description comparator for warnings
	typeComparator := &property.Type{}
	typeComparator.SetEnforcement(config.EnforcementPolicyError)

	descriptionComparator := &property.Description{}
	descriptionComparator.SetEnforcement(config.EnforcementPolicyWarn)

	validator := New(
		WithComparators(typeComparator, descriptionComparator),
		WithUnhandledEnforcementPolicy(config.EnforcementPolicyError),
	)

	// CRD A has some issues that will be common with B (and thus subtracted)
	crdA := createCRD(
		withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "string", Description: "Original description"},
		})),
		withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "integer", Description: "Updated description"}, // type change + description change
		})),
	)

	// CRD B has the same common issues PLUS new issues that should be reported
	crdB := createCRD(
		withVersion("v1", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "string", Description: "Original description"},
			"field2": {Type: "string", Description: "New field original"}, // New field in v1
		})),
		withVersion("v2", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "integer", Description: "Updated description"}, // Same changes as A (will be subtracted)
			"field2": {Type: "integer", Description: "New field changed"},   // New field with different type/description
		})),
		// New version v3 that doesn't exist in A
		withVersion("v3", true, createSimpleSchema(map[string]apiextensionsv1.JSONSchemaProps{
			"field1": {Type: "boolean", Description: "Different type and description"}, // New version issues
		})),
	)

	result := validator.Validate(crdA, crdB)

	// Define expected structure with counts instead of exact messages
	type expectedComparison struct {
		ExpectedErrors   int
		ExpectedWarnings int
	}

	expected := map[string]map[string]map[string]expectedComparison{
		"v1 -> v2": {
			"^.field1": {
				"type":        {ExpectedErrors: 0, ExpectedWarnings: 0}, // Common issue, subtracted
				"description": {ExpectedErrors: 0, ExpectedWarnings: 0}, // Common issue, subtracted
				"unhandled":   {ExpectedErrors: 0, ExpectedWarnings: 0},
			},
			"^.field2": {
				"type":        {ExpectedErrors: 1, ExpectedWarnings: 0}, // New field issue
				"description": {ExpectedErrors: 0, ExpectedWarnings: 1}, // New field issue
				"unhandled":   {ExpectedErrors: 0, ExpectedWarnings: 0},
			},
		},
		"v1 -> v3": {
			"^.field1": {
				"type":        {ExpectedErrors: 1, ExpectedWarnings: 0}, // New version pair
				"description": {ExpectedErrors: 0, ExpectedWarnings: 1}, // New version pair
				"unhandled":   {ExpectedErrors: 0, ExpectedWarnings: 0},
			},
			"^.field2": {
				"type":        {ExpectedErrors: 1, ExpectedWarnings: 0}, // New field in new version
				"description": {ExpectedErrors: 0, ExpectedWarnings: 1}, // New field in new version
				"unhandled":   {ExpectedErrors: 0, ExpectedWarnings: 0},
			},
		},
		"v2 -> v3": {
			"^.field1": {
				"type":        {ExpectedErrors: 1, ExpectedWarnings: 0}, // New version pair
				"description": {ExpectedErrors: 0, ExpectedWarnings: 1}, // New version pair
				"unhandled":   {ExpectedErrors: 0, ExpectedWarnings: 0},
			},
			"^.field2": {
				"type":        {ExpectedErrors: 1, ExpectedWarnings: 0}, // New field changes
				"description": {ExpectedErrors: 0, ExpectedWarnings: 1}, // New field changes
				"unhandled":   {ExpectedErrors: 0, ExpectedWarnings: 0},
			},
		},
	}

	// Verify structure matches expectations
	for versionPair, expectedVersionResults := range expected {
		actualVersionResults, exists := result[versionPair]
		require.True(t, exists, "Expected version pair %s", versionPair)

		for fieldPath, expectedFieldResults := range expectedVersionResults {
			actualFieldResults, exists := actualVersionResults[fieldPath]
			require.True(t, exists, "Expected field path %s in %s", fieldPath, versionPair)

			// Create map for easier lookup
			actualByName := make(map[string]validations.ComparisonResult)
			for _, comp := range actualFieldResults {
				actualByName[comp.Name] = comp
			}

			for comparatorName, expectedComp := range expectedFieldResults {
				actualComp, exists := actualByName[comparatorName]
				require.True(t, exists, "Expected comparator %s for %s in %s",
					comparatorName, fieldPath, versionPair)

				assert.Len(t, actualComp.Errors, expectedComp.ExpectedErrors,
					"Wrong error count for %s/%s/%s", versionPair, fieldPath, comparatorName)
				assert.Len(t, actualComp.Warnings, expectedComp.ExpectedWarnings,
					"Wrong warning count for %s/%s/%s", versionPair, fieldPath, comparatorName)
			}
		}
	}
}

func TestValidator_New_DefaultValues(t *testing.T) {
	validator := New()

	assert.Equal(t, config.ConversionPolicyNone, validator.conversionPolicy)
	assert.Equal(t, config.EnforcementPolicyError, validator.unhandledEnforcement)
	assert.Empty(t, validator.comparators)
}

func TestValidator_New_WithOptions(t *testing.T) {
	typeComparator := &property.Type{}
	typeComparator.SetEnforcement(config.EnforcementPolicyError)

	validator := New(
		WithComparators(typeComparator),
		WithConversionPolicy(config.ConversionPolicyIgnore),
		WithUnhandledEnforcementPolicy(config.EnforcementPolicyWarn),
	)

	assert.Equal(t, config.ConversionPolicyIgnore, validator.conversionPolicy)
	assert.Equal(t, config.EnforcementPolicyWarn, validator.unhandledEnforcement)
	require.Len(t, validator.comparators, 1)
	assert.Equal(t, typeComparator, validator.comparators[0])
}

func TestValidator_Options_DefaultHandling(t *testing.T) {
	// Test that empty policies get set to defaults
	validator := New(
		WithConversionPolicy(""),
		WithUnhandledEnforcementPolicy(""),
	)

	assert.Equal(t, config.ConversionPolicyNone, validator.conversionPolicy)
	assert.Equal(t, config.EnforcementPolicyError, validator.unhandledEnforcement)
}

func TestValidator_Validate_EmptySchemas(t *testing.T) {
	validator := New(WithUnhandledEnforcementPolicy(config.EnforcementPolicyError))

	// Test with nil schemas
	crdA := createCRD(
		withVersion("v1", true, nil),
		withVersion("v2", true, nil),
	)

	crdB := createCRD(
		withVersion("v1", true, nil),
		withVersion("v2", true, nil),
	)

	result := validator.Validate(crdA, crdB)

	// Should handle nil schemas gracefully and generate expected structure
	require.NotNil(t, result, "Expected non-nil result")

	// Should have exactly one version pair: v1 -> v2
	assert.Len(t, result, 1, "Expected exactly one version pair")

	versionResults, exists := result["v1 -> v2"]
	require.True(t, exists, "Expected v1 -> v2 version pair")
	assert.Len(t, versionResults, 0, "Expected no version results")
}

func Test_numUnidirectionalPermutations(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected int
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: 0,
		},
		{
			name:     "single element",
			input:    []string{"a"},
			expected: 0,
		},
		{
			name:     "two elements",
			input:    []string{"a", "b"},
			expected: 1,
		},
		{
			name:     "three elements",
			input:    []string{"a", "b", "c"},
			expected: 3,
		},
		{
			name:     "four elements",
			input:    []string{"a", "b", "c", "d"},
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := numUnidirectionalPermutations(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

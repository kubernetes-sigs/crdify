// Copyright 2026 The Kubernetes Authors.
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

package property

import (
	"errors"
	"strings"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/crdify/pkg/config"
	internaltesting "sigs.k8s.io/crdify/pkg/validations/internal/testing"
)

func TestAllOf(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			Flagged:              false,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "no diff different order, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
					{MaxLength: ptr.To[int64](100)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MaxLength: ptr.To[int64](100)},
					{MinLength: ptr.To[int64](1)},
				},
			},
			Flagged:              false,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "both empty, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{},
			},
			Flagged:              false,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "net new allOf constraint, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			Flagged:              true,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "net new allOf constraint from nil, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			Flagged:              true,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "net new allOf constraint, addition policy Allow, not flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			Flagged: false,
			ComparableValidation: &AllOf{
				AllOfConfig: AllOfConfig{AdditionPolicy: AdditionPolicyAllow},
			},
		},
		{
			Name: "subschema added to existing allOf, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
					{MaxLength: ptr.To[int64](100)},
				},
			},
			Flagged:              true,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "subschema added to existing allOf, addition policy set to Disallow, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
					{MaxLength: ptr.To[int64](100)},
				},
			},
			Flagged: true,
			ComparableValidation: &AllOf{
				AllOfConfig: AllOfConfig{AdditionPolicy: AdditionPolicyDisallow},
			},
		},
		{
			Name: "subschema added to existing allOf, addition policy set to Allow, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
					{MaxLength: ptr.To[int64](100)},
				},
			},
			Flagged: false,
			ComparableValidation: &AllOf{
				AllOfConfig: AllOfConfig{AdditionPolicy: AdditionPolicyAllow},
			},
		},
		{
			Name: "subschema removed from existing allOf, flagged by default",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
					{MaxLength: ptr.To[int64](100)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			Flagged:              true,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "subschema removed from existing allOf, allowed via config",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
					{MaxLength: ptr.To[int64](100)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			Flagged: false,
			ComparableValidation: &AllOf{
				AllOfConfig: AllOfConfig{RemovalPolicy: RemovalPolicyAllow},
			},
		},
		{
			Name: "subschema changed, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](5)},
				},
			},
			Flagged:              true,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "subschema changed, removal allowed, flagged via addition",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](5)},
				},
			},
			Flagged: true,
			ComparableValidation: &AllOf{
				AllOfConfig: AllOfConfig{RemovalPolicy: RemovalPolicyAllow},
			},
		},
		{
			Name: "all subschemas removed, flagged by default",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			New:                  &apiextensionsv1.JSONSchemaProps{},
			Flagged:              true,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "all subschemas removed, allowed via config",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			New:     &apiextensionsv1.JSONSchemaProps{},
			Flagged: false,
			ComparableValidation: &AllOf{
				AllOfConfig: AllOfConfig{RemovalPolicy: RemovalPolicyAllow},
			},
		},
		{
			Name: "nested description changed in subschema, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"port": {Description: "old", Minimum: ptr.To[float64](1)},
						},
					},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"port": {Description: "new", Minimum: ptr.To[float64](1)},
						},
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "required and enum order changed in subschema, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{
						Required: []string{"a", "b"},
						Enum:     []apiextensionsv1.JSON{{Raw: []byte(`"x"`)}, {Raw: []byte(`"y"`)}},
					},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{
						Required: []string{"b", "a"},
						Enum:     []apiextensionsv1.JSON{{Raw: []byte(`"y"`)}, {Raw: []byte(`"x"`)}},
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "invalid subschema JSON, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{Enum: []apiextensionsv1.JSON{{Raw: []byte(`{`)}}},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{Enum: []apiextensionsv1.JSON{{Raw: []byte(`{`)}}},
				},
			},
			Flagged:              true,
			ComparableValidation: &AllOf{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &AllOf{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestAllOfReportsRemovalAndAddition(t *testing.T) {
	old := &apiextensionsv1.JSONSchemaProps{
		AllOf: []apiextensionsv1.JSONSchemaProps{
			{MinLength: ptr.To[int64](1)},
		},
	}
	updated := &apiextensionsv1.JSONSchemaProps{
		AllOf: []apiextensionsv1.JSONSchemaProps{
			{MinLength: ptr.To[int64](5)},
		},
	}

	allOf := &AllOf{}
	allOf.SetEnforcement(config.EnforcementPolicyError)

	result := allOf.Compare(old, updated)
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}

	if !strings.HasPrefix(result.Errors[0], ErrRemovedAllOf.Error()) {
		t.Fatalf("expected removal error first, got %q", result.Errors[0])
	}

	if !strings.HasPrefix(result.Errors[1], ErrAddedAllOf.Error()) {
		t.Fatalf("expected addition error second, got %q", result.Errors[1])
	}
}

func TestValidateAllOfConfig(t *testing.T) {
	testcases := []struct {
		name          string
		cfg           *AllOfConfig
		wantErr       error
		wantAddPolicy AdditionPolicy
		wantRemPolicy RemovalPolicy
	}{
		{
			name: "nil config",
			cfg:  nil,
		},
		{
			name:          "defaults both policies",
			cfg:           &AllOfConfig{},
			wantAddPolicy: AdditionPolicyDisallow,
			wantRemPolicy: RemovalPolicyDisallow,
		},
		{
			name:          "allows valid addition policy Allow",
			cfg:           &AllOfConfig{AdditionPolicy: AdditionPolicyAllow},
			wantAddPolicy: AdditionPolicyAllow,
			wantRemPolicy: RemovalPolicyDisallow,
		},
		{
			name:          "allows valid addition policy Disallow",
			cfg:           &AllOfConfig{AdditionPolicy: AdditionPolicyDisallow},
			wantAddPolicy: AdditionPolicyDisallow,
			wantRemPolicy: RemovalPolicyDisallow,
		},
		{
			name:    "invalid addition policy",
			cfg:     &AllOfConfig{AdditionPolicy: "invalid"},
			wantErr: errUnknownAdditionPolicy,
		},
		{
			name:          "allows valid removal policy Allow",
			cfg:           &AllOfConfig{RemovalPolicy: RemovalPolicyAllow},
			wantAddPolicy: AdditionPolicyDisallow,
			wantRemPolicy: RemovalPolicyAllow,
		},
		{
			name:          "allows valid removal policy Disallow",
			cfg:           &AllOfConfig{RemovalPolicy: RemovalPolicyDisallow},
			wantAddPolicy: AdditionPolicyDisallow,
			wantRemPolicy: RemovalPolicyDisallow,
		},
		{
			name:    "invalid removal policy",
			cfg:     &AllOfConfig{RemovalPolicy: "invalid"},
			wantErr: errUnknownRemovalPolicy,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateAllOfConfig(tc.cfg)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.cfg == nil {
				return
			}

			if tc.wantAddPolicy != "" && tc.cfg.AdditionPolicy != tc.wantAddPolicy {
				t.Fatalf("expected addition policy %q, got %q", tc.wantAddPolicy, tc.cfg.AdditionPolicy)
			}

			if tc.wantRemPolicy != "" && tc.cfg.RemovalPolicy != tc.wantRemPolicy {
				t.Fatalf("expected removal policy %q, got %q", tc.wantRemPolicy, tc.cfg.RemovalPolicy)
			}
		})
	}
}

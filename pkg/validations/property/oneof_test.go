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

package property

import (
	"errors"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	internaltesting "sigs.k8s.io/crdify/pkg/validations/internal/testing"
)

func TestOneOf(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			Flagged:              false,
			ComparableValidation: &OneOf{},
		},
		{
			Name: "no diff different order, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "string"},
					{Type: "integer"},
				},
			},
			Flagged:              false,
			ComparableValidation: &OneOf{},
		},
		{
			Name: "new oneOf constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			Flagged:              true,
			ComparableValidation: &OneOf{},
		},
		{
			Name: "removed oneOf subschema, removal policy not set, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "string"},
				},
			},
			Flagged:              true,
			ComparableValidation: &OneOf{},
		},
		{
			Name: "removed oneOf subschema, removal policy set to Disallow, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "string"},
				},
			},
			Flagged: true,
			ComparableValidation: &OneOf{
				OneOfConfig: OneOfConfig{
					RemovalPolicy: RemovalPolicyDisallow,
				},
			},
		},
		{
			Name: "removed oneOf subschema, removal policy set to Allow, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "string"},
				},
			},
			Flagged: false,
			ComparableValidation: &OneOf{
				OneOfConfig: OneOfConfig{
					RemovalPolicy: RemovalPolicyAllow,
				},
			},
		},
		{
			Name: "new allowed oneOf subschema added, addition policy not set, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			Flagged:              true,
			ComparableValidation: &OneOf{},
		},
		{
			Name: "new allowed oneOf subschema added, addition policy set to Disallow, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			Flagged: true,
			ComparableValidation: &OneOf{
				OneOfConfig: OneOfConfig{
					AdditionPolicy: AdditionPolicyDisallow,
				},
			},
		},
		{
			Name: "new allowed oneOf subschema added, addition policy set to Allow, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			Flagged: false,
			ComparableValidation: &OneOf{
				OneOfConfig: OneOfConfig{
					AdditionPolicy: AdditionPolicyAllow,
				},
			},
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
			ComparableValidation: &OneOf{},
		},
		{
			Name: "different field changed with oneOf, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
				OneOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			Flagged:              false,
			ComparableValidation: &OneOf{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestValidateOneOfConfig(t *testing.T) {
	testcases := []struct {
		name               string
		cfg                *OneOfConfig
		wantErr            error
		wantAdditionPolicy AdditionPolicy
		wantRemovalPolicy  RemovalPolicy
	}{
		{
			name: "nil config",
			cfg:  nil,
		},
		{
			name:               "defaults addition policy",
			cfg:                &OneOfConfig{},
			wantAdditionPolicy: AdditionPolicyDisallow,
			wantRemovalPolicy:  RemovalPolicyDisallow,
		},
		{
			name:               "allows valid addition policies",
			cfg:                &OneOfConfig{AdditionPolicy: AdditionPolicyAllow},
			wantAdditionPolicy: AdditionPolicyAllow,
			wantRemovalPolicy:  RemovalPolicyDisallow,
		},
		{
			name:    "invalid addition policy",
			cfg:     &OneOfConfig{AdditionPolicy: "invalid"},
			wantErr: errUnknownAdditionPolicy,
		},
		{
			name:               "allows valid removal policies",
			cfg:                &OneOfConfig{RemovalPolicy: RemovalPolicyAllow},
			wantAdditionPolicy: AdditionPolicyDisallow,
			wantRemovalPolicy:  RemovalPolicyAllow,
		},
		{
			name:    "invalid removal policy",
			cfg:     &OneOfConfig{RemovalPolicy: "invalid"},
			wantErr: errUnknownRemovalPolicy,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateOneOfConfig(tc.cfg)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.cfg != nil && tc.cfg.AdditionPolicy != tc.wantAdditionPolicy {
				t.Fatalf("expected addition policy %q, got %q", tc.wantAdditionPolicy, tc.cfg.AdditionPolicy)
			}

			if tc.cfg != nil && tc.cfg.RemovalPolicy != tc.wantRemovalPolicy {
				t.Fatalf("expected removal policy %q, got %q", tc.wantRemovalPolicy, tc.cfg.RemovalPolicy)
			}
		})
	}
}

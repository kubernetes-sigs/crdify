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
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/ptr"
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
			Name: "net new allOf constraint, addition policy Allow does not suppress, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				AllOf: []apiextensionsv1.JSONSchemaProps{
					{MinLength: ptr.To[int64](1)},
				},
			},
			Flagged: true,
			ComparableValidation: &AllOf{
				AllOfConfig: AllOfConfig{AdditionPolicy: AllOfAdditionPolicyAllow},
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
				AllOfConfig: AllOfConfig{AdditionPolicy: AllOfAdditionPolicyDisallow},
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
				AllOfConfig: AllOfConfig{AdditionPolicy: AllOfAdditionPolicyAllow},
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
				AllOfConfig: AllOfConfig{RemovalPolicy: AllOfRemovalPolicyAllow},
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
				AllOfConfig: AllOfConfig{RemovalPolicy: AllOfRemovalPolicyAllow},
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
				AllOfConfig: AllOfConfig{RemovalPolicy: AllOfRemovalPolicyAllow},
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
			ComparableValidation: &AllOf{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestValidateAllOfConfig(t *testing.T) {
	testcases := []struct {
		name          string
		cfg           *AllOfConfig
		wantErr       error
		wantAddPolicy AllOfAdditionPolicy
		wantRemPolicy AllOfRemovalPolicy
	}{
		{
			name: "nil config",
			cfg:  nil,
		},
		{
			name:          "defaults both policies",
			cfg:           &AllOfConfig{},
			wantAddPolicy: AllOfAdditionPolicyDisallow,
			wantRemPolicy: AllOfRemovalPolicyDisallow,
		},
		{
			name:          "allows valid addition policy Allow",
			cfg:           &AllOfConfig{AdditionPolicy: AllOfAdditionPolicyAllow},
			wantAddPolicy: AllOfAdditionPolicyAllow,
			wantRemPolicy: AllOfRemovalPolicyDisallow,
		},
		{
			name:          "allows valid addition policy Disallow",
			cfg:           &AllOfConfig{AdditionPolicy: AllOfAdditionPolicyDisallow},
			wantAddPolicy: AllOfAdditionPolicyDisallow,
			wantRemPolicy: AllOfRemovalPolicyDisallow,
		},
		{
			name:    "invalid addition policy",
			cfg:     &AllOfConfig{AdditionPolicy: "invalid"},
			wantErr: errUnknownAllOfAdditionPolicy,
		},
		{
			name:          "allows valid removal policy Allow",
			cfg:           &AllOfConfig{RemovalPolicy: AllOfRemovalPolicyAllow},
			wantAddPolicy: AllOfAdditionPolicyDisallow,
			wantRemPolicy: AllOfRemovalPolicyAllow,
		},
		{
			name:          "allows valid removal policy Disallow",
			cfg:           &AllOfConfig{RemovalPolicy: AllOfRemovalPolicyDisallow},
			wantAddPolicy: AllOfAdditionPolicyDisallow,
			wantRemPolicy: AllOfRemovalPolicyDisallow,
		},
		{
			name:    "invalid removal policy",
			cfg:     &AllOfConfig{RemovalPolicy: "invalid"},
			wantErr: errUnknownAllOfRemovalPolicy,
		},
	}

	for _, tc := range testcases {
		tc := tc
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

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

func TestNullable(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Nullable: true,
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Nullable: true,
			},
			Flagged:              false,
			ComparableValidation: &Nullable{},
		},
		{
			Name: "nullable allowed, flagged by default",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				Nullable: true,
			},
			Flagged:              true,
			ComparableValidation: &Nullable{},
		},
		{
			Name: "nullable allowed via config",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				Nullable: true,
			},
			Flagged: false,
			ComparableValidation: &Nullable{
				NullableConfig: NullableConfig{AdditionPolicy: NullableAdditionPolicyAllow},
			},
		},
		{
			Name: "nullable removed, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Nullable: true,
			},
			New:                  &apiextensionsv1.JSONSchemaProps{},
			Flagged:              true,
			ComparableValidation: &Nullable{},
		},
		{
			Name: "nullable removed allowed via config",
			Old: &apiextensionsv1.JSONSchemaProps{
				Nullable: true,
			},
			New:     &apiextensionsv1.JSONSchemaProps{},
			Flagged: false,
			ComparableValidation: &Nullable{
				NullableConfig: NullableConfig{RemovalPolicy: NullableRemovalPolicyAllow},
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
			ComparableValidation: &Nullable{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestValidateNullableConfig(t *testing.T) {
	testcases := []struct {
		name          string
		cfg           *NullableConfig
		wantErr       error
		wantAddPolicy NullableAdditionPolicy
		wantRemPolicy NullableRemovalPolicy
	}{
		{
			name: "nil config",
			cfg:  nil,
		},
		{
			name:          "defaults addition policy",
			cfg:           &NullableConfig{},
			wantAddPolicy: NullableAdditionPolicyDisallow,
			wantRemPolicy: NullableRemovalPolicyDisallow,
		},
		{
			name:          "allows valid addition policy",
			cfg:           &NullableConfig{AdditionPolicy: NullableAdditionPolicyAllow},
			wantAddPolicy: NullableAdditionPolicyAllow,
			wantRemPolicy: NullableRemovalPolicyDisallow,
		},
		{
			name:    "invalid addition policy mentions valid values",
			cfg:     &NullableConfig{AdditionPolicy: "invalid"},
			wantErr: errUnknownNullableAdditionPolicy,
		},
		{
			name:          "allows valid removal policy",
			cfg:           &NullableConfig{RemovalPolicy: NullableRemovalPolicyAllow},
			wantAddPolicy: NullableAdditionPolicyDisallow,
			wantRemPolicy: NullableRemovalPolicyAllow,
		},
		{
			name:    "invalid removal policy mentions valid values",
			cfg:     &NullableConfig{RemovalPolicy: "invalid"},
			wantErr: errUnknownNullableRemovalPolicy,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateNullableConfig(tc.cfg)
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

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

func TestExclusiveMaximum(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			Flagged:              false,
			ComparableValidation: &ExclusiveMaximum{},
		},
		{
			Name: "tightened flagged by default",
			Old: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum:          ptr.To(10.0),
				ExclusiveMaximum: true,
			},
			Flagged:              true,
			ComparableValidation: &ExclusiveMaximum{},
		},
		{
			Name: "tightened allowed via config",
			Old: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum:          ptr.To(10.0),
				ExclusiveMaximum: true,
			},
			Flagged: false,
			ComparableValidation: &ExclusiveMaximum{
				ExclusiveMaximumConfig: ExclusiveMaximumConfig{AdditionPolicy: ExclusiveMaximumAdditionPolicyAllow},
			},
		},
		{
			Name: "loosened flagged by default",
			Old: &apiextensionsv1.JSONSchemaProps{
				Maximum:          ptr.To(10.0),
				ExclusiveMaximum: true,
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			Flagged:              true,
			ComparableValidation: &ExclusiveMaximum{},
		},
		{
			Name: "loosening allowed via config",
			Old: &apiextensionsv1.JSONSchemaProps{
				Maximum:          ptr.To(10.0),
				ExclusiveMaximum: true,
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			Flagged: false,
			ComparableValidation: &ExclusiveMaximum{
				ExclusiveMaximumConfig: ExclusiveMaximumConfig{RemovalPolicy: ExclusiveMaximumRemovalPolicyAllow},
			},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestValidateExclusiveMaximumConfig(t *testing.T) {
	testcases := []struct {
		name         string
		cfg          *ExclusiveMaximumConfig
		wantErr      error
		wantAddition ExclusiveMaximumAdditionPolicy
		wantRemoval  ExclusiveMaximumRemovalPolicy
	}{
		{
			name: "nil config",
			cfg:  nil,
		},
		{
			name:         "defaults policies",
			cfg:          &ExclusiveMaximumConfig{},
			wantAddition: ExclusiveMaximumAdditionPolicyDisallow,
			wantRemoval:  ExclusiveMaximumRemovalPolicyDisallow,
		},
		{
			name:         "allows valid addition policy",
			cfg:          &ExclusiveMaximumConfig{AdditionPolicy: ExclusiveMaximumAdditionPolicyAllow},
			wantAddition: ExclusiveMaximumAdditionPolicyAllow,
			wantRemoval:  ExclusiveMaximumRemovalPolicyDisallow,
		},
		{
			name:         "allows valid removal policy",
			cfg:          &ExclusiveMaximumConfig{RemovalPolicy: ExclusiveMaximumRemovalPolicyAllow},
			wantAddition: ExclusiveMaximumAdditionPolicyDisallow,
			wantRemoval:  ExclusiveMaximumRemovalPolicyAllow,
		},
		{
			name:    "invalid addition policy mentions valid values",
			cfg:     &ExclusiveMaximumConfig{AdditionPolicy: "invalid"},
			wantErr: errUnknownExclusiveMaximumAdditionPolicy,
		},
		{
			name:    "invalid removal policy mentions valid values",
			cfg:     &ExclusiveMaximumConfig{RemovalPolicy: "invalid"},
			wantErr: errUnknownExclusiveMaximumRemovalPolicy,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateExclusiveMaximumConfig(tc.cfg)
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

			if tc.cfg.AdditionPolicy != tc.wantAddition {
				t.Fatalf("expected addition policy %q, got %q", tc.wantAddition, tc.cfg.AdditionPolicy)
			}

			if tc.cfg.RemovalPolicy != tc.wantRemoval {
				t.Fatalf("expected removal policy %q, got %q", tc.wantRemoval, tc.cfg.RemovalPolicy)
			}
		})
	}
}

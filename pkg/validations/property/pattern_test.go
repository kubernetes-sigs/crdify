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

func TestPattern(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Pattern: "^[a-z]+$",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Pattern: "^[a-z]+$",
			},
			Flagged:              false,
			ComparableValidation: &Pattern{},
		},
		{
			Name: "pattern added, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				Pattern: "^[a-z]+$",
			},
			Flagged:              true,
			ComparableValidation: &Pattern{},
		},
		{
			Name: "pattern changed, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Pattern: "^[a-z]+$",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Pattern: "^[A-Z]+$",
			},
			Flagged:              true,
			ComparableValidation: &Pattern{},
		},
		{
			Name: "pattern removed, flagged by default",
			Old: &apiextensionsv1.JSONSchemaProps{
				Pattern: "^[a-z]+$",
			},
			New:                  &apiextensionsv1.JSONSchemaProps{},
			Flagged:              true,
			ComparableValidation: &Pattern{},
		},
		{
			Name: "pattern removed, allowed via config",
			Old: &apiextensionsv1.JSONSchemaProps{
				Pattern: "^[a-z]+$",
			},
			New:     &apiextensionsv1.JSONSchemaProps{},
			Flagged: false,
			ComparableValidation: &Pattern{
				PatternConfig: PatternConfig{RemovalPolicy: PatternRemovalPolicyAllow},
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
			ComparableValidation: &Pattern{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestValidatePatternConfig(t *testing.T) {
	testcases := []struct {
		name              string
		cfg               *PatternConfig
		wantErr           error
		wantRemovalPolicy PatternRemovalPolicy
	}{
		{
			name: "nil config",
			cfg:  nil,
		},
		{
			name:              "defaults removal policy",
			cfg:               &PatternConfig{},
			wantRemovalPolicy: PatternRemovalPolicyDisallow,
		},
		{
			name:              "allows valid removal policy",
			cfg:               &PatternConfig{RemovalPolicy: PatternRemovalPolicyAllow},
			wantRemovalPolicy: PatternRemovalPolicyAllow,
		},
		{
			name:    "invalid removal policy mentions valid values",
			cfg:     &PatternConfig{RemovalPolicy: "invalid"},
			wantErr: errUnknownPatternRemovalPolicy,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePatternConfig(tc.cfg)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.cfg != nil && tc.cfg.RemovalPolicy != tc.wantRemovalPolicy {
				t.Fatalf("expected removal policy %q, got %q", tc.wantRemovalPolicy, tc.cfg.RemovalPolicy)
			}
		})
	}
}

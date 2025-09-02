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
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestAnyOf(t *testing.T) {
	{
		tests := []struct {
			name    string
			in      *AnyOfConfig
			wantErr bool
		}{
			{
				name:    "nil config",
				in:      nil,
				wantErr: true,
			},
			{
				name: "empty subschemas",
				in: &AnyOfConfig{
					AdditionPolicy: AdditionPolicyAllow,
					Subschemas:     nil,
				},
				wantErr: true,
			},
			{
				name: "valid config with allow",
				in: &AnyOfConfig{
					AdditionPolicy: AdditionPolicyAllow,
					Subschemas: []apiextensionsv1.JSONSchemaProps{
						{},
					},
				},
				wantErr: false,
			},
			{
				name: "valid config with disallow",
				in: &AnyOfConfig{
					AdditionPolicy: AdditionPolicyDisallow,
					Subschemas: []apiextensionsv1.JSONSchemaProps{
						{},
					},
				},
				wantErr: false,
			},
			{
				name: "invalid addition policy",
				in: &AnyOfConfig{
					AdditionPolicy: "invalid",
					Subschemas: []apiextensionsv1.JSONSchemaProps{
						{},
					},
				},
				wantErr: true,
			},
			{
				name: "empty addition policy defaults to disallow",
				in: &AnyOfConfig{
					AdditionPolicy: "",
					Subschemas: []apiextensionsv1.JSONSchemaProps{
						{},
					},
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := ValidateAnyOfConfig(tt.in)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateAnyOfConfig() error = %v, wantErr %v", err, tt.wantErr)
				}
				// Check defaulting behavior
				if tt.in != nil && tt.in.AdditionPolicy == "" && !tt.wantErr {
					if tt.in.AdditionPolicy != AdditionPolicyDisallow {
						t.Errorf("AdditionPolicy not defaulted to Disallow, got %q", tt.in.AdditionPolicy)
					}
				}
			})
		}
	}
}

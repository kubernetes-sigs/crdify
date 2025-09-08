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

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/crdify/pkg/config"
)

func TestOneOf(t *testing.T) {
	intSchema := &apiextensionsv1.JSONSchemaProps{Type: "integer"}
	stringSchema := &apiextensionsv1.JSONSchemaProps{Type: "string"}

	tests := []struct {
		name            string
		config          OneOfConfig
		oldSchema       *apiextensionsv1.JSONSchemaProps
		newSchema       *apiextensionsv1.JSONSchemaProps
		expectError     bool
		expectErrorMsgs []string
	}{
		{
			name:            "net new oneOf",
			config:          OneOfConfig{},
			oldSchema:       &apiextensionsv1.JSONSchemaProps{},
			newSchema:       &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*intSchema}},
			expectError:     true,
			expectErrorMsgs: []string{"oneOf constraint added when there was none previously"},
		},
		{
			name:            "removed oneOf",
			config:          OneOfConfig{},
			oldSchema:       &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*intSchema, *stringSchema}},
			newSchema:       &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*intSchema}},
			expectError:     true,
			expectErrorMsgs: []string{"allowed oneOf schemas removed"},
		},
		{
			name:            "added oneOf, disallowed",
			config:          OneOfConfig{},
			oldSchema:       &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*intSchema}},
			newSchema:       &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*intSchema, *stringSchema}},
			expectError:     true,
			expectErrorMsgs: []string{"allowed oneOf schemas added"},
		},
		{
			name:        "no change",
			config:      OneOfConfig{},
			oldSchema:   &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*intSchema, *stringSchema}},
			newSchema:   &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*stringSchema, *intSchema}},
			expectError: false,
		},
		{
			name:        "valid change with other fields",
			config:      OneOfConfig{},
			oldSchema:   &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*intSchema}, Description: "old"},
			newSchema:   &apiextensionsv1.JSONSchemaProps{OneOf: []apiextensionsv1.JSONSchemaProps{*intSchema}, Description: "new"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comparator := &OneOf{OneOfConfig: tt.config}
			comparator.SetEnforcement(config.EnforcementPolicyError)

			result := comparator.Compare(tt.oldSchema, tt.newSchema)

			if tt.expectError {
				require.NotEmpty(t, result.Errors)
				for _, msg := range tt.expectErrorMsgs {
					require.Contains(t, result.Errors[0], msg)
				}
			} else {
				require.Empty(t, result.Errors)
			}
			// Ensure the field is cleared after handling
			require.Nil(t, tt.newSchema.OneOf)
		})
	}
}

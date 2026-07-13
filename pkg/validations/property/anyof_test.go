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
	internaltesting "sigs.k8s.io/crdify/pkg/validations/internal/testing"
)

func TestAnyOf(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			Flagged:              false,
			ComparableValidation: &AnyOf{},
		},
		{
			Name: "no diff different order, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "string"},
					{Type: "integer"},
				},
			},
			Flagged:              false,
			ComparableValidation: &AnyOf{},
		},
		{
			Name: "new anyOf constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			Flagged:              true,
			ComparableValidation: &AnyOf{},
		},
		{
			Name: "removed anyOf subschema, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "string"},
				},
			},
			Flagged:              true,
			ComparableValidation: &AnyOf{},
		},
		{
			Name: "new anyOf subschema added, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
					{Type: "string"},
				},
			},
			Flagged:              false,
			ComparableValidation: &AnyOf{},
		},
		{
			Name: "anyOf subschema changed, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "string"},
				},
			},
			Flagged:              true,
			ComparableValidation: &AnyOf{},
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
			ComparableValidation: &AnyOf{},
		},
		{
			Name: "different field changed with anyOf, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
				AnyOf: []apiextensionsv1.JSONSchemaProps{
					{Type: "integer"},
				},
			},
			Flagged:              false,
			ComparableValidation: &AnyOf{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

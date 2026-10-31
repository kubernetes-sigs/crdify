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

package runner

import (
	"encoding/json"
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

func createCRD(name, replicasType string) *apiextensionsv1.CustomResourceDefinition {
	return &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Scope: apiextensionsv1.NamespaceScoped,
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"replicas": {Type: replicasType},
							},
						},
					},
				},
			},
		},
	}
}

func TestRunByName(t *testing.T) {
	cfg, err := config.Load("")
	require.NoError(t, err)

	run, err := New(cfg, DefaultRegistry())
	require.NoError(t, err)

	tests := []struct {
		name             string
		oldCrds          []*apiextensionsv1.CustomResourceDefinition
		newCrds          []*apiextensionsv1.CustomResourceDefinition
		expectedNames    []string
		expectedFailures []string
	}{
		{
			name:             "pairs by name whatever the order",
			oldCrds:          []*apiextensionsv1.CustomResourceDefinition{createCRD("widgets", "integer"), createCRD("gadgets", "integer")},
			newCrds:          []*apiextensionsv1.CustomResourceDefinition{createCRD("gadgets", "string"), createCRD("widgets", "integer")},
			expectedNames:    []string{"gadgets", "widgets"},
			expectedFailures: []string{"gadgets"},
		},
		{
			name:             "no changes",
			oldCrds:          []*apiextensionsv1.CustomResourceDefinition{createCRD("widgets", "integer"), createCRD("gadgets", "integer")},
			newCrds:          []*apiextensionsv1.CustomResourceDefinition{createCRD("widgets", "integer"), createCRD("gadgets", "integer")},
			expectedNames:    []string{"gadgets", "widgets"},
			expectedFailures: []string{},
		},
		{
			name:             "removed CustomResourceDefinition is a failure",
			oldCrds:          []*apiextensionsv1.CustomResourceDefinition{createCRD("widgets", "integer"), createCRD("gadgets", "integer")},
			newCrds:          []*apiextensionsv1.CustomResourceDefinition{createCRD("widgets", "integer")},
			expectedNames:    []string{"gadgets", "widgets"},
			expectedFailures: []string{"gadgets"},
		},
		{
			name:             "added CustomResourceDefinition is not validated",
			oldCrds:          []*apiextensionsv1.CustomResourceDefinition{createCRD("widgets", "integer")},
			newCrds:          []*apiextensionsv1.CustomResourceDefinition{createCRD("widgets", "integer"), createCRD("gadgets", "integer")},
			expectedNames:    []string{"widgets"},
			expectedFailures: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := run.RunByName(tt.oldCrds, tt.newCrds)

			assert.ElementsMatch(t, tt.expectedNames, slices.Collect(maps.Keys(results)))

			failures := []string{}

			for name, result := range results {
				if result.HasFailures() {
					failures = append(failures, name)
				}
			}

			assert.ElementsMatch(t, tt.expectedFailures, failures)
			assert.Equal(t, len(tt.expectedFailures) > 0, results.HasFailures())
		})
	}
}

func TestResultsByNameRender(t *testing.T) {
	results := ResultsByName{
		"widgets": {},
		"gadgets": {
			CRDValidation: []validations.ComparisonResult{{Name: "removed", Errors: []string{"CustomResourceDefinition removed"}}},
		},
	}

	plaintext, err := results.Render(FormatPlainText)
	require.NoError(t, err)
	assert.Equal(t, "gadgets:\n- removed - ERROR - CustomResourceDefinition removed\n", plaintext)

	markdown, err := results.Render(FormatMarkdown)
	require.NoError(t, err)
	assert.Equal(t, "## gadgets\n\n- **removed** - `ERROR` - CustomResourceDefinition removed\n\n", markdown)

	jsonOut, err := results.Render(FormatJSON)
	require.NoError(t, err)

	decoded := map[string]*Results{}
	require.NoError(t, json.Unmarshal([]byte(jsonOut), &decoded))
	assert.ElementsMatch(t, []string{"gadgets", "widgets"}, slices.Collect(maps.Keys(decoded)))
	assert.True(t, decoded["gadgets"].HasFailures())

	yamlOut, err := results.Render(FormatYAML)
	require.NoError(t, err)
	assert.Contains(t, yamlOut, "gadgets:")
	assert.Contains(t, yamlOut, "widgets:")

	_, err = results.Render("unknown")
	require.ErrorIs(t, err, errUnknownRenderFormat)
}

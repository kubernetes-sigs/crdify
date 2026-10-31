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

package decode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func crdYAML(name, kind string) string {
	return fmt.Sprintf(`apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: %s
spec:
  group: example.com
  names:
    kind: %s
`, name, kind)
}

func TestCRDs(t *testing.T) {
	widgets := crdYAML("widgets.example.com", "Widget")
	gadgets := crdYAML("gadgets.example.com", "Gadget")

	tests := []struct {
		name          string
		content       string
		expectedKinds []string
		expectedErr   error
	}{
		{
			name:          "single document",
			content:       widgets,
			expectedKinds: []string{"Widget"},
		},
		{
			name:          "single document with a leading separator",
			content:       "---\n" + widgets,
			expectedKinds: []string{"Widget"},
		},
		{
			name:          "multiple documents",
			content:       widgets + "---\n" + gadgets,
			expectedKinds: []string{"Widget", "Gadget"},
		},
		{
			name:          "documents of other kinds and empty documents are skipped",
			content:       "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: example\n---\n# comment only\n---\n---\n" + widgets,
			expectedKinds: []string{"Widget"},
		},
		{
			name:          "JSON",
			content:       `{"apiVersion": "apiextensions.k8s.io/v1", "kind": "CustomResourceDefinition", "metadata": {"name": "widgets.example.com"}, "spec": {"names": {"kind": "Widget"}}}`,
			expectedKinds: []string{"Widget"},
		},
		{
			name:        "no CustomResourceDefinition",
			content:     "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: example\n",
			expectedErr: errNoCRD,
		},
		{
			name:        "empty content",
			content:     "",
			expectedErr: errNoCRD,
		},
		{
			name:        "duplicate name",
			content:     widgets + "---\n" + gadgets + "---\n" + widgets,
			expectedErr: errDuplicateName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crds, err := CRDs([]byte(tt.content))
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, crds)

				return
			}

			require.NoError(t, err)

			kinds := []string{}
			for _, crd := range crds {
				kinds = append(kinds, crd.Spec.Names.Kind)
			}

			assert.Equal(t, tt.expectedKinds, kinds)
		})
	}
}

func TestCRDsInvalidContent(t *testing.T) {
	for name, content := range map[string]string{
		"document is not an object":               "- not an object\n",
		"CustomResourceDefinition does not parse": "apiVersion: apiextensions.k8s.io/v1\nkind: CustomResourceDefinition\nspec: 5\n",
	} {
		t.Run(name, func(t *testing.T) {
			crds, err := CRDs([]byte(content))
			require.Error(t, err)
			assert.Nil(t, crds)
		})
	}
}

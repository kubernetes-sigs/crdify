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

package composite

import (
	"context"
	"net/url"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/crdify/pkg/loaders/file"
	"sigs.k8s.io/crdify/pkg/loaders/scheme"
)

const twoCRDs = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: widgets.example.com
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: gadgets.example.com
`

// singleLoader is a Loader that does not implement MultiLoader.
type singleLoader struct{}

func (singleLoader) Load(_ context.Context, location *url.URL) (*apiextensionsv1.CustomResourceDefinition, error) {
	return &apiextensionsv1.CustomResourceDefinition{ObjectMeta: metav1.ObjectMeta{Name: location.Hostname()}}, nil
}

func names(crds []*apiextensionsv1.CustomResourceDefinition) []string {
	out := []string{}
	for _, crd := range crds {
		out = append(out, crd.Name)
	}

	return out
}

func TestLoadAll(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(filesystem, "/crds.yaml", []byte(twoCRDs), 0o600))

	loader := NewComposite(map[string]Loader{
		scheme.SchemeFile: file.New(filesystem),
		"single":          singleLoader{},
	})

	t.Run("MultiLoader returns every CustomResourceDefinition", func(t *testing.T) {
		crds, err := loader.LoadAll(t.Context(), "file:///crds.yaml")
		require.NoError(t, err)
		assert.Equal(t, []string{"widgets.example.com", "gadgets.example.com"}, names(crds))
	})

	t.Run("Load still returns the first CustomResourceDefinition", func(t *testing.T) {
		crd, err := loader.Load(t.Context(), "file:///crds.yaml")
		require.NoError(t, err)
		assert.Equal(t, "widgets.example.com", crd.Name)
	})

	t.Run("Loader without LoadAll returns its single CustomResourceDefinition", func(t *testing.T) {
		crds, err := loader.LoadAll(t.Context(), "single://foos.example.com")
		require.NoError(t, err)
		assert.Equal(t, []string{"foos.example.com"}, names(crds))
	})

	t.Run("unknown scheme", func(t *testing.T) {
		_, err := loader.LoadAll(t.Context(), "unknown://foos.example.com")
		require.ErrorIs(t, err, errNoLoader)
	})
}

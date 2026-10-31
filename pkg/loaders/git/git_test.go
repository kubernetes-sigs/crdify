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

package git

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestLoadAndLoadAll(t *testing.T) {
	dir := t.TempDir()

	repo, err := gogit.PlainInit(dir, false)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "crds.yaml"), []byte(twoCRDs), 0o600))

	worktree, err := repo.Worktree()
	require.NoError(t, err)
	_, err = worktree.Add("crds.yaml")
	require.NoError(t, err)
	_, err = worktree.Commit("add CRDs", &gogit.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
	})
	require.NoError(t, err)

	// Load and LoadAll open the git repository in the current directory.
	t.Chdir(dir)

	location, err := url.Parse("git://HEAD?path=crds.yaml")
	require.NoError(t, err)

	crds, err := New().LoadAll(t.Context(), location)
	require.NoError(t, err)
	require.Len(t, crds, 2)
	assert.Equal(t, "widgets.example.com", crds[0].Name)
	assert.Equal(t, "gadgets.example.com", crds[1].Name)

	crd, err := New().Load(t.Context(), location)
	require.NoError(t, err)
	assert.Equal(t, "widgets.example.com", crd.Name)
}

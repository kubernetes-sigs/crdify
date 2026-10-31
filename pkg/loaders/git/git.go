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

package git

import (
	"context"
	"fmt"
	"io"
	"net/url"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/crdify/pkg/loaders/decode"
)

// Git is a Loader implementation for loading a CustomResourceDefinition
// from a git repository.
type Git struct{}

// New returns a new instance of the Git Loader.
func New() *Git {
	return &Git{}
}

// Load loads the CustomResourceDefinition from the git revision and file path specified in the URL.
// It reads a query key named 'path' for the file path and uses the hostname for the revision.
// For example, 'git://main?path=foo/bar/file.yaml' would source the CustomResourceDefinition from the
// main branch of the git repository using the file 'foo/bar/file.yaml'.
func (g *Git) Load(_ context.Context, location *url.URL) (*apiextensionsv1.CustomResourceDefinition, error) {
	repo, hash, err := resolveRevision(location)
	if err != nil {
		return nil, err
	}

	crd, err := LoadCRDFileFromRepositoryWithRef(repo, hash, location.Query().Get("path"))
	if err != nil {
		return nil, fmt.Errorf("loading CRD: %w", err)
	}

	return crd, nil
}

// LoadAll determines the git revision and file path like Load does
// and loads every CustomResourceDefinition in the file.
// Empty documents and documents of any other kind are skipped.
func (g *Git) LoadAll(_ context.Context, location *url.URL) ([]*apiextensionsv1.CustomResourceDefinition, error) {
	repo, hash, err := resolveRevision(location)
	if err != nil {
		return nil, err
	}

	filePath := location.Query().Get("path")

	content, err := readFileFromRepositoryWithRef(repo, hash, filePath)
	if err != nil {
		return nil, fmt.Errorf("loading CRDs: %w", err)
	}

	crds, err := decode.CRDs(content)
	if err != nil {
		return nil, fmt.Errorf("decoding content of file %q from repo with ref %v: %w", filePath, hash, err)
	}

	return crds, nil
}

// resolveRevision opens the git repository in the current directory and resolves
// the revision named by the hostname of the provided URL.
func resolveRevision(location *url.URL) (*gogit.Repository, *plumbing.Hash, error) {
	repo, err := gogit.PlainOpen("")
	if err != nil {
		return nil, nil, fmt.Errorf("opening repository: %w", err)
	}

	rev := plumbing.Revision(location.Hostname())

	hash, err := repo.ResolveRevision(rev)
	if err != nil {
		return nil, nil, fmt.Errorf("calculating hash for revision %q: %w", rev, err)
	}

	return repo, hash, nil
}

// LoadCRDFileFromRepositoryWithRef loads a CustomResourceDefinition from the provided
// git.Repository using the provided git ref and file name.
func LoadCRDFileFromRepositoryWithRef(repo *gogit.Repository, ref *plumbing.Hash, filename string) (*apiextensionsv1.CustomResourceDefinition, error) {
	crdBytes, err := readFileFromRepositoryWithRef(repo, ref, filename)
	if err != nil {
		return nil, err
	}

	loadedCRD := &apiextensionsv1.CustomResourceDefinition{}

	err = yaml.Unmarshal(crdBytes, loadedCRD)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling content of blob for file %q from repo with ref %v: %w", filename, ref, err)
	}

	return loadedCRD, nil
}

// readFileFromRepositoryWithRef reads the content of the file with the provided name
// from the provided git.Repository using the provided git ref.
func readFileFromRepositoryWithRef(repo *gogit.Repository, ref *plumbing.Hash, filename string) ([]byte, error) {
	commit, err := repo.CommitObject(*ref)
	if err != nil {
		return nil, fmt.Errorf("getting commit object from repo for ref %v: %w", ref, err)
	}

	tree, err := repo.TreeObject(commit.TreeHash)
	if err != nil {
		return nil, fmt.Errorf("getting tree object from repo for tree hash %v: %w", commit.TreeHash, err)
	}

	file, err := tree.File(filename)
	if err != nil {
		return nil, fmt.Errorf("getting file %q from repo for tree hash %v: %w", filename, commit.TreeHash, err)
	}

	reader, err := file.Reader()
	if err != nil {
		return nil, fmt.Errorf("getting reader for blob for file %q from repo with ref %v: %w", filename, commit.TreeHash, err)
	}
	//nolint:errcheck
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading content of blob for file %q from repo with ref %v: %w", filename, commit.TreeHash, err)
	}

	return content, nil
}

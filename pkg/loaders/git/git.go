package git

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/go-git/go-git/v5"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// TODO: Support remote git references
type Git struct{}

func NewGit() *Git {
	return &Git{}
}

func (g *Git) Load(ctx context.Context, location url.URL) (*apiextensionsv1.CustomResourceDefinition, error) {
	filePath := location.Query().Get("path")
	repo, err := gogit.PlainOpen("")
	if err != nil {
		return nil, fmt.Errorf("opening repository: %w", err)
	}

	rev := plumbing.Revision(location.Hostname())
	hash, err := repo.ResolveRevision(rev)
	if err != nil {
		return nil, fmt.Errorf("calculating hash for revision %q: %w", rev, err)
	}

	crd, err := LoadCRDFileFromRepositoryWithRef(repo, hash, filePath)
	if err != nil {
		return nil, fmt.Errorf("loading CRD: %w", err)
	}

	return crd, nil
}

func LoadCRDFileFromRepositoryWithRef(repo *git.Repository, ref *plumbing.Hash, filename string) (*apiextensionsv1.CustomResourceDefinition, error) {
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
	defer reader.Close()

	crdBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading content of blob for file %q from repo with ref %v: %w", filename, commit.TreeHash, err)
	}

	loadedCRD := &apiextensionsv1.CustomResourceDefinition{}
	err = yaml.Unmarshal(crdBytes, loadedCRD)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling content of blob for file %q from repo with ref %v: %w", filename, commit.TreeHash, err)
	}

	return loadedCRD, nil
}

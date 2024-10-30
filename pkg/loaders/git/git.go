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
)

// TODO: Support remote git references
type Git struct{}

func NewGit() *Git {
	return &Git{}
}

func (g *Git) Load(ctx context.Context, location url.URL) (*apiextensionsv1.CustomResourceDefinition, error) {
	filePath := location.Query().Get("path")
	// TODO: get only the specific file. Helpful link:
	// https://stackoverflow.com/questions/73561564/how-to-checkout-a-specific-single-file-to-inspect-it-using-go-git
	repo, err := gogit.PlainOpen("")
	if err != nil {
		return nil, fmt.Errorf("opening repository: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("getting worktree: %w", err)
	}

	originalRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("getting HEAD: %w", err)
	}

	defer func() {
		if err := wt.Checkout(&gogit.CheckoutOptions{Hash: originalRef.Hash()}); err != nil {
			fmt.Println("WARNING: failed to checkout your original working commit after loading:", err)
		}
	}()

	rev := plumbing.Revision(location.Hostname())
	hash, err := repo.ResolveRevision(rev)
	if err != nil {
		return nil, fmt.Errorf("calculating hash for revision %q: %w", rev, err)
	}
	err = wt.Checkout(&gogit.CheckoutOptions{Hash: *hash})
	if err != nil {
		return nil, fmt.Errorf("checking out hash %v: %w", hash, err)
	}

	file, err := wt.Filesystem.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file %q in revision %q: %w", filePath, rev, err)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", filePath, err)
	}

	crd := &apiextensionsv1.CustomResourceDefinition{}
	err = yaml.Unmarshal(fileBytes, crd)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling contents of file %q: %w", filePath, err)
	}

	return crd, nil
}

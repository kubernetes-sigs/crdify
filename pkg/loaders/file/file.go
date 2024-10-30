package file

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type File struct {
	filesystem afero.Fs
}

func NewFile(filesystem afero.Fs) *File {
	return &File{
		filesystem: filesystem,
	}
}

func (f *File) Load(ctx context.Context, location url.URL) (*apiextensionsv1.CustomResourceDefinition, error) {
	filePath, err := filepath.Abs(path.Join(location.Hostname(), location.Path))
	if err != nil {
		return nil, fmt.Errorf("ensuring absolute path: %w", err)
	}

	file, err := f.filesystem.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", filePath, err)
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

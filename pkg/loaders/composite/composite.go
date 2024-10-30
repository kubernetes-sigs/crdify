package composite

import (
	"context"
	"fmt"
	"net/url"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Loader interface {
	Load(context.Context, url.URL) (*apiextensionsv1.CustomResourceDefinition, error)
}

type Composite struct {
	loaders map[string]Loader
}

type CompositeOption func(*Composite)

func WithLoaders(loaders map[string]Loader) CompositeOption {
	return func(c *Composite) {
		c.loaders = loaders
	}
}

func NewComposite(opts ...CompositeOption) *Composite {
	composite := &Composite{
		loaders: map[string]Loader{},
	}

	for _, opt := range opts {
		opt(composite)
	}

	return composite
}

func (c *Composite) Load(ctx context.Context, location url.URL) (*apiextensionsv1.CustomResourceDefinition, error) {
	loader, ok := c.loaders[location.Scheme]
	if !ok {
		return nil, fmt.Errorf("no loader found for scheme %q", location.Scheme)
	}

	return loader.Load(ctx, location)
}

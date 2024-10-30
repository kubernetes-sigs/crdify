package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Kubernetes struct{}

func NewKubernetes() *Kubernetes {
	return &Kubernetes{}
}

func (k *Kubernetes) Load(ctx context.Context, location url.URL) (*apiextensionsv1.CustomResourceDefinition, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	apiextensionsClient, err := apiextensionsv1client.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating CustomResourceDefinition client: %w", err)
	}
	crdClient := apiextensionsClient.CustomResourceDefinitions()

	err = ValidateHostname(location.Hostname())
	if err != nil {
		return nil, fmt.Errorf("validating hostname: %w", err)
	}

	return crdClient.Get(ctx, location.Hostname(), v1.GetOptions{})
}

func ValidateHostname(hostname string) error {
	if hostname == "" {
		return errors.New("hostname is empty. The hostname should be the name of the CustomResourceDefinition to load.")
	}

	if errs := validation.IsDNS1123Subdomain(hostname); len(errs) > 0 {
		actualErrs := []error{}
		for _, errString := range errs {
			actualErrs = append(actualErrs, errors.New(errString))
		}
		err := errors.Join(actualErrs...)
		return fmt.Errorf("hostname %q is not a valid CustomResourceDefinition name: %w . CustomResourceDefinition names are required to be valid DNS Subdomains as outlined in https://tools.ietf.org/html/rfc1123", hostname, err)
	}

	return nil
}

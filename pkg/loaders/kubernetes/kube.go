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
	"k8s.io/client-go/rest"
)

// KubeConfigFunc is a function with no input parameters that returns a rest.Config
// for building a client to interact with a Kubernetes cluster or an error.
type KubeConfigFunc func() (*rest.Config, error)

// Kubernetes is a Loader implementation for sourcing a CustomResourceDefinition from a Kubernetes cluster.
type Kubernetes struct {
	// cfgFunc is a function to source the rest.Config for building
	// a client for fetching a CustomResourceDefinition from a Kubernetes cluster.
	// We use a function here so that a configuration for interacting with a Kubernetes cluster
	// is only run when this loader has been intentionally called.
	cfgFunc KubeConfigFunc
}

// New returns a new instance of a Kubernetes Loader, configured with
// the provided function for loading a Kubeconfig for interacting with a
// Kubernetes cluster.
func New(cfgFunc KubeConfigFunc) *Kubernetes {
	return &Kubernetes{
		cfgFunc: cfgFunc,
	}
}

// Load loads a CustomResourceDefinition from a Kubernetes cluster using the same configurations as tools like
// kubectl. It uses the hostname of the provided URL as the name of the CustomResourceDefinition to fetch from the cluster.
func (k *Kubernetes) Load(ctx context.Context, location *url.URL) (*apiextensionsv1.CustomResourceDefinition, error) {
	cfg, err := k.cfgFunc()
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

	crd, err := crdClient.Get(ctx, location.Hostname(), v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting CustomResourceDefinition: %w", err)
	}

	return crd, nil
}

// ValidateHostname validates that the provided hostname of a URL
// is a valid RFC1123 DNS Subdomain name, which is what valid
// CustomResourceDefinition resource names must follow.
func ValidateHostname(hostname string) error {
	if hostname == "" {
		return errEmptyHostname
	}

	if errs := validation.IsDNS1123Subdomain(hostname); len(errs) > 0 {
		actualErrs := []error{}

		for _, errString := range errs {
			//nolint:err113
			actualErrs = append(actualErrs, errors.New(errString))
		}

		err := errors.Join(actualErrs...)

		return fmt.Errorf("hostname %q is not a valid CustomResourceDefinition name: %w . CustomResourceDefinition names are required to be valid DNS Subdomains as outlined in https://tools.ietf.org/html/rfc1123", hostname, err)
	}

	return nil
}

var errEmptyHostname = errors.New("hostname is empty - the hostname should be the name of the CustomResourceDefinition to load")

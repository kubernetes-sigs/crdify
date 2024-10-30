package crd

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Scope struct{}

func (s *Scope) Name() string {
	return "Scope"
}

func (s *Scope) Validate(old, new *apiextensionsv1.CustomResourceDefinition) error {
	if old.Spec.Scope != new.Spec.Scope {
		return fmt.Errorf("scope changed from %q to %q", old.Spec.Scope, new.Spec.Scope)
	}
	return nil
}

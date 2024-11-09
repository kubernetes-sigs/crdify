package crd

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Scope struct{}

func (s *Scope) Name() string {
	return "scope"
}

func (s *Scope) Validate(old, new *apiextensionsv1.CustomResourceDefinition) ValidationResult {
	vr := &validationResult{
		Validation: s.Name(),
	}
	if old.Spec.Scope != new.Spec.Scope {
		vr.Err = fmt.Sprintf("scope changed from %q to %q", old.Spec.Scope, new.Spec.Scope)
	}
	return vr
}

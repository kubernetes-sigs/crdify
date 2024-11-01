package crd

import (
	"fmt"

	"github.com/everettraven/crd-diff/pkg/validations/results"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Scope struct{}

func (s *Scope) Name() string {
	return "Scope"
}

func (s *Scope) Validate(old, new *apiextensionsv1.CustomResourceDefinition) *results.Result {
	if old.Spec.Scope != new.Spec.Scope {
		return &results.Result{
			Error:      fmt.Errorf("scope changed from %q to %q", old.Spec.Scope, new.Spec.Scope),
			Subresults: []*results.Result{},
		}
	}
	return nil
}

package scope

import (
	"testing"

	internaltesting "github.com/everettraven/crd-diff/pkg/validations/internal/testing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestScope(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.CustomResourceDefinition]{
		{
			Name: "no scope change, not flagged",
			Old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Scope: apiextensionsv1.ClusterScoped,
				},
			},
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Scope: apiextensionsv1.ClusterScoped,
				},
			},
			Flagged:              false,
			ComparableValidation: &Scope{},
		},
		{
			Name: "scope change, flagged",
			Old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Scope: apiextensionsv1.ClusterScoped,
				},
			},
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Scope: apiextensionsv1.NamespaceScoped,
				},
			},
			Flagged:              true,
			ComparableValidation: &Scope{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

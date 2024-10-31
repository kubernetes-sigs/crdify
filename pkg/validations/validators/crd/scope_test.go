package crd

import (
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestScope(t *testing.T) {
	for _, tc := range []struct {
		name        string
		old         *apiextensionsv1.CustomResourceDefinition
		new         *apiextensionsv1.CustomResourceDefinition
		shouldError bool
		scope       *Scope
	}{
		{
			name: "no scope change, no error",
			old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Scope: apiextensionsv1.ClusterScoped,
				},
			},
			new: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Scope: apiextensionsv1.ClusterScoped,
				},
			},
			scope: &Scope{},
		},
		{
			name: "scope change, error",
			old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Scope: apiextensionsv1.ClusterScoped,
				},
			},
			new: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Scope: apiextensionsv1.NamespaceScoped,
				},
			},
			shouldError: true,
			scope:       &Scope{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.scope.Validate(tc.old, tc.new)
			require.Equal(t, tc.shouldError, err != nil)
		})
	}
}

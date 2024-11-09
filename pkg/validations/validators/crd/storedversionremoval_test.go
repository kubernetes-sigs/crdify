package crd

import (
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestStoredVersionRemoval(t *testing.T) {
	for _, tc := range []struct {
		name        string
		old         *apiextensionsv1.CustomResourceDefinition
		new         *apiextensionsv1.CustomResourceDefinition
		shouldError bool
		svr         *StoredVersionRemoval
	}{
		{
			name: "no stored versions, no error",
			new: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
						},
					},
				},
			},
			old: &apiextensionsv1.CustomResourceDefinition{},
			svr: &StoredVersionRemoval{},
		},
		{
			name: "stored versions, no stored version removed, no error",
			new: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
						},
						{
							Name: "v1alpha2",
						},
					},
				},
			},
			old: &apiextensionsv1.CustomResourceDefinition{
				Status: apiextensionsv1.CustomResourceDefinitionStatus{
					StoredVersions: []string{
						"v1alpha1",
					},
				},
			},
			svr: &StoredVersionRemoval{},
		},
		{
			name: "stored versions, stored version removed, error",
			new: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha2",
						},
					},
				},
			},
			old: &apiextensionsv1.CustomResourceDefinition{
				Status: apiextensionsv1.CustomResourceDefinitionStatus{
					StoredVersions: []string{
						"v1alpha1",
					},
				},
			},
			shouldError: true,
			svr:         &StoredVersionRemoval{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.svr.Validate(tc.old, tc.new)
			require.Equal(t, tc.shouldError, result.Error(0) != nil)
		})
	}
}

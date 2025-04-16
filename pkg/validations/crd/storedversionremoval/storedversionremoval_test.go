package storedversionremoval

import (
	"testing"

	internaltesting "github.com/everettraven/crd-diff/pkg/validations/internal/testing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestStoredVersionRemoval(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.CustomResourceDefinition]{
		{
			Name: "no stored versions, not flagged",
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
						},
					},
				},
			},
			Old:                  &apiextensionsv1.CustomResourceDefinition{},
			Flagged:              false,
			ComparableValidation: &StoredVersionRemoval{},
		},
		{
			Name: "stored versions, no stored version removed, not flagged",
			New: &apiextensionsv1.CustomResourceDefinition{
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
			Old: &apiextensionsv1.CustomResourceDefinition{
				Status: apiextensionsv1.CustomResourceDefinitionStatus{
					StoredVersions: []string{
						"v1alpha1",
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &StoredVersionRemoval{},
		},
		{
			Name: "stored versions, stored version removed, flagged",
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha2",
						},
					},
				},
			},
			Old: &apiextensionsv1.CustomResourceDefinition{
				Status: apiextensionsv1.CustomResourceDefinitionStatus{
					StoredVersions: []string{
						"v1alpha1",
					},
				},
			},
			Flagged:              true,
			ComparableValidation: &StoredVersionRemoval{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

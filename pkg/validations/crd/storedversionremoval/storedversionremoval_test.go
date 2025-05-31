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

package storedversionremoval

import (
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	internaltesting "sigs.k8s.io/crdify/pkg/validations/internal/testing"
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

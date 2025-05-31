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

package existingfieldremoval

import (
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	internaltesting "sigs.k8s.io/crdify/pkg/validations/internal/testing"
)

func TestExistingFieldRemoval(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.CustomResourceDefinition]{
		{
			Name: "no existing field removed, not flagged",
			Old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &ExistingFieldRemoval{},
		},
		{
			Name: "existing field removed, flagged",
			Old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
										"fieldTwo": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			Flagged:              true,
			ComparableValidation: &ExistingFieldRemoval{},
		},
		{
			Name: "new version is added with the field removed, not flagged",
			Old: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
										"fieldTwo": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			New: &apiextensionsv1.CustomResourceDefinition{
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name: "v1alpha1",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
										"fieldTwo": {
											Type: "string",
										},
									},
								},
							},
						},
						{
							Name: "v1alpha2",
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"fieldOne": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &ExistingFieldRemoval{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

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

package property

import (
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/ptr"
	internaltesting "sigs.k8s.io/crdify/pkg/validations/internal/testing"
)

func TestMinimum(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			Flagged:              false,
			ComparableValidation: &Minimum{},
		},
		{
			Name: "new minimum constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			Flagged:              true,
			ComparableValidation: &Minimum{},
		},
		{
			Name: "minimum constraint decreased, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(20.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			Flagged:              false,
			ComparableValidation: &Minimum{},
		},
		{
			Name: "minimum constraint increased, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(20.0),
			},
			Flagged:              true,
			ComparableValidation: &Minimum{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &Minimum{},
		},
		{
			Name: "exclusiveMinimum changed from false to true, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Minimum:          ptr.To(10.0),
				ExclusiveMinimum: false,
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Minimum:          ptr.To(10.0),
				ExclusiveMinimum: true,
			},
			Flagged:              true,
			ComparableValidation: &Minimum{},
		},
		{
			Name: "exclusiveMinimum changed from true to false, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Minimum:          ptr.To(10.0),
				ExclusiveMinimum: true,
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Minimum:          ptr.To(10.0),
				ExclusiveMinimum: false,
			},
			Flagged:              false,
			ComparableValidation: &Minimum{},
		},
		{
			Name: "net new exclusiveMinimum, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Minimum:          ptr.To(10.0),
				ExclusiveMinimum: true,
			},
			Flagged:              true,
			ComparableValidation: &Minimum{},
		},
		{
			Name: "no diff exclusiveMinimum, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Minimum:          ptr.To(10.0),
				ExclusiveMinimum: true,
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Minimum:          ptr.To(10.0),
				ExclusiveMinimum: true,
			},
			Flagged:              false,
			ComparableValidation: &Minimum{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestMinItems(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MinItems{},
		},
		{
			Name: "new minItems constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MinItems{},
		},
		{
			Name: "minItems constraint decreased, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MinItems{},
		},
		{
			Name: "minItems constraint increased, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(20)),
			},
			Flagged:              true,
			ComparableValidation: &MinItems{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &MinItems{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestMinLength(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MinLength{},
		},
		{
			Name: "new minLength constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MinLength{},
		},
		{
			Name: "minLength constraint decreased, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MinLength{},
		},
		{
			Name: "minLength constraint increased, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(20)),
			},
			Flagged:              true,
			ComparableValidation: &MinLength{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &MinLength{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestMinProperties(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MinProperties{},
		},
		{
			Name: "new minProperties constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MinProperties{},
		},
		{
			Name: "minProperties constraint decreased, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MinProperties{},
		},
		{
			Name: "minProperties constraint increased, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(20)),
			},
			Flagged:              true,
			ComparableValidation: &MinProperties{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &MinProperties{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

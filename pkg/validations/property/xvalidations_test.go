// Copyright The Kubernetes Authors.
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

	"k8s.io/utils/ptr"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	internaltesting "sigs.k8s.io/crdify/pkg/validations/internal/testing"
)

// propsWithRules is a helper that returns a JSONSchemaProps populated with the provided ValidationRules
func propsWithRules(rules ...apiextensionsv1.ValidationRule) *apiextensionsv1.JSONSchemaProps {
	return &apiextensionsv1.JSONSchemaProps{XValidations: rules}
}

func TestXValidations(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "identical rules, not flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			Flagged:              false,
			ComparableValidation: &XValidations{},
		},
		{
			Name:                 "both empty, not flagged",
			Old:                  &apiextensionsv1.JSONSchemaProps{},
			New:                  &apiextensionsv1.JSONSchemaProps{},
			Flagged:              false,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "rule added from empty, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "all rules removed, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			New:                  &apiextensionsv1.JSONSchemaProps{},
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "rule added, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
				apiextensionsv1.ValidationRule{Rule: "self.size() > 0", Message: "must not be empty"},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "rule removed, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
				apiextensionsv1.ValidationRule{Rule: "self.size() > 0", Message: "must not be empty"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "rule expression modified, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self != oldSelf", Message: "field is immutable"},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "optionalOldSelf added, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", OptionalOldSelf: ptr.To(true)},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "optionalOldSelf removed, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", OptionalOldSelf: ptr.To(true)},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "optionalOldSelf changed true->false, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", OptionalOldSelf: ptr.To(true)},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", OptionalOldSelf: ptr.To(false)},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "reason added, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", Reason: ptr.To(apiextensionsv1.FieldValueForbidden)},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "reason removed, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", Reason: ptr.To(apiextensionsv1.FieldValueForbidden)},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "reason changed, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", Reason: ptr.To(apiextensionsv1.FieldValueForbidden)},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", Reason: ptr.To(apiextensionsv1.FieldValueInvalid)},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "message changed, not flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "original immutable message"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "different message text"},
			),
			Flagged:              false,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "messageExpression changed, not flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", MessageExpression: "'original expr'"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", MessageExpression: "'different expr'"},
			),
			Flagged:              false,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "fieldPath changed, not flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", FieldPath: ".field"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable", FieldPath: ".changedField"},
			),
			Flagged:              false,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "multiple rules with identical Rule expression, not flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "A: field is immutable", OptionalOldSelf: ptr.To(true)},
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "B: field is immutable", OptionalOldSelf: ptr.To(false)},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "A: field is immutable", OptionalOldSelf: ptr.To(true)},
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "B: field is immutable", OptionalOldSelf: ptr.To(false)},
			),
			Flagged:              false,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "one of multiple same-Rule rules removed, flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "A: field is immutable", OptionalOldSelf: ptr.To(true)},
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "B: field is immutable", OptionalOldSelf: ptr.To(false)},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "A: field is immutable", OptionalOldSelf: ptr.To(true)},
			),
			Flagged:              true,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "rules reordered, not flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self.size() > 0", Message: "must not be empty"},
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
				apiextensionsv1.ValidationRule{Rule: "self.matches('[a-z]+')", Message: "must be lowercase"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self.matches('[a-z]+')", Message: "must be lowercase"},
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
				apiextensionsv1.ValidationRule{Rule: "self.size() > 0", Message: "must not be empty"},
			),
			Flagged:              false,
			ComparableValidation: &XValidations{},
		},
		{
			Name: "duplicate identical rules collapse, not flagged",
			Old: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			New: propsWithRules(
				apiextensionsv1.ValidationRule{Rule: "self == oldSelf", Message: "field is immutable"},
			),
			Flagged:              false,
			ComparableValidation: &XValidations{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

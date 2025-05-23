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

package testing

import (
	"testing"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/validations"
	"github.com/stretchr/testify/assert"
)

// ComparableValidation is a generic interface that represents
// a validation that can perform comparisons for the provided
// validations.Comparable.
type ComparableValidation[T validations.Comparable] interface {
	validations.Validation
	validations.Comparator[T]
}

// Testcase defines a single test case for a comparable validation.
type Testcase[T validations.Comparable] struct {
	Name                 string
	Old                  *T
	New                  *T
	Flagged              bool
	ComparableValidation ComparableValidation[T]
}

// RunTestcases runs the provided set of test cases for the comparable validation.
func RunTestcases[T validations.Comparable](t *testing.T, testcases ...Testcase[T]) {
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			val := testcase.ComparableValidation

			t.Log("with enforcement policy error")
			val.SetEnforcement(config.EnforcementPolicyError)
			result := val.Compare(copyComparable(testcase.Old), copyComparable(testcase.New))
			assert.Equal(t, testcase.Flagged, len(result.Errors) > 0, "unexpected state", "result errors", result.Errors)

			t.Log("with enforcement policy warn")
			val.SetEnforcement(config.EnforcementPolicyWarn)
			result = val.Compare(copyComparable(testcase.Old), copyComparable(testcase.New))
			assert.Equal(t, testcase.Flagged, len(result.Warnings) > 0, "unexpected state", "result warnings", result.Warnings)

			t.Log("with enforcement policy none")
			val.SetEnforcement(config.EnforcementPolicyNone)
			result = val.Compare(copyComparable(testcase.Old), copyComparable(testcase.New))
			assert.True(t, len(result.Errors) == 0, "unexpected state", "result errors", result.Errors)
			assert.True(t, len(result.Warnings) == 0, "unexpected state", "result warnings", result.Warnings)
		})
	}
}

func copyComparable[T validations.Comparable](in *T) *T {
	cIn := *in
	return &cIn
}

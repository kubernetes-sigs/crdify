package testing

import (
	"testing"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/validations"
	"github.com/stretchr/testify/assert"
)

type ComparableValidation[T validations.Comparable] interface {
	validations.Validation
	validations.Comparator[T]
}

type DeepCopyableComparable[T validations.Comparable] interface {
	DeepCopy() *T
}

type Testcase[T validations.Comparable] struct {
	Name                 string
	Old                  *T
	New                  *T
	Flagged              bool
	ComparableValidation ComparableValidation[T]
}

func RunTestcases[T validations.Comparable](t *testing.T, testcases ...Testcase[T]) {
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			val := testcase.ComparableValidation

			t.Log("with enforcement policy error")
			val.SetEnforcement(config.EnforcementPolicyError)
			result := val.Compare(copy(testcase.Old), copy(testcase.New))
			assert.Equal(t, testcase.Flagged, len(result.Errors) > 0, "unexpected state", "result errors", result.Errors)

			t.Log("with enforcement policy warn")
			val.SetEnforcement(config.EnforcementPolicyWarn)
			result = val.Compare(copy(testcase.Old), copy(testcase.New))
			assert.Equal(t, testcase.Flagged, len(result.Warnings) > 0, "unexpected state", "result warnings", result.Warnings)

			t.Log("with enforcement policy none")
			val.SetEnforcement(config.EnforcementPolicyNone)
			result = val.Compare(copy(testcase.Old), copy(testcase.New))
			assert.True(t, len(result.Errors) == 0, "unexpected state", "result errors", result.Errors)
			assert.True(t, len(result.Warnings) == 0, "unexpected state", "result warnings", result.Warnings)
		})
	}
}

func copy[T validations.Comparable](in *T) *T {
	cIn := *in
	return &cIn
}

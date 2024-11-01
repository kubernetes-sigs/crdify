package property

import (
	"fmt"

	"github.com/everettraven/crd-diff/pkg/validations/results"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Required struct{}

func (r *Required) Name() string {
	return "Required"
}

func (r *Required) Validate(diff Diff) (bool, *results.Result) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		newProperty.Required = []string{}
		oldProperty.Required = []string{}
		return NewDiff(oldProperty, newProperty)
	}

	oldRequired := sets.New(diff.Old().Required...)
	newRequired := sets.New(diff.New().Required...)
	diffRequired := newRequired.Difference(oldRequired)
	var err error

	if diffRequired.Len() > 0 {
		err = fmt.Errorf("new required fields %v added", diffRequired.UnsortedList())
	}

	return IsHandled(diff, reset), &results.Result{
		Error:      err,
		Subresults: []*results.Result{},
	}
}

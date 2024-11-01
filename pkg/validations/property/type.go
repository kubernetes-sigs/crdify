package property

import (
	"fmt"

	"github.com/everettraven/crd-diff/pkg/validations/results"
)

type Type struct{}

func (t *Type) Name() string {
	return "Type"
}

func (t *Type) Validate(diff Diff) (bool, *results.Result) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Type = ""
		newProperty.Type = ""
		return NewDiff(oldProperty, newProperty)
	}

	var err error
	if diff.Old().Type != diff.New().Type {
		err = fmt.Errorf("type changed from %q to %q", diff.Old().Type, diff.New().Type)
	}

	return IsHandled(diff, reset), &results.Result{
		Error:      err,
		Subresults: []*results.Result{},
	}
}

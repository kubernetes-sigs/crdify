package property

import (
	"fmt"
)

type Type struct{}

func (t *Type) Name() string {
	return "type"
}

func (t *Type) Validate(diff Diff) (bool, error) {
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

	return IsHandled(diff, reset), err
}

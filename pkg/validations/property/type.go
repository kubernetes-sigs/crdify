package property

import (
	"fmt"
)

type Type struct{}

func (t *Type) Name() string {
	return "Type"
}

func (t *Type) Validate(diff PropertyDiff) (bool, error) {
	reset := func(diff PropertyDiff) PropertyDiff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Type = ""
		newProperty.Type = ""
		return NewPropertyDiff(oldProperty, newProperty)
	}

	var err error
	if diff.Old().Type != diff.New().Type {
		err = fmt.Errorf("type changed from %q to %q", diff.Old().Type, diff.New().Type)
	}

	return IsHandled(diff, reset), err
}

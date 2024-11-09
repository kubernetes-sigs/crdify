package property

import (
	"cmp"
	"fmt"
)

func minVerification[T cmp.Ordered](older *T, newer *T) error {
	var err error
	switch {
	case older == nil && newer != nil:
		err = fmt.Errorf("constraint %v added when there were no restrictions previously", *newer)
	case older != nil && newer != nil && *newer > *older:
		err = fmt.Errorf("constraint increased from %v to %v", *older, *newer)
	}
	return err
}

type Minimum struct{}

func (m *Minimum) Name() string {
	return "minimum"
}

func (m *Minimum) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Minimum = nil
		newProperty.Minimum = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := minVerification(diff.Old().Minimum, diff.New().Minimum)
	if err != nil {
		err = fmt.Errorf("minimum: %s", err.Error())
	}

	return IsHandled(diff, reset), err
}

type MinItems struct{}

func (m *MinItems) Name() string {
	return "minItems"
}

func (m *MinItems) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MinItems = nil
		newProperty.MinItems = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := minVerification(diff.Old().MinItems, diff.New().MinItems)
	if err != nil {
		err = fmt.Errorf("minItems: %s", err.Error())
	}

	return IsHandled(diff, reset), err
}

type MinLength struct{}

func (m *MinLength) Name() string {
	return "minLength"
}

func (m *MinLength) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MinLength = nil
		newProperty.MinLength = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := minVerification(diff.Old().MinLength, diff.New().MinLength)
	if err != nil {
		err = fmt.Errorf("minLength: %s", err.Error())
	}

	return IsHandled(diff, reset), err
}

type MinProperties struct{}

func (m *MinProperties) Name() string {
	return "minProperties"
}

func (m *MinProperties) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MinProperties = nil
		newProperty.MinProperties = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := minVerification(diff.Old().MinProperties, diff.New().MinProperties)
	if err != nil {
		err = fmt.Errorf("minProperties: %s", err.Error())
	}

	return IsHandled(diff, reset), err
}

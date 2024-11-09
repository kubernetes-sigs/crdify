package property

import (
	"cmp"
	"fmt"
)

func maxVerification[T cmp.Ordered](older *T, newer *T) error {
	var err error
	switch {
	case older == nil && newer != nil:
		err = fmt.Errorf("constraint %v added when there were no restrictions previously", *newer)
	case older != nil && newer != nil && *newer < *older:
		err = fmt.Errorf("constraint decreased from %v to %v", *older, *newer)
	}
	return err
}

type Maximum struct{}

func (m *Maximum) Name() string {
	return "maximum"
}

func (m *Maximum) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Maximum = nil
		newProperty.Maximum = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := maxVerification(diff.Old().Maximum, diff.New().Maximum)
	if err != nil {
		err = fmt.Errorf("maximum: %s", err.Error())
	}

	return IsHandled(diff, reset), err
}

type MaxItems struct{}

func (m *MaxItems) Name() string {
	return "maxItems"
}

func (m *MaxItems) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MaxItems = nil
		newProperty.MaxItems = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := maxVerification(diff.Old().MaxItems, diff.New().MaxItems)
	if err != nil {
		err = fmt.Errorf("maxItems: %s", err.Error())
	}

	return IsHandled(diff, reset), err
}

type MaxLength struct{}

func (m *MaxLength) Name() string {
	return "maxLength"
}

func (m *MaxLength) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MaxLength = nil
		newProperty.MaxLength = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := maxVerification(diff.Old().MaxLength, diff.New().MaxLength)
	if err != nil {
		err = fmt.Errorf("maxLength: %s", err.Error())
	}

	return IsHandled(diff, reset), err
}

type MaxProperties struct{}

func (m *MaxProperties) Name() string {
	return "maxProperties"
}

func (m *MaxProperties) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MaxProperties = nil
		newProperty.MaxProperties = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := maxVerification(diff.Old().MaxProperties, diff.New().MaxProperties)
	if err != nil {
		err = fmt.Errorf("maxProperties: %s", err.Error())
	}

	return IsHandled(diff, reset), err
}

package property

import (
	"cmp"
	"fmt"
)

type MaxVerificationAdditionEnforcement string

const (
	MaxVerificationAdditionEnforcementStrict = "Strict"
	MaxVerificationAdditionEnforcementNone   = "None"
)

type MaxVerificationDecreaseEnforcement string

const (
	MaxVerificationDecreaseEnforcementStrict = "Strict"
	MaxVerificationDecreaseEnforcementNone   = "None"
)

// MaxOptions is an abstraction for the common
// options for all the "Maximum" related constraints
// on CRD properties.
type MaxOptions struct {
	// AdditionEnforcement is the enforcement strategy to be used when
	// evaluating if adding a new maximum constraint for a CRD property
	// is considered incompatible.
	//
	// Known values are "Strict" and "None".
	//
	// When set to "Strict", adding a new maximum constraint to a CRD property
	// will be considered incompatible. Defaults to "Strict" when
	// unknown values are provided.
	//
	// When set to "None", adding a new maximum constraint to a CRD property
	// will not be considered an incompatible change.
	AdditionEnforcement MaxVerificationAdditionEnforcement `json:"additionEnforcement"`

	// DecreaseEnforcement is the enforcement strategy to be used when
	// evaluating if decreasing the maximum constraint for a CRD property
	// is considered incompatible.
	//
	// Known values are "Strict" and "None".
	//
	// When set to "Strict", decreasing a maximum constraint to a CRD property
	// will be considered incompatible. Defaults to "Strict" when
	// unknown values are provided.
	//
	// When set to "None", decreasing a maximum constraint for a CRD property
	// will not be considered an incompatible change.
	DecreaseEnforcement MaxVerificationDecreaseEnforcement `json:"decreaseEnforcement"`
}

func MaxVerification[T cmp.Ordered](older, newer *T, maxOptions MaxOptions) error {
	var err error
	switch {
	case older == nil && newer != nil && maxOptions.AdditionEnforcement != MaxVerificationAdditionEnforcementNone:
		err = fmt.Errorf("constraint %v added when there were no restrictions previously", *newer)
	case older != nil && newer != nil && *newer < *older && maxOptions.DecreaseEnforcement != MaxVerificationDecreaseEnforcementNone:
		err = fmt.Errorf("constraint decreased from %v to %v", *older, *newer)
	}
	return err
}

type Maximum struct {
	MaxOptions
}

func (m *Maximum) Name() string {
	return "maximum"
}

func (m *Maximum) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Maximum = nil
		newProperty.Maximum = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := MaxVerification(diff.Old().Maximum, diff.New().Maximum, m.MaxOptions)
	if err != nil {
		err = fmt.Errorf("maximum: %s", err.Error())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

type MaxItems struct {
	MaxOptions
}

func (m *MaxItems) Name() string {
	return "maxItems"
}

func (m *MaxItems) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MaxItems = nil
		newProperty.MaxItems = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := MaxVerification(diff.Old().MaxItems, diff.New().MaxItems, m.MaxOptions)
	if err != nil {
		err = fmt.Errorf("maxItems: %s", err.Error())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

type MaxLength struct {
	MaxOptions
}

func (m *MaxLength) Name() string {
	return "maxLength"
}

func (m *MaxLength) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MaxLength = nil
		newProperty.MaxLength = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := MaxVerification(diff.Old().MaxLength, diff.New().MaxLength, m.MaxOptions)
	if err != nil {
		err = fmt.Errorf("maxLength: %s", err.Error())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

type MaxProperties struct {
	MaxOptions
}

func (m *MaxProperties) Name() string {
	return "maxProperties"
}

func (m *MaxProperties) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MaxProperties = nil
		newProperty.MaxProperties = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := MaxVerification(diff.Old().MaxProperties, diff.New().MaxProperties, m.MaxOptions)
	if err != nil {
		err = fmt.Errorf("maxProperties: %s", err.Error())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

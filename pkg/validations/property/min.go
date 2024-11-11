package property

import (
	"cmp"
	"fmt"
)

type MinVerificationAdditionEnforcement string

const (
	MinVerificationAdditionEnforcementStrict = "Strict"
	MinVerificationAdditionEnforcementNone   = "None"
)

type MinVerificationIncreaseEnforcement string

const (
	MinVerificationIncreaseEnforcementStrict = "Strict"
	MinVerificationIncreaseEnforcementNone   = "None"
)

// MinOptions is an abstraction for the common
// options for all the "Minimum" related constraints
// on CRD properties.
type MinOptions struct {
	// AdditionEnforcement is the enforcement strategy to be used when
	// evaluating if adding a new minimum constraint for a CRD property
	// is considered incompatible.
	//
	// Known values are "Strict" and "None".
	//
	// When set to "Strict", adding a new minimum constraint to a CRD property
	// will be considered incompatible. Defaults to "Strict" when
	// unknown values are provided.
	//
	// When set to "None", adding a new minimum constraint to a CRD property
	// will not be considered an incompatible change.
	AdditionEnforcement MinVerificationAdditionEnforcement `json:"additionEnforcement"`

	// IncreaseEnforcement is the enforcement strategy to be used when
	// evaluating if increaseing the minimum constraint for a CRD property
	// is considered incompatible.
	//
	// Known values are "Strict" and "None".
	//
	// When set to "Strict", increasing a minimum constraint to a CRD property
	// will be considered incompatible. Defaults to "Strict" when
	// unknown values are provided.
	//
	// When set to "None", increasing a minimum constraint for a CRD property
	// will not be considered an incompatible change.
	IncreaseEnforcement MinVerificationIncreaseEnforcement `json:"increaseEnforcement"`
}

func MinVerification[T cmp.Ordered](older, newer *T, minOptions MinOptions) error {
	var err error
	switch {
	case older == nil && newer != nil && minOptions.AdditionEnforcement != MinVerificationAdditionEnforcementNone:
		err = fmt.Errorf("constraint %v added when there were no restrictions previously", *newer)
	case older != nil && newer != nil && *newer > *older && minOptions.IncreaseEnforcement != MinVerificationIncreaseEnforcementNone:
		err = fmt.Errorf("constraint increased from %v to %v", *older, *newer)
	}
	return err
}

type Minimum struct {
	MinOptions
}

func (m *Minimum) Name() string {
	return "minimum"
}

func (m *Minimum) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Minimum = nil
		newProperty.Minimum = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := MinVerification(diff.Old().Minimum, diff.New().Minimum, m.MinOptions)
	if err != nil {
		err = fmt.Errorf("minimum: %s", err.Error())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

type MinItems struct {
	MinOptions
}

func (m *MinItems) Name() string {
	return "minItems"
}

func (m *MinItems) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MinItems = nil
		newProperty.MinItems = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := MinVerification(diff.Old().MinItems, diff.New().MinItems, m.MinOptions)
	if err != nil {
		err = fmt.Errorf("minItems: %s", err.Error())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

type MinLength struct {
	MinOptions
}

func (m *MinLength) Name() string {
	return "minLength"
}

func (m *MinLength) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MinLength = nil
		newProperty.MinLength = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := MinVerification(diff.Old().MinLength, diff.New().MinLength, m.MinOptions)
	if err != nil {
		err = fmt.Errorf("minLength: %s", err.Error())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

type MinProperties struct {
	MinOptions
}

func (m *MinProperties) Name() string {
	return "minProperties"
}

func (m *MinProperties) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.MinProperties = nil
		newProperty.MinProperties = nil
		return NewDiff(oldProperty, newProperty)
	}

	err := MinVerification(diff.Old().MinProperties, diff.New().MinProperties, m.MinOptions)
	if err != nil {
		err = fmt.Errorf("minProperties: %s", err.Error())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

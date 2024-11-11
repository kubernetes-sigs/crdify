package property

import (
	"bytes"
	"fmt"
)

type DefaultValidationChangeEnforcement string

const (
	DefaultValidationChangeEnforcementStrict = "Strict"
	DefaultValidationChangeEnforcementNone   = "None"
)

type DefaultValidationRemovalEnforcement string

const (
	DefaultValidationRemovalEnforcementStrict = "Strict"
	DefaultValidationRemovalEnforcementNone   = "None"
)

type DefaultValidationAdditionEnforcement string

const (
	DefaultValidationAdditionEnforcementStrict = "Strict"
	DefaultValidationAdditionEnforcementNone   = "None"
)

// Default is a Validation that can be used to identify
// incompatible changes to the default value of CRD properties
type Default struct {
	// ChangeEnforcement is the enforcement strategy that should be used
	// when evaluating if a change to the default value of a property
	// is considered incompatible.
	//
	// Known enforcement strategies are "Strict" and "None".
	//
	// When set to "Strict", any changes to the default value of a property
	// is considered incompatible.
	//
	// When set to "None", changes to the default value of a property are
	// not considered incompatible.
	//
	// If set to an unknown value, the "Strict" enforcement strategy
	// will be used.
	ChangeEnforcement DefaultValidationChangeEnforcement

	// RemovalEnforcement is the enforcement strategy that should be used
	// when evaluating if the removal of the default value of a property
	// is considered incompatible.
	//
	// Known enforcement strategies are "Strict" and "None".
	//
	// When set to "Strict", removal of the default value of a property
	// is considered incompatible.
	//
	// When set to "None", removal of the default value of a property is
	// not considered incompatible.
	//
	// If set to an unknown value, the "Strict" enforcement strategy
	// will be used.
	RemovalEnforcement DefaultValidationRemovalEnforcement

	// AdditionEnforcement is the enforcement strategy that should be used
	// when evaluating if the addition of a default value for a property
	// is considered incompatible.
	//
	// Known enforcement strategies are "Strict" and "None".
	//
	// When set to "Strict", addition of a default value for a property
	// is considered incompatible.
	//
	// When set to "None", addition of a default value for a property is
	// not considered incompatible.
	//
	// If set to an unknown value, the "Strict" enforcement strategy
	// will be used.
	AdditionEnforcement DefaultValidationAdditionEnforcement
}

func (d *Default) Name() string {
	return "default"
}

func (d *Default) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Default = nil
		newProperty.Default = nil
		return NewDiff(oldProperty, newProperty)
	}
	var err error

	switch {
	case diff.Old().Default == nil && diff.New().Default != nil && d.AdditionEnforcement != DefaultValidationAdditionEnforcementNone:
		err = fmt.Errorf("default value %q added when there was no default previously", string(diff.New().Default.Raw))
	case diff.Old().Default != nil && diff.New().Default == nil && d.RemovalEnforcement != DefaultValidationRemovalEnforcementNone:
		err = fmt.Errorf("default value %q removed", string(diff.Old().Default.Raw))
	case diff.Old().Default != nil && diff.New().Default != nil && !bytes.Equal(diff.Old().Default.Raw, diff.New().Default.Raw) && d.ChangeEnforcement != DefaultValidationChangeEnforcementNone:
		err = fmt.Errorf("default value changed from %q to %q", string(diff.Old().Default.Raw), string(diff.New().Default.Raw))
	}

    resetDiff, handled := IsHandled(diff, reset) 
	return resetDiff, handled, err
}

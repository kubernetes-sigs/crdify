package property

import (
	"fmt"
)

type TypeValidationChangeEnforcement string

const (
	TypeValidationChangeEnforcementStrict = "Strict"
	TypeValidationChangeEnforcementNone   = "None"
)

type Type struct {
	// ChangeEnforcement is the enforcement strategy to be used
	// when evaluating if a change to the type of a CRD property
	// is considered an incompatible change
	//
	// Known values are "Strict" and "None".
	//
	// When set to "Strict", changes to the type of a CRD property
	// are considered incompatible changes. "Strict" is the default
	// enforcement strategy and is used when unknown values are specified.
	//
	// When set to "None", changes to the type of a CRD property
	// are not considered incompatible changes.
	ChangeEnforcement TypeValidationChangeEnforcement `json:"changeEnforcement"`
}

func (t *Type) Name() string {
	return "type"
}

func (t *Type) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Type = ""
		newProperty.Type = ""
		return NewDiff(oldProperty, newProperty)
	}

	var err error
	if diff.Old().Type != diff.New().Type && t.ChangeEnforcement != TypeValidationChangeEnforcementNone {
		err = fmt.Errorf("type changed from %q to %q", diff.Old().Type, diff.New().Type)
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

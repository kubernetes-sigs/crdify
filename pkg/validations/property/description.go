package property

import (
	"fmt"
)

type DescriptionValidationChangeEnforcement string

const (
	DescriptionValidationChangeEnforcementStrict = "Strict"
	DescriptionValidationChangeEnforcementNone   = "None"
)

type DescriptionValidationRemovalEnforcement string

const (
	DescriptionValidationRemovalEnforcementStrict = "Strict"
	DescriptionValidationRemovalEnforcementNone   = "None"
)

type DescriptionValidationAdditionEnforcement string

const (
	DescriptionValidationAdditionEnforcementStrict = "Strict"
	DescriptionValidationAdditionEnforcementNone   = "None"
)

type Description struct {
	// ChangeEnforcement is the enforcement strategy that should be used
	// when evaluating if a change to the description of a property
	// is considered incompatible.
	//
	// Known enforcement strategies are "Strict" and "None".
	//
	// When set to "Strict", any changes to the description of a property
	// is considered incompatible.
	//
	// When set to "None", changes to the description of a property are
	// not considered incompatible.
	//
	// If set to an unknown value, the "Strict" enforcement strategy
	// will be used.
	ChangeEnforcement DescriptionValidationChangeEnforcement

	// RemovalEnforcement is the enforcement strategy that should be used
	// when evaluating if the removal of the description of a property
	// is considered incompatible.
	//
	// Known enforcement strategies are "Strict" and "None".
	//
	// When set to "Strict", removal of the description of a property
	// is considered incompatible.
	//
	// When set to "None", removal of the description of a property is
	// not considered incompatible.
	//
	// If set to an unknown value, the "Strict" enforcement strategy
	// will be used.
	RemovalEnforcement DescriptionValidationRemovalEnforcement

	// AdditionEnforcement is the enforcement strategy that should be used
	// when evaluating if the addition of a description for a property
	// is considered incompatible.
	//
	// Known enforcement strategies are "Strict" and "None".
	//
	// When set to "Strict", addition of a description for a property
	// is considered incompatible.
	//
	// When set to "None", addition of a description for a property is
	// not considered incompatible.
	//
	// If set to an unknown value, the "Strict" enforcement strategy
	// will be used.
	AdditionEnforcement DescriptionValidationAdditionEnforcement
}

func (d *Description) Name() string {
	return "description"
}

func (d *Description) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		newProperty.Description = ""
		oldProperty.Description = ""
		return NewDiff(oldProperty, newProperty)
	}

	oldDescription := diff.Old().Description
	newDescription := diff.New().Description
	var err error

	switch {
	case oldDescription == "" && newDescription != "" && d.AdditionEnforcement != DescriptionValidationAdditionEnforcementNone:
		err = fmt.Errorf("description %q added when there was no description previously", newDescription)
	case oldDescription != "" && newDescription == "" && d.RemovalEnforcement != DescriptionValidationRemovalEnforcementNone:
		err = fmt.Errorf("description %q removed", oldDescription)
	case oldDescription != "" && newDescription != "" && oldDescription != newDescription && d.ChangeEnforcement != DescriptionValidationChangeEnforcementNone:
		err = fmt.Errorf("description changed from %q to %q", oldDescription, newDescription)
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

package property

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
)

type RequiredValidationNewEnforcement string

const (
	RequiredValidationNewEnforcementStrict = "Strict"
	RequiredValidationNewEnforcementNone   = "None"
)

type RequiredValidationRemovalEnforcement string

const (
	RequiredValidationRemovalEnforcementStrict = "Strict"
	RequiredValidationRemovalEnforcementNone   = "None"
)

type Required struct {
	// NewEnforcement is the enforcement strategy to use when
	// evaluating if adding a new required field to a CRD
	// property is considered an incompatible change.
	//
	// Known values are "Strict" and "None".
	//
	// When set to "Strict", adding a new required field
	// will be considered an incompatible change. "Strict"
	// is the default enforcement strategy and is used when
	// unknown values are specified.
	//
	// When set to "None", adding a new required field
	// will not be considered an incompatible change
	NewEnforcement RequiredValidationNewEnforcement `json:"newEnforcement"`

	// RemovalEnforcement is the enforcement strategy to use when
	// evaluating if removing a required field to a CRD
	// property is considered an incompatible change.
	//
	// Known values are "Strict" and "None".
	//
	// When set to "Strict", removing a required field
	// will be considered an incompatible change. "Strict"
	// is the default enforcement strategy and is used when
	// unknown values are specified.
	//
	// When set to "None", removing a required field
	// will not be considered an incompatible change
	RemovalEnforcement RequiredValidationNewEnforcement `json:"removalEnforcement"`
}

func (r *Required) Name() string {
	return "required"
}

func (r *Required) Validate(diff Diff) (Diff, bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		newProperty.Required = []string{}
		oldProperty.Required = []string{}
		return NewDiff(oldProperty, newProperty)
	}

	oldRequired := sets.New(diff.Old().Required...)
	newRequired := sets.New(diff.New().Required...)
	addedRequired := newRequired.Difference(oldRequired)
	removedRequired := oldRequired.Difference(newRequired)
	var err error

	if addedRequired.Len() > 0 && r.NewEnforcement != RequiredValidationNewEnforcementNone {
		err = fmt.Errorf("new required fields %v added", addedRequired.UnsortedList())
	}
	if removedRequired.Len() > 0 && r.NewEnforcement != RequiredValidationRemovalEnforcementNone {
		err = fmt.Errorf("required fields %v removed", removedRequired.UnsortedList())
	}

	resetDiff, handled := IsHandled(diff, reset)
	return resetDiff, handled, err
}

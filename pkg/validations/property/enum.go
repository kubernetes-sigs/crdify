package property

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type EnumValidationRemovalEnforcement string

const (
	EnumValidationRemovalEnforcementStrict = "Strict"
	EnumValidationRemovalEnforcementNone   = "None"
)

type EnumValidationAdditionEnforcement string

const (
	EnumValidationAdditionEnforcementStrict                  = "Strict"
	EnumValidationAdditionEnforcementIfPreviouslyConstrained = "IfPreviouslyConstrained"
	EnumValidationAdditionEnforcementNone                    = "None"
)

// Enum is a Validation that can be used to identify
// incompatible changes to the enum values of CRD properties
type Enum struct {
	// RemovalEnforcement is the enforcement strategy that should be used
	// when evaluating if the removal of an allowed enum value for a property
	// is considered incompatible.
	//
	// Known enforcement strategies are "Strict" and "None".
	//
	// When set to "Strict", removal of an enum value for a property
	// is considered incompatible.
	//
	// When set to "None", removal of an enum value for a property is
	// not considered incompatible.
	//
	// If set to an unknown value, the "Strict" enforcement strategy
	// will be used.
	RemovalEnforcement EnumValidationRemovalEnforcement

	// AdditionEnforcement is the enforcement strategy that should be used
	// when evaluating if the addition of an allowed enum value for a property
	// is considered incompatible.
	//
	// Known enforcement strategies are "Strict", "NotPreviouslyConstrained", and "None".
	//
	// When set to "Strict", addition of any new enum values for a property
	// is considered incompatible, including when there were previously no enum constraints.
	//
	// When set to "IfPreviouslyConstrained", addition any number of enum values for
	// a property when it was not previously constrained by enum values is considered incompatible.
	// Addition of enum values for a property, when it was previously constrained by enum values,
	// is considered compatible with this enforcement strategy
	//
	// When set to "None", addition of a new enum value for a property is
	// not considered incompatible.
	//
	// If set to an unknown value, the "NotPreviouslyConstrained" enforcement strategy
	// will be used.
	AdditionEnforcement EnumValidationAdditionEnforcement
}

func (e *Enum) Name() string {
	return "enum"
}

func (e *Enum) Validate(diff Diff) (bool, error) {
	reset := func(diff Diff) Diff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Enum = []apiextensionsv1.JSON{}
		newProperty.Enum = []apiextensionsv1.JSON{}
		return NewDiff(oldProperty, newProperty)
	}

	oldEnums := sets.New[string]()
	for _, json := range diff.Old().Enum {
		oldEnums.Insert(string(json.Raw))
	}

	newEnums := sets.New[string]()
	for _, json := range diff.New().Enum {
		newEnums.Insert(string(json.Raw))
	}
	removedEnums := oldEnums.Difference(newEnums)
	addedEnums := newEnums.Difference(oldEnums)
	var err error

	switch {
	case oldEnums.Len() == 0 && newEnums.Len() > 0 && e.AdditionEnforcement != EnumValidationAdditionEnforcementNone:
		err = fmt.Errorf("enum constraints %v added when there were no restrictions previously", newEnums.UnsortedList())
	case removedEnums.Len() > 0 && e.RemovalEnforcement != EnumValidationRemovalEnforcementNone:
		err = fmt.Errorf("enums %v removed from the set of previously allowed values", removedEnums.UnsortedList())
	case addedEnums.Len() > 0 && e.AdditionEnforcement == EnumValidationAdditionEnforcementStrict:
		err = fmt.Errorf("enums %v added to the set of allowed values", addedEnums.UnsortedList())
	}

	return IsHandled(diff, reset), err
}

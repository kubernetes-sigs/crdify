package property

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Enum struct{}

func (e *Enum) Name() string {
	return "Enum"
}

func (e *Enum) Validate(diff PropertyDiff) (bool, error) {
	reset := func(diff PropertyDiff) PropertyDiff {
		oldProperty := diff.Old()
		newProperty := diff.New()
		oldProperty.Enum = []apiextensionsv1.JSON{}
		newProperty.Enum = []apiextensionsv1.JSON{}
		return NewPropertyDiff(oldProperty, newProperty)
	}

	oldEnums := sets.New[string]()
	for _, json := range diff.Old().Enum {
		oldEnums.Insert(string(json.Raw))
	}

	newEnums := sets.New[string]()
	for _, json := range diff.New().Enum {
		newEnums.Insert(string(json.Raw))
	}
	diffEnums := oldEnums.Difference(newEnums)
	var err error

	switch {
	case oldEnums.Len() == 0 && newEnums.Len() > 0:
		err = fmt.Errorf("enum constraints %v added when there were no restrictions previously", newEnums.UnsortedList())
	case diffEnums.Len() > 0:
		err = fmt.Errorf("enums %v removed from the set of previously allowed values", diffEnums.UnsortedList())
	}

	return IsHandled(diff, reset), err
}

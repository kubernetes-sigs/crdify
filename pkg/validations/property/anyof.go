// Copyright 2025 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package property

import (
	"errors"
	"fmt"
	"slices"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

const anyOfValidationName = "anyOf"

var (
	_ validations.Validation                                  = (*AnyOf)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*AnyOf)(nil)
)

// RegisterAnyOf registers the AnyOf validation
// with the provided validation registry.
func RegisterAnyOf(registry validations.Registry) {
	registry.Register(anyOfValidationName, anyOfFactory)
}

// anyOfFactory is a function used to initialize an AnyOf validation.
func anyOfFactory(_ map[string]interface{}) (validations.Validation, error) {
	return &AnyOf{}, nil
}

// AnyOf is a Validation that can be used to identify
// incompatible changes to the anyOf constraint of CRD properties.
type AnyOf struct {
	enforcement config.EnforcementPolicy
}

// Name returns the name of the AnyOf validation.
func (a *AnyOf) Name() string {
	return anyOfValidationName
}

// SetEnforcement sets the EnforcementPolicy for the AnyOf validation.
func (a *AnyOf) SetEnforcement(policy config.EnforcementPolicy) {
	a.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the anyOf constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.AnyOf field will be reset to 'nil' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (a *AnyOf) Compare(oldSchema, newSchema *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	oldSchemas := sets.New[string]()

	for i := range oldSchema.AnyOf {
		normalizeSchema(&oldSchema.AnyOf[i])
		oldSchemas.Insert(oldSchema.AnyOf[i].String())
	}

	newSchemas := sets.New[string]()

	for i := range newSchema.AnyOf {
		normalizeSchema(&newSchema.AnyOf[i])
		newSchemas.Insert(newSchema.AnyOf[i].String())
	}

	removedSchemas := oldSchemas.Difference(newSchemas)
	addedSchemas := newSchemas.Difference(oldSchemas)

	var err error

	switch {
	case oldSchemas.Len() == 0 && newSchemas.Len() > 0:
		newSchemaSlice := newSchemas.UnsortedList()
		slices.Sort(newSchemaSlice)
		err = fmt.Errorf("%w: %v", ErrNetNewAnyOfConstraint, newSchemaSlice)
	case removedSchemas.Len() > 0 && addedSchemas.Len() > 0:
		removedSchemaSlice := removedSchemas.UnsortedList()
		slices.Sort(removedSchemaSlice)

		addedSchemaSlice := addedSchemas.UnsortedList()
		slices.Sort(addedSchemaSlice)

		err = fmt.Errorf("%w: removed %v, added %v", ErrChangedAnyOf, removedSchemaSlice, addedSchemaSlice)
	case removedSchemas.Len() > 0:
		removedSchemaSlice := removedSchemas.UnsortedList()
		slices.Sort(removedSchemaSlice)
		err = fmt.Errorf("%w: %v", ErrRemovedAnyOf, removedSchemaSlice)
	}

	oldSchema.AnyOf = nil
	newSchema.AnyOf = nil

	return validations.HandleErrors(a.Name(), a.enforcement, err)
}

var (
	// ErrNetNewAnyOfConstraint represents an error state where a net new anyOf constraint was added to a property.
	ErrNetNewAnyOfConstraint = errors.New("anyOf constraint added when there was none previously")
	// ErrRemovedAnyOf represents an error state where at least one previously allowed anyOf subschema was removed
	// from the anyOf constraint on a property.
	ErrRemovedAnyOf = errors.New("allowed anyOf schemas removed")
	// ErrChangedAnyOf represents an error state where an anyOf subschema was changed.
	ErrChangedAnyOf = errors.New("allowed anyOf schemas changed")
)

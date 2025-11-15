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

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

var (
	_ validations.Validation                                  = (*Not)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*Not)(nil)
)

const notValidationName = "not"

// RegisterNot registers the Not validation
// with the provided validation registry.
func RegisterNot(registry validations.Registry) {
	registry.Register(notValidationName, notFactory)
}

// notFactory is a function used to initialize a Not validation
// implementation based on the provided configuration.
func notFactory(_ map[string]any) (validations.Validation, error) {
	return &Not{}, nil
}

// Not is a Validation that can be used to identify
// incompatible changes to the "not" constraints of CRD properties.
// The "not" constraint means that a property's value cannot match the subschema.
// Adding or changing "not" is more restrictive and breaking.
// Removing "not" is less restrictive and OK.
type Not struct {
	enforcement config.EnforcementPolicy
}

// Name returns the name of the Not validation.
func (n *Not) Name() string {
	return notValidationName
}

// SetEnforcement sets the EnforcementPolicy for the Not validation.
func (n *Not) SetEnforcement(policy config.EnforcementPolicy) {
	n.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the "not" constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.Not field will be reset to 'nil' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (n *Not) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	var err error

	switch {
	case a.Not == nil && b.Not != nil:
		// Adding a "not" constraint is more restrictive (breaking change)
		err = fmt.Errorf("%w", ErrNotConstraintAdded)
	case a.Not != nil && b.Not != nil:
		// Changing a "not" constraint is breaking
		if !equality.Semantic.DeepEqual(a.Not, b.Not) {
			err = fmt.Errorf("%w", ErrNotConstraintChanged)
		}
	case a.Not != nil && b.Not == nil:
		// Removing a "not" constraint is less restrictive (OK, not an error)
	}

	a.Not = nil
	b.Not = nil

	return validations.HandleErrors(n.Name(), n.enforcement, err)
}

var (
	// ErrNotConstraintAdded represents an error state where a "not" constraint was added to a property.
	ErrNotConstraintAdded = errors.New("not constraint added when there was none previously")
	// ErrNotConstraintChanged represents an error state where a "not" constraint was changed.
	ErrNotConstraintChanged = errors.New("not constraint changed")
)
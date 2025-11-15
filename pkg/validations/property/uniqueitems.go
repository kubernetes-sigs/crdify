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
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

var (
	_ validations.Validation                                  = (*UniqueItems)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*UniqueItems)(nil)
)

const uniqueItemsValidationName = "uniqueItems"

// RegisterUniqueItems registers the UniqueItems validation
// with the provided validation registry.
func RegisterUniqueItems(registry validations.Registry) {
	registry.Register(uniqueItemsValidationName, uniqueItemsFactory)
}

// uniqueItemsFactory is a function used to initialize a UniqueItems validation
// implementation based on the provided configuration.
func uniqueItemsFactory(_ map[string]any) (validations.Validation, error) {
	return &UniqueItems{}, nil
}

// UniqueItems is a Validation that can be used to identify
// incompatible changes to the uniqueItems constraints of CRD properties.
// The uniqueItems constraint enforces that lists should only contain unique items.
// Going from non-unique to unique is more restrictive and breaking.
// Going from unique to non-unique is less restrictive and OK.
type UniqueItems struct {
	enforcement config.EnforcementPolicy
}

// Name returns the name of the UniqueItems validation.
func (u *UniqueItems) Name() string {
	return uniqueItemsValidationName
}

// SetEnforcement sets the EnforcementPolicy for the UniqueItems validation.
func (u *UniqueItems) SetEnforcement(policy config.EnforcementPolicy) {
	u.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the uniqueItems constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.UniqueItems field will be reset to 'false' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (u *UniqueItems) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	var err error

	// Going from non-unique (false) to unique (true) is breaking
	if !a.UniqueItems && b.UniqueItems {
		err = fmt.Errorf("%w", ErrUniqueItemsConstraintAdded)
	}
	// Going from unique (true) to non-unique (false) is OK (less restrictive)
	// No error needed for this case

	a.UniqueItems = false
	b.UniqueItems = false

	return validations.HandleErrors(u.Name(), u.enforcement, err)
}

var (
	// ErrUniqueItemsConstraintAdded represents an error state where a uniqueItems constraint was added to a property.
	ErrUniqueItemsConstraintAdded = errors.New("uniqueItems constraint added when there was none previously")
)
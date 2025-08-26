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
	"reflect"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

const anyOfValidationName = "anyOf"

var (
	_ validations.Validation                                  = (*Enum)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*Enum)(nil)
	// ErrNilAnyOfConfig is returned when the AnyOfConfig is nil.
	ErrNilAnyOfConfig = errors.New("AnyOfConfig must not be nil")
	// ErrAtLeastOneSubschema is returned when an unknown AdditionPolicy is provided.
	ErrAtLeastOneSubschema = errors.New("AnyOfConfig must contain at least one subschema")
)

// RegisterAnyOf registers the AnyOf validation
// with the provided validation registry.
func RegisterAnyOf(registry validations.Registry) {
	registry.Register(anyOfValidationName, anyOfFactory)
}

// anyOfFactory is a function used to initialize an Enum validation
// implementation based on the provided configuration.
func anyOfFactory(cfg map[string]interface{}) (validations.Validation, error) {
	anyOfCfg := &AnyOfConfig{}

	err := ConfigToType(cfg, anyOfCfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	err = ValidateAnyOfConfig(anyOfCfg)
	if err != nil {
		return nil, fmt.Errorf("validating anyOf config: %w", err)
	}

	return &AnyOf{AnyOfConfig: *anyOfCfg}, nil
}

// ValidateAnyOfConfig validates the provided AnyOfConfig
// setting default values where appropriate.
// Currently the defaulting behavior defaults the
// AnyOfConfig.AdditionPolicy to AdditionPolicyDisallow
// if it is set to the empty string ("").
func ValidateAnyOfConfig(in *AnyOfConfig) error {
	if in == nil {
		return fmt.Errorf("%w", ErrNilAnyOfConfig)
	}

	if len(in.Subschemas) == 0 {
		return fmt.Errorf("%w", ErrAtLeastOneSubschema)
	}

	switch in.AdditionPolicy {
	case AdditionPolicyAllow, AdditionPolicyDisallow:
		// valid
	case AdditionPolicy(""):
		// default to disallow
		in.AdditionPolicy = AdditionPolicyDisallow
	default:
		return fmt.Errorf("%w: %q", errUnknownAdditionPolicy, in.AdditionPolicy)
	}

	return nil
}

// AnyOfConfig contains additional configurations for the AnyOf validation.
type AnyOfConfig struct {
	// additionPolicy is how adding enums to an existing set of
	// enums should be treated.
	// Allowed values are Allow and Disallow.
	// When set to Allow, adding new values to an existing set
	// of enums will not be flagged.
	// When set to Disallow, adding new values to an existing
	// set of enums will be flagged.
	// Defaults to Disallow.
	AdditionPolicy AdditionPolicy `json:"additionPolicy,omitempty"`
	// Subschemas is a list of subschemas that are part of the anyOf constraint.
	Subschemas []apiextensionsv1.JSONSchemaProps `json:"subschemas,omitempty"`
}

// AnyOf is a validation implementation for the "anyOf" constraint in JSONSchemaProps.
type AnyOf struct {
	// AnyOfConfig is the set of additional configuration options
	AnyOfConfig
	enforcement config.EnforcementPolicy
}

// Name NewAnyOf creates a new AnyOf validation instance with the provided configuration.
func (a *AnyOf) Name() string {
	return anyOfValidationName
}

// SetEnforcement sets the EnforcementPolicy for the Required validation.
func (a *AnyOf) SetEnforcement(policy config.EnforcementPolicy) {
	a.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the format constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.format field will be reset to '""' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (a *AnyOf) Compare(old, newProp *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	oldList := old.AnyOf
	newList := newProp.AnyOf

	var err error

	switch {
	case len(oldList) == 0 && len(newList) > 0:
		err = fmt.Errorf("%w : net-new anyOf constraint added", ErrNetNewAnyOfConstraint)
	case len(oldList) > 0 && len(newList) == 0:
		err = fmt.Errorf("%w : anyOf constraint removed", ErrRemovedAnyOfConstraint)
	case !anyOfEqual(oldList, newList):
		err = fmt.Errorf("%w : anyOf subschemas changed", ErrChangedAnyOfConstraint)
	}

	old.AnyOf = nil
	newProp.AnyOf = nil

	return validations.HandleErrors(a.Name(), a.enforcement, err)
}

// anyOfEqual compares two slices of JSONSchemaProps for deep equality.
func anyOfEqual(a, b []apiextensionsv1.JSONSchemaProps) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !equalJSONSchemaProps(&a[i], &b[i]) {
			return false
		}
	}

	return true
}

// equalJSONSchemaProps performs a deep comparison of two JSONSchemaProps.
// Replace with a proper deep equality check as needed.
func equalJSONSchemaProps(a, b *apiextensionsv1.JSONSchemaProps) bool {
	return reflect.DeepEqual(a, b)
}

var (
	// ErrNetNewAnyOfConstraint represents an error state where a net new anyOf constraint was added to a property.
	ErrNetNewAnyOfConstraint = errors.New("anyOf constraint added when there was none previously")
	// ErrRemovedAnyOfConstraint represents an error state where the anyOf constraint was removed from a property.
	ErrRemovedAnyOfConstraint = errors.New("anyOf constraint removed")
	// ErrChangedAnyOfConstraint represents an error state where the anyOf subschemas changed.
	ErrChangedAnyOfConstraint = errors.New("anyOf subschemas changed")
)

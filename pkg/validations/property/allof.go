// Copyright 2026 The Kubernetes Authors.
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
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

const allOfValidationName = "allOf"

var (
	_ validations.Validation                                  = (*AllOf)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*AllOf)(nil)
)

// RegisterAllOf registers the AllOf validation
// with the provided validation registry.
func RegisterAllOf(registry validations.Registry) {
	registry.Register(allOfValidationName, allOfFactory)
}

// allOfFactory is a function used to initialize an AllOf validation
// implementation based on the provided configuration.
func allOfFactory(cfg map[string]interface{}) (validations.Validation, error) {
	allOfCfg := &AllOfConfig{}

	err := ConfigToType(cfg, allOfCfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	err = ValidateAllOfConfig(allOfCfg)
	if err != nil {
		return nil, fmt.Errorf("validating allOf config: %w", err)
	}

	return &AllOf{AllOfConfig: *allOfCfg}, nil
}

// ValidateAllOfConfig validates the provided AllOfConfig
// setting default values where appropriate.
func ValidateAllOfConfig(in *AllOfConfig) error {
	if in == nil {
		return nil
	}

	switch in.AdditionPolicy {
	case AdditionPolicyAllow, AdditionPolicyDisallow:
		// do nothing, valid case
	case AdditionPolicy(""):
		in.AdditionPolicy = AdditionPolicyDisallow
	default:
		return fmt.Errorf("%w : %q", errUnknownAdditionPolicy, in.AdditionPolicy)
	}

	switch in.RemovalPolicy {
	case RemovalPolicyAllow, RemovalPolicyDisallow:
		// do nothing, valid case
	case RemovalPolicy(""):
		in.RemovalPolicy = RemovalPolicyDisallow
	default:
		return fmt.Errorf("%w : %q", errUnknownRemovalPolicy, in.RemovalPolicy)
	}

	return nil
}

// AllOfConfig contains additional configurations for the AllOf validation.
type AllOfConfig struct {
	// additionPolicy is how adding subschemas to an allOf
	// constraint should be treated, including adding an allOf
	// constraint where there was none previously.
	// Allowed values are Allow and Disallow.
	// When set to Allow, adding new subschemas will not be flagged.
	// When set to Disallow, adding new subschemas will be flagged.
	// Defaults to Disallow.
	AdditionPolicy AdditionPolicy `json:"additionPolicy,omitempty"`
	// removalPolicy is how removing subschemas from an existing allOf
	// constraint should be treated.
	// Allowed values are Allow and Disallow.
	// When set to Allow, removing subschemas will not be flagged.
	// When set to Disallow, removing subschemas will be flagged.
	// Defaults to Disallow.
	RemovalPolicy RemovalPolicy `json:"removalPolicy,omitempty"`
}

// AllOf is a Validation that can be used to identify
// incompatible changes to the allOf constraint of CRD properties.
type AllOf struct {
	AllOfConfig

	enforcement config.EnforcementPolicy
}

// Name returns the name of the AllOf validation.
func (ao *AllOf) Name() string {
	return allOfValidationName
}

// SetEnforcement sets the EnforcementPolicy for the AllOf validation.
func (ao *AllOf) SetEnforcement(policy config.EnforcementPolicy) {
	ao.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the allOf constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.AllOf field will be reset to 'nil' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (ao *AllOf) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	oldSchemas, oldErr := schemaKeys(a.AllOf)
	newSchemas, newErr := schemaKeys(b.AllOf)

	if err := errors.Join(oldErr, newErr); err != nil {
		a.AllOf = nil
		b.AllOf = nil

		return validations.HandleErrors(ao.Name(), ao.enforcement, err)
	}

	removedSchemas := oldSchemas.Difference(newSchemas)
	addedSchemas := newSchemas.Difference(oldSchemas)

	var errs []error

	if removedSchemas.Len() > 0 && ao.RemovalPolicy != RemovalPolicyAllow {
		removedSchemaSlice := removedSchemas.UnsortedList()
		slices.Sort(removedSchemaSlice)
		errs = append(errs, fmt.Errorf("%w: %v", ErrRemovedAllOf, removedSchemaSlice))
	}

	if addedSchemas.Len() > 0 && ao.AdditionPolicy != AdditionPolicyAllow {
		addedSchemaSlice := addedSchemas.UnsortedList()
		slices.Sort(addedSchemaSlice)
		errs = append(errs, fmt.Errorf("%w: %v", ErrAddedAllOf, addedSchemaSlice))
	}

	a.AllOf = nil
	b.AllOf = nil

	return validations.HandleErrors(ao.Name(), ao.enforcement, errs...)
}

var (
	// ErrRemovedAllOf represents an error state where at least one subschema was removed
	// from the allOf constraint on a property, loosening validation.
	ErrRemovedAllOf = errors.New("allOf schemas removed")
	// ErrAddedAllOf represents an error state where at least one subschema was added
	// to the allOf constraint on a property, tightening validation.
	ErrAddedAllOf = errors.New("allOf schemas added")
)

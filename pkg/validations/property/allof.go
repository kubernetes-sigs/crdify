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
	"k8s.io/apimachinery/pkg/util/sets"
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

// ValidateAllOfConfig ensures provided AllOfConfig is valid and defaults missing values.
func ValidateAllOfConfig(in *AllOfConfig) error {
	if in == nil {
		return nil
	}

	switch in.AdditionPolicy {
	case AllOfAdditionPolicyAllow, AllOfAdditionPolicyDisallow:
	case AllOfAdditionPolicy(""):
		in.AdditionPolicy = AllOfAdditionPolicyDisallow
	default:
		return fmt.Errorf("%w : %q (valid values: %q, %q)", errUnknownAllOfAdditionPolicy, in.AdditionPolicy, AllOfAdditionPolicyAllow, AllOfAdditionPolicyDisallow)
	}

	switch in.RemovalPolicy {
	case AllOfRemovalPolicyAllow, AllOfRemovalPolicyDisallow:
	case AllOfRemovalPolicy(""):
		in.RemovalPolicy = AllOfRemovalPolicyDisallow
	default:
		return fmt.Errorf("%w : %q (valid values: %q, %q)", errUnknownAllOfRemovalPolicy, in.RemovalPolicy, AllOfRemovalPolicyAllow, AllOfRemovalPolicyDisallow)
	}

	return nil
}

var (
	errUnknownAllOfAdditionPolicy = errors.New("unknown addition policy")
	errUnknownAllOfRemovalPolicy  = errors.New("unknown removal policy")
)

// AllOfAdditionPolicy represents how adding allOf subschemas should be evaluated.
type AllOfAdditionPolicy string

const (
	// AllOfAdditionPolicyAllow treats adding allOf subschemas as compatible.
	AllOfAdditionPolicyAllow AllOfAdditionPolicy = "Allow"
	// AllOfAdditionPolicyDisallow treats adding allOf subschemas as incompatible.
	AllOfAdditionPolicyDisallow AllOfAdditionPolicy = "Disallow"
)

// AllOfRemovalPolicy represents how removing allOf subschemas should be evaluated.
type AllOfRemovalPolicy string

const (
	// AllOfRemovalPolicyAllow treats removing allOf subschemas as compatible.
	AllOfRemovalPolicyAllow AllOfRemovalPolicy = "Allow"
	// AllOfRemovalPolicyDisallow treats removing allOf subschemas as incompatible.
	AllOfRemovalPolicyDisallow AllOfRemovalPolicy = "Disallow"
)

// AllOfConfig contains additional configuration for the AllOf validation.
type AllOfConfig struct {
	// AdditionPolicy dictates whether adding allOf subschemas is compatible.
	// Allowed values are Allow and Disallow. Defaults to Disallow.
	AdditionPolicy AllOfAdditionPolicy `json:"additionPolicy,omitempty"`
	// RemovalPolicy dictates whether removing allOf subschemas is compatible.
	// Allowed values are Allow and Disallow. Defaults to Disallow.
	RemovalPolicy AllOfRemovalPolicy `json:"removalPolicy,omitempty"`
}

// AllOf is a Validation that can be used to identify
// incompatible changes to the allOf constraints of CRD properties.
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
	oldSchemas := sets.New[string]()

	for i := range a.AllOf {
		normalizeAllOfSchema(&a.AllOf[i])
		oldSchemas.Insert(a.AllOf[i].String())
	}

	newSchemas := sets.New[string]()

	for i := range b.AllOf {
		normalizeAllOfSchema(&b.AllOf[i])
		newSchemas.Insert(b.AllOf[i].String())
	}

	removedSchemas := oldSchemas.Difference(newSchemas)
	addedSchemas := newSchemas.Difference(oldSchemas)

	var err error

	switch {
	case oldSchemas.Len() == 0 && newSchemas.Len() > 0:
		newSchemaSlice := newSchemas.UnsortedList()
		slices.Sort(newSchemaSlice)
		err = fmt.Errorf("%w : %v", ErrNetNewAllOfConstraint, newSchemaSlice)
	case removedSchemas.Len() > 0 && ao.RemovalPolicy != AllOfRemovalPolicyAllow:
		removedSchemaSlice := removedSchemas.UnsortedList()
		slices.Sort(removedSchemaSlice)
		err = fmt.Errorf("%w : %v", ErrAllOfConstraintRemoved, removedSchemaSlice)
	case addedSchemas.Len() > 0 && ao.AdditionPolicy != AllOfAdditionPolicyAllow:
		addedSchemaSlice := addedSchemas.UnsortedList()
		slices.Sort(addedSchemaSlice)
		err = fmt.Errorf("%w : %v", ErrAllOfConstraintAdded, addedSchemaSlice)
	}

	a.AllOf = nil
	b.AllOf = nil

	return validations.HandleErrors(ao.Name(), ao.enforcement, err)
}

// normalizeAllOfSchema zeroes non-structural fields on a schema
// so that only validation-relevant differences are compared.
func normalizeAllOfSchema(schema *apiextensionsv1.JSONSchemaProps) {
	schema.Description = ""
	schema.Title = ""
	schema.Example = nil
	schema.ExternalDocs = nil
}

var (
	// ErrNetNewAllOfConstraint represents an error state where a net new allOf constraint was added to a property.
	ErrNetNewAllOfConstraint = errors.New("allOf constraint added when there was none previously")
	// ErrAllOfConstraintAdded represents an error state where new subschemas were added to an existing allOf constraint.
	ErrAllOfConstraintAdded = errors.New("allOf subschema(s) added")
	// ErrAllOfConstraintRemoved represents an error state where subschemas were removed from an existing allOf constraint.
	ErrAllOfConstraintRemoved = errors.New("allOf subschema(s) removed")
)

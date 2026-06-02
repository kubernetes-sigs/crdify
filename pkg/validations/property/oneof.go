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

const oneOfValidationName = "oneOf"

var (
	_ validations.Validation                                  = (*OneOf)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*OneOf)(nil)
)

// RegisterOneOf registers the OneOf validation
// with the provided validation registry.
func RegisterOneOf(registry validations.Registry) {
	registry.Register(oneOfValidationName, oneOfFactory)
}

// oneOfFactory is a function used to initialize a OneOf validation
// implementation based on the provided configuration.
func oneOfFactory(cfg map[string]interface{}) (validations.Validation, error) {
	oneOfCfg := &OneOfConfig{}

	err := ConfigToType(cfg, oneOfCfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	err = ValidateOneOfConfig(oneOfCfg)
	if err != nil {
		return nil, fmt.Errorf("validating oneOf config: %w", err)
	}

	return &OneOf{OneOfConfig: *oneOfCfg}, nil
}

// ValidateOneOfConfig validates the provided OneOfConfig
// setting default values where appropriate.
func ValidateOneOfConfig(in *OneOfConfig) error {
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

// RemovalPolicy is used to represent how a validation
// should determine compatibility of removing existing constraints.
type RemovalPolicy string

const (
	// RemovalPolicyAllow signals that removing an existing constraint
	// should be considered a compatible change.
	RemovalPolicyAllow RemovalPolicy = "Allow"

	// RemovalPolicyDisallow signals that removing an existing constraint
	// should be considered an incompatible change.
	RemovalPolicyDisallow RemovalPolicy = "Disallow"
)

var errUnknownRemovalPolicy = errors.New("unknown removal policy")

// OneOfConfig contains additional configurations for the OneOf validation.
type OneOfConfig struct {
	// additionPolicy is how adding subschemas to an existing oneOf
	// constraint should be treated.
	// Allowed values are Allow and Disallow.
	// When set to Allow, adding new subschemas to an existing
	// oneOf constraint will not be flagged.
	// When set to Disallow, adding new subschemas to an existing
	// oneOf constraint will be flagged.
	// Defaults to Disallow.
	AdditionPolicy AdditionPolicy `json:"additionPolicy,omitempty"`
	// removalPolicy is how removing subschemas from an existing oneOf
	// constraint should be treated.
	// Allowed values are Allow and Disallow.
	// When set to Allow, removing subschemas from an existing
	// oneOf constraint will not be flagged.
	// When set to Disallow, removing subschemas from an existing
	// oneOf constraint will be flagged.
	// Defaults to Disallow.
	RemovalPolicy RemovalPolicy `json:"removalPolicy,omitempty"`
}

// OneOf is a Validation that can be used to identify
// incompatible changes to the oneOf constraint of CRD properties.
type OneOf struct {
	OneOfConfig

	enforcement config.EnforcementPolicy
}

// Name returns the name of the OneOf validation.
func (o *OneOf) Name() string {
	return oneOfValidationName
}

// SetEnforcement sets the EnforcementPolicy for the OneOf validation.
func (o *OneOf) SetEnforcement(policy config.EnforcementPolicy) {
	o.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the oneOf constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.OneOf field will be reset to 'nil' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (o *OneOf) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	oldSchemas := sets.New[string]()

	for i := range a.OneOf {
		normalizeSchema(&a.OneOf[i])
		oldSchemas.Insert(a.OneOf[i].String())
	}

	newSchemas := sets.New[string]()

	for i := range b.OneOf {
		normalizeSchema(&b.OneOf[i])
		newSchemas.Insert(b.OneOf[i].String())
	}

	removedSchemas := oldSchemas.Difference(newSchemas)
	addedSchemas := newSchemas.Difference(oldSchemas)

	var err error

	switch {
	case oldSchemas.Len() == 0 && newSchemas.Len() > 0:
		err = o.checkNetNewOneOf(a, newSchemas)
	case removedSchemas.Len() > 0 && o.RemovalPolicy != RemovalPolicyAllow:
		removedSchemaSlice := removedSchemas.UnsortedList()
		slices.Sort(removedSchemaSlice)
		err = fmt.Errorf("%w: %v", ErrRemovedOneOf, removedSchemaSlice)
	case addedSchemas.Len() > 0 && o.AdditionPolicy != AdditionPolicyAllow:
		addedSchemaSlice := addedSchemas.UnsortedList()
		slices.Sort(addedSchemaSlice)
		err = fmt.Errorf("%w: %v", ErrAddedOneOf, addedSchemaSlice)
	}

	a.OneOf = nil
	b.OneOf = nil

	return validations.HandleErrors(o.Name(), o.enforcement, err)
}

// checkNetNewOneOf evaluates whether a net-new oneOf constraint is
// an incompatible change. When additionPolicy is Allow and the
// pre-existing property schema is preserved as one of the new oneOf
// entries, the change is considered compatible.
func (o *OneOf) checkNetNewOneOf(old *apiextensionsv1.JSONSchemaProps, newSchemas sets.Set[string]) error {
	newSchemaSlice := newSchemas.UnsortedList()
	slices.Sort(newSchemaSlice)

	if o.AdditionPolicy != AdditionPolicyAllow {
		return fmt.Errorf("%w: %v", ErrNetNewOneOfConstraint, newSchemaSlice)
	}

	if !preExistingSchemaPreserved(old, newSchemas) {
		return fmt.Errorf("%w: %v", ErrNetNewOneOfPreExistingSchemaNotPreserved, newSchemaSlice)
	}

	return nil
}

// preExistingSchemaPreserved checks whether the old property schema
// (with its OneOf field cleared) appears as one of the entries in the
// new oneOf constraint. This is used to determine if a net-new oneOf
// is a safe loosening change for writer-focused APIs.
func preExistingSchemaPreserved(old *apiextensionsv1.JSONSchemaProps, newSchemas sets.Set[string]) bool {
	oldCopy := old.DeepCopy()
	oldCopy.OneOf = nil
	normalizeSchema(oldCopy)

	return newSchemas.Has(oldCopy.String())
}

// normalizeSchema zeroes non-structural fields on a schema
// so that only validation-relevant differences are compared.
func normalizeSchema(schema *apiextensionsv1.JSONSchemaProps) {
	schema.Description = ""
	schema.Title = ""
	schema.Example = nil
	schema.ExternalDocs = nil
}

var (
	// ErrNetNewOneOfConstraint represents an error state where a net new oneOf constraint was added to a property.
	ErrNetNewOneOfConstraint = errors.New("oneOf constraint added when there was none previously")
	// ErrNetNewOneOfPreExistingSchemaNotPreserved represents an error state where a net new oneOf constraint was added
	// but the pre-existing property schema is not preserved as one of the new oneOf entries, breaking both readers and writers.
	ErrNetNewOneOfPreExistingSchemaNotPreserved = errors.New("oneOf constraint added and pre-existing property schema is not preserved in the new oneOf entries")
	// ErrRemovedOneOf represents an error state where at least one previously allowed oneOf subschema was removed
	// from the oneOf constraint on a property.
	ErrRemovedOneOf = errors.New("allowed oneOf schemas removed")
	// ErrAddedOneOf represents an error state where at least one oneOf subschema, that was not previously allowed,
	// was added to the oneOf constraint on a property.
	ErrAddedOneOf = errors.New("allowed oneOf schemas added")
)

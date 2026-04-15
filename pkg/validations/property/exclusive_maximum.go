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

//nolint:dupl
package property

import (
	"errors"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

var (
	_ validations.Validation                                  = (*ExclusiveMaximum)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*ExclusiveMaximum)(nil)
)

const exclusiveMaximumValidationName = "exclusiveMaximum"

// RegisterExclusiveMaximum registers the ExclusiveMaximum validation
// with the provided validation registry.
func RegisterExclusiveMaximum(registry validations.Registry) {
	registry.Register(exclusiveMaximumValidationName, exclusiveMaximumFactory)
}

// exclusiveMaximumFactory is a function used to initialize an ExclusiveMaximum validation
// implementation based on the provided configuration.
func exclusiveMaximumFactory(cfg map[string]interface{}) (validations.Validation, error) {
	exclusiveCfg := &ExclusiveMaximumConfig{}

	err := ConfigToType(cfg, exclusiveCfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	err = ValidateExclusiveMaximumConfig(exclusiveCfg)
	if err != nil {
		return nil, fmt.Errorf("validating exclusiveMaximum config: %w", err)
	}

	return &ExclusiveMaximum{ExclusiveMaximumConfig: *exclusiveCfg}, nil
}

// ValidateExclusiveMaximumConfig ensures provided ExclusiveMaximumConfig is valid and defaults missing values.
func ValidateExclusiveMaximumConfig(in *ExclusiveMaximumConfig) error {
	if in == nil {
		return nil
	}

	switch in.AdditionPolicy {
	case ExclusiveMaximumAdditionPolicyAllow, ExclusiveMaximumAdditionPolicyDisallow:
		// valid entries
	case ExclusiveMaximumAdditionPolicy(""):
		in.AdditionPolicy = ExclusiveMaximumAdditionPolicyDisallow
	default:
		return fmt.Errorf("%w : %q (valid values: %q, %q)", errUnknownExclusiveMaximumAdditionPolicy, in.AdditionPolicy, ExclusiveMaximumAdditionPolicyAllow, ExclusiveMaximumAdditionPolicyDisallow)
	}

	switch in.RemovalPolicy {
	case ExclusiveMaximumRemovalPolicyAllow, ExclusiveMaximumRemovalPolicyDisallow:
		// valid entries
	case ExclusiveMaximumRemovalPolicy(""):
		in.RemovalPolicy = ExclusiveMaximumRemovalPolicyDisallow
	default:
		return fmt.Errorf("%w : %q (valid values: %q, %q)", errUnknownExclusiveMaximumRemovalPolicy, in.RemovalPolicy, ExclusiveMaximumRemovalPolicyAllow, ExclusiveMaximumRemovalPolicyDisallow)
	}

	return nil
}

var errUnknownExclusiveMaximumAdditionPolicy = errors.New("unknown addition policy")
var errUnknownExclusiveMaximumRemovalPolicy = errors.New("unknown removal policy")

// ExclusiveMaximumAdditionPolicy represents how adding the exclusiveMaximum constraint should be evaluated.
type ExclusiveMaximumAdditionPolicy string

const (
	// ExclusiveMaximumAdditionPolicyAllow treats adding exclusiveMaximum when it was previously absent as compatible.
	ExclusiveMaximumAdditionPolicyAllow ExclusiveMaximumAdditionPolicy = "Allow"
	// ExclusiveMaximumAdditionPolicyDisallow treats adding exclusiveMaximum when it was previously absent as incompatible.
	ExclusiveMaximumAdditionPolicyDisallow ExclusiveMaximumAdditionPolicy = "Disallow"
)

// ExclusiveMaximumRemovalPolicy represents how loosening the exclusiveMaximum constraint should be evaluated.
type ExclusiveMaximumRemovalPolicy string

const (
	// ExclusiveMaximumRemovalPolicyAllow treats loosening exclusiveMaximum as compatible.
	ExclusiveMaximumRemovalPolicyAllow ExclusiveMaximumRemovalPolicy = "Allow"
	// ExclusiveMaximumRemovalPolicyDisallow treats loosening exclusiveMaximum as incompatible.
	ExclusiveMaximumRemovalPolicyDisallow ExclusiveMaximumRemovalPolicy = "Disallow"
)

// ExclusiveMaximumConfig contains additional configuration for the ExclusiveMaximum validation.
type ExclusiveMaximumConfig struct {
	// AdditionPolicy dictates whether adding exclusiveMaximum when it was previously absent is compatible.
	// Allowed values are Allow and Disallow. Defaults to Disallow.
	AdditionPolicy ExclusiveMaximumAdditionPolicy `json:"additionPolicy,omitempty"`
	// RemovalPolicy dictates whether loosening exclusiveMaximum is compatible.
	// Allowed values are Allow and Disallow. Defaults to Disallow.
	RemovalPolicy ExclusiveMaximumRemovalPolicy `json:"removalPolicy,omitempty"`
}

// ExclusiveMaximum is a Validation that can be used to identify
// incompatible changes to the exclusiveMaximum constraint of CRD properties.
type ExclusiveMaximum struct {
	ExclusiveMaximumConfig
	enforcement config.EnforcementPolicy
}

// Name returns the name of the ExclusiveMaximum validation.
func (e *ExclusiveMaximum) Name() string {
	return exclusiveMaximumValidationName
}

// SetEnforcement sets the EnforcementPolicy for the ExclusiveMaximum validation.
func (e *ExclusiveMaximum) SetEnforcement(policy config.EnforcementPolicy) {
	e.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the exclusiveMaximum constraint of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.ExclusiveMaximum field will be reset to 'false' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (e *ExclusiveMaximum) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	var err error

	switch {
	case a.ExclusiveMaximum == b.ExclusiveMaximum:
		// nothing to do
	case !a.ExclusiveMaximum && b.ExclusiveMaximum && e.AdditionPolicy != ExclusiveMaximumAdditionPolicyAllow:
		err = fmt.Errorf("%w : %t -> %t", ErrExclusiveMaximumActivated, a.ExclusiveMaximum, b.ExclusiveMaximum)
	case a.ExclusiveMaximum && !b.ExclusiveMaximum && e.RemovalPolicy != ExclusiveMaximumRemovalPolicyAllow:
		err = fmt.Errorf("%w : %t -> %t", ErrExclusiveMaximumRemoved, a.ExclusiveMaximum, b.ExclusiveMaximum)
	}

	a.ExclusiveMaximum = false
	b.ExclusiveMaximum = false

	return validations.HandleErrors(e.Name(), e.enforcement, err)
}

// ErrExclusiveMaximumActivated represents an error state when a property transitions from inclusive to exclusive maximum.
var ErrExclusiveMaximumActivated = errors.New("exclusive maximum activated when it was not previously")

// ErrExclusiveMaximumRemoved represents an error state when a property transitions from exclusive to inclusive maximum.
var ErrExclusiveMaximumRemoved = errors.New("exclusive maximum removed when it was not previously")

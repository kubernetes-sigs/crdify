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

const nullableValidationName = "nullable"

var (
	_ validations.Validation                                  = (*Nullable)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*Nullable)(nil)
)

// RegisterNullable registers the Nullable validation
// with the provided validation registry.
func RegisterNullable(registry validations.Registry) {
	registry.Register(nullableValidationName, nullableFactory)
}

// nullableFactory is a function used to initialize a Nullable validation
// implementation based on the provided configuration.
func nullableFactory(cfg map[string]interface{}) (validations.Validation, error) {
	nullableCfg := &NullableConfig{}

	err := ConfigToType(cfg, nullableCfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	err = ValidateNullableConfig(nullableCfg)
	if err != nil {
		return nil, fmt.Errorf("validating nullable config: %w", err)
	}

	return &Nullable{NullableConfig: *nullableCfg}, nil
}

// ValidateNullableConfig validates the provided NullableConfig
// setting default values where appropriate.
// Currently the defaulting behavior defaults the
// NullableConfig.ToNullablePolicy to ToNullablePolicyDisallow
// if it is set to the empty string ("").
func ValidateNullableConfig(in *NullableConfig) error {
	if in == nil {
		// nothing to validate
		return nil
	}

	switch in.ToNullablePolicy {
	case ToNullablePolicyAllow, ToNullablePolicyDisallow:
		// do nothing, valid case
	case ToNullablePolicy(""):
		// default to disallow
		in.ToNullablePolicy = ToNullablePolicyDisallow
	default:
		return fmt.Errorf("%w : %q", errUnknownToNullablePolicy, in.ToNullablePolicy)
	}

	return nil
}

var errUnknownToNullablePolicy = errors.New("unknown to nullable policy")

// ToNullablePolicy is used to represent how the Nullable validation
// should determine compatibility of changing from not-nullable to nullable.
type ToNullablePolicy string

const (
	// ToNullablePolicyAllow signals that changing from not-nullable to nullable
	// should be considered a compatible change.
	ToNullablePolicyAllow ToNullablePolicy = "Allow"

	// ToNullablePolicyDisallow signals that changing from not-nullable to nullable
	// should be considered an incompatible change.
	ToNullablePolicyDisallow ToNullablePolicy = "Disallow"
)

// NullableConfig contains additional configurations for the Nullable validation.
type NullableConfig struct {
	// ToNullablePolicy controls how changing from non-nullable to nullable is treated.
	// Allowed values are Allow and Disallow.
	// When set to Allow, changing from not-nullable to nullable will not be flagged.
	// When set to Disallow, changing from not-nullable to nullable will be flagged.
	// Defaults to Disallow.
	ToNullablePolicy ToNullablePolicy `json:"toNullablePolicy,omitempty"`
}

// Nullable is a Validation that can be used to identify
// incompatible changes to the nullable constraint of a property.
type Nullable struct {
	// NullableConfig is the set of additional configuration options
	// for the Nullable validation.
	NullableConfig

	// enforcement is the EnforcementPolicy that this validation
	// should use when performing its validation logic
	enforcement config.EnforcementPolicy
}

// Name returns the name of the Nullable validation.
func (n *Nullable) Name() string {
	return nullableValidationName
}

// SetEnforcement sets the EnforcementPolicy for the Nullable validation.
func (n *Nullable) SetEnforcement(policy config.EnforcementPolicy) {
	n.enforcement = policy
}

// Compare compares the nullable constraint of two JSONSchemaProps objects.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.Nullable field will be reset to 'false' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (n *Nullable) Compare(oldSchema, newSchema *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	var err error

	// Check for incompatible changes to nullable constraint
	switch {
	case !oldSchema.Nullable && newSchema.Nullable:
		if n.ToNullablePolicy == ToNullablePolicyDisallow {
			err = ErrMadeNullable
		}
	case oldSchema.Nullable && !newSchema.Nullable:
		err = ErrMadeNonNullable
	}

	// Reset fields to indicate they've been handled
	oldSchema.Nullable = false
	newSchema.Nullable = false

	return validations.HandleErrors(n.Name(), n.enforcement, err)
}

var (
	// ErrMadeNonNullable represents an error state where a property was changed
	// from nullable to non-nullable, which is a breaking change.
	ErrMadeNonNullable = errors.New("property changed from nullable to non-nullable")
	// ErrMadeNullable represents an error state where a property was changed
	// from non-nullable to nullable, which is a breaking change.
	ErrMadeNullable = errors.New("property changed from non-nullable to nullable")
)

// Copyright The Kubernetes Authors.
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
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/set"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

var (
	_ validations.Validation                                  = (*XValidations)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*XValidations)(nil)
)

const validationName = "xvalidations"

// RegisterXValidations registers the XValidations validation
// with the provided validation registry.
func RegisterXValidations(registry validations.Registry) {
	registry.Register(validationName, factoryXValidations)
}

// factoryXValidations is a function used to initialize an XValidations validation
// implementation based on the provided configuration.
func factoryXValidations(cfg map[string]interface{}) (validations.Validation, error) {
	xvCfg := &XValidationsConfig{}

	if err := ConfigToType(cfg, xvCfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	err := ValidateXValidationsConfig(xvCfg)
	if err != nil {
		return nil, fmt.Errorf("validating XValidations config: %w", err)
	}

	return &XValidations{XValidationsConfig: *xvCfg}, nil
}

// ValidateXValidationsConfig ensures provided XValidationsConfig is valid and defaults missing values.
func ValidateXValidationsConfig(in *XValidationsConfig) error {
	if in == nil {
		return nil
	}

	switch in.AdditionPolicy {
	case XValidationsAdditionPolicyAllow, XValidationsAdditionPolicyDisallow:
		// valid entries
	case XValidationsAdditionPolicy(""):
		in.AdditionPolicy = XValidationsAdditionPolicyDisallow
	default:
		return fmt.Errorf("%w : %q (valid values: %q, %q)", errUnknownXValidationsAdditionPolicy, in.AdditionPolicy, XValidationsAdditionPolicyAllow, XValidationsAdditionPolicyDisallow)
	}

	switch in.RemovalPolicy {
	case XValidationsRemovalPolicyAllow, XValidationsRemovalPolicyDisallow:
		// valid entries
	case XValidationsRemovalPolicy(""):
		in.RemovalPolicy = XValidationsRemovalPolicyDisallow
	default:
		return fmt.Errorf("%w : %q (valid values: %q, %q)", errUnknownXValidationsRemovalPolicy, in.RemovalPolicy, XValidationsRemovalPolicyAllow, XValidationsRemovalPolicyDisallow)
	}

	return nil
}

var errUnknownXValidationsAdditionPolicy = errors.New("unknown addition policy")
var errUnknownXValidationsRemovalPolicy = errors.New("unknown removal policy")

// XValidationsAdditionPolicy represents how adding a new x-kubernetes-validation rule should be evaluated.
type XValidationsAdditionPolicy string

const (
	// XValidationsAdditionPolicyAllow treats adding a new x-kubernetes-validation rule as compatible.
	XValidationsAdditionPolicyAllow XValidationsAdditionPolicy = "Allow"
	// XValidationsAdditionPolicyDisallow treats adding a new x-kubernetes-validation rule as incompatible.
	XValidationsAdditionPolicyDisallow XValidationsAdditionPolicy = "Disallow"
)

// XValidationsRemovalPolicy represents how removing an existing x-kubernetes-validation rule should be evaluated.
type XValidationsRemovalPolicy string

const (
	// XValidationsRemovalPolicyAllow treats removing an existing x-kubernetes-validation rule as compatible.
	XValidationsRemovalPolicyAllow XValidationsRemovalPolicy = "Allow"
	// XValidationsRemovalPolicyDisallow treats removing an existing x-kubernetes-validation rule as incompatible.
	XValidationsRemovalPolicyDisallow XValidationsRemovalPolicy = "Disallow"
)

// XValidationsConfig contains additional configuration for the XValidations validation.
type XValidationsConfig struct {
	// AdditionPolicy dictates whether adding a new x-kubernetes-validation rule is compatible.
	// Allowed values are Allow and Disallow. Defaults to Disallow.
	AdditionPolicy XValidationsAdditionPolicy `json:"additionPolicy,omitempty"`
	// RemovalPolicy dictates whether removing an existing x-kubernetes-validation rule is compatible.
	// Allowed values are Allow and Disallow. Defaults to Disallow.
	RemovalPolicy XValidationsRemovalPolicy `json:"removalPolicy,omitempty"`
}

// XValidations is a Validation that can be used to identify
// incompatible changes to the x-kubernetes-validations of CRD properties.
type XValidations struct {
	XValidationsConfig
	enforcement config.EnforcementPolicy
}

// Name returns the name of the XValidations validation.
func (x *XValidations) Name() string {
	return validationName
}

// SetEnforcement sets the EnforcementPolicy for the XValidations validation.
func (x *XValidations) SetEnforcement(policy config.EnforcementPolicy) {
	x.enforcement = policy
}

// hashValidationRule returns a stable SHA-256 hash of the semantically
// significant fields of a ValidationRule. Only Rule, Reason, and OptionalOldSelf are included.
func hashValidationRule(validationRule apiextensionsv1.ValidationRule) string {
	strToHash := fmt.Sprintf("%v", validationRule.Rule)

	if validationRule.Reason != nil {
		strToHash += fmt.Sprintf("%v", *validationRule.Reason)
	}

	if validationRule.OptionalOldSelf != nil {
		strToHash += fmt.Sprintf("%v", *validationRule.OptionalOldSelf)
	}

	h := sha256.New()

	h.Write([]byte(strToHash))

	return string(h.Sum(nil))
}

// Compare compares an old and a new JSONSchemaProps, checking for
// incompatible changes to the x-kubernetes-validations of a property.
func (x *XValidations) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	oldXValidations := map[string]apiextensionsv1.ValidationRule{}
	newXValidations := map[string]apiextensionsv1.ValidationRule{}

	oldHashes := set.New[string]()
	newHashes := set.New[string]()

	for _, oldVal := range a.XValidations {
		hash := hashValidationRule(oldVal)
		oldXValidations[hash] = oldVal
		oldHashes.Insert(hash)
	}

	for _, newVal := range b.XValidations {
		hash := hashValidationRule(newVal)
		newXValidations[hash] = newVal
		newHashes.Insert(hash)
	}

	errs := []error{}

	oldRemovedHashes := oldHashes.Difference(newHashes)
	for _, hash := range oldRemovedHashes.SortedList() {
		r, _ := json.Marshal(oldXValidations[hash])
		errs = append(errs, fmt.Errorf("%w: %s", ErrXValidationRemoved, string(r)))
	}

	newAddedHashes := newHashes.Difference(oldHashes)
	for _, hash := range newAddedHashes.SortedList() {
		r, _ := json.Marshal(newXValidations[hash])
		errs = append(errs, fmt.Errorf("%w: %s", ErrXValidationAdded, string(r)))
	}

	a.XValidations = nil
	b.XValidations = nil

	return validations.HandleErrors(x.Name(), x.enforcement, errs...)
}

var (
	// ErrXValidationAdded represents an error state when an x-kubernetes-validation rule is added to a property.
	ErrXValidationAdded = errors.New("x-kubernetes-validations added")
	// ErrXValidationRemoved represents an error state when an x-kubernetes-validation rule is removed from a property.
	ErrXValidationRemoved = errors.New("x-kubernetes-validations removed")
)

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

package served

import (
	"fmt"
	"slices"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	versionhelper "k8s.io/apimachinery/pkg/version"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

// Validator validates Kubernetes CustomResourceDefinitions using the configured validations.
type Validator struct {
	comparators          []validations.Comparator[apiextensionsv1.JSONSchemaProps]
	conversionPolicy     config.ConversionPolicy
	unhandledEnforcement config.EnforcementPolicy
}

// ValidatorOption configures a Validator.
type ValidatorOption func(*Validator)

// WithComparators configures a Validator with the provided JSONSchemaProps Comparators.
// Each call to WithComparators is a replacement, not additive.
func WithComparators(comparators ...validations.Comparator[apiextensionsv1.JSONSchemaProps]) ValidatorOption {
	return func(v *Validator) {
		v.comparators = comparators
	}
}

// WithUnhandledEnforcementPolicy sets the unhandled enforcement policy for the validator.
func WithUnhandledEnforcementPolicy(policy config.EnforcementPolicy) ValidatorOption {
	return func(v *Validator) {
		if policy == "" {
			policy = config.EnforcementPolicyError
		}

		v.unhandledEnforcement = policy
	}
}

// WithConversionPolicy sets the conversion policy for the validator.
func WithConversionPolicy(policy config.ConversionPolicy) ValidatorOption {
	return func(v *Validator) {
		if policy == "" {
			policy = config.ConversionPolicyNone
		}

		v.conversionPolicy = policy
	}
}

// New creates a new Validator to validate the served versions of an old and new CustomResourceDefinition
// configured with the provided ValidatorOptions.
func New(opts ...ValidatorOption) *Validator {
	validator := &Validator{
		comparators:          []validations.Comparator[apiextensionsv1.JSONSchemaProps]{},
		conversionPolicy:     config.ConversionPolicyNone,
		unhandledEnforcement: config.EnforcementPolicyError,
	}

	for _, opt := range opts {
		opt(validator)
	}

	return validator
}

// Validate runs the validations configured in the Validator.
func (v *Validator) Validate(_, b *apiextensionsv1.CustomResourceDefinition) map[string]map[string][]validations.ComparisonResult {
	result := map[string]map[string][]validations.ComparisonResult{}

	// If conversion webhook is specified and conversion policy is ignore, pass check
	if v.conversionPolicy == config.ConversionPolicyIgnore && b.Spec.Conversion != nil && b.Spec.Conversion.Strategy == apiextensionsv1.WebhookConverter {
		return result
	}

	servedVersions := []apiextensionsv1.CustomResourceDefinitionVersion{}

	for _, version := range b.Spec.Versions {
		if version.Served {
			servedVersions = append(servedVersions, version)
		}
	}

	slices.SortFunc(servedVersions, func(a, b apiextensionsv1.CustomResourceDefinitionVersion) int {
		return versionhelper.CompareKubeAwareVersionStrings(a.Name, b.Name)
	})

	for i, oldVersion := range servedVersions[:len(servedVersions)-1] {
		for _, newVersion := range servedVersions[i+1:] {
			resultVersion := fmt.Sprintf("%s <-> %s", oldVersion.Name, newVersion.Name)
			result[resultVersion] = validations.CompareVersions(oldVersion, newVersion, v.unhandledEnforcement, v.comparators...)
		}
	}

	return result
}

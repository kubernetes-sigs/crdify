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
	"bytes"
	"fmt"
	"maps"
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
	versionComparison    config.ServedVersionComparison
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

// WithVersionComparison sets the version comparison for the validator.
func WithVersionComparison(versionComparison config.ServedVersionComparison) ValidatorOption {
	return func(v *Validator) {
		if versionComparison == "" {
			versionComparison = config.ServedVersionComparisonAll
		}

		v.versionComparison = versionComparison
	}
}

// New creates a new Validator to validate the served versions of an old and new CustomResourceDefinition
// configured with the provided ValidatorOptions.
func New(opts ...ValidatorOption) *Validator {
	validator := &Validator{
		comparators:          []validations.Comparator[apiextensionsv1.JSONSchemaProps]{},
		conversionPolicy:     config.ConversionPolicyNone,
		unhandledEnforcement: config.EnforcementPolicyError,
		versionComparison:    config.ServedVersionComparisonAll,
	}

	for _, opt := range opts {
		opt(validator)
	}

	return validator
}

// Validate runs the validations configured in the Validator.
func (v *Validator) Validate(a, b *apiextensionsv1.CustomResourceDefinition) map[string]map[string][]validations.ComparisonResult {
	result := map[string]map[string][]validations.ComparisonResult{}

	// If we are not comparing any versions, pass check
	if v.versionComparison == config.ServedVersionComparisonNone {
		return result
	}

	// If conversion webhook is specified and conversion policy is ignore, pass check
	if v.conversionPolicy == config.ConversionPolicyIgnore && b.Spec.Conversion != nil && b.Spec.Conversion.Strategy == apiextensionsv1.WebhookConverter {
		return result
	}

	bPairs := buildVersionPairs(b)

	// If we are only comparing pairs where a schema change was detected between old and new CRD
	// filter the list of pairs by removing pairs whose schemas match between old and new.
	if v.versionComparison == config.ServedVersionComparisonOnlyDiff {
		aPairs := buildVersionPairs(a)

		maps.DeleteFunc(bPairs, func(bNames [2]string, bVersions [2]apiextensionsv1.CustomResourceDefinitionVersion) bool {
			if aVersions, ok := aPairs[bNames]; ok {
				if versionSchemasEqual(aVersions[0], bVersions[0]) && versionSchemasEqual(aVersions[1], bVersions[1]) {
					return true
				}
			}

			return false
		})
	}

	// Compare the served version pairs
	for versionNames, versions := range bPairs {
		resultVersion := fmt.Sprintf("%s <-> %s", versionNames[0], versionNames[1])
		result[resultVersion] = validations.CompareVersions(versions[0], versions[1], v.unhandledEnforcement, v.comparators...)
	}

	return result
}

func buildVersionPairs(crd *apiextensionsv1.CustomResourceDefinition) map[[2]string][2]apiextensionsv1.CustomResourceDefinitionVersion {
	servedVersions := make(map[string]apiextensionsv1.CustomResourceDefinitionVersion, len(crd.Spec.Versions))

	for _, version := range crd.Spec.Versions {
		if version.Served && version.Schema != nil {
			servedVersions[version.Name] = version
		}
	}

	servedVersionNames := slices.Sorted(maps.Keys(servedVersions))
	slices.SortFunc(servedVersionNames, versionhelper.CompareKubeAwareVersionStrings)

	n := len(servedVersionNames)
	pairs := make(map[[2]string][2]apiextensionsv1.CustomResourceDefinitionVersion, ((n-1)*n)/2)

	for i, iName := range servedVersionNames[:len(servedVersionNames)-1] {
		iVersion := servedVersions[iName]

		for _, newVersion := range servedVersionNames[i+1:] {
			jVersion := servedVersions[newVersion]
			pairs[[2]string{iName, newVersion}] = [2]apiextensionsv1.CustomResourceDefinitionVersion{iVersion, jVersion}
		}
	}

	return pairs
}

func versionSchemasEqual(x, y apiextensionsv1.CustomResourceDefinitionVersion) bool {
	if (x.Schema == nil) != (y.Schema == nil) {
		return false
	}

	if x.Schema == nil && y.Schema == nil {
		return true
	}

	xData, xErr := x.Schema.Marshal()
	yData, yErr := y.Schema.Marshal()

	if xErr != nil || yErr != nil {
		return false
	}

	return bytes.Equal(xData, yData)
}

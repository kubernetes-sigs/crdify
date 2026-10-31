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

package runner

import (
	"fmt"
	"maps"
	"slices"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
	"sigs.k8s.io/crdify/pkg/validators/crd"
	"sigs.k8s.io/crdify/pkg/validators/version/same"
	"sigs.k8s.io/crdify/pkg/validators/version/served"
)

// Runner is a utility struct for running
// - Whole CRD validations
// - Same version validations
// - Served version validations.
type Runner struct {
	crdValidator           *crd.Validator
	sameVersionValidator   *same.Validator
	servedVersionValidator *served.Validator
}

// New returns a new instance of a Runner using the provided Config and validations.Registry
// to build the CRD, same version, and served version validators.
// It returns an error if any errors are encountered.
func New(cfg *config.Config, registry validations.Registry) (*Runner, error) {
	initialValidations, err := validations.LoadValidationsFromRegistry(registry)
	if err != nil {
		return nil, fmt.Errorf("loading validations from registry: %w", err)
	}

	configuredValidations, err := validations.ConfigureValidations(initialValidations, registry, *cfg)
	if err != nil {
		return nil, fmt.Errorf("configuring validations: %w", err)
	}

	vals := slices.Collect(maps.Values(configuredValidations))

	crdComparators := validations.ComparatorsForValidations[apiextensionsv1.CustomResourceDefinition](vals...)
	propertyComparators := validations.ComparatorsForValidations[apiextensionsv1.JSONSchemaProps](vals...)

	return &Runner{
		crdValidator:           crd.New(crd.WithComparators(crdComparators...)),
		sameVersionValidator:   same.New(same.WithComparators(propertyComparators...), same.WithUnhandledEnforcementPolicy(cfg.UnhandledEnforcement)),
		servedVersionValidator: served.New(served.WithComparators(propertyComparators...), served.WithUnhandledEnforcementPolicy(cfg.UnhandledEnforcement), served.WithConversionPolicy(cfg.Conversion)),
	}, nil
}

// Run executes all the validators and collects the results into a utility struct for
// reporting and evaluating the results.
func (i *Runner) Run(oldCrd, newCrd *apiextensionsv1.CustomResourceDefinition) *Results {
	return &Results{
		CRDValidation:           i.crdValidator.Validate(oldCrd, newCrd),
		SameVersionValidation:   i.sameVersionValidator.Validate(oldCrd, newCrd),
		ServedVersionValidation: i.servedVersionValidator.Validate(oldCrd, newCrd),
	}
}

// RunByName pairs the old and new CustomResourceDefinitions by name, executes all the validators
// on each pair and returns the results keyed by name.
// A CustomResourceDefinition that is only in oldCrds has been removed and is reported as an error.
// A CustomResourceDefinition that is only in newCrds has been added and is not validated.
func (i *Runner) RunByName(oldCrds, newCrds []*apiextensionsv1.CustomResourceDefinition) ResultsByName {
	newCrdsByName := make(map[string]*apiextensionsv1.CustomResourceDefinition, len(newCrds))
	for _, newCrd := range newCrds {
		newCrdsByName[newCrd.Name] = newCrd
	}

	results := make(ResultsByName, len(oldCrds))

	for _, oldCrd := range oldCrds {
		newCrd, ok := newCrdsByName[oldCrd.Name]
		if !ok {
			results[oldCrd.Name] = &Results{
				CRDValidation: []validations.ComparisonResult{{Name: "removed", Errors: []string{"CustomResourceDefinition removed"}}},
			}

			continue
		}

		results[oldCrd.Name] = i.Run(oldCrd, newCrd)
	}

	return results
}

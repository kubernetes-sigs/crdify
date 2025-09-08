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
	"encoding/json"
	"errors"
	"fmt"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/crdify/pkg/config"
	"sigs.k8s.io/crdify/pkg/validations"
)

const oneOfValidationName = "oneOf"

var (
	_ validations.Validation                                  = (*OneOf)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*OneOf)(nil)
)

// RegisterOneOf registers the OneOf validation with the provided validation registry.
func RegisterOneOf(registry validations.Registry) {
	registry.Register(oneOfValidationName, oneOfFactory)
}

// oneOfFactory initializes an OneOf validation from configuration.
func oneOfFactory(cfg map[string]interface{}) (validations.Validation, error) {
	oneOfCfg := &OneOfConfig{}
	if err := ConfigToType(cfg, oneOfCfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if err := ValidateOneOfConfig(oneOfCfg); err != nil {
		return nil, fmt.Errorf("validating oneOf config: %w", err)
	}
	return &OneOf{OneOfConfig: *oneOfCfg}, nil
}

// ValidateOneOfConfig validates the OneOfConfig, setting defaults.
func ValidateOneOfConfig(in *OneOfConfig) error {
	return nil
}

// OneOfConfig contains configurations for the OneOf validation.
type OneOfConfig struct {
}

// OneOf validation identifies incompatible changes to oneOf constraints.
type OneOf struct {
	OneOfConfig
	enforcement config.EnforcementPolicy
}

// Name returns the name of the validation.
func (o *OneOf) Name() string {
	return oneOfValidationName
}

// SetEnforcement sets the enforcement policy.
func (o *OneOf) SetEnforcement(policy config.EnforcementPolicy) {
	o.enforcement = policy
}

// Compare checks for incompatible changes in the oneOf constraint.
func (o *OneOf) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	oldSubSchemas := sets.New[string]()
	for _, schema := range a.OneOf {
		marshalled, err := marshallSchema(schema)
		if err != nil {
			return validations.HandleErrors(o.Name(), o.enforcement, fmt.Errorf("failed to marshal old oneOf subschema: %w", err))
		}
		oldSubSchemas.Insert(marshalled)
	}

	newSubSchemas := sets.New[string]()
	for _, schema := range b.OneOf {
		marshalled, err := marshallSchema(schema)
		if err != nil {
			return validations.HandleErrors(o.Name(), o.enforcement, fmt.Errorf("failed to marshal new oneOf schema: %w", err))
		}
		newSubSchemas.Insert(marshalled)
	}

	var errs []error
	if oldSubSchemas.Len() == 0 && newSubSchemas.Len() > 0 {
		errs = append(errs, ErrNetNewOneOfConstraint)
	}

	removed := oldSubSchemas.Difference(newSubSchemas)
	if removed.Len() > 0 {
		errs = append(errs, fmt.Errorf("%w: %v", ErrRemovedOneOf, sets.List(removed)))
	}

	added := newSubSchemas.Difference(oldSubSchemas)
	if added.Len() > 0 {
		errs = append(errs, fmt.Errorf("%w: %v", ErrAddedOneOf, sets.List(added)))
	}

	a.OneOf = nil
	b.OneOf = nil

	return validations.HandleErrors(o.Name(), o.enforcement, utilerrors.NewAggregate(errs))
}

// marshallSchema converts a schema reprsented as aJSONSchemaProps into a JSON string which captures the structure of the schema.
func marshallSchema(schema apiextensionsv1.JSONSchemaProps) (string, error) {
	// Use a copy to avoid modifying the original
	schemaCopy := schema.DeepCopy()
	// Remove fields that don't affect the structure of the schema
	schemaCopy.Description = ""
	schemaCopy.Example = nil

	bytes, err := json.Marshal(schemaCopy)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %w", err)
	}
	return string(bytes), nil
}

var (
	ErrNetNewOneOfConstraint = errors.New("oneOf constraint added when there was none previously")
	ErrRemovedOneOf          = errors.New("allowed oneOf schemas removed")
	ErrAddedOneOf            = errors.New("allowed oneOf schemas added")
)

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

var (
	_ validations.Validation                                  = (*Format)(nil)
	_ validations.Comparator[apiextensionsv1.JSONSchemaProps] = (*Format)(nil)
)

const formatValidationName = "format"

// RegisterFormat registers the Format validation
// with the provided validation registry.
func RegisterFormat(registry validations.Registry) {
	registry.Register(formatValidationName, formatFactory)
}

// formatFactory is a function used to initialize a Format validation
// implementation based on the provided configuration.
func formatFactory(_ map[string]interface{}) (validations.Validation, error) {
	return &Format{}, nil
}

// Format is a Validation that can be used to identify
// incompatible changes to the format constraints of CRD properties.
type Format struct {
	enforcement config.EnforcementPolicy
}

// Name returns the name of the Format validation.
func (t *Format) Name() string {
	return formatValidationName
}

// SetEnforcement sets the EnforcementPolicy for the Format validation.
func (t *Format) SetEnforcement(policy config.EnforcementPolicy) {
	t.enforcement = policy
}

// Compare compares an old and a new JSONSchemaProps, checking for incompatible changes to the format constraints of a property.
// In order for callers to determine if diffs to a JSONSchemaProps have been handled by this validation
// the JSONSchemaProps.format field will be reset to '""' as part of this method.
// It is highly recommended that only copies of the JSONSchemaProps to compare are provided to this method
// to prevent unintentional modifications.
func (t *Format) Compare(a, b *apiextensionsv1.JSONSchemaProps) validations.ComparisonResult {
	var err error
	if a.Format != b.Format {
		err = fmt.Errorf("%w : %q -> %q", ErrFormatChanged, a.Format, b.Format)
	}

	a.Format = ""
	b.Format = ""

	return validations.HandleErrors(t.Name(), t.enforcement, err)
}

// ErrFormatChanged represents an error state when a property Format changed.
var ErrFormatChanged = errors.New("format changed")

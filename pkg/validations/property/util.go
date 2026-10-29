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
	"bytes"
	"encoding/json"
	"fmt"
	"slices"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

// ConfigToType is a utility function for unmarshalling unstructured data (map[string]interface{})
// into a concrete type.
func ConfigToType[T any](in map[string]interface{}, out *T) error {
	jsonBytes, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshalling input to JSON: %w", err)
	}

	err = json.Unmarshal(jsonBytes, out)
	if err != nil {
		return fmt.Errorf("unmarshalling input to output: %w", err)
	}

	return nil
}

// normalizeSchema recursively zeroes non-structural fields on a schema
// and sorts lists where order is irrelevant so that only
// validation-relevant differences are compared.
func normalizeSchema(schema *apiextensionsv1.JSONSchemaProps) {
	if schema == nil {
		return
	}

	schema.Description = ""
	schema.Title = ""
	schema.Example = nil
	schema.ExternalDocs = nil

	slices.Sort(schema.Required)
	slices.SortFunc(schema.Enum, func(a, b apiextensionsv1.JSON) int {
		return bytes.Compare(a.Raw, b.Raw)
	})

	normalizeSchemaSlice(schema.AllOf)
	normalizeSchemaSlice(schema.OneOf)
	normalizeSchemaSlice(schema.AnyOf)
	normalizeSchema(schema.Not)
	normalizeSchemaMap(schema.Properties)
	normalizeSchemaMap(schema.PatternProperties)
	normalizeSchemaMap(schema.Definitions)

	if schema.Items != nil {
		normalizeSchema(schema.Items.Schema)
		normalizeSchemaSlice(schema.Items.JSONSchemas)
	}

	if schema.AdditionalProperties != nil {
		normalizeSchema(schema.AdditionalProperties.Schema)
	}

	if schema.AdditionalItems != nil {
		normalizeSchema(schema.AdditionalItems.Schema)
	}

	for _, dependency := range schema.Dependencies {
		normalizeSchema(dependency.Schema)
	}
}

func normalizeSchemaSlice(schemas []apiextensionsv1.JSONSchemaProps) {
	for i := range schemas {
		normalizeSchema(&schemas[i])
	}
}

func normalizeSchemaMap[M ~map[string]apiextensionsv1.JSONSchemaProps](schemas M) {
	for name, schema := range schemas {
		normalizeSchema(&schema)
		schemas[name] = schema
	}
}

// schemaKeys returns the set of comparable representations of the provided schemas.
func schemaKeys(schemas []apiextensionsv1.JSONSchemaProps) (sets.Set[string], error) {
	keys := sets.New[string]()

	for i := range schemas {
		key, err := schemaKey(&schemas[i])
		if err != nil {
			return nil, err
		}

		keys.Insert(key)
	}

	return keys, nil
}

// schemaKey normalizes the provided schema and returns a comparable
// representation of it. JSON is used instead of String() because String()
// prints nested pointer fields as memory addresses.
func schemaKey(schema *apiextensionsv1.JSONSchemaProps) (string, error) {
	normalizeSchema(schema)

	out, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("marshalling schema to JSON: %w", err)
	}

	return string(out), nil
}

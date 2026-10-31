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

package decode

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CRDs decodes every CustomResourceDefinition in the provided YAML or JSON content,
// which may hold several YAML documents. Documents of any other kind and empty
// documents are skipped.
// It returns an error if the content holds no CustomResourceDefinition or holds
// more than one CustomResourceDefinition with the same name.
func CRDs(content []byte) ([]*apiextensionsv1.CustomResourceDefinition, error) {
	reader := yaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(content)))
	crds := []*apiextensionsv1.CustomResourceDefinition{}
	names := sets.New[string]()

	for {
		document, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("reading YAML document: %w", err)
		}

		typeMeta := &metav1.TypeMeta{}

		err = yaml.Unmarshal(document, typeMeta)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling YAML document: %w", err)
		}

		if typeMeta.Kind != "CustomResourceDefinition" {
			continue
		}

		crd := &apiextensionsv1.CustomResourceDefinition{}

		err = yaml.Unmarshal(document, crd)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling CustomResourceDefinition: %w", err)
		}

		if names.Has(crd.Name) {
			return nil, fmt.Errorf("%w : %q", errDuplicateName, crd.Name)
		}

		names.Insert(crd.Name)
		crds = append(crds, crd)
	}

	if len(crds) == 0 {
		return nil, errNoCRD
	}

	return crds, nil
}

var (
	errNoCRD         = errors.New("no CustomResourceDefinition found")
	errDuplicateName = errors.New("more than one CustomResourceDefinition with the same name")
)

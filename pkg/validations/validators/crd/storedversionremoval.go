package crd

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type StoredVersionRemoval struct{}

func (svr *StoredVersionRemoval) Name() string {
	return "StoredVersionRemoval"
}

func (svr *StoredVersionRemoval) Validate(old, new *apiextensionsv1.CustomResourceDefinition) error {
	newVersions := sets.New[string]()
	for _, version := range new.Spec.Versions {
		if !newVersions.Has(version.Name) {
			newVersions.Insert(version.Name)
		}
	}

	for _, storedVersion := range old.Status.StoredVersions {
		if !newVersions.Has(storedVersion) {
			return fmt.Errorf("stored version %q removed", storedVersion)
		}
	}
	return nil
}

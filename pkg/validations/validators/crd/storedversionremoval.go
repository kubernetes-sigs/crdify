package crd

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type StoredVersionRemoval struct{}

func (svr *StoredVersionRemoval) Name() string {
	return "storedVersionRemoval"
}

func (svr *StoredVersionRemoval) Validate(old, new *apiextensionsv1.CustomResourceDefinition) ValidationResult {
	newVersions := sets.New[string]()
	for _, version := range new.Spec.Versions {
		if !newVersions.Has(version.Name) {
			newVersions.Insert(version.Name)
		}
	}

	removedVersions := []string{}
	for _, storedVersion := range old.Status.StoredVersions {
		if !newVersions.Has(storedVersion) {
			removedVersions = append(removedVersions, storedVersion)
		}
	}

	vr := &validationResult{
		Validation: svr.Name(),
	}

	if len(removedVersions) > 0 {
		vr.Err = fmt.Sprintf("stored versions %v removed", removedVersions)
	}

	return vr
}

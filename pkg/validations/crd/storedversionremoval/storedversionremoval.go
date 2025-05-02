package storedversionremoval

import (
	"fmt"

	"github.com/everettraven/crd-diff/pkg/config"
	"github.com/everettraven/crd-diff/pkg/validations"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	_ validations.Validation                                           = (*StoredVersionRemoval)(nil)
	_ validations.Comparator[apiextensionsv1.CustomResourceDefinition] = (*StoredVersionRemoval)(nil)
)

const name = "storedVersionRemoval"

// Register registers the StoredVersionRemoval validation
// with the provided validation registry
func Register(registry validations.Registry) {
	registry.Register(name, factory)
}

// factory is a function used to initialize a StoredVersionRemoval validation
// implementation based on the provided configuration.
func factory(_ map[string]interface{}) (validations.Validation, error) {
	return &StoredVersionRemoval{}, nil
}

// StoredVersionRemoval is a validations.Validation implementation
// used to check if any versions existing in the set of stored versions
// has been removed in the new instance of the CustomResourceDefinition
type StoredVersionRemoval struct {
	// enforcement is the EnforcementPolicy that this validation
	// should use when performing its validation logic
	enforcement config.EnforcementPolicy
}

// Name returns the name of the StoredVersionRemoval validation
func (svr *StoredVersionRemoval) Name() string {
	return name
}

// SetEnforcement sets the EnforcementPolicy for the StoredVersionRemoval validation
func (svr *StoredVersionRemoval) SetEnforcement(enforcement config.EnforcementPolicy) {
	svr.enforcement = enforcement
}

// Compare compares an old and a new CustomResourceDefintion, checking for removal of
// any stored versions present in the old CustomResourceDefinition in the new instance
// of the CustomResourceDefinition
func (svr *StoredVersionRemoval) Compare(old, new *apiextensionsv1.CustomResourceDefinition) validations.ComparisonResult {
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

	var err error
	if len(removedVersions) > 0 {
		err = fmt.Errorf("stored versions %v removed", removedVersions)
	}

	return validations.HandleErrors(svr.Name(), svr.enforcement, err)
}

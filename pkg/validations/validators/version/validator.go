package version

import (
	"errors"
	"fmt"
	"slices"

	"github.com/everettraven/crd-diff/pkg/validations/property"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	versionhelper "k8s.io/apimachinery/pkg/version"
)

type FailureMode string

const (
	FailureModeOpen   = "Open"
	FailureModeClosed = "Closed"
)

type Validator struct {
	sameVersionConfig   SameVersionConfig
	servedVersionConfig ServedVersionConfig
}

type SameVersionConfig struct {
	UnhandledFailureMode FailureMode
	Skip                 bool
	Validations          []property.Validation
}

type ServedVersionConfig struct {
	UnhandledFailureMode FailureMode
	Skip                 bool
	IgnoreConversion     bool
	Validations          []property.Validation
}

type ValidatorOption func(*Validator)

func WithSameVersionConfig(cfg SameVersionConfig) ValidatorOption {
	return func(v *Validator) {
		v.sameVersionConfig = cfg
	}
}

func WithServedVersionConfig(cfg ServedVersionConfig) ValidatorOption {
	return func(v *Validator) {
		v.servedVersionConfig = cfg
	}
}

func NewValidator(opts ...ValidatorOption) *Validator {
	validator := &Validator{
		sameVersionConfig: SameVersionConfig{
			UnhandledFailureMode: FailureModeClosed,
			Skip:                 false,
			Validations:          []property.Validation{},
		},
		servedVersionConfig: ServedVersionConfig{
			UnhandledFailureMode: FailureModeClosed,
			Skip:                 false,
			IgnoreConversion:     false,
			Validations:          []property.Validation{},
		},
	}

	for _, opt := range opts {
		opt(validator)
	}
	return validator
}

func (v *Validator) Name() string {
	return "Version"
}

func (v *Validator) Validate(old, new *apiextensionsv1.CustomResourceDefinition) error {
	errs := []error{}
	if !v.sameVersionConfig.Skip {
		err := v.ValidateSameVersions(old, new)
		if err != nil {
			errs = append(errs, fmt.Errorf("validating same versions: %w", err))
		}
	}

	if !v.servedVersionConfig.Skip {
		err := v.ValidateServedVersions(new)
		if err != nil {
			errs = append(errs, fmt.Errorf("validating new served versions: %w", err))
		}
	}
	return errors.Join(errs...)
}

func (v *Validator) ValidateSameVersions(old, new *apiextensionsv1.CustomResourceDefinition) error {
	errs := []error{}
	for _, oldVersion := range old.Spec.Versions {
		newVersion := GetCRDVersionByName(new, oldVersion.Name)
		// in this case, there is nothing to compare. Generally, the removal
		// of an existing version is a breaking change. It may be considered safe
		// if there are no CRs stored at that version or migration has successfully
		// occurred. Since the safety of this varies and we don't have explicit
		// knowledge of this we assume a separate check will be in place to capture
		// this as a breaking change if a user desires.
		if newVersion == nil {
			continue
		}

		errs = append(errs, CompareVersions(oldVersion, *newVersion, v.sameVersionConfig.UnhandledFailureMode, v.sameVersionConfig.Validations))
	}
	return errors.Join(errs...)
}

func (v *Validator) ValidateServedVersions(crd *apiextensionsv1.CustomResourceDefinition) error {
	// If conversion webhook is specified, pass check
	if crd.Spec.Conversion != nil && crd.Spec.Conversion.Strategy == apiextensionsv1.WebhookConverter {
		return nil
	}

	servedVersions := []apiextensionsv1.CustomResourceDefinitionVersion{}
	for _, version := range crd.Spec.Versions {
		if version.Served {
			servedVersions = append(servedVersions, version)
		}
	}

	slices.SortFunc(servedVersions, func(a, b apiextensionsv1.CustomResourceDefinitionVersion) int {
		return versionhelper.CompareKubeAwareVersionStrings(a.Name, b.Name)
	})

	errs := []error{}
	for i, oldVersion := range servedVersions[:len(servedVersions)-1] {
		for _, newVersion := range servedVersions[i+1:] {
			errs = append(errs, CompareVersions(oldVersion, newVersion, v.servedVersionConfig.UnhandledFailureMode, v.servedVersionConfig.Validations))
		}
	}
	return errors.Join(errs...)
}

func CompareVersions(old, new apiextensionsv1.CustomResourceDefinitionVersion, failureMode FailureMode, validations []property.Validation) error {
	oldFlattened := FlattenCRDVersion(old)
	newFlattened := FlattenCRDVersion(new)

	diffs := FlattenedCRDVersionDiff(oldFlattened, newFlattened)
	errs := []error{}
	for property, diff := range diffs {
		handled := false
		for _, validation := range validations {
			ok, err := validation.Validate(diff)
			if err != nil {
				errs = append(errs, fmt.Errorf("old version %q compared to new version %q, property %q failed validation %q: %w", old.Name, new.Name, property, validation.Name(), err))
			}
			// if the validation handled this difference continue to the next difference
			if ok {
				handled = true
				break
			}
		}

		if failureMode == FailureModeClosed && !handled {
			errs = append(errs, fmt.Errorf("old version %q compared to new version %q, property %q has unknown change, refusing to determine that change is safe", old.Name, new.Name, property))
		}
	}
	return errors.Join(errs...)
}

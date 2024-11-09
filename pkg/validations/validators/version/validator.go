package version

import (
	"errors"
	"fmt"
	"slices"

	"github.com/everettraven/crd-diff/pkg/validations/property"
	"github.com/everettraven/crd-diff/pkg/validations/validators/crd"
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
	return "version"
}

func (v *Validator) Validate(old, new *apiextensionsv1.CustomResourceDefinition) crd.ValidationResult {
	result := Result{
        Validation: v.Name(),
		SameVersionResults:   []VersionCompareResult{},
		ServedVersionResults: []VersionCompareResult{},
	}
	if !v.sameVersionConfig.Skip {
		result.SameVersionResults = v.ValidateSameVersions(old, new)
	}

	if !v.servedVersionConfig.Skip {
		result.ServedVersionResults = v.ValidateServedVersions(new)
	}
	return &result
}

func (v *Validator) ValidateSameVersions(old, new *apiextensionsv1.CustomResourceDefinition) []VersionCompareResult {
	vcrs := []VersionCompareResult{}
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

		vcrs = append(vcrs, CompareVersions(oldVersion, *newVersion, v.sameVersionConfig.UnhandledFailureMode, v.sameVersionConfig.Validations))
	}
	return vcrs
}

func (v *Validator) ValidateServedVersions(crd *apiextensionsv1.CustomResourceDefinition) []VersionCompareResult {
	// If conversion webhook is specified, pass check
	if !v.servedVersionConfig.IgnoreConversion && crd.Spec.Conversion != nil && crd.Spec.Conversion.Strategy == apiextensionsv1.WebhookConverter {
		return []VersionCompareResult{}
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

	vcrs := []VersionCompareResult{}
	for i, oldVersion := range servedVersions[:len(servedVersions)-1] {
		for _, newVersion := range servedVersions[i+1:] {
			vcrs = append(vcrs, CompareVersions(oldVersion, newVersion, v.servedVersionConfig.UnhandledFailureMode, v.servedVersionConfig.Validations))
		}
	}
	return vcrs
}

func CompareVersions(old, new apiextensionsv1.CustomResourceDefinitionVersion, failureMode FailureMode, validations []property.Validation) VersionCompareResult {
	oldFlattened := FlattenCRDVersion(old)
	newFlattened := FlattenCRDVersion(new)

	diffs := FlattenedCRDVersionDiff(oldFlattened, newFlattened)
	result := VersionCompareResult{
		VersionA:               old.Name,
		VersionB:               new.Name,
		PropertyCompareResults: []PropertyCompareResult{},
	}

	for property, diff := range diffs {
		errs := ComparePropertyDiff(diff, failureMode, validations)
		result.PropertyCompareResults = append(result.PropertyCompareResults, PropertyCompareResult{
			Property: property,
			Errors: convert(errs, func(s error) string {
				return s.Error()
			}),
		})
	}

	return result
}

func ComparePropertyDiff(diff property.Diff, failureMode FailureMode, validations []property.Validation) []error {
	errs := []error{}
	handled := false
	for _, validation := range validations {
		ok, err := validation.Validate(diff)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s validation failed: %w", validation.Name(), err))
		}
		// if the validation handled this difference continue to the next difference
		if ok {
			handled = true
			break
		}
	}

	if failureMode == FailureModeClosed && !handled {
		errs = append(errs, errors.New("unknown change(s), refusing to determine that change is safe"))
	}
	return errs
}

func convert[S any, E any](s []S, convertFunc func(S) E) []E {
	converted := make([]E, len(s))
	for i, val := range s {
		converted[i] = convertFunc(val)
	}
	return converted
}

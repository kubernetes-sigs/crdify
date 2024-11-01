package version

import (
	"errors"
	"fmt"
	"slices"

	"github.com/everettraven/crd-diff/pkg/validations/property"
	"github.com/everettraven/crd-diff/pkg/validations/results"
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

func (v *Validator) Validate(old, new *apiextensionsv1.CustomResourceDefinition) *results.Result {
	result := &results.Result{
		Subresults: []*results.Result{},
	}
	if !v.sameVersionConfig.Skip {
		res := v.ValidateSameVersions(old, new)
		if res != nil {
			result.Subresults = append(result.Subresults, res)
			if res.Error != nil {
				result.Error = errors.New("checking versions for compatibility")
			}
		}
	}

	if !v.servedVersionConfig.Skip {
		res := v.ValidateServedVersions(new)
		if res != nil {
			result.Subresults = append(result.Subresults, res)
			if res.Error != nil {
				result.Error = errors.New("checking versions for compatibility")
			}
		}
	}
	return result
}

func (v *Validator) ValidateSameVersions(old, new *apiextensionsv1.CustomResourceDefinition) *results.Result {
	result := &results.Result{
		Subresults: []*results.Result{},
	}
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

		res := CompareVersions(oldVersion, *newVersion, v.sameVersionConfig.UnhandledFailureMode, v.sameVersionConfig.Validations)
		if res != nil {
			subResult := &results.Result{
				Subresults: []*results.Result{res},
			}
			if res.Error != nil {
				result.Error = errors.New("comparing same versions")
				subResult.Error = fmt.Errorf("comparing version %q to version %q", oldVersion.Name, newVersion.Name)
			}
			result.Subresults = append(result.Subresults, subResult)
		}
	}
	return result
}

func (v *Validator) ValidateServedVersions(crd *apiextensionsv1.CustomResourceDefinition) *results.Result {
	// If conversion webhook is specified, pass check
	if !v.servedVersionConfig.IgnoreConversion && crd.Spec.Conversion != nil && crd.Spec.Conversion.Strategy == apiextensionsv1.WebhookConverter {
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

	result := &results.Result{
		Subresults: []*results.Result{},
	}
	for i, oldVersion := range servedVersions[:len(servedVersions)-1] {
		for _, newVersion := range servedVersions[i+1:] {
			res := CompareVersions(oldVersion, newVersion, v.servedVersionConfig.UnhandledFailureMode, v.servedVersionConfig.Validations)
			if res != nil {
				subResult := &results.Result{
					Subresults: []*results.Result{res},
				}
				if res.Error != nil {
					result.Error = errors.New("comparing served versions")
					subResult.Error = fmt.Errorf("comparing version %q to version %q", oldVersion.Name, newVersion.Name)
				}
				result.Subresults = append(result.Subresults, subResult)
			}
		}
	}
	return result
}

func CompareVersions(old, new apiextensionsv1.CustomResourceDefinitionVersion, failureMode FailureMode, validations []property.Validation) *results.Result {
	oldFlattened := FlattenCRDVersion(old)
	newFlattened := FlattenCRDVersion(new)

	diffs := FlattenedCRDVersionDiff(oldFlattened, newFlattened)
	result := &results.Result{
		Subresults: []*results.Result{},
	}
	for property, diff := range diffs {
		res := ComparePropertyDiff(diff, failureMode, validations)
		if res != nil {
			subResult := &results.Result{
				Subresults: []*results.Result{res},
			}
			if res.Error != nil {
				result.Error = errors.New("property validation failures")
				propError := fmt.Errorf("property %q", property)
				subResult.Error = propError
			}
			result.Subresults = append(result.Subresults, subResult)
		}
	}
	return result
}

func ComparePropertyDiff(diff property.Diff, failureMode FailureMode, validations []property.Validation) *results.Result {
	result := &results.Result{
		Subresults: []*results.Result{},
	}
	handled := false
	for _, validation := range validations {
		ok, res := validation.Validate(diff)
		if res != nil {
			subResult := &results.Result{
				Subresults: []*results.Result{res},
			}
			if res.Error != nil {
				result.Error = errors.New("failed validations")
				subResult.Error = fmt.Errorf("%q validation failed", validation.Name())
			}
			result.Subresults = append(result.Subresults, subResult)
		}
		// if the validation handled this difference continue to the next difference
		if ok {
			handled = true
			break
		}
	}

	if failureMode == FailureModeClosed && !handled {
		result.Error = errors.New("validation failed")
		result.Subresults = append(result.Subresults, &results.Result{
			Error:      errors.New("unknown change(s), refusing to determine that change is safe"),
			Subresults: []*results.Result{},
		})
	}
	return result
}

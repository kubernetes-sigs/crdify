package property

import (
	"testing"

	internaltesting "github.com/everettraven/crd-diff/pkg/validations/internal/testing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/ptr"
)

func TestMaximum(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			Flagged:              false,
			ComparableValidation: &Maximum{},
		},
		{
			Name: "new maximum constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			Flagged:              true,
			ComparableValidation: &Maximum{},
		},
		{
			Name: "maximum constraint decreased, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(20.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			Flagged:              true,
			ComparableValidation: &Maximum{},
		},
		{
			Name: "maximum constraint increased, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(20.0),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(30.0),
			},
			Flagged:              false,
			ComparableValidation: &Maximum{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &Maximum{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestMaxItems(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MaxItems{},
		},
		{
			Name: "new maxItems constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MaxItems{},
		},
		{
			Name: "maxItems constraint decreased, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MaxItems{},
		},
		{
			Name: "maxItems constraint increased, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(30)),
			},
			Flagged:              false,
			ComparableValidation: &MaxItems{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &MaxItems{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestMaxLength(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MaxLength{},
		},
		{
			Name: "new maxLength constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MaxLength{},
		},
		{
			Name: "maxLength constraint decreased, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MaxLength{},
		},
		{
			Name: "maxLength constraint increased, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(30)),
			},
			Flagged:              false,
			ComparableValidation: &MaxLength{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &MaxLength{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

func TestMaxProperties(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			Flagged:              false,
			ComparableValidation: &MaxProperties{},
		},
		{
			Name: "new maxProperties constraint, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MaxProperties{},
		},
		{
			Name: "maxProperties constraint decreased, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			Flagged:              true,
			ComparableValidation: &MaxProperties{},
		},
		{
			Name: "maxProperties constraint increased, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(20)),
			},
			New: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(30)),
			},
			Flagged:              false,
			ComparableValidation: &MaxProperties{},
		},
		{
			Name: "different field changed, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			Flagged:              false,
			ComparableValidation: &MaxProperties{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

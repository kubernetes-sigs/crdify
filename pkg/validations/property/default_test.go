package property

import (
	"testing"

	internaltesting "github.com/everettraven/crd-diff/pkg/validations/internal/testing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestDefault(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			Flagged:              false,
			ComparableValidation: &Default{},
		},
		{
			Name: "new default value, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			Flagged:              true,
			ComparableValidation: &Default{},
		},
		{
			Name: "default value removed, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			New:                  &apiextensionsv1.JSONSchemaProps{},
			Flagged:              true,
			ComparableValidation: &Default{},
		},
		{
			Name: "default value changed, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("bar"),
				},
			},
			Flagged:              true,
			ComparableValidation: &Default{},
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
			ComparableValidation: &Default{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

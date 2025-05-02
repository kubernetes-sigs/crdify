package property

import (
	"testing"

	internaltesting "github.com/everettraven/crd-diff/pkg/validations/internal/testing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestRequired(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged ",
			Old: &apiextensionsv1.JSONSchemaProps{
				Required: []string{
					"foo",
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Required: []string{
					"foo",
				},
			},
			Flagged:              false,
			ComparableValidation: &Required{},
		},
		{
			Name: "new required field, flagged",
			Old:  &apiextensionsv1.JSONSchemaProps{},
			New: &apiextensionsv1.JSONSchemaProps{
				Required: []string{
					"foo",
				},
			},
			Flagged:              true,
			ComparableValidation: &Required{},
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
			ComparableValidation: &Required{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

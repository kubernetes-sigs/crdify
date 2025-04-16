package property

import (
	"testing"

	internaltesting "github.com/everettraven/crd-diff/pkg/validations/internal/testing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestType(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Type: "string",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Type: "string",
			},
			Flagged:              false,
			ComparableValidation: &Type{},
		},
		{
			Name: "type changed, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Type: "string",
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Type: "integer",
			},
			Flagged:              true,
			ComparableValidation: &Type{},
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
			ComparableValidation: &Type{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

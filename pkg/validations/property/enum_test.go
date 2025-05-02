package property

import (
	"testing"

	internaltesting "github.com/everettraven/crd-diff/pkg/validations/internal/testing"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestEnum(t *testing.T) {
	testcases := []internaltesting.Testcase[apiextensionsv1.JSONSchemaProps]{
		{
			Name: "no diff, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &Enum{},
		},
		{
			Name: "new enum constraint, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			Flagged:              true,
			ComparableValidation: &Enum{},
		},
		{
			Name: "new allowed enum value added, addition policy not set, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			Flagged:              true,
			ComparableValidation: &Enum{},
		},
		{
			Name: "new allowed enum value added, addition policy set to Disallow, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			Flagged: true,
			ComparableValidation: &Enum{
				EnumConfig: EnumConfig{
					AdditionPolicy: AdditionPolicyDisallow,
				},
			},
		},
		{
			Name: "new allowed enum value added, addition policy set to Allow, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			Flagged: false,
			ComparableValidation: &Enum{
				EnumConfig: EnumConfig{
					AdditionPolicy: AdditionPolicyAllow,
				},
			},
		},
		{
			Name: "removed enum value, flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
					{
						Raw: []byte("bar"),
					},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("bar"),
					},
				},
			},
			Flagged:              true,
			ComparableValidation: &Enum{},
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
			ComparableValidation: &Enum{},
		},
		{
			Name: "different field changed with enum, not flagged",
			Old: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			New: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
				Enum: []apiextensionsv1.JSON{
					{
						Raw: []byte("foo"),
					},
				},
			},
			Flagged:              false,
			ComparableValidation: &Enum{},
		},
	}

	internaltesting.RunTestcases(t, testcases...)
}

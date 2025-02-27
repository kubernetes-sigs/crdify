package property

import (
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestDescription(t *testing.T) {
	type testcase struct {
		name        string
		oldProperty *apiextensionsv1.JSONSchemaProps
		newProperty *apiextensionsv1.JSONSchemaProps
		err         error
		handled     bool
		description *Description
	}

	for _, tc := range []testcase{
		{
			name:        "no description, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{},
			err:         nil,
			handled:     true,
			description: &Description{},
		},
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Description: "foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Description: "foo",
			},
			err:         nil,
			handled:     true,
			description: &Description{},
		},
		{
			name:        "new description, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Description: "new foo",
			},
			err:         nil,
			handled:     true,
			description: &Description{},
		},
		{
			name: "description removed, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Description: "old foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{},
			err:         nil,
			handled:     true,
			description: &Description{},
		},
		{
			name: "description changed, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Description: "old foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Description: "new foo",
			},
			err:         nil,
			handled:     true,
			description: &Description{},
		},
		{
			name: "different field changed, no error, not handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			err:         nil,
			handled:     false,
			description: &Description{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, handled, err := tc.description.Validate(NewDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

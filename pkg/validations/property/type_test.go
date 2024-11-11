package property

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestType(t *testing.T) {
	type testcase struct {
		name           string
		oldProperty    *apiextensionsv1.JSONSchemaProps
		newProperty    *apiextensionsv1.JSONSchemaProps
		err            error
		handled        bool
		typeValidation *Type
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Type: "string",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Type: "string",
			},
			err:            nil,
			handled:        true,
			typeValidation: &Type{},
		},
		{
			name: "type changed, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Type: "string",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Type: "integer",
			},
			err:            errors.New("type changed from \"string\" to \"integer\""),
			handled:        true,
			typeValidation: &Type{},
		},
		{
			name: "type changed, change enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Type: "string",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Type: "integer",
			},
			err:     nil,
			handled: true,
			typeValidation: &Type{
				ChangeEnforcement: TypeValidationChangeEnforcementNone,
			},
		},
		{
			name: "different field changed, no error, not handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			err:            nil,
			handled:        false,
			typeValidation: &Type{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, handled, err := tc.typeValidation.Validate(NewDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

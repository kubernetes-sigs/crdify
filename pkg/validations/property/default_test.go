package property

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestDefault(t *testing.T) {
	type testcase struct {
		name              string
		oldProperty       *apiextensionsv1.JSONSchemaProps
		newProperty       *apiextensionsv1.JSONSchemaProps
		err               error
		handled           bool
		defaultValidation *Default
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			err:               nil,
			handled:           true,
			defaultValidation: &Default{},
		},
		{
			name:        "new default value, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			err:               errors.New("default value \"foo\" added when there was no default previously"),
			handled:           true,
			defaultValidation: &Default{},
		},
		{
			name:        "new default value, addition enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			err:     nil,
			handled: true,
			defaultValidation: &Default{
				AdditionEnforcement: DefaultValidationAdditionEnforcementNone,
			},
		},
		{
			name: "default value removed, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			newProperty:       &apiextensionsv1.JSONSchemaProps{},
			err:               errors.New("default value \"foo\" removed"),
			handled:           true,
			defaultValidation: &Default{},
		},
		{
			name: "default value removed, removal enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{},
			err:         nil,
			handled:     true,
			defaultValidation: &Default{
				RemovalEnforcement: DefaultValidationRemovalEnforcementNone,
			},
		},
		{
			name: "default value changed, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("bar"),
				},
			},
			err:               errors.New("default value changed from \"foo\" to \"bar\""),
			handled:           true,
			defaultValidation: &Default{},
		},
		{
			name: "default value changed, change enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("foo"),
				},
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Default: &apiextensionsv1.JSON{
					Raw: []byte("bar"),
				},
			},
			err:     nil,
			handled: true,
			defaultValidation: &Default{
				ChangeEnforcement: DefaultValidationChangeEnforcementNone,
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
			err:               nil,
			handled:           false,
			defaultValidation: &Default{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handled, err := tc.defaultValidation.Validate(NewDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

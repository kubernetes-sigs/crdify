package property

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/ptr"
)

func TestMaximum(t *testing.T) {
	type testcase struct {
		name        string
		oldProperty *apiextensionsv1.JSONSchemaProps
		newProperty *apiextensionsv1.JSONSchemaProps
		err         error
		handled     bool
		maximum     *Maximum
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			err:     nil,
			handled: true,
			maximum: &Maximum{},
		},
		{
			name:        "new maximum constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			err:     errors.New("maximum: constraint 10 added when there were no restrictions previously"),
			handled: true,
			maximum: &Maximum{},
		},
		{
			name: "maximum constraint decreased, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(20.0),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(10.0),
			},
			err:     errors.New("maximum: constraint decreased from 20 to 10"),
			handled: true,
			maximum: &Maximum{},
		},
		{
			name: "maximum constraint increased, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(20.0),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Maximum: ptr.To(30.0),
			},
			err:     nil,
			handled: true,
			maximum: &Maximum{},
		},
		{
			name: "different field changed, no error, not handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			err:     nil,
			handled: false,
			maximum: &Maximum{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handled, err := tc.maximum.Validate(NewPropertyDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

func TestMaxItems(t *testing.T) {
	type testcase struct {
		name        string
		oldProperty *apiextensionsv1.JSONSchemaProps
		newProperty *apiextensionsv1.JSONSchemaProps
		err         error
		handled     bool
		maxItems    *MaxItems
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			err:      nil,
			handled:  true,
			maxItems: &MaxItems{},
		},
		{
			name:        "new maxItems constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			err:      errors.New("maxItems: constraint 10 added when there were no restrictions previously"),
			handled:  true,
			maxItems: &MaxItems{},
		},
		{
			name: "maxItems constraint decreased, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(20)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			err:      errors.New("maxItems: constraint decreased from 20 to 10"),
			handled:  true,
			maxItems: &MaxItems{},
		},
		{
			name: "maxitems constraint increased, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxItems: ptr.To(int64(20)),
			},
			err:      nil,
			handled:  true,
			maxItems: &MaxItems{},
		},
		{
			name: "different field changed, no error, not handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			err:      nil,
			handled:  false,
			maxItems: &MaxItems{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handled, err := tc.maxItems.Validate(NewPropertyDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

func TestMaxLength(t *testing.T) {
	type testcase struct {
		name        string
		oldProperty *apiextensionsv1.JSONSchemaProps
		newProperty *apiextensionsv1.JSONSchemaProps
		err         error
		handled     bool
		maxLength   *MaxLength
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			err:       nil,
			handled:   true,
			maxLength: &MaxLength{},
		},
		{
			name:        "new maxLength constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			err:       errors.New("maxLength: constraint 10 added when there were no restrictions previously"),
			handled:   true,
			maxLength: &MaxLength{},
		},
		{
			name: "maxLength constraint decreased, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(20)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			err:       errors.New("maxLength: constraint decreased from 20 to 10"),
			handled:   true,
			maxLength: &MaxLength{},
		},
		{
			name: "maxLength constraint increased, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxLength: ptr.To(int64(20)),
			},
			err:       nil,
			handled:   true,
			maxLength: &MaxLength{},
		},
		{
			name: "different field changed, no error, not handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			err:       nil,
			handled:   false,
			maxLength: &MaxLength{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handled, err := tc.maxLength.Validate(NewPropertyDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

func TestMaxProperties(t *testing.T) {
	type testcase struct {
		name          string
		oldProperty   *apiextensionsv1.JSONSchemaProps
		newProperty   *apiextensionsv1.JSONSchemaProps
		err           error
		handled       bool
		maxProperties *MaxProperties
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			err:           nil,
			handled:       true,
			maxProperties: &MaxProperties{},
		},
		{
			name:        "new maxProperties constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			err:           errors.New("maxProperties: constraint 10 added when there were no restrictions previously"),
			handled:       true,
			maxProperties: &MaxProperties{},
		},
		{
			name: "maxProperties constraint decreased, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(20)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			err:           errors.New("maxProperties: constraint decreased from 20 to 10"),
			handled:       true,
			maxProperties: &MaxProperties{},
		},
		{
			name: "maxProperties constraint increased, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MaxProperties: ptr.To(int64(20)),
			},
			err:           nil,
			handled:       true,
			maxProperties: &MaxProperties{},
		},
		{
			name: "different field changed, no error, not handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "foo",
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				ID: "bar",
			},
			err:           nil,
			handled:       false,
			maxProperties: &MaxProperties{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handled, err := tc.maxProperties.Validate(NewPropertyDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

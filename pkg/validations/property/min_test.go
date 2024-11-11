package property

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/ptr"
)

func TestMinimum(t *testing.T) {
	type testcase struct {
		name        string
		oldProperty *apiextensionsv1.JSONSchemaProps
		newProperty *apiextensionsv1.JSONSchemaProps
		err         error
		handled     bool
		minimum     *Minimum
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			err:     nil,
			handled: true,
			minimum: &Minimum{},
		},
		{
			name:        "new minimum constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			err:     errors.New("minimum: constraint 10 added when there were no restrictions previously"),
			handled: true,
			minimum: &Minimum{},
		},
		{
			name:        "new minimum constraint, addition enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			err:     nil,
			handled: true,
			minimum: &Minimum{
				MinOptions: MinOptions{
					AdditionEnforcement: MinVerificationAdditionEnforcementNone,
				},
			},
		},
		{
			name: "minimum constraint decreased, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(20.0),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			err:     nil,
			handled: true,
			minimum: &Minimum{},
		},
		{
			name: "minimum constraint increased, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(20.0),
			},
			err:     errors.New("minimum: constraint increased from 10 to 20"),
			handled: true,
			minimum: &Minimum{},
		},
		{
			name: "minimum constraint increased, increase enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(10.0),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				Minimum: ptr.To(20.0),
			},
			err:     nil,
			handled: true,
			minimum: &Minimum{
				MinOptions: MinOptions{
					IncreaseEnforcement: MinVerificationIncreaseEnforcementNone,
				},
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
			err:     nil,
			handled: false,
			minimum: &Minimum{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, handled, err := tc.minimum.Validate(NewDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

func TestMinLength(t *testing.T) {
	type testcase struct {
		name        string
		oldProperty *apiextensionsv1.JSONSchemaProps
		newProperty *apiextensionsv1.JSONSchemaProps
		err         error
		handled     bool
		minLength   *MinLength
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			err:       nil,
			handled:   true,
			minLength: &MinLength{},
		},
		{
			name:        "new minLength constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			err:       errors.New("minLength: constraint 10 added when there were no restrictions previously"),
			handled:   true,
			minLength: &MinLength{},
		},
		{
			name:        "new minLength constraint, addition enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			err:     nil,
			handled: true,
			minLength: &MinLength{
				MinOptions: MinOptions{
					AdditionEnforcement: MinVerificationAdditionEnforcementNone,
				},
			},
		},
		{
			name: "minLength constraint decreased, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(20)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			err:       nil,
			handled:   true,
			minLength: &MinLength{},
		},
		{
			name: "minLength constraint increased, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(20)),
			},
			err:       errors.New("minLength: constraint increased from 10 to 20"),
			handled:   true,
			minLength: &MinLength{},
		},
		{
			name: "minLength constraint increased, increase enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinLength: ptr.To(int64(20)),
			},
			err:     nil,
			handled: true,
			minLength: &MinLength{
				MinOptions: MinOptions{
					IncreaseEnforcement: MinVerificationIncreaseEnforcementNone,
				},
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
			err:       nil,
			handled:   false,
			minLength: &MinLength{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, handled, err := tc.minLength.Validate(NewDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

func TestMinItems(t *testing.T) {
	type testcase struct {
		name        string
		oldProperty *apiextensionsv1.JSONSchemaProps
		newProperty *apiextensionsv1.JSONSchemaProps
		err         error
		handled     bool
		minItems    *MinItems
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			err:      nil,
			handled:  true,
			minItems: &MinItems{},
		},
		{
			name:        "new minItems constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			err:      errors.New("minItems: constraint 10 added when there were no restrictions previously"),
			handled:  true,
			minItems: &MinItems{},
		},
		{
			name:        "new minItems constraint, addition enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			err:     nil,
			handled: true,
			minItems: &MinItems{
				MinOptions: MinOptions{
					AdditionEnforcement: MinVerificationAdditionEnforcementNone,
				},
			},
		},
		{
			name: "minItems constraint decreased, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(20)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			err:      nil,
			handled:  true,
			minItems: &MinItems{},
		},
		{
			name: "minItems constraint increased, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(20)),
			},
			err:      errors.New("minItems: constraint increased from 10 to 20"),
			handled:  true,
			minItems: &MinItems{},
		},
		{
			name: "minItems constraint increased, increase enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinItems: ptr.To(int64(20)),
			},
			err:     nil,
			handled: true,
			minItems: &MinItems{
				MinOptions: MinOptions{
					IncreaseEnforcement: MinVerificationIncreaseEnforcementNone,
				},
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
			err:      nil,
			handled:  false,
			minItems: &MinItems{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, handled, err := tc.minItems.Validate(NewDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

func TestMinProperties(t *testing.T) {
	type testcase struct {
		name          string
		oldProperty   *apiextensionsv1.JSONSchemaProps
		newProperty   *apiextensionsv1.JSONSchemaProps
		err           error
		handled       bool
		minProperties *MinProperties
	}

	for _, tc := range []testcase{
		{
			name: "no diff, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			err:           nil,
			handled:       true,
			minProperties: &MinProperties{},
		},
		{
			name:        "new minProperties constraint, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			err:           errors.New("minProperties: constraint 10 added when there were no restrictions previously"),
			handled:       true,
			minProperties: &MinProperties{},
		},
		{
			name:        "new minProperties constraint, addition enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			err:     nil,
			handled: true,
			minProperties: &MinProperties{
				MinOptions: MinOptions{
					AdditionEnforcement: MinVerificationAdditionEnforcementNone,
				},
			},
		},
		{
			name: "minProperties constraint decreased, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(20)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			err:           nil,
			handled:       true,
			minProperties: &MinProperties{},
		},
		{
			name: "minProperties constraint increased, error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(20)),
			},
			err:           errors.New("minProperties: constraint increased from 10 to 20"),
			handled:       true,
			minProperties: &MinProperties{},
		},
		{
			name: "minProperties constraint increased, increase enforcement set to None, no error, handled",
			oldProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(10)),
			},
			newProperty: &apiextensionsv1.JSONSchemaProps{
				MinProperties: ptr.To(int64(20)),
			},
			err:     nil,
			handled: true,
			minProperties: &MinProperties{
				MinOptions: MinOptions{
					IncreaseEnforcement: MinVerificationIncreaseEnforcementNone,
				},
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
			err:           nil,
			handled:       false,
			minProperties: &MinProperties{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, handled, err := tc.minProperties.Validate(NewDiff(tc.oldProperty, tc.newProperty))
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.handled, handled)
		})
	}
}

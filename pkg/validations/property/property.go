package property

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type PropertyDiff interface {
	Old() *apiextensionsv1.JSONSchemaProps
	New() *apiextensionsv1.JSONSchemaProps
}

func NewPropertyDiff(old, new *apiextensionsv1.JSONSchemaProps) PropertyDiff {
	return &propertyDiff{
		old: old,
		new: new,
	}
}

type propertyDiff struct {
	old *apiextensionsv1.JSONSchemaProps
	new *apiextensionsv1.JSONSchemaProps
}

func (pd *propertyDiff) Old() *apiextensionsv1.JSONSchemaProps {
	return pd.old.DeepCopy()
}

func (pd *propertyDiff) New() *apiextensionsv1.JSONSchemaProps {
	return pd.new.DeepCopy()
}

type PropertyValidation interface {
	Validate(PropertyDiff) (bool, error)
	Name() string
}

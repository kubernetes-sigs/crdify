package property

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Diff interface {
	Old() *apiextensionsv1.JSONSchemaProps
	New() *apiextensionsv1.JSONSchemaProps
}

func NewDiff(old, new *apiextensionsv1.JSONSchemaProps) Diff {
	return &diff{
		old: old,
		new: new,
	}
}

type diff struct {
	old *apiextensionsv1.JSONSchemaProps
	new *apiextensionsv1.JSONSchemaProps
}

func (pd *diff) Old() *apiextensionsv1.JSONSchemaProps {
	return pd.old.DeepCopy()
}

func (pd *diff) New() *apiextensionsv1.JSONSchemaProps {
	return pd.new.DeepCopy()
}

type Validation interface {
	Validate(Diff) (bool, error)
	Name() string
}

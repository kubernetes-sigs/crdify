package runner

import (
	"github.com/everettraven/crd-diff/pkg/validations"
	"github.com/everettraven/crd-diff/pkg/validations/crd/existingfieldremoval"
	"github.com/everettraven/crd-diff/pkg/validations/crd/scope"
	"github.com/everettraven/crd-diff/pkg/validations/crd/storedversionremoval"
	"github.com/everettraven/crd-diff/pkg/validations/property"
)

var defaultRegistry = validations.NewRegistry()

func init() {
	existingfieldremoval.Register(defaultRegistry)
	scope.Register(defaultRegistry)
	storedversionremoval.Register(defaultRegistry)
	property.RegisterDefault(defaultRegistry)
	property.RegisterEnum(defaultRegistry)
	property.RegisterMaximum(defaultRegistry)
	property.RegisterMaxItems(defaultRegistry)
	property.RegisterMaxLength(defaultRegistry)
	property.RegisterMaxProperties(defaultRegistry)
	property.RegisterMinimum(defaultRegistry)
	property.RegisterMinItems(defaultRegistry)
	property.RegisterMinLength(defaultRegistry)
	property.RegisterMinProperties(defaultRegistry)
	property.RegisterRequired(defaultRegistry)
	property.RegisterType(defaultRegistry)
}

// DefaultRegistry returns a pre-configured validations.Registry
func DefaultRegistry() validations.Registry {
	return defaultRegistry
}

package config

import (
	"github.com/everettraven/crd-diff/pkg/validations/property"
	"github.com/everettraven/crd-diff/pkg/validations/validators/crd"
	"github.com/everettraven/crd-diff/pkg/validations/validators/version"
)

var StrictConfig = Config{
	Checks: Checks{
		CRD: StrictCRDChecks,
		Version: VersionChecks{
			SameVersion:   StrictSameVersionChecks,
			ServedVersion: StrictServedVersionChecks,
		},
	},
}

var StrictCRDChecks = CRDChecks{
	Scope: CheckConfig{
		Enabled: true,
	},
	ExistingFieldRemoval: CheckConfig{
		Enabled: true,
	},
	StoredVersionRemoval: CheckConfig{
		Enabled: true,
	},
}

var StrictSameVersionChecks = SameVersionCheckConfig{
	CheckConfig: CheckConfig{
		Enabled: true,
	},
	VersionCheckConfig: StrictVersionCheckConfig,
}

var StrictServedVersionChecks = ServedVersionCheckConfig{
	CheckConfig: CheckConfig{
		Enabled: true,
	},
	IgnoreConversion:   false,
	VersionCheckConfig: StrictVersionCheckConfig,
}

var StrictVersionCheckConfig = VersionCheckConfig{
	UnhandledFailureMode: version.FailureModeClosed,
	PropertyCheckConfig:  StrictPropertyCheckConfig,
}

var StrictPropertyCheckConfig = PropertyCheckConfig{
	Enum: EnumCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	Default: DefaultCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	Required: RequiredCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	Type: TypeCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	Maximum: MaxCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	MaxItems: MaxCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	MaxProperties: MaxCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	MaxLength: MaxCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	Minimum: MinCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	MinItems: MinCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	MinProperties: MinCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	MinLength: MinCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
}

type Config struct {
	Checks Checks `yaml:"checks"`
}

type Checks struct {
	CRD     CRDChecks     `yaml:"crd"`
	Version VersionChecks `yaml:"version"`
}

type CRDChecks struct {
	Scope                CheckConfig `yaml:"scope"`
	ExistingFieldRemoval CheckConfig `yaml:"existingFieldRemoval"`
	StoredVersionRemoval CheckConfig `yaml:"storedVersionRemoval"`
}

type VersionChecks struct {
	SameVersion   SameVersionCheckConfig   `yaml:"sameVersion"`
	ServedVersion ServedVersionCheckConfig `yaml:"servedVersion"`
}

type SameVersionCheckConfig struct {
	CheckConfig
	VersionCheckConfig
}

type ServedVersionCheckConfig struct {
	CheckConfig
	VersionCheckConfig
	IgnoreConversion bool `yaml:"ignoreConversion"`
}

type VersionCheckConfig struct {
	PropertyCheckConfig
	UnhandledFailureMode version.FailureMode `yaml:"unhandledFailureMode"`
}

type PropertyCheckConfig struct {
	Enum          EnumCheckConfig     `yaml:"enum"`
	Default       DefaultCheckConfig  `yaml:"default"`
	Required      RequiredCheckConfig `yaml:"required"`
	Type          TypeCheckConfig     `yaml:"type"`
	Maximum       MaxCheckConfig      `yaml:"maximum"`
	MaxItems      MaxCheckConfig      `yaml:"maxItems"`
	MaxProperties MaxCheckConfig      `yaml:"maxProperties"`
	MaxLength     MaxCheckConfig      `yaml:"maxLength"`
	Minimum       MinCheckConfig      `yaml:"minimum"`
	MinItems      MinCheckConfig      `yaml:"minItems"`
	MinProperties MinCheckConfig      `yaml:"minProperties"`
	MinLength     MinCheckConfig      `yaml:"minLength"`
}

type CheckConfig struct {
	Enabled bool
}

type EnumCheckConfig struct {
	CheckConfig
}

type DefaultCheckConfig struct {
	CheckConfig
}

type RequiredCheckConfig struct {
	CheckConfig
}

type TypeCheckConfig struct {
	CheckConfig
}

type MaxCheckConfig struct {
	CheckConfig
}

type MinCheckConfig struct {
	CheckConfig
}

func ValidatorForConfig(cfg Config) *crd.Validator {
	validations := ValidationsForCRDChecks(cfg.Checks.CRD)
	validations = append(validations, VersionValidationForVersionChecks(cfg.Checks.Version))
	return crd.NewValidator(crd.WithValidations(validations...))
}

func ValidationsForCRDChecks(checks CRDChecks) []crd.Validation {
	validations := []crd.Validation{}
	if checks.Scope.Enabled {
		validations = append(validations, &crd.Scope{})
	}

	if checks.ExistingFieldRemoval.Enabled {
		validations = append(validations, &crd.ExistingFieldRemoval{})
	}

	if checks.StoredVersionRemoval.Enabled {
		validations = append(validations, &crd.StoredVersionRemoval{})
	}

	return validations
}

func VersionValidationForVersionChecks(checks VersionChecks) *version.Validator {
	return version.NewValidator(
		version.WithSameVersionConfig(
			SameVersionConfigForSameVersionCheckConfig(checks.SameVersion),
		),
		version.WithServedVersionConfig(
			ServedVersionConfigForServedVersionCheckConfig(checks.ServedVersion),
		),
	)
}

func SameVersionConfigForSameVersionCheckConfig(cfg SameVersionCheckConfig) version.SameVersionConfig {
	svc := version.SameVersionConfig{
		Skip:                 !cfg.Enabled,
		UnhandledFailureMode: cfg.UnhandledFailureMode,
		Validations:          PropertyValidationsForPropertyCheckConfig(cfg.PropertyCheckConfig),
	}

	return svc
}

func ServedVersionConfigForServedVersionCheckConfig(cfg ServedVersionCheckConfig) version.ServedVersionConfig {
	svc := version.ServedVersionConfig{
		Skip:                 !cfg.Enabled,
		UnhandledFailureMode: cfg.UnhandledFailureMode,
		Validations:          PropertyValidationsForPropertyCheckConfig(cfg.PropertyCheckConfig),
		IgnoreConversion:     cfg.IgnoreConversion,
	}

	return svc
}

func PropertyValidationsForPropertyCheckConfig(cfg PropertyCheckConfig) []property.Validation {
	validations := []property.Validation{}
	if cfg.Enum.Enabled {
		validations = append(validations, &property.Enum{})
	}

	if cfg.Default.Enabled {
		validations = append(validations, &property.Default{})
	}

	if cfg.Required.Enabled {
		validations = append(validations, &property.Required{})
	}

	if cfg.Type.Enabled {
		validations = append(validations, &property.Type{})
	}

	if cfg.Maximum.Enabled {
		validations = append(validations, &property.Maximum{})
	}

	if cfg.MaxItems.Enabled {
		validations = append(validations, &property.MaxItems{})
	}

	if cfg.MaxLength.Enabled {
		validations = append(validations, &property.MaxLength{})
	}

	if cfg.MaxProperties.Enabled {
		validations = append(validations, &property.MaxProperties{})
	}

	if cfg.Minimum.Enabled {
		validations = append(validations, &property.Minimum{})
	}

	if cfg.MinItems.Enabled {
		validations = append(validations, &property.MinItems{})
	}

	if cfg.MinLength.Enabled {
		validations = append(validations, &property.MinLength{})
	}

	if cfg.MinProperties.Enabled {
		validations = append(validations, &property.MinProperties{})
	}

	return validations
}

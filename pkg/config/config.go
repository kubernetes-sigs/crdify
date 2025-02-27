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
		RemovalEnforcement:  property.EnumValidationRemovalEnforcementStrict,
		AdditionEnforcement: property.EnumValidationAdditionEnforcementStrict,
	},
	Default: DefaultCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		ChangeEnforcement:   property.DefaultValidationChangeEnforcementStrict,
		RemovalEnforcement:  property.DefaultValidationRemovalEnforcementStrict,
		AdditionEnforcement: property.DefaultValidationAdditionEnforcementStrict,
	},
	Description: DescriptionCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
	},
	Required: RequiredCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		NewEnforcement: property.RequiredValidationNewEnforcementStrict,
	},
	Type: TypeCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		ChangeEnforcement: property.TypeValidationChangeEnforcementStrict,
	},
	Maximum: MaxCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		MaxOptions: property.MaxOptions{
			AdditionEnforcement: property.MaxVerificationAdditionEnforcementStrict,
			DecreaseEnforcement: property.MaxVerificationDecreaseEnforcementStrict,
		},
	},
	MaxItems: MaxCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		MaxOptions: property.MaxOptions{
			AdditionEnforcement: property.MaxVerificationAdditionEnforcementStrict,
			DecreaseEnforcement: property.MaxVerificationDecreaseEnforcementStrict,
		},
	},
	MaxProperties: MaxCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		MaxOptions: property.MaxOptions{
			AdditionEnforcement: property.MaxVerificationAdditionEnforcementStrict,
			DecreaseEnforcement: property.MaxVerificationDecreaseEnforcementStrict,
		},
	},
	MaxLength: MaxCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		MaxOptions: property.MaxOptions{
			AdditionEnforcement: property.MaxVerificationAdditionEnforcementStrict,
			DecreaseEnforcement: property.MaxVerificationDecreaseEnforcementStrict,
		},
	},
	Minimum: MinCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		MinOptions: property.MinOptions{
			AdditionEnforcement: property.MinVerificationAdditionEnforcementStrict,
			IncreaseEnforcement: property.MinVerificationIncreaseEnforcementStrict,
		},
	},
	MinItems: MinCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		MinOptions: property.MinOptions{
			AdditionEnforcement: property.MinVerificationAdditionEnforcementStrict,
			IncreaseEnforcement: property.MinVerificationIncreaseEnforcementStrict,
		},
	},
	MinProperties: MinCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		MinOptions: property.MinOptions{
			AdditionEnforcement: property.MinVerificationAdditionEnforcementStrict,
			IncreaseEnforcement: property.MinVerificationIncreaseEnforcementStrict,
		},
	},
	MinLength: MinCheckConfig{
		CheckConfig: CheckConfig{
			Enabled: true,
		},
		MinOptions: property.MinOptions{
			AdditionEnforcement: property.MinVerificationAdditionEnforcementStrict,
			IncreaseEnforcement: property.MinVerificationIncreaseEnforcementStrict,
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
	Enum          EnumCheckConfig        `yaml:"enum"`
	Default       DefaultCheckConfig     `yaml:"default"`
	Description   DescriptionCheckConfig `yaml:"description"`
	Required      RequiredCheckConfig    `yaml:"required"`
	Type          TypeCheckConfig        `yaml:"type"`
	Maximum       MaxCheckConfig         `yaml:"maximum"`
	MaxItems      MaxCheckConfig         `yaml:"maxItems"`
	MaxProperties MaxCheckConfig         `yaml:"maxProperties"`
	MaxLength     MaxCheckConfig         `yaml:"maxLength"`
	Minimum       MinCheckConfig         `yaml:"minimum"`
	MinItems      MinCheckConfig         `yaml:"minItems"`
	MinProperties MinCheckConfig         `yaml:"minProperties"`
	MinLength     MinCheckConfig         `yaml:"minLength"`
}

type CheckConfig struct {
	Enabled bool `json:"enabled"`
}

type EnumCheckConfig struct {
	CheckConfig
	RemovalEnforcement  property.EnumValidationRemovalEnforcement  `json:"removalEnforcement"`
	AdditionEnforcement property.EnumValidationAdditionEnforcement `json:"additionEnforcement"`
}

type DefaultCheckConfig struct {
	CheckConfig
	ChangeEnforcement   property.DefaultValidationChangeEnforcement   `json:"changeEnforcement"`
	RemovalEnforcement  property.DefaultValidationRemovalEnforcement  `json:"removalEnforcement"`
	AdditionEnforcement property.DefaultValidationAdditionEnforcement `json:"additionEnforcement"`
}

type DescriptionCheckConfig struct {
	CheckConfig
}

type RequiredCheckConfig struct {
	CheckConfig
	NewEnforcement property.RequiredValidationNewEnforcement `json:"newEnforcement"`
}

type TypeCheckConfig struct {
	CheckConfig
	ChangeEnforcement property.TypeValidationChangeEnforcement `json:"changeEnforcement"`
}

type MaxCheckConfig struct {
	CheckConfig
	property.MaxOptions
}

type MinCheckConfig struct {
	CheckConfig
	property.MinOptions
}

func ValidatorForConfig(cfg Config) *crd.Validator {
	validations := ValidationsForCRDChecks(cfg.Checks.CRD)
	validations = append(validations, VersionValidationForVersionChecks(cfg.Checks.Version))
	return crd.NewValidator(crd.WithValidations(validations...))
}

func ValidationsForCRDChecks(checks CRDChecks) []crd.Validation {
	var validations []crd.Validation
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
	var validations []property.Validation
	if cfg.Enum.Enabled {
		validations = append(validations, &property.Enum{
			RemovalEnforcement:  cfg.Enum.RemovalEnforcement,
			AdditionEnforcement: cfg.Enum.AdditionEnforcement,
		})
	}

	if cfg.Default.Enabled {
		validations = append(validations, &property.Default{
			ChangeEnforcement:   cfg.Default.ChangeEnforcement,
			RemovalEnforcement:  cfg.Default.RemovalEnforcement,
			AdditionEnforcement: cfg.Default.AdditionEnforcement,
		})
	}

	if cfg.Description.Enabled {
		validations = append(validations, &property.Description{})
	}

	if cfg.Required.Enabled {
		validations = append(validations, &property.Required{
			NewEnforcement: cfg.Required.NewEnforcement,
		})
	}

	if cfg.Type.Enabled {
		validations = append(validations, &property.Type{
			ChangeEnforcement: cfg.Type.ChangeEnforcement,
		})
	}

	if cfg.Maximum.Enabled {
		validations = append(validations, &property.Maximum{
			MaxOptions: cfg.Maximum.MaxOptions,
		})
	}

	if cfg.MaxItems.Enabled {
		validations = append(validations, &property.MaxItems{
			MaxOptions: cfg.MaxItems.MaxOptions,
		})
	}

	if cfg.MaxLength.Enabled {
		validations = append(validations, &property.MaxLength{
			MaxOptions: cfg.MaxLength.MaxOptions,
		})
	}

	if cfg.MaxProperties.Enabled {
		validations = append(validations, &property.MaxProperties{
			MaxOptions: cfg.MaxProperties.MaxOptions,
		})
	}

	if cfg.Minimum.Enabled {
		validations = append(validations, &property.Minimum{
			MinOptions: cfg.Minimum.MinOptions,
		})
	}

	if cfg.MinItems.Enabled {
		validations = append(validations, &property.MinItems{
			MinOptions: cfg.MinItems.MinOptions,
		})
	}

	if cfg.MinLength.Enabled {
		validations = append(validations, &property.MinLength{
			MinOptions: cfg.MinLength.MinOptions,
		})
	}

	if cfg.MinProperties.Enabled {
		validations = append(validations, &property.MinProperties{
			MinOptions: cfg.MinProperties.MinOptions,
		})
	}

	return validations
}

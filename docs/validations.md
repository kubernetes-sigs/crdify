# Validations

`crd-diff` has three different types of validations:

- Global validations - which are validations that are run against the entire CustomResourceDefinition
- Version validations - which are validations that are run against the versions that are specified in the CustomResourceDefinition
- Property validations - which are validations that are run against the schema properties in a given version-pair of the CustomResourceDefinition

## Global Validations

### Scope
Compares the old and new CustomResourceDefinitions to verify that their
scopes are the same. A CRD going from `Cluster` to `Namespace` scope (or vice-versa) is considered a breaking change.

In a configuration file, this validation can be disabled by setting

```yaml
checks:
  crd:
    scope:
      enabled: false
```

### ExistingFieldRemoval

Evaluates all versions of the old and new CustomResourceDefinitions to
verify that not existing fields have been removed from the CRD schemas. Removing an existing field means
that clients relying on that field will no longer be able to read or write to it and is considered a breaking change.

In a configuration file, this validation can be disabled by setting

```yaml
checks:
  crd:
    existingFieldRemoval:
      enabled: false
```

### StoredVersionRemoval

Ensures that no stored versions of the CustomResourceDefinition have been removed in the new CRD. This validation really only has significance when the old CRD is sourced from a
Kubernetes cluster with some CustomResources (CRs) already present on the cluster (i.e stored in etcd).
Kubernetes itself won't let you make a change where you drop a stored version because all existing stored
data _must_ be migrated to a newer version before the old version is removed.

In a configuration file, this validation can be disabled by setting

```yaml
checks:
  crd:
    storedVersionRemoval:
      enabled: false
```

## Version Validations

### SameVersion

Ensures compatibility between the same version of the old and new CustomResourceDefinition by running property validations with a same version version-pair (i.e `v1alpha1` from the old CRD and `v1alpha1` from the new CRD).

In a configuration file this validation can be:

- Enabled/Disabled
- Configured to pass/fail on any unhandled failures (an unhandled failure is a change to a property in which no property validation exists)
- Configured to run a specific subset of known property validations

Example configuration for this validation:

```yaml
checks:
  sameVersion:
    enabled: { true || false }
    # Open == Pass on unhandled failures, Closed == Fail on unhandled failures
    unhandledFailureMode: { Open || Closed }
    # property validation configurations go here
```

### ServedVersion

Ensures compatibility between all served versions of the new CustomResourceDefinition by running property validations with a served version-pair (i.e served version `v1alpha1` and served version `v1` from the new CRD).

In a configuration file this validation can be:

- Enabled/Disabled
- Configured to pass/fail on any unhandled failures (an unhandled failure is a change to a property in which no property validation exists)
- Configured to automatically pass on conversion strategy being specified.
- Configured to run a specific subset of known property validations

Example configuration for this validation:

```yaml
checks:
  versions:
    servedVersion:
      enabled: { true || false }
      # Open == Pass on unhandled failures, Closed == Fail on unhandled failures
      unhandledFailureMode: { Open || Closed }
      # Setting this to true means all property validations are run even if a conversion strategy is
      # specified. Setting this to false means if a conversion strategy is specified this validation
      # will automatically pass without running any property validations.
      ignoreConversion: { true || false }
      # property validation configurations go here
```
## Property Validations
### Enum

Validates compatibility of changes to enum constraints on a property.

Incompatible changes are generally:

- Adding enum constraints when there were none previously
- Removing a previously valid enum
- Adding a new enum value

Depending on the circumstances, adding/removing an enum may also be considered an incompatible change. APIs must set very prescriptive field descriptions to indicate
how clients should react to changes to allowed enum values.

In a configuration file, this validation can be configured to:

- Allow adding enum constraints when there were none previously
- Prevent adding new enum values
- Allow removal of a previously valid enum

Example configuration for this validation:

```yaml
checks:
  versions:
   # in this example we are configuring the sameVersion validation's enum property validation.
   # it is possible to configure this property validation separately in both the sameVersion and
   # servedVersion validations.
   sameVersion:
    enum:
      enabled: { true || false }
      # additionEnforcement is the enforcement strategy that should be used
      # when evaluating if the addition of an allowed enum value for a property
      # is considered incompatible.
      #
      # Known enforcement strategies are "Strict", "IfPreviouslyConstrained", and "None".
      #
      # When set to "Strict", addition of any new enum values for a property
      # is considered incompatible, including when there were previously no enum constraints.
      #
      # When set to "IfPreviouslyConstrained", addition any number of enum values for
      # a property when it was not previously constrained by enum values is considered incompatible.
      # Addition of enum values for a property, when it was previously constrained by enum values,
      # is considered compatible with this enforcement strategy
      #
      # When set to "None", addition of a new enum value for a property is
      # not considered incompatible.
      #
      # If set to an unknown value, the "NotPreviouslyConstrained" enforcement strategy
      # will be used.
      additionEnforcement: { Strict || IfPreviouslyConstrained || None }
      # removalEnforcement is the enforcement strategy that should be used
      # when evaluating if the removal of an allowed enum value for a property
      # is considered incompatible.
      #
      # Known enforcement strategies are "Strict" and "None".
      #
      # When set to "Strict", removal of an enum value for a property
      # is considered incompatible.
      #
      # When set to "None", removal of an enum value for a property is
      # not considered incompatible.
      #
      # If set to an unknown value, the "Strict" enforcement strategy
      # will be used.
      removalEnforcement: { Strict || None }
```

### Default

Validates compatibility of changes to a property's default value. Changes to default
values of properties may result in breaking expectations of clients and users that rely
on specific defaulting behaviors of an API.

Incompatible changes are generally:

- Removing the default value
- Changing the default value
- Adding a default value when one did not exist previously

Depending on how it is done, adding a default value may be considered a compatible
change. Generally the addition of a default value must be done such that the default value
is semantically equivalent to an empty value.


In a configuration file, this validation can be configured to:

- Be enabled/disabled
- Prevent changes to the default value
- Prevent the removal of the default value
- Prevent the addition of a default value when one did not previously exist

Example configuration for this validation:

```yaml
checks:
  version:
   # in this example we are configuring the sameVersion validation's enum property validation.
   # it is possible to configure this property validation separately in both the sameVersion and
   # servedVersion validations.
    sameVersion:
      default:
        enabled: { true || false }
        # changeEnforcement is the enforcement strategy that should be used
        # when evaluating if a change to the default value of a property
        # is considered incompatible.
        #
        # Known enforcement strategies are "Strict" and "None".
        #
        # When set to "Strict", any changes to the default value of a property
        # is considered incompatible.
        #
        # When set to "None", changes to the default value of a property are
        # not considered incompatible.
        #
        # If set to an unknown value, the "Strict" enforcement strategy
        # will be used.
        changeEnforcement: { Strict || None }
        # removalEnforcement is the enforcement strategy that should be used
        # when evaluating if the removal of the default value of a property
        # is considered incompatible.
        #
        # Known enforcement strategies are "Strict" and "None".
        #
        # When set to "Strict", removal of the default value of a property
        # is considered incompatible.
        #
        # When set to "None", removal of the default value of a property is
        # not considered incompatible.
        #
        # If set to an unknown value, the "Strict" enforcement strategy
        # will be used.
        removalEnforcement: { Strict || None }
        # AdditionEnforcement is the enforcement strategy that should be used
        # when evaluating if the addition of a default value for a property
        # is considered incompatible.
        #
        # Known enforcement strategies are "Strict" and "None".
        #
        # When set to "Strict", addition of a default value for a property
        # is considered incompatible.
        #
        # When set to "None", addition of a default value for a property is
        # not considered incompatible.
        #
        # If set to an unknown value, the "Strict" enforcement strategy
        # will be used.
        additionEnforcement: { Strict || None }
```

### Maximum, MaxLength, MaxItems, and MaxProperties

Validates compatibility of changes to the property constraints related to maximum
allowed values. There are individual validations for the following maximum value property constraints:

- Maximum
- MaxItems
- MaxLength
- MaxProperties

All of these individual validations are based on a common maximum validation and as such share
the same general configuration options and compatibility evaluations.

Generally, incompatible changes for these property constraints are:

- Adding a maximum value constraint when one did not exist previously
- Decreasing a maximum value constraint

As with other validations, these validations can be configured in a configuration file to:

- Be enabled/disabled
- Prevent the addition of a new maximum value constraint to a property
- Prevent the decreasing of a maximum value constraint

Example configuration for this validation:

```yaml
checks:
  version:
   # in this example we are configuring the sameVersion validation's enum property validation.
   # it is possible to configure this property validation separately in both the sameVersion and
   # servedVersion validations.
    sameVersion:
      # configuring the maximum validation here, but the same configuration options
      # exist for the maxItems, maxLength, and maxProperties validations
      maximum:
        enabled: { true || false }
        # additionEnforcement is the enforcement strategy to be used when
        # evaluating if adding a new maximum constraint for a CRD property
        # is considered incompatible.
        #
        # Known values are "Strict" and "None".
        #
        # When set to "Strict", adding a new maximum constraint to a CRD property
        # will be considered incompatible. Defaults to "Strict" when
        # unknown values are provided.
        #
        # When set to "None", adding a new maximum constraint to a CRD property
        # will not be considered an incompatible change.
        additionEnforcement: { Strict || None }
        # decreaseEnforcement is the enforcement strategy to be used when
        # evaluating if decreasing the maximum constraint for a CRD property
        # is considered incompatible.
        #
        # Known values are "Strict" and "None".
        #
        # When set to "Strict", decreasing a maximum constraint to a CRD property
        # will be considered incompatible. Defaults to "Strict" when
        # unknown values are provided.
        #
        # When set to "None", decreasing a maximum constraint for a CRD property
        # will not be considered an incompatible change.
        decreaseEnforcement: { Strict || None }
```

### Minimum, MinLength, MinItems, MinProperties

Validates compatibility of changes to the property constraints related to minimum
allowed values. There are individual validations for the following minimum value property constraints:

- Minimum
- MinItems
- MinLength
- MinProperties

All of these individual validations are based on a common minimum validation and as such share
the same general configuration options and compatibility evaluations.

Generally, incompatible changes for these property constraints are:

- Adding a minimum value constraint when one did not exist previously
- Increasing a minimum value constraint

As with other validations, these validations can be configured in a configuration file to:

- Be enabled/disabled
- Prevent the addition of a new minimum value constraint to a property
- Prevent the increasing of a minimum value constraint

Example configuration for this validation:

```yaml
checks:
  version:
   # in this example we are configuring the sameVersion validation's enum property validation.
   # it is possible to configure this property validation separately in both the sameVersion and
   # servedVersion validations.
    sameVersion:
      # configuring the minimum validation here, but the same configuration options
      # exist for the minItems, minLength, and minProperties validations
      minimum:
        enabled: { true || false }
        # additionEnforcement is the enforcement strategy to be used when
        # evaluating if adding a new minimum constraint for a CRD property
        # is considered incompatible.
        #
        # Known values are "Strict" and "None".
        #
        # When set to "Strict", adding a new minimum constraint to a CRD property
        # will be considered incompatible. Defaults to "Strict" when
        # unknown values are provided.
        #
        # When set to "None", adding a new minimum constraint to a CRD property
        # will not be considered an incompatible change.
        additionEnforcement: { Strict || None }
        # increaseEnforcement is the enforcement strategy to be used when
        # evaluating if increasing the minimum constraint for a CRD property
        # is considered incompatible.
        #
        # Known values are "Strict" and "None".
        #
        # When set to "Strict", increasing a minimum constraint to a CRD property
        # will be considered incompatible. Defaults to "Strict" when
        # unknown values are provided.
        #
        # When set to "None", increasing a minimum constraint for a CRD property
        # will not be considered an incompatible change.
        increaseEnforcement: { Strict || None }
```

### Required

Validates compatibility of required fields. It is generally considered an incompatible
change to add new required fields. Adding new required fields breaks client and user
expectations of the fields required for a version and makes existing resources stored in
etcd invalid.

As with the other validations, this validation can be configured to allow the addition of
a new required field and be enabled/disabled.

Example configuration for this validation:

```yaml
checks:
  version:
   # in this example we are configuring the sameVersion validation's enum property validation.
   # it is possible to configure this property validation separately in both the sameVersion and
   # servedVersion validations.
    sameVersion:
      required:
        enabled: { true || false }
        # newEnforcement is the enforcement strategy to use when
        # evaluating if adding a new required field to a CRD
        # property is considered an incompatible change.
        #
        # Known values are "Strict" and "None".
        #
        # When set to "Strict", adding a new required field
        # will be considered an incompatible change. "Strict"
        # is the default enforcement strategy and is used when
        # unknown values are specified.
        #
        # When set to "None", adding a new required field
        # will not be considered an incompatible change
        newEnforcement: { Strict || None }
```

### Type

Validates compatibility of property types. It is generally considered a breaking change
to change the type of a property as it breaks client/user expectations and makes existing
stored instances of the resource invalid.

As with the other validations, this validation can be configured to be enabled/disabled
and to prevent/allow changes to a property's type.

Example configuration for this validation:

```yaml
checks:
  version:
   # in this example we are configuring the sameVersion validation's enum property validation.
   # it is possible to configure this property validation separately in both the sameVersion and
   # servedVersion validations.
    sameVersion:
      type:
        enabled: { true || false }
        # changeEnforcement is the enforcement strategy to be used
        # when evaluating if a change to the type of a CRD property
        # is considered an incompatible change
        #
        # Known values are "Strict" and "None".
        #
        # When set to "Strict", changes to the type of a CRD property
        # are considered incompatible changes. "Strict" is the default
        # enforcement strategy and is used when unknown values are specified.
        #
        # When set to "None", changes to the type of a CRD property
        # are not considered incompatible changes.
        changeEnforcement: { Strict || None }
```

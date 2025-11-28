# Validations

`crdify` has three different "validators" that run:

- CRD Validator - runs validations, known as "Global Validations", that verify compatibility of changes across the entire CRD.
- Same Version Validator - runs validations, known as "Property Validations", that verify compatibility of changes to properties within the same versions.
- Served Version Validator - runs validations, known as "Property Validations", that verify compatibility of changes to properties across served versions.

## Global Validations

### scope

Compares the old and new CustomResourceDefinitions to verify that their
scopes are the same. A CRD going from `Cluster` to `Namespace` scope (or vice-versa) is considered a breaking change.

### existingFieldRemoval

Evaluates all versions of the old and new CustomResourceDefinitions to
verify that not existing fields have been removed from the CRD schemas. Removing an existing field means
that clients relying on that field will no longer be able to read or write to it and is considered a breaking change.

### storedVersionRemoval

Ensures that no stored versions of the CustomResourceDefinition have been removed in the new CRD. This validation really only has significance when the old CRD is sourced from a
Kubernetes cluster with some CustomResources (CRs) already present on the cluster (i.e stored in etcd).
Kubernetes itself won't let you make a change where you drop a stored version because all existing stored
data _must_ be migrated to a newer version before the old version is removed.

## Property Validations

### enum

Validates compatibility of changes to enum constraints on a property.

Incompatible changes are:

- Adding enum constraints when there were none previously
- Removing a previously valid enum
- Adding a new enum value

Depending on the circumstances, adding an enum _may_ be considered a compatible change. APIs must set very prescriptive field descriptions to indicate
how clients should react to changes to allowed enum values.

#### Configuration

The `enum` validation has unique configuration options that can be used to change how it determines compatibility of a change to enum constraints on a property:

- `additionPolicy` - used to configure how compatibility is determined when adding new allowed enums to an existing set of enum constraints. Allowed values are `Allow` and `Disallow`. When set to `Allow`, adding a new enum value is considered a compatible change. When set to `Disallow`, adding a new enum value is considered an incompatible change. The default is `Disallow`.

An example of configuring the `enum` validation to allow adding a new enum value:

```yaml
validations:
  - name: enum
    enforcement: Error
    configuration:
      additionPolicy: Allow
```

### default

Validates compatibility of changes to a property's default value. Changes to default
values of properties may result in breaking expectations of clients and users that rely
on specific defaulting behaviors of an API.

Incompatible changes are:

- Removing the default value
- Changing the default value
- Adding a default value when one did not exist previously

### maximum, maxLength, maxItems, and maxProperties

Validates compatibility of changes to the property constraints related to maximum
allowed values. There are individual validations for the following maximum value property constraints:

- Maximum
- MaxItems
- MaxLength
- MaxProperties

All of these individual validations are based on a common maximum validation and as such share
the same general configuration options and compatibility evaluations.

Incompatible changes for these property constraints are:

- Adding a maximum value constraint when one did not exist previously
- Decreasing a maximum value constraint

### minimum, minLength, minItems, minProperties

Validates compatibility of changes to the property constraints related to minimum
allowed values. There are individual validations for the following minimum value property constraints:

- Minimum
- MinItems
- MinLength
- MinProperties

All of these individual validations are based on a common minimum validation and as such share
the same general configuration options and compatibility evaluations.

Incompatible changes for these property constraints are:

- Adding a minimum value constraint when one did not exist previously
- Increasing a minimum value constraint

### required

Validates compatibility of required fields. It is an incompatible
change to add new required fields. Adding new required fields breaks client and user
expectations of the fields required for a version and makes existing resources stored in
etcd invalid.

### type

Validates compatibility of property types. It is considered a breaking change
to change the type of a property as it breaks client/user expectations and makes existing
stored instances of the resource invalid.

### description

Validates compatibility of changes to a property description. While most changes to the
description of a property are _generally_ safe, it is important to note that changing
the semantics of a field _is_ a breaking change as it breaks expectations clients/users
have made about what configuring the property does.

### format

Validates compatibility of changes to a property's format. Changing the format of a property
generally results in a change to the validation performed on the property values and is a
breaking change for clients/users existing expectations based on the previous format. 

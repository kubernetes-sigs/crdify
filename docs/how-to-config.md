
# Documentation of `config.yaml`

This documents each configurable option in `config.yaml`, including enforcement levels and functionality. Each option enforces different constraints on CRD compatibility, with varying levels of strictness.

## `checks` Section

The configuration file consists of two main sections under `checks`: **crd** and **version**. Both sections contain options to control validation logic.

---

### 1. `crd` Section

#### **scope**
- **enabled**: Enables scope checks for CRDs. Scope, defined as either `Namespaced` or `Cluster`, affects the CRD’s reach within the cluster. Changes to scope can lead to incompatibility issues if not validated.

#### **existingFieldRemoval**
- **enabled**: Ensures fields in older versions of the CRD schema are not removed in newer versions. This avoids potential issues where existing resources depend on removed fields.

#### **storedVersionRemoval**
- **enabled**: Checks that stored versions of the CRD remain accessible across upgrades. Removing stored versions can lead to data unavailability.

---

### 2. `version` Section

This section includes schema validation rules for each versioned field within CRDs, covering properties and constraints.

#### **sameVersion**
- **enabled**: Enables compatibility checks within the same CRD version.
- **unhandledFailureMode**: Specifies handling for unexpected validation failures. `"Closed"` indicates strict enforcement, where failures are considered critical.

##### Property-Specific Validation Rules under `sameVersion`

Each rule targets a specific property within CRD schemas.

- **enum**
  - **enabled**: Enables validation for enumerated values within fields, ensuring consistent allowed values.
  - **removalEnforcement**: Specifies the strategy for handling removed enum values. `"Strict"` means any removal will be flagged as incompatible.
  - **additionEnforcement**: Governs the addition of new enum values, with `"Strict"` indicating added values must align with compatibility guidelines.

- **default**
  - **enabled**: Validates default values for fields, ensuring default behaviors remain stable.
  - **changeEnforcement**: Controls how changes to default values are handled. `"Strict"` flags any change as incompatible.
  - **removalEnforcement**: Handles cases where default values are removed, with `"Strict"` indicating removal is flagged as incompatible.
  - **additionEnforcement**: Governs adding default values, with `"Strict"` enforcing compatibility checks.

- **required**
  - **enabled**: Checks for new fields marked as `required` to prevent unexpected changes in client expectations.

- **type**
  - **enabled**: Enforces compatibility in data types (e.g., `string` vs. `integer`) to prevent schema mismatches.

- **maximum**
  - **enabled**: Validates `maximum` constraints in numerical fields, ensuring compatibility with set limits.

- **maxItems**
  - **enabled**: Enforces `maxItems` constraints in array fields, ensuring that maximum item counts are not altered in a way that breaks compatibility.

- **maxProperties**
  - **enabled**: Validates the maximum number of properties in `object` fields, enforcing consistency in structure.

- **maxLength**
  - **enabled**: Checks that maximum lengths in string fields remain within compatible limits.

- **minimum**
  - **enabled**: Enforces minimum constraints on numerical fields.

- **minItems**
  - **enabled**: Enforces `minItems` constraints in array fields, ensuring consistency in required array lengths.

- **minProperties**
  - **enabled**: Checks for minimum property constraints in `object` fields.

- **minLength**
  - **enabled**: Enforces minimum length constraints in string fields, ensuring shorter values are not introduced unexpectedly.

#### **servedVersion**
- **enabled**: Ensures compatibility for served CRD versions available in the Kubernetes API.
- **unhandledFailureMode**: Controls response to failures, with `"Closed"` indicating strict enforcement.
- **ignoreConversion**: (Boolean) If `true`, skips conversion checks, allowing structural differences if versions are convertible.

##### Property-Specific Validation Rules under `servedVersion`

These rules follow the same logic as `sameVersion`, validating properties to ensure compatibility in served versions.

- **enum**
- **default**
- **required**
- **type**
- **maximum**
- **maxItems**
- **maxProperties**
- **maxLength**
- **minimum**
- **minItems**
- **minProperties**
- **minLength**

---

### Usage

To apply these configurations, define a `config.yaml` with the above options to meet your CRD compatibility requirements. Adjust enforcement levels based on your specific needs for compatibility and backward support.

For further customization, refer to the codebase and `pkg/validations/property` for enforcement details.

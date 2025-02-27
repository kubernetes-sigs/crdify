# Configuration
`crd-diff` can be customized using a YAML configuration file.

By default, `crd-diff` uses this configuration:

```yaml
checks:
  crd:
    scope:
      enabled: true
    existingFieldRemoval:
      enabled: true
    storedVersionRemoval:
      enabled: true
  version:
    sameVersion:
      enabled: true
      unhandledFailureMode: "Closed"
      enum:
        enabled: true
        removalEnforcement: Strict
        additionEnforcement: Strict
      default:
        enabled: true
        changeEnforcement: Strict
        removalEnforcement: Strict
        additionEnforcement: Strict
      description:
        enabled: true
      required:
        enabled: true
        newEnforcement: Strict
      type:
        enabled: true
        changeEnforcement: Strict
      maximum:
        enabled: true
        additionEnforcement: Strict
        decreaseEnforcement: Strict
      maxItems:
        enabled: true
        additionEnforcement: Strict
        decreaseEnforcement: Strict
      maxProperties:
        enabled: true
        additionEnforcement: Strict
        decreaseEnforcement: Strict
      maxLength:
        enabled: true
        additionEnforcement: Strict
        decreaseEnforcement: Strict
      minimum:
        enabled: true
        additionEnforcement: Strict
        increaseEnforcement: Strict
      minItems:
        enabled: true
        additionEnforcement: Strict
        increaseEnforcement: Strict
      minProperties:
        enabled: true
        additionEnforcement: Strict
        increaseEnforcement: Strict
      minLength:
        enabled: true
        additionEnforcement: Strict
        increaseEnforcement: Strict
    servedVersion:
      enabled: true
      unhandledFailureMode: "Closed"
      ignoreConversion: false
      enum:
        enabled: true
        removalEnforcement: Strict
        additionEnforcement: Strict
      default:
        enabled: true
        changeEnforcement: Strict
        removalEnforcement: Strict
        additionEnforcement: Strict
      description:
        enabled: true
      required:
        enabled: true
        newEnforcement: Strict
      type:
        enabled: true
        changeEnforcement: Strict
      maximum:
        enabled: true
        additionEnforcement: Strict
        decreaseEnforcement: Strict
      maxItems:
        enabled: true
        additionEnforcement: Strict
        decreaseEnforcement: Strict
      maxProperties:
        enabled: true
        additionEnforcement: Strict
        decreaseEnforcement: Strict
      maxLength:
        enabled: true
        additionEnforcement: Strict
        decreaseEnforcement: Strict
      minimum:
        enabled: true
        additionEnforcement: Strict
        increaseEnforcement: Strict
      minItems:
        enabled: true
        additionEnforcement: Strict
        increaseEnforcement: Strict
      minProperties:
        enabled: true
        additionEnforcement: Strict
        increaseEnforcement: Strict
      minLength:
        enabled: true
        additionEnforcement: Strict
        increaseEnforcement: Strict
```

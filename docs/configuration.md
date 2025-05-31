# Configuration
`crdify` can be customized using a YAML configuration file.

An example of configuring a validation via the configuration file:

```yaml
validations:
  - name: someValidation # the name of the validation you wish to configure
    enforcement: Error # the level of enforcement you'd like for the validation. Options are Error, Warn, and None.
    configuration: # additional configuration options unique to each validation.
        foo: bar
```

By default all validations run with the `Error` enforcement policy and may be configured to use the `Warn` or `None` policies instead.
An enforcement policy of `Error` means that if the validation detects an incompatible change it will manifest as an error and the program will exit with a non-zero exit code.
An enforcement policy of `Warn` means that if the validation detects an incompatible change it will manifest as a warning and the program will exit with a zero exit code.
An enforcement policy of `None` means that if the validation detects an incompatible change it will be ignored and the program will exit with a zero exit code.

Each validation may have additional defaults for their individual configuration options.

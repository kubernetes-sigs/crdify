# Getting Started

## Installing `crd-diff`
Currently, the only way to install the `crd-diff` tool is to use the `go install` command:

```sh
go install github.com/everettraven/crd-diff@{revision}
```

Replace `{revision}` with a tag, commit, or `latest` to build and install the tool from source at that particular revision.

## General Usage

```sh
crd-diff is a tool for evaluating changes to Kubernetes CustomResourceDefinitions
to help cluster administrators, gitops practitioners, and Kubernetes extension developers identify
changes that might result in a negative impact to clusters and/or users.

Example use cases:
    Evaluating a change in a CustomResourceDefinition on a Kubernetes Cluster with one in a file:
        $ crd-diff kube://{crd-name} file://{filepath}

    Evaluating a change from file to file:
        $ crd-diff file://{filepath} file://{filepath}

    Evaluating a change from git ref to git ref:
            $ crd-diff git://{ref}?path={filepath} git://{ref}?path={filepath}

Usage:
  crd-diff <old> <new> [flags]
  crd-diff [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     installed version of crd-diff

Flags:
      --config string   the filepath to load the check configurations from
  -h, --help            help for crd-diff
  -o, --output string   the format the output should take when incompatibilities are identified. May be one of plaintext, json, yaml (default "plaintext")

Use "crd-diff [command] --help" for more information about a command.
```

The `<old>` and `<new>` arguments are required and should be the sourcing information for the old and new
`CustomResourceDefinition` YAML

The supported sources are:

- `kube://{name}`
- `git://{ref}?path={filepath}`
- `file://{filepath}`

An example of using `crd-diff` to compare a `CustomResourceDefinition` on a Kubernetes cluster to the same one in a local file:

```sh
crd-diff kube://memcacheds.cache.example.com file://crd.yaml
```

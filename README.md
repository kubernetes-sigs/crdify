# crd-diff
`crd-diff` is a CLI tool for comparing Kubernetes `CustomResourceDefinition` resources (CRDs) for differences.
It checks for incompatible changes to help:
- Cluster administrators protect CRDs on their clusters from breaking changes
- GitOps practitioners prevent CRDs with breaking changes being committed
- Developers of Kubernetes extension identify when changes to their CRDs are incompatible

## Usage
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

## Installation

`crd-diff` can be installed by running:
```sh
go install github.com/everettraven/crd-diff@latest
```

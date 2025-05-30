# Getting Started

## Installing `crdify`
Currently, the only way to install the `crdify` tool is to use the `go install` command:

```sh
go install sigs.k8s.io/crdify@{revision}
```

Replace `{revision}` with a tag, commit, or `latest` to build and install the tool from source at that particular revision.

## General Usage

```sh
crdify is a tool for evaluating changes to Kubernetes CustomResourceDefinitions
to help cluster administrators, gitops practitioners, and Kubernetes extension developers identify
changes that might result in a negative impact to clusters and/or users.

Example use cases:
    Evaluating a change in a CustomResourceDefinition on a Kubernetes Cluster with one in a file:
        $ crdify kube://{crd-name} file://{filepath}

    Evaluating a change from file to file:
        $ crdify file://{filepath} file://{filepath}

    Evaluating a change from git ref to git ref:
            $ crdify git://{ref}?path={filepath} git://{ref}?path={filepath}

Usage:
  crdify <old> <new> [flags]
  crdify [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     installed version of crdify

Flags:
      --config string   the filepath to load the check configurations from
  -h, --help            help for crdify
  -o, --output string   the format the output should take when incompatibilities are identified. May be one of plaintext, json, yaml (default "plaintext")

Use "crdify [command] --help" for more information about a command.
```

The `<old>` and `<new>` arguments are required and should be the sourcing information for the old and new
`CustomResourceDefinition` YAML

The supported sources are:

- `kube://{name}`
- `git://{ref}?path={filepath}`
- `file://{filepath}`

An example of using `crdify` to compare a `CustomResourceDefinition` on a Kubernetes cluster to the same one in a local file:

```sh
crdify kube://memcacheds.cache.example.com file://crd.yaml
```

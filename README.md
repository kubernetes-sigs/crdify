# crd-diff
`crd-diff` is a CLI tool for comparing Kubernetes `CustomResourceDefinition` resources (CRDs) for differences.
It checks for incompatible changes to help:
- Cluster administrators protect CRDs on their clusters from breaking changes
- GitOps practitioners prevent CRDs with breaking changes being committed
- Developers of Kubernetes extension identify when changes to their CRDs are incompatible

# Usage
```sh
Usage:
  crd-diff <old> <new> [flags]

Flags:
  -h, --help   help for crd-diff
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

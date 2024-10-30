# crd-diff
`crd-diff` is a CLI tool for comparing Kubernetes `CustomResourceDefinition` resources (CRDs) for differences.
It checks for incompatible changes to help:
- Cluster administrators protect CRDs on their clusters from breaking changes
- GitOps practitioners prevent CRDs with breaking changes being committed
- Developers of Kubernetes extension identify when changes to their CRDs are incompatible


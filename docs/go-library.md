# Go Library

crdify may be used as a CLI or as a Go library. For an example of using crdify as a library, see the
[OLM v1 preflight checks for CRD compatibility](https://github.com/operator-framework/operator-controller/blob/f17f3c5728d330511f970ce501b709321ba54b09/internal/operator-controller/rukpak/preflights/crdupgradesafety/crdupgradesafety.go#L67-L107).

To communicate library changes, crdify adheres to the [go package version numbers convention](https://go.dev/doc/modules/version-numbers).
Keep in mind that, while we will try to keep the API as stable as possible, v0.x.y versions may introduce breaking changes.

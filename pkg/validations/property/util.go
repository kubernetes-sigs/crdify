package property

import (
	"k8s.io/apimachinery/pkg/api/equality"
)

type ResetFunc func(diff PropertyDiff) PropertyDiff

func IsHandled(diff PropertyDiff, reset ResetFunc) bool {
	resetDiff := reset(diff)
	return equality.Semantic.DeepEqual(resetDiff.Old(), resetDiff.New())
}

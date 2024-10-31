package property

import (
	"k8s.io/apimachinery/pkg/api/equality"
)

type ResetFunc func(diff Diff) Diff

func IsHandled(diff Diff, reset ResetFunc) bool {
	resetDiff := reset(diff)
	return equality.Semantic.DeepEqual(resetDiff.Old(), resetDiff.New())
}

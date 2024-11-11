package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/everettraven/crd-diff/pkg/validations/validators/crd"
	"sigs.k8s.io/yaml"
)

var _ crd.ValidationResult = (*Result)(nil)

type Result struct {
	Validation           string                 `json:"validation"`
	SameVersionResults   []VersionCompareResult `json:"sameVersion,omitempty"`
	ServedVersionResults []VersionCompareResult `json:"servedVersion,omitempty"`
}

func (r *Result) Error(printDepth int) error {
	sameVersionErrs := []error{}
	for _, vcr := range r.SameVersionResults {
		sameVersionErrs = append(sameVersionErrs, vcr.Error(printDepth+1))
	}
	var sameVersionErr error
	if errors.Join(sameVersionErrs...) != nil {
		var out strings.Builder
		out.WriteString(strings.Repeat("\t", printDepth))
		out.WriteString("comparing same versions:\n")
		for _, err := range sameVersionErrs {
			if err != nil {
				out.WriteString(strings.Repeat("\t", printDepth+1))
				out.WriteString(err.Error())
				out.WriteString("\n")
			}
		}
		sameVersionErr = errors.New(out.String())
	}

	servedVersionErrs := []error{}
	for _, vcr := range r.ServedVersionResults {
		servedVersionErrs = append(servedVersionErrs, vcr.Error(printDepth+1))
	}
	var servedVersionErr error
	if errors.Join(servedVersionErrs...) != nil {
		var out strings.Builder
		out.WriteString("comparing served versions:\n")
		for _, err := range servedVersionErrs {
			if err != nil {
				out.WriteString(strings.Repeat("\t", printDepth+1))
				out.WriteString(err.Error())
				out.WriteString("\n")
			}
		}
		servedVersionErr = errors.New(out.String())
	}

	return errors.Join(sameVersionErr, servedVersionErr)
}

func (r *Result) JSON() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Result) YAML() ([]byte, error) {
	return yaml.Marshal(r)
}

type VersionCompareResult struct {
	VersionA               string                  `json:"versionA"`
	VersionB               string                  `json:"versionB"`
	PropertyCompareResults []PropertyCompareResult `json:"propertyComparisons"`
}

func (vcr *VersionCompareResult) Error(printDepth int) error {
	propertyErrs := []error{}
	for _, pcr := range vcr.PropertyCompareResults {
		propertyErrs = append(propertyErrs, pcr.Error(printDepth+1))
	}
	if errors.Join(propertyErrs...) != nil {
		var out strings.Builder
		out.WriteString(fmt.Sprintf("comparing version %q with version %q:\n", vcr.VersionA, vcr.VersionB))

		for _, err := range propertyErrs {
			if err != nil {
				out.WriteString(strings.Repeat("\t", printDepth+1))
				out.WriteString(err.Error())
				out.WriteString("\n")
			}
		}
		return errors.New(out.String())
	}
	return nil
}

type PropertyCompareResult struct {
	Property string   `json:"property"`
	Errors   []string `json:"errors"`
}

func (pcr *PropertyCompareResult) Error(printDepth int) error {
	if len(pcr.Errors) > 0 {
		var out strings.Builder
		out.WriteString(fmt.Sprintf("comparing property %q:\n", pcr.Property))

		for _, err := range pcr.Errors {
			out.WriteString(strings.Repeat("\t", printDepth+1))
			out.WriteString(err)
			out.WriteString("\n")
		}
		return errors.New(out.String())
	}
	return nil
}

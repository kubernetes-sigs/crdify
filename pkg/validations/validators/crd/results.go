package crd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"
)

type ValidationResult interface {
	Error(printDepth int) error
	JSON() ([]byte, error)
	YAML() ([]byte, error)
}

type validationResult struct {
	Validation string `json:"validation"`
	Err        string `json:"error,omitempty"`
}

func (vr *validationResult) Error(printDepth int) error {
	if vr.Err != "" {
		errStr := fmt.Sprintf("%s validation failed: %s", vr.Validation, vr.Err)
		return errors.New(strings.Join([]string{
			strings.Repeat("\t", printDepth),
			errStr,
		}, ""))
	}
	return nil
}

func (vr *validationResult) JSON() ([]byte, error) {
	return json.Marshal(vr)
}

func (vr *validationResult) YAML() ([]byte, error) {
	return yaml.Marshal(vr)
}

type ValidatorResult struct {
	ValidationResults []ValidationResult `json:"results"`
}

func (vr *ValidatorResult) Error(printDepth int) error {
	validationErrors := []error{}
	for _, validationResult := range vr.ValidationResults {
		validationErrors = append(validationErrors, validationResult.Error(printDepth+1))
	}

	return errors.Join(validationErrors...)
}

func (vr *ValidatorResult) JSON() ([]byte, error) {
	return json.Marshal(vr)
}

func (vr *ValidatorResult) YAML() ([]byte, error) {
	return yaml.Marshal(vr)
}

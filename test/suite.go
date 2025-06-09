// Copyright 2025 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/crdify/pkg/runner"
	"sigs.k8s.io/crdify/pkg/validations"
)

var (
	errMismatchedOutput = errors.New("output does not match expected")
	errMissingTestFiles = errors.New("missing expected test files")
)

func main() {
	if err := newTestCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func newTestCommand() *cobra.Command {
	var (
		update  bool
		binary  string
		testDir string
	)

	command := &cobra.Command{
		Use:   "test [options]",
		Short: "runs the crdify e2e test suite",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTests(binary, testDir, update)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	command.Flags().BoolVar(&update, "update", false, "whether or not to update the expected output golden files")
	command.Flags().StringVar(&binary, "binary", "./bin/crdify", "location of the crdify binary to execute")
	command.Flags().StringVar(&testDir, "test-dir", "./test", "location of the tests to execute")

	return command
}

func runTests(binary string, testDir string, update bool) error {
	tests, err := testsForDirectory(testDir)
	if err != nil {
		return err
	}

	errs := []error{}

	for _, test := range tests {
		fmt.Println("running test", test.name, "source A", test.sourceA, "source B", test.sourceB, "expected", test.expected, "update?", update, "binary", binary)
		err := executeTest(test, binary, update)
		if err != nil {
			errs = append(errs, fmt.Errorf("executing test %q: %w", test.name, err))
		}
	}

	return errors.Join(errs...)
}

type test struct {
	name     string
	sourceA  string
	sourceB  string
	expected string
}

const (
	sourceAName  = "a.yaml"
	sourceBName  = "b.yaml"
	expectedName = "expected.json"
)

func testsForDirectory(testDir string) ([]test, error) {
	tests := []test{}

	dirEntries, err := os.ReadDir(testDir)
	if err != nil {
		return nil, fmt.Errorf("reading test directory %q: %w", testDir, err)
	}

	errs := []error{}

	for _, entry := range dirEntries {
		// all tests will have their own subdirectory
		if !entry.IsDir() {
			continue
		}

		base := filepath.Join(testDir, entry.Name())
		sourceA := filepath.Join(base, sourceAName)
		sourceB := filepath.Join(base, sourceBName)
		expected := filepath.Join(base, expectedName)

		if !hasFile(sourceA) || !hasFile(sourceB) || !hasFile(expected) {
			errs = append(errs, fmt.Errorf("test %q : %w", entry.Name(), errMissingTestFiles))
			continue
		}

		tests = append(tests, test{
			name:     entry.Name(),
			sourceA:  sourceA,
			sourceB:  sourceB,
			expected: expected,
		})
	}

	return tests, errors.Join(errs...)
}

func hasFile(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func executeTest(t test, binary string, update bool) error {
	//nolint:gosec
	cmd := exec.Command(binary, fmt.Sprintf("file://%s", t.sourceA), fmt.Sprintf("file://%s", t.sourceB), "--output=json")

	outBytes, err := cmd.Output()
	if err != nil {
		// ignore errors when exit code is 1 - this is expected for tests that check for failures.
		if cmd.ProcessState.ExitCode() != 1 {
			return fmt.Errorf("failed to run command %q: %w | Output, if any: %s", cmd.String(), err, outBytes)
		}
	}

	if update {
		err := os.WriteFile(t.expected, outBytes, os.FileMode(0555))
		if err != nil {
			return fmt.Errorf("updating golden file %q: %w", t.expected, err)
		}

		return nil
	}

	var outJson runner.Results

	err = json.Unmarshal(outBytes, &outJson)
	if err != nil {
		return fmt.Errorf("unmarshalling output: %w", err)
	}

	expectedBytes, err := os.ReadFile(t.expected)
	if err != nil {
		return fmt.Errorf("reading contents of %q containing expected output: %w", t.expected, err)
	}

	var expectedJson runner.Results

	err = json.Unmarshal(expectedBytes, &expectedJson)
	if err != nil {
		return fmt.Errorf("unmarshalling expected output: %w", err)
	}

	if !compareSemantic(expectedJson, outJson) {
		return fmt.Errorf("%w : %s", errMismatchedOutput, cmp.Diff(expectedJson, outJson))
	}

	return nil
}

func compareSemantic(a, b runner.Results) bool {
	if !compareComparisonResultSemantic(a.CRDValidation, b.CRDValidation) {
		return false
	}

	if !compareVersionedComparisonResultsSemantic(a.SameVersionValidation, b.SameVersionValidation) {
		return false
	}

	if !compareVersionedComparisonResultsSemantic(a.ServedVersionValidation, a.ServedVersionValidation) {
		return false
	}

	return true
}

func compareComparisonResultSemantic(a, b []validations.ComparisonResult) bool {
	// do initial set comparison to make sure a and b have the same entries
	aSet := sets.New[string]()
	bSet := sets.New[string]()

	for _, res := range a {
		aSet.Insert(res.Name)
	}

	for _, res := range b {
		bSet.Insert(res.Name)
	}

	if !aSet.Equal(bSet) {
		return false
	}

	type result struct {
		errs     []string
		warnings []string
	}
	resultSet := map[string]result{}

	// build the expected set
	for _, res := range a {
		resultSet[res.Name] = result{
			errs:     res.Errors,
			warnings: res.Warnings,
		}
	}

	// do the comparison against actual
	for _, res := range b {
		expect := resultSet[res.Name]

		if !compareArraySemantic(expect.errs, res.Errors) {
			return false
		}

		if !compareArraySemantic(expect.warnings, res.Warnings) {
			return false
		}
	}

	return true
}

func compareVersionedComparisonResultsSemantic(a, b map[string]map[string][]validations.ComparisonResult) bool {
	aSet := sets.New[string]()
	bSet := sets.New[string]()

	for k := range a {
		aSet.Insert(k)
	}

	for k := range b {
		bSet.Insert(k)
	}

	if !aSet.Equal(bSet) {
		return false
	}

	for k, v := range b {
		expect := a[k]

		if !comparePropertyComparisonResultsSemantic(expect, v) {
			return false
		}
	}

	return true
}

func comparePropertyComparisonResultsSemantic(a, b map[string][]validations.ComparisonResult) bool {
	aSet := sets.New[string]()
	bSet := sets.New[string]()

	for k := range a {
		aSet.Insert(k)
	}

	for k := range b {
		bSet.Insert(k)
	}

	if !aSet.Equal(bSet) {
		return false
	}

	for k, v := range b {
		expect := a[k]

		if !compareComparisonResultSemantic(expect, v) {
			return false
		}
	}

	return true
}

func compareArraySemantic[T comparable](a, b []T) bool {
	aSet := sets.New(a...)
	bSet := sets.New(b...)

	return aSet.Equal(bSet)
}

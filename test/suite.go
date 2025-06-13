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
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/crdify/pkg/runner"
	"sigs.k8s.io/crdify/pkg/slices"
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
		return performUpdate(t, outBytes)
	}

	return performComparison(t, outBytes)
}

func performUpdate(t test, outBytes []byte) error {
	// First perform a comparison of the output and existing expected state.
	// If we receive an error here that means they are not semantically equal
	// and we should update the expected state file.
	// If we do not receive an error here, then we know that the output and the
	// existing expected state for this test is semantically equivalent
	// and we should not perform an update.
	// This helps to reduce churn in the expected state files because the output
	// of crdify is non-deterministic.
	err := performComparison(t, outBytes)
	if err == nil {
		return nil
	}

	err = os.WriteFile(t.expected, outBytes, os.FileMode(0555))
	if err != nil {
		return fmt.Errorf("updating golden file %q: %w", t.expected, err)
	}

	return nil
}

func performComparison(t test, outBytes []byte) error {
	var outJSON runner.Results

	err := json.Unmarshal(outBytes, &outJSON)
	if err != nil {
		return fmt.Errorf("unmarshalling output: %w", err)
	}

	expectedBytes, err := os.ReadFile(t.expected)
	if err != nil {
		return fmt.Errorf("reading contents of %q containing expected output: %w", t.expected, err)
	}

	var expectedJSON runner.Results

	err = json.Unmarshal(expectedBytes, &expectedJSON)
	if err != nil {
		return fmt.Errorf("unmarshalling expected output: %w", err)
	}

	if err := compareSemantic(expectedJSON, outJSON); err != nil {
		return fmt.Errorf("%w : %w", errMismatchedOutput, err)
	}

	return nil
}

func compareSemantic(a, b runner.Results) error {
	if err := compareComparisonResultSemantic(a.CRDValidation, b.CRDValidation); err != nil {
		return fmt.Errorf("comparing CRD validations: %w", err)
	}

	if err := compareVersionedComparisonResultsSemantic(a.SameVersionValidation, b.SameVersionValidation); err != nil {
		return fmt.Errorf("comparing same version validations: %w", err)
	}

	if err := compareVersionedComparisonResultsSemantic(a.ServedVersionValidation, a.ServedVersionValidation); err != nil {
		return fmt.Errorf("comparing served version validations: %w", err)
	}

	return nil
}

func compareComparisonResultSemantic(a, b []validations.ComparisonResult) error {
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
		return fmt.Errorf("expected comparison set %v does not match actual %v", aSet, bSet) //nolint:err113
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

		expectedErrsNormalized := normalizeStringSlice(expect.errs...)
		actualErrsNormalized := normalizeStringSlice(res.Errors...)

		if !compareArraySemantic(expectedErrsNormalized, actualErrsNormalized) {
			return fmt.Errorf("validation %q: expected error set %v does not match actual %v", res.Name, expect.errs, res.Errors) //nolint:err113
		}

		expectedWarnsNormalized := normalizeStringSlice(expect.warnings...)
		actualWarnsNormalized := normalizeStringSlice(res.Warnings...)

		if !compareArraySemantic(expectedWarnsNormalized, actualWarnsNormalized) {
			return fmt.Errorf("validation %q: expected warning set %v does not match actual %v", res.Name, expect.warnings, res.Warnings) //nolint:err113
		}
	}

	return nil
}

func normalizeStringSlice(in ...string) []string {
	return slices.Translate(normalizeWhitespace, in...)
}

// normalizeWhitespace normalizes a given string by splitting the
// string on all whitespace and rejoining them with a new line.
// This reduces flakiness in comparing things like diffs generated
// by the `unhandled` validation that is meant to generate
// human readable diffs and is nondeterministic in the whitespacing
// it outputs.
//
// An example of a normalized string:
// "A quick brown fox" becomes "A\nquick\nbrown\nfox".
func normalizeWhitespace(in string) string {
	fields := strings.Fields(in)
	return strings.Join(fields, "\n")
}

func compareVersionedComparisonResultsSemantic(a, b map[string]map[string][]validations.ComparisonResult) error {
	aSet := sets.New[string]()
	bSet := sets.New[string]()

	for k := range a {
		aSet.Insert(k)
	}

	for k := range b {
		bSet.Insert(k)
	}

	if !aSet.Equal(bSet) {
		return fmt.Errorf("expected version set %v does not match actual %v", aSet, bSet) //nolint:err113
	}

	for k, v := range b {
		expect := a[k]

		if err := comparePropertyComparisonResultsSemantic(expect, v); err != nil {
			return fmt.Errorf("comparing property validation results for version %q: %w", k, err)
		}
	}

	return nil
}

func comparePropertyComparisonResultsSemantic(a, b map[string][]validations.ComparisonResult) error {
	aSet := sets.New[string]()
	bSet := sets.New[string]()

	for k := range a {
		aSet.Insert(k)
	}

	for k := range b {
		bSet.Insert(k)
	}

	if !aSet.Equal(bSet) {
		return fmt.Errorf("expected property validation set %v does not match actual %v", aSet, bSet) //nolint:err113
	}

	for k, v := range b {
		expect := a[k]

		if err := compareComparisonResultSemantic(expect, v); err != nil {
			return fmt.Errorf("comparing results for property %q: %w", k, err)
		}
	}

	return nil
}

func compareArraySemantic[T comparable](a, b []T) bool {
	aSet := sets.New(a...)
	bSet := sets.New(b...)

	return aSet.Equal(bSet)
}

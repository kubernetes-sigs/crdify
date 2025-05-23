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
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/everettraven/crd-diff/pkg/runner"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/yaml"
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
		Short: "runs the crd-diff e2e test suite",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTests(binary, testDir, update)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	command.Flags().BoolVar(&update, "update", false, "whether or not to update the expected output golden files")
	command.Flags().StringVar(&binary, "binary", "./bin/crd-diff", "location of the crd-diff binary to execute")
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
	expectedName = "expected.yaml"
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
	cmd := exec.Command(binary, fmt.Sprintf("file://%s", t.sourceA), fmt.Sprintf("file://%s", t.sourceB), "--output=yaml")

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

	var outYaml runner.Results

	err = yaml.Unmarshal(outBytes, outYaml)
	if err != nil {
		return fmt.Errorf("unmarshalling output: %w", err)
	}

	expectedBytes, err := os.ReadFile(t.expected)
	if err != nil {
		return fmt.Errorf("reading contents of %q containing expected output: %w", t.expected, err)
	}

	var expectedYaml runner.Results

	err = yaml.Unmarshal(expectedBytes, expectedYaml)
	if err != nil {
		return fmt.Errorf("unmarshalling expected output: %w", err)
	}

	if diff := cmp.Diff(expectedYaml, outYaml); diff != "" {
		return fmt.Errorf("%w : %s", errMismatchedOutput, diff)
	}

	return nil
}

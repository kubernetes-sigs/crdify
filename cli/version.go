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

package cli

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/lipgloss"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// NewVersionCommand returns a new cobra.Command
// for printing the current version of crd-diff.
func NewVersionCommand() *cobra.Command {
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "installed version of crdify",
		Run: func(cmd *cobra.Command, args []string) {
			var out strings.Builder
			fig := figure.NewFigure("crdify", "rounded", true)
			out.WriteString(fig.String() + "\n\n")

			settingNameStyle := lipgloss.NewStyle().Bold(true)
			if bi, ok := debug.ReadBuildInfo(); ok {
				out.WriteString(fmt.Sprintf("%s: %s\n\n", settingNameStyle.Render("version"), bi.Main.Version))
				for _, setting := range bi.Settings {
					name := settingNameStyle.Render(setting.Key)
					out.WriteString(fmt.Sprintf("%s: %s\n", name, setting.Value))
				}
			} else {
				out.WriteString("unable to read build info")
			}
			fmt.Print(out.String())
		},
	}

	return versionCommand
}

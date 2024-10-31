package cli

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/lipgloss"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "installed version of crd-diff",
		Run: func(cmd *cobra.Command, args []string) {
			var out strings.Builder
			fig := figure.NewFigure("crd-diff", "rounded", true)
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

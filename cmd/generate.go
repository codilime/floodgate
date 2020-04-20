package cmd

import (
	"io"

	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/spf13/cobra"
)

// generateOptions store render command options
type generateOptions struct {
	outDir string
}

func NewGenerateCmd(out io.Writer) *cobra.Command {
	options := generateOptions{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate an example config file for Floodgate",
		Long:  "Generate an example config file for Floodgate. Note: --config flag is ignored.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(cmd, options)
		},
	}
	cmd.Flags().StringVarP(&options.outDir, "output-dir", "o", "", "output directory (required)")
	cmd.MarkFlagRequired("output-dir")
	return cmd
}

func runGenerate(cmd *cobra.Command, options generateOptions) error {
	resourceManager := rm.ResourceManager{}
	resourceManager.SaveStringToFile(options.outDir, exampleConfig)
	return nil
}

var exampleConfig = `endpoint: https://127.0.0.1/api/v1
insecure: true
auth:
  user: admin
  password: VRCm9L80yO3FHTKeVthtxknfGq1b10WqInKoBFqozphGcrGi
libraries:
  - sponnet
resources:
  - resources/applications
  - resources/pipelines
  - resources/pipelinetemplates`

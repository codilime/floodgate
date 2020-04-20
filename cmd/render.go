package cmd

import (
	"io"

	c "github.com/codilime/floodgate/config"
	fl "github.com/codilime/floodgate/parser/fileloader"
	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/spf13/cobra"
)

// renderOptions store render command options
type renderOptions struct {
	outDir string
}

// NewRenderCmd create new render command
func NewRenderCmd(out io.Writer) *cobra.Command {
	options := renderOptions{}
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render Jsonnet files",
		Long:  "Render Jsonnet files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRender(cmd, options)
		},
	}
	cmd.Flags().StringVarP(&options.outDir, "output-dir", "o", "", "output directory (required)")
	cmd.MarkFlagRequired("output-dir")
	return cmd
}

func runRender(cmd *cobra.Command, options renderOptions) error {
	flags := cmd.InheritedFlags()
	configPath, err := flags.GetString("config")
	if err != nil {
		return err
	}
	config, err := c.LoadConfig(configPath)
	if err != nil {
		return err
	}
	resourceManager := rm.ResourceManager{}
	resourceManager.Init(config, rm.FileLoaders(fl.NewJsonnetLoader(config.Libraries...)))
	if err := resourceManager.SaveResources(options.outDir); err != nil {
		return err
	}
	return nil
}

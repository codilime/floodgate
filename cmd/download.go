package cmd

import (
	c "github.com/codilime/floodgate/config"
	"github.com/spf13/cobra"
	"io"
)

// downloadOptions store download command options
type downloadOptions struct {
	projectName string
	outDir      string
}

// NewDownloadCmd create new download command
func NewDownloadCmd(out io.Writer) *cobra.Command {
	options := downloadOptions{}
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download spinnaker resources",
		Long:  "Download all resources related to provided project name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDownload(cmd, options)
		},
	}
	cmd.Flags().StringVarP(&options.projectName, "project-name", "p", "", "project name (required)")
	cmd.Flags().StringVarP(&options.outDir, "output-dir", "o", "", "output directory (required)")

	cmd.MarkFlagRequired("project-name")
	cmd.MarkFlagRequired("output-dir")
	return cmd
}

func runDownload(cmd *cobra.Command, options downloadOptions) error {
	flags := cmd.InheritedFlags()
	configPath, err := flags.GetString("config")
	if err != nil {
		return err
	}
	//config, err :=
	_, err = c.LoadConfig(configPath)
	if err != nil {
		return err
	}

	return nil
}

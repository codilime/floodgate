package cmd

import (
	c "github.com/codilime/floodgate/config"
	swagger "github.com/codilime/floodgate/gateapi"
	rm "github.com/codilime/floodgate/resourcemanager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
)

type executeOptions struct {
	webhookName string
}

func NewExecuteCmd(out io.Writer) *cobra.Command {
	options := executeOptions{}
	cmd := &cobra.Command{
		Use:   "execute [webhook name]",
		Short: "Trigger spinnaker webhook and track execution status",
		Long:  "Trigger spinnaker webhook and track execution status",
		RunE: func(cmd *cobra.Command, args []string) error {
			options.webhookName = args[0]

			return runExecute(cmd, options)
		},
		Args: cobra.MinimumNArgs(1),
	}

	return cmd
}

func runExecute(cmd *cobra.Command, options executeOptions) error {
	flags := cmd.InheritedFlags()
	configPath, err := flags.GetString("config")
	if err != nil {
		return err
	}
	config, err := c.LoadConfig(configPath)
	if err != nil {
		return err
	}
	config.Merge(cfg)

	resourceManager := &rm.ResourceManager{}
	if err := resourceManager.Init(config); err != nil {
		return err
	}
	client := resourceManager.GetClient()

	log.Infof("triggering '%s'", options.webhookName)

	var opts *swagger.WebhooksUsingPOSTOpts
	payload, resp, err := client.WebhookControllerApi.WebhooksUsingPOST(client.Context, options.webhookName, "webhook", opts)
	if err != nil {
		return err
	}

	if resp != nil {
		data := payload.(map[string]interface{})
		if data["eventProcessed"].(bool) {
			log.Info("event processed successfully")
			log.Infof("execution id is %s", data["eventId"])
		}
	}

	return nil
}

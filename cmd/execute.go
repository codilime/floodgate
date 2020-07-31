package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antihax/optional"
	c "github.com/codilime/floodgate/config"
	swagger "github.com/codilime/floodgate/gateapi"
	gc "github.com/codilime/floodgate/gateclient"
	rm "github.com/codilime/floodgate/resourcemanager"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type executeOptions struct {
	webhookName    string
	wait           bool
	waitTime       int
	parametersFile string
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

	cmd.Flags().BoolVarP(&options.wait, "wait", "w", false, "wait for pipeline execution to finish")
	cmd.Flags().IntVarP(&options.waitTime, "wait-time", "t", 30, "wait time between status check")
	cmd.Flags().StringVarP(&options.parametersFile, "parameters", "p", "", "path to parameters json file")

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

	var opts swagger.WebhooksUsingPOSTOpts

	if _, err := os.Stat(options.parametersFile); !os.IsNotExist(err) {
		jsonFile, err := os.Open(options.parametersFile)
		if err != nil {
			return err
		}
		defer jsonFile.Close()
		byteVal, _ := ioutil.ReadAll(jsonFile)
		var result map[string]interface{}
		err = json.Unmarshal(byteVal, &result)
		if err != nil {
			return err
		}

		opts.Event = optional.NewInterface(result)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "triggering '%s'\n", options.webhookName)

	payload, resp, err := client.WebhookControllerApi.WebhooksUsingPOST(client.Context, options.webhookName,
		"webhook", &opts)
	if err != nil {
		return err
	}

	if resp == nil {
		return errors.New("spinnaker API response is empty")
	}
	data := payload.(map[string]interface{})
	if !data["eventProcessed"].(bool) {
		return errors.New("event not processed")
	}

	fmt.Fprintf(cmd.OutOrStdout(), "event processed successfully\n")
	fmt.Fprintf(cmd.OutOrStdout(), "execution id is %s\n", data["eventId"])

	for options.wait {
		status, err := executionStatus(client, data["eventId"].(string))
		if err != nil {
			return err
		}

		switch status {
		case "NOT_STARTED":
			fmt.Fprintf(cmd.OutOrStdout(), "waiting for pipeline to start\n")
		case "RUNNING":
			fmt.Fprintf(cmd.OutOrStdout(), "pipeline is still running\n")
		case "SUCCEEDED":
			fmt.Fprintf(cmd.OutOrStdout(), "pipeline succeeded\n")
		default:
			fmt.Fprintf(cmd.OutOrStdout(), "something went wrong\n")
		}

		if status == "SUCCEEDED" {
			break
		}

		time.Sleep(time.Duration(options.waitTime) * time.Second)
	}

	return nil
}

func executionStatus(spinnakerAPI *gc.GateapiClient, eventId string) (string, error) {
	var opts swagger.SearchForPipelineExecutionsByTriggerUsingGETOpts
	opts.TriggerTypes = optional.NewString("webhook")
	opts.EventId = optional.NewString(eventId)

	payload, resp, err := spinnakerAPI.ExecutionsControllerApi.SearchForPipelineExecutionsByTriggerUsingGET(spinnakerAPI.Context,
		"*", &opts)
	if err != nil {
		return "", err
	}

	if resp != nil {
		data := payload[0].(map[string]interface{})

		if status, ok := data["status"].(string); ok {
			return status, nil
		}
	}

	return "", nil
}

package gateclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/codilime/floodgate/config"
	gateapi "github.com/codilime/floodgate/gateapi"
	"golang.org/x/net/context"
)

// GateapiClient is a Client for Gate API which has instance-specific information.
type GateapiClient struct {
	// Gate API Client
	*gateapi.APIClient

	// request context
	Context context.Context
}

// NewGateapiClient creates instance of Gate API wrapper based on Floodgate's Config
func NewGateapiClient(floodgateConfig *config.Config) *GateapiClient {
	var gateHTTPClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: floodgateConfig.Insecure},
		},
	}

	cfg := gateapi.NewConfiguration()
	cfg.BasePath = floodgateConfig.Endpoint
	cfg.HTTPClient = gateHTTPClient
	client := gateapi.NewAPIClient(cfg)

	auth := context.WithValue(context.Background(), gateapi.ContextBasicAuth, gateapi.BasicAuth{
		UserName: floodgateConfig.Auth.User,
		Password: floodgateConfig.Auth.Password,
	})

	return &GateapiClient{
		APIClient: client,
		Context:   auth,
	}
}

// WaitForSuccessfulTask function is waiting for task to finish
func (c GateapiClient) WaitForSuccessfulTask(checkTask map[string]interface{}, maxRetries int) error {
	taskID := strings.Split(checkTask["ref"].(string), "/")[2]

	task, resp, err := c.TaskControllerApi.GetTaskUsingGET1(c.Context, taskID)

	retry := 0
	for (checkTask == nil || !isTaskCompleted(task)) && (retry < maxRetries) {
		log.Print(retry)
		retry++
		time.Sleep(time.Duration(retry*retry) * time.Second)
		task, resp, err = c.TaskControllerApi.GetTaskUsingGET1(c.Context, taskID)
	}
	log.Print(task)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("encountered an error while polling for task: %s", taskID)
	}
	if !isTaskSuccessful(task) {
		return fmt.Errorf("encountered an error in task: %s", taskID)
	}

	return nil
}

func isTaskCompleted(task map[string]interface{}) bool {
	status, exists := task["status"]
	if !exists {
		return false
	}

	switch status.(string) {
	case
		"SUCCEEDED",
		"STOPPED",
		"SKIPPED",
		"TERMINAL",
		"FAILED_CONTINUE":
		return true
	}
	return false
}

func isTaskSuccessful(task map[string]interface{}) bool {
	status, exists := task["status"]
	if !exists {
		return false
	}

	switch status.(string) {
	case
		"SUCCEEDED",
		"STOPPED",
		"SKIPPED":
		return true
	}
	return false
}

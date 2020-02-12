package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"cl-gitlab.intra.codilime.com/spinops/floodgate/cmd/gateclient"
	"cl-gitlab.intra.codilime.com/spinops/floodgate/config"
	gateapi "cl-gitlab.intra.codilime.com/spinops/floodgate/gateapi"
)

func main() {
	floodgateConfig := &config.Config{}
	configFile, err := ioutil.ReadFile(config.DefaultLocation)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(configFile, &floodgateConfig)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Config: %v\n", floodgateConfig)
	}
	client := gateclient.NewGateapiClient(floodgateConfig)
	opts := gateapi.GetAllApplicationsUsingGETOpts{}
	app, _, err := client.ApplicationControllerApi.GetAllApplicationsUsingGET(client.Context, &opts)
	if err == nil {
		log.Print(app)
	} else {
		log.Fatal(err)
	}
}

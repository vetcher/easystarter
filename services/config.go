package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/vetcher/easystarter/util"
)

const SERVICES_JSON = "services.json"

var (
	configFile = flag.String("file", SERVICES_JSON, "Path to makefile of microservices inside Dir")
)

type ServiceConfig struct {
	Name   string   `json:"name"`   // Name of service
	Target string   `json:"target"` // Path to Makefile inside Dir
	Args   []string `json:"args"`   // Command line arguments for service
	Dir    string   `json:"dir"`    // Full path to directory with microservice
}

func validateConfig(config *ServiceConfig) (errArray []error) {
	if config.Name == "" {
		errArray = append(errArray, fmt.Errorf("field `name` is not provided for service"))
	}
	if config.Target == "" {
		errArray = append(errArray, fmt.Errorf("field `target` is not provided for %v service", config.Name))
	}
	return
}

func loadServicesConfiguration() ([]*ServiceConfig, error) {
	var configs []*ServiceConfig
	raw, err := ioutil.ReadFile(*configFile)
	if err != nil {
		// Lookup in home directory for configuration SERVICES_JSON file
		usr, err1 := user.Current()
		if err1 != nil {
			return nil, fmt.Errorf("can't get current user: %v", err1)
		}
		raw1, err1 := ioutil.ReadFile(filepath.Join(usr.HomeDir, SERVICES_JSON))
		if err1 != nil {
			return nil, err
		}
		raw = raw1
	}
	err = json.Unmarshal(raw, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func loadServices() error {
	configs, err := loadServicesConfiguration()
	if err != nil {
		return fmt.Errorf("can't load services configuration: %v", err)
	}
	var configurationErrors []error
	for _, config := range configs {
		errs := validateConfig(config)
		if len(errs) > 0 {
			configurationErrors = append(configurationErrors, errs...)
		} else {
			err := ServiceManager.RegisterService(config)
			if err != nil {
				configurationErrors = append(configurationErrors, err)
			}
		}
	}
	return util.ComposeErrors(configurationErrors)
}

func init() {

}

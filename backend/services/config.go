package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/vetcher/easystarter/backend/util"
)

const SERVICES_JSON = "services.json"

var (
	configFile = flag.String("config", SERVICES_JSON, "Path to configuration file with services parameters")
)

type ServiceConfig struct {
	Name    string   `json:"name"`    // Name of service
	Target  string   `json:"target"`  // Path to Makefile inside Dir/Name/
	Args    []string `json:"args"`    // Command line arguments for service
	Dir     string   `json:"dir"`     // Full path to directory with service
	Version string   `json:"version"` // Git tag, which can be represented as semantic versioning http://semver.org/
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

func loadAllConfigsFrom(fileName string) ([]*ServiceConfig, error) {
	var allConfigs []*ServiceConfig
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("can't read %s: %v", fileName, err)
	}
	err = json.Unmarshal(raw, &allConfigs)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s: %v", fileName, err)
	}
	return allConfigs, nil
}

func loadServicesConfiguration(allFlag bool, serviceNames []string) ([]*ServiceConfig, error) {
	allConfigs, err := loadAllConfigsFrom(*configFile)
	if err != nil {
		usr, err1 := user.Current()
		if err1 != nil {
			return nil, fmt.Errorf("can't get current user: %v", err1)
		}
		allConfigs, err1 = loadAllConfigsFrom(filepath.Join(usr.HomeDir, SERVICES_JSON))
		if err1 != nil {
			return nil, err
		}
	}
	if allFlag {
		return allConfigs, nil
	}
	var configs []*ServiceConfig
	for _, name := range serviceNames {
		if cfg := findConfigInConfigsByName(name, allConfigs); cfg != nil {
			configs = append(configs, cfg)
		}
	}
	return configs, nil
}

func findConfigInConfigsByName(name string, configs []*ServiceConfig) *ServiceConfig {
	for _, cfg := range configs {
		if cfg.Name == name {
			return cfg
		}
	}
	return nil
}

func loadServicesConfigurations(allFlag bool, svcNames []string) error {
	configs, err := loadServicesConfiguration(allFlag, svcNames)
	if err != nil {
		return fmt.Errorf("can't load services configuration: %v", err)
	}
	var configurationErrors []error
	for _, config := range configs {
		errs := validateConfig(config)
		if len(errs) > 0 {
			configurationErrors = append(configurationErrors, errs...)
		} else {
			err := serviceManager.RegisterService(config)
			if err != nil {
				configurationErrors = append(configurationErrors, err)
			}
		}
	}
	return util.ComposeErrors(configurationErrors)
}

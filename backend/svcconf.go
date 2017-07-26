package backend

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"path/filepath"

	"os/user"

	"github.com/kpango/glg"
)

const SERVICES_JSON = "services.json"

type ServiceConfig struct {
	Target string   `json:"target"` // Path to Makefile inside Dir
	Args   []string `json:"args"`   // Command line arguments for service
	Dir    string   `json:"dir"`    // Full path to directory with microservice
}

func loadServicesConfiguration() (error, bool) {
	configs := make(map[string]*ServiceConfig)
	raw, err := ioutil.ReadFile(SERVICES_JSON)
	if err != nil {
		// Lookup in home directory for configuration SERVICES_JSON file
		usr, err1 := user.Current()
		if err1 != nil {
			glg.Fatalf("Can't get current user: %v", err1)
		}
		raw1, err1 := ioutil.ReadFile(filepath.Join(usr.HomeDir, SERVICES_JSON))
		if err1 != nil {
			return err, false
		}
		raw = raw1
	}
	err = json.Unmarshal(raw, &configs)
	if err != nil {
		return err, false
	}
	allServices = make(map[string]*service)
	for svcName, svc := range configs {
		if svc.Target == "" {
			glg.Warnf("Field `target` is not provided for %v service", svcName)
		} else {
			allServices[svcName] = &service{
				Name:   svcName,
				Args:   svc.Args,
				Target: filepath.Clean(svc.Target),
				Dir:    filepath.Clean(svc.Dir),
			}
		}
	}
	return nil, false
}

func init() {
	err, fatal := loadServicesConfiguration()
	if err != nil {
		glg.Errorf("Can't load services: %v.", err)
		if fatal {
			os.Exit(FATAL_WHEN_LOAD_SERVICES_CONFIGURATION)
		}
	}
}

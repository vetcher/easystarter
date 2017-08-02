package services

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/gosuri/uitable"
	"github.com/vetcher/easystarter/util"
)

type ServiceRepository struct {
	services map[string]Service
	versions map[string]string
}

func NewRepository() *ServiceRepository {
	return &ServiceRepository{
		services: make(map[string]Service),
		versions: make(map[string]string),
	}
}

func (f *ServiceRepository) Names() []string {
	var keys []string
	for key := range f.services {
		keys = append(keys, key)
	}
	return keys
}

func (f *ServiceRepository) GetService(svcName string) (Service, error) {
	svc, exist := f.services[svcName]
	if !exist {
		return nil, fmt.Errorf("%s %s", CANT_FIND_SERVICE, svcName)
	} else {
		return svc, nil
	}
}

func (f *ServiceRepository) registerService(config *ServiceConfig) error {
	svc := &goService{
		info: ServiceInfo{
			Name:    config.Name,
			Dir:     config.Dir,
			Target:  config.Target,
			Args:    config.Args,
			Version: config.Version,
		},
	}
	if service, ok := f.services[config.Name]; ok && service.IsRunning() {
		return errors.New("service is running")
	} else {
		f.services[config.Name] = svc
	}
	return nil
}

func (f *ServiceRepository) RegisterService(config *ServiceConfig) error {
	return f.registerService(config)
}

func (f *ServiceRepository) String() string {
	if f.services != nil {
		table := uitable.New()
		table.Wrap = true
		for name, svc := range f.services {
			info := svc.Info()
			table.AddRow("Service:", fmt.Sprintf("%s:%s", name, f.versions[name]))
			table.AddRow("Dir:", util.StringOrEmpty(info.Dir))
			table.AddRow("Args:", info.Args)
			table.AddRow("Target:", util.StringOrEmpty(info.Target))
			table.AddRow("")
		}
		return table.String()
	} else {
		return ""
	}
}

func SwitchVersion(service Service) error {
	info := service.Info()
	cmd := exec.Command("git", "checkout", info.Version)
	cmd.Dir = info.Dir
	return cmd.Run()
}

func CallMakeDepGen(service Service) error {
	info := service.Info()
	cmd := exec.Command("make", "dep", "gen")
	cmd.Dir = info.Dir
	return cmd.Run()
}

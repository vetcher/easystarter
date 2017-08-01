package services

import (
	"fmt"
)

type ServiceRepository struct {
	services map[string]Service
}

func NewRepository() *ServiceRepository {
	return &ServiceRepository{
		services: make(map[string]Service),
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
	var err error = nil
	_, ok := f.services[config.Name]
	if ok {
		err = fmt.Errorf("overwrites %v", config.Name)
	}
	svc := &goService{
		SvcName: config.Name,
		Dir:     config.Dir,
		Target:  config.Target,
		Args:    config.Args,
	}
	f.services[config.Name] = svc
	return err
}

func (f *ServiceRepository) RegisterService(config *ServiceConfig) error {
	return f.registerService(config)
}

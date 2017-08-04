package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"os"

	"github.com/vetcher/easystarter/backend"
)

type LogicHandler func([]byte) interface{}

func NewHTTPWrapper(wrapped LogicHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			fmt.Fprint(writer, err)
			return
		}
		resp := wrapped(body)
		data, err := json.Marshal(resp)
		if err != nil {
			_, err := writer.Write([]byte(err.Error()))
			if err != nil {
				log.Println(err)
			}
			return
		}
		_, err = writer.Write(data)
		if err != nil {
			log.Println(err)
		}
	}
}

func HandleServicesInfo(input []byte) interface{} {
	var allFlag bool
	err := json.Unmarshal(input, allFlag)
	if err != nil {
		return err
	}
	return serviceManager.Info(allFlag)
}

func HandleReloadEnv(_ []byte) interface{} {
	return backend.ReloadEnvironment()
}

func HandleGetEnv(input []byte) interface{} {
	var allFlag bool
	err := json.Unmarshal(input, allFlag)
	if err != nil {
		return err
	}
	var env []string
	if allFlag {
		env, _ = backend.AllEnvironmentString()
	} else {
		env = backend.CurrentEnvironmentString()
	}
	return env
}

func HandleStartServices(input []byte) interface{} {
	var svcNames []string
	err := json.Unmarshal(input, svcNames)
	if err != nil {
		return err
	}
	return serviceManager.Start(svcNames...)
}

func HandleStopServices(input []byte) interface{} {
	var svcNames []string
	err := json.Unmarshal(input, svcNames)
	if err != nil {
		return err
	}
	return serviceManager.Stop(svcNames...)
}

func HandleKillServices(input []byte) interface{} {
	var svcNames []string
	err := json.Unmarshal(input, svcNames)
	if err != nil {
		return err
	}
	return serviceManager.Kill(svcNames...)
}

func HandleRestartServices(input []byte) interface{} {
	var svcNames []string
	err := json.Unmarshal(input, svcNames)
	if err != nil {
		return err
	}
	return serviceManager.Restart(svcNames...)
}

func HandleLoadAllServices(_ []byte) interface{} {
	return loadServicesConfigurations(true, nil)
}

func HandleLoadServices(input []byte) interface{} {
	var svcNames []string
	err := json.Unmarshal(input, svcNames)
	if err != nil {
		return err
	}
	return loadServicesConfigurations(false, svcNames)
}

func HandleAllServicesNames(_ []byte) interface{} {
	return serviceManager.AllServicesNames()
}

func StopServer(w http.ResponseWriter, _ *http.Request) {
	err := serviceManager.Stop(serviceManager.AllServicesNames()...)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	log.Println("Stop")
	os.Exit(0)
}

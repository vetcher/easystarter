package backend

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"os/user"
	"path/filepath"

	"github.com/kpango/glg"
	"gopkg.in/ini.v1"
)

const ENV_SECTION = ""

var (
	envFileName = flag.String("config", "env.ini", "File with configuration parameters")

	environment       *ini.File
	globalEnvFilePath string
)

func CurrentEnvironmentString() string {
	formattedStr := []string{""}
	for _, key := range environment.Section(ENV_SECTION).Keys() {
		formattedStr = append(formattedStr, fmt.Sprintf("%v=%v", key.Name(), key.Value()))
	}
	return strings.Join(formattedStr, "\n")
}

func AllEnvironmentString() string {
	cmd := exec.Command("printenv")
	s, err := cmd.Output()
	if err != nil {
		glg.Errorf("Can't print environment: %v", err)
		return ""
	}
	return string(s)
}

func SetupEnv() bool {
	var err error
	environment, err = LoadEnv()
	if err != nil {
		glg.Warnf("Can't load environment: %v", err)
		return false
	}
	err = ExpandEnv(environment)
	if err != nil {
		glg.Errorf("%v", err)
		return false
	}
	return true
}

func CreateEnvFile() {
	file, err := os.Create(globalEnvFilePath)
	if err != nil {
		glg.Warnf("Can't create file `%v`: %v.", *envFileName, err.Error())
	}
	glg.Printf("File `%v` was created.", glg.Green(globalEnvFilePath))
	glg.Print("You can add some environment variables.")
	file.Close()
}

func LoadEnv() (*ini.File, error) {
	file, err := os.Open(*envFileName)
	if err != nil {
		// Lookup in home directory for configuration global env file
		usr, err1 := user.Current()
		if err1 != nil {
			glg.Fatalf("Can't get current user: %v", err1)
		}
		globalEnvFilePath = filepath.Join(usr.HomeDir, *envFileName)
		file, err = os.Open(globalEnvFilePath)
		if err != nil {
			CreateEnvFile()
		} else {
			file.Close()
		}
		envConfig, err := ini.Load(globalEnvFilePath)
		if err != nil {
			return nil, err
		}
		return envConfig, nil
	} else {
		file.Close()
	}
	envConfig, err := ini.Load(*envFileName)
	if err != nil {
		return nil, err
	}
	return envConfig, nil
}

func ExpandEnv(cfg *ini.File) error {
	keys := cfg.Section("").Keys()
	var err error
	for _, key := range keys {
		env, ok := os.LookupEnv(key.Name())
		if ok {
			err = os.Setenv(key.Name(), fmt.Sprintf("%s:%s", env, os.ExpandEnv(key.Value())))
		} else {
			err = os.Setenv(key.Name(), os.ExpandEnv(key.Value()))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

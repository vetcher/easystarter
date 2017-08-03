package backend

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kpango/glg"
	"gopkg.in/ini.v1"
)

const (
	ENV_SECTION = ""
	ENV_INI     = "env.ini"
)

var (
	envFileName = flag.String("env", ENV_INI, "File with configuration parameters")

	environment   *ini.File
	globalEnvFile string
)

func init() {
	usr, err1 := user.Current()
	if err1 != nil {
		glg.Fatalf("Can't get current user: %v", err1)
	}
	globalEnvFile = filepath.Join(usr.HomeDir, ENV_INI)
}

func CurrentEnvironmentString() []string {
	var formattedStr []string
	for _, key := range environment.Section(ENV_SECTION).Keys() {
		formattedStr = append(formattedStr, fmt.Sprintf("%v=%v", key.Name(), key.Value()))
	}
	return formattedStr
}

func AllEnvironmentString() ([]string, error) {
	cmd := exec.Command("printenv")
	s, err := cmd.Output()
	if err != nil {
		return []string{}, fmt.Errorf("can't print environment: %v", err)
	}
	return strings.Split(string(s), "\n"), nil
}

func setupEnv() error {
	var err error
	environment, err = loadEnv()
	if err != nil {
		return fmt.Errorf("can't setup environment: %v", err)
	}
	err = expandEnv(environment)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func createEnvFile() error {
	file, err := os.Create(globalEnvFile)
	if err != nil {
		return fmt.Errorf("can't create file `%s`: %v.", globalEnvFile, err)
	}
	file.Close()
	return nil
}

func loadEnv() (*ini.File, error) {
	file, err := os.Open(*envFileName)
	if err != nil {
		// Lookup in home directory for configuration global env file
		file, err = os.Open(globalEnvFile)
		if err != nil {
			err := createEnvFile()
			if err != nil {
				return nil, err
			}
		} else {
			file.Close()
		}
		envConfig, err := ini.Load(globalEnvFile)
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

func expandEnv(cfg *ini.File) error {
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

func ReloadEnvironment() error {
	if err := setupEnv(); err != nil {
		return fmt.Errorf("can't load environment: %v", err)
	}
	return nil
}

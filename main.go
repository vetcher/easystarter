package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"strings"

	"github.com/vetcher/easystarter/services"
	"gopkg.in/ini.v1"
)

// TODO: commands: start, stop, kill, list services

const InputTip string = "InputTip"

func CreateEnvFile() {
	file, err := os.Create("env.ini")
	if err != nil {
		fmt.Printf("[!] Can't create file `env.ini` because of %v.\n", err.Error())
	}
	fmt.Printf("[?] File `env.ini` was created.\n[?] You can add some environment variables.\n")
	file.Close()
}

func LoadEnv() (*ini.File, error) {
	file, err := os.Open("env.ini")
	if err != nil {
		CreateEnvFile()
	}
	file.Close()
	envConfig, err := ini.Load("env.ini")
	if err != nil {
		return nil, err
	}
	return envConfig, nil
}

func ExpandEnv(cfg *ini.File) error {
	keys := cfg.Section("").Keys()
	var err error
	for _, key := range keys {
		err = os.Setenv(key.Name(), key.Value())
		if err != nil {
			return err
		}
	}
	return nil
}

var Environment *ini.File

func PrintEnvironment(env *ini.File) {
	if env != nil {
		for _, key := range env.Section("").Keys() {
			fmt.Printf("%v=%v\n", key.Name(), key.Value())
		}
	} else {
		cmd := exec.Command("printenv")
		s, err := cmd.Output()
		if err != nil {
			fmt.Printf("[!] Can't print environment because of: %v.\n", err)
			return
		}
		fmt.Println(string(s))
	}
}

func CommandManager(command string, args ...string) {
	switch command {
	case "start", "s":
		if len(args) > 0 {
			svcName := args[0]
			if svcName == "all" {
				services.StartAll()
			} else {
				svc := services.NewService(svcName, args[1:]...)
				if svc != nil {
					svc.Start()
				}
			}
		} else {
			fmt.Println("[?] Specify service name.")
		}
	case "stop", "kill", "k":
		if len(args) > 0 {
			services.StopService(args[0])
		} else {
			fmt.Println("[?] Specify service name.")
		}
	case "ps", "list":
		fmt.Println(services.ListServices())
	case "help", "h":
		fmt.Println(InputTip)
	case "reenv", "reload env":
		if SetupEnv() {
			fmt.Println("[_] Environment was reloaded.\n")
		}
	case "env":
		PrintEnvironment(Environment)
	case "env all", "all env":
		PrintEnvironment(nil)
	default:
		fmt.Printf("`%v` is wrong command, try to `help`.\n", command)
	}
}

func InfiniteLoop() {
	stdin := bufio.NewScanner(os.Stdin)
	for fmt.Print("->"); stdin.Scan(); fmt.Print("->") {
		text := stdin.Text()
		commands := strings.Split(text, " ")
		switch commands[0] {
		case "exit", "e", "ext", "out", "end", "break":
			return
		case "":
			continue
		default:
			CommandManager(commands[0], commands[1:]...)
		}
	}
}

func SetupEnv() bool {
	var err error
	Environment, err = LoadEnv()
	if err != nil {
		fmt.Printf("[!] Can't load environment because of %v.\n", err)
		return false
	}
	err = ExpandEnv(Environment)
	if err != nil {
		fmt.Printf("[!] There is an error: %v.\n", err)
		return false
	}
	return true
}

func init() {
	if !SetupEnv() {
		fmt.Println("I'm out, can't setup env")
		os.Exit(0)
	}
}

func main() {
	defer fmt.Println("I'm out")
	InfiniteLoop()
}

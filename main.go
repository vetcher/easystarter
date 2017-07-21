package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"strings"

	"flag"

	"github.com/vetcher/easystarter/printer"
	"github.com/vetcher/easystarter/services"
	"gopkg.in/ini.v1"
)

// TODO: commands: kill
// TODO: add customizable service logging
// TODO: specify service version

const (
	VERSION          string = "0.1"
	MainTip          string = "MainTip"
	WelcomeTip       string = "Easy Starter " + VERSION
	MKDIR_PERMISSION        = 0777
)

func init() {
	if !SetupEnv() {
		printer.Print("!", "I'm out, can't setup env")
		os.Exit(0)
	}
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", MKDIR_PERMISSION)
	}
}

func CreateEnvFile() {
	file, err := os.Create("env.ini")
	if err != nil {
		printer.Printf("!", "Can't create file `env.ini` because of %v.", err.Error())
	}
	printer.Print("?", "File `env.ini` was created.")
	printer.Print("?", "You can add some environment variables.")
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
			printer.Printf("I", "%v=%v", key.Name(), key.Value())
		}
	} else {
		cmd := exec.Command("printenv")
		s, err := cmd.Output()
		if err != nil {
			printer.Printf("!", "Can't print environment because %v.", err)
			return
		}
		fmt.Println(string(s))
	}
}

func CommandManager(command string, args ...string) {
	switch command {
	case "start", "s", "up":
		if len(args) > 0 {
			svcName := args[0]
			if svcName == "-all" {
				services.StartAll(args[1:]...)
			} else {
				svc := services.NewService(svcName, args[1:]...)
				if svc != nil {
					svc.Start()
				}
			}
		} else {
			printer.Print("?", "Specify service name.")
		}
	case "reload", "r":
		if len(args) > 0 {
			svcName := args[0]
			if svcName == "-all" {
				services.ReloadServices(args[1:]...)
			} else if svcName == "-env" {
				if SetupEnv() {
					printer.Print("I", "Environment was reloaded.")
				}
			} else {
				services.ReloadService(svcName)
			}
		} else {
			printer.Print("?", "Specify service name.")
		}
	case "stop", "kill", "k", "down":
		if len(args) > 0 {
			if args[0] == "-all" {
				services.StopAll()
			} else {
				services.StopService(args[0])
			}
		} else {
			printer.Print("?", "Specify service name.")
		}
	case "ps", "list", "ls":
		printer.PrintRaw(services.ListServices())
	case "help", "h":
		printer.PrintRaw(WelcomeTip)
		printer.PrintRaw(MainTip)
	case "reenv":
		if SetupEnv() {
			printer.Print("I", "Environment was reloaded.")
		}
	case "env", "vars":
		if len(args) > 0 {
			if args[0] == "-all" {
				PrintEnvironment(nil)
			} else if args[0] == "-reload" {
				if SetupEnv() {
					printer.Print("I", "Environment was reloaded.")
				}
			} else {
				PrintEnvironment(Environment)
			}
		}
	case "version", "v":
		printer.Print(VERSION)
	default:
		printer.Printf("?", "`%v` is wrong command, try to `help`.\n", command)
	}
}

func InfiniteLoop() {
	stdin := bufio.NewScanner(os.Stdin)
	for fmt.Print("->"); stdin.Scan(); fmt.Print("->") {
		text := stdin.Text()
		commands := strings.Split(text, " ")
		switch commands[0] {
		case "exit", "e", "ext", "out", "end", "break", "close", "quit":
			services.StopAll()
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
		printer.Printf("!", "Can't load environment because %v.", err)
		return false
	}
	err = ExpandEnv(Environment)
	if err != nil {
		printer.Printf("!", "There is an error: %v.", err)
		return false
	}
	return true
}

func main() {
	flag.Parse()
	defer printer.Print("!", "I'm out")
	printer.PrintRaw(WelcomeTip)
	InfiniteLoop()
}

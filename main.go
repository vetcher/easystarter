package main

import (
	"bufio"
	"os"

	"strings"

	"flag"

	"fmt"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
	"github.com/vetcher/easystarter/commands"
)

// TODO: commands: separate stop and kill
// TODO: specify service version
// TODO: add cleaning command

const (
	VERSION          = "0.2"
	WelcomeTip       = "Easy Starter " + VERSION
	MKDIR_PERMISSION = 0777

	START_CMD   = "start"
	STOP_CMD    = "stop"
	RESTART_CMD = "restart"
	PS_CMD      = "ps"
	ENV_CMD     = "env"
	EXIT_CMD    = "exit"
	VERSION_CMD = "ver"
)

func init() {
	if !backend.SetupEnv() {
		glg.Fatal("I'm out, can't setup env")
		os.Exit(0)
	}
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", MKDIR_PERMISSION)
	}
}

var allCommands = map[string]commands.Command{
	START_CMD:   commands.StartCommand{},
	STOP_CMD:    commands.StopCommand{},
	PS_CMD:      commands.PSCommand{},
	ENV_CMD:     commands.EnvCommand{},
	RESTART_CMD: commands.RestartCommand{},
}

func InfiniteLoop() {
	stdin := bufio.NewScanner(os.Stdin)
	for fmt.Print("->"); stdin.Scan(); fmt.Print("->") {
		text := stdin.Text()
		inputCommands := strings.Split(text, " ")
		command, ok := allCommands[inputCommands[0]]
		if ok {
			command.Exec(inputCommands[1:]...)
		} else {
			switch inputCommands[0] {
			case EXIT_CMD:
				backend.StopAllServices()
				return
			case "":
				continue
			case VERSION_CMD, "v":
				glg.Print(VERSION)
			default:
				glg.Printf("`%v` is wrong command, try to `help`.\n", inputCommands[0])
			}
		}
	}
}

func main() {
	flag.Parse()
	defer glg.Print("I'm out")
	glg.Print(WelcomeTip)
	InfiniteLoop()
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/services"
	"github.com/vetcher/easystarter/client/commands"
)

// TODO: specify service version
// TODO: add cleaning command

const (
	VERSION          = "0.6"
	WelcomeTip       = "Easy Starter " + VERSION
	MKDIR_PERMISSION = 0777

	CMD_START   = "start"
	CMD_STOP    = "stop"
	CMD_RESTART = "restart"
	CMD_PS      = "ps"
	CMD_ENV     = "env"
	CMD_EXIT    = "exit"
	CMD_VERSION = "version"
	CMD_KILL    = "kill"
	CMD_LOGS    = "logs"
	CMD_CFG     = "cfg"

	EXIT_CODE_SETUP_ENV_ERR = 1 + iota
	EXIT_CODE_INIT_LOGS_DIR_ERR
)

var (
	isStartOnStartup = flag.Bool("s", false, "Start all services after startup. Same as enter `start -all` after run program")
)

func init() {
	if err := <-services.ServeReloadEnv(); err != nil {
		glg.Fatalf("Exit, can't setup env: %v", err)
		os.Exit(EXIT_CODE_SETUP_ENV_ERR)
	}
	_, err := os.Stat("logs")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("logs", MKDIR_PERMISSION)
		} else {
			glg.Fatal(err)
			os.Exit(EXIT_CODE_INIT_LOGS_DIR_ERR)
		}
	}
}

func handleSignals(stopCommand commands.Command) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	for sig := range sigChan {
		glg.Print("Stop all services")
		_ = stopCommand.Validate("-all")
		err := stopCommand.Exec()
		if err != nil {
			glg.Error(err)
		}
		glg.Print("Terminate")
		os.Exit(int(sig.(syscall.Signal)))
	}
}

func main() {
	allCommands := map[string]commands.Command{
		CMD_START:   &commands.StartCommand{},
		CMD_STOP:    &commands.StopCommand{},
		CMD_PS:      &commands.PSCommand{},
		CMD_ENV:     &commands.EnvCommand{},
		CMD_RESTART: &commands.RestartCommand{},
		CMD_VERSION: &commands.VersionCommand{VERSION},
		CMD_EXIT:    &commands.ExitCommand{},
		"":          &commands.EmptyCommand{},
		CMD_KILL:    &commands.KillCommand{},
		CMD_LOGS:    &commands.LogsCommand{},
		CMD_CFG:     &commands.CfgCommand{},
	}
	flag.Parse()
	go handleSignals(allCommands[CMD_STOP])
	glg.Print(WelcomeTip)
	if *isStartOnStartup {
		glg.Print("start -all")
		_ = allCommands[CMD_START].Validate("-all")
		_ = allCommands[CMD_START].Exec()
	}
	stdin := bufio.NewScanner(os.Stdin)
	for fmt.Print("-> "); stdin.Scan(); fmt.Print("-> ") {
		text := stdin.Text()
		inputCommands := strings.Split(text, " ")
		command, ok := allCommands[inputCommands[0]]
		if ok {
			err := command.Validate(inputCommands[1:]...)
			if err != nil {
				glg.Errorf("Validation error: %v", err)
				continue
			}
			err = command.Exec()
			if err != nil {
				_ = allCommands[CMD_STOP].Validate("-all")
				err1 := allCommands[CMD_STOP].Exec()
				if err1 != nil {
					glg.Error(err1)
				}
				glg.Error(err)
				return
			}
		} else {
			glg.Printf("`%v` is wrong command", inputCommands[0])
		}
	}
}
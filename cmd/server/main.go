package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
	"github.com/vetcher/easystarter/commands"
)

// TODO: specify service version
// TODO: add cleaning command
// TODO: open logs

const (
	VERSION          = "0.2"
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

	EXIT_CODE_SETUP_ENV_ERR = 1 + iota
	EXIT_CODE_INIT_LOGS_DIR_ERR
)

var UNIX_SOCKET_FILE = filepath.Join(os.TempDir(), "easystarter.sock")

func init() {
	if !backend.SetupEnv() {
		glg.Fatal("I'm out, can't setup env")
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

var allCommands = map[string]commands.Command{
	CMD_START:   &commands.StartCommand{},
	CMD_STOP:    &commands.StopCommand{},
	CMD_PS:      &commands.PSCommand{},
	CMD_ENV:     &commands.EnvCommand{},
	CMD_RESTART: &commands.RestartCommand{},
	CMD_VERSION: &commands.VersionCommand{VERSION},
	CMD_EXIT:    &commands.ExitCommand{},
	"":          &commands.EmptyCommand{},
	CMD_KILL:    &commands.KillCommand{},
}

func commandHandler(c net.Conn, errc chan string, done chan struct{}) {
	buf := make([]byte, 512)
	n, err := c.Read(buf)
	if err != nil {
		errc <- "Internal server error, can't read"
		return
	}
	text := string(buf[:n])
	inputCommands := strings.Split(text, " ")
	command, ok := allCommands[inputCommands[0]]
	if ok {
		err := command.Validate(inputCommands[1:]...)
		if err != nil {
			resp := fmt.Sprintf("Validation error: %v", err)
			glg.Errorf(resp)
			errc <- resp
			return
		}
		err = command.Exec(inputCommands[1:]...)
		if err != nil {
			glg.Error(err)
			errc <- err.Error()
			return
		}
		errc <- ""
		return
	} else {
		formattedErr := fmt.Sprintf("`%v` is wrong command, try to `help`.", inputCommands[0])
		glg.Print(formattedErr)
		errc <- formattedErr
		return
	}
}

func responseWriter(c net.Conn, respc chan string, done chan struct{}) {
	resp, ok := <-respc
	println(resp)
	if !ok {
		glg.Log("Chan is closed")
		return
	}
	_, err := c.Write([]byte(resp))
	if err != nil {
		glg.Fatalf("Writing err: %v", err)
	}
}

func main() {
	flag.Parse()
	l, err := net.Listen("unix", UNIX_SOCKET_FILE)
	if err != nil {
		glg.Fatal(err)
		return
	}
	defer l.Close()
	defer os.Remove(UNIX_SOCKET_FILE)
	for {
		conn, err := l.Accept()
		if err != nil {
			glg.Errorf("Accept error: %v", err)
		} else {
			respChan := make(chan string)
			doneChan := make(chan struct{})
			go commandHandler(conn, respChan, doneChan)
			go responseWriter(conn, respChan, doneChan)
		}
	}
}

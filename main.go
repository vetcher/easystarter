package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/ini.v1"
)

const InputTip string = "InputTip"

type Service struct {
	Name string
}

var currentDir, _ = os.Getwd()

func CreateEnvFile() {
	file, err := os.Create("env.ini")
	if err != nil {
		fmt.Printf("[!] Can't create file `env.ini` because of %v.\n", err.Error())
	}
	fmt.Printf("[?] File `env.ini` was created.\n[?] You can add some enviroment variables.")
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

func NewService() {

}

func InfiniteLoop() {
	scaner := bufio.NewScanner(os.Stdin)
	expandEnv, err := LoadEnv()
	if err != nil {
		fmt.Printf("[!] Can't load environment because of %v.\n", err)
	}
	for fmt.Print("->"); scaner.Scan(); fmt.Print("->") {
		text := scaner.Text()
		switch text {
		case "start", "s":
			err := os.Setenv("$TESTENV", "omg")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(os.Getenv("TESTENV"))
				go fmt.Println("with go " + os.Getenv("TESTENV"))
				fmt.Println("With exec")
				cmd := exec.Command("printenv", "$TESTENV")
				s, err := cmd.Output()
				fmt.Println("TESTENV: "+string(s), err)
				go func() {
					fmt.Println("With exec and go")
					cmd := exec.Command("printenv", "$TESTENV")
					s, err := cmd.Output()
					fmt.Println("TESTENV: "+string(s), err)
				}()
			}

		case "help", "h":
			fmt.Println(InputTip)
		case "exit", "e", "ext", "out", "end", "break":
			fmt.Println("I'm out")
			os.Exit(0)
		default:
			fmt.Printf("`%v` is wrong command, try to `help`\n", text)
		}
	}
}

func main() {
	InfiniteLoop()
}
